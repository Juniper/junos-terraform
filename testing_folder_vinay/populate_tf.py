from lxml import etree
import sys
import re

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
    return re.sub(r'[-]', '_', tag)  # Replace hyphens with underscores

def parse_element(element):
    """
    Recursively parse an XML element into a nested structure for HCL conversion.
    - Empty tags become `True`
    - Tags with only text return scalar value
    - Tags with children return list[dict] to enforce HCL consistency
    """

    # Case: <vlan-tagging></vlan-tagging> → True
    if (not element.text or not element.text.strip()) and len(element) == 0:
        return True

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
        child_data = parse_element(child)
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
    resource_block = 'resource "junos-vsrx" "config_vsrx" {\n'
    resource_block += '  resource_name = "example_resource"\n'

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
        with open(xml_file, 'r') as file:
            xml_content = file.read()

        # Wrap the XML content with a single root element
        wrapped_xml = f"<root>{xml_content}</root>"

        # Parse the wrapped XML content
        tree = etree.fromstring(wrapped_xml)
        
        # Now we can parse the child elements under the root
        parsed_data = {}
        for elem in tree:
            tag = normalize_tag(elem.tag)
            parsed_data[tag] = parse_element(elem)

        # Generate the HCL output
        hcl_block = generate_hcl_resources(parsed_data)
        print(hcl_block)

    except Exception as e:
        print(f"Error parsing XML: {e}")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python xml_to_hcl.py <xml_file>")
        sys.exit(1)
    xml_file = sys.argv[1]
    parse_xml_to_hcl(xml_file)
