# Proposal

## Why

The trimmed schema that the generators hand to the downstream tools is written in two different forms: `jtaf-provider --generic` writes compact, gzipped `trimmed_schema.json.gz`, while `jtaf-provider` (Jinja2) and `jtaf-ansible` write indented `trimmed_schema.json`, and every consumer (`jtaf-xml2tf`, `jtaf-xml2yaml`, `filter_json_using_xml`) only reads the plain form. The documented generic workflow (`build-generic.sh` then `convert.sh`) therefore fails on a file that was never written, and a full Junos model as indented JSON is hundreds of MB. One artifact, in one form, shared by the Terraform and Ansible pipelines, removes the seam.

## What Changes

- **BREAKING**: every generator writes the schema as `trimmed_schema.json.gz` (compact JSON, gzipped); no generator writes a plain `trimmed_schema.json`.
  - Terraform: `jtaf-provider` in both the Jinja2 and `--generic` paths (the generic path already does).
  - Ansible: `jtaf-ansible` (and so `jtaf-yang2ansible`, which drives it).
- Every consumer reads the schema through one shared loader in `jtaf_common` that detects gzip by content (the `1f 8b` magic bytes), not by file extension, and accepts `-` for stdin. Plain JSON is still accepted, so existing provider and role directories keep working.
  - Terraform: `jtaf-xml2tf -j`.
  - Ansible: `jtaf-xml2yaml -j`, including override mode.
  - Both: `jtaf-provider -j` / `jtaf-ansible -j` (the pyang JSON, via `filter_json_using_xml`).
- The Go side needs no change: `generic/embed.go` and `cmd/compileschema` already sniff gzip. The Jinja2 provider inlines the schema in Go source and never reads the file.
- CI workflows, example scripts, READMEs and the demo point `-j` at `trimmed_schema.json.gz`, and both workflows (Go provider and Ansible provider) exercise the gzipped handoff.
- Tests: the loader is unit-tested with plain, gzipped and stdin input; the end-to-end workflow tests for `jtaf-yang2go` and `jtaf-yang2ansible` assert the gzipped file and that the downstream tool consumes it.

Out of scope: the attribute-name case difference between `generic.SanitizeName` (lowercases) and `normalize_tag` (does not), which is a separate change.

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `code-generation`: the schema emitted with a generated Terraform provider is `trimmed_schema.json.gz` (compact, gzipped) in both generation paths, and the output structure lists it.
- `cli-tools`: `jtaf-provider` and `jtaf-ansible` write `trimmed_schema.json.gz`; `jtaf-xml2tf`, `jtaf-xml2yaml`, `jtaf-provider` and `jtaf-ansible` read a schema given with `-j` as plain or gzipped JSON, detected by content, through a shared `jtaf_common` loader.
- `override-mode`: the unmodeled-path check reads the schema in the same form as the rest of the Ansible pipeline (wording only: `trimmed_schema.json` → the trimmed schema).
- `testing`: the workflow tests cover the gzipped handoff from generator to consumer for both pipelines.

## Impact

- **Python CLI** (`junosterraform/`): `jtaf-provider` (Jinja2 path, step 7), `jtaf-ansible` (step 5) write gzipped; `jtaf-xml2tf`, `jtaf-xml2yaml` (`load_schema`), `jtaf_common.filter_json_using_xml` read through the new `jtaf_common` loader. `-j` help text on all four tools.
- **Go**: none.
- **CI**: `.github/workflows/go-terraform-provider.yml` and `.github/workflows/ansible-provider.yml` pass `trimmed_schema.json.gz` to `jtaf-xml2tf` / `jtaf-xml2yaml`.
- **Scripts and docs**: `examples/providers/convert.sh`, `build-generic-full.sh`, `examples/ansible/convert-ansible.sh`, `demo-override-mode.sh`, `README-terraform.md`, `README-ansible.md`, `grouping-hosts-file.md`, `examples/DEMO-GENERIC-PROVIDER.md`, `changelog.md`.
- **Tests**: `junosterraform/tests/test_workflow.py`, `test_script_coverage.py`, `test_hierarchical_groups.py`, `test_ansible_override.py` reference the plain filename; new loader tests in `test_jtaf_common.py`.
- **Users**: regenerated providers and roles contain `trimmed_schema.json.gz`, and the `-j` argument to `jtaf-xml2tf` / `jtaf-xml2yaml` changes accordingly. Directories generated before this change still load, as the loader accepts plain JSON.
