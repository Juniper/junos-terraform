# Tasks

## 1. Schema Loading & Embedding

- [x] 1.1 Create `terraform_provider/generic/embed.go` with `go:embed` for `trimmed_schema.json.gz` (or raw fallback)
- [x] 1.2 Implement `LoadSchema()` function that decompresses and parses JSON into `patch.ProcessSchema` index
- [x] 1.3 Add startup initialization in provider Configure that calls LoadSchema once and stores the index

## 2. Dynamic Schema Builder

- [x] 2.1 Create `terraform_provider/generic/schema_builder.go` with `BuildSchema(nodes []patch.SchemaNode) schema.Schema`
- [x] 2.2 Implement recursive `buildAttributes(children []SchemaNode) map[string]schema.Attribute` (leaf→String, leaf-list→List, container/list→ListNested)
- [x] 2.3 Add name sanitization (dashes/dots → underscores) and include `resource_name` required attribute
- [x] 2.4 Write unit tests: single leaf, nested container, leaf-list, name sanitization

## 3. Plan-to-XML Bridge

- [x] 3.1 Create `terraform_provider/generic/plan_xml_bridge.go` with `PlanToXML(ctx, plan tftypes.Value, index map[string]NodeInfo) ([]byte, error)`
- [x] 3.2 Implement recursive tree walk: read plan attributes using schema index, emit XML elements with original YANG names
- [x] 3.3 Handle null/unset attributes (skip), leaf-lists (repeated elements), and containers (recurse)
- [x] 3.4 Write unit tests: simple leaf plan→XML, nested containers, null handling

## 4. XML-to-State Bridge

- [x] 4.1 Create `terraform_provider/generic/xml_state_bridge.go` with `XMLToState(ctx, xmlBytes []byte, index map[string]NodeInfo) (tftypes.Value, error)`
- [x] 4.2 Implement recursive XML parse: build `tftypes.Value` tree from XML elements using schema for type info
- [x] 4.3 Handle missing elements (→ null), repeated elements (→ list), and nested containers
- [x] 4.4 Write unit tests: device XML → state values, missing leaves → null, leaf-lists

## 5. Generic CRUD Resource

- [x] 5.1 Create `terraform_provider/generic/resource.go` implementing `resource.Resource` interface
- [x] 5.2 Implement `Schema()` using schema builder output
- [x] 5.3 Implement `Create()`: PlanToXML → SendDirectTransaction → commit → ReadStateFromDevice
- [x] 5.4 Implement `Read()`: get-config → XMLToState → set state
- [x] 5.5 Implement `Update()`: PlanToXML + StateToXML → ComputeDiff → CreateDiffPatch → edit-config → commit
- [x] 5.6 Implement `Delete()`: StateToXML vs empty → ComputeDiff → CreateDiffPatch (deletes) → edit-config → commit
- [x] 5.7 Wire resource into `provider.go` Resources() list

## 6. jtaf-provider CLI Update

- [x] 6.1 Modify `jtaf-provider` to copy generic provider Go source instead of rendering Jinja2 template
- [x] 6.2 Emit `trimmed_schema.json` (and optionally .gz) into output directory
- [x] 6.3 Update `go.mod.j2` template (or generate go.mod) to reference generic package
- [x] 6.4 Verify end-to-end: `jtaf-yang2go -p ... -t vqfx` → output dir → `go build .` succeeds

## 7. Integration Testing

- [x] 7.1 Build provider with small filtered schema (existing EVPN test XMLs) and verify `terraform plan` works
- [x] 7.2 Build provider with full QFX 18.2 schema (no -x) and verify compile + `terraform init` succeeds
- [x] 7.3 Test CRUD cycle against mock NETCONF server with generic provider
- [ ] 7.4 Add CI workflow step that builds generic provider and runs `terraform validate`
