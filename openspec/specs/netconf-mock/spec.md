# NETCONF Mock Server

Python-based multi-device NETCONF-over-SSH simulator for testing providers without live hardware. Lives at `netconf_mock/`.

## Purpose

Simulates Junos NETCONF endpoints so the Terraform provider can be tested end-to-end without real devices. Each mock device maintains independent candidate/running configuration state.

## Architecture

```
Terraform Provider
    ↓ (SSH/NETCONF)
netconf_mock_server.py
    ↓
Per-device state:
  ├── running_groups (committed config)
  └── candidate_groups (staged config)
```

## Files

| File | Purpose |
|------|---------|
| `netconf_mock_server.py` | Async SSH server with NETCONF subsystem, per-device state |
| `tests/test_netconf_mock_server.py` | Mock server behavior validation |

## Supported RPCs

| RPC | Behavior |
|-----|----------|
| `<load-configuration>` | Stores config XML into candidate groups |
| `<edit-config operation="delete">` | Marks groups for deletion in candidate |
| `<commit>` | Applies candidate to running state |
| `<discard-changes>` | Resets candidate to running state |
| `<get-configuration>` | Returns running config XML |

## Features

- **Multi-device** — Each device gets its own SSH port and independent state
- **Group-based config** — Tracks configuration by group name (matches Junos apply-groups)
- **State snapshots** — Dumps state for post-test debugging
- **Auth bypass** — Optional authentication disable for CI environments
- **Candidate/running lifecycle** — Full two-stage commit model matching real Junos

## Configuration

Devices are configured as a list, each with:
- `name` — Device identifier
- `port` — SSH listen port
- `initial_config` — Optional starting XML configuration

## Usage

```python
# Start mock server
from netconf_mock.netconf_mock_server import NetconfMockServer

server = NetconfMockServer(devices=[
    {"name": "vpaa", "port": 10830, "initial_config": "<configuration>...</configuration>"}
])
await server.start()

# Provider connects to localhost:10830 as if it were a real device
```

## Tests

- `netconf_mock/tests/test_netconf_mock_server.py` — Device state management, RPC handling, multi-device scenarios
