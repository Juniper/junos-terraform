# CLI Tools

Python command-line tools that drive the JTAF pipeline — from YANG models to generated providers and Terraform/Ansible artifacts.

## Entry Points

| Tool | Script | Purpose |
|------|--------|---------|
| `jtaf-provider` | `junosterraform/jtaf-provider` | JSON schema + XML configs → Go Terraform provider |
| `jtaf-yang2go` | `junosterraform/jtaf-yang2go` | End-to-end YANG → provider pipeline (wraps pyang + jtaf-provider) |
| `jtaf-xml2tf` | `junosterraform/jtaf-xml2tf` | Junos XML config → HCL `.tf` files |
| `jtaf-xml2yaml` | `junosterraform/jtaf-xml2yaml` | Junos XML config → YAML representation |
| `jtaf-yang2ansible` | `junosterraform/jtaf-yang2ansible` | YANG + XML → Ansible role structure |
| `jtaf-ansible` | `junosterraform/jtaf-ansible` | JSON schema + XML → Ansible artifacts |
| `jtaf-pyang-plugindir` | `jtaf_pyang_plugin/jtaf-pyang-plugindir` | Prints plugin directory path for pyang `--plugindir` |

## Core Module — `junosterraform/jtaf_common.py`

Shared functions used by all CLI tools:

| Function | Purpose |
|----------|---------|
| `filter_json_using_xml(schema, xml)` | Filters YANG-derived JSON schema to only paths present in XML config |
| `load_and_merge_xmls(xml_file_list)` | Merges multiple XML config files into single tree |
| `get_xpaths(root)` | Extracts all XPath expressions from XML tree recursively |
| `check_path(paths, node)` | Validates if node's path exists in allowed paths list |
| `check_for_choice(elem)` | Extracts YANG choice children from element |
| `walk_schema(paths, node)` | Recursive walker filtering schema by allowed paths |
| `remove_tags_by_name(root, tag_names)` | Removes XML elements by tag name (in-place) |

## Pipeline Flow

```
jtaf-yang2go (all-in-one):
  1. Calls pyang with jtaf plugin → JSON schema
  2. Calls jtaf-provider internally

jtaf-provider:
  1. Reads JSON schema (-j) and XML configs (-x)
  2. Filters schema to only paths in XML via filter_json_using_xml()
  3. Renders Go templates (Jinja2) → terraform-provider-junos-{type}/

jtaf-xml2tf:
  1. Reads XML config and type lookup map
  2. Flattens XML hierarchy into HCL resource blocks
  3. Outputs .tf files
```

## Inputs & Outputs

| Tool | Inputs | Outputs |
|------|--------|---------|
| `jtaf-provider` | `-j` JSON schema, `-x` XML configs, `-t` device type | `terraform-provider-junos-{type}/` directory |
| `jtaf-yang2go` | `-p` YANG paths/files, `-x` XML configs, `-t` device type | `terraform-provider-junos-{type}/` directory |
| `jtaf-xml2tf` | `-config` XML file, `-type` device type | `.tf` files in `examples/terraform_files/` |
| `jtaf-xml2yaml` | XML file + schema resources | YAML output |

## Dependencies

- `pyang>=2.5` — YANG parser
- `jinja2>=3.1` — Template rendering
- `click>=8.0` — CLI argument parsing
- `lxml>=4.9` — XML parsing
- `pyaml>=25.7.0` — YAML serialization

## Tests

- `junosterraform/unit_tests/test_jtaf_common.py` — XPath extraction, nesting, filtering
- `junosterraform/unit_tests/test_workflow.py` — End-to-end jtaf-yang2go workflow
- `junosterraform/unit_tests/test_xml2tf_flatten.py` — XML→HCL conversion with apply-groups
