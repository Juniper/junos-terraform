#!/usr/bin/env python3
from lxml import etree
import sys
import re
import os

def convert_to_hcl(value, indent=2):
    """
    Converts Python dictionary or list into Terraform-compatible HCL syntax.
    """
    spaces = " " * indent
    if isinstance(value, dict):
        return "{\n" + "\n".join([f'{spaces}{k} = {convert_to_hcl(v, indent + 2)}' for k, v in value.items()]) + f'\n{spaces[:-2]}}}'
    elif isinstance(value, list):
        return "[\n" + ",\n".join([spaces + convert_to_hcl(item, indent + 2) for item in value]) + f'\n{spaces[:-2]}]'
    elif isinstance(value, bool):
        return "true" if value else "false"
    elif isinstance(value, (int, float)):
        return str(value)
    else:
        return f'"{value}"'

def normalize_tag(tag):
    return tag.replace('-', '_').replace('.', '_')  # Replace hyphens with underscores, Also replace dots


def parse_element(element, explicit_empty_tags):
    """
    Recursively parse an XML element into a nested structure for HCL conversion.
    - Empty tags become `True`
    - Tags with only text return scalar value
    - Tags with children return list[dict] to enforce HCL consistency
    """
    tag_name = etree.QName(element.tag).localname
    has_text = element.text is not None and element.text.strip() != ""
    has_children = len(element) > 0

    # Self-closing or empty tag
    if not has_text and not has_children:
        if tag_name in explicit_empty_tags:
            print(f"Explicit empty tag → [{{}}]: <{tag_name}>")
            return [{}]
        else:
            print(f"Self-closing tag → empty string: <{tag_name}>")
            return ""

    # # Case: <vlan-tagging></vlan-tagging> → True
    # if (not element.text or not element.text.strip()) and len(element) == 0:
    #     print(f"Empty tag: {element.tag}")
    #     return ""
    
    # Case: <name>ge-0/0/3</name> → "ge-0/0/3"
    if element.text and element.text.strip() and len(element) == 0:
        text_value = element.text.strip()
        try:
            return int(text_value)
        except ValueError:
            return text_value

    # Default case: has children
    parsed = {}
    for child in element:
        child_data = parse_element(child, explicit_empty_tags)
        tag = normalize_tag(child.tag)

        # If tag repeats, group into list
        if tag in parsed:
            if not isinstance(parsed[tag], list):
                parsed[tag] = [parsed[tag]]
            parsed[tag].append(child_data)
        else:
            parsed[tag] = child_data

    return [parsed]  # Always return a list of dict for elements with children


def generate_hcl_resources(parsed_data):
    """
    Generates a single Terraform-style resource block combining all top-level elements.
    """
    resource_block = 'resource "junos-<device-type>_Apply_Groups" "<device_hostname>" {\n'
    resource_block += '  resource_name = "JTAF_<device_hostname>"\n'

    for key, value in parsed_data.items():
        hcl_value = convert_to_hcl(value, indent=4)
        resource_block += f'  {key} = {hcl_value}\n'

    resource_block += "}\n"
    return resource_block


def parse_xml_to_hcl(xml_file):
    """
    Parses an XML file and converts it to Terraform-style HCL resources.
    """
    try:
        # Read the XML file content
        print(f"Reading XML file: {xml_file}")
        with open(xml_file, 'r') as file:
            xml_content = file.read()

        explicit_empty_tags = set()
        for match in re.finditer(r"<(\w[\w\-]*)>\s*</\1>", xml_content):
            explicit_empty_tags.add(match.group(1))

        # Wrap the XML content with a single root element
        wrapped_xml = f"<root>{xml_content}</root>"

        # Parse the wrapped XML content
        tree = etree.fromstring(wrapped_xml)
        
        # Now we can parse the child elements under the root
        parsed_data = {}
        for elem in tree:
            tag = normalize_tag(elem.tag)
            parsed_data[tag] = parse_element(elem, explicit_empty_tags)

        # Generate the HCL output
        hcl_block = generate_hcl_resources(parsed_data)
        return hcl_block

    except Exception as e:
        print(f"Error parsing XML: {e}")

def main():
    if len(sys.argv) != 2:
        print("Usage: populate_tf <config-file.xml>")
        sys.exit(1)

    xml_file = sys.argv[1]
    hcl_output = parse_xml_to_hcl(xml_file)

    if hcl_output:
        output_dir = "testbed"
        os.makedirs(output_dir, exist_ok=True)
        output_path = os.path.join(output_dir, "main.tf")
        with open(output_path, "w") as f:
            f.write(hcl_output)
        print(f"\nTerraform configuration written to {os.path.relpath(output_path)}\n")

if __name__ == "__main__":
    main()
