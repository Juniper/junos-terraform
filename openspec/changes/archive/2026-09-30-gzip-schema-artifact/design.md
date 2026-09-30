# Design

## Context

See proposal.md — Why. The schema artifact flows through two parallel pipelines that share `jtaf_common`:

```
                    pyang JSON (stdin or file)
                              │
          ┌───────────────────┴───────────────────┐
          ▼                                       ▼
   jtaf-provider  -j                       jtaf-ansible  -j
   (Jinja2 | --generic)                    (group | override)
          │  writes                               │  writes
          ▼                                       ▼
 terraform-provider-junos-<t>/          ansible-provider-junos-<t>/
   trimmed_schema.json.gz   ◀── one form ──▶  trimmed_schema.json.gz
          │  read by                              │  read by
          ▼                                       ▼
   jtaf-xml2tf  -j                         jtaf-xml2yaml  -j
```

Current state:

- Writers: `jtaf-provider --generic` writes `trimmed_schema.json.gz` (compact, `gzip.open(..., "wt")`); `jtaf-provider` (Jinja2, step 7) and `jtaf-ansible` (step 5) write `trimmed_schema.json` with `indent=2`.
- Readers: `jtaf-xml2tf` (`json.load(open(...))`), `jtaf-xml2yaml.load_schema`, and `jtaf_common.filter_json_using_xml` (`-` or file) all read plain JSON only, each with its own few lines.
- Go: `generic/embed.go` and `cmd/compileschema` already sniff the gzip magic bytes and accept either form. The Jinja2 provider inlines the schema into Go source (`TrimmedSchemaJSON`) at render time and never opens the file. The Go side is therefore not touched.

Constraints:

- The `-j -` stdin path must keep working (it is how `jtaf-yang2go` and `jtaf-yang2ansible` pipe pyang into the generators).
- Provider and role directories generated before this change exist in users' workspaces and CI caches; they carry plain `trimmed_schema.json`.
- A full Junos model is hundreds of MB as indented JSON; the generic path already had to go compact + gzip for that reason.

## Goals / Non-Goals

**Goals:**
- One on-disk form and one filename for the schema across both pipelines.
- One code path for reading it, so the accepted forms cannot drift between tools again.
- No flag day for existing directories: old plain files still load.

**Non-Goals:**
- Changing what the schema contains, or the JSON structure.
- Reconciling attribute-name sanitisation between `generic.SanitizeName` and `normalize_tag` (separate change).
- Changing the compiled `schema.bin.gz` or anything on the Go side.
- Making Go read `trimmed_schema.json.gz` from disk at runtime (it embeds; nothing to do).

## Decisions

### D1. Detect gzip by content, not by extension

The shared loader reads the first two bytes and treats `1f 8b` as gzip; anything else is parsed as JSON directly.

- **Why**: matches what `generic/embed.go` and `cmd/compileschema` already do, so Python and Go agree; works for stdin, which has no extension; tolerates a user who renames or decompresses the file; makes backward compatibility free.
- **Alternative — by `.gz` extension**: simpler to read, but cannot handle stdin, and a mis-named file fails with a confusing JSON decode error rather than just working.
- **Alternative — try gzip, fall back to JSON on error**: works, but hides real corruption behind a fallback; sniffing is explicit.

Implementation: one function in `jtaf_common`, e.g. `load_schema_json(path_or_dash) -> dict`, handling `-` (read `sys.stdin.buffer`), open in binary, sniff, decompress via `gzip.decompress`, `json.loads`. On a decode failure raise a `ValueError` naming the input, which the CLIs surface via `parser.error` / `sys.exit`. `jtaf-xml2yaml.load_schema` keeps its `configuration`-node validation but delegates the read.

### D2. Writers emit compact JSON, gzipped, as `trimmed_schema.json.gz` only

- **Why compact**: the generic path already writes compact; once gzipped the file is not human-diffable regardless of indentation, so indentation costs bytes for nothing. Uniform output also means the same schema produces a byte-identical file from both `jtaf-provider` paths.
- **Why not dual-emit (`.json` and `.json.gz`)**: keeps a second artifact that can drift and doubles disk for full models; the loader's plain-JSON acceptance already covers compatibility for *reading*. The proposal deliberately makes this a **BREAKING** filename change on output.
- **Alternative — keep `.json` name but gzip the bytes**: content-sniffing readers would cope, but a `.json` file that is not JSON surprises every other tool (editors, `jq`, `python -m json.tool`). The `.gz` suffix is honest.

A small shared writer, e.g. `write_schema_json_gz(resources, path)`, used by all three writer sites, so `separators=(",", ":")` and the gzip settings are defined once.

### D3. Keep the `-j` argument semantics; change the help text

`-j` stays the single schema input on all four tools; only the documented filename changes. No new flag (e.g. `--gzip`) is introduced — the format is inferred (D1).

### D4. Tests assert absence of the old file

Workflow tests assert `trimmed_schema.json.gz` exists **and** `trimmed_schema.json` does not, so a regression that quietly resumes dual-writing is caught.

## Risks / Trade-offs

- [Users' scripts pass `-j .../trimmed_schema.json` and the file is gone after regeneration] → BREAKING is called out in proposal and changelog; the error is an immediate `FileNotFoundError` naming the path, not silent misbehaviour; all in-repo scripts, CI and docs are updated in the same change.
- [`jtaf-xml2yaml.load_schema` and `filter_json_using_xml` are imported by tests directly] → keep their names and signatures; only the body delegates to the shared loader.
- [Stdin sniffing needs the binary stream] → read `sys.stdin.buffer`, not `sys.stdin`; a unit test covers gzipped stdin.
- [Content sniffing misclassifies a JSON file that happens to start with `1f 8b`] → impossible; valid JSON cannot start with those bytes.
- [`build-generic-full.sh` reports `du -sh trimmed_schema.json`] → update to `.gz`; cosmetic but would fail under `set -e`.

## Migration Plan

1. Land loader + writer helpers in `jtaf_common` with unit tests (readers accept both forms; nothing observable changes yet).
2. Switch the three writer sites to the helper; update workflow tests to assert the new file and absence of the old.
3. Update CI workflows, example scripts, READMEs, demo and `changelog.md` to the `.gz` filename.
4. Rollback: revert writer sites to plain output; readers still accept both, so a partial rollback cannot break consumers.

Existing generated directories need no action to be *read*; regenerating them produces the new filename.
