# Tasks

## 1. Unify on the schema-driven generator

- [x] 1.1 Make `jtaf-provider` without `--generic` run the same generation as `--generic`, differing only in that the
      schema is trimmed with `-x` first; verify `jtaf-provider -j <schema> -x <xml> -t t1` produces `main.go`,
      `embed_schema.go`, `schema.bin.gz` and `trimmed_schema.json.gz`, and no `resource_config_provider.go`
- [x] 1.2 Delete `junosterraform/templates/resource_config_provider.go.j2`, `provider.go.j2` and `config.go.j2` plus
      the rendering helpers only they used; verify `grep -r "resource_config_provider.go.j2" junosterraform` returns
      nothing and `pytest junosterraform/tests/` passes after updating `test_script_coverage.py` and `test_workflow.py`
- [x] 1.3 Reject `--generic` with `-x` with a message naming the conflict, and update the `--generic` help text to say
      it embeds the untrimmed model; verify both behaviours with a test in `junosterraform/tests/test_script_coverage.py`
- [x] 1.4 Assert the two invocations produce identical Go source: add a test generating a provider with `-x` and one
      with `--generic` from the same model and comparing every `.go` file byte for byte except `embed_schema.go`
- [x] 1.5 Collapse `examples/providers/test-both-providers.sh` and `build.sh`/`build-generic.sh` to the single
      generator, keeping a trimmed and an untrimmed build; verify both scripts run to a compiled binary
- [x] 1.6 Record the breaking change in `changelog.md` and update `openspec/specs`-adjacent docs (`README-terraform.md`,
      `terraform_provider/README.md`) that describe template rendering; verify the documented commands run as written

## 2. Groups in the schema

- [x] 2.1 Add `--groups` to `jtaf-provider`, defaulting to excluding the `groups` subtree and the top-level
      `apply-groups` leaf-list the way `drop_version()` drops `version`; verify a test asserts neither name appears in
      `trimmed_schema.json.gz` without the flag and both appear with it
- [x] 2.2 Reject `--groups` combined with `--exclude groups` with a message naming the conflict; verify with a test
- [x] 2.3 Mark the `apply-groups` leaf-list as ordered in the schema it produces so the patch engine's `ordered-by user`
      handling applies; verify a test reads the compiled schema and asserts the ordered flag
- [x] 2.4 Stop `jtaf_common.filter_json_using_xml()` from removing `apply-groups`, and let `<groups>` children
      contribute XPaths, when groups are requested; verify a test trims a schema against group-based XML and asserts
      the group body is present and the unused parts of the group subtree are not
- [x] 2.5 Add `--groups` pass-through to `jtaf-yang2go`; verify a test asserts the flag reaches the `jtaf-provider`
      argument list, mirroring the existing `--exclude` pass-through test
- [x] 2.6 Document `--groups` in `README-terraform.md` and the tools' help text, including that
      `--generic --groups` embeds the group subtree in full; verify the help text matches the documentation

## 3. Groups in the provider

- [x] 3.1 Remove the `groups`/`apply-groups` skip from `attributeNodes()` in `terraform_provider/generic/names.go` so
      attribute listing is purely structural; verify `go test ./generic/...` passes with a test building a schema that
      contains `groups` and asserting the resource exposes `groups` and `apply_groups`
- [x] 3.2 Verify a schema without `groups` yields a resource with no such attribute and no error, with a test
- [x] 3.3 Confirm the read path requests the committed configuration without inheritance and add a test asserting the
      `get-configuration` RPC carries no `inherit` attribute
- [x] 3.4 Delete the legacy group bookkeeping in `terraform_provider/netconf/client.go` (`groupStrXML`,
      `getGroupXMLStr`, `applyGroupsList`, `applyGroupsMutex`, `sendApplyGroupsLocked`, `addToApplyGroupsList`,
      `sortApplyGroupsList`, `MarshalGroup` and their call sites); verify `go build ./...` and `go test ./...` pass
- [x] 3.5 Add a provider-level test asserting a group declared without an `apply_groups` entry is sent without any
      reference being synthesised

## 4. Patch engine

- [x] 4.1 Add leaf-map and diff tests for a single leaf changing inside `groups[name=X]` and assert the emitted edit
      targets only that leaf below its `<groups>` element
- [x] 4.2 Add a test for the same path in two groups with different values and assert changing one leaves the other
      untouched, and that neither affects the identically named base-hierarchy path
- [x] 4.3 Make removing a managed group emit one container-level delete rather than per-leaf deletes; verify with a
      coalescing test in `terraform_provider/patch/`
- [x] 4.4 Add ordered-leaf-list tests for `apply-groups`: reordering produces a non-empty diff that results in the
      declared order, an unchanged order produces an empty diff, and a removed entry is deleted
- [x] 4.5 Add a test with a wildcard group key (`<name>ge-*</name>`) asserting the key survives leaf-map encoding and
      round-trips through the emitted patch

## 5. XML to Terraform conversion

- [ ] 5.1 Add `--groups` to `jtaf-xml2tf` so it keeps `<groups>` and `<apply-groups>` instead of flattening applied
      groups into the base hierarchy; verify a test in `junosterraform/tests/test_xml2tf_flatten.py` asserts the
      preserved output contains a `groups` block and an `apply_groups` list and copies nothing into the base hierarchy
- [ ] 5.2 Verify the default conversion output is unchanged, with a test comparing against the current flattened result
- [ ] 5.3 Document the flag in `README-terraform.md` alongside the provider flag; verify the documented command runs

## 6. Mock NETCONF server

- [ ] 6.1 Change the mock to hold one configuration tree with `groups` inside it instead of blob state keyed by group
      name, so patches can edit inside a group; verify `pytest netconf_mock/tests/` passes with a new test applying an
      `edit-config` to a leaf under `<groups>`
- [ ] 6.2 Add mock tests for deleting a group and for `apply-groups` ordering, asserting the committed tree matches what
      a device would hold

## 7. Group-based example

- [ ] 7.1 Add `examples/evpn-vxlan-dc-groups/` mirroring `examples/evpn-vxlan-dc/` device for device, with each device's
      non-identity configuration inside one role-named group and the base hierarchy holding host name, management
      addressing and the `apply-groups` reference; verify each XML parses and the device set matches the base example
- [ ] 7.2 Add `examples/providers/build-groups.sh` and `convert-groups.sh` using `--groups`; verify they produce a
      compiled provider whose schema contains `groups` and `.tf` files containing a `groups` block per device
- [ ] 7.3 Add the generated `.tf` output under its own directory alongside `examples/terraform_files/`; verify
      `terraform validate` passes against the group-based provider
- [ ] 7.4 Document the example in `examples/DEMO-GENERIC-PROVIDER.md` or a sibling document, including the build,
      convert and apply commands; verify the documented commands run as written

## 8. End-to-end verification

- [ ] 8.1 Run the group-based example end to end against the mock: build, convert, apply, then plan again and assert no
      changes
- [ ] 8.2 Change a leaf inside a group on the mock out of band, then plan and apply, asserting the drift is reported and
      restored
- [ ] 8.3 Remove a group from the configuration, apply, and assert the group is deleted while unmanaged groups and the
      base hierarchy are untouched
- [ ] 8.4 Run the existing base-configuration example end to end and assert its behaviour and output are unchanged
- [ ] 8.5 Wire the group-based run into the Terraform end-to-end workflow with the same bounded parallelism the existing
      runs use; verify the workflow passes
