"""
Regression comparison test for generic vs legacy provider generation modes.

This test ensures both generation modes produce valid, compilable Go providers
from the same JSON schema input.
"""
import json
import os
import shutil
import subprocess
import tempfile
import pytest


SMALL_SCHEMA = {
    "root": {
        "children": [
            {
                "name": "system",
                "type": "container",
                "path": "",
                "children": [
                    {"name": "host-name", "type": "leaf", "leaf-type": "string", "path": "system"},
                    {"name": "domain-name", "type": "leaf", "leaf-type": "string", "path": "system"},
                ],
            },
            {
                "name": "interfaces",
                "type": "list",
                "key": "name",
                "path": "",
                "children": [
                    {"name": "name", "type": "leaf", "leaf-type": "string", "path": "interfaces"},
                    {"name": "description", "type": "leaf", "leaf-type": "string", "path": "interfaces"},
                    {"name": "mtu", "type": "leaf", "leaf-type": "string", "path": "interfaces"},
                ],
            },
        ]
    }
}


@pytest.fixture
def schema_file(tmp_path):
    """Write the test schema to a temporary JSON file."""
    schema_path = tmp_path / "test_schema.json"
    schema_path.write_text(json.dumps(SMALL_SCHEMA, indent=2))
    return str(schema_path)


def _find_jtaf_provider():
    """Locate jtaf-provider executable."""
    exe = shutil.which("jtaf-provider")
    if exe:
        return exe
    # Try the build/scripts directory
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    candidate = os.path.join(repo_root, "build", "scripts-3.9", "jtaf-provider")
    if os.path.exists(candidate):
        return candidate
    return None


def _find_go():
    """Locate go executable."""
    return shutil.which("go")


def _generate_provider(jtaf_exe, schema_file, device_type, mode, work_dir):
    """Run jtaf-provider and return the output directory path."""
    cmd = [jtaf_exe, "-j", schema_file, "-t", device_type, "-m", mode]
    result = subprocess.run(cmd, capture_output=True, text=True, cwd=work_dir)
    assert result.returncode == 0, f"jtaf-provider failed ({mode}): {result.stderr}"
    output_dir = os.path.join(work_dir, f"terraform-provider-junos-{device_type}")
    assert os.path.isdir(output_dir), f"Output dir not created: {output_dir}"
    return output_dir


def _go_build(go_exe, provider_dir):
    """Run 'go build ./...' in the provider directory."""
    # Remove any venv dirs that might have been copied
    venv_dir = os.path.join(provider_dir, "venv")
    if os.path.exists(venv_dir):
        shutil.rmtree(venv_dir)
    result = subprocess.run(
        [go_exe, "build", "./..."],
        capture_output=True,
        text=True,
        cwd=provider_dir,
    )
    return result


class TestRegressionComparison:
    """Compare generic and legacy generation modes."""

    def test_generic_mode_generates_valid_go(self, schema_file, tmp_path):
        """Generic mode produces a compilable Go provider."""
        jtaf_exe = _find_jtaf_provider()
        if not jtaf_exe:
            pytest.skip("jtaf-provider not found on PATH")
        go_exe = _find_go()
        if not go_exe:
            pytest.skip("go not found on PATH")

        work_dir = str(tmp_path / "generic")
        os.makedirs(work_dir)
        provider_dir = _generate_provider(jtaf_exe, schema_file, "gentest", "generic", work_dir)

        result = _go_build(go_exe, provider_dir)
        assert result.returncode == 0, f"go build failed (generic):\n{result.stderr}"

    def test_legacy_mode_generates_valid_go(self, schema_file, tmp_path):
        """Legacy mode produces a compilable Go provider."""
        jtaf_exe = _find_jtaf_provider()
        if not jtaf_exe:
            pytest.skip("jtaf-provider not found on PATH")
        go_exe = _find_go()
        if not go_exe:
            pytest.skip("go not found on PATH")

        work_dir = str(tmp_path / "legacy")
        os.makedirs(work_dir)
        provider_dir = _generate_provider(jtaf_exe, schema_file, "legtest", "legacy", work_dir)

        result = _go_build(go_exe, provider_dir)
        assert result.returncode == 0, f"go build failed (legacy):\n{result.stderr}"

    def test_generic_mode_has_trimmed_schema_json(self, schema_file, tmp_path):
        """Generic mode embeds the schema as trimmed_schema.json."""
        jtaf_exe = _find_jtaf_provider()
        if not jtaf_exe:
            pytest.skip("jtaf-provider not found on PATH")

        work_dir = str(tmp_path / "generic_schema")
        os.makedirs(work_dir)
        provider_dir = _generate_provider(jtaf_exe, schema_file, "schematest", "generic", work_dir)

        # Verify trimmed_schema.json exists
        trimmed = os.path.join(provider_dir, "trimmed_schema.json")
        assert os.path.exists(trimmed), "trimmed_schema.json not found in generic output"

        # Verify it contains our schema structure
        with open(trimmed) as f:
            data = json.load(f)
        assert "root" in data
        assert "children" in data["root"]

    def test_generic_mode_no_resource_config_provider(self, schema_file, tmp_path):
        """Generic mode should NOT contain resource_config_provider.go (legacy file)."""
        jtaf_exe = _find_jtaf_provider()
        if not jtaf_exe:
            pytest.skip("jtaf-provider not found on PATH")

        work_dir = str(tmp_path / "generic_no_legacy")
        os.makedirs(work_dir)
        provider_dir = _generate_provider(jtaf_exe, schema_file, "noltest", "generic", work_dir)

        legacy_file = os.path.join(provider_dir, "resource_config_provider.go")
        assert not os.path.exists(legacy_file), "resource_config_provider.go should not exist in generic mode"

    def test_generic_smaller_source_than_legacy_with_more_resources(self, schema_file, tmp_path):
        """
        With a larger schema, generic mode source should be smaller than legacy.
        The generic runtime is fixed-size; legacy grows with schema complexity.
        """
        jtaf_exe = _find_jtaf_provider()
        if not jtaf_exe:
            pytest.skip("jtaf-provider not found on PATH")

        # Create a larger schema with many interfaces-like resources
        large_schema = {"root": {"children": []}}
        for i in range(20):
            large_schema["root"]["children"].append({
                "name": f"resource-{i}",
                "type": "list",
                "key": "name",
                "path": "",
                "children": [
                    {"name": "name", "type": "leaf", "leaf-type": "string", "path": f"resource-{i}"},
                    {"name": "description", "type": "leaf", "leaf-type": "string", "path": f"resource-{i}"},
                    {"name": "value", "type": "leaf", "leaf-type": "string", "path": f"resource-{i}"},
                    {"name": "enabled", "type": "leaf", "leaf-type": "empty", "path": f"resource-{i}"},
                ],
            })

        large_schema_file = str(tmp_path / "large_schema.json")
        with open(large_schema_file, "w") as f:
            json.dump(large_schema, f)

        # Generate both modes
        gen_work = str(tmp_path / "gen_large")
        leg_work = str(tmp_path / "leg_large")
        os.makedirs(gen_work)
        os.makedirs(leg_work)

        gen_dir = _generate_provider(jtaf_exe, large_schema_file, "genlarge", "generic", gen_work)
        leg_dir = _generate_provider(jtaf_exe, large_schema_file, "leglarge", "legacy", leg_work)

        # Count Go source lines (excluding test files)
        def count_go_lines(d):
            total = 0
            for root, _, files in os.walk(d):
                for fname in files:
                    if fname.endswith(".go") and not fname.endswith("_test.go"):
                        with open(os.path.join(root, fname)) as fh:
                            total += sum(1 for _ in fh)
            return total

        gen_lines = count_go_lines(gen_dir)
        leg_lines = count_go_lines(leg_dir)

        # Generic should have fewer source lines for a 20-resource schema
        assert gen_lines < leg_lines, (
            f"Generic ({gen_lines} lines) should be smaller than legacy ({leg_lines} lines) "
            f"for a 20-resource schema"
        )
