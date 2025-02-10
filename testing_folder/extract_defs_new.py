#!/usr/bin/env python3
import xml.etree.ElementTree as ElementTree
from copy import copy
import json
import argparse
import sys
from jinja2 import Template

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

def check_kids(paths, elem, node_parent, current_path):
    if isinstance(node_parent[-2], dict):
        if "kids" in node_parent[-2].keys():
            if isinstance(elem, dict):
                elem_path = ''
                if current_path == '':
                    elem_path = elem['name']
                else:
                    elem_path = current_path + "/" + elem["name"]
            if elem_path in paths or elem_path == 'configuration':
                return True
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
        result = {}
        parent.append(node)
        for k in node.keys():
            if emit_data:
                result[k] = walk_schema(paths, node[k], parent)
        parent.pop()
    elif isinstance(node, list):
        result = []
        parent.append(node)
        for elem in node:
            if check_kids(paths, elem, parent, current_path):
                result.append(walk_schema(paths, elem, parent))
        parent.pop()
    else:
        result = node
    return result

def filter_json_using_xml(schema, xml):
    with open(schema) as f:
        schema = json.loads(f.read())
    with open(xml) as f:
        xml_text = f.read()
    root = ElementTree.fromstring(f"<root>{xml_text}</root>")
    paths = unique_xpaths(get_xpaths(root))
    return walk_schema(paths, schema)

def main():
    # other arguments
    parser = argparse.ArgumentParser(exit_on_error=True)
    parser.add_argument('-j', '--json-schema', required=True, help='specify the json schema file')
    parser.add_argument('-x', '--xml-config', required=True, help='specify the xml config file')
    args = parser.parse_args()
    resources = filter_json_using_xml(args.json_schema, args.xml_config)
    # print(json.dumps(resources, indent=2))
    with open('go_template.j2') as jinja_tmpl:
        tmpl = Template(jinja_tmpl.read())
    print(tmpl.render(data=resources))
    
# run main()
if __name__ == "__main__":
    main()
