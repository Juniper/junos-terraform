# Prompt for validating issue #168 on a NITA instance

Paste the block below into the coding agent once you are on the NITA box with
this branch checked out. Fill in the device names/IPs first.

---

You are helping validate Juniper/junos-terraform issue #168, which proposes
removing the build-time `jtaf-provider --yang-compliant` flag and instead
having the generic provider normalise YANG choices at runtime: accept either
encoding on read, and **always write the paired form**
(`<choice-ident>exact</choice-ident><choice-value/>`) rather than the native
form (`<exact/>`).

That proposal rests on one assumption: every Junos platform and release accepts
a paired-form write and means exactly the same thing by it as by the native
form. It was originally tested on a single vMX-class device at 26.2R1.7. It has
since been confirmed on a physical MX240 at 21.1R1.11, a vMX at 23.2R2.21 and
an MX960 at 25.2R1.9 — but **only on the MX family**. The open question is
whether it also holds on SRX and PTX, where routing-policy handling may differ.

The probe is already written: `tools/choice-probe/choice_probe.py`.

Please do the following.

1. Set up:
   ```bash
   cd tools/choice-probe
   python3 -m venv .venv && . .venv/bin/activate
   pip install -r requirements.txt
   ```

2. Create `devices.json` (gitignored) from `devices.example.json` with the NITA
   devices. Include at least one **vSRX** and one **vPTX**; add vMX as a control
   that should reproduce the known-good result:
   ```json
   {
     "vsrx1": {"host": "<IP>", "username": "<USER>", "password": "<PASS>"},
     "vptx1": {"host": "<IP>", "username": "<USER>", "password": "<PASS>"},
     "vmx1":  {"host": "<IP>", "username": "<USER>", "password": "<PASS>"}
   }
   ```
   You can leave `password` empty and export `JTAF_PROBE_PASSWORD` instead.

3. Run the read-only probe first. It commits nothing: every case is loaded into
   a locked candidate and discarded, and it refuses to touch a device whose
   candidate already holds uncommitted changes.
   ```bash
   python choice_probe.py --devices devices.json
   ```

4. Then, **only on lab devices that are safe to commit on**, run the knob stage.
   It commits `system services netconf yang-compliant` using `commit confirmed`
   (auto-reverts) and then reverts explicitly:
   ```bash
   python choice_probe.py --devices devices.json --targets vsrx1 vptx1 --commit-knob
   ```

5. Afterwards, confirm each device is clean: knob off, no `JTAF-CHOICE-PROBE`
   policy in the committed config, no pending candidate changes.

Report back, per device, model and Junos version plus:

- is the `yang-compliant` knob supported?
- are **paired** writes accepted for route-filter `exact`,
  `prefix-length-range`, `orlonger`+trailing action, `upto`, `community add`
  and `martians`?
- do native and paired produce an **identical** candidate diff? If any case
  differs, paste both diffs — that is the result that would block the proposal.
- read encoding with the knob off, and with it on.
- with the knob on, is a delete naming only the address rejected, and one
  carrying the full compound key accepted?

Finally, state plainly whether issue #168's "always write the paired form" rule
holds on SRX and PTX, or whether it must become platform-aware. Flag any device
left in a modified state.

Do not modify anything under `terraform_provider/` or `junosterraform/`; this is
a validation run only.
