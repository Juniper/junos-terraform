# Spec Delta

## ADDED Requirements

### Requirement: Schema input is read by content, not extension
Every CLI tool that takes a JTAF JSON schema with `-j` (`jtaf-provider`, `jtaf-ansible`, `jtaf-xml2tf`, `jtaf-xml2yaml`) SHALL accept the schema either as plain JSON or as gzip-compressed JSON, detecting gzip from the file content (the `1f 8b` magic bytes), regardless of the file's extension. `-j -` SHALL read the schema from stdin with the same detection. All four tools SHALL share one loader so the accepted forms cannot diverge.

#### Scenario: Gzipped schema file
- **WHEN** `-j` names a gzip-compressed JSON file, whatever its extension
- **THEN** the tool decompresses and parses it, and behaves as it would with the equivalent plain file

#### Scenario: Plain schema file from an earlier generator
- **WHEN** `-j` names a plain `trimmed_schema.json` written before this change
- **THEN** the tool parses it unchanged, so existing provider and role directories keep working

#### Scenario: Gzipped schema on stdin
- **WHEN** `-j -` is given and stdin carries gzip-compressed JSON
- **THEN** the tool decompresses and parses it

#### Scenario: Not JSON
- **WHEN** `-j` names a file that is neither JSON nor gzip-compressed JSON
- **THEN** the tool exits with a non-zero status and an error naming the file; it does not silently proceed with an empty schema

### Requirement: Generators write the schema gzipped
`jtaf-provider` and `jtaf-ansible` SHALL write the trimmed schema to their output directory as `trimmed_schema.json.gz` (compact JSON, gzip-compressed) and SHALL NOT write a plain `trimmed_schema.json`. `jtaf-yang2go` and `jtaf-yang2ansible`, which drive them, therefore produce the same file.

#### Scenario: Terraform provider output
- **WHEN** `jtaf-yang2go -p YANG -x CONFIG -t vqfx` completes
- **THEN** `terraform-provider-junos-vqfx/trimmed_schema.json.gz` exists and no `trimmed_schema.json` exists

#### Scenario: Ansible role output
- **WHEN** `jtaf-yang2ansible -p YANG -x CONFIG -t vqfx` completes
- **THEN** `ansible-provider-junos-vqfx/trimmed_schema.json.gz` exists and no `trimmed_schema.json` exists

#### Scenario: Help text
- **WHEN** a user runs any of the four tools with `--help`
- **THEN** the `-j` description names `trimmed_schema.json.gz` and states that plain JSON is also accepted
