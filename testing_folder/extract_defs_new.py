#!/usr/bin/env python3
import xml.etree.ElementTree as ElementTree
from copy import copy
import json
import argparse
import sys

def get_paths(root):
    # defined a recursive funciton to walk the xml and populate result[]
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

def unique_paths(paths):
    path_dict = {}
    result = []
    for path in paths:
        path_dict["/".join(path)] = path
    for key in path_dict.keys():
        result.append(path_dict[key])
    return result

def kid_by_name(node, name):
    if "kids" in node.keys():
        for kid in node["kids"]:
            if "name" in kid.keys() and kid["name"] == name:
                return kid
def get_base(schema):
    root = schema["root"]
    conf = kid_by_name(root, "configuration")
    if conf == None:
        conf = root["kids"][0]["configuration"]
    return conf
def get_def(schema, path):
    kid = get_base(schema)
    for elem in path:
        kid = kid_by_name(kid, elem)
        if kid == None:
            break
    return kid

def calculate_path(node, parent):
    pass

def get_path(parent):
    path = ''
    for i in parent:
        if isinstance(i, dict) and "name" in i:
            path += "/"+i["name"]
    tmp_list = path.split("/")[3:]
    path = "/".join(tmp_list)
    return path
def check_path(node):
    path = get_path(node)
    # print(path)
    if path in paths or path =='':
        return True
    # for match_path in paths:
    #     if path == match_path:
    #         return True
    return False

def walk_schema(node, indent = '', parent = [], parent_flag = False):
    indent += '  '
    flag = check_path(parent)
    current_path = get_path(parent)
    if not flag:
        return
    # print(flag)
    if isinstance(node, dict):
        dict_len = len(node.keys())
        local_indent = ''
        if len(parent) > 0 and isinstance(parent[-1], list):
            local_indent = indent
        if flag:
            print(local_indent+"{")
        parent.append(node)
        k_count = 1
        for k in node.keys():
            if flag:
                print(indent+"  "+f'"{k}": ', end='')
                
            walk_schema(node[k], indent+"  ", parent, flag)
            if k_count < dict_len:
                if flag:
                    print(",")
            k_count += 1
        parent.pop()
        if flag:
            print("\n"+indent+"}", end='')
    elif isinstance(node, list):
        list_len = len(node)
        if flag:
            print("[")
        parent.append(node)
        i = 1
        for elem in node:
            # print(get_path(elem))
            # print(elem)
            include_elem = False
            if isinstance(parent[-2], dict):
                if "kids" in parent[-2].keys():
                    if isinstance(elem, dict):
                        elem_path = ''
                        if current_path == '':
                            elem_path = elem['name']
                        else:
                            elem_path = current_path+"/"+elem["name"]
                        # print(elem_path, type(parent[-2]))
                    if elem_path in paths or elem_path == 'configuration':
                        include_elem = True
            else:
                include_elem = True
            # print(include_elem)
            if include_elem:
                walk_schema(elem, indent, parent, flag)
                if i < list_len:
                    if flag:
                        print(",")
            i += 1
        parent.pop()
        if flag:
            print("\n"+indent+"]", end='')
    else:
        if isinstance(node, str):
            if parent_flag and flag:
                print(f'"{node}"', end='')
        elif isinstance(node, bool):
            if parent_flag:
                if node:
                    print("true")
                else:
                    print("false")
        elif node == None:
            if parent_flag:
                print('null')
        else:
            if parent_flag:
                print(f'"{node}"')

def get_xml_config_resources(schema, xml):
    global paths
    with open(schema) as f:
        schema = json.loads(f.read())
    with open(xml) as f:
        xml_text = f.read()
    root = ElementTree.fromstring(f"<root>{xml_text}</root>")
    resources = []
    paths = unique_paths(get_paths(root))
    del paths[-1]
    walk_schema(schema)
    return resources

def main():
    # other arguments
    parser = argparse.ArgumentParser(exit_on_error=True)
    parser.add_argument('-j', '--json-schema', required=True, help='specify the json schema file')
    parser.add_argument('-x', '--xml-config', required=True, help='specify the xml config file')
    args = parser.parse_args()
    resources = get_xml_config_resources(args.json_schema, args.xml_config)
    print(json.dumps(resources, indent=2))
# run main()
if __name__ == "__main__":
    main()
