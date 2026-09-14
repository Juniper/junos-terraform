## ADDED Requirements

### Requirement: Build Terraform schema from JSON at startup
The provider SHALL recursively walk `trimmed_schema.json` at startup and produce a `schema.Schema` object that Terraform can use for plan validation and state management.

#### Scenario: Leaf node becomes StringAttribute
- **WHEN** the schema JSON contains a node with `"type": "leaf"`
- **THEN** the builder SHALL produce a `schema.StringAttribute{Optional: true}` with the node name (dashes replaced with underscores) as the attribute key

#### Scenario: Leaf-list node becomes ListAttribute of strings
- **WHEN** the schema JSON contains a node with `"type": "leaf-list"`
- **THEN** the builder SHALL produce a `schema.ListAttribute{ElementType: types.StringType, Optional: true}`

#### Scenario: Container node becomes ListNestedAttribute
- **WHEN** the schema JSON contains a node with `"type": "container"` and children
- **THEN** the builder SHALL produce a `schema.ListNestedAttribute` whose `NestedObject.Attributes` map is recursively built from the container's children

#### Scenario: List node becomes ListNestedAttribute
- **WHEN** the schema JSON contains a node with `"type": "list"` and children
- **THEN** the builder SHALL produce a `schema.ListNestedAttribute` whose `NestedObject.Attributes` map is recursively built from the list's children

#### Scenario: Top-level configuration children
- **WHEN** the schema JSON has `root.children[0]` as the configuration node
- **THEN** the builder SHALL produce one `ListNestedAttribute` per direct child of configuration (excluding `groups` and `apply-groups`)

#### Scenario: Resource name attribute
- **WHEN** the schema is built
- **THEN** the schema SHALL include a `resource_name` StringAttribute with `Required: true` and `PlanModifiers: [stringplanmodifier.RequiresReplace()]`

### Requirement: Handle name sanitization
The builder SHALL convert YANG node names to valid Terraform attribute names by replacing dashes and dots with underscores.

#### Scenario: Dash in name
- **WHEN** a YANG node is named `apply-groups`
- **THEN** the Terraform attribute key SHALL be `apply_groups`

#### Scenario: Dot in name
- **WHEN** a YANG node is named `802.1x`
- **THEN** the Terraform attribute key SHALL be `802_1x`
