# Spec Delta

## Purpose

Builds the complete set of generic, groups-aware Terraform providers for a named Junos OS and Junos OS Evolved release
from the Juniper YANG repository, lays them out per release for Terraform to use, and keeps the build machine's disk
and CPU usable for other work while doing so.

## ADDED Requirements

### Requirement: The release matrix is built from the Juniper YANG repository
The release build SHALL read the YANG modules for a release from a clone of the Juniper YANG repository at a given git
reference, independent of that clone's working tree. A release SHALL be identified by Juniper's release name
(`25.4R1`, `25.4R1-EVO`). For a Junos OS release the build SHALL enumerate every family directory that carries
configuration modules and SHALL build one provider per family; for a Junos OS Evolved release it SHALL build one
provider from the unified configuration tree. Every provider SHALL be generated with the whole model embedded and
configuration groups managed. A release name that does not exist at the reference SHALL fail the build naming the
release and the reference.

#### Scenario: Junos OS 25.4R1 builds five family providers
- **WHEN** the release build is run for Junos release `25.4R1` against a Juniper YANG clone whose reference publishes it
- **THEN** it produces providers `junos-mx`, `junos-es`, `junos-ex`, `junos-nfx` and `junos-qfx`, each embedding the
  full model of its family with the `groups` subtree and `apply-groups` leaf-list exposed

#### Scenario: Junos OS Evolved 25.4R1-EVO builds one provider
- **WHEN** the release build is run for EVO release `25.4R1-EVO`
- **THEN** it produces the single provider `junos-evo`, embedding the full unified EVO model with groups exposed

#### Scenario: The working tree of the YANG clone does not matter
- **WHEN** the Juniper YANG clone's working tree has the release's files deleted or is checked out at another release
- **THEN** the build still exports the release's modules from the given reference and completes

#### Scenario: Unknown release
- **WHEN** the release build is asked for a release name the reference does not publish
- **THEN** it exits non-zero before generating anything, naming the release and the reference it searched

### Requirement: The latest release of each branch is discovered and named
The release build SHALL determine, from the Juniper YANG repository at the given reference, the newest Junos OS release
and the newest Junos OS Evolved release, ordering releases by major version, minor version and R-number and treating
names ending in `-EVO` as the Evolved branch. These SHALL be the defaults when no release is named. The build SHALL
be able to print the two names and the repository commit they come from without building anything, and SHALL accept
an explicit release name for either branch instead.

#### Scenario: Defaults resolve to the newest published releases
- **WHEN** the repository's reference publishes `24.4R1`, `24.4R2`, `25.2R1`, `25.4R1` and their `-EVO` counterparts
  and no release is named
- **THEN** the build selects `25.4R1` for Junos OS and `25.4R1-EVO` for Junos OS Evolved

#### Scenario: R-number ordering within a version
- **WHEN** both `24.4R1` and `24.4R2` are published and no newer version exists
- **THEN** `24.4R2` is selected over `24.4R1`

#### Scenario: Naming the releases without building
- **WHEN** the user asks the build which releases it would build
- **THEN** it prints the Junos OS release, the Evolved release and the repository commit hash, and creates no outputs

#### Scenario: Pinning a release
- **WHEN** the user names `25.2R1` as the Junos OS release
- **THEN** the Junos family providers are built from `25.2R1` while the Evolved release remains whatever was selected
  or named for that branch

### Requirement: Any subset of the matrix can be built
The release build SHALL build the whole matrix, only the Junos OS families, only the Evolved provider, or any single
named family, and SHALL rebuild a provider only when its inputs (the release, the reference, the family's modules, the
exclusions, or the JTAF tooling that generates it) have changed since it was last built. Subtrees MAY be excluded
from every provider or from one named provider; the exclusions in effect for a provider SHALL be recorded with it.

#### Scenario: One family
- **WHEN** the user asks for only `qfx`
- **THEN** only the `junos-qfx` provider is generated and compiled

#### Scenario: Unchanged provider is not rebuilt
- **WHEN** the whole matrix has been built and the user runs the build again with nothing changed
- **THEN** no provider is regenerated or recompiled

#### Scenario: Exclusions for one provider
- **WHEN** a subtree is excluded for the `nfx` provider only and the matrix is built
- **THEN** the `junos-nfx` provider is generated without that subtree, every other provider is generated with it, and
  the `junos-nfx` manifest entry lists the excluded path

#### Scenario: A model the generator refuses is built with a documented exclusion
- **WHEN** a family's model has two sibling nodes that map to the same Terraform attribute name, as the 25.4R1 NFX
  `hash-function` choice does
- **THEN** the release build's default exclusions for that family leave out the offending subtree at each place it
  occurs, including under `groups`, the provider builds, and the manifest records the paths

### Requirement: Outputs are laid out per release with a manifest and dev overrides
For each release built, the release build SHALL place under a release-named output directory: every provider binary
named `terraform-provider-junos-<family>`, each provider's `trimmed_schema.json.gz` so Terraform files can be generated
with the matching schema, a manifest recording for each provider the release, the YANG repository commit, the family,
the exclusions applied, the embedded schema fingerprint, the binary's SHA-256 and the build time, and a Terraform CLI
configuration whose `dev_overrides` map each provider's source address to that directory. Provider names SHALL NOT
contain the release; two releases SHALL be able to coexist on disk. A combined CLI configuration covering every built
release SHALL also be written; where two releases carry the same provider, the newest release SHALL be mapped and the
conflict reported.

#### Scenario: Release directory contents
- **WHEN** Junos release `25.4R1` has been built
- **THEN** the release directory holds five binaries, five schema archives, a manifest with one entry per provider, and
  a CLI configuration with one `dev_overrides` entry per provider

#### Scenario: Terraform uses the release set
- **WHEN** Terraform is run with the release's CLI configuration and a configuration that requires `juniper/junos-qfx`
- **THEN** Terraform loads the provider from the release directory without consulting a registry

#### Scenario: Schema fingerprint is reproducible from the manifest
- **WHEN** a Terraform file generated by `jtaf-xml2tf` from a release's `trimmed_schema.json.gz` is planned against
  that release's provider
- **THEN** the provider accepts it without the schema-mismatch warning, and the fingerprint in the plan matches the
  manifest entry

#### Scenario: Two releases coexist
- **WHEN** `25.2R1` and `25.4R1` have both been built
- **THEN** both release directories remain complete, and the combined CLI configuration maps `junos-qfx` to `25.4R1`
  and reports that `25.2R1` also provides it

### Requirement: Preflight checks fail fast and keep temporary files off the root filesystem
Before generating anything the release build SHALL verify that the tools it needs are available and that the
filesystem holding its build and output directories has at least a configurable amount of free space, failing with a
message that names what is missing or how much space is short and how to reclaim it. All temporary files the build
and its Go compilations create SHALL be placed under the build directory, not the system temporary directory.

#### Scenario: Missing tool
- **WHEN** `pyang` or `jtaf-yang2go` or `go` is not on the PATH
- **THEN** the build exits non-zero before exporting any modules, naming the missing tool

#### Scenario: Insufficient disk space
- **WHEN** the build filesystem has less free space than the configured minimum
- **THEN** the build exits non-zero, reports the free and required amounts, and points at the cleanup target

#### Scenario: Go work directories stay on the build filesystem
- **WHEN** a provider is compiled
- **THEN** the compiler's temporary work directory and the pipeline's temporary files are created under the build
  directory

### Requirement: Cleanup reclaims only what is unused and never disturbs running work
The release build SHALL provide a cleanup that removes stopped containers, dangling images, unreferenced image build
cache and this repository's stale provider build output, and SHALL report the space reclaimed. It SHALL NOT stop,
remove or otherwise affect a running container, and SHALL NOT remove tagged images that no container uses, the Go
build cache or the Go module cache unless each is separately and explicitly requested. A dry-run mode SHALL list every
removal that would be made, and the files and containers it names, without removing anything. When the container
runtime is absent, cleanup SHALL skip the container steps and still perform the local ones.

#### Scenario: Default cleanup beside running containers
- **WHEN** cleanup runs while containers are running and others have exited
- **THEN** the exited containers, dangling images and unreferenced build cache are removed, the running containers and
  the images they use are untouched, and the reclaimed space is reported

#### Scenario: Tagged images need an explicit request
- **WHEN** cleanup runs without the unused-images option
- **THEN** no tagged image is removed, even if no container uses it

#### Scenario: Aggressive options are opt-in
- **WHEN** cleanup is run with the unused-images option and the Go build cache option
- **THEN** images unused by any container (optionally only those older than a given age) and the Go build cache are
  removed, and the Go module cache is still kept

#### Scenario: Dry run
- **WHEN** cleanup is run in dry-run mode
- **THEN** it prints each container, image and directory it would remove and the commands it would run, and nothing
  is removed

#### Scenario: No container runtime
- **WHEN** cleanup runs on a machine without the container runtime on the PATH
- **THEN** it reports that container cleanup was skipped and still removes stale provider build output

### Requirement: The build shares the machine
The release build SHALL run its CPU- and memory-heavy steps at reduced scheduling priority and SHALL bound how many
providers it generates concurrently to a configurable number, defaulting to two.

#### Scenario: Default concurrency
- **WHEN** the whole matrix is built without a concurrency setting
- **THEN** at most two providers are generated or compiled at any one time

#### Scenario: Reduced priority
- **WHEN** modules are processed or a provider is compiled
- **THEN** those processes run at a lower scheduling priority than the default
