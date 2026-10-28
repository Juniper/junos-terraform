"""Tests for the release matrix build helpers (Makefile, scripts/release/).

Nothing here needs the Juniper YANG clone, Docker, Go or the JTAF generators:
the repository layouts are built as fixtures on disk, the container runtime is
a stub on PATH, and the Makefile is only ever dry-run.
"""

import importlib.util
import re
import shutil
import pathlib
import subprocess
import sys

import pytest

REPO_ROOT = pathlib.Path(__file__).resolve().parents[2]
RELEASE_DIR = REPO_ROOT / "scripts" / "release"


def _load_release_module():
    spec = importlib.util.spec_from_file_location(
        "jtaf_release", RELEASE_DIR / "jtaf_release.py")
    module = importlib.util.module_from_spec(spec)
    # dataclasses resolves the defining module through sys.modules.
    sys.modules[spec.name] = module
    spec.loader.exec_module(module)
    return module


jr = _load_release_module()


# --------------------------------------------------------------------------
# Fixtures: the repository's two layouts, as directories of empty modules.
# --------------------------------------------------------------------------

def _touch(path: pathlib.Path, *names: str) -> None:
    path.mkdir(parents=True, exist_ok=True)
    for name in names:
        (path / name).write_text("module x { }\n")


def make_native_junos(root: pathlib.Path, release: str = "25.4R1") -> pathlib.Path:
    major = release.split("R")[0]
    tree = root / major / release / "native" / "conf-and-rpcs"
    _touch(tree / "common" / "models",
           "junos-common-types@2025-01-01.yang", "junos-configuration-metadata.yang")
    for fam in ("junos", "junos-es", "junos-ex", "junos-nfx", "junos-qfx"):
        _touch(tree / fam / "conf" / "models", f"{fam}-conf-root@2025-01-01.yang")
        _touch(tree / fam / "conf-with-extensions" / "models", f"{fam}-conf-root@2025-01-01.yang")
        _touch(tree / fam / "rpc" / "models", f"{fam}-rpc-system@2025-01-01.yang")
        _touch(tree / fam / "genstate" / "models", f"{fam}-state-system@2025-01-01.yang")
    (tree / "README.txt").write_text("Juniper YANG modules\n")
    _touch(root / major / release / "openconfig" / "models", "openconfig-interfaces.yang")
    _touch(root / major / release / "ietf", "ietf-yang-types.yang")
    return tree


def make_native_evo(root: pathlib.Path, release: str = "25.4R1-EVO") -> pathlib.Path:
    major = release.split("R")[0]
    tree = root / major / release / "native" / "conf-and-rpcs"
    _touch(tree / "common" / "models", "junos-common-types@2025-01-01.yang")
    _touch(tree / "conf" / "models", "junos-conf-root@2025-01-01.yang",
           "junos-conf-system@2025-01-01.yang")
    _touch(tree / "conf-with-extensions" / "models", "junos-conf-root@2025-01-01.yang")
    _touch(tree / "rpc" / "models", "junos-rpc-system@2025-01-01.yang")
    _touch(tree / "genstate" / "models", "junos-state-system@2025-01-01.yang")
    return tree


def make_flat_junos(root: pathlib.Path, release: str = "22.4R1") -> pathlib.Path:
    major = release.split("R")[0]
    tree = root / major / release
    _touch(tree / "common", "junos-common-types@2022-01-01.yang")
    for fam in ("junos", "junos-es", "junos-ex", "junos-nfx", "junos-qfx"):
        _touch(tree / fam / "conf", f"{fam}-conf-root@2022-01-01.yang")
        _touch(tree / fam / "rpc", f"{fam}-rpc-system@2022-01-01.yang")
    (tree / "README.txt").write_text("Juniper YANG modules\n")
    return tree


@pytest.fixture
def yang_dir(tmp_path):
    root = tmp_path / "yang"
    make_flat_junos(root, "22.4R1")
    make_native_junos(root, "24.4R1")
    make_native_junos(root, "24.4R2")
    make_native_junos(root, "25.2R1")
    make_native_junos(root, "25.4R1")
    make_native_evo(root, "24.4R1-EVO")
    make_native_evo(root, "24.4R2-EVO")
    make_native_evo(root, "25.2R1-EVO")
    make_native_evo(root, "25.4R1-EVO")
    # Things the real repository has that are not releases.
    _touch(root / "14.2", "configuration.yang")
    _touch(root / "16.1" / "operational", "clear.yang")
    (root / "README.md").write_text("# yang\n")
    return root


def run_helper(*args, check=True):
    return subprocess.run(
        [sys.executable, str(RELEASE_DIR / "jtaf_release.py"), *args],
        capture_output=True, text=True, check=check)


# --------------------------------------------------------------------------
# Release names and ordering
# --------------------------------------------------------------------------

def test_parse_release_forms():
    assert jr.parse_release("25.4R1").key == (25, 4, 1, 0, 0)
    assert jr.parse_release("25.4R1").evo is False
    assert jr.parse_release("25.4R1-EVO").evo is True
    assert jr.parse_release("24.2R1.17").key == (24, 2, 1, 0, 17)
    assert jr.parse_release("23.4R2-S1").key == (23, 4, 2, 1, 0)
    assert jr.parse_release("25.4R1").major == "25.4"
    assert jr.parse_release("operational") is None
    assert jr.parse_release("25.4") is None


def test_latest_selects_newest_per_branch():
    names = ["24.4R1", "24.4R2", "25.2R1", "25.4R1",
             "24.4R1-EVO", "24.4R2-EVO", "25.2R1-EVO", "25.4R1-EVO", "24.2R1.17"]
    releases = [jr.parse_release(n) for n in names]
    assert jr.latest(releases, evo=False).name == "25.4R1"
    assert jr.latest(releases, evo=True).name == "25.4R1-EVO"


def test_r_number_orders_within_a_version():
    r1, r2 = jr.parse_release("24.4R1"), jr.parse_release("24.4R2")
    assert r1 < r2
    assert jr.latest([r2, r1], evo=False).name == "24.4R2"
    assert sorted([r2, r1])[-1].name == "24.4R2"


def test_discover_skips_non_release_directories(yang_dir):
    warnings = []
    releases = jr.discover(jr.DirTree(str(yang_dir)), warn=warnings.append)
    assert [r.name for r in releases if not r.evo] == [
        "22.4R1", "24.4R1", "24.4R2", "25.2R1", "25.4R1"]
    assert [r.name for r in releases if r.evo] == [
        "24.4R1-EVO", "24.4R2-EVO", "25.2R1-EVO", "25.4R1-EVO"]
    assert warnings == ["skipping 16.1/operational: not a release name"]


def test_latest_cli(yang_dir):
    assert run_helper("--dir", str(yang_dir), "latest", "--branch", "junos").stdout.strip() == "25.4R1"
    assert run_helper("--dir", str(yang_dir), "latest", "--branch", "evo").stdout.strip() == "25.4R1-EVO"
    out = run_helper("--dir", str(yang_dir), "latest", "--json")
    assert '"junos": "25.4R1"' in out.stdout
    assert '"evo": "25.4R1-EVO"' in out.stdout
    assert '"commit": "local"' in out.stdout
    assert "warning: skipping 16.1/operational" in out.stderr


# --------------------------------------------------------------------------
# Layout resolution and families
# --------------------------------------------------------------------------

def test_resolve_native_junos(yang_dir):
    layout = jr.resolve(jr.DirTree(str(yang_dir)), "25.4R1")
    assert layout.tree == "25.4/25.4R1/native/conf-and-rpcs"
    assert layout.layout == "native"
    assert layout.common == "common/models"
    assert layout.families == {
        "mx": "junos/conf/models",
        "es": "junos-es/conf/models",
        "ex": "junos-ex/conf/models",
        "nfx": "junos-nfx/conf/models",
        "qfx": "junos-qfx/conf/models",
    }
    assert layout.family_dir("mx") == "junos"


def test_resolve_native_evo(yang_dir):
    layout = jr.resolve(jr.DirTree(str(yang_dir)), "25.4R1-EVO")
    assert layout.tree == "25.4/25.4R1-EVO/native/conf-and-rpcs"
    assert layout.common == "common/models"
    assert layout.families == {"evo": "conf/models"}


def test_resolve_flat_layout(yang_dir):
    layout = jr.resolve(jr.DirTree(str(yang_dir)), "22.4R1")
    assert layout.tree == "22.4/22.4R1"
    assert layout.layout == "flat"
    assert layout.common == "common"
    assert layout.families["qfx"] == "junos-qfx/conf"
    assert list(layout.families) == ["mx", "es", "ex", "nfx", "qfx"]


def test_resolve_unknown_release_names_release_and_source(yang_dir):
    with pytest.raises(jr.ReleaseError) as e:
        jr.resolve(jr.DirTree(str(yang_dir)), "99.9R9")
    assert "99.9R9" in str(e.value)
    assert str(yang_dir) in str(e.value)

    out = run_helper("--dir", str(yang_dir), "resolve", "99.9R9", check=False)
    assert out.returncode == 2
    assert "99.9R9" in out.stderr and str(yang_dir) in out.stderr


def test_families_cli(yang_dir):
    assert run_helper("--dir", str(yang_dir), "families", "25.4R1").stdout.split() == [
        "mx", "es", "ex", "nfx", "qfx"]
    assert run_helper("--dir", str(yang_dir), "families", "25.4R1-EVO").stdout.split() == ["evo"]


def test_makevars_output(yang_dir):
    out = run_helper("--dir", str(yang_dir), "makevars", "25.4R1", "--prefix", "JUNOS").stdout
    lines = out.splitlines()
    assert "JUNOS_RELEASE := 25.4R1" in lines
    assert "JUNOS_TREE := 25.4/25.4R1/native/conf-and-rpcs" in lines
    assert "JUNOS_COMMON := common/models" in lines
    assert "JUNOS_TYPES := mx es ex nfx qfx" in lines
    assert "JUNOS_CONF_qfx := junos-qfx/conf/models" in lines
    assert "JUNOS_CONF_mx := junos/conf/models" in lines

    out = run_helper("--dir", str(yang_dir), "makevars", "25.4R1-EVO", "--prefix", "EVO").stdout
    assert "EVO_TYPES := evo" in out.splitlines()
    assert "EVO_CONF_evo := conf/models" in out.splitlines()


def test_export_from_directory_copies_only_conf_and_common(yang_dir, tmp_path):
    dest = tmp_path / "exported"
    out = run_helper("--dir", str(yang_dir), "export", "25.4R1", str(dest))
    assert "exported 25.4R1" in out.stdout
    assert sorted(p.name for p in dest.iterdir()) == [
        "common", "junos", "junos-es", "junos-ex", "junos-nfx", "junos-qfx"]
    assert (dest / "junos-qfx" / "conf" / "models").is_dir()
    assert not (dest / "junos-qfx" / "rpc").exists()
    assert not (dest / "junos-qfx" / "conf-with-extensions").exists()


# --------------------------------------------------------------------------
# Tooling identity
# --------------------------------------------------------------------------

@pytest.fixture
def tooling_root(tmp_path):
    root = tmp_path / "tooling"
    (root / "junosterraform").mkdir(parents=True)
    (root / "junosterraform" / "jtaf-provider").write_text("print('provider')\n")
    (root / "junosterraform" / "jtaf-yang2go").write_text("print('yang2go')\n")
    (root / "junosterraform" / "jtaf_common.py").write_text("X = 1\n")
    (root / "jtaf_pyang_plugin").mkdir()
    (root / "jtaf_pyang_plugin" / "jtaf_json.py").write_text("Y = 2\n")
    (root / "jtaf_pyang_plugin" / "__init__.py").write_text("")
    tp = root / "terraform_provider"
    (tp / "generic").mkdir(parents=True)
    (tp / "tools").mkdir()
    (tp / "main.go").write_text("package main\n")
    (tp / "generic" / "serve.go").write_text("package generic\n")
    (tp / "generic" / "serve_test.go").write_text("package generic // test\n")
    (tp / "tools" / "lint.go").write_text("package tools\n")
    (tp / "go.mod").write_text("module terraform_provider\n")
    (tp / "go.sum").write_text("")
    return root


def test_tooling_id_is_stable_and_tracks_source_edits(tooling_root):
    first = jr.tooling_id(tooling_root)
    assert len(first) == 12
    assert jr.tooling_id(tooling_root) == first

    # Tests and the tools directory are not part of what a provider is built from.
    (tooling_root / "terraform_provider" / "generic" / "serve_test.go").write_text("changed\n")
    (tooling_root / "terraform_provider" / "tools" / "lint.go").write_text("changed\n")
    assert jr.tooling_id(tooling_root) == first

    (tooling_root / "terraform_provider" / "generic" / "serve.go").write_text("package generic // v2\n")
    second = jr.tooling_id(tooling_root)
    assert second != first

    (tooling_root / "junosterraform" / "jtaf_common.py").write_text("X = 2\n")
    assert jr.tooling_id(tooling_root) not in (first, second)


def test_tooling_id_cli(tooling_root):
    out = run_helper("tooling-id", "--root", str(tooling_root)).stdout.strip()
    assert out == jr.tooling_id(tooling_root)


def test_tooling_id_covers_the_installed_generators():
    names = [rel for rel, _ in jr.tooling_files()]
    assert "junosterraform/jtaf-provider" in names
    assert "junosterraform/jtaf-yang2go" in names
    assert "jtaf_pyang_plugin/jtaf_json.py" in names
    assert any(n.startswith("terraform_provider/") and n.endswith(".go") for n in names)
    assert not any(n.endswith("_test.go") for n in names)


# --------------------------------------------------------------------------
# Manifest and Terraform CLI configuration
# --------------------------------------------------------------------------

def _gz_schema(path: pathlib.Path, fingerprint: str) -> None:
    import gzip
    import json
    payload = {"configuration": {"children": {"fingerprint": {"type": "string"}}},
               "fingerprint": fingerprint}
    with gzip.open(path, "wb") as f:
        f.write(json.dumps(payload).encode())


def make_built_release(root: pathlib.Path, release: str, types, fingerprint="full-groups-abc123def456"):
    """A release's build and dist directories as the Makefile leaves them."""
    build = root / "build" / release
    dist = root / "dist" / release
    (dist / "bin").mkdir(parents=True)
    (dist / "schemas").mkdir()
    build.mkdir(parents=True)
    for t in types:
        (dist / "bin" / f"terraform-provider-junos-{t}").write_bytes(f"ELF {release} {t}".encode())
        _gz_schema(dist / "schemas" / f"junos-{t}.trimmed_schema.json.gz", fingerprint)
        conf = "conf/models" if t == "evo" else ("junos/conf/models" if t == "mx" else f"junos-{t}/conf/models")
        (build / f"inputs-{t}.txt").write_text(
            f"release={release}\ncommit=96ad7bad\ntype={t}\nmodules=common/models {conf}\n"
            f"exclude=tenants logical-systems\nflags=\ntooling=d6c723e70095\n")
    return build, dist


def test_manifest_contents(tmp_path):
    import hashlib
    import json
    build, dist = make_built_release(tmp_path, "25.4R1", ["qfx", "mx"])
    manifest = jr.write_manifest(build, dist, source="git:/x/yang@HEAD", commit="96ad7bad", tooling="d6c723e70095")

    on_disk = json.loads((dist / "manifest.json").read_text())
    assert on_disk == manifest
    assert manifest["release"] == "25.4R1"
    assert manifest["branch"] == "junos"
    assert manifest["commit"] == "96ad7bad"
    assert manifest["source"] == "git:/x/yang@HEAD"
    assert manifest["tooling_id"] == "d6c723e70095"
    assert manifest["built_at"].endswith("Z")
    assert [p["type"] for p in manifest["providers"]] == ["mx", "qfx"]

    qfx = manifest["providers"][1]
    assert qfx["name"] == "junos-qfx"
    assert qfx["family"] == "junos-qfx"
    assert qfx["binary"] == "bin/terraform-provider-junos-qfx"
    assert qfx["sha256"] == hashlib.sha256(b"ELF 25.4R1 qfx").hexdigest()
    assert qfx["schema"] == "schemas/junos-qfx.trimmed_schema.json.gz"
    assert qfx["fingerprint"] == "full-groups-abc123def456"
    assert qfx["exclude"] == ["tenants", "logical-systems"]
    assert qfx["modules"] == ["common/models", "junos-qfx/conf/models"]
    assert manifest["providers"][0]["family"] == "junos"


def test_manifest_evo_and_cli(tmp_path):
    import json
    build, dist = make_built_release(tmp_path, "25.4R1-EVO", ["evo"], fingerprint="full-groups-evo")
    out = run_helper("manifest", str(build), str(dist), "--source", "s", "--commit", "c", "--tooling-id", "t")
    assert "[25.4R1-EVO] manifest: junos-evo" in out.stdout
    manifest = json.loads((dist / "manifest.json").read_text())
    assert manifest["branch"] == "evo"
    assert manifest["providers"][0]["family"] == "evo"
    assert manifest["providers"][0]["fingerprint"] == "full-groups-evo"


def test_manifest_without_binaries_fails(tmp_path):
    (tmp_path / "dist" / "25.4R1" / "bin").mkdir(parents=True)
    out = run_helper("manifest", str(tmp_path / "build"), str(tmp_path / "dist" / "25.4R1"), check=False)
    assert out.returncode == 2
    assert "no provider binaries" in out.stderr


def test_release_terraformrc(tmp_path):
    build, dist = make_built_release(tmp_path, "25.4R1", ["qfx", "es"])
    jr.write_manifest(build, dist, "s", "c", "t")
    path = jr.write_release_terraformrc(dist, "hashicorp")
    text = path.read_text()
    assert path == dist / "terraformrc"
    assert "provider_installation {" in text and "dev_overrides {" in text and "direct {}" in text
    bin_dir = str((dist / "bin").resolve())
    assert '"registry.terraform.io/hashicorp/junos-qfx"' in text
    assert '"registry.terraform.io/hashicorp/junos-es"' in text
    assert text.count(bin_dir) == 2
    assert "# Junos release 25.4R1" in text


def test_combined_terraformrc_newest_release_wins(tmp_path):
    for release in ("25.2R1", "25.4R1"):
        build, dist = make_built_release(tmp_path, release, ["qfx", "es"])
        jr.write_manifest(build, dist, "s", "c", "t")
    build, dist = make_built_release(tmp_path, "25.4R1-EVO", ["evo"])
    jr.write_manifest(build, dist, "s", "c", "t")

    path, chosen = jr.write_combined_terraformrc(tmp_path / "dist", "hashicorp")
    text = path.read_text()
    assert path == tmp_path / "dist" / "terraformrc"
    assert chosen == {"junos-qfx": "25.4R1", "junos-es": "25.4R1", "junos-evo": "25.4R1-EVO"}
    newest = str((tmp_path / "dist" / "25.4R1" / "bin").resolve())
    older = str((tmp_path / "dist" / "25.2R1" / "bin").resolve())
    assert f'"registry.terraform.io/hashicorp/junos-qfx" = "{newest}"' in text.replace("  ", " ").replace("  ", " ") \
        or (newest in text and '"registry.terraform.io/hashicorp/junos-qfx"' in text)
    assert older not in text.split("dev_overrides")[1]
    assert "# junos-qfx: 25.4R1  (also built: 25.2R1)" in text
    assert "# junos-evo: 25.4R1-EVO" in text
    assert text.split("dev_overrides")[1].count("registry.terraform.io/") == 3


def test_combined_terraformrc_cli(tmp_path):
    build, dist = make_built_release(tmp_path, "25.4R1", ["qfx"])
    jr.write_manifest(build, dist, "s", "c", "t")
    out = run_helper("terraformrc", "--combined", str(tmp_path / "dist"), "--namespace", "juniper")
    assert "junos-qfx=25.4R1" in out.stdout
    assert '"registry.terraform.io/juniper/junos-qfx"' in (tmp_path / "dist" / "terraformrc").read_text()


# --------------------------------------------------------------------------
# Preflight
# --------------------------------------------------------------------------

def run_preflight(env, tmp_path):
    full_env = {"PATH": "/usr/bin:/bin", "HOME": str(tmp_path), **env}
    return subprocess.run(["bash", str(RELEASE_DIR / "preflight.sh")],
                          capture_output=True, text=True, env=full_env, cwd=str(tmp_path))


def test_preflight_names_missing_tool(tmp_path):
    out = run_preflight({"REQUIRED_TOOLS": "bash definitely-not-a-tool-xyz", "BUILD_DIR": str(tmp_path / "b")}, tmp_path)
    assert out.returncode == 1
    assert "missing tool(s): definitely-not-a-tool-xyz" in out.stderr
    assert "bash" not in out.stderr.split("missing tool(s):")[1].splitlines()[0].replace("definitely-not-a-tool-xyz", "")


def test_preflight_disk_shortfall_names_amounts_and_cleanup(tmp_path):
    out = run_preflight({"REQUIRED_TOOLS": "bash", "MIN_FREE_GB": "100000",
                         "BUILD_DIR": str(tmp_path / "b")}, tmp_path)
    assert out.returncode == 1
    assert "need at least 100000 GB" in out.stderr
    assert "GB free on the filesystem holding" in out.stderr
    assert "make clean-space" in out.stderr
    assert "AUTO_CLEAN=1" in out.stderr


def test_preflight_ok_and_memory_warning(tmp_path):
    out = run_preflight({"REQUIRED_TOOLS": "bash", "MIN_FREE_GB": "0", "BUILD_DIR": str(tmp_path / "b"),
                         "JOBS": "100000"}, tmp_path)
    assert out.returncode == 0, out.stderr
    assert out.stdout.startswith("preflight: ok")
    assert "warning" in out.stderr and "JOBS=1" in out.stderr
    assert (tmp_path / "b" / "tmp").is_dir()


def test_preflight_auto_clean_runs_cleanup_once_then_rechecks(tmp_path):
    # A cleanup stub that records its invocation and makes the next df check
    # pass by lowering the threshold it reads back.
    stub = tmp_path / "clean.sh"
    stub.write_text("#!/bin/bash\necho cleaned >> \"$HOME/clean.log\"\n")
    stub.chmod(0o755)
    # First check: impossible threshold. To let the recheck pass, the stub
    # cannot change MIN_FREE_GB; instead run with a threshold just above
    # what is free so the first check fails and verify the stub ran once.
    free_gb = int(subprocess.run(["df", "-Pk", str(tmp_path)], capture_output=True, text=True)
                  .stdout.splitlines()[1].split()[3]) // (1024 * 1024)
    out = run_preflight({"REQUIRED_TOOLS": "bash", "MIN_FREE_GB": str(free_gb + 1000), "AUTO_CLEAN": "1",
                         "BUILD_DIR": str(tmp_path / "b"), "CLEAN_SPACE_SCRIPT": str(stub)}, tmp_path)
    assert out.returncode == 1
    assert (tmp_path / "clean.log").read_text() == "cleaned\n"
    assert "AUTO_CLEAN=1, running clean-space" in out.stdout
    assert "need at least" in out.stderr


# --------------------------------------------------------------------------
# Cleanup, against a stub container runtime
# --------------------------------------------------------------------------

STUB_DOCKER = r'''#!/bin/bash
log="$STUB_LOG"
echo "$*" >> "$log"
case "$1 $2" in
  "info ") exit 0 ;;
  "ps -q") echo "aaa111" ;;
  "ps -a") printf '%s\n' "bbb222  old-builder  Exited (0) 5 days ago  (builder:1)" \
                        "ccc333  old-docs  Exited (0) 4 days ago  (docs:2)" ;;
  "images -f") echo "ddd444  <none>:<none>  120MB" ;;
  "images --format") echo "eee555  junos-sim:latest  2GB  3 days ago" ;;
  "system df") echo "TYPE TOTAL ACTIVE SIZE RECLAIMABLE" ;;
  *prune*) echo "stub: prune must not run in a dry run" >&2; exit 99 ;;
esac
'''


def run_clean(env, tmp_path, with_docker=True):
    bin_dir = tmp_path / "stubbin"
    bin_dir.mkdir(exist_ok=True)
    if with_docker:
        docker = bin_dir / "docker"
        docker.write_text(STUB_DOCKER)
        docker.chmod(0o755)
    log = tmp_path / "docker.log"
    log.write_text("")
    full_env = {"PATH": f"{bin_dir}:/usr/bin:/bin", "HOME": str(tmp_path), "STUB_LOG": str(log), **env}
    out = subprocess.run(["bash", str(RELEASE_DIR / "clean_space.sh")],
                         capture_output=True, text=True, env=full_env, cwd=str(tmp_path))
    return out, log.read_text()


def test_clean_space_dry_run_names_exited_not_running(tmp_path):
    build = tmp_path / "build"
    (build / "tmp").mkdir(parents=True)
    (build / "tmp" / "junk").write_text("x")
    (build / "25.2R1").mkdir()
    (build / "25.4R1").mkdir()
    out, log = run_clean({"DRY_RUN": "1", "BUILD_DIR": str(build), "KEEP_RELEASES": "25.4R1 25.4R1-EVO"}, tmp_path)
    assert out.returncode == 0, out.stderr

    assert "old-builder" in out.stdout and "old-docs" in out.stdout
    assert "aaa111" not in out.stdout
    assert "1 running container(s) will not be touched" in out.stdout
    assert "ddd444" in out.stdout
    assert "docker container prune -f --filter until=24h" in out.stdout
    assert "docker image prune -f" in out.stdout
    assert "docker builder prune -f --filter until=24h" in out.stdout
    assert "image prune -a" not in out.stdout
    assert "go clean" not in out.stdout
    assert "prune" not in log
    assert f"would remove stale build directory {build}/25.2R1/" in out.stdout
    assert f"keeping {build}/25.4R1/" in out.stdout
    assert (build / "25.2R1").is_dir() and (build / "tmp" / "junk").is_file()
    assert "dry run complete; nothing removed" in out.stdout


def test_clean_space_opt_ins_are_listed_only_when_asked(tmp_path):
    build = tmp_path / "build"
    build.mkdir()
    env = {"DRY_RUN": "1", "BUILD_DIR": str(build), "PRUNE_UNUSED_IMAGES": "1", "CLEAN_GO_CACHE": "1",
           "PRUNE_IMAGE_AGE": "72h"}
    # go is not on the stub PATH, so the Go cache step reports it cannot run.
    out, log = run_clean(env, tmp_path)
    assert out.returncode == 0, out.stderr
    assert "docker image prune -a -f --filter until=72h" in out.stdout
    assert "junos-sim:latest" in out.stdout
    assert "CLEAN_GO_CACHE=1 but go is not available" in out.stdout
    assert "modcache" not in out.stdout
    assert "prune" not in log


def test_clean_space_without_runtime_still_cleans_local(tmp_path):
    build = tmp_path / "build"
    (build / "tmp").mkdir(parents=True)
    (build / "tmp" / "junk").write_text("x")
    (build / "old").mkdir()
    out, _ = run_clean({"BUILD_DIR": str(build), "KEEP_RELEASES": "25.4R1", "DOCKER": "no-such-runtime"},
                       tmp_path, with_docker=False)
    assert out.returncode == 0, out.stderr
    assert "container cleanup skipped: no-such-runtime is not available" in out.stdout
    assert not (build / "tmp" / "junk").exists()
    assert not (build / "old").exists()
    assert "reclaimed" in out.stdout


def test_clean_space_keeps_all_release_dirs_when_none_selected(tmp_path):
    build = tmp_path / "build"
    (build / "25.2R1").mkdir(parents=True)
    out, _ = run_clean({"BUILD_DIR": str(build), "KEEP_RELEASES": "", "DOCKER": "no-such-runtime"},
                       tmp_path, with_docker=False)
    assert out.returncode == 0, out.stderr
    assert "keeping every release build directory" in out.stdout
    assert (build / "25.2R1").is_dir()


# --------------------------------------------------------------------------
# The Makefile, dry-run against a fixture repository
# --------------------------------------------------------------------------

@pytest.mark.skipif(shutil.which("make") is None, reason="GNU make is not installed")
def test_make_dry_run_covers_the_matrix(tmp_path):
    root = tmp_path / "yang"
    make_native_junos(root, "25.4R1")
    make_native_evo(root, "25.4R1-EVO")
    out = subprocess.run(
        ["make", "-n", "all", f"YANG_DIR={root}", "JUNOS_RELEASE=25.4R1", "EVO_RELEASE=25.4R1-EVO",
         "YANG_COMMIT=0000000", "TOOLING_ID=test", f"BUILD_DIR={tmp_path}/build", f"DIST_DIR={tmp_path}/dist"],
        capture_output=True, text=True, cwd=str(REPO_ROOT))
    assert out.returncode == 0, out.stderr
    types = sorted(re.findall(r"jtaf-yang2go --generic --groups .*? -t (\w+)", out.stdout))
    assert types == ["es", "evo", "ex", "mx", "nfx", "qfx"]
    assert out.stdout.count("nice -n 10") >= 12          # six generations, six compilations
    assert "preflight.sh" in out.stdout
    assert f"manifest {tmp_path}/build/25.4R1 {tmp_path}/dist/25.4R1" in out.stdout
    assert f"manifest {tmp_path}/build/25.4R1-EVO {tmp_path}/dist/25.4R1-EVO" in out.stdout
    assert f"terraformrc --combined {tmp_path}/dist" in out.stdout
    # Only NFX carries exclusions by default (see test_make_per_family_exclusions).
    excluding = [line for line in out.stdout.splitlines() if "--exclude" in line and "jtaf-yang2go" in line]
    assert len(excluding) == 1 and "-t nfx " in excluding[0]


@pytest.mark.skipif(shutil.which("make") is None, reason="GNU make is not installed")
def test_make_dry_run_single_family_and_exclude(tmp_path):
    root = tmp_path / "yang"
    make_native_junos(root, "25.4R1")
    common = [f"YANG_DIR={root}", "JUNOS_RELEASE=25.4R1", "EVO_RELEASE=", "YANG_COMMIT=0000000",
              "TOOLING_ID=test", f"BUILD_DIR={tmp_path}/build", f"DIST_DIR={tmp_path}/dist"]
    out = subprocess.run(["make", "-n", "qfx", "EXCLUDE=logical-systems tenants", *common],
                         capture_output=True, text=True, cwd=str(REPO_ROOT))
    assert out.returncode == 0, out.stderr
    assert re.findall(r"-t (\w+)", out.stdout) == ["qfx"]
    assert "--exclude logical-systems --exclude tenants" in out.stdout

    out = subprocess.run(["make", "-n", "junos", "FAMILIES=mx es", *common],
                         capture_output=True, text=True, cwd=str(REPO_ROOT))
    assert out.returncode == 0, out.stderr
    assert sorted(re.findall(r"-t (\w+)", out.stdout)) == ["es", "mx"]

    out = subprocess.run(["make", "-n", "junos", "FAMILIES=mx srx", *common],
                         capture_output=True, text=True, cwd=str(REPO_ROOT))
    assert out.returncode != 0
    assert "FAMILIES has srx" in out.stderr


@pytest.mark.skipif(shutil.which("make") is None, reason="GNU make is not installed")
def test_make_per_family_exclusions(tmp_path):
    root = tmp_path / "yang"
    make_native_junos(root, "25.4R1")
    common = [f"YANG_DIR={root}", "JUNOS_RELEASE=25.4R1", "EVO_RELEASE=", "YANG_COMMIT=0000000",
              "TOOLING_ID=test", f"BUILD_DIR={tmp_path}/build", f"DIST_DIR={tmp_path}/dist"]

    def generation(out, type_):
        return next(line for line in out.stdout.splitlines()
                    if "jtaf-yang2go" in line and f"-t {type_} " in line)

    # NFX carries the hash-function exclusions by default, in the base
    # hierarchy and under groups; the other families carry none.
    out = subprocess.run(["make", "-n", "junos", *common], capture_output=True, text=True, cwd=str(REPO_ROOT))
    assert out.returncode == 0, out.stderr
    nfx = generation(out, "nfx")
    excluded = re.findall(r"--exclude (\S+)", nfx)
    assert len(excluded) == 8
    assert "forwarding-options/enhanced-hash-key/hash-function" in excluded
    assert "groups/tenants/routing-instances/instance/forwarding-options/enhanced-hash-key/hash-function" in excluded
    assert all(p.endswith("/hash-function") for p in excluded)
    assert "--exclude" not in generation(out, "qfx")

    # Global EXCLUDE applies to every provider and adds to the per-family set;
    # EXCLUDE_<type> can be cleared or replaced.
    out = subprocess.run(["make", "-n", "junos", "EXCLUDE=tenants", "EXCLUDE_qfx=poe", "EXCLUDE_nfx=", *common],
                         capture_output=True, text=True, cwd=str(REPO_ROOT))
    assert out.returncode == 0, out.stderr
    assert re.findall(r"--exclude (\S+)", generation(out, "nfx")) == ["tenants"]
    assert re.findall(r"--exclude (\S+)", generation(out, "qfx")) == ["tenants", "poe"]
    assert re.findall(r"--exclude (\S+)", generation(out, "mx")) == ["tenants"]
    # The recorded inputs carry the effective exclusions, so the manifest can.
    assert "'exclude=tenants poe'" in out.stdout


@pytest.mark.skipif(shutil.which("make") is None, reason="GNU make is not installed")
def test_make_utility_goals_do_not_resolve_releases(tmp_path):
    # A git that fails loudly: utility goals must never reach it.
    bin_dir = tmp_path / "bin"
    bin_dir.mkdir()
    (bin_dir / "git").write_text("#!/bin/sh\necho 'git must not be called' >&2\nexit 1\n")
    (bin_dir / "git").chmod(0o755)
    env = {**dict(PATH=f"{bin_dir}:/usr/bin:/bin"), "HOME": str(tmp_path)}
    for goal in ("help", "-n latest JUNOS_RELEASE=x EVO_RELEASE=y"):
        out = subprocess.run(["make", *goal.split()], capture_output=True, text=True,
                             cwd=str(REPO_ROOT), env=env)
        assert out.returncode == 0, (goal, out.stderr)
        assert "git must not be called" not in out.stderr
    assert "Targets:" in subprocess.run(["make"], capture_output=True, text=True,
                                        cwd=str(REPO_ROOT), env=env).stdout
