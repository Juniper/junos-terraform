# Spec Delta

## ADDED Requirements

### Requirement: Schema handoff is tested for both pipelines
The test suite SHALL verify that the schema written by each generator is consumed by its downstream tool, in the gzipped form, for the Terraform pipeline (`jtaf-yang2go` → `jtaf-xml2tf`) and the Ansible pipeline (`jtaf-yang2ansible` → `jtaf-xml2yaml`), and SHALL verify the shared schema loader with plain, gzipped and stdin input.

#### Scenario: Loader unit tests
- **WHEN** the Python unit tests run
- **THEN** the shared loader is tested with a plain JSON file, a gzipped JSON file (including one without a `.gz` extension), gzipped JSON on stdin, and a file that is neither, which fails with an error naming the file

#### Scenario: Terraform workflow test
- **WHEN** the end-to-end workflow test for `jtaf-yang2go` runs
- **THEN** it asserts that `trimmed_schema.json.gz` exists and no `trimmed_schema.json` exists in the provider directory, and that `jtaf-xml2tf` given that file produces `.tf` output

#### Scenario: Ansible workflow test
- **WHEN** the end-to-end workflow test for `jtaf-yang2ansible` runs
- **THEN** it asserts that `trimmed_schema.json.gz` exists and no `trimmed_schema.json` exists in the role directory, and that `jtaf-xml2yaml` given that file produces host and group variables

#### Scenario: CI
- **WHEN** the Go provider and Ansible provider GitHub workflows run
- **THEN** each passes the generator's `trimmed_schema.json.gz` to its downstream tool and the run succeeds
