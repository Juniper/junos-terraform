# Design

## Context

See proposal.md — Why. The pieces that shape the approach:

- `jtaf-yang2go --generic --groups -p <common-dir> <conf-modules…> -t <type>` is the whole generation step. It writes
  `terraform-provider-junos-<type>/` into the current directory (deleting any existing one first), embeds the schema
  compiled and gzipped, stamps a fingerprint `full-groups-<hash of module names and revisions>` into the binary and
  into `trimmed_schema.json.gz`, and exits non-zero on failure. `jtaf-xml2tf` writes `source = "hashicorp/junos-<type>"`
  into the Terraform files it generates, so that is the source address a `dev_overrides` entry has to match.
- The Juniper YANG repository has two layouts. Through 22.4: `<major>/<release>/{common,junos,junos-es,…}/conf/*.yang`.
  From 23.4: `<major>/<release>/native/conf-and-rpcs/{common,junos,…}/conf/models/*.yang`, with `common/models/`, and
  Evolved releases as a sibling `<release>-EVO` whose configuration tree is `conf/models/` with no family level.
  Release names are `X.YRn`, occasionally with a build suffix (`24.2R1.17`). The clone on this machine (`../yang`) has
  every release in its objects but its working tree is partly deleted; HEAD is `96ad7bad`.
- The build machine: 8 cores, 31 GB RAM; `/` is 96% full with `/tmp` on it; `/mnt/extra` holds the sources, Go's
  caches and Docker's root; four containers belonging to another experiment are running; `GOTMPDIR` is unset.
- `--generic --groups` on a full model is measured at about 5.2 GB of Terraform core memory and 13.7 s per plan
  (`openspec/specs/config-groups`, `2026-10-01-group-based-config/design.md`). The provider itself stays around
  210 MB. That cost is paid by whoever plans with these providers, not by the build.
- Measured on this machine during the spike (task 1.2/1.3), `nice -n 10`, nothing else of ours running:

  | Provider | `jtaf-yang2go --generic --groups` | peak RSS | `go build` (warm cache) | `schema.bin.gz` | binary |
  |----------|-----------------------------------|----------|-------------------------|-----------------|--------|
  | QFX 25.4R1 (33 modules) | 2 min 46 s | 6.3 GB | 8 s, 0.3 GB | 2.3 MB | 24.7 MB |
  | EVO 25.4R1-EVO (36 modules) | 3 min 29 s | 8.1 GB | ~10 s | 3.1 MB | 25.4 MB |

  Generation, not compilation, is the expensive step, and its memory is what bounds `JOBS`: two in flight is
  16 GB, four would not fit beside anything else on a 31 GB machine. Disk is not a concern — a release's exported
  modules are ~20 MB and a built provider tree ~30 MB.
- The full matrix (task 8.2), `make all` with `JOBS=2` on the same machine while fifteen unrelated containers ran:
  12 min 38 s of wall time over two runs (the first stopped at the NFX clash, Decision 8; the second finished NFX
  and EVO), peak RSS 10.4 GB, `build/release/` 347 MB and `dist/` 158 MB afterwards. A second `make all` with
  nothing changed takes 1.5 s. Per family the 25.4R1 generation times were MX 3:20, ES 3:30, EX 1:10 (6 modules),
  NFX 4:00, QFX 2:46, EVO 3:29.

## Goals / Non-Goals

**Goals:**
- One command builds the six providers for the newest Junos OS and Junos OS Evolved release, and exactly the parts
  that need rebuilding on later runs.
- Everything the build decides — release, commit, family, exclusions, fingerprint — is written down next to the
  binaries, so a set can be reproduced or told apart from another.
- The build is a well-behaved tenant of a shared machine: it checks before it starts, keeps its temporary files off
  the root filesystem, runs heavy work at low priority, and its cleanup cannot reach anything that is running.
- Works for older releases and both repository layouts, so the small 18.2R3 model remains available as a fast smoke
  test of the whole pipeline.

**Non-Goals:**
- Changing `jtaf-yang2go`, `jtaf-provider`, the generated provider or the example scripts.
- Publishing to a Terraform registry or producing a versioned filesystem mirror; `dev_overrides` is the project's
  installation model and this change stays with it.
- Cross-compiling for other platforms (the provider binaries are for the build host's `GOOS/GOARCH`).
- Deciding what to exclude from the model for anyone: the matrix is fully generic as asked; `EXCLUDE` is passed
  through verbatim.
- Freeing space on the root filesystem beyond the build's own use of it (nothing the build creates goes there).

## Decisions

### 1. A GNU `Makefile` at the repository root, with small helpers under `scripts/release/`

Make gives the matrix its natural shape — one rule per provider, prerequisites on the exported modules, `-j` for
bounded concurrency — and is already on every build machine. The Makefile holds the variables, the targets and the
dependency graph; anything that needs parsing or JSON goes to a helper:

| Helper | Language | Why |
|--------|----------|-----|
| `scripts/release/jtaf_release.py` | Python | `latest`, `resolve` (release → tree path, layout), `families`, `tooling-id`, `manifest`, `terraformrc`: sorting release names, walking `git ls-tree` output, reading the fingerprint out of a gzipped schema and writing JSON are all easier and more testable in Python; the repository's tests are pytest. |
| `scripts/release/preflight.sh` | bash | Tool and disk checks; a dozen lines of `command -v` and `df`. |
| `scripts/release/clean_space.sh` | bash | Drives the `docker` CLI; its dry run prints the commands it would run, which is clearer as a shell script than as subprocess calls. |

*Alternatives considered:* a Python-only driver (`jtaf-release build …`) — loses make's incremental rebuilds and
`-j` for free, and adds a console script the CLI spec would have to cover; `mk/*.mk` includes — more files for no
structural gain at this size; extending `examples/providers/build-*.sh` — those are pinned demos, and a release
build is not an example.

### 2. Modules come out of git objects, not the working tree

`git -C $(YANG_REPO) archive --format=tar $(YANG_REF):<release-tree> <common-path> <family-conf-path> | tar -x -C
build/release/<release>/yang/` exports only the modules a release needs (its `common` and the families'
configuration modules — not `rpc`, `genstate`, `*-with-extensions` or `openconfig`), into a directory the Makefile
owns. The commit hash is `git rev-parse $(YANG_REF)^{commit}` and is recorded everywhere the release name is.

This is what makes the build independent of the clone's checkout state, lets two releases be built from one clone
without switching branches, and keeps the exported set to tens of MB rather than the 374 MB of a whole release tree.
`jtaf_release.py resolve` maps a release name to its tree path and layout: it looks for `<major>/<release>/native/
conf-and-rpcs` and falls back to `<major>/<release>`; a family's modules are at `conf/models` if that exists, else
`conf`; common likewise. `YANG_DIR=<path>` bypasses git altogether for a plain directory that already has that
layout.

*Alternative considered:* requiring the clone to be checked out at the release — fails today on this machine, and
cannot build Junos and EVO from different commits.

### 3. Release discovery and family mapping

`jtaf_release.py latest --branch junos|evo` lists `<major>` directories at the reference, then their entries, parses
`^(\d+)\.(\d+)R(\d+)(?:\.(\d+))?(?:-S(\d+))?(-EVO)?$`, and takes the maximum tuple per branch (`-EVO` → evo). Names
that do not parse are ignored with a warning. `make latest` prints both names and the commit. `JUNOS_RELEASE` and
`EVO_RELEASE` default to those lookups (lazily, via `?=`/`$(shell)`, so setting them on the command line costs no
git call) and either can be pinned independently; `EVO_RELEASE=` (empty) or `JUNOS_RELEASE=` skips that branch.

Family directories become provider types by stripping the `junos-` prefix; the bare `junos` family (MX, ACX, T, some
EX) becomes `mx`, because `jtaf-provider` hard-codes the `junos-` prefix of the output name and `junos-junos` is
unusable. The Evolved tree is one family, `evo`. The inverse mapping (`mx` → `junos`, `evo` → `conf`, else
`junos-<type>`) lives in the Makefile as a function so a single family target (`make qfx`) needs no helper call.
`FAMILIES` defaults to what `jtaf_release.py families` finds at the release (today `mx es ex nfx qfx`), so a family
Juniper adds later appears without a Makefile change; it can be overridden to build a subset.

*Alternative considered:* release in the provider name (`junos-qfx-25-4r1`) — Terraform forbids dots in type names,
every `required_providers` block would change on a release bump, and the fingerprint already distinguishes schemas.
Keeping Juniper's names verbatim for the type (`-t junos-es` → `terraform-provider-junos-junos-es`) — rejected for
the doubled prefix.

### 4. Incremental rebuilds through content stamps

`jtaf-provider` deletes and recreates its output directory on every run, so the generated tree cannot be make's
prerequisite. Each stage instead depends on a small file whose content is the stage's inputs, rewritten only when
that content changes (the usual `cmp`-then-`mv` idiom):

```
build/release/<rel>/yang/.inputs            = release, commit (or YANG_DIR), module paths
build/release/<rel>/yang/.exported          ← .inputs            : git archive | tar
build/release/<rel>/.inputs-<type>          = release, commit, type, EXCLUDE, tooling-id
build/release/<rel>/.generated-<type>       ← .exported, .inputs-<type> : jtaf-yang2go
build/release/<rel>/terraform-provider-junos-<type>/terraform-provider-junos-<type>
                                            ← .generated-<type> : go build
dist/<rel>/bin/terraform-provider-junos-<type>, dist/<rel>/schemas/junos-<type>.trimmed_schema.json.gz
                                            ← the binary        : copy
dist/<rel>/manifest.json, dist/<rel>/terraformrc, dist/terraformrc
                                            ← all of a release's bins : jtaf_release.py manifest / terraformrc
```

`tooling-id` is a hash over the generator scripts, `jtaf_common.py`, the pyang plugin and `terraform_provider/**`
(non-test Go sources) as installed, so editing the generator or the provider source invalidates every provider,
while a `go build` with nothing changed is a cache hit and costs seconds.

### 5. Output layout: `dist/<release>/`, provider names without the release

```
dist/
  25.4R1/
    bin/terraform-provider-junos-{mx,es,ex,nfx,qfx}
    schemas/junos-{mx,es,ex,nfx,qfx}.trimmed_schema.json.gz
    manifest.json
    terraformrc
  25.4R1-EVO/
    bin/terraform-provider-junos-evo
    schemas/junos-evo.trimmed_schema.json.gz
    manifest.json
    terraformrc
  terraformrc          # every release; newest wins per provider, conflicts listed in a comment
```

Junos OS and Evolved are separate releases with separate names and commits, so they get separate directories even
when built together. `terraformrc` is the `provider_installation { dev_overrides { "registry.terraform.io/
hashicorp/junos-<type>" = "<absolute dist>/<rel>/bin" } direct {} }` block the existing `examples/example-terraformrc`
uses, with the namespace in `PROVIDER_NAMESPACE ?= hashicorp` to follow `jtaf-xml2tf`. The schemas are shipped because
`jtaf-xml2tf -j` needs the provider's own schema to generate files the provider will accept; the manifest records
the fingerprint read from each schema so a mismatch can be diagnosed from the files alone.

`manifest.json`: `{"release", "commit", "yang_repo", "built_at", "tooling_id", "providers": [{"type", "family",
"binary", "sha256", "schema", "fingerprint", "exclude": [...]}]}`.

*Alternative considered:* installing into `$GOPATH/bin` as the example scripts do — one directory cannot hold two
releases of the same provider, and nothing records which release is there.

### 6. Preflight, temporary files and priority

`make preflight` (a prerequisite of every build target, skipped by `make -n`) runs `scripts/release/preflight.sh`:
`command -v` for `git pyang jtaf-yang2go jtaf-provider go tar python3`; `df` on `$(BUILD_DIR)` against `MIN_FREE_GB`
(default 5 — a release set is a few hundred MB of modules and binaries plus Go's work and cache growth; 5 GB leaves
room for two providers in flight); a warning, not a failure, when `MemAvailable` is under `JOBS × 8 GB` (the EVO
model peaks at 8.1 GB during generation). On a space
shortfall it prints free, required and `make clean-space`, or with `AUTO_CLEAN=1` runs the default cleanup and checks
again before failing.

The Makefile exports `TMPDIR` and `GOTMPDIR` as `$(abspath $(BUILD_DIR))/tmp`, so `go build`'s work directories, the
compiler's temporary files and anything the Python tools spill land on the build filesystem. `GOFLAGS` is left to the
environment; `GOCACHE`/`GOMODCACHE` are not changed, since sharing them with the rest of the machine is what makes
the second provider's build fast.

Every `jtaf-yang2go` and `go build` invocation is prefixed with `$(NICE)` = `nice -n $(NICENESS)` (default 10), plus
`ionice -c 2 -n 7` when `ionice` exists. `make all` runs `preflight`, then `+$(MAKE) -j$(JOBS) junos evo` (default
`JOBS=2`), then the combined `terraformrc`. Two is the default because each provider's generation peaks at 6–8 GB
(pyang's parse of a full family model, then the schema compiler); four would still fit this machine's 23 GB available
only with nothing else running, and the machine is shared.

### 7. Cleanup: a short allow-list, and everything else opt-in

`scripts/release/clean_space.sh` does, in order, reporting `df` on the build filesystem before and after:

1. **Local**: remove `build/release/<release>/` for releases other than the ones currently selected, and anything
   under `build/release/tmp/`. (Example-provider output under `examples/providers/` and the repository root is left
   alone unless `CLEAN_EXAMPLE_PROVIDERS=1`: it is small and someone may be mid-demo.)
2. **Containers**, only if `docker` is on the PATH and the daemon answers: `docker container prune -f --filter
   until=$(PRUNE_UNTIL)` (default `24h`) — Docker's own definition of prune is stopped containers only, and the age
   filter keeps something that just exited out of reach; `docker image prune -f` (dangling only, no `-a`);
   `docker builder prune -f --filter until=$(PRUNE_UNTIL)` (unreferenced build cache). Never `docker system prune`,
   never `docker volume prune`, never `stop`/`kill`/`rm` on a running container.
3. **Opt-in**: `PRUNE_UNUSED_IMAGES=1` adds `docker image prune -a -f [--filter until=$(PRUNE_IMAGE_AGE)]`, which
   removes tagged images no container uses — on this machine ~15.6 GB, including images other work on the machine
   relies on, which is why it is not default. `CLEAN_GO_CACHE=1` runs `go clean -cache` (2.9 GB; other Go builds recompile).
   `CLEAN_GO_MODCACHE=1` runs `go clean -modcache` (18 GB; every project re-downloads) and is the only step that asks
   for confirmation unless `FORCE=1`. `OWN_TMP=1` removes the invoking user's own entries under `/tmp` older than
   seven days (`find -user $(id -u) -mtime +7`), nothing belonging to anyone else.

`DRY_RUN=1` replaces every removal with a listing: `docker ps -a --filter status=exited --filter until=…`, `docker
images -f dangling=true`, `docker system df`, and the directories matched locally, followed by the exact commands
that would run. The Makefile exposes it as `make clean-space DRY_RUN=1`.

*Alternative considered:* `docker system prune -a` — one command, but it takes volumes, networks and every unused
image, and cannot be made to respect "someone else's work" on a shared host.

### 8. Per-family exclusions, and the NFX `hash-function` clash

`EXCLUDE` applies to every provider; `EXCLUDE_<type>` adds paths for one. Both feed the provider's input stamp,
so changing either regenerates exactly that provider, and the manifest records the effective list.

The matrix needed it on the first full run: `jtaf-provider` rejects the 25.4R1 NFX model with
`…/enhanced-hash-key/hash-function: CRC16 and crc16 both map to attribute 'crc16'`. The `hash-function` choice has
a `leaf CRC16` (BISYNC) and a `container crc16` as sibling cases — the other four families have only the leaf — and
`validate_attribute_names` refuses a schema where two siblings share a Terraform attribute name, since the generic
provider could not tell which the user meant. Walking the pyang JSON found exactly eight occurrences: the choice
under `forwarding-options`, `routing-instances/instance/…`, `logical-systems/routing-instances/instance/…` and
`tenants/routing-instances/instance/…`, each once in the base hierarchy and once under `groups` (which `--exclude`
does not reach on its own, since it removes explicit paths). `EXCLUDE_nfx` defaults to those eight, so the NFX
provider lacks the ECMP hash-function knob and nothing else.

*Alternatives considered:* teaching the generator to disambiguate (`crc16` and `crc16_container`, say) — the
attribute name is also what the provider maps back to an XML element, so it is a change to the Go provider's naming
and its inverse, with a fingerprint change, and belongs in its own change against `code-generation`; excluding
`enhanced-hash-key` — broader than necessary; building NFX without `--groups` — breaks the one-rule matrix for a knob.

### 9. Tests

`junosterraform/tests/test_release_build.py` loads `jtaf_release.py` by path (as `test_script_coverage.py` does for
the CLIs) and covers: release parsing and ordering over a fixture list (including `24.2R1.17` and `-EVO`); family
enumeration over fixture `ls-tree` listings for both layouts; `resolve` for a 22.4 and a 25.4 release; manifest entry
from a temp directory holding a fake binary and a gzipped `{"fingerprint": …}`; `terraformrc` conflict handling with
two releases. `clean_space.sh` is run with `DRY_RUN=1` and a stub `docker` on `PATH` that answers `ps`/`images`/
`system df` from canned text and logs every call; the test asserts the exited containers are named, the running one
is not, and no `prune` call was made. `make -n all JUNOS_RELEASE=25.4R1 EVO_RELEASE=25.4R1-EVO FAMILIES="mx es ex nfx
qfx" YANG_COMMIT=0000000` is parsed for six `-t <type>` invocations (skipped when `make` is absent). Nothing in the
tests needs Docker, Go, pyang or the YANG clone.

## Risks / Trade-offs

- [pyang or the JTAF plugin may choke on something new in the 25.4R1 models (`junos-configuration-metadata.yang`
  is new, and the models have grown since 18.2)] → The first task is a one-family spike (`make qfx`) before the
  matrix; a failure is a bug report against the plugin, not this change, and `JUNOS_RELEASE=18.2R3` keeps a known-good
  model available for exercising the build itself. *Outcome:* pyang and the plugin handled all six models; the one
  failure was the generator's attribute-name check on NFX (Decision 8), handled with a per-family exclusion.
- [`EXCLUDE_nfx`'s default names paths that another release may not have, and `--exclude` of a missing path is an
  error] → The paths are in every NFX release that has `tenants` and `enhanced-hash-key`; for one that does not,
  `EXCLUDE_nfx=` (or a shorter list) is a command-line variable, and the error names the missing path.
- [A full-model pyang run holds 6–8 GB; two in parallel could squeeze a machine with less than this one's 31 GB] →
  `JOBS` defaults to 2, preflight warns on low `MemAvailable`, and `JOBS=1` is one variable away.
- [`--generic --groups` providers cost ~5 GB per plan on the consumer's side] → Stated in the proposal and README;
  `EXCLUDE="logical-systems tenants dynamic-profiles"` is the documented lever and is recorded in the manifest so a
  trimmed set is not mistaken for a full one.
- [`container prune` removes someone's deliberately kept stopped container] → 24 h age filter by default, dry run
  available, and only `docker`'s own prune semantics (stopped only). Running containers are structurally out of reach.
- [`tooling-id` hashing `terraform_provider/**` makes every edit to the provider source invalidate six builds] →
  Intended: a provider built from stale generator or provider source is the harder bug to find. `go build` after
  regeneration is mostly cache hits.
- [`make -n` evaluates `$(shell …)` defaults, so a dry run without pinned releases still calls git] → Those calls are
  reads; the tests pin the variables, and `make latest` is the documented way to see the defaults.
- [Release names outside `X.YRn[.b][-Sn][-EVO]` would be skipped by discovery] → Warned about, and any release can
  still be built by name as long as `resolve` finds its tree.
- [Root filesystem stays at 96%] → Not this build's doing and not touched; the build keeps everything, including
  temporary files, on the build filesystem. `clean-space` with `OWN_TMP=1` removes the invoking user's own `/tmp`
  entries older than seven days as a courtesy, nothing else there.

## Migration Plan

Additive: new files only, no behaviour change for existing scripts, CI or generated code. Rolling back is deleting
the Makefile, `scripts/release/`, the test and the README section. Outputs live in ignored directories.

## Open Questions

- Whether to also exercise `make qfx JUNOS_RELEASE=18.2R3` in CI as a pipeline smoke test — cheap (the model is in
  the repository), but it would be the first CI job to depend on GNU make. Can be decided after the first local run.

## Follow-ups observed, outside this change

Running the built 25.4R1 QFX provider through a plan → apply → plan cycle against a lab device (task 8.3, extended)
worked — the apply changed exactly the host name and the re-plan was a no-op — and surfaced two things that belong to
other capabilities:

- **`<undocumented>` wrappers hide list keys from the reader** (`patch-engine`). Junos renders hidden statements in
  its XML wrapped in `<undocumented>…</undocumented>`; a `syslog file … contents` entry whose facility is hidden comes
  back as `<contents><undocumented><name>local0</name></undocumented><any/></contents>`. The reader does not see the
  key, indexes the entry as `contents/any`, and the diff against a plan that has `contents[name=local0]` is a
  delete-and-recreate the device rejects (`statement creation failed`). Unwrapping `<undocumented>` when reading
  device configuration would fix it.
- **`jtaf-xml2tf` keeps an empty container when every child was dropped as not in the model** (`cli-tools`). A device
  statement whose children the model lacks (`system commit { xpath; constraints }` on that device) becomes
  `commit = [{}]` in the Terraform file; the provider reads the device's `commit` as absent, so every plan is a
  one-attribute in-place update until the block is removed by hand. The tool already warns about the dropped paths;
  dropping the emptied parent too would make the warning the whole story.
