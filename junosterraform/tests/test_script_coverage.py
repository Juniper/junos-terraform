import importlib.util
import json
import pathlib
import runpy
import shutil
import sys
from importlib.machinery import SourceFileLoader

import pytest


REPO_ROOT = pathlib.Path(__file__).resolve().parents[2]
JUNOS_DIR = REPO_ROOT / "junosterraform"


def test_terraform_workflow_uses_shared_apply_parallelism():
    workflow = (REPO_ROOT / ".github" / "workflows" / "go-terraform-provider.yml").read_text()

    assert "TF_APPLY_PARALLELISM: 3" in workflow
    apply_commands = [line.strip() for line in workflow.splitlines() if "terraform apply " in line]
    assert len(apply_commands) == 4
    assert all('-parallelism="$TF_APPLY_PARALLELISM"' in command for command in apply_commands)
    assert "-parallelism=1" not in workflow


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

    go_dir = tmp_path / "go"
    go_dir.mkdir()
    (go_dir / "a.go").write_text('import "terraform_provider/netconf"')
    mod.rewrite_import_prefixes(str(go_dir), "terraform_provider", "terraform-provider-junos-qfx")
    assert "terraform-provider-junos-qfx/netconf" in (go_dir / "a.go").read_text()


def test_jtaf_provider_exclude_schema_paths():
    mod = _load_script("jtaf-provider", "jtaf_provider_exclude_mod")
    resources = {"root": {"children": [{"name": "configuration", "children": [
        {"name": "groups"},
        {"name": "system", "children": [
            {"name": "host-name"},
            {"name": "services", "children": [{"name": "ssh"}, {"name": "web-management"}]},
        ]},
    ]}]}}
    mod.exclude_schema_paths(resources, ["groups", "system/services/web-management"])
    config = resources["root"]["children"][0]
    assert [c["name"] for c in config["children"]] == ["system"]
    services = config["children"][0]["children"][1]
    assert [c["name"] for c in services["children"]] == ["ssh"]
    for bad in ["no-such-section", "system/no-such/ssh", "/"]:
        with pytest.raises(ValueError):
            mod.exclude_schema_paths(resources, [bad])


def test_jtaf_provider_validate_attribute_names():
    mod = _load_script("jtaf-provider", "jtaf_provider_names_mod")

    def schema(*children):
        return {"root": {"children": [{"name": "configuration", "children": list(children)}]}}

    mod.validate_attribute_names(schema(
        {"name": "firewall", "type": "container", "children": [
            {"name": "AH_header", "type": "leaf"}, {"name": "ESP_header", "type": "leaf"},
        ]},
    ))
    with pytest.raises(ValueError, match="both map to attribute 'a_b'"):
        mod.validate_attribute_names(schema({"name": "a-b", "type": "leaf"}, {"name": "a_b", "type": "leaf"}))
    with pytest.raises(ValueError, match="configuration/system/802.1x"):
        mod.validate_attribute_names(schema(
            {"name": "system", "type": "container", "children": [{"name": "802.1x", "type": "leaf"}]},
        ))
    # Nodes in different cases of a choice are siblings once the choice is flattened.
    with pytest.raises(ValueError, match="configuration/system: a-b and a_b"):
        mod.validate_attribute_names(schema(
            {"name": "system", "type": "container", "children": [
                {"name": "c", "type": "choice", "children": [
                    {"name": "x", "type": "case", "children": [{"name": "a-b", "type": "leaf"}]},
                    {"name": "y", "type": "case", "children": [{"name": "a_b", "type": "leaf"}]},
                ]},
            ]},
        ))
    with pytest.raises(ValueError, match="no configuration node"):
        mod.validate_attribute_names({"root": {"children": []}})


def test_jtaf_provider_drop_version():
    mod = _load_script("jtaf-provider", "jtaf_provider_version_mod")
    resources = {"root": {"children": [{"name": "configuration", "children": [
        {"name": "version"},
        {"name": "system", "children": [{"name": "ntp", "children": [{"name": "server", "children": [{"name": "version"}]}]}]},
    ]}]}}
    mod.drop_version(resources)
    config = resources["root"]["children"][0]
    assert [c["name"] for c in config["children"]] == ["system"]
    server = config["children"][0]["children"][0]["children"][0]
    assert [c["name"] for c in server["children"]] == ["version"]
    # A schema without version is left as it is
    mod.drop_version(resources)
    assert [c["name"] for c in config["children"]] == ["system"]


def test_jtaf_provider_main_smoke(tmp_path, monkeypatch):
    mod = _load_script("jtaf-provider", "jtaf_provider_main_mod")

    resources = {"root": {"children": [{"name": "configuration", "children": [
        {"name": "system", "type": "container", "children": [{"name": "host-name", "type": "leaf"}]},
    ]}]}}
    schema = tmp_path / "schema.json"
    schema.write_text(json.dumps(resources))
    xml = tmp_path / "cfg.xml"
    xml.write_text("<configuration><system/></configuration>")

    monkeypatch.chdir(tmp_path)
    monkeypatch.setattr(mod, "filter_json_using_xml", lambda _s, _x: resources)
    # Compiling the schema needs Go; the copy/emit steps are what is under test.
    monkeypatch.setattr(mod, "_compile_schema", lambda *_: None)

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
    assert (out / "main.go").exists()
    assert (out / "embed_schema.go").exists()
    assert (out / "go.mod").exists()
    assert (out / "trimmed_schema.json.gz").exists()
    assert not (out / "trimmed_schema.json").exists()
    # No Go source is rendered from a template any more.
    assert not (out / "resource_config_provider.go").exists()
    from junosterraform.jtaf_common import load_schema_json
    assert load_schema_json(str(out / "trimmed_schema.json.gz")) == resources


def test_trimmed_and_untrimmed_builds_share_provider_source(tmp_path, monkeypatch):
    """-x and --generic choose the schema, not the provider implementation."""
    resources = {"root": {"children": [{"name": "configuration", "children": [
        {"name": "system", "type": "container", "children": [{"name": "host-name", "type": "leaf"}]},
    ]}]}}
    schema = tmp_path / "schema.json"
    schema.write_text(json.dumps(resources))
    xml = tmp_path / "cfg.xml"
    xml.write_text("<configuration><system><host-name>r1</host-name></system></configuration>")

    monkeypatch.chdir(tmp_path)

    for name, argv in (
        ("trimmed", ["-j", str(schema), "-x", str(xml), "-t", "same"]),
        ("untrimmed", ["-j", str(schema), "-t", "same", "--generic"]),
    ):
        mod = _load_script("jtaf-provider", f"jtaf_provider_{name}_mod")
        monkeypatch.setattr(mod, "filter_json_using_xml", lambda _s, _x: resources)
        monkeypatch.setattr(mod, "_compile_schema", lambda *_: None)
        monkeypatch.setattr(sys, "argv", ["jtaf-provider", *argv])
        mod.main()
        shutil.copytree(tmp_path / "terraform-provider-junos-same", tmp_path / name)
        shutil.rmtree(tmp_path / "terraform-provider-junos-same")

    trimmed_go = sorted(p.relative_to(tmp_path / "trimmed")
                        for p in (tmp_path / "trimmed").rglob("*.go"))
    untrimmed_go = sorted(p.relative_to(tmp_path / "untrimmed")
                          for p in (tmp_path / "untrimmed").rglob("*.go"))
    assert trimmed_go == untrimmed_go
    assert trimmed_go, "expected Go source in the generated provider"

    for rel in trimmed_go:
        a = (tmp_path / "trimmed" / rel).read_bytes()
        b = (tmp_path / "untrimmed" / rel).read_bytes()
        assert a == b, f"{rel} differs between the trimmed and untrimmed builds"


def test_jtaf_provider_generic_writes_gzipped_schema(tmp_path, monkeypatch):
    mod = _load_script("jtaf-provider", "jtaf_provider_generic_mod")
    from junosterraform.jtaf_common import load_schema_json

    resources = {"root": {"children": [{"name": "configuration", "children": [
        {"name": "system", "type": "container", "children": [{"name": "host-name", "type": "leaf"}]},
    ]}]}}
    schema = tmp_path / "schema.json"
    schema.write_text(json.dumps(resources))
    monkeypatch.chdir(tmp_path)
    # Compiling the schema needs Go; the copy/emit steps are what is under test.
    monkeypatch.setattr(mod, "_compile_schema", lambda *_: None)
    monkeypatch.setattr(sys, "argv", ["jtaf-provider", "-j", str(schema), "-t", "srx", "--generic"])
    mod.main()

    out = tmp_path / "terraform-provider-junos-srx"
    assert (out / "main.go").exists()
    assert (out / "embed_schema.go").exists()
    assert (out / "trimmed_schema.json.gz").exists()
    assert not (out / "trimmed_schema.json").exists()
    assert load_schema_json(str(out / "trimmed_schema.json.gz")) == resources


def test_xml2tf_build_type_map_flattens_choice_case_for_xml_path():
    mod = _load_script("jtaf-xml2tf", "jtaf_xml2tf_choice_mod")
    schema = {
        "name": "configuration",
        "type": "container",
        "children": [
            {
                "name": "routing-options",
                "type": "container",
                "children": [
                    {
                        "name": "static",
                        "type": "list",
                        "children": [
                            {
                                "name": "route",
                                "type": "list",
                                "children": [
                                    {
                                        "name": "next_hop",
                                        "type": "choice",
                                        "children": [
                                            {
                                                "name": "case_1",
                                                "type": "case",
                                                "children": [
                                                    {"name": "next-hop", "type": "leaf-list"},
                                                ],
                                            },
                                        ],
                                    },
                                ],
                            },
                        ],
                    },
                ],
            },
        ],
    }

    type_map = mod.build_type_map(schema)
    xml_leaf_list_path = "routing_options/static/route/next_hop"
    assert type_map[xml_leaf_list_path]["type"] == "leaf-list"

    xml = mod.etree.fromstring("<next-hop>100.123.0.1</next-hop>")
    assert mod.parse_element(
        xml,
        explicit_empty_tags=set(),
        type_lookup=type_map,
        parent_path="routing-options/static/route",
    ) == ["100.123.0.1"]


def test_xml2yaml_resolves_duplicate_tags_by_xml_path(tmp_path):
    mod = _load_script("jtaf-xml2yaml", "jtaf_xml2yaml_path_mod")
    schema = {
        "root": {
            "name": "root",
            "children": [
                {
                    "name": "configuration",
                    "type": "container",
                    "children": [
                        {
                            "name": "routing-options",
                            "type": "container",
                            "children": [{"name": "address", "type": "leaf"}],
                        },
                        {
                            "name": "interfaces",
                            "type": "container",
                            "children": [
                                {
                                    "name": "interface",
                                    "type": "list",
                                    "children": [
                                        {"name": "name", "type": "leaf"},
                                        {
                                            "name": "unit",
                                            "type": "list",
                                            "children": [
                                                {"name": "name", "type": "leaf"},
                                                {
                                                    "name": "family",
                                                    "type": "container",
                                                    "children": [
                                                        {
                                                            "name": "inet",
                                                            "type": "container",
                                                            "children": [
                                                                {
                                                                    "name": "address",
                                                                    "type": "list",
                                                                    "children": [
                                                                        {"name": "name", "type": "leaf"},
                                                                    ],
                                                                },
                                                            ],
                                                        },
                                                    ],
                                                },
                                            ],
                                        },
                                    ],
                                },
                            ],
                        },
                    ],
                },
            ],
        },
    }
    xml_file = tmp_path / "device.xml"
    xml_file.write_text(
        "<rpc-reply><configuration><interfaces><interface><name>ge-0/0/0</name>"
        "<unit><name>0</name><family><inet><address><name>192.0.2.1/24</name>"
        "</address></inet></family></unit></interface></interfaces></configuration></rpc-reply>",
        encoding="utf-8",
    )

    _, payload, _ = mod.parse_xml_to_payload(str(xml_file), schema)

    address = payload["interfaces"]["interface"][0]["unit"][0]["family"]["inet"]["address"][0]
    assert address == {"name": "192.0.2.1/24"}


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

    # The same schema gzipped, as the generators write it, gives the same .tf.
    from junosterraform.jtaf_common import write_schema_json_gz
    gz_schema = tmp_path / "trimmed_schema.json.gz"
    write_schema_json_gz(schema, str(gz_schema))
    out_gz = tmp_path / "tf_gz"
    argv[argv.index(str(schema_file))] = str(gz_schema)
    argv[argv.index(str(out))] = str(out_gz)
    monkeypatch.setattr(sys, "argv", argv)
    mod.main()
    assert (out_gz / "leaf1.tf").read_text() == (out / "leaf1.tf").read_text()


def test_xml2tf_rejects_non_json_schema(tmp_path, monkeypatch):
    mod = _load_script("jtaf-xml2tf", "jtaf_xml2tf_badschema_mod")
    bad = tmp_path / "trimmed_schema.json"
    bad.write_text("<configuration/>")
    xml_file = tmp_path / "leaf1.xml"
    xml_file.write_text("<configuration/>")
    monkeypatch.setattr(sys, "argv", [
        "jtaf-xml2tf", "-j", str(bad), "-x", str(xml_file), "-t", "qfx", "-d", str(tmp_path / "tf"),
    ])
    with pytest.raises(SystemExit) as exc:
        mod.main()
    assert exc.value.code == 2


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


def test_yang2go_passes_generic_and_exclude(tmp_path, monkeypatch):
    yang_file = tmp_path / "a.yang"
    yang_file.write_text("module a { namespace \"x\"; prefix x; }")

    import subprocess

    commands = []

    class _RecordingPopen(_FakePopen):
        def __init__(self, cmd, **kwargs):
            commands.append(cmd)
            super().__init__(cmd, **kwargs)

    monkeypatch.setattr(subprocess, "Popen", _RecordingPopen)
    monkeypatch.setattr(sys, "argv", [
        "jtaf-yang2go", "-p", str(yang_file), "-t", "srx", "--generic",
        "--exclude", "groups", "--exclude", "system/services/web-management",
    ])
    runpy.run_path(str(JUNOS_DIR / "jtaf-yang2go"), run_name="__main__")

    provider = next(c for c in commands if c[0] == "jtaf-provider")
    assert provider[provider.index("--generic")] == "--generic"
    excludes = [provider[i + 1] for i, a in enumerate(provider) if a == "--exclude"]
    assert excludes == ["groups", "system/services/web-management"]


@pytest.mark.parametrize(
    ("script_name", "argv"),
    [
        ("jtaf-yang2go", ["-p", "a.yang", "-x", "config.xml", "-t", "qfx", "--generic"]),
        ("jtaf-provider", ["-j", "-", "-x", "config.xml", "-t", "qfx", "--generic"]),
        ("jtaf-provider", ["-j", "-", "--xml-config", "config.xml", "-t", "qfx", "--generic"]),
    ],
    ids=["yang2go", "provider-short-option", "provider-long-option"],
)
def test_generic_generation_rejects_xml_filter(script_name, argv, monkeypatch, capsys):
    monkeypatch.setattr(sys, "argv", [script_name, *argv])

    with pytest.raises(SystemExit) as exc_info:
        runpy.run_path(str(JUNOS_DIR / script_name), run_name="__main__")

    assert exc_info.value.code == 2
    assert "not allowed with argument" in capsys.readouterr().err


def test_provider_schema_scope_help(monkeypatch, capsys):
    monkeypatch.setattr(sys, "argv", ["jtaf-provider", "--help"])

    with pytest.raises(SystemExit) as exc_info:
        runpy.run_path(str(JUNOS_DIR / "jtaf-provider"), run_name="__main__")

    assert exc_info.value.code == 0
    help_text = " ".join(capsys.readouterr().out.split())
    assert "Embed the untrimmed model rather than trimming it to XML" in help_text
    assert "trim the embedded schema to" in help_text


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
