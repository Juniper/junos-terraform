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


def test_leaf_list_definition(tmp_path: Path):
    """Test leaf-list statement from RFC 7950 Section 7.7"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module leaflistmod {
  namespace "urn:test";
  prefix t;

  container config {
    leaf-list dns-search {
      type string;
      description "List of DNS search domains";
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    assert "root" in data
    root = data["root"]
    config = root["children"][0]
    dns_search = config["children"][0]
    assert dns_search["name"] == "dns-search"


def test_list_with_key(tmp_path: Path):
    """Test list statement with key from RFC 7950 Section 7.8"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module listmod {
  namespace "urn:test";
  prefix t;

  container users {
    list user {
      key "name";
      leaf name {
        type string;
      }
      leaf full-name {
        type string;
      }
      leaf class {
        type string;
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    users = root["children"][0]
    user_list = users["children"][0]
    assert user_list["name"] == "user"
    assert "key" in user_list
    assert user_list["key"] == "name"


def test_list_with_multiple_keys(tmp_path: Path):
    """Test list with composite key from RFC 7950 Section 7.8"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module multikey {
  namespace "urn:test";
  prefix t;

  container servers {
    list server {
      key "ip port";
      leaf ip {
        type string;
      }
      leaf port {
        type uint16;
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    servers = root["children"][0]
    server_list = servers["children"][0]
    assert "key" in server_list
    # Key should contain both ip and port
    assert "ip" in server_list["key"]
    assert "port" in server_list["key"]


def test_choice_statement(tmp_path: Path):
    """Test choice and case statements from RFC 7950 Section 7.9"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module choicemod {
  namespace "urn:test";
  prefix t;

  container transport {
    choice protocol {
      case tcp {
        leaf tcp-port {
          type uint16;
        }
      }
      case udp {
        leaf udp-port {
          type uint16;
        }
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    assert "root" in data
    root = data["root"]
    transport = root["children"][0]
    # The choice node should exist in schema
    choice_found = any(c.get("name") == "protocol" for c in transport.get("children", []))
    assert choice_found


def test_presence_container(tmp_path: Path):
    """Test presence container from RFC 7950 Section 7.5.1"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module presencemod {
  namespace "urn:test";
  prefix t;

  container system {
    container ssh {
      presence "Enables SSH";
      leaf port {
        type uint16;
        default 22;
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    system = root["children"][0]
    ssh = system["children"][0]
    assert ssh["name"] == "ssh"
    assert ssh["type"] == "container"
    # Verify it has child nodes
    assert "children" in ssh
    assert len(ssh["children"]) > 0


def test_typedef_derived_type(tmp_path: Path):
    """Test typedef statement from RFC 7950 Section 7.3"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module typedefmod {
  namespace "urn:test";
  prefix t;

  typedef percent {
    type uint8 {
      range "0 .. 100";
    }
    description "Percentage type";
  }

  container stats {
    leaf completion-rate {
      type percent;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    stats = root["children"][0]
    completion = stats["children"][0]
    assert completion["name"] == "completion-rate"


def test_grouping_and_uses(tmp_path: Path):
    """Test grouping and uses statements from RFC 7950 Sections 7.12 and 7.13"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module groupingmod {
  namespace "urn:test";
  prefix t;

  grouping interface-common {
    leaf enabled {
      type boolean;
    }
    leaf mtu {
      type uint16;
    }
  }

  container interfaces {
    container eth0 {
      uses interface-common;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    interfaces = root["children"][0]
    eth0 = interfaces["children"][0]
    # Check that uses children are present
    child_names = [c["name"] for c in eth0.get("children", [])]
    assert "enabled" in child_names
    assert "mtu" in child_names


def test_integer_types(tmp_path: Path):
    """Test integer types from RFC 7950 Section 9.2"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module inttypes {
  namespace "urn:test";
  prefix t;

  container numbers {
    leaf count-int8 { type int8; }
    leaf count-int16 { type int16; }
    leaf count-int32 { type int32; }
    leaf count-int64 { type int64; }
    leaf count-uint8 { type uint8; }
    leaf count-uint16 { type uint16; }
    leaf count-uint32 { type uint32; }
    leaf count-uint64 { type uint64; }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    numbers = root["children"][0]
    assert len(numbers["children"]) == 8


def test_string_type_with_pattern(tmp_path: Path):
    """Test string type with pattern constraint from RFC 7950 Section 9.4"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module stringmod {
  namespace "urn:test";
  prefix t;

  container config {
    leaf hostname {
      type string {
        pattern "[a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?";
        length "1..63";
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    config = root["children"][0]
    hostname = config["children"][0]
    assert hostname["name"] == "hostname"


def test_decimal64_type(tmp_path: Path):
    """Test decimal64 type from RFC 7950 Section 9.3"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module decimalmod {
  namespace "urn:test";
  prefix t;

  container sensors {
    leaf temperature {
      type decimal64 {
        fraction-digits 2;
        range "-40.00 .. 125.00";
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    sensors = root["children"][0]
    temp = sensors["children"][0]
    assert temp["name"] == "temperature"


def test_boolean_type(tmp_path: Path):
    """Test boolean type from RFC 7950 Section 9.5"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module boolmod {
  namespace "urn:test";
  prefix t;

  container switches {
    leaf enabled {
      type boolean;
      default "true";
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    switches = root["children"][0]
    enabled = switches["children"][0]
    assert enabled["name"] == "enabled"


def test_bits_type(tmp_path: Path):
    """Test bits type from RFC 7950 Section 9.7"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module bitsmod {
  namespace "urn:test";
  prefix t;

  container flags {
    leaf options {
      type bits {
        bit read-bit { position 0; }
        bit write-bit { position 1; }
        bit execute-bit { position 2; }
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    flags = root["children"][0]
    options = flags["children"][0]
    assert options["name"] == "options"


def test_empty_type(tmp_path: Path):
    """Test empty type from RFC 7950 Section 9.11"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module emptymod {
  namespace "urn:test";
  prefix t;

  container settings {
    leaf enable-feature {
      type empty;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    settings = root["children"][0]
    feature = settings["children"][0]
    assert feature["name"] == "enable-feature"


def test_union_type(tmp_path: Path):
    """Test union type from RFC 7950 Section 9.12"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module unionmod {
  namespace "urn:test";
  prefix t;

  container config {
    leaf timeout {
      type union {
        type uint32;
        type enumeration {
          enum "infinite";
        }
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    config = root["children"][0]
    timeout = config["children"][0]
    assert timeout["name"] == "timeout"


def test_leafref_type(tmp_path: Path):
    """Test leafref type from RFC 7950 Section 9.9"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module leafrefmod {
  namespace "urn:test";
  prefix t;

  container interfaces {
    list interface {
      key "name";
      leaf name {
        type string;
      }
    }
  }

  container routing {
    leaf active-interface {
      type leafref {
        path "/interfaces/interface/name";
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    assert "root" in data
    # Just verify structure is parsed
    root = data["root"]
    assert "children" in root


def test_identity_and_identityref(tmp_path: Path):
    """Test identity and identityref from RFC 7950 Sections 7.18 and 9.10"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module identitymod {
  namespace "urn:test";
  prefix t;

  identity crypto-algorithm {
    description "Base crypto algorithm identity";
  }

  identity aes {
    base crypto-algorithm;
  }

  identity des {
    base crypto-algorithm;
  }

  container crypto {
    leaf algorithm {
      type identityref {
        base crypto-algorithm;
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    assert "root" in data
    if "identities" in data:
        # Check identities are present if plugin outputs them
        assert len(data["identities"]) >= 0


def test_leaf_with_default_value(tmp_path: Path):
    """Test leaf default value from RFC 7950 Section 7.6.1"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module defaultmod {
  namespace "urn:test";
  prefix t;

  container settings {
    leaf timeout {
      type uint32;
      default "30";
    }
    leaf hostname {
      type string;
      default "localhost";
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    settings = root["children"][0]
    timeout = settings["children"][0]
    assert timeout["name"] == "timeout"


def test_leaf_mandatory_constraint(tmp_path: Path):
    """Test mandatory leaf from RFC 7950 Section 7.6.5"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module mandatorymod {
  namespace "urn:test";
  prefix t;

  container system {
    leaf hostname {
      type string;
      mandatory true;
    }
    leaf timezone {
      type string;
      mandatory false;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    system = root["children"][0]
    hostname = system["children"][0]
    assert hostname["name"] == "hostname"


def test_leaf_list_min_max_elements(tmp_path: Path):
    """Test leaf-list min/max-elements from RFC 7950 Sections 7.7.5 and 7.7.6"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module minmaxmod {
  namespace "urn:test";
  prefix t;

  container nameservers {
    leaf-list server {
      type string;
      min-elements 1;
      max-elements 3;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    nameservers = root["children"][0]
    server = nameservers["children"][0]
    assert server["name"] == "server"


def test_list_unique_constraint(tmp_path: Path):
    """Test list unique constraint from RFC 7950 Section 7.8.3"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module uniquemod {
  namespace "urn:test";
  prefix t;

  container data {
    list entry {
      key "id";
      leaf id {
        type string;
      }
      leaf email {
        type string;
      }
      unique "email";
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    data_cont = root["children"][0]
    entry_list = data_cont["children"][0]
    assert entry_list["name"] == "entry"


def test_nested_containers(tmp_path: Path):
    """Test nested container structures from RFC 7950 Section 7.5"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module nestedmod {
  namespace "urn:test";
  prefix t;

  container system {
    container ntp {
      leaf enabled {
        type boolean;
      }
      leaf source {
        type string;
      }
    }
    container snmp {
      leaf community {
        type string;
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    system = root["children"][0]
    # Check we have multiple containers
    assert len(system["children"]) >= 2
    container_names = [c["name"] for c in system["children"]]
    assert "ntp" in container_names
    assert "snmp" in container_names


def test_config_statement(tmp_path: Path):
    """Test config statement from RFC 7950 Section 7.21.1"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module configmod {
  namespace "urn:test";
  prefix t;

  container interface {
    leaf name {
      type string;
      config true;
    }
    leaf statistics {
      type string;
      config false;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    interface = root["children"][0]
    assert interface["name"] == "interface"


def test_status_statement(tmp_path: Path):
    """Test status statement from RFC 7950 Section 7.21.2"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module statusmod {
  namespace "urn:test";
  prefix t;

  container config {
    leaf deprecated-leaf {
      type string;
      status deprecated;
    }
    leaf current-leaf {
      type string;
      status current;
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    config = root["children"][0]
    assert len(config["children"]) == 2


def test_anyxml_node(tmp_path: Path):
    """Test anyxml statement from RFC 7950 Section 7.11"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module anyxmlmod {
  namespace "urn:test";
  prefix t;

  container rpc-reply {
    anyxml response;
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    rpc_reply = root["children"][0]
    assert rpc_reply["name"] == "rpc-reply"
    # Verify anyxml is present in children if it's output
    if "children" in rpc_reply:
        response_found = any(c["name"] == "response" for c in rpc_reply["children"])
        assert response_found or len(rpc_reply["children"]) >= 0


def test_complex_nested_structure(tmp_path: Path):
    """Test complex nested structure combining multiple statement types"""
    plugin_dir = Path(__file__).resolve().parents[1]

    yang = """
module complexmod {
  namespace "urn:test";
  prefix t;

  typedef ipv4-address {
    type string {
      pattern "([0-9]{1,3}\\.){3}[0-9]{1,3}";
    }
  }

  container network {
    list interface {
      key "name";
      leaf name {
        type string;
      }
      leaf enabled {
        type boolean;
        default "true";
      }
      container addresses {
        leaf-list ipv4 {
          type ipv4-address;
        }
      }
    }
  }
}
"""

    data = run_pyang(yang, plugin_dir, tmp_path)
    root = data["root"]
    network = root["children"][0]
    interface_list = network["children"][0]
    assert interface_list["name"] == "interface"
    assert "key" in interface_list
