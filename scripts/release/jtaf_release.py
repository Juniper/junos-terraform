#!/usr/bin/env python3
"""Release helper for the provider matrix build (see the Makefile).

Answers the questions the Makefile cannot: which releases the Juniper YANG
repository publishes, where a release's modules live in either of the
repository's two layouts, which families it carries, and what identifies the
tooling a provider was generated with. Later subcommands write the manifest
and the Terraform CLI configuration for a built release.

The repository is read through ``git ls-tree`` at a reference, so the state
of its working tree does not matter; ``--dir`` reads a plain directory laid
out the same way instead.
"""

from __future__ import annotations

import argparse
import functools
import gzip
import hashlib
import importlib.util
import json
import os
import re
import subprocess
import sys
from dataclasses import dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Optional

REPO_ROOT = Path(__file__).resolve().parents[2]

# X.YRn, optionally with a build number (24.2R1.17), a service release
# (-S2) and the Evolved marker.
RELEASE_RE = re.compile(r"^(\d+)\.(\d+)R(\d+)(?:\.(\d+))?(?:-S(\d+))?(-EVO)?$")
MAJOR_RE = re.compile(r"^\d+\.\d+$")

# The bare "junos" family (MX, ACX, T) cannot be a provider type: jtaf-provider
# hard-codes the "junos-" prefix of the output name.
BARE_FAMILY_TYPE = "mx"
EVO_TYPE = "evo"


class ReleaseError(Exception):
    pass


# --------------------------------------------------------------------------
# Trees: the same questions asked of git objects or of a plain directory.
# --------------------------------------------------------------------------

class Tree:
    """A read-only tree with ``ls`` returning name -> 'tree'|'blob'."""

    label = ""

    def ls(self, path: str) -> dict[str, str]:
        raise NotImplementedError

    def is_dir(self, path: str) -> bool:
        parent, _, name = path.rstrip("/").rpartition("/")
        if not name:
            return True
        return self.ls(parent).get(name) == "tree"

    def dirs(self, path: str) -> list[str]:
        return sorted(n for n, kind in self.ls(path).items() if kind == "tree")

    def commit(self) -> str:
        return "local"


class GitTree(Tree):
    def __init__(self, repo: str, ref: str = "HEAD"):
        self.repo = repo
        self.ref = ref
        self.label = f"{repo} at {ref}"

    def _git(self, *args: str) -> subprocess.CompletedProcess:
        return subprocess.run(["git", "-C", self.repo, *args],
                              capture_output=True, text=True)

    @functools.lru_cache(maxsize=None)
    def ls(self, path: str) -> dict[str, str]:
        spec = f"{self.ref}:{path}" if path else self.ref
        out = self._git("ls-tree", spec)
        if out.returncode != 0:
            return {}
        entries = {}
        for line in out.stdout.splitlines():
            meta, _, name = line.partition("\t")
            parts = meta.split()
            if len(parts) >= 2 and name:
                entries[name] = parts[1]
        return entries

    def commit(self) -> str:
        out = self._git("rev-parse", f"{self.ref}^{{commit}}")
        if out.returncode != 0:
            raise ReleaseError(
                f"cannot resolve {self.ref} in {self.repo}: {out.stderr.strip()}")
        return out.stdout.strip()

    def archive(self, tree: str, paths: list[str]) -> bytes:
        spec = f"{self.ref}:{tree}" if tree else self.ref
        out = subprocess.run(
            ["git", "-C", self.repo, "archive", "--format=tar", spec, *paths],
            capture_output=True)
        if out.returncode != 0:
            raise ReleaseError(
                f"git archive {spec} failed: {out.stderr.decode().strip()}")
        return out.stdout


class DirTree(Tree):
    def __init__(self, root: str):
        self.root = Path(root)
        self.label = str(self.root)

    def ls(self, path: str) -> dict[str, str]:
        base = self.root / path if path else self.root
        if not base.is_dir():
            return {}
        return {p.name: "tree" if p.is_dir() else "blob" for p in base.iterdir()}


def open_tree(args: argparse.Namespace) -> Tree:
    if getattr(args, "dir", None):
        return DirTree(args.dir)
    repo = args.repo or os.environ.get("YANG_REPO") or str(REPO_ROOT.parent / "yang")
    if not Path(repo).exists():
        raise ReleaseError(
            f"YANG repository not found at {repo}; clone "
            f"https://github.com/Juniper/yang there or set YANG_REPO")
    return GitTree(repo, args.ref or "HEAD")


# --------------------------------------------------------------------------
# Releases
# --------------------------------------------------------------------------

@dataclass(frozen=True, order=True)
class Release:
    key: tuple
    name: str = field(compare=False)
    evo: bool = field(compare=False)

    @property
    def major(self) -> str:
        return f"{self.key[0]}.{self.key[1]}"


def parse_release(name: str) -> Optional[Release]:
    m = RELEASE_RE.match(name)
    if not m:
        return None
    major, minor, rnum, build, service, evo = m.groups()
    key = (int(major), int(minor), int(rnum), int(service or 0), int(build or 0))
    return Release(key=key, name=name, evo=bool(evo))


def discover(tree: Tree, warn=None) -> list[Release]:
    """Every release the tree publishes, from ``<major>/<release>`` directories."""
    found = []
    for major in tree.dirs(""):
        if not MAJOR_RE.match(major):
            continue
        for name in tree.dirs(major):
            rel = parse_release(name)
            if rel is None:
                if warn:
                    warn(f"skipping {major}/{name}: not a release name")
                continue
            found.append(rel)
    return sorted(found)


def latest(releases: list[Release], evo: bool) -> Optional[Release]:
    matching = [r for r in releases if r.evo == evo]
    return max(matching) if matching else None


def major_dir_of(name: str, tree: Tree) -> Optional[str]:
    rel = parse_release(name)
    if rel is not None:
        return rel.major
    # A name the pattern does not know: find the major directory that has it.
    for major in tree.dirs(""):
        if MAJOR_RE.match(major) and name in tree.dirs(major):
            return major
    return None


@dataclass
class Layout:
    release: str
    major: str
    tree: str                 # path within the repository
    layout: str               # "native" (>= 23.4) or "flat" (<= 22.4)
    common: str               # relative to tree
    families: dict[str, str]  # provider type -> conf modules dir, relative to tree

    def family_dir(self, type_: str) -> str:
        return self.families[type_].split("/")[0]

    def as_dict(self) -> dict:
        return {"release": self.release, "major": self.major, "tree": self.tree,
                "layout": self.layout, "common": self.common,
                "families": self.families}


def _conf_dir(tree: Tree, base: str) -> Optional[str]:
    for candidate in (f"{base}/conf/models", f"{base}/conf"):
        if tree.is_dir(candidate) and tree.ls(candidate):
            return candidate
    return None


def _type_for(family_dir: str) -> str:
    if family_dir == "junos":
        return BARE_FAMILY_TYPE
    if family_dir.startswith("junos-"):
        return family_dir[len("junos-"):]
    return family_dir


def _families(tree: Tree, root: str) -> dict[str, str]:
    """Provider type -> configuration modules directory, relative to root."""
    families: dict[str, str] = {}
    if tree.is_dir(f"{root}/conf"):
        # Evolved: one unified configuration tree, no family level.
        conf = _conf_dir(tree, root)
        if conf:
            families[EVO_TYPE] = conf[len(root) + 1:]
        return families
    for entry in tree.dirs(root):
        if entry == "common" or entry.endswith("-with-extensions"):
            continue
        conf = _conf_dir(tree, f"{root}/{entry}")
        if conf is not None:
            families[_type_for(entry)] = conf[len(root) + 1:]
    return families


def resolve(tree: Tree, name: str) -> Layout:
    major = major_dir_of(name, tree)
    if major is None:
        raise ReleaseError(f"release {name} not found in {tree.label}")

    native = f"{major}/{name}/native/conf-and-rpcs"
    flat = f"{major}/{name}"
    if tree.is_dir(native):
        root, layout = native, "native"
    elif tree.is_dir(flat):
        root, layout = flat, "flat"
    else:
        raise ReleaseError(f"release {name} not found in {tree.label}")

    common = next((c for c in ("common/models", "common")
                   if tree.is_dir(f"{root}/{c}") and tree.ls(f"{root}/{c}")), None)
    if common is None:
        raise ReleaseError(f"release {name} in {tree.label} has no common modules")

    families = _families(tree, root)
    if not families:
        raise ReleaseError(
            f"release {name} in {tree.label} has no configuration modules")
    return Layout(name, major, root, layout, common, families)


def make_vars(layout: Layout, prefix: str) -> str:
    lines = [
        f"{prefix}_RELEASE := {layout.release}",
        f"{prefix}_TREE := {layout.tree}",
        f"{prefix}_LAYOUT := {layout.layout}",
        f"{prefix}_COMMON := {layout.common}",
        f"{prefix}_TYPES := {' '.join(layout.families)}",
    ]
    for type_, conf in layout.families.items():
        lines.append(f"{prefix}_CONF_{type_} := {conf}")
    return "\n".join(lines) + "\n"


# --------------------------------------------------------------------------
# Tooling identity
# --------------------------------------------------------------------------

def _package_dir(name: str, root: Optional[Path]) -> Optional[Path]:
    if root is not None:
        p = root / name
        return p if p.is_dir() else None
    spec = importlib.util.find_spec(name)
    if spec and spec.origin:
        return Path(spec.origin).parent
    p = REPO_ROOT / name
    return p if p.is_dir() else None


def tooling_files(root: Optional[Path] = None) -> list[tuple[str, Path]]:
    """The generator and provider sources a provider's content depends on."""
    files: list[tuple[str, Path]] = []
    plan = {
        "junosterraform": lambda d: [d / "jtaf-provider", d / "jtaf-yang2go",
                                     d / "jtaf_common.py"],
        "jtaf_pyang_plugin": lambda d: [d / "jtaf_json.py", d / "__init__.py"],
        "terraform_provider": lambda d: (
            [p for p in d.rglob("*.go")
             if not p.name.endswith("_test.go") and "tools" not in p.parts]
            + [d / "go.mod", d / "go.sum"]),
    }
    for pkg, select in plan.items():
        d = _package_dir(pkg, root)
        if d is None:
            continue
        for p in select(d):
            if p.is_file():
                files.append((f"{pkg}/{p.relative_to(d).as_posix()}", p))
    return sorted(files)


def tooling_id(root: Optional[Path] = None) -> str:
    h = hashlib.sha256()
    files = tooling_files(root)
    if not files:
        raise ReleaseError("no JTAF tooling found to identify")
    for rel, p in files:
        h.update(rel.encode())
        h.update(b"\0")
        h.update(p.read_bytes())
        h.update(b"\0")
    return h.hexdigest()[:12]


# --------------------------------------------------------------------------
# Manifest and Terraform CLI configuration for a built release
# --------------------------------------------------------------------------

BINARY_PREFIX = "terraform-provider-junos-"
FINGERPRINT_RE = re.compile(rb'"fingerprint"\s*:\s*"([^"]+)"')


def _utc(ts: Optional[float] = None) -> str:
    dt = datetime.fromtimestamp(ts, timezone.utc) if ts else datetime.now(timezone.utc)
    return dt.replace(microsecond=0).isoformat().replace("+00:00", "Z")


def _sha256(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()


def schema_fingerprint(schema_gz: Path) -> Optional[str]:
    """The fingerprint jtaf-provider wrote into a (gzipped) schema.

    The key is the last one in the document, so the last match is the one at
    the top level rather than any node that happens to be called fingerprint.
    """
    with gzip.open(schema_gz, "rb") as f:
        data = f.read()
    matches = FINGERPRINT_RE.findall(data)
    return matches[-1].decode() if matches else None


def _read_stamp(path: Path) -> dict[str, str]:
    """A ``key=value`` per line stamp written by the Makefile."""
    values: dict[str, str] = {}
    if path.is_file():
        for line in path.read_text().splitlines():
            key, _, value = line.partition("=")
            values[key.strip()] = value.strip()
    return values


def _provider_entry(build_dir: Path, dist_dir: Path, binary: Path) -> dict:
    type_ = binary.name[len(BINARY_PREFIX):]
    # What the generated output was built from, recorded by the Makefile.
    stamp = _read_stamp(build_dir / f"inputs-{type_}.txt")
    modules = stamp.get("modules", "").split()
    conf = modules[-1] if modules else ""
    family_dir = conf.split("/")[0] if conf else ""
    schema_rel = Path("schemas") / f"junos-{type_}.trimmed_schema.json.gz"
    schema = dist_dir / schema_rel
    return {
        "type": type_,
        "name": f"junos-{type_}",
        "family": EVO_TYPE if family_dir == "conf" else (family_dir or None),
        "modules": modules,
        "exclude": stamp.get("exclude", "").split(),
        "binary": str(Path("bin") / binary.name),
        "sha256": _sha256(binary),
        "size": binary.stat().st_size,
        "built_at": _utc(binary.stat().st_mtime),
        "schema": str(schema_rel) if schema.is_file() else None,
        "fingerprint": schema_fingerprint(schema) if schema.is_file() else None,
    }


def write_manifest(build_dir: Path, dist_dir: Path, source: str, commit: str,
                   tooling: str) -> dict:
    release = dist_dir.name
    binaries = sorted((dist_dir / "bin").glob(f"{BINARY_PREFIX}*"))
    if not binaries:
        raise ReleaseError(f"no provider binaries under {dist_dir / 'bin'}")
    rel = parse_release(release)
    manifest = {
        "release": release,
        "branch": "evo" if (rel.evo if rel else release.endswith("-EVO")) else "junos",
        "commit": commit,
        "source": source,
        "tooling_id": tooling,
        "generator": "jtaf-yang2go --generic --groups",
        "built_at": _utc(),
        "providers": [_provider_entry(build_dir, dist_dir, b) for b in binaries],
    }
    _write_atomic(dist_dir / "manifest.json", json.dumps(manifest, indent=2) + "\n")
    return manifest


def _write_atomic(path: Path, text: str) -> None:
    tmp = path.with_suffix(path.suffix + ".tmp")
    tmp.write_text(text)
    os.replace(tmp, path)


def _overrides_block(entries: dict[str, str]) -> str:
    width = max((len(k) for k in entries), default=0) + 2
    lines = ["provider_installation {", "  dev_overrides {"]
    for source, directory in entries.items():
        lines.append(f'    {json.dumps(source):{width}} = {json.dumps(directory)}')
    lines += ["  }", "  direct {}", "}", ""]
    return "\n".join(lines)


def write_release_terraformrc(dist_dir: Path, namespace: str) -> Path:
    manifest = json.loads((dist_dir / "manifest.json").read_text())
    bin_dir = str((dist_dir / "bin").resolve())
    entries = {f"registry.terraform.io/{namespace}/{p['name']}": bin_dir
               for p in manifest["providers"]}
    header = (f"# Junos release {manifest['release']} (YANG commit {manifest['commit']}), "
              f"generated by make.\n"
              f"# Use with: TF_CLI_CONFIG_FILE={(dist_dir / 'terraformrc').resolve()}\n")
    path = dist_dir / "terraformrc"
    _write_atomic(path, header + _overrides_block(entries))
    return path


def write_combined_terraformrc(dist_root: Path, namespace: str) -> tuple[Path, dict]:
    """One CLI configuration over every built release; the newest wins a name."""
    chosen: dict[str, tuple[Release, str]] = {}
    also: dict[str, list[str]] = {}
    for manifest_path in sorted(dist_root.glob("*/manifest.json")):
        manifest = json.loads(manifest_path.read_text())
        release = parse_release(manifest["release"])
        if release is None:
            continue
        bin_dir = str((manifest_path.parent / "bin").resolve())
        for p in manifest["providers"]:
            name = p["name"]
            current = chosen.get(name)
            if current is None or release > current[0]:
                if current is not None:
                    also.setdefault(name, []).append(current[0].name)
                chosen[name] = (release, bin_dir)
            else:
                also.setdefault(name, []).append(release.name)
    if not chosen:
        raise ReleaseError(f"no built releases under {dist_root}")

    lines = ["# Every release built under this directory; the newest release provides each name.",
             f"# Use with: TF_CLI_CONFIG_FILE={(dist_root / 'terraformrc').resolve()}"]
    for name in sorted(chosen):
        line = f"# {name}: {chosen[name][0].name}"
        if name in also:
            line += f"  (also built: {', '.join(sorted(also[name]))})"
        lines.append(line)
    entries = {f"registry.terraform.io/{namespace}/{name}": chosen[name][1]
               for name in sorted(chosen)}
    path = dist_root / "terraformrc"
    _write_atomic(path, "\n".join(lines) + "\n" + _overrides_block(entries))
    return path, {name: chosen[name][0].name for name in chosen}


# --------------------------------------------------------------------------
# CLI
# --------------------------------------------------------------------------

def _warn(msg: str) -> None:
    print(f"jtaf_release: warning: {msg}", file=sys.stderr)


def cmd_latest(args) -> int:
    tree = open_tree(args)
    releases = discover(tree, warn=None if args.quiet else _warn)
    junos = latest(releases, evo=False)
    evo = latest(releases, evo=True)
    if args.json:
        print(json.dumps({"junos": junos.name if junos else None,
                          "evo": evo.name if evo else None,
                          "commit": tree.commit(),
                          "source": tree.label}, indent=2))
        return 0
    if args.branch is None:
        print(f"junos:  {junos.name if junos else '(none)'}")
        print(f"evo:    {evo.name if evo else '(none)'}")
        print(f"commit: {tree.commit()}")
        print(f"source: {tree.label}")
        return 0
    want = evo if args.branch == "evo" else junos
    if want is None:
        raise ReleaseError(f"no {args.branch} release found in {tree.label}")
    print(want.name)
    return 0


def cmd_commit(args) -> int:
    print(open_tree(args).commit())
    return 0


def cmd_resolve(args) -> int:
    layout = resolve(open_tree(args), args.release)
    if args.field:
        value = layout.as_dict()[args.field]
        print(" ".join(value) if isinstance(value, dict) else value)
    else:
        print(json.dumps(layout.as_dict(), indent=2))
    return 0


def cmd_families(args) -> int:
    layout = resolve(open_tree(args), args.release)
    if args.json:
        print(json.dumps(layout.families, indent=2))
    else:
        print(" ".join(layout.families))
    return 0


def cmd_makevars(args) -> int:
    layout = resolve(open_tree(args), args.release)
    sys.stdout.write(make_vars(layout, args.prefix))
    return 0


def cmd_export(args) -> int:
    """Extract a release's common and configuration modules into a directory."""
    tree = open_tree(args)
    layout = resolve(tree, args.release)
    paths = [layout.common] + sorted(set(layout.families.values()))
    dest = Path(args.dest)
    dest.mkdir(parents=True, exist_ok=True)
    if isinstance(tree, GitTree):
        tar = tree.archive(layout.tree, paths)
        subprocess.run(["tar", "-x", "-C", str(dest)], input=tar, check=True)
    else:
        import shutil
        for rel in paths:
            src = Path(tree.root) / layout.tree / rel
            target = dest / rel
            if target.exists():
                shutil.rmtree(target)
            shutil.copytree(src, target)
    count = sum(1 for _ in dest.rglob("*.yang"))
    print(f"exported {layout.release}: {count} modules from {len(paths)} "
          f"directories into {dest}")
    return 0


def cmd_tooling_id(args) -> int:
    root = Path(args.root) if args.root else None
    if args.list:
        for rel, p in tooling_files(root):
            print(f"{_sha256(p)[:12]}  {rel}  ({p})")
    print(tooling_id(root))
    return 0


def cmd_manifest(args) -> int:
    manifest = write_manifest(Path(args.build_dir), Path(args.dist_dir),
                              args.source, args.commit, args.tooling_id)
    names = ", ".join(p["name"] for p in manifest["providers"])
    print(f"[{manifest['release']}] manifest: {names}")
    return 0


def cmd_terraformrc(args) -> int:
    if args.combined:
        path, chosen = write_combined_terraformrc(Path(args.dist_dir), args.namespace)
        print(f"terraformrc: {path} ({', '.join(f'{n}={r}' for n, r in sorted(chosen.items()))})")
    else:
        path = write_release_terraformrc(Path(args.dist_dir), args.namespace)
        print(f"[{Path(args.dist_dir).name}] terraformrc: {path}")
    return 0


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    parser.add_argument("--repo", help="Juniper YANG clone (default: $YANG_REPO or ../yang)")
    parser.add_argument("--ref", default="HEAD", help="git reference to read (default HEAD)")
    parser.add_argument("--dir", help="read a plain directory laid out like the repository instead of git")
    parser.add_argument("-q", "--quiet", action="store_true", help="suppress warnings")
    sub = parser.add_subparsers(dest="command", required=True)

    p = sub.add_parser("latest", help="newest release of each branch")
    p.add_argument("--branch", choices=["junos", "evo"], help="print only this branch's release name")
    p.add_argument("--json", action="store_true", help="both branches and the commit as JSON")
    p.set_defaults(func=cmd_latest)

    p = sub.add_parser("commit", help="commit hash the reference resolves to")
    p.set_defaults(func=cmd_commit)

    p = sub.add_parser("resolve", help="where a release's modules are")
    p.add_argument("release")
    p.add_argument("--field", choices=["release", "major", "tree", "layout", "common", "families"])
    p.set_defaults(func=cmd_resolve)

    p = sub.add_parser("families", help="provider types a release builds")
    p.add_argument("release")
    p.add_argument("--json", action="store_true")
    p.set_defaults(func=cmd_families)

    p = sub.add_parser("makevars", help="a release's layout as Makefile variables")
    p.add_argument("release")
    p.add_argument("--prefix", required=True)
    p.set_defaults(func=cmd_makevars)

    p = sub.add_parser("export", help="extract a release's modules into a directory")
    p.add_argument("release")
    p.add_argument("dest")
    p.set_defaults(func=cmd_export)

    p = sub.add_parser("tooling-id", help="identifier of the JTAF tooling in use")
    p.add_argument("--root", help="resolve the packages under this directory instead of as installed")
    p.add_argument("--list", action="store_true", help="also list the files and their hashes")
    p.set_defaults(func=cmd_tooling_id)

    p = sub.add_parser("manifest", help="write dist/<release>/manifest.json from what is built")
    p.add_argument("build_dir", help="the release's build directory (stamps)")
    p.add_argument("dist_dir", help="the release's dist directory (bin/, schemas/)")
    p.add_argument("--source", default="", help="where the modules came from")
    p.add_argument("--commit", default="", help="YANG repository commit")
    p.add_argument("--tooling-id", default="", help="identifier from tooling-id")
    p.set_defaults(func=cmd_manifest)

    p = sub.add_parser("terraformrc", help="write a Terraform CLI configuration with dev_overrides")
    p.add_argument("dist_dir", help="a release's dist directory, or the dist root with --combined")
    p.add_argument("--combined", action="store_true", help="cover every release under dist_dir")
    p.add_argument("--namespace", default="hashicorp", help="provider source namespace (jtaf-xml2tf writes hashicorp)")
    p.set_defaults(func=cmd_terraformrc)

    return parser


def main(argv=None) -> int:
    args = build_parser().parse_args(argv)
    try:
        return args.func(args)
    except ReleaseError as e:
        print(f"jtaf_release: {e}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    sys.exit(main())
