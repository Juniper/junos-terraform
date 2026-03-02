import json
import shutil
import subprocess
import sys
from pathlib import Path

import pytest


pytest.importorskip("pyang")


def run_pyang(yang_text: str, plugin_dir: Path, tmp_path: Path):
    """Write YANG text to a file, run pyang with the jtaf plugin, return parsed JSON."""
    yang_file = tmp_path / "test.yang"
    yang_file.write_text(yang_text)

    runner = shutil.which("pyang")
    if runner:
        cmd = [runner, "--plugindir", str(plugin_dir), "-f", "jtaf", str(yang_file)]
    else:
        # fallback to python -m pyang
        cmd = [sys.executable, "-m", "pyang", "--plugindir", str(plugin_dir), "-f", "jtaf", str(yang_file)]

    proc = subprocess.run(cmd, capture_output=True, text=True)
    if proc.returncode != 0:
        pytest.fail(f"pyang failed: {proc.stderr}\ncmd: {cmd}")

    return json.loads(proc.stdout)


def test_simple_container_outputs_children(tmp_path: Path):
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module testmod {
  namespace "urn:test";
  prefix t;

  container top {
    leaf name {
      type string;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)

    assert "root" in data
    root = data["root"]
    assert "children" in root
    # top container should be first child
    top = root["children"][0]
    assert top["name"] == "top"
    # ensure leaf exists under container
    names = [c["name"] for c in top.get("children", [])]
    assert "name" in names


def test_enum_type_emits_enums(tmp_path: Path):
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module enummod {
  namespace "urn:enum";
  prefix e;

  container top {
    leaf status {
      type enumeration {
        enum up { value 0; }
        enum down { value 1; }
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)

    root = data["root"]
    top = root["children"][0]
    # find the 'status' leaf
    status = None
    for c in top.get("children", []):
        if c.get("name") == "status":
            status = c
            break

    assert status is not None, "status leaf not found"
    assert "enums" in status
    assert isinstance(status["enums"], list) and len(status["enums"]) >= 1
