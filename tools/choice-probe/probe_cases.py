"""Configuration snippets used to probe how a device encodes a YANG choice.

Junos returns a choice in one of two forms depending on
`set system services netconf yang-compliant`:

    native   <exact/>
    paired   <choice-ident>exact</choice-ident><choice-value/>

Issue #168 proposes always writing the paired form. These snippets test that
claim: whether each form is accepted, and whether both produce the same
configuration.
"""

POLICY = "JTAF-CHOICE-PROBE"
COMMUNITY = "JTAF-PROBE-C"
PROBE_PREFIX_A = "10.0.0.0/8"
PROBE_PREFIX_B = "10.1.0.0/16"
PROBE_MARTIAN = "192.0.2.0/24"


def policy(from_body="", then_body=""):
    return (f"<config><configuration><policy-options>"
            f"<policy-statement><name>{POLICY}</name>"
            f"<term><name>t1</name>"
            f"<from>{from_body}</from><then>{then_body}</then>"
            f"</term></policy-statement></policy-options></configuration></config>")


def route_filter_native(address, case, value=""):
    inner = f"<{case}/>" if value == "" else f"<{case}>{value}</{case}>"
    return f"<route-filter><address>{address}</address>{inner}</route-filter>"


def route_filter_paired(address, case, value=""):
    return (f"<route-filter><address>{address}</address>"
            f"<choice-ident>{case}</choice-ident>"
            f"<choice-value>{value}</choice-value></route-filter>")


def martians(paired):
    choice = ("<choice-ident>orlonger</choice-ident><choice-value></choice-value>"
              if paired else "<orlonger/>")
    return (f"<config><configuration><routing-options><martians>"
            f"<address>{PROBE_MARTIAN}</address>{choice}<allow/>"
            f"</martians></routing-options></configuration></config>")


def community(paired):
    choice = ("<choice-ident>add</choice-ident><choice-value></choice-value>"
              if paired else "<add/>")
    return policy("", f"<community>{choice}<community-name>{COMMUNITY}</community-name></community>")


KNOB_SUPPORTED = ("<config><configuration><system><services><netconf>"
                  "<yang-compliant/></netconf></services></system>"
                  "</configuration></config>")

KNOB_ON = KNOB_SUPPORTED
KNOB_OFF = ('<config><configuration><system><services><netconf>'
            '<yang-compliant operation="delete"/>'
            "</netconf></services></system></configuration></config>")

# label -> (native XML, paired XML). Both members must mean the same thing to
# the device for issue #168's "always write paired" rule to be safe.
PAIRS = {
    "route-filter exact": (
        policy(route_filter_native(PROBE_PREFIX_A, "exact"), "<accept/>"),
        policy(route_filter_paired(PROBE_PREFIX_A, "exact"), "<accept/>"),
    ),
    "route-filter prefix-length-range": (
        policy(route_filter_native(PROBE_PREFIX_A, "prefix-length-range", "/24-/24"), "<accept/>"),
        policy(route_filter_paired(PROBE_PREFIX_A, "prefix-length-range", "/24-/24"), "<accept/>"),
    ),
    "route-filter orlonger + reject": (
        policy(route_filter_native(PROBE_PREFIX_A, "orlonger"), "<reject/>"),
        policy(route_filter_paired(PROBE_PREFIX_A, "orlonger"), "<reject/>"),
    ),
    "route-filter upto": (
        policy(route_filter_native(PROBE_PREFIX_A, "upto", "/24"), "<accept/>"),
        policy(route_filter_paired(PROBE_PREFIX_A, "upto", "/24"), "<accept/>"),
    ),
    "community add": (community(paired=False), community(paired=True)),
    "martians orlonger allow": (martians(paired=False), martians(paired=True)),
}

# Read-back probe: several cases at once, so one fetch shows the encoding.
READ_PROBE = policy(
    route_filter_native(PROBE_PREFIX_A, "exact")
    + route_filter_native(PROBE_PREFIX_B, "prefix-length-range", "/24-/24"),
    f"<community><add/><community-name>{COMMUNITY}</community-name></community>",
)

READ_FILTER = (f"<configuration><policy-options><policy-statement>"
               f"<name>{POLICY}</name></policy-statement></policy-options></configuration>")

# With the knob on the match type is part of the entry's key, so a delete that
# names only the address should be rejected.
DELETE_ADDRESS_ONLY = policy(
    f'<route-filter operation="delete"><address>{PROBE_PREFIX_A}</address></route-filter>')
DELETE_FULL_KEY = policy(
    f'<route-filter operation="delete"><address>{PROBE_PREFIX_A}</address>'
    "<choice-ident>exact</choice-ident><choice-value></choice-value></route-filter>")
