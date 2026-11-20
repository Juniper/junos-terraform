## ADDED Requirements

### Requirement: Convert Terraform plan to Junos XML
The provider SHALL convert a Terraform plan value into a valid Junos `<configuration>` XML document by walking the plan's `tftypes.Value` tree using the schema index.

#### Scenario: Leaf value to XML element
- **WHEN** the plan contains a string value at path `interfaces/interface/name` with value `"ge-0/0/0"`
- **THEN** the output XML SHALL contain `<name>ge-0/0/0</name>` nested under `<interface>` under `<interfaces>`

#### Scenario: Leaf-list values to repeated XML elements
- **WHEN** the plan contains a list of strings at path `interfaces/interface/apply-groups` with values `["GRP1", "GRP2"]`
- **THEN** the output XML SHALL contain `<apply-groups>GRP1</apply-groups><apply-groups>GRP2</apply-groups>`

#### Scenario: Null/unset attributes are omitted
- **WHEN** the plan contains a null value for an attribute
- **THEN** the output XML SHALL NOT include any element for that attribute

#### Scenario: Container nesting produces correct XML hierarchy
- **WHEN** the plan has nested objects (container within container)
- **THEN** the output XML SHALL produce properly nested elements matching the YANG hierarchy

### Requirement: Convert device XML to Terraform state
The provider SHALL convert a Junos `<configuration>` XML response into a Terraform state `tftypes.Value` tree by walking the XML using the schema index.

#### Scenario: XML leaf to state string
- **WHEN** the device XML contains `<name>ge-0/0/0</name>` at path `interfaces/interface/name`
- **THEN** the state value for that attribute SHALL be `types.StringValue("ge-0/0/0")`

#### Scenario: Missing XML element becomes null in state
- **WHEN** the device XML does not contain an element for a schema-defined leaf
- **THEN** the state value for that attribute SHALL be `types.StringNull()`

#### Scenario: Repeated XML elements become list in state
- **WHEN** the device XML contains multiple `<apply-groups>` elements
- **THEN** the state value SHALL be a `types.List` containing each element's text as a string

#### Scenario: Attribute name uses underscores in state
- **WHEN** the device XML contains element `<apply-groups>`
- **THEN** the state attribute key SHALL be `apply_groups` (YANG name with dashes converted)

### Requirement: Use original YANG names for XML tag generation
The bridge SHALL use the original YANG node names (with dashes/dots) for XML element tags, not the sanitized Terraform attribute names.

#### Scenario: Terraform attribute back to XML
- **WHEN** converting plan attribute `apply_groups` to XML
- **THEN** the XML tag SHALL be `<apply-groups>` (reverse name mapping via schema index)
