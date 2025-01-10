#!/usr/bin/ python3

import xml.etree.ElementTree as ElementTree
from copy import copy
import json
import argparse
import sys

def get_paths(root):
    # defined a recursive funciton to walk the xml and populate result[]
    def recurse_children(node, result = [], path = []):
        for child in node:
            path.append(child.tag)
            recurse_children(child, result, path)
        # if this is a terminal, capture the path    
        if (len(node)) == 0:
            # need to make a copy of the path array because arrays are refs
            result.append(copy(path))
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
            elif kid["type"] == "choice":
                for m in kid['kids']:
                    if m['kids'][0]['name'] == name:
                        return m['kids'][0]

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

def get_xml_config_resources(schema, xml):
    with open(schema) as f:
        schema = json.loads(f.read())
    with open(xml) as f:
        xml_text = f.read()
    root = ElementTree.fromstring(f"<root>{xml_text}</root>")
    resources = []
    for path in unique_paths(get_paths(root)):
        resource = { "path": path }
        d = get_def(schema, path)
        d_type = d['type']
        if d != None:
            resource["def"] = d
            resource['type']=d_type
            resources.append(resource)

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

