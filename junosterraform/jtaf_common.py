import gzip
import hashlib
import json
import re
import sys
import xml.etree.ElementTree as ElementTree
from typing import Any, Union


GZIP_MAGIC = b"\x1f\x8b"


def load_schema_json(path: str) -> Any:
    """Load a JTAF JSON schema from a file, or stdin when path is '-'.

    The schema may be plain JSON or gzip-compressed JSON; gzip is detected from
    the content (the 1f 8b magic bytes), not the file extension, matching what
    the Go provider's embed.go and cmd/compileschema do.
    """
    if path == "-":
        data = sys.stdin.buffer.read()
        name = "<stdin>"
    else:
        with open(path, "rb") as f:
            data = f.read()
        name = path
    if data[:2] == GZIP_MAGIC:
        try:
            data = gzip.decompress(data)
        except (OSError, EOFError) as e:
            raise ValueError(f"{name}: not a valid gzip file ({e})") from None
    try:
        return json.loads(data)
    except ValueError as e:
        raise ValueError(f"{name}: not valid JSON ({e})") from None


def write_schema_json_gz(resources: Any, path: str) -> None:
    """Write a JTAF JSON schema to path as compact, gzip-compressed JSON."""
    with gzip.open(path, "wt", encoding="utf-8") as f:
        json.dump(resources, f, separators=(",", ":"))


def get_xpaths(root: ElementTree.Element) -> dict[str, bool]:
    """Extract all XPath expressions from an XML tree.

    Args:
        root: Root element of the XML tree to traverse.

    Returns:
        Dictionary with XPath strings as keys and True as values.
    """
    # defined a recursive function to walk the xml and populate result[]
    def recurse_children(node: ElementTree.Element, result: dict[str, bool] = {}, path: list[str] = []) -> dict[str, bool]:
        """Recursively walk XML children and collect slash-separated paths."""
        for child in node:
            path.append(child.tag)
            recurse_children(child, result, path)
        result['/'.join(path)] = True
        # before we leave here remove the last path element
        if len(path) > 0:
            del path[-1]
        return result
    # run the search and return the result
    return recurse_children(root)


def unique_xpaths(paths: dict[str, bool]) -> list[str]:
    """Extract unique XPath strings while preserving input path semantics.

    Args:
        paths: Dictionary with XPath strings as keys.

    Returns:
        List of unique XPath strings.
    """
    path_dict = {}
    result = []

    for path in paths:
        # Only add non-empty paths
        if path:
            path_dict[path] = path

    for key in path_dict.keys():
        result.append(path_dict[key])
    return result


def get_path(parent: list[Any]) -> str:
    """Construct a path string from a parent hierarchy list.

    Args:
        parent: List of objects (typically dicts) with 'name' fields.

    Returns:
        Path string constructed from element names, skipping first 3 levels.
    """
    path = ''
    for i in parent:
        if isinstance(i, dict) and "name" in i:
            path += "/" + i["name"]
    tmp_list = path.split("/")[3:]
    path = "/".join(tmp_list)
    return path


def check_path(paths: list[str], node: list[Any]) -> bool:
    """Check if computed path from node exists in allowed paths list.

    Args:
        paths: List of allowed path strings.
        node: Hierarchy list to compute path from.

    Returns:
        True if computed path is in paths list or is empty, False otherwise.
    """
    path = get_path(node)
    if path in paths or path == '':
        return True
    return False


def check_for_choice(elem: dict[str, Any]) -> list[Any]:
    """Extract choice children from element if element type is 'choice'.

    Args:
        elem: Element dict to inspect for choice type.

    Returns:
        List of choice child dicts, empty list if not a choice or no children.
    """
    cases = []
    children = []
    if elem["type"] == "choice":
        # node of type choice but do not have any children e.g. vstp-flooding-option
        if "children" not in elem.keys():
            return children
        for k in elem["children"]:
            cases.append(k)
        for case in cases:
            if "children" in case.keys():
                children.append(case["children"])
    return children


def check_for_enums(elem: dict[str, Any], node_parent: list[Any]) -> list[Any]:
    """Transform enum entries for choice-ident leaf elements.

    Args:
        elem: Element dict with potential 'enums' field.
        node_parent: Parent hierarchy for context checking.

    Returns:
        List of transformed enum dicts (each with 'name', 'type', 'leaf-type'),
        empty list if not a choice-ident or has no enums.
    """
    cases = []
    children = []
    if elem["name"] == 'choice-ident' and elem["type"] == 'leaf':
        for k in elem["enums"]:
            # this fix for especially for community-name which is present in enums as well child outside enums
            if k['id'] not in node_parent[-2]['key']:
                cases.append(k)
        for case in cases:
            tmp_dict = {}
            tmp_dict['name'] = case['id']
            tmp_dict['type'] = 'leaf'
            tmp_dict['leaf-type'] = 'string'
            children.append(tmp_dict)
    return children


def _calc_elem_path(current_path: str, elem: dict[str, Any]) -> str:
    """Return full path string for an element based on current_path."""
    if not isinstance(elem, dict):
        return ''
    if current_path == '':
        return elem.get('name', '')
    return current_path + "/" + elem.get('name', '')


def _match_choices(paths: list[str], current_path: str, choices: list[Any]) -> list[Any]:
    """Return subset of choice options whose path is in the allowed paths."""
    matches = []
    for choice in choices:
        temp_path = current_path + "/" + choice[0].get("name", '')
        if temp_path in paths:
            matches.append(choice)
    return matches


def _match_enums(paths: list[str], current_path: str, enums: list[Any]) -> list[Any]:
    """Return subset of enum entries whose path is in the allowed paths."""
    matches = []
    for enum in enums:
        temp_path = current_path + "/" + enum.get("name", '')
        if temp_path in paths:
            matches.append(enum)
    return matches


def check_children(paths: list[str], elem: dict[str, Any], node_parent: list[Any],
                   current_path: str) -> Union[list[Any], bool]:  # noqa: C901
    """Determine whether an element should be included based on paths.

    Returns True if element matches directly, a list of matching
    choice/enums, or False otherwise.
    """
    if not isinstance(node_parent[-2], dict):
        return True

    parent = node_parent[-2]
    if "children" not in parent.keys():
        return True

    elem_path = _calc_elem_path(current_path, elem)
    if elem_path in paths or elem_path == 'configuration':
        # direct match in configuration paths
        return True

    # try choices and enums
    choices = check_for_choice(elem)
    if choices:
        matched = _match_choices(paths, current_path, choices)
        if matched:
            return matched

    enums = check_for_enums(elem, node_parent)
    if enums:
        matched = _match_enums(paths, current_path, enums)
        if matched:
            return matched

    return False


def flatten_choice_idents(node: Any) -> int:
    """Replace choice-ident/choice-value pairs with a leaf per choice case.

    A device writes a YANG choice as an element naming the case -- <add/>,
    <exact/> -- not as the choice-ident/choice-value pair the model declares.
    Trimming already produces that form, through check_for_enums, for the cases
    the XML uses. Doing the same here keeps a trimmed and an untrimmed schema
    describing the same node the same way, so Terraform files generated against
    one are valid against the other.

    The list key is left alone: it still names choice-ident and choice-value,
    and a key leaf that is absent is skipped when the entry is identified.

    Returns the number of nodes rewritten.
    """
    rewritten = 0

    if isinstance(node, list):
        for item in node:
            rewritten += flatten_choice_idents(item)
        return rewritten

    if not isinstance(node, dict):
        return rewritten

    children = node.get("children")
    if isinstance(children, list):
        ident = next(
            (c for c in children
             if isinstance(c, dict) and c.get("name") == "choice-ident"
             and c.get("type") == "leaf" and c.get("enums")),
            None,
        )
        if ident is not None:
            key = node.get("key", "")
            kept = [c for c in children
                    if not (isinstance(c, dict)
                            and c.get("name") in ("choice-ident", "choice-value"))]
            # An enum that is already a sibling leaf, as community-name and
            # martians' allow are, is that leaf rather than a case of the
            # choice; emitting it twice collides on one attribute name.
            taken = {c.get("name") for c in kept if isinstance(c, dict)}
            cases = [e for e in ident["enums"]
                     if e.get("id") not in key and e.get("id") not in taken]
            kept.extend({"name": case["id"], "type": "leaf", "leaf-type": "string"}
                        for case in cases)
            node["children"] = kept
            children = kept
            rewritten += 1

        for child in children:
            rewritten += flatten_choice_idents(child)

    for key, value in node.items():
        if key != "children" and isinstance(value, (dict, list)):
            rewritten += flatten_choice_idents(value)

    return rewritten


def yang_release_id(yang_files: list[str]) -> str:
    """Return a short identifier for the YANG modules a provider was built from.

    Derived from each module's name and newest revision date, so a provider
    built from one Junos release can be told from a provider built from
    another.

    Junos reuses a revision date across releases: every module in 18.2R3 here
    is revision 2019-01-01. Two releases that share revision dates therefore
    produce the same identifier, and a mismatch between them is not detected.
    Hashing the modules' contents would catch that, at the cost of changing
    the identifier whenever a file is touched at all.
    """
    module = re.compile(r"^\s*(?:sub)?module\s+([\w.-]+)", re.MULTILINE)
    revision = re.compile(r"^\s*revision\s+\"?(\d{4}-\d{2}-\d{2})\"?", re.MULTILINE)

    seen: dict[str, str] = {}
    for path in yang_files:
        try:
            text = open(path, encoding="utf-8", errors="replace").read()
        except OSError:
            continue
        name = module.search(text)
        if not name:
            continue
        revs = revision.findall(text)
        seen[name.group(1)] = max(revs) if revs else "none"

    if not seen:
        return "unknown"
    joined = ",".join(f"{name}@{rev}" for name, rev in sorted(seen.items()))
    return hashlib.sha256(joined.encode("utf-8")).hexdigest()[:12]


def schema_fingerprint(scope: str, release: str) -> str:
    """Return the identifier a provider and its Terraform files have to share."""
    return f"{scope}-{release}"


def remove_tags_by_name(root: ElementTree.Element, tag_names: list[str]) -> None:
    """Recursively remove all elements with specified tag names from tree.

    Args:
        root: Root element of XML tree.
        tag_names: List of tag names to remove.
    """
    # Recursively remove all nodes with matching tag names
    for tag in tag_names:
        for elem in root.findall(f".//{tag}"):
            parent = find_parent(root, elem)
            if parent is not None:
                parent.remove(elem)


def find_parent(root: ElementTree.Element,
                child: ElementTree.Element) -> Union[ElementTree.Element, None]:
    """Find and return the parent element of a given child element.

    Args:
        root: Root element to search from.
        child: Child element to find parent of.

    Returns:
        Parent element if found, None otherwise.
    """
    # Find parent of a given element
    for parent in root.iter():
        for elem in parent:
            if elem is child:
                return parent
    return None


def _walk_dict(paths: list[str], node: dict, parent: list[Any],
               emit_data: bool, current_path: str) -> dict:
    """Process dictionary node in schema walk."""
    node['path'] = current_path
    result: dict = {}
    parent.append(node)
    for k in node.keys():
        if not emit_data:
            break
        # Node with empty children list is continued as type container
        # or list just that no children element is missing
        tmp_result = walk_schema(paths, node[k], parent)
        if isinstance(tmp_result, list) and len(tmp_result) == 0:
            continue
        result[k] = walk_schema(paths, tmp_result, parent)
    parent.pop()
    return result


def _walk_list(paths: list[str], node: list, parent: list[Any],
               current_path: str) -> list:
    """Process list node in schema walk."""
    result: list = []
    parent.append(node)
    for elem in node:
        # UPDATE: This code now handles choice options -->
        # vlan_tagging and vlan_id now included
        result_val = check_children(paths, elem, parent, current_path)
        if isinstance(result_val, list):
            for item in result_val:
                if isinstance(item, list):
                    item[0]['path'] = current_path
                    result.append(item[0])
                elif isinstance(item, dict):
                    item['path'] = current_path
                    result.append(item)
        elif isinstance(result_val, dict):
            result_val['path'] = current_path
            result.append(result_val)
        elif isinstance(result_val, bool) and result_val:
            result.append(walk_schema(paths, elem, parent))
    parent.pop()
    return result


def walk_schema(paths: list[str], node: Any,
                parent: list[Any] = []) -> Union[str, Any]:
    """Recursive schema walker that dispatches by node type."""
    emit_data = check_path(paths, parent)
    current_path = get_path(parent)
    if isinstance(node, dict):
        return _walk_dict(paths, node, parent, emit_data, current_path)
    if isinstance(node, list):
        return _walk_list(paths, node, parent, current_path)
    return node


# Method which starts the walk
def filter_json_using_xml(schema: str,
                          xml: Union[ElementTree.Element, str],
                          groups: bool = False) -> str:
    """Filter JSON schema based on paths extracted from XML configuration.

    Reads schema from file or stdin, extracts configuration XPaths from XML,
    and returns filtered schema containing only relevant configuration sections.

    Args:
        schema: Path to JSON schema file or '-' to read from stdin.
        xml: Path to XML file or ElementTree element containing configuration.
        groups: Keep apply-groups, so that configuration inside <groups>
            contributes paths and the group body survives trimming.

    Returns:
        Filtered schema structure as dict (will be serialized by caller).
    """
    schema = load_schema_json(schema)

    # Check if xml is a single file or merged xml element
    if isinstance(xml, str):
        with open(xml) as f:
            xml_text = f.read()
        # Parse XML under temporary root
        root = ElementTree.fromstring(f"<root>{xml_text}</root>")
    else:
        root = xml

    # Try to find the <configuration> node under <rpc-reply>
    config_node = root.find(".//configuration")
    if config_node is not None:
        # set the config node to be root
        root = config_node

    # remove the top-level version node (the Junos release, not configuration); a leaf named
    # version elsewhere (system ntp server, protocols igmp interface, ...) is configuration
    for elem in root.findall("version"):
        root.remove(elem)
    drop = ["versions", "model"] if groups else ["versions", "model", "apply-groups"]
    remove_tags_by_name(root, drop)

    # find the unique paths
    paths = unique_xpaths(get_xpaths(root))

    new_schema = walk_schema(paths, schema)
    return new_schema


def load_and_merge_xmls(xml_file_list: list[str]) -> ElementTree.Element:
    """Load and merge configuration elements from multiple XML files.

    Args:
        xml_file_list: List of paths to XML files to merge.

    Returns:
        Root element containing merged configuration.

    Raises:
        ValueError: If any file does not contain a <configuration> element.
    """
    merged_config = ElementTree.Element("configuration")

    # Parse and find <configuration> in each file
    for path in xml_file_list:
        with open(path) as f:
            raw = f.read()
        root = ElementTree.fromstring(f"<root>{raw}</root>")
        config = root.find(".//configuration")
        if config is None:
            raise ValueError(f"No <configuration> found in {path}")
        for child in config:
            merged_config.append(child)

    return merged_config


# Using trimmed json, build mapping of path and path type
def build_type_map(node, parent_path=""):
    type_map = {}

    if node.get("type") in {"choice", "case"}:
        for child in node.get("children", []):
            type_map.update(build_type_map(child, parent_path))
        return type_map

    # Normalize node name
    node_name = normalize_tag(node["name"])
    path = f"{parent_path}/{node_name}" if parent_path else node_name

    # Strip leading 'root/configuration' to match XML traversal
    for prefix in ["root/configuration/", "configuration/"]:
        if path.startswith(prefix):
            path = path[len(prefix):]

    # Store type info
    node_info = {"type": node.get("type")}
    if "leaf-type" in node:
        node_info["leaf-type"] = node["leaf-type"]
    if "types" in node:
        node_info["types"] = node["types"]

    type_map[path] = node_info

    # Recurse into children
    for child in node.get("children", []):
        type_map.update(build_type_map(child, path))

    return type_map


def normalize_tag(tag: str) -> str:
    """Prepare XML tag to be used as YAML key by replacing '-' and '.' with '_'.

    Lower-cased to agree with the provider's SanitizeName: a few Junos names
    carry capitals, and an attribute the provider names in lower case has to
    be written that way in the Terraform files.
    """
    return tag.replace("-", "_").replace(".", "_").lower()
