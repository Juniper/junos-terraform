# Spec Delta

## ADDED Requirements

### Requirement: The release build's helpers are tested without a build machine
The test suite SHALL verify the release build's release discovery and ordering, family enumeration, manifest contents
and cleanup dry run using fixtures, without requiring a Juniper YANG clone, a container runtime, Go or the JTAF
generators to be present, and SHALL verify that the build's dry run enumerates every provider in the matrix.

#### Scenario: Release ordering
- **WHEN** the discovery helper is given the release names `24.4R1`, `24.4R2`, `25.2R1`, `25.4R1`, `24.4R1-EVO`,
  `24.4R2-EVO`, `25.2R1-EVO` and `25.4R1-EVO`
- **THEN** it selects `25.4R1` for Junos OS and `25.4R1-EVO` for Junos OS Evolved, and places `24.4R2` after `24.4R1`

#### Scenario: Family enumeration
- **WHEN** the helper is given a Junos OS release tree with `common`, `junos`, `junos-es`, `junos-ex`, `junos-nfx`,
  `junos-qfx` and a `README.txt`, and an Evolved release tree with `common`, `conf` and `conf-with-extensions`
- **THEN** it reports the five Junos families mapped to `mx`, `es`, `ex`, `nfx`, `qfx` and the single Evolved family
  `evo`, and ignores `common`, the `-with-extensions` directories and the README

#### Scenario: Manifest contents
- **WHEN** the manifest helper is given a built provider directory and its binary
- **THEN** the manifest entry carries the release, commit, family, exclusions, fingerprint read from the schema
  archive, the binary's SHA-256 and a build time

#### Scenario: Cleanup dry run against a stubbed runtime
- **WHEN** the cleanup runs in dry-run mode with a stub container runtime reporting one running and two exited
  containers
- **THEN** its output names the two exited containers and no running one, lists the commands it would run, and the
  stub records no removal

#### Scenario: Build dry run covers the matrix
- **WHEN** the build is dry-run for the whole matrix against a fixture release listing
- **THEN** the planned commands name all six providers once each
