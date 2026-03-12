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

from copy import deepcopy
from typing import Any, Dict, Optional

try:
    from ansible.errors import AnsibleFilterError
except ImportError:
    class AnsibleFilterError(Exception):
        """Fallback error type when ansible is not installed."""


class FilterModule:
    """JTAF filters for Ansible."""

    def filters(self):
        """Return the filter-name to callable mapping for Ansible registration."""
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

            # Process nested structures first
            processed = {}
            for key, value in data.items():
                if key.startswith('_merge'):
                    # Keep meta-directive keys for now (will remove later)
                    processed[key] = value
                else:
                    processed[key] = self._process_directives(value)

            # Apply directive if present (mostly for documentation/transparency)
            # The actual merge behavior is handled at the Ansible combine level
            if directive:
                processed['_applied_directive'] = directive

            return processed

        if isinstance(data, list):
            return [self._process_directives(item) for item in data]

        return data


def _merge_replace(_base: Any, override: Any) -> Any:
    """Always return override value."""
    return override


def _merge_keep_parent(base: Any, _override: Any) -> Any:
    """Always keep parent/base value."""
    return base


def _merge_recursive(base: Any, override: Any) -> Any:
    """Deep-merge dicts, otherwise return override."""
    if isinstance(base, dict) and isinstance(override, dict):
        result = deepcopy(base)
        result.update(override)
        return result
    return override


def _merge_append(base: Any, override: Any) -> Any:
    """Append override to base with list-friendly coercion."""
    if isinstance(base, list) and isinstance(override, list):
        return base + override
    if isinstance(base, list):
        return base + [override]
    return [base, override]


def _merge_prepend(base: Any, override: Any) -> Any:
    """Prepend override to base with list-friendly coercion."""
    if isinstance(base, list) and isinstance(override, list):
        return override + base
    if isinstance(base, list):
        return [override] + base
    return [override, base]


def _merge_extend(base: Any, override: Any) -> Any:
    """Strict list-only append operation."""
    if not isinstance(base, list) or not isinstance(override, list):
        raise AnsibleFilterError(
            f"'extend' directive requires both values to be lists, "
            f"got {type(base).__name__} and {type(override).__name__}"
        )
    return base + override


DIRECTIVE_HANDLERS = {
    'replace': _merge_replace,
    'keep_parent': _merge_keep_parent,
    'merge_recursive': _merge_recursive,
    'append': _merge_append,
    'prepend': _merge_prepend,
    'extend': _merge_extend,
}


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
    normalized = directive or 'replace'
    handler = DIRECTIVE_HANDLERS.get(normalized)
    if handler is None:
        raise AnsibleFilterError(f"Unknown merge directive: {normalized}")
    return handler(base, override)
