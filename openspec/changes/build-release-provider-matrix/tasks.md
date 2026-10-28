# Tasks

## 1. Spike: the 25.4R1 models go through the pipeline

- [x] 1.1 Export `common/models` and `junos-qfx/conf/models` for `25.4R1` from `../yang` with `git archive HEAD:25.4/25.4R1/native/conf-and-rpcs … | tar -x` into a scratch directory under `build/release/`, and verify the extracted tree has 4 common and 33 QFX modules
- [x] 1.2 Run `jtaf-yang2go --generic --groups -p <common> <qfx conf>/*.yang -t qfx` from the scratch directory with `TMPDIR`/`GOTMPDIR` pointed under `build/release/tmp`, under `/usr/bin/time -v`, then `go build` the result; verify the binary runs (`./terraform-provider-junos-qfx -h` or `--debug` banner), and record wall time, peak RSS, schema size and binary size in the design's Context for the preflight threshold and the `JOBS` default
- [x] 1.3 Repeat 1.2 for `25.4R1-EVO` (`common/models` + `conf/models`, `-t evo`) and verify it builds; if pyang or the plugin rejects a 25.4R1 or EVO construct, capture the error in design.md Risks and file it separately before continuing

## 2. Release helper: discovery, resolution, families

- [x] 2.1 Create `scripts/release/jtaf_release.py` with `parse_release()` and `latest()` over a list of names per the regex in design.md Decision 3, a `--branch junos|evo` CLI reading `git ls-tree` at `--repo`/`--ref`; verify with `junosterraform/tests/test_release_build.py` that `24.4R1 24.4R2 25.2R1 25.4R1` and their `-EVO` forms select `25.4R1`/`25.4R1-EVO`, `24.4R2` sorts after `24.4R1`, `24.2R1.17` parses, and an unparseable name is skipped with a warning
- [x] 2.2 Add `resolve <release>` returning the tree path, common path and layout (`native/conf-and-rpcs` + `models` for ≥23.4, flat for ≤22.4), erroring with release and ref when absent; verify with tests over fixture `ls-tree` listings for `22.4R1`, `25.4R1` and `25.4R1-EVO`, and an unknown release
- [x] 2.3 Add `families <release>` mapping family directories to types (`junos`→`mx`, `junos-X`→`X`, EVO `conf`→`evo`), ignoring `common`, `README.txt` and `*-with-extensions`; verify with tests that the Junos fixture yields `mx es ex nfx qfx` and the EVO fixture yields `evo`
- [x] 2.4 Add `tooling-id` hashing the generator scripts, `jtaf_common.py`, the pyang plugin and non-test `terraform_provider/**` Go sources as installed; verify with a test that the id changes when a copied source file is edited and is stable otherwise
- [x] 2.5 Run `python3 scripts/release/jtaf_release.py latest --repo ../yang --branch junos` and `--branch evo` against the real clone and verify they print `25.4R1` and `25.4R1-EVO`

## 3. Makefile: export, generate, build, dist

- [x] 3.1 Create the root `Makefile` with the variables from design.md (`YANG_REPO`, `YANG_REF`, `YANG_DIR`, `JUNOS_RELEASE`, `EVO_RELEASE`, `FAMILIES`, `EXCLUDE`, `JOBS`, `NICENESS`, `MIN_FREE_GB`, `BUILD_DIR`, `DIST_DIR`, `PROVIDER_NAMESPACE`), lazy git-backed defaults, the family↔directory function, exported `TMPDIR`/`GOTMPDIR`, and `make latest`; verify `make latest` prints both releases and the commit and creates nothing, and `make -n latest JUNOS_RELEASE=x EVO_RELEASE=y` makes no git call (`strace -f -e execve` or a stub `git` on PATH)
- [x] 3.2 Add the export stage: the `yang/.inputs` content stamp and the `yang/.exported` rule running `git archive … | tar -x` (or copying from `YANG_DIR`); verify `make build/release/25.4R1/yang/.exported` extracts only `common` and the families' conf modules and that a second run is a no-op
- [x] 3.3 Add the generate and build stages per family: `.inputs-<type>` stamp, `.generated-<type>` rule running `$(NICE) jtaf-yang2go --generic --groups $(EXCLUDE flags) -p … -t <type>` inside the release build dir, and the binary rule running `$(NICE) go build`; wire `make <type>` for each family and `evo`; verify `make qfx` produces `build/release/25.4R1/terraform-provider-junos-qfx/terraform-provider-junos-qfx`, a rerun does nothing, and changing `EXCLUDE` regenerates
- [x] 3.4 Add the dist stage copying each binary to `dist/<rel>/bin/` and its schema to `dist/<rel>/schemas/junos-<type>.trimmed_schema.json.gz`, plus `junos`, `evo` and `all` (`preflight`, then `+$(MAKE) -j$(JOBS) junos evo`, then `terraformrc`) aggregate targets; verify `make -n all JUNOS_RELEASE=25.4R1 EVO_RELEASE=25.4R1-EVO FAMILIES="mx es ex nfx qfx" YANG_COMMIT=0000000` lists exactly six `-t` invocations, and add that check to `test_release_build.py` (skipped without `make`)
- [x] 3.5 Add `make clean` (build dir) and `make distclean` (also `dist/`), and verify `git status --porcelain` stays clean after a build because `/build/` and `/dist/` are already ignored

## 4. Manifest and Terraform CLI configuration

- [x] 4.1 Add `manifest <release-build-dir> <dist-dir>` to `jtaf_release.py` writing `dist/<rel>/manifest.json` with release, commit, repo, built_at, tooling_id and one entry per provider (type, family, binary, sha256, schema, fingerprint read from the gzipped schema, exclude list); verify with a test over a temp dir holding a fake binary and a gzipped `{"fingerprint": "full-groups-abc"}`
- [x] 4.2 Add `terraformrc` writing `dist/<rel>/terraformrc` (one `dev_overrides` entry per provider using `$(PROVIDER_NAMESPACE)`, absolute paths, `direct {}`) and the combined `dist/terraformrc` where the newest release wins per provider and conflicts are listed in a comment; verify with a test over two fixture release dirs sharing `junos-qfx`
- [x] 4.3 Wire both into the Makefile so a release's manifest and terraformrc depend on all of its bins and `make all` ends with the combined file; verify after `make qfx` that `dist/25.4R1/manifest.json` has one entry whose sha256 matches `sha256sum` of the bin and whose fingerprint equals the string in `trimmed_schema.json.gz`

## 5. Preflight and sharing the machine

- [x] 5.1 Create `scripts/release/preflight.sh` checking `git pyang jtaf-yang2go jtaf-provider go tar python3`, free GB on `$(BUILD_DIR)`'s filesystem against `MIN_FREE_GB`, and `MemAvailable` against `JOBS × 8 GB` (warning only), with messages naming the missing tool or the free/required amounts and `make clean-space`; verify with tests that `PATH=/nonexistent` fails naming the tool and `MIN_FREE_GB=100000` fails naming both amounts and the cleanup target
- [x] 5.2 Add `AUTO_CLEAN=1` handling: on a shortfall run the default cleanup then re-check; verify with a test using a stub `df` that the cleanup script is invoked once and the second check passes
- [x] 5.3 Make `preflight` a prerequisite of every build target, prefix `jtaf-yang2go` and `go build` with `nice -n $(NICENESS)` and `ionice -c 2 -n 7` when available; verify during a `make qfx` that `ps -o ni,cls,args` shows the pyang/go processes at niceness 10 and best-effort class 7, and that `ls build/release/tmp` holds the `go-build*` work directory while compiling

## 6. Cleanup

- [x] 6.1 Create `scripts/release/clean_space.sh` with the local step (stale `build/release/<other-release>/`, `build/release/tmp/*`; `CLEAN_EXAMPLE_PROVIDERS=1` for the example output), the container steps (`container prune --filter until=$PRUNE_UNTIL`, `image prune` without `-a`, `builder prune --filter until=`) run only when `docker` is present and the daemon answers, and `df` before/after; verify with a test using a stub `docker` that reports one running and two exited containers that the dry run names only the exited ones and the stub logs no `prune`
- [x] 6.2 Add the opt-ins `PRUNE_UNUSED_IMAGES=1` (`image prune -a`, optional `PRUNE_IMAGE_AGE`), `CLEAN_GO_CACHE=1`, `CLEAN_GO_MODCACHE=1` (confirmation unless `FORCE=1`) and `OWN_TMP=1`; verify with tests that without the flags the dry run lists no `-a` prune and no `go clean`, and with them it lists exactly those commands
- [x] 6.3 Add the no-runtime path (report "container cleanup skipped", still do the local step) and the `make clean-space` target passing the variables through; verify with a test running the script with `docker` absent from `PATH`, and by running `make clean-space DRY_RUN=1` on this machine and checking it lists the exited containers and the dangling images, and none of the running containers

## 7. Documentation

- [x] 7.1 Add `make help` listing targets and variables with their defaults; verify it is the default when `make` is run with an unknown goal hint and that every target in the Makefile appears in it
- [x] 7.2 Add a "Building a release set" section to `README-terraform.md`: what `make all` builds (the six providers for `25.4R1` / `25.4R1-EVO`), the `dist/<release>/` layout, pointing Terraform at `dist/terraformrc` via `TF_CLI_CONFIG_FILE`, generating `.tf` files with the shipped schemas, pinning releases, `EXCLUDE` and the plan-time memory cost of `--generic --groups`, `clean-space` and its opt-ins; verify every command in the section runs as written (`make latest`, `make qfx`, `make clean-space DRY_RUN=1`)

## 8. Integration on the build machine

- [x] 8.1 Run `make clean-space DRY_RUN=1`, review, then `make clean-space`; verify the exited containers and dangling images are gone, `docker ps` still lists every container that was running before, and the reported reclaimed space is positive
- [x] 8.2 Run `make all` (default `JOBS=2`) and record wall time, peak memory and `du -sh dist/`; verify `dist/25.4R1/bin` holds five binaries and `dist/25.4R1-EVO/bin` one, both manifests have the right counts and the commit `96ad7bad…`, and a second `make all` finishes in under a minute with no regeneration
- [x] 8.3 Generate Terraform files for one device with `jtaf-xml2tf -j dist/25.4R1/schemas/junos-qfx.trimmed_schema.json.gz -x examples/evpn-vxlan-dc-groups/dc1/dc1-spine1.xml -t qfx …`, then with `TF_CLI_CONFIG_FILE=$PWD/dist/terraformrc` run `terraform init` and `terraform validate` in that directory; verify Terraform reports the dev override and validates without a schema-fingerprint warning
- [x] 8.4 Run the whole pytest suite and `make -n all` once more; verify both pass, then update design.md's Context with the measured numbers from 1.2 and 8.2
