# CLI Tools Specification

Python command-line tools that drive the JTAF pipeline — from YANG models to generated Terraform/Ansible providers.

## Purpose

Define the behavior of the JTAF command-line tools, including schema file handling and generated schema artifacts.

## Overview

| Tool | Purpose |
|------|---------|
| `jtaf-provider` | JSON schema + XML configs → Go Terraform provider |
| `jtaf-yang2go` | End-to-end YANG → provider (wraps pyang + jtaf-provider) |
| `jtaf-xml2tf` | Junos XML config → HCL `.tf` files |
| `jtaf-xml2yaml` | Junos XML config → YAML representation |
| `jtaf-yang2ansible` | YANG + XML → Ansible role structure |
| `jtaf-ansible` | JSON schema + XML → Ansible artifacts |

## Requirements

### Requirement: Schema input supports plain and gzip JSON
All CLI tools that accept a JTAF JSON schema with `-j` SHALL accept plain JSON or gzip-compressed JSON, detecting gzip from the file content rather than its extension. `-j -` SHALL read either form from stdin. The tools SHALL use a shared loader, and invalid schema input SHALL fail with an error naming the input.

#### Scenario: Gzipped schema file
- **WHEN** `-j` names a gzip-compressed JSON file, regardless of extension
- **THEN** the tool decompresses and parses it as the equivalent plain JSON schema

#### Scenario: Gzipped schema on stdin
- **WHEN** `-j -` is given and stdin contains gzip-compressed JSON
- **THEN** the tool decompresses and parses it

#### Scenario: Invalid schema input
- **WHEN** `-j` names a file that is neither plain JSON nor gzip-compressed JSON
- **THEN** the tool exits with a non-zero status and an error naming the input

### Requirement: Provider and role generators write compressed schemas
`jtaf-provider` and `jtaf-ansible` SHALL write their trimmed schema as compact gzip-compressed JSON named
`trimmed_schema.json.gz` and SHALL NOT write `trimmed_schema.json`. This applies whether the schema is trimmed to XML
or embedded untrimmed.

#### Scenario: Generated schema handoff
- **WHEN** either generator completes
- **THEN** its output directory contains `trimmed_schema.json.gz`, and downstream tools can consume that file

#### Scenario: Schema help text
- **WHEN** a user checks the `-j` help text for any of the four tools
- **THEN** it names `trimmed_schema.json.gz` and states that plain JSON is also accepted

### Requirement: Generic provider generation excludes XML filtering
The `jtaf-yang2go` and `jtaf-provider` commands SHALL reject invocations that combine `--generic` with `-x` or `--xml-config`. Generic provider generation SHALL consume the unfiltered YANG schema; XML files remain valid inputs to downstream configuration conversion tools.

#### Scenario: End-to-end generic generation rejects XML filtering
- **WHEN** `jtaf-yang2go` is invoked with both `--generic` and `-x`
- **THEN** it exits with an argument error before invoking pyang or generating a provider

#### Scenario: Direct generic provider generation rejects XML filtering
- **WHEN** `jtaf-provider` is invoked with both `--generic` and `-x` or `--xml-config`
- **THEN** it exits with an argument error before loading or generating a provider

#### Scenario: Generic schema is converted using XML downstream
- **WHEN** a full generic provider schema is passed to `jtaf-xml2tf` or `jtaf-xml2yaml` with XML configuration inputs
- **THEN** the downstream tool uses those XML files for configuration conversion without applying provider-generation XML filtering

---

### Requirement: Groups option on the Terraform tools
`jtaf-provider` SHALL accept a `--groups` flag that keeps the `groups` subtree and the `apply-groups` leaf-list in the
generated schema. `jtaf-yang2go` SHALL accept the same flag and pass it to `jtaf-provider`. `jtaf-xml2tf` SHALL accept
a `--groups` flag that preserves the group hierarchy instead of flattening applied groups into the base configuration.
Omitting the flag SHALL preserve today's behaviour in every tool. The Ansible tools SHALL NOT accept the flag.

#### Scenario: Flag reaches the generator
- **WHEN** `jtaf-yang2go --groups` runs
- **THEN** it invokes `jtaf-provider` with `--groups`, and the resulting `trimmed_schema.json.gz` contains the `groups`
  subtree and the `apply-groups` leaf-list

#### Scenario: Flag is absent
- **WHEN** any of the three tools runs without `--groups`
- **THEN** its output is byte-identical to the output of the same invocation before this change

#### Scenario: Conversion preserves groups
- **WHEN** `jtaf-xml2tf --groups` converts an XML configuration containing `<groups>` and `<apply-groups>`
- **THEN** the `.tf` output contains a `groups` block and an `apply_groups` list rather than the flattened result

### Requirement: Schema scope flag
`jtaf-provider --generic` SHALL select the scope of the embedded schema rather than a provider implementation: with
`--generic` the untrimmed model is embedded, and without it the model is trimmed to the XML given with `-x`. The two
invocations SHALL produce the same provider source and SHALL differ only in the embedded schema.

#### Scenario: Same source, different schema
- **WHEN** a provider is generated with `--generic` and another is generated with `-x`
- **THEN** their Go source files are identical and only the embedded schema differs

#### Scenario: Trimming still requires XML
- **WHEN** `--generic` is combined with `-x`
- **THEN** the tool exits with a non-zero status and an error stating the two are mutually exclusive

## jtaf-provider

**Location:** `junosterraform/jtaf-provider`

### Behaviors

#### Input Handling

- **Given** `-j` is set to `'-'`, **When** the tool runs, **Then** JSON schema is read from stdin (enables piping from pyang)
- **Given** `-j` points to a file path, **When** the tool runs, **Then** JSON schema is loaded from that file
- **Given** `-x` specifies multiple XML files, **When** the tool runs, **Then** `load_and_merge_xmls()` merges them into a single XML tree before filtering
- **Given** `-t` is `vqfx-evpn-vxlan`, **When** the tool runs, **Then** output directory is `terraform-provider-junos-vqfx-evpn-vxlan/`

#### Processing Pipeline

1. **When** schema JSON and XML configs are loaded, **Then** `filter_json_using_xml()` retains only schema paths that appear in the XML config
2. **When** filtering completes, **Then** Jinja2 renders `resource_config_provider.go.j2` using the filtered resources
3. **When** the output directory already exists, **Then** it is deleted via `shutil.rmtree()` and recreated via `shutil.copytree()`
4. **When** the Go provider is created, **Then** `ensure_go_module_name()` sets `go.mod` to `module terraform-provider-junos-{type}`
5. **When** module name is set, **Then** `rewrite_import_prefixes()` replaces all `"terraform_provider/"` imports with `"terraform-provider-junos-{type}/"`
6. **When** code generation completes, **Then** `trimmed_schema.json.gz` (compact, gzipped JSON) is written to the output directory for the downstream tools; `-j` on all tools accepts plain or gzipped JSON, detected by content

#### Error Handling

- **Given** `-j` or `-x` not provided, **When** the tool runs, **Then** argparse exits with error code 2
- **Given** a Jinja2 template is missing, **When** rendering is attempted, **Then** `FileNotFoundError` is raised
- **Given** invalid JSON is provided, **When** template rendering starts, **Then** `TemplateError` is raised and execution stops

#### Expected Output

```
Plugin created in terraform-provider-junos-{type}/resource_config_provider.go
Updated provider.go with type junos-{type}
Updated config.go with type junos-{type}
Updated go.mod with type junos-{type}
```

---

## jtaf-yang2go

**Location:** `junosterraform/jtaf-yang2go`

### Behaviors

#### Input Classification

- **Given** `-p` argument is a directory, **When** arguments are parsed, **Then** it is added to pyang search paths via `-p`
- **Given** `-p` argument is a `.yang` file, **When** arguments are parsed, **Then** it is passed directly to pyang as a file argument
- **Given** `-p` argument is neither a directory nor `.yang` file, **When** arguments are parsed, **Then** a stderr warning is printed and the argument is ignored

#### Processing Pipeline

1. **When** arguments are classified, **Then** pyang is invoked: `pyang --plugindir {plugin_dir} -f jtaf -p {search_paths} {yang_files}`
2. **Given** pyang succeeds and produces output, **When** pyang returns, **Then** stdout is piped to `jtaf-provider -j - -x {files} -t {type}`
3. **Given** pyang fails (returncode != 0), **When** pyang returns, **Then** "pyang failed or produced no output, skipping jtaf-provider" is printed; no provider is generated
4. **Given** pyang produces no stdout, **When** pyang returns, **Then** `jtaf-provider` is skipped gracefully

---

## jtaf-xml2tf

**Location:** `junosterraform/jtaf-xml2tf`

### Behaviors

#### Element Parsing Rules

- **Given** an XML element has no text AND no children AND its tag is in `explicit_empty_tags`, **When** parsed, **Then** return `[{}]` (presence semantics)
- **Given** an XML element has no text AND no children AND tag is NOT in `explicit_empty_tags`, **When** parsed, **Then** return `""`
- **Given** an element has text but no children, **When** parsed AND `type_lookup[path]["type"] == "leaf-list"`, **Then** return `[value]` (wrapped in list)
- **Given** an element has text but no children, **When** parsed AND type is scalar, **Then** attempt `int(text)`, else return string
- **Given** an element has children, **When** parsed, **Then** build dict from children and return `[dict]`
- **Given** a child tag repeats under a parent, **When** parsed, **Then** values are collected into a list

#### Schema Validation

- **Given** an XML path has no corresponding entry in `type_lookup`, **When** parsing encounters it, **Then**:
  - Log to stderr: `"XML config '{path}' not defined in YANG. Removing it..."`
  - Return `None` (path excluded from HCL output)
  - Path added to `missing_xpaths` set (warning printed once per path)

#### Provider Block Generation

- **Given** `include_provider_block=True`, **When** `create_provider_block()` called, **Then** return full `required_providers` Terraform block
- **Given** `junos-{device_type}` already exists in block, **When** `edit_provider_block()` called, **Then** return block unchanged
- **Given** a new device type, **When** `edit_provider_block()` called, **Then** append to `required_providers` dict

---

## jtaf-xml2yaml

**Location:** `junosterraform/jtaf-xml2yaml`

### Behaviors

#### Type-Based Conversion

| Given Resource Type | When Element Processed | Then Output |
|---|---|---|
| `"list"` | Children present | Create list, append children as dicts |
| `"leaf-list"` | Text present | Create list, each occurrence is separate item |
| `"container"` | Children present | Create dict of children |
| `"container"` | No children | Return `""` |
| `"leaf"` | No text | Return `True` (presence) |
| `"leaf"` | Text present | Return string or int |

#### Tag Normalization

- **Given** a tag contains `-` or `.`, **When** `prepare_tag()` is called, **Then** characters are replaced with `_` (e.g., `product-name` → `product_name`)

#### Device Type Detection

- **Given** `system.product_name` contains "qfx", **When** `detect_device_type()` runs, **Then** return `"qfx"`
- **Given** `chassis.product_name` contains "mx", **When** detection runs, **Then** return `"mx"`
- Detection checks in order: qfx, srx, mx, ptx, acx, ex
- **Given** no product name matches, **When** detection completes, **Then** return `None`

---

## Core Module — jtaf_common.py

**Location:** `junosterraform/jtaf_common.py`

### Key Behaviors

#### filter_json_using_xml(schema, xml)

- **Given** a full YANG JSON schema and an XML config, **When** called, **Then** return JSON containing only paths present in the XML
- **Given** an XML element at path `/configuration/interfaces/interface`, **When** filtering, **Then** the corresponding JSON node and all its ancestors are retained
- **Given** a JSON node has no corresponding XML path, **When** filtering, **Then** it is pruned from output

#### load_and_merge_xmls(xml_file_list)

- **Given** a list of XML file paths, **When** called, **Then** merge all XML trees into a single `ElementTree.Element`
- **Given** two XML files define the same path, **When** merged, **Then** later file's values override earlier ones
- **Given** an XML file is malformed, **When** loading, **Then** `ElementTree.ParseError` is raised

#### get_xpaths(root)

- **Given** an XML root element, **When** called, **Then** recursively extract all XPaths as `dict[str, bool]`
- **Given** nested elements `configuration/interfaces/interface/name`, **When** extracted, **Then** all intermediate paths are included

---

## Testing

### How to Run

```bash
# All Python CLI tool tests
pytest junosterraform/tests/ -v

# Specific test file
pytest junosterraform/tests/test_jtaf_common.py -v

# Single test function
pytest junosterraform/tests/test_jtaf_common.py::test_get_xpaths -v

# With coverage
pytest junosterraform/tests/ --cov=junosterraform --cov-report=term-missing

# Workflow end-to-end test (requires YANG files)
pytest junosterraform/tests/test_workflow.py -v
```

### Test Coverage Map

| Test File | What It Validates |
|-----------|-------------------|
| `test_jtaf_common.py` | XPath extraction, path construction, schema filtering, merge behavior |
| `test_workflow.py` | End-to-end jtaf-yang2go: YANG parse → JSON → provider generation |
| `test_xml2tf_flatten.py` | XML→HCL conversion, apply-groups flattening, empty-tag handling |
| `test_hierarchical_groups.py` | Hierarchical group merging, inheritance, override behavior |
| `test_script_coverage.py` | CLI entry point invocability, argument parsing, error exits |

### Key Test Assertions

- `test_get_xpaths()`: Verifies all XML paths are extracted including intermediate containers
- `test_filter_json_using_xml()`: Schema nodes not in XML are pruned; matched nodes retained
- `test_workflow()`: Full pipeline produces valid Go provider directory with expected files
- `test_xml2tf_flatten()`: Verifies `apply-groups` are flattened before HCL generation
