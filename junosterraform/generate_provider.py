#!/usr/bin/env python3
import xml.etree.ElementTree as ElementTree
from copy import copy
import json
import argparse
import sys
from jinja2 import Template
from junosterraform.resource_config_provider_go import render_template
import os
import shutil
import pkg_resources


def get_xpaths(root):
    # defined a recursive function to walk the xml and populate result[]
    def recurse_children(node, result = {}, path = []):
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
 
def unique_xpaths(paths):
    path_dict = {}
    result = []
    for path in paths:
        path_dict["/".join(path)] = path
    for key in path_dict.keys():
        result.append(path_dict[key])
    return result
 
def get_path(parent):
    path = ''
    for i in parent:
        if isinstance(i, dict) and "name" in i:
            path += "/" + i["name"]
    tmp_list = path.split("/")[3:]
    path = "/".join(tmp_list)
    return path
 
def check_path(paths, node):
    path = get_path(node)
    if path in paths or path =='':
        return True
    return False
 
def check_for_choice(elem):
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

def check_for_enums(elem, node_parent):
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
        
def check_children(paths, elem, node_parent, current_path):
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
 
def walk_schema(paths, node, parent = []):
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
 
# # Method which starts the walk
# def filter_json_using_xml(schema, xml):
#     with open(schema) as f:
#         schema = json.loads(f.read())
#     with open(xml) as f:
#         xml_text = f.read()
#     root = ElementTree.fromstring(f"<root>{xml_text}</root>")
#     paths = unique_xpaths(get_xpaths(root))
#     print(paths)
#     return walk_schema(paths, schema)

# Method which starts the walk
def filter_json_using_xml(schema, xml):
    if schema == "-":
        schema = json.load(sys.stdin)
    else:
        with open(schema) as f:
            schema = json.load(f)
    with open(xml) as f:
        xml_text = f.read()

    # Parse XML under temporary root
    root = ElementTree.fromstring(f"<root>{xml_text}</root>")

    # Try to find the <configuration> node under <rpc-reply>
    config_node = root.find(".//configuration")
    if config_node is not None:
        paths = unique_xpaths(get_xpaths(config_node))
    else:
        # Fall back to full tree if <configuration> not found
        paths = unique_xpaths(get_xpaths(root))

    return walk_schema(paths, schema)
 
# Main Method
def main():
    # other arguments
    parser = argparse.ArgumentParser(exit_on_error=True)
    parser.add_argument('-j', '--json-schema', required=True, help='specify the json schema file')
    parser.add_argument('-x', '--xml-config', required=True, help='specify the xml config file')
    parser.add_argument('-t', '--type', required=True, help='device type (i.e. vsrx, mx960, ex4200, etc)')
    args = parser.parse_args()

    # Step 1: Filter the schema using the config
    resources = filter_json_using_xml(args.json_schema, args.xml_config)

    # Step 2: Render the template into Go code
    output = render_template(data=resources)

    # Step 3: Prepare new output directory based on type
    package_dir = os.path.dirname(pkg_resources.resource_filename('junosterraform', '__init__.py'))
    base_dir = f"{package_dir}/terraform_provider"
    new_dir = f"terraform-provider-junos-{args.type}"

    # Step 4: copy the template go files into place
    shutil.copytree(base_dir, new_dir, dirs_exist_ok=True)

    # Step 5: Save rendered Go file into new directory
    output_path = os.path.join(new_dir, "resource_config_provider.go")
    with open(output_path, "w") as f:
        f.write(render_template(data=resources).lstrip())

    print(f"Plugin created in {os.path.relpath(output_path)}\n")

    # Step 6: Update provider.go with correct type
    provider_path = os.path.join(new_dir, "provider.go")
    if os.path.exists(provider_path):
        with open(provider_path, "r") as f:
            provider_content = f.read()

        updated_content = provider_content.replace(
            'resp.TypeName = "junos-<device-type>"',
            f'resp.TypeName = "junos-{args.type}"'
        )

        with open(provider_path, "w") as f:
            f.write(updated_content)

        print(f"Updated provider.go with type junos-{args.type}")
    else:
        print(f"Warning: provider.go not found in {new_dir}")
    
# run main()
if __name__ == "__main__":
    main()
