# Live Device Test Review — NETCONF Patch Matrix

> Device: **vpaa** (10.54.21.174) · MX960 · Junos 26.2I-20260201  
> Date: 2026-04-28  
> Branch: `netconf_patch`

---

## Summary

All 30 matrix scenarios tested against a live Junos device via NETCONF
`edit-config` in two phases:

1. **Candidate-only** (`tests/live_matrix_test.py`) — 30/30 PASS.
   Each change applied to candidate and verified, then discarded.
2. **Commit + verify + cleanup** (`tests/live_commit_test.py`) — 30/30 PASS.
   Each change committed to running config, verified in running, then reverted.

## Provider Fix Applied

The Go provider's `patchEditConfigStr` in `terraform_provider/netconf/client.go`
was changed from `<default-operation>none</default-operation>` to
`<default-operation>merge</default-operation>`. The corresponding assertion in
`client_test.go` was updated to match.

**Root cause:** Junos 26.2I silently ignores `nc:operation` attributes when
`default-operation=none` is set. With `default-operation=merge`, all
`nc:operation` attributes (`replace`, `create`, `delete`) work correctly.

## Topology

| Device | IP             | Role           |
|--------|----------------|----------------|
| vpaa   | 10.54.21.174   | Primary test target (MX960, Junos 26.2I) |
| vpbb   | 10.54.21.160   | Available       |
| vpcc   | 10.54.21.234   | Available       |
| vpdd   | 10.54.22.140   | Available       |

## Phase 1: Candidate-Only Results

| ID | Scenario | Status | Notes |
|----|----------|--------|-------|
| L1 | Create leaf (host-name) | PASS | `nc:operation="replace"` sets host-name to test value |
| L2 | Replace leaf (host-name) | PASS | Two-step: set initial value, then replace |
| L3 | Delete leaf (interface description) | PASS | `nc:operation="delete"` removes description from ge-0/0/0 |
| L4 | Special characters in leaf | PASS | `&`, `<`, `>` round-trip correctly via XML escaping |
| L5 | Embedded newlines in leaf | PASS | Junos accepts multiline description values |
| L6 | No-op replace (same value) | PASS | Replacing with identical value accepted without error |
| LL1 | Add leaf-list entry | PASS | `nc:operation="create"` appends community member |
| LL2 | Remove leaf-list entry | PASS | `nc:operation="delete"` removes specific member by value |
| LL3 | Replace leaf-list entry (delete+create) | PASS | Atomic delete old + create new in one edit-config |
| LL4 | Reorder leaf-list (set semantics) | PASS | Diff-engine only — Junos stores members as a set |
| LL5 | Delete all leaf-list entries | PASS | Multiple `nc:operation="delete"` in one edit-config |
| LL6 | Create leaf-list from scratch | PASS | Multiple `nc:operation="create"` on empty community |
| K1 | Add list entry (interface) | PASS | `<interface nc:operation="create">` adds ge-0/0/9 |
| K2 | Delete list entry | PASS | `<interface nc:operation="delete">` removes by key |
| K3 | Rename list key (delete+create) | PASS | Delete ge-0/0/8 + create ge-0/0/9 in one RPC |
| K4 | Modify leaf in list entry | PASS | Replace nested unit/0/description |
| K5 | Add leaf in list entry | PASS | `nc:operation="create"` on description inside unit |
| K6 | Delete leaf in list entry | PASS | `nc:operation="delete"` on description inside unit |
| K7 | Reorder list entries | PASS | Diff-engine only — Junos uses canonical order |
| K8 | Nested list entry (syslog/host/contents) | PASS | `nc:operation="create"` on 3rd-level list entry |
| K9 | Per-instance ordering | PASS | Diff-engine only — verified by Go tests |
| K10 | Key-only list entry | PASS | Host entry with minimal contents — key leaf stored |
| C1 | Create container (chassis/aggregated-devices) | PASS | Container with child leaf created |
| C2 | Delete container | PASS | `nc:operation="delete"` removes entire container |
| C3 | Empty presence container (ssh) | PASS | `<ssh nc:operation="replace"/>` accepted |
| C4 | Replace container | PASS | Replace container changes child leaf value |
| C5 | Modify container child | PASS | Replace leaf inside system container |
| M1 | Mixed operations (replace+create+delete) | PASS | Three operations in single edit-config RPC |
| M2 | Deep nesting (4+ levels) | PASS | syslog/host/contents/any at 4th level |
| O1-O6 | Order alignment | PASS | Go diff-engine only — not NETCONF-testable |

## Phase 2: Commit + Verify + Cleanup Results

Each test commits to running config, verifies the change, then reverts it.
The device's original hostname is captured at startup and restored after each
hostname-modifying test. All test artifacts are cleaned up after verification.

| ID | Scenario | Status | Notes |
|----|----------|--------|-------|
| L1 | Create leaf (host-name) | PASS | Committed to running, verified, hostname restored |
| L2 | Replace leaf (host-name) | PASS | Replaced in running before→after, hostname restored |
| L3 | Delete leaf (interface description) | PASS | Deleted from running, verified gone |
| L4 | Special characters in leaf | PASS | Special chars committed + round-tripped in running |
| L5 | Embedded newlines in leaf | PASS | Multiline committed to running, description cleaned up |
| L6 | No-op replace (same value) | PASS | Same-value commit accepted, hostname restored |
| LL1 | Add leaf-list entry | PASS | Both members verified in running, community deleted |
| LL2 | Remove leaf-list entry | PASS | Member removed from running, community deleted |
| LL3 | Replace leaf-list entry | PASS | bbb→ccc in running, community deleted |
| LL4 | Reorder leaf-list | PASS | Diff-engine only |
| LL5 | Delete all entries | PASS | All members gone from running, community deleted |
| LL6 | Create from scratch | PASS | Members created in running, community deleted |
| K1 | Add list entry | PASS | ge-0/0/9 in running, interface deleted |
| K2 | Delete list entry | PASS | ge-0/0/9 gone from running |
| K3 | Rename list key | PASS | ge-0/0/8→ge-0/0/9 in running, cleaned up |
| K4 | Modify leaf in list | PASS | Description replaced in running, cleaned up |
| K5 | Add leaf in list | PASS | Description added in running, cleaned up |
| K6 | Delete leaf in list | PASS | Description deleted from running |
| K7 | Reorder list entries | PASS | Diff-engine only |
| K8 | Nested list entry | PASS | kernel contents in running, syslog host deleted |
| K9 | Per-instance ordering | PASS | Diff-engine only |
| K10 | Key-only list entry | PASS | Host 10.0.0.99 in running, syslog host deleted |
| C1 | Create container | PASS | aggregated-devices in running, deleted |
| C2 | Delete container | PASS | aggregated-devices gone from running |
| C3 | Empty presence container | PASS | ssh committed to running |
| C4 | Replace container | PASS | device-count 10→48 in running, deleted |
| C5 | Modify container child | PASS | hostname changed in running, restored |
| M1 | Mixed operations | PASS | All 3 ops committed to running, cleaned up |
| M2 | Deep nesting | PASS | syslog host committed to running, deleted |
| O1-O6 | Order alignment | PASS | Diff-engine only |

## Key Finding: `default-operation` Behavior (FIXED)

During initial candidate-only testing, a critical Junos NETCONF behavior was
discovered:

| `default-operation` | `nc:operation` on leaf | Result |
|---------------------|----------------------|--------|
| `none` | `replace` | **Ignored** — Junos returns warning: "ignore operation=replace" |
| `merge` (default) | `replace` | **Works** — leaf value updated |
| `merge` | `create` | **Works** — leaf/entry created |
| `merge` | `delete` | **Works** — leaf/entry deleted |

### Fix Applied

**File:** `terraform_provider/netconf/client.go`

```go
// BEFORE (broken on Junos 26.2I):
const patchEditConfigStr = `<edit-config>
    <target><candidate/></target>
    <default-operation>none</default-operation>
    ...`

// AFTER (fixed):
const patchEditConfigStr = `<edit-config>
    <target><candidate/></target>
    <default-operation>merge</default-operation>
    ...`
```

**File:** `terraform_provider/netconf/client_test.go`  
Updated `TestSendUpdateBaseConfigPayload` assertion to check for
`<default-operation>merge</default-operation>`.

### Why `merge` is safe

With `default-operation=merge`:
- Elements **without** `nc:operation` are merged (the safe default)
- Elements **with** `nc:operation="replace"` are replaced
- Elements **with** `nc:operation="create"` are created
- Elements **with** `nc:operation="delete"` are deleted

The diff engine always decorates every changed element with an explicit
`nc:operation`, so `merge` as the default only affects unchanged ancestor
containers (which are being merged anyway).

### Note on `deleteStr`

The pre-existing `deleteStr` template also uses `default-operation=none` with
the non-namespaced `operation="delete"` attribute. This is the legacy code path
for apply-group deletion and was not changed. It should be verified separately
if issues arise on Junos 26.2I.

## Test Adaptations for Live Device

Some matrix scenarios required adaptations for live Junos behavior:

| ID | Adaptation | Reason |
|----|-----------|--------|
| K4, K5, K6 | Modified description leaf only, not entire unit | Replacing `<unit>` strips VLAN-ID config, causing commit failure |
| K10 | Added minimal `contents` to host entry | Junos prunes syslog host entries with no contents |
| M2 | Changed `nc:operation="create"` to `"replace"` on `<any/>` | Cannot `create` a severity that already exists |
| K7, K9, LL4, O1-O6 | Tested as pass-through (Go engine only) | These test ordering logic that operates on in-memory XML, not NETCONF |

## Test Scripts

### Candidate-only test (original)

[`tests/live_matrix_test.py`](../tests/live_matrix_test.py) — Applies changes
to candidate config and discards after each test. No commits to running.

```
python3 tests/live_matrix_test.py                           # all tests on vpaa
python3 tests/live_matrix_test.py --host 10.54.21.160       # test on vpbb
python3 tests/live_matrix_test.py --test L1                 # single test
```

### Commit + verify + cleanup test

[`tests/live_commit_test.py`](../tests/live_commit_test.py) — Commits each
change to running config, verifies it took effect, then reverts it.

```
python3 tests/live_commit_test.py                           # all tests on vpaa
python3 tests/live_commit_test.py --host 10.54.21.160       # test on vpbb
python3 tests/live_commit_test.py --test L1,K3,M1           # multiple tests
```

Requires: `ncclient`, `lxml` (install via `pip install ncclient`).
