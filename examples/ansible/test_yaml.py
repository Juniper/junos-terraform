#!/usr/bin/env python3
import yaml
from jinja2 import Environment, FileSystemLoader,Undefined
import os
import argparse


class SilentUndefined(Undefined):
    def _fail_with_undefined_error(self, *args, **kwargs):
        return ''


# Load YAML data
def load_yaml(file_path):
    with open(file_path, 'r') as file:
        return yaml.safe_load(file)

# Render Jinja2 template
def render_template(template_path, context):
    env = Environment(loader=FileSystemLoader(os.path.dirname(template_path)),undefined=SilentUndefined)
    template_name = os.path.basename(template_path)
    template = env.get_template(template_name)
    return template.render(context,undefined=SilentUndefined)

# Main function to test YAML and Jinja2
def main(yaml_file, template_file):
    # Load YAML data
    yaml_data = load_yaml(yaml_file)

    # Render the template with YAML data
    xml_output = render_template(template_file, yaml_data)

    # Save the output to an XML file
    with open('output.xml', 'w') as xml_file:
        xml_file.write(xml_output)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Process YAML and Jinja2 template files.')
    parser.add_argument('yaml_file', help='Path to the YAML file')
    parser.add_argument('template_file', help='Path to the Jinja2 template file')
    args = parser.parse_args()

    main(args.yaml_file, args.template_file)