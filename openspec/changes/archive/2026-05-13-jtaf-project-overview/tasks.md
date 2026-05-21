# Tasks

## 1. Generation Pipeline
- [x] 1.1 pyang plugin converts YANG AST to JSON schema tree
- [x] 1.2 jtaf-provider filters JSON using XML configs and renders Go code via Jinja2
- [x] 1.3 jtaf-yang2go wraps pyang + jtaf-provider into single command
- [x] 1.4 Post-generation: import rewriting, go.mod module naming, trimmed_schema.json emission

## 2. Terraform Provider Runtime
- [x] 2.1 Provider struct with Configure() establishing NETCONF client
- [x] 2.2 Config resource CRUD: Create (load-configuration), Read (get-configuration), Update (diff + edit-config), Delete (edit-config delete)
- [x] 2.3 File resource for dev/testing mode
- [x] 2.4 NETCONF client over SSH with RPC framing

## 3. Patch Engine (netconf_patch branch)
- [x] 3.1 LeafMapWithSchema: flatten XML to path→value map using schema-derived list keys
- [x] 3.2 ComputeDiff: set comparison producing add/delete/modify entries
- [x] 3.3 CreateDiffPatchWithSchema: two-pass XML generation (deletes first, then adds/modifies)
- [x] 3.4 AlignXMLOrderToReference: reorder patch XML to match device ordering
- [x] 3.5 ProcessSchema: load trimmed_schema.json into runtime index

## 4. Testing Infrastructure
- [x] 4.1 Mock NETCONF server with per-device state isolation
- [x] 4.2 Go unit tests: patch engine, provider core, NETCONF client
- [x] 4.3 Python unit tests: pyang plugin, mock server
- [x] 4.4 E2E: Terraform plan/apply against mock server

## 5. Ansible Output
- [x] 5.1 jtaf-ansible / jtaf-yang2ansible: generate roles + playbooks from YANG + XML
- [x] 5.2 jtaf-xml2yaml: convert XML configs to host_vars/group_vars with merge-safe repeated runs
- [x] 5.3 jtaf-xml2tf: generate .tf test files from XML configs

## 6. Documentation and Specs
- [x] 6.1 OpenSpec behavioral specs for all 7 components
- [x] 6.2 Project proposal capturing intent, scope, and approach
- [x] 6.3 Technical design document with architecture diagrams and component details
- [x] 6.4 OpenSpec README with setup, commands reference, and workflow examples
