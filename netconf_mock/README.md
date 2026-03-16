# NETCONF Mock Server

`netconf_mock_server.py` is a stateful, multi-device NETCONF-over-SSH simulator used for Terraform/provider integration tests.

It is designed to emulate the core Junos config flow needed by this repo:
- Per-device listeners on independent ports
- Candidate vs running config lifecycle
- `load-configuration`, `edit-config` delete, `commit`, `discard-changes`
- `get-configuration` responses by group
- Optional per-device state dump for post-test debugging

## Requirements

- Python 3.10+
- `asyncssh`

Install dependency:

```bash
python -m pip install asyncssh
```

## Quick Start

From repository root:

```bash
python netconf_mock/netconf_mock_server.py \
  --host 127.0.0.1 \
  --username ci-user \
  --password ci-password \
  --disable-auth \
  --device dc1-leaf1:8301 \
  --device dc1-leaf2:8302 \
  --log-level INFO
```

The process runs until it receives `SIGINT`/`SIGTERM` (for example `Ctrl+C`).

## Using a Devices File

Instead of repeating `--device`, provide a file with one entry per line:

```text
dc1-leaf1:8301
dc1-leaf2:8302
# comments are allowed
```

Run with:

```bash
python netconf_mock/netconf_mock_server.py \
  --host 127.0.0.1 \
  --username ci-user \
  --password ci-password \
  --devices-file ci-evpn-vxlan-tf/mock-devices.txt
```

## CLI Options

- `--host`: bind address for all listeners (default: `127.0.0.1`)
- `--username`: SSH username required for auth (default: `ci-user`)
- `--password`: SSH password required for auth (default: `ci-password`)
- `--device`: `<device-name>:<port>` listener spec; may be repeated
- `--devices-file`: file with one `<device-name>:<port>` per line
- `--disable-auth`: disable SSH auth checks (useful for mock-only compatibility runs)
- `--state-dump`: optional JSON file written on shutdown with device state/history
- `--log-level`: logging level (`DEBUG`, `INFO`, `WARNING`, `ERROR`)

At least one of `--device` or `--devices-file` is required.

## State Model

Each device tracks:
- `running_groups`
- `candidate_groups`
- `submitted_xml_by_group`
- `rpc_log`
- `history`

Behavior:
- candidate starts as a deep copy of running
- `load-configuration` updates candidate for the target group
- `edit-config` with `operation="delete"` removes group from candidate
- `commit` copies candidate to running
- `discard-changes` resets candidate from running
- `get-configuration` returns running config for requested group

## CI Usage Pattern

Typical integration sequence:

```bash
python netconf_mock/netconf_mock_server.py \
  --host 127.0.0.1 \
  --username ci-user \
  --password ci-password \
  --devices-file ci-evpn-vxlan-tf/mock-devices.txt \
  --state-dump netconf-mock-state.json \
  --log-level DEBUG > netconf-mock.log 2>&1 &

echo $! > .netconf_mock.pid
```

Later, stop it cleanly:

```bash
kill "$(cat .netconf_mock.pid)"
```

## Troubleshooting

- Auth failures: verify provider credentials match `--username` and `--password`.
- Connection failures: verify device ports in `--device`/`--devices-file` and listener readiness.
- Unexpected behavior: run with `--log-level DEBUG` and inspect `netconf-mock.log`.
- Need post-run evidence: set `--state-dump` and inspect the JSON snapshot.
