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
def debug_print(parent):
    path = ''
    for i in parent:
        # if isinstance(i, list):
        #     # pass
        #     print("/kids", file=sys.stderr, end='')
        if isinstance(i, dict) and "name" in i:
            path += "/"+i["name"]
            # print(, file=sys.stderr, end='')
    tmp_list = path.split("/")[3:]
    path = "/".join(tmp_list)
    print(path+"\n", file=sys.stderr)
def walk_schema(node, indent = '', parent = []):
    indent += '  '
    if isinstance(node, dict):
        dict_len = len(node.keys())
        local_indent = ''
        debug_print(parent)
        if len(parent) > 0 and isinstance(parent[-1], list):
            local_indent = indent
        print(local_indent+"{")
        parent.append(node)
        k_count = 1
        for k in node.keys():
            print(indent+"  "+f'"{k}": ', end='')
            walk_schema(node[k], indent+"  ", parent)
            if k_count < dict_len:
                print(",")
            k_count += 1
        parent.pop()
        print("\n"+indent+"}", end='')
    elif isinstance(node, list):
        list_len = len(node)
        print("[")
        parent.append(node)
        i = 1
        for elem in node:
            walk_schema(elem, indent, parent)
            if i < list_len:
                print(",")
            i += 1
        parent.pop()
        print("\n"+indent+"]", end='')
    else:
        if isinstance(node, str):
            print(f'"{node}"', end='')
        elif isinstance(node, bool):
            if node:
                print("true")
            else:
                print("false")
        elif node == None:
            print('null')
        else:
            print(f'"{node}"')
def get_xml_config_resources(schema, xml):
    with open(schema) as f:
        schema = json.loads(f.read())
    with open(xml) as f:
        xml_text = f.read()
    root = ElementTree.fromstring(f"<root>{xml_text}</root>")
    resources = []
    paths = get_paths(root)
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
