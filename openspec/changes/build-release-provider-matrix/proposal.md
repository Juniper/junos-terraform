# Proposal

## Why

JTAF can build a generic, groups-aware provider for any one Junos family, but only by hand: the scripts under
`examples/providers/` are pinned to the 18.2R3 QFX and SRX models copied into `examples/yang/`, and nothing enumerates
the families Juniper publishes or the release they are published for. Producing the set of providers for a current
release means six hand-run `jtaf-yang2go` invocations against a 900 MB YANG checkout, on a build machine whose root
filesystem is 96% full and whose Docker root carries 15 GB of images nothing uses.

The Juniper YANG repository (`../yang`, HEAD `96ad7bad`, "Publish content of 24.4R2, 25.2R1 and 25.4R1 releases") now
carries the two latest releases this change targets:

| Branch | Release | Models | Families |
|--------|---------|--------|----------|
| Junos OS | **25.4R1** | `25.4/25.4R1/native/conf-and-rpcs/` | `junos` (MX/ACX/T), `junos-es` (SRX), `junos-ex`, `junos-nfx`, `junos-qfx` |
| Junos OS Evolved | **25.4R1-EVO** | `25.4/25.4R1-EVO/native/conf-and-rpcs/` | one unified `conf/` tree |

## What Changes

- Add a top-level `Makefile` that builds the full provider matrix for a named release pair — by default the newest
  Junos and the newest EVO release found in the Juniper YANG repository — as `--generic --groups` providers, one per
  family: `junos-mx`, `junos-es`, `junos-ex`, `junos-nfx`, `junos-qfx` for Junos OS and `junos-evo` for Junos OS
  Evolved. Any single family, the Junos set, or the EVO set can be built on its own.
- Read the YANG modules straight out of the Juniper repository's git objects (`git archive <ref>:<release-path>`), so
  the build does not depend on the state of that repository's working tree and does not need the repository checked
  out at the release being built. The release name, the YANG commit and the family are recorded with every provider.
- Lay the results out per release under `dist/<release>/` — the provider binaries, each provider's
  `trimmed_schema.json.gz` (which `jtaf-xml2tf` needs to generate matching Terraform files), a `manifest.json`, and a
  generated `.terraformrc` with `dev_overrides` so Terraform can use the set without a registry. Releases coexist on
  disk; switching is a matter of which `.terraformrc` is active. Provider names stay `junos-<family>`; the release
  lives in the path and in the embedded schema fingerprint, not in the provider name.
- Add a `preflight` step that fails early when a tool is missing or free disk space on the build filesystem is below
  a threshold, and route Go's and the pipeline's temporary files to the build filesystem — on this machine `/tmp` sits
  on the 499 MB-free root filesystem and `go build` would otherwise write its work directories there.
- Add a `clean-space` target that reclaims only what is demonstrably unused: stopped containers, dangling images,
  unreferenced Docker build cache, and stale JTAF provider output from earlier builds. It never touches a running
  container. Reclaiming tagged images that no container uses, trimming the Go build cache, and purging the Go module
  cache are each behind an explicit opt-in, and every removal can be previewed with a dry run.
- Run the heavy steps (pyang, schema compilation, `go build`) at reduced scheduling priority and with bounded
  parallelism, so a matrix build shares the machine with other builds instead of starving them.
- Add tests for the release discovery, the manifest, and the cleanup's dry run, and document the workflow in
  `README-terraform.md`.

Nothing changes in the CLI tools, the generated provider, or the existing example scripts. `--generic --groups` is the
expensive combination (`jtaf-provider` warns that it roughly doubles the model, about 5 GB of Terraform core memory at
plan time against a full model); the matrix builds it as requested, and an `EXCLUDE` variable passes `--exclude`
subtrees (for example `logical-systems tenants dynamic-profiles`) through to every provider for anyone who needs a
lighter set.

## Capabilities

### New Capabilities
- `release-build`: Building the complete set of generic, groups-aware providers for a named Junos OS and Junos OS
  Evolved release from the Juniper YANG repository — release discovery, family enumeration, per-release output layout
  with manifest and `dev_overrides`, preflight resource checks, and build-machine cleanup that is safe to run beside
  other work.

### Modified Capabilities
- `testing`: The test suite SHALL cover the release build's helpers — release discovery and ordering, family
  enumeration, manifest contents, and the cleanup dry run — without a Docker daemon or a YANG checkout present.

## Impact

- New: `Makefile` (repository root), helper scripts under `scripts/release/`, tests under `junosterraform/tests/`.
- Modified: `README-terraform.md` (how to build a release set and point Terraform at it). The outputs land under
  `build/` and `dist/`, which `.gitignore` already covers.
- Unchanged: `jtaf-yang2go`, `jtaf-provider`, `terraform_provider/`, `examples/providers/*.sh`, CI workflows.
- Dependencies: GNU make, git, pyang, Go and the installed JTAF CLIs; Docker only when `clean-space` is asked to prune
  it. A clone of `Juniper/yang` reachable as `YANG_REPO` (default `../yang`).
- Build machine: a full matrix build needs a few GB of free space on the build filesystem (the exported configuration
  modules, six provider build trees at ~25 MB each, Go's work directories and build-cache growth) and several GB of
  RAM per concurrent provider while pyang and the schema compiler run; the default parallelism is two.
