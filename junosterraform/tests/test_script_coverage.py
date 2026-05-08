import importlib.util
import json
import pathlib
import runpy
import sys
from importlib.machinery import SourceFileLoader

import pytest


REPO_ROOT = pathlib.Path(__file__).resolve().parents[2]
JUNOS_DIR = REPO_ROOT / "junosterraform"


def _load_script(script_name: str, module_name: str):
    script_path = JUNOS_DIR / script_name
    loader = SourceFileLoader(module_name, str(script_path))
    spec = importlib.util.spec_from_loader(module_name, loader)
    module = importlib.util.module_from_spec(spec)
    assert spec is not None and spec.loader is not None
    spec.loader.exec_module(module)
    return module


def test_jtaf_provider_helpers(tmp_path):
    mod = _load_script("jtaf-provider", "jtaf_provider_mod")

    tpl = tmp_path / "in.j2"
    tpl.write_text("hello {{ data.name }}")
    out = tmp_path / "out.txt"
    mod.render_template_and_write(str(tpl), str(out), {"name": "world"}, "sample")
    assert out.read_text() == "hello world"

    go_dir = tmp_path / "go"
    go_dir.mkdir()
    (go_dir / "a.go").write_text('import "terraform_provider/netconf"')
    mod.rewrite_import_prefixes(str(go_dir), "terraform_provider", "terraform-provider-junos-qfx")
    assert "terraform-provider-junos-qfx/netconf" in (go_dir / "a.go").read_text()


def test_jtaf_provider_main_smoke(tmp_path, monkeypatch):
    mod = _load_script("jtaf-provider", "jtaf_provider_main_mod")

    schema = tmp_path / "schema.json"
    schema.write_text(json.dumps({"root": {"children": []}}))
    xml = tmp_path / "cfg.xml"
    xml.write_text("<configuration><system/></configuration>")

    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(mod, "filter_json_using_xml", lambda _s, _x: {"root": {"children": []}})

    # Keep Template usage deterministic and lightweight.
    class _FakeTemplate:
        def __init__(self, _src):
            pass

        def render(self, **kwargs):
            if "data" in kwargs and isinstance(kwargs["data"], dict) and "device_type" in kwargs["data"]:
                return "package main\n"
            return "package main\n"

    monkeypatch.setattr(mod, "Template", _FakeTemplate)

    argv = [
        "jtaf-provider",
        "-j",
        str(schema),
        "-x",
        str(xml),
        "-t",
        "qfx",
    ]
    monkeypatch.setattr(sys, "argv", argv)
    mod.main()

    out = tmp_path / "terraform-provider-junos-qfx"
    assert (out / "resource_config_provider.go").exists()
    assert (out / "provider.go").exists()
    assert (out / "go.mod").exists()
    assert (out / "config.go").exists()
    assert (out / "trimmed_schema.json").exists()


def test_xml2tf_helpers_and_main(tmp_path, monkeypatch):
    mod = _load_script("jtaf-xml2tf", "jtaf_xml2tf_mod")

    assert mod.normalize_tag("host-name.v4") == "host_name_v4"
    assert mod.convert_to_hcl({"a": [1, True, "x"]}).startswith("{")

    type_map = mod.build_type_map(
        {
            "name": "configuration",
            "children": [{"name": "system", "type": "container", "children": [{"name": "host-name", "type": "leaf"}]}],
        }
    )
    assert "system" in type_map
    assert "system/host_name" in type_map

    schema = {
        "root": {
            "children": [
                {
                    "name": "configuration",
                    "children": [
                        {
                            "name": "system",
                            "type": "container",
                            "children": [
                                {"name": "host-name", "type": "leaf"},
                                {"name": "services", "type": "container", "children": [{"name": "ssh", "type": "leaf"}]},
                            ],
                        }
                    ],
                }
            ]
        }
    }
    schema_file = tmp_path / "trimmed_schema.json"
    schema_file.write_text(json.dumps(schema))
    xml_file = tmp_path / "leaf1.xml"
    xml_file.write_text(
        "<configuration><system><host-name>leaf1</host-name>"
        "<services><ssh/></services></system></configuration>"
    )

    out = tmp_path / "tf"
    argv = [
        "jtaf-xml2tf",
        "-j",
        str(schema_file),
        "-x",
        str(xml_file),
        "-t",
        "qfx",
        "-d",
        str(out),
    ]
    monkeypatch.setattr(sys, "argv", argv)
    mod.main()

    assert (out / "providers.tf").exists()
    assert (out / "leaf1.tf").exists()
    assert "provider \"junos-qfx\"" in (out / "providers.tf").read_text()


def test_template_filter_module_and_merge_function():
    filters_mod = _load_script("templates/jtaf_filters.py", "jtaf_filters_mod")
    fm = filters_mod.FilterModule()

    filters = fm.filters()
    assert "jtaf_apply_merge_directives" in filters
    assert fm.extract_directive({"_merge_directive": "append"}) == "append"
    assert fm.extract_directive("x") is None

    cleaned = fm.remove_meta_keys({"a": 1, "_merge_directive": "replace", "b": [{"_merge_x": True, "y": 2}]})
    assert "_merge_directive" not in cleaned
    assert cleaned["b"][0]["y"] == 2

    assert filters_mod.jtaf_merge_with_directive([1], [2], "append") == [1, 2]
    assert filters_mod.jtaf_merge_with_directive([1], [2], "prepend") == [2, 1]
    assert filters_mod.jtaf_merge_with_directive({"a": 1}, {"b": 2}, "merge_recursive") == {"a": 1, "b": 2}
    with pytest.raises(Exception):
        filters_mod.jtaf_merge_with_directive(1, 2, "unknown")


class _FakePopen:
    def __init__(self, cmd, stdout=None, stderr=None, stdin=None):
        self.cmd = cmd
        self.returncode = 0
        self._stdout = b"{}"
        self._stderr = b""

    def communicate(self, input=None):
        if self.cmd and self.cmd[0] in {"jtaf-provider", "jtaf-ansible"}:
            return (b"ok", b"")
        return (self._stdout, self._stderr)


class _FailPopen(_FakePopen):
    def __init__(self, cmd, stdout=None, stderr=None, stdin=None):
        super().__init__(cmd, stdout=stdout, stderr=stderr, stdin=stdin)
        self.returncode = 1
        self._stdout = b""


def test_yang2go_and_yang2ansible_scripts(tmp_path, monkeypatch):
    yang_dir = tmp_path / "yang"
    yang_dir.mkdir()
    yang_file = yang_dir / "a.yang"
    yang_file.write_text("module a { namespace \"x\"; prefix x; }")
    xml_file = tmp_path / "cfg.xml"
    xml_file.write_text("<configuration/>")

    import subprocess

    monkeypatch.setattr(subprocess, "Popen", _FakePopen)

    monkeypatch.setattr(
        sys,
        "argv",
        ["jtaf-yang2go", "-p", str(yang_dir), str(yang_file), "-x", str(xml_file), "-t", "qfx"],
    )
    runpy.run_path(str(JUNOS_DIR / "jtaf-yang2go"), run_name="__main__")

    monkeypatch.setattr(
        sys,
        "argv",
        ["jtaf-yang2ansible", "-p", str(yang_dir), str(yang_file), "-x", str(xml_file), "-t", "qfx"],
    )
    runpy.run_path(str(JUNOS_DIR / "jtaf-yang2ansible"), run_name="__main__")

    # failure branch for pyang
    monkeypatch.setattr(subprocess, "Popen", _FailPopen)
    monkeypatch.setattr(sys, "argv", ["jtaf-yang2go", "-p", str(yang_file), "-x", str(xml_file), "-t", "qfx"])
    runpy.run_path(str(JUNOS_DIR / "jtaf-yang2go"), run_name="__main__")
