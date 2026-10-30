#!/usr/bin/env python3
"""Probe how Junos devices encode YANG choices, for issue #168.

Issue #168 proposes that the generic provider stop using the build-time
`--yang-compliant` flag and instead normalise on read while *always writing the
paired form* (`<choice-ident>exact</choice-ident><choice-value/>`). That rule is
only safe if every platform and release accepts a paired write and means the
same thing by it as by the native form (`<exact/>`).

This probes that claim on real devices.

By default nothing is committed: each case is loaded into a *locked* candidate,
inspected, and discarded. The candidate is left untouched if it already holds
someone else's uncommitted changes.

Stage 5 is the exception and is opt-in (--commit-knob): the knob only takes
effect for NETCONF sessions after a commit. It uses `commit confirmed` so the
device reverts itself even if this script is killed, then reverts explicitly.

Examples
--------
    # Everything read-only, every device in the inventory
    python choice_probe.py --devices devices.json

    # Named devices only
    python choice_probe.py --devices devices.json --targets vsrx1 vptx1

    # One device, password prompted
    python choice_probe.py --host 10.1.2.3 --user admin

    # Include the committing knob-on stage on a lab device you own
    python choice_probe.py --devices devices.json --targets vsrx1 --commit-knob
"""
from __future__ import annotations

import argparse
import re
import sys
import time

from ncclient.operations import RPCError

import probe_cases as cases
from probe_lib import (DEFAULT_PORT, candidate_diff, connect, facts,
                       knob_committed, load_then_diff, port_open,
                       resolve_targets, try_load)

SEPARATOR = "=" * 72


def classify(xml):
    if "choice-ident" in xml:
        return "PAIRED"
    if re.search(r"<(exact|orlonger|upto|longer|through)\s*/>", xml):
        return "NATIVE"
    return "UNKNOWN"


def banner(text):
    print(f"\n{SEPARATOR}\n{text}\n{SEPARATOR}")


class ProbeResult:
    def __init__(self, name):
        self.name = name
        self.model = self.version = "?"
        self.knob_supported = None
        self.writes = {}
        self.equivalence = {}
        self.read_form_off = None
        self.read_form_on = None
        self.delete_address_only = None
        self.delete_full_key = None
        self.skipped = None
        self.error = None

    @property
    def paired_writes_ok(self):
        paired = [v for k, v in self.writes.items() if "paired" in k]
        return bool(paired) and all(paired)

    @property
    def all_equivalent(self):
        return bool(self.equivalence) and all(self.equivalence.values())


def stage_writes(conn, result):
    """Is each form accepted at all?"""
    print("\n-- write acceptance (candidate only, discarded) --")
    ok, detail = try_load(conn, cases.KNOB_SUPPORTED)
    result.knob_supported = ok
    print(f"  {'ACCEPT' if ok else 'REJECT':<6} yang-compliant knob supported  {detail[:90]}")

    for label, (native, paired) in cases.PAIRS.items():
        for form, cfg in (("native", native), ("paired", paired)):
            ok, detail = try_load(conn, cfg)
            result.writes[f"{label} [{form}]"] = ok
            print(f"  {'ACCEPT' if ok else 'REJECT':<6} {label} [{form}]  {detail[:90]}")


def stage_equivalence(conn, result):
    """Do native and paired produce the same configuration?"""
    print("\n-- native vs paired produce the same config? --")
    for label, (native_cfg, paired_cfg) in cases.PAIRS.items():
        native, nerr = load_then_diff(conn, native_cfg)
        paired, perr = load_then_diff(conn, paired_cfg)
        if native is None or paired is None:
            result.equivalence[label] = False
            print(f"  ERROR  {label}: native={nerr[:60]} paired={perr[:60]}")
            continue
        same = native == paired
        result.equivalence[label] = same
        print(f"  {'SAME  ' if same else 'DIFFER'} {label}")
        if not same:
            print("    native:")
            print("\n".join("      " + x for x in native.splitlines()))
            print("    paired:")
            print("\n".join("      " + x for x in paired.splitlines()))


def stage_read(conn, result, attr):
    """What encoding comes back on read, with the knob as it currently is?"""
    print("\n-- read-back encoding --")
    try:
        conn.edit_config(target="candidate", config=cases.READ_PROBE)
    except RPCError as exc:
        print(f"  could not stage read probe: {exc}")
        return
    try:
        xml = conn.get_config(source="candidate",
                              filter=("subtree", cases.READ_FILTER)).data_xml
        form = classify(xml)
        setattr(result, attr, form)
        body = " ".join(xml[xml.find("<route-filter"):xml.rfind("</term>")].split())
        print(f"  encoding: {form}")
        print(f"  {body[:320]}")
    finally:
        try:
            conn.discard_changes()
        except RPCError:
            pass


def stage_deletes(conn, result):
    """With the knob on, a delete should need the whole compound key."""
    print("\n-- delete semantics --")
    try:
        conn.edit_config(target="candidate", config=cases.READ_PROBE)
    except RPCError as exc:
        print(f"  could not stage entries to delete: {exc}")
        return
    for label, cfg, attr in (
        ("delete naming address only", cases.DELETE_ADDRESS_ONLY, "delete_address_only"),
        ("delete with full compound key", cases.DELETE_FULL_KEY, "delete_full_key"),
    ):
        try:
            conn.edit_config(target="candidate", config=cfg)
            setattr(result, attr, True)
            print(f"  ACCEPT {label}")
        except RPCError as exc:
            setattr(result, attr, False)
            print(f"  REJECT {label}: {' '.join(str(exc).split())[:110]}")
    try:
        conn.discard_changes()
    except RPCError:
        pass


def set_knob(target, config_xml, confirmed_minutes=None):
    with connect(target) as conn:
        conn.lock(target="candidate")
        try:
            conn.discard_changes()
            conn.edit_config(target="candidate", config=config_xml)
            if confirmed_minutes:
                conn.commit(confirmed=True, timeout=str(confirmed_minutes * 60))
            else:
                conn.commit()
        finally:
            conn.unlock(target="candidate")


def stage_knob_on(target, result, confirm_minutes):
    """Commit the knob, observe knob-on behaviour, then revert. Opt-in."""
    banner(f"{target.name}: COMMITTING yang-compliant (auto-reverts in "
           f"{confirm_minutes} min)")
    try:
        set_knob(target, cases.KNOB_ON, confirmed_minutes=confirm_minutes)
    except Exception as exc:  # noqa: BLE001 - report and leave device alone
        print(f"  could not commit the knob: {type(exc).__name__}: {str(exc)[:160]}")
        return

    try:
        time.sleep(3)  # a fresh session picks up the new encoding
        with connect(target) as conn:
            print(f"  knob now committed: {knob_committed(conn)}")
            conn.lock(target="candidate")
            try:
                conn.discard_changes()
                stage_read(conn, result, "read_form_on")
                stage_deletes(conn, result)
            finally:
                try:
                    conn.discard_changes()
                finally:
                    conn.unlock(target="candidate")
    finally:
        print("\n-- reverting knob --")
        for attempt in range(3):
            try:
                set_knob(target, cases.KNOB_OFF)
                with connect(target) as conn:
                    print(f"  knob still committed: {knob_committed(conn)}")
                break
            except Exception as exc:  # noqa: BLE001 - must keep trying to revert
                print(f"  revert attempt {attempt + 1} failed: {str(exc)[:120]}")
                time.sleep(5)
        else:
            print(f"  !! REVERT FAILED - `commit confirmed` reverts the device "
                  f"within {confirm_minutes} minutes")


def probe(target, commit_knob, confirm_minutes):
    result = ProbeResult(target.name)
    if not port_open(target.host, target.port):
        result.error = f"no route to {target.host}:{target.port}"
        banner(f"{target.name}: UNREACHABLE ({result.error})")
        return result

    with connect(target) as conn:
        info = facts(conn)
        result.model, result.version = info["model"], info["version"]
        knob = knob_committed(conn)
        banner(f"{target.name}  {info['model']}  {info['version']}  "
               f"yang-compliant={'ON' if knob else 'OFF'}")

        conn.lock(target="candidate")
        try:
            pending = candidate_diff(conn)
            if pending:
                result.skipped = "candidate holds uncommitted changes"
                print("  SKIP: the candidate already has uncommitted changes; "
                      "not touching this device.")
                print("\n".join("    " + x for x in pending.splitlines()[:20]))
                return result
            stage_writes(conn, result)
            stage_equivalence(conn, result)
            stage_read(conn, result, "read_form_on" if knob else "read_form_off")
            if knob:
                stage_deletes(conn, result)
        finally:
            try:
                conn.discard_changes()
            finally:
                conn.unlock(target="candidate")

    if commit_knob and not result.skipped:
        with connect(target) as conn:
            already_on = knob_committed(conn)
        if already_on:
            print("\n  knob already ON; skipping the commit stage.")
        else:
            stage_knob_on(target, result, confirm_minutes)
    return result


def summarise(results):
    banner("SUMMARY")
    header = (f"{'device':<14}{'model':<10}{'version':<13}{'knob':<6}"
              f"{'paired w':<10}{'equiv':<8}{'read off':<10}{'read on':<9}")
    print(header)
    print("-" * len(header))
    for r in results:
        if r.error or r.skipped:
            print(f"{r.name:<14}{(r.error or r.skipped)[:60]}")
            continue
        print(f"{r.name:<14}{r.model:<10}{r.version:<13}"
              f"{str(r.knob_supported):<6}"
              f"{('yes' if r.paired_writes_ok else 'NO'):<10}"
              f"{('yes' if r.all_equivalent else 'NO'):<8}"
              f"{str(r.read_form_off or '-'):<10}{str(r.read_form_on or '-'):<9}")

    print("\nIssue #168 is safe on a device when: paired writes are accepted, "
          "and native/paired are equivalent.")
    bad = [r for r in results if not (r.error or r.skipped)
           and not (r.paired_writes_ok and r.all_equivalent)]
    if bad:
        print("\nPROBLEM on: " + ", ".join(r.name for r in bad)
              + "\n  -> #168's 'always write paired' rule does NOT hold here.")
    else:
        print("\nAll probed devices support the proposal's write rule.")
    for r in results:
        if r.delete_address_only is False and r.delete_full_key:
            print(f"  {r.name}: with the knob on, deletes need the full "
                  f"compound key (address alone is rejected).")


def main():
    p = argparse.ArgumentParser(
        description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("--devices", help="inventory JSON (or set JTAF_PROBE_DEVICES)")
    p.add_argument("--targets", nargs="*", help="inventory names to probe; default all")
    p.add_argument("--host", help="probe a single host instead of an inventory")
    p.add_argument("--user", help="username for --host (or set JTAF_PROBE_USER)")
    p.add_argument("--port", type=int, default=None,
                   help=f"NETCONF port (default {DEFAULT_PORT})")
    p.add_argument("--commit-knob", action="store_true",
                   help="also commit yang-compliant to observe the knob-on read "
                        "form; uses commit confirmed and reverts afterwards")
    p.add_argument("--confirm-minutes", type=int, default=3,
                   help="commit confirmed auto-revert window (default 3)")
    args = p.parse_args()

    targets = resolve_targets(names=args.targets, devices_path=args.devices,
                              host=args.host, username=args.user, port=args.port)
    print(f"probing {len(targets)} device(s); "
          f"{'COMMITTING knob where asked' if args.commit_knob else 'no commits'}")

    results = []
    for t in targets:
        try:
            results.append(probe(t, args.commit_knob, args.confirm_minutes))
        except Exception as exc:  # noqa: BLE001 - keep surveying other devices
            r = ProbeResult(t.name)
            r.error = f"{type(exc).__name__}: {str(exc)[:140]}"
            results.append(r)
            print(f"\n  ERROR on {t.name}: {r.error}")
    summarise(results)
    return 0


if __name__ == "__main__":
    sys.exit(main())
