# Tasks

## 1. Full model loading
- [x] 1.1 Accept numeric-string range bounds (`patch.NumRange`)
- [x] 1.2 Lowercase attribute names; `ValidateNames` rejects invalid or colliding names
- [x] 1.3 Report schema load errors as diagnostics
- [x] 1.4 Flatten `choice`/`case`; build the index from parsed nodes (`patch.FlattenChoices`, `patch.BuildSchemaIndex`)

## 2. Code generation
- [x] 2.1 `jtaf-provider --exclude PATH`, passed on by `jtaf-yang2go`
- [x] 2.2 Embed the schema compact and gzipped
- [x] 2.3 Generate `main.go` calling `generic.Serve`
- [x] 2.4 Always drop the top-level `version` (the release, not configuration), as trimming does

## 3. Converters
- [x] 3.1 `ValueToConfig`, `ConfigToValue` with list keys; remove the flat-path bridges

## 4. Provider
- [x] 4.1 Create/Read/Update/Delete against the device
- [x] 4.2 Serve tfprotov6 with terraform-plugin-go
- [x] 4.3 Unit tests with a fake NETCONF client
- [x] 4.4 Take an update's state from the configuration it already read

## 5. Verification
- [x] 5.1 Plan against an SRX300 (26.2R1.7) with a full model: no changes, existing state loads
- [x] 5.2 Apply a change against a device: an SRX300, adding and removing an interface description, one commit each
- [ ] 5.3 CI step building a generic provider (go-generator-replacement 7.4)
