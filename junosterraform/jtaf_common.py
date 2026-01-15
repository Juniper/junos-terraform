import xml.etree.ElementTree as ElementTree
import json
import sys
from typing import Any, Union

def get_xpaths(root:ElementTree.Element) -> dict[str, bool]:
    # defined a recursive function to walk the xml and populate result[]
    def recurse_children(node:ElementTree.Element, result:dict[str,bool] = {}, path:list[str] = [])-> dict[str,bool]:
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
 
def unique_xpaths(paths:dict[str, bool]) -> list[str]:
    path_dict = {}
    result = []

    for path in paths:
        if path == "groups/name":
            # Skip the groups/name path entirely
            continue
        elif path.startswith("groups"):
            # Remove the "groups/" prefix from other paths
            path = path[len("groups/"):]
        
        # Only add non-empty paths
        if path:
            path_dict[path] = path
    
    for key in path_dict.keys():
        result.append(path_dict[key])
    return result
 
def get_path(parent:list[Any]) -> str:
    path = ''
    for i in parent:
        if isinstance(i, dict) and "name" in i:
            path += "/" + i["name"]
    tmp_list = path.split("/")[3:]
    path = "/".join(tmp_list)
    return path
 
def check_path(paths:list[str], node:list[Any])-> bool:
    path = get_path(node)
    if path in paths or path =='':
        return True
    return False
 
def check_for_choice(elem:dict[str,Any]) -> list[Any]:
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

def check_for_enums(elem:dict[str,Any], node_parent:list[Any])-> list[Any]:
    cases =[]
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
        
def check_children(paths:list[str], elem:dict[str,Any], node_parent:list[Any], current_path:str)-> Union[list[Any], bool]:
    if isinstance(node_parent[-2], dict):
        if "children" in node_parent[-2].keys():
            if isinstance(elem, dict):
                elem_path = ''
                if current_path == '':
                    elem_path = elem['name']
                else:
                    elem_path = current_path + "/" + elem["name"]
 
            if elem_path in paths or elem_path == 'configuration':
                    # This code handles the logic for paths that are directly in the config
                    # and are not choices or enums
                return True
            else:
                # This code handles the logic for paths that aren't directly in the config but follow type 'choice' which leads to that path
                choices = check_for_choice(elem)
                enums = check_for_enums(elem, node_parent)
                if choices:
                    choice_list =[]
                    for choice in choices:
                        temp_path = current_path + "/" + choice[0]["name"]
                        if temp_path in paths:
                            choice_list.append(choice)
                    return choice_list   
                if enums:
                    enums_list = []
                    for enum in enums:
                        temp_path = current_path + "/" + enum["name"]
                        if temp_path in paths:
                            enums_list.append(enum)
                    return enums_list
        else:
            return True
    else:
        return True
    return False

def remove_tags_by_name(root:ElementTree.Element, tag_names:list[str]) -> None:
    # Recursively remove all nodes with matching tag names 
    for tag in tag_names:
        for elem in root.findall(f".//{tag}"):
            parent = find_parent(root, elem)
            if parent is not None:
                parent.remove(elem)

def find_parent(root:ElementTree.Element, child:ElementTree.Element)-> Union[ElementTree.Element, None]:
    # Find parent of a given element 
    for parent in root.iter():
        for elem in parent:
            if elem is child:
                return parent
    return None
 
def walk_schema(paths:list[str], node:Any, parent:list[Any] = [])-> Union[str, Any]:
    result = None
    emit_data = check_path(paths, parent)
    current_path = get_path(parent)
    if isinstance(node, dict):
        node['path'] = current_path
        result = {}
        parent.append(node)
        for k in node.keys():
            if emit_data:                
                # Node with empty children list is continued as type container or list just that no children element is missing'       
                tmp_result = walk_schema(paths, node[k], parent)
                if isinstance(tmp_result, list) and len(tmp_result) == 0:
                    continue
                else:
                    result[k] = walk_schema(paths, tmp_result, parent)

        parent.pop()
    elif isinstance(node, list):
        result = []
        parent.append(node)
        for elem in node:
            # UPDATE: This code now handles choice options --> vlan_tagging and vlan_id now included
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
            else:
                if isinstance(result_val, bool):
                    if result_val:                          
                        result.append(walk_schema(paths, elem, parent))

        parent.pop()
    else:
        result = node
    return result
 
# Method which starts the walk
def filter_json_using_xml(schema:str, xml: Union[ElementTree.Element,str]) -> str:
    if schema == "-":
        schema = json.load(sys.stdin)
    else:
        with open(schema) as f:
            schema = json.load(f)

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

    # find and remove any version node
    remove_tags_by_name(root, ["version", "versions", "model", "apply-groups"])

    # find the unique paths
    paths = unique_xpaths(get_xpaths(root))

    new_schema = walk_schema(paths, schema)
    return new_schema
 

def load_and_merge_xmls(xml_file_list:list[str]) -> ElementTree.Element:
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
