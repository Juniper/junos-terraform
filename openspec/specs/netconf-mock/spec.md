# NETCONF Mock Server Specification

Stateful multi-device NETCONF-over-SSH simulator for Terraform CI testing. Lives at `netconf_mock/`.

## Architecture

```
Terraform Provider (via SSH:830)
    ↓ SSH subsystem "netconf"
AsyncSSH Server (one TCP socket per device)
    ↓ framed messages (]]>]]> delimiter)
DeviceSession (RPC dispatch)
    ↓ mutates
DeviceState (candidate/running config lifecycle)
```

---

## DeviceState Model

```python
@dataclass
class DeviceState:
    name: str
    running_groups: dict[str, str]           # Committed config by group name
    candidate_groups: dict[str, str]         # Staged config (load/edit-config target)
    deleted_candidate_groups: set[str]       # Groups marked for deletion
    submitted_xml_by_group: dict[str, str]   # Last submitted XML per group
    rpc_log: list[str]                       # All received RPC payloads
    history: list[dict[str, str]]            # Structured operation log
```

### State Lifecycle

- **Given** a device starts, **When** `__post_init__` runs, **Then** `candidate_groups` is deep-copied from `running_groups`
- **Given** `running_groups` is empty at start, **When** device initializes, **Then** candidate is also empty
- **Given** `snapshot()` is called, **When** serialized, **Then** returns JSON-serializable dict of all fields

---

## SSH Authentication

### Behaviors

- **Given** valid `username` and `password`, **When** client connects, **Then** `validate_password` returns True
- **Given** invalid credentials, **When** client connects, **Then** `validate_password` returns False and connection rejected
- **Given** `disable_auth=True`, **When** `begin_auth()` runs, **Then** returns False (auth bypassed, all users accepted)
- **Given** a session is requested, **When** `session_requested()` runs, **Then** create `DeviceSession` bound to this device's state

---

## NETCONF Framing

- **Given** subsystem "netconf" requested, **When** `subsystem_requested()` runs, **Then** allow (return True)
- **Given** any other subsystem, **When** `subsystem_requested()` runs, **Then** reject (return False)
- **Given** subsystem starts, **When** `session_started()` runs, **Then** send NETCONF `<hello>` with capabilities and session-id
- **Given** data arrives, **When** `]]>]]>` delimiter found in buffer, **Then** split on delimiter and dispatch each complete RPC
- **Given** partial data without delimiter, **When** buffered, **Then** wait for more data

---

## RPC Dispatch

### Message-ID Extraction

- **Given** RPC has `message-id="123"`, **When** parsed, **Then** extract `"123"` via XML attribute lookup
- **Given** XML parse fails, **When** fallback regex used, **Then** extract message-id from attribute pattern
- **Given** no message-id attribute, **When** extraction attempted, **Then** return empty string

### Dispatch Order

Each incoming RPC is checked in this order:
1. `<edit-config>` with `nc:operation` attributes → patch mode
2. `<edit-config>` with group `operation="delete"` → group delete mode
3. `<load-configuration>` → load candidate
4. `<commit/>` → promote candidate to running
5. `<discard-changes/>` → restore candidate from running
6. `<get-configuration>` → return running config
7. Unknown → `<ok/>` reply

---

## load-configuration

### Behaviors

- **Given** XML payload with `<configuration><groups><name>GROUP</name>...</groups></configuration>`, **When** `_handle_load_configuration()` runs, **Then** store under `candidate_groups["GROUP"]`
- **Given** action is not specified (or "replace"), **When** groups extracted, **Then** overwrite group config entirely
- **Given** `action="merge"`, **When** groups extracted, **Then** merge incoming elements into existing group config using `_merge_group_configuration()`
- **Given** multiple groups in one payload, **When** processing, **Then** each group stored independently
- **Given** `<configuration-set>` with `set groups X ...` lines, **When** XML groups not found, **Then** extract set-style commands and store per group
- **Given** direct (non-group) configuration, **When** payload has `<configuration>` without `<groups>`, **Then** wrap in a groups structure using `_wrap_direct_configuration()`
- **Given** group previously deleted, **When** new load-configuration arrives for it, **Then** discard from `deleted_candidate_groups` set

---

## edit-config — Patch Mode

### Entry Conditions

- **Given** RPC contains `<edit-config>` AND `<config><configuration>` with any element having `nc:operation` attribute, **When** dispatched, **Then** enter patch mode
- **Given** `<edit-config>` is a whole-group `operation="delete"`, **When** dispatched, **Then** do NOT enter patch mode (handled separately)

### Patch Target Resolution

- **Given** patch elements reference a specific group, **When** resolving target, **Then** use that group
- **Given** no group specified, **When** `_resolve_patch_group_name()` runs, **Then**:
  1. If "base-config" exists in known groups, use it
  2. Else if only one group known, use it
  3. Else use alphabetically first group
  4. Else default to "base-config"

### Recursive Merge (`_apply_patch_element`)

- **Given** a patch element with `operation="delete"`, **When** matching target found, **Then** remove target from parent
- **Given** a patch element with `operation="delete"`, **When** no matching target found, **Then** no-op (idempotent)
- **Given** a patch element without operation (merge), **When** no matching target exists, **Then** create new subelement
- **Given** a leaf element (no children), **When** applied, **Then** set `target.text = patch_elem.text`
- **Given** a container element with children, **When** applied, **Then** recurse into each child

### Target Matching (`_find_matching_patch_target`)

| Given | When | Then |
|-------|------|------|
| Element has key children (e.g., `<name>`) | Searching candidates | Match by key values (`_find_candidate_by_keys`) |
| Leaf element (no children, no keys) | Searching candidates | Match by text value (`_find_leaf_value_candidate`) |
| Container with children but no keys | Single candidate exists | Use that candidate |
| Container with children but no keys | Multiple candidates exist | Return None (ambiguous) |

### Leaf-List Value Matching

- **Given** `operation="create"` and no existing leaf matches value, **When** searching, **Then** return None (triggers append)
- **Given** merge with multiple siblings of same tag (leaf-list pattern), **When** no value match, **Then** return None (append new entry)
- **Given** single candidate with same tag, **When** searching, **Then** return that candidate (overwrite)

### Keyed Entry Deletion

- **Given** patch element has key children AND a child with `operation="delete"` that matches the key, **When** `_deletes_matched_keyed_entry()` returns True, **Then** remove entire parent entry from config

---

## edit-config — Group Delete Mode

### Behaviors

- **Given** `<edit-config>` with `<groups operation="delete">`, **When** dispatched, **Then** add group name to `deleted_candidate_groups` and remove from `candidate_groups`
- **Given** `<apply-groups>` also present with `operation="delete"`, **When** processing, **Then** also record in history

---

## commit

### Behaviors

- **Given** `<commit/>` received, **When** candidate has changes, **Then**:
  1. Copy all `candidate_groups` to `running_groups`
  2. For each group in `deleted_candidate_groups`, remove from `running_groups`
  3. Clear `deleted_candidate_groups`
  4. Reply `<ok/>`
- **Given** `<commit/>` received, **When** candidate is unchanged, **Then** still reply `<ok/>` (idempotent)

---

## discard-changes

### Behaviors

- **Given** `<discard-changes/>` received, **When** candidate has modifications, **Then**:
  1. Deep-copy `running_groups` back to `candidate_groups`
  2. Clear `deleted_candidate_groups`
  3. Reply `<ok/>`

---

## get-configuration

### Behaviors

- **Given** `<get-configuration>` with group filter, **When** group exists in running, **Then** return `<configuration>` XML with that group's content
- **Given** `<get-configuration>` for non-existent group, **When** queried, **Then** return `<configuration/>`
- **Given** request asks for full (non-group) view, **When** `_full_configuration_for_group()` runs, **Then** strip `<groups>` wrapper and return inner config elements

---

## Multi-Device Support

- **Given** `--device dc1-leaf1:8301 --device dc1-leaf2:8302`, **When** server starts, **Then** one `DeviceState` and TCP listener per device
- **Given** each device, **When** SSH connections arrive on its port, **Then** operations are isolated to that device's state
- **Given** `--host 127.0.0.1`, **When** listeners bind, **Then** all devices listen on that IP

---

## CLI Arguments

```bash
python netconf_mock/netconf_mock_server.py \
  --host 127.0.0.1 \
  --username ci-user \
  --password ci-password \
  --device dc1-leaf1:8301 \
  --device dc1-leaf2:8302
```

| Argument | Required | Description |
|----------|----------|-------------|
| `--host` | Yes | Bind address for all device listeners |
| `--username` | Yes | SSH username for authentication |
| `--password` | Yes | SSH password for authentication |
| `--device NAME:PORT` | Yes (1+) | Device name and port pair |
| `--disable-auth` | No | Skip SSH authentication |

---

## Testing

### How to Run

```bash
# All mock server tests
pytest netconf_mock/tests/ -v

# Specific test file
pytest netconf_mock/tests/test_netconf_mock_server.py -v

# Single test
pytest netconf_mock/tests/test_netconf_mock_server.py::test_commit_promotes_candidate -v

# With coverage
pytest netconf_mock/tests/ --cov=netconf_mock --cov-report=term-missing
```

### Test Coverage Map

| Test Area | Given | When | Then |
|-----------|-------|------|------|
| Load config | XML with groups payload | `_handle_load_configuration` | Group stored in candidate |
| Load merge | `action="merge"` with existing config | Load runs | Elements merged not replaced |
| Edit-config patch | Patch with `operation="delete"` | Patch applied | Target leaf removed |
| Edit-config patch | Patch with `operation="create"` | Patch applied | New leaf appended |
| Edit-config patch | Nested patch elements | Patch applied recursively | Deep merge correct |
| Group delete | `<groups operation="delete">` | Edit-config dispatch | Group removed from candidate |
| Commit | Candidate modified | `<commit/>` sent | Running updated from candidate |
| Discard | Candidate modified | `<discard-changes/>` sent | Candidate reset to running |
| Get-config | Group exists in running | `<get-configuration>` | Correct XML returned |
| Get-config | Group doesn't exist | `<get-configuration>` | Empty config returned |
| Multi-device | Two devices started | Independent RPCs sent | States isolated |
| Auth | Wrong password | SSH connect | Connection rejected |
| Framing | Partial message | Data buffered | No dispatch until delimiter |
| Message-ID | Valid message-id attr | Extraction | Correct ID in reply |
| Patch target | Key-based matching | `_find_matching_patch_target` | Correct entry found |
| Leaf-list | Create new value | Patch applied | Value appended (not overwritten) |

### Integration Test Pattern

```python
import asyncio
from netconf_mock.netconf_mock_server import DeviceState, DeviceSession

def test_commit_promotes_candidate():
    state = DeviceState(name="test-dev")
    state.candidate_groups["my-group"] = "<configuration><groups><name>my-group</name><interfaces/></groups></configuration>"

    # Simulate commit
    state.running_groups = copy.deepcopy(state.candidate_groups)
    state.deleted_candidate_groups.clear()

    assert "my-group" in state.running_groups
    assert state.running_groups["my-group"] == state.candidate_groups["my-group"]
```
