#!/usr/bin/env python3
"""
JTAF Ansible Filters for hierarchical YAML merge with merge directives.

Supports _merge_directive meta-instructions within YAML variables to control
how merging proceeds during playbook execution.

Merge directives:
  _merge_directive: "replace"           # Replace parent value (default)
  _merge_directive: "append"            # Append to parent list
  _merge_directive: "prepend"           # Prepend to parent list
  _merge_directive: "extend"            # Extend parent list
  _merge_directive: "merge_recursive"   # Deep merge dicts
  _merge_directive: "keep_parent"       # Use parent, ignore this override
"""

try:
    from ansible.errors import AnsibleFilterError
except ImportError:
    class AnsibleFilterError(Exception):
        """Fallback filter error when ansible is unavailable."""
from copy import deepcopy
from typing import Any, Dict, Optional


class FilterModule:
    """JTAF filters for Ansible."""

    def filters(self):
        return {
            'jtaf_apply_merge_directives': self.apply_merge_directives,
            'jtaf_extract_directive': self.extract_directive,
            'jtaf_remove_meta': self.remove_meta_keys,
        }

    @staticmethod
    def extract_directive(data: Any) -> Optional[str]:
        """Extract _merge_directive value from a dict if present."""
        if isinstance(data, dict):
            return data.get('_merge_directive')
        return None

    @staticmethod
    def remove_meta_keys(data: Any) -> Any:
        """Recursively remove all _merge_* keys from data structure."""
        if isinstance(data, dict):
            return {
                k: FilterModule.remove_meta_keys(v)
                for k, v in data.items()
                if not k.startswith('_merge')
            }
        if isinstance(data, list):
            return [FilterModule.remove_meta_keys(item) for item in data]
        return data

    def apply_merge_directives(self, jtaf_effective: Dict[str, Any]) -> Dict[str, Any]:
        """
        Process _merge_directive meta-instructions throughout the data structure.

        This filter walks through jtaf_effective and processes _merge_directive
        keys to determine how values should be handled.

        Example YAML with _merge_directive:
            routing:
              bgp:
                _merge_directive: replace    # Replace entire BGP block
                local_as: 65001

            interfaces:
              _merge_directive: append       # Append to parent interfaces list
              - name: eth0
                mtu: 1500
        """
        result = deepcopy(jtaf_effective)
        return self._process_directives(result)

    def _process_directives(self, data: Any) -> Any:
        """Recursively process merge directives in data structure."""
        if isinstance(data, dict):
            # Check if this dict has a merge directive
            directive = data.get('_merge_directive')
            merge_list_directives = data.get('_merge_list_directives', {})

            # Process nested structures first
            processed = {}
            for key, value in data.items():
                if key.startswith('_merge'):
                    # Keep meta-directive keys for now (will remove later)
                    processed[key] = value
                else:
                    processed[key] = self._process_directives(value)

            # Apply _merge_list_directives: de-duplicate lists by key
            for list_key, list_directive_config in merge_list_directives.items():
                if list_key in processed and isinstance(processed[list_key], list):
                    list_merge_directive = list_directive_config.get('_merge_directive')
                    if list_merge_directive == 'merge_by_key':
                        merge_key = list_directive_config.get('_merge_key', 'name')
                        processed[list_key] = self._merge_list_by_key(
                            processed[list_key], merge_key
                        )

            if directive:
                processed['_applied_directive'] = directive

            return processed

        if isinstance(data, list):
            return [self._process_directives(item) for item in data]

        return data

    @staticmethod
    def _merge_list_by_key(lst: list, key: str) -> list:
        """De-duplicate a list of dicts by a key field, merging duplicate entries.

        Later entries (higher priority) override earlier entries for scalar values.
        Nested dictionaries retain fields from both entries. Sub-lists at any depth
        are recursively merged by the same key (e.g. syslog file contents).
        Items lacking the key field are appended unchanged.
        """
        seen: Dict[Any, Any] = {}
        order = []

        for item in lst:
            if isinstance(item, dict) and key in item:
                k = item[key]
                if k in seen:
                    seen[k] = FilterModule._merge_dict_by_key(seen[k], item, key)
                else:
                    seen[k] = deepcopy(item)
                    order.append(k)
            else:
                # Items without the key field are appended as-is
                order.append(id(item))
                seen[id(item)] = item

        return [seen[k] for k in order]

    @staticmethod
    def _merge_dict_by_key(base: dict, override: dict, key: str) -> dict:
        merged = deepcopy(base)
        for name, value in override.items():
            previous = merged.get(name)
            if isinstance(previous, dict) and isinstance(value, dict):
                merged[name] = FilterModule._merge_dict_by_key(previous, value, key)
            elif isinstance(previous, list) and isinstance(value, list):
                merged[name] = FilterModule._merge_list_by_key(previous + value, key)
            else:
                merged[name] = deepcopy(value)
        return merged


def _dispatch_merge(base: Any, override: Any, directive: Optional[str] = None) -> Any:
    directive = directive or 'replace'

    handlers = {
        'replace': _merge_replace,
        'keep_parent': _merge_keep_parent,
        'merge_recursive': _merge_recursive,
        'append': _merge_append,
        'prepend': _merge_prepend,
        'extend': _merge_extend,
    }

    try:
        handler = handlers[directive]
    except KeyError as exc:
        raise AnsibleFilterError(f"Unknown merge directive: {directive}") from exc

    return handler(base, override)


def _merge_replace(base: Any, override: Any) -> Any:
    del base
    return override


def _merge_keep_parent(base: Any, override: Any) -> Any:
    del override
    return base


def _merge_recursive(base: Any, override: Any) -> Any:
    if isinstance(base, dict) and isinstance(override, dict):
        result = deepcopy(base)
        result.update(override)
        return result
    return override


def _merge_append(base: Any, override: Any) -> Any:
    if isinstance(base, list) and isinstance(override, list):
        return base + override
    if isinstance(base, list):
        return base + [override]
    return [base, override]


def _merge_prepend(base: Any, override: Any) -> Any:
    if isinstance(base, list) and isinstance(override, list):
        return override + base
    if isinstance(base, list):
        return [override] + base
    return [override, base]


def _merge_extend(base: Any, override: Any) -> Any:
    if not isinstance(base, list) or not isinstance(override, list):
        raise AnsibleFilterError(
            f"'extend' directive requires both values to be lists, "
            f"got {type(base).__name__} and {type(override).__name__}"
        )
    return base + override


def jtaf_merge_with_directive(base: Any, override: Any, directive: Optional[str] = None) -> Any:
    """
    Merge two values according to a merge directive.

    Args:
        base: Base/parent value
        override: Override/child value
        directive: Merge directive ('replace', 'append', 'prepend', 'extend', 'merge_recursive')

    Returns:
        Merged value
    """
    return _dispatch_merge(base, override, directive)
