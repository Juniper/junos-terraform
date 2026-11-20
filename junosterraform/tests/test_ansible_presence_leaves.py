from pathlib import Path
import xml.etree.ElementTree as ET

from jinja2 import ChainableUndefined, Environment
import pytest


TEMPLATE = Path(__file__).resolve().parents[1] / "templates" / "ansible.j2"


def render_configuration(values, nested=False):
    leaves = [
        {"name": "mtu-discovery", "type": "leaf", "leaf-type": "empty"},
        {"name": "description", "type": "leaf", "leaf-type": "string"},
        {"name": "enabled", "type": "leaf", "leaf-type": "boolean"},
    ]
    for leaf in leaves:
        leaf["path"] = "protocols/group" if nested else "protocols"
    children = (
        [{"name": "group", "type": "list", "path": "protocols", "children": leaves}]
        if nested else leaves
    )
    schema = {"root": {"children": [{
        "name": "configuration", "children": [{
            "name": "protocols", "type": "container", "path": "", "children": children,
        }],
    }]}}
    environment = Environment(undefined=ChainableUndefined)
    generated = environment.from_string(TEMPLATE.read_text()).render(data=schema)
    generated = "\n".join(
        line for line in generated.splitlines() if not line.startswith("#jinja2:")
    )
    variables = {"group": [values]} if nested else values
    output = environment.from_string(generated).render(
        jtaf_effective={"protocols": variables}
    )
    return ET.fromstring(output)


@pytest.mark.parametrize("nested", [False, True])
@pytest.mark.parametrize("value", [True, "true", "TRUE", "True"])
def test_empty_leaf_true_renders_without_text(value, nested):
    root = render_configuration({"mtu_discovery": value}, nested)
    leaves = root.findall(".//mtu-discovery")
    assert len(leaves) == 1
    assert leaves[0].text is None
    assert len(leaves[0]) == 0


@pytest.mark.parametrize("nested", [False, True])
@pytest.mark.parametrize(
    "values",
    [{}, {"mtu_discovery": False},
     {"mtu_discovery": "false"},
     {"mtu_discovery": "FALSE"}],
)
def test_empty_leaf_false_or_undefined_is_omitted(values, nested):
    root = render_configuration(values, nested)
    assert root.find(".//mtu-discovery") is None


@pytest.mark.parametrize("value", ["true", "false", "unchanged"])
def test_nonempty_leaf_strings_are_not_normalized(value):
    root = render_configuration({"description": value, "enabled": value})
    assert root.findtext(".//description") == value
    assert root.findtext(".//enabled") == value
