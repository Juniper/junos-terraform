import xml.etree.ElementTree as ET
from pathlib import Path

from jinja2 import ChainableUndefined, Environment
import pytest

TEMPLATE = Path(__file__).resolve().parents[1] / "templates" / "ansible.j2"


def render(values, presence=False):
    ingress = {
        "name": "ingress", "type": "container",
        "path": "routing-options/chained-composite-next-hop",
        "children": [{
            "name": "evpn", "type": "leaf", "leaf-type": "empty",
            "path": "routing-options/chained-composite-next-hop/ingress",
        }],
    }
    if presence:
        ingress["presence"] = "Enable ingress"
    schema = {"root": {"children": [{
        "name": "configuration", "children": [{
            "name": "routing-options", "type": "container", "path": "",
            "children": [{
                "name": "chained-composite-next-hop", "type": "container",
                "path": "routing-options", "children": [ingress],
            }],
        }],
    }]}}
    environment = Environment(undefined=ChainableUndefined)
    generated = environment.from_string(TEMPLATE.read_text()).render(data=schema)
    generated = "\n".join(
        line for line in generated.splitlines() if not line.startswith("#jinja2:")
    )
    return ET.fromstring(environment.from_string(generated).render(
        jtaf_effective=values
    ))


@pytest.mark.parametrize("values", [
    {},
    {"routing_options": {}},
    {"routing_options": {"chained_composite_next_hop": {}}},
    {"routing_options": {"chained_composite_next_hop": {"ingress": {}}}},
    {"routing_options": {"chained_composite_next_hop": {
        "ingress": {"evpn": False},
    }}},
    {"routing_options": {"chained_composite_next_hop": {
        "ingress": {"unsupported": "value"},
    }}},
])
def test_empty_nonpresence_ancestors_are_omitted(values):
    root = render(values)
    assert root.find(".//routing-options") is None


@pytest.mark.parametrize("value", [True, "true"])
def test_nonempty_evpn_hierarchy_is_preserved(value):
    root = render({"routing_options": {
        "chained_composite_next_hop": {"ingress": {"evpn": value}},
    }})
    leaf = root.find("./groups/routing-options/chained-composite-next-hop/ingress/evpn")
    assert leaf is not None
    assert leaf.text is None


def test_empty_presence_container_preserves_its_ancestors():
    root = render({"routing_options": {
        "chained_composite_next_hop": {"ingress": {}},
    }}, presence=True)
    ingress = root.find("./groups/routing-options/chained-composite-next-hop/ingress")
    assert ingress is not None
    assert not list(ingress)


@pytest.mark.parametrize("value", [False, None])
def test_disabled_presence_container_is_omitted(value):
    root = render({"routing_options": {
        "chained_composite_next_hop": {"ingress": value},
    }}, presence=True)
    assert root.find(".//routing-options") is None
