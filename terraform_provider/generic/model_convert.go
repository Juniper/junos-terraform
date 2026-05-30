package generic

import (
	"context"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ModelToXMLBytes converts Terraform plan attributes into XML config bytes.
// It walks the schema tree and extracts values from the plan object map.
func ModelToXMLBytes(ctx context.Context, attrs map[string]attr.Value, idx *SchemaIndex, diags *diag.Diagnostics) []byte {
	root := &patch.Node{
		Tag:   "configuration",
		Attrs: make(map[string]string),
	}

	for _, schemaNode := range idx.TopLevel {
		name := normalizeName(schemaNode.Name)
		val, ok := attrs[name]
		if !ok || val.IsNull() || val.IsUnknown() {
			continue
		}
		children := modelValueToNodes(ctx, val, schemaNode, diags)
		if diags.HasError() {
			return nil
		}
		root.Children = append(root.Children, children...)
	}

	xmlBytes, err := SerializeTree(root)
	if err != nil {
		diags.AddError("Failed to serialize XML", err.Error())
		return nil
	}
	return xmlBytes
}

// XMLBytesToModel converts device XML config bytes into Terraform state attribute values.
func XMLBytesToModel(ctx context.Context, xmlBytes []byte, idx *SchemaIndex, diags *diag.Diagnostics) map[string]attr.Value {
	tree, err := patch.BuildTree(xmlBytes)
	if err != nil {
		diags.AddError("Failed to parse XML", err.Error())
		return nil
	}

	result := make(map[string]attr.Value, len(idx.TopLevel))
	for _, schemaNode := range idx.TopLevel {
		name := normalizeName(schemaNode.Name)
		result[name] = xmlNodeToModelValue(ctx, tree, schemaNode, diags)
		if diags.HasError() {
			return nil
		}
	}
	return result
}

// modelValueToNodes converts a single Terraform attribute value into XML Node(s)
// based on the schema node type.
func modelValueToNodes(ctx context.Context, val attr.Value, schema patch.SchemaNode, diags *diag.Diagnostics) []*patch.Node {
	switch schema.Type {
	case "leaf":
		return leafToNodes(val, schema)
	case "leaf-list":
		return leafListToNodes(ctx, val, schema, diags)
	case "container":
		return containerToNodes(ctx, val, schema, diags)
	case "list":
		return listToNodes(ctx, val, schema, diags)
	default:
		return leafToNodes(val, schema)
	}
}

func leafToNodes(val attr.Value, schema patch.SchemaNode) []*patch.Node {
	if val.IsNull() || val.IsUnknown() {
		return nil
	}
	sv, ok := val.(basetypes.StringValue)
	if !ok {
		return nil
	}
	node := &patch.Node{
		Tag:   schema.Name,
		Attrs: make(map[string]string),
		Text:  sv.ValueString(),
	}
	return []*patch.Node{node}
}

func leafListToNodes(ctx context.Context, val attr.Value, schema patch.SchemaNode, diags *diag.Diagnostics) []*patch.Node {
	if val.IsNull() || val.IsUnknown() {
		return nil
	}
	lv, ok := val.(basetypes.ListValue)
	if !ok {
		return nil
	}
	var nodes []*patch.Node
	for _, elem := range lv.Elements() {
		sv, ok := elem.(basetypes.StringValue)
		if !ok {
			continue
		}
		nodes = append(nodes, &patch.Node{
			Tag:   schema.Name,
			Attrs: make(map[string]string),
			Text:  sv.ValueString(),
		})
	}
	return nodes
}

func containerToNodes(ctx context.Context, val attr.Value, schema patch.SchemaNode, diags *diag.Diagnostics) []*patch.Node {
	if val.IsNull() || val.IsUnknown() {
		return nil
	}

	// Container is modeled as ListNestedAttribute (max 1 element) for backward
	// compatibility with generated .tf files that use list syntax: attr = [{...}]
	lv, ok := val.(basetypes.ListValue)
	if !ok {
		// May be a string presence marker for empty containers
		if sv, ok := val.(basetypes.StringValue); ok && !sv.IsNull() {
			return []*patch.Node{{
				Tag:   schema.Name,
				Attrs: make(map[string]string),
			}}
		}
		return nil
	}

	elems := lv.Elements()
	if len(elems) == 0 {
		return nil
	}

	ov, ok := elems[0].(basetypes.ObjectValue)
	if !ok {
		return nil
	}

	node := &patch.Node{
		Tag:   schema.Name,
		Attrs: make(map[string]string),
	}

	objAttrs := ov.Attributes()
	for _, child := range schema.Children {
		childName := normalizeName(child.Name)
		childVal, exists := objAttrs[childName]
		if !exists || childVal.IsNull() || childVal.IsUnknown() {
			continue
		}
		childNodes := modelValueToNodes(ctx, childVal, child, diags)
		if diags.HasError() {
			return nil
		}
		node.Children = append(node.Children, childNodes...)
	}

	return []*patch.Node{node}
}

func listToNodes(ctx context.Context, val attr.Value, schema patch.SchemaNode, diags *diag.Diagnostics) []*patch.Node {
	if val.IsNull() || val.IsUnknown() {
		return nil
	}

	lv, ok := val.(basetypes.ListValue)
	if !ok {
		return nil
	}

	var nodes []*patch.Node
	for _, elem := range lv.Elements() {
		ov, ok := elem.(basetypes.ObjectValue)
		if !ok {
			continue
		}

		node := &patch.Node{
			Tag:   schema.Name,
			Attrs: make(map[string]string),
		}

		objAttrs := ov.Attributes()
		for _, child := range schema.Children {
			childName := normalizeName(child.Name)
			childVal, exists := objAttrs[childName]
			if !exists || childVal.IsNull() || childVal.IsUnknown() {
				continue
			}
			childNodes := modelValueToNodes(ctx, childVal, child, diags)
			if diags.HasError() {
				return nil
			}
			node.Children = append(node.Children, childNodes...)
		}

		nodes = append(nodes, node)
	}
	return nodes
}

// xmlNodeToModelValue converts XML tree children into a Terraform attribute value
// for a given schema node.
func xmlNodeToModelValue(ctx context.Context, xmlRoot *patch.Node, schema patch.SchemaNode, diags *diag.Diagnostics) attr.Value {
	switch schema.Type {
	case "leaf":
		return xmlLeafToValue(xmlRoot, schema)
	case "leaf-list":
		return xmlLeafListToValue(ctx, xmlRoot, schema)
	case "container":
		return xmlContainerToValue(ctx, xmlRoot, schema, diags)
	case "list":
		return xmlListToValue(ctx, xmlRoot, schema, diags)
	default:
		return xmlLeafToValue(xmlRoot, schema)
	}
}

func findChildrenByTag(parent *patch.Node, tag string) []*patch.Node {
	var found []*patch.Node
	for _, child := range parent.Children {
		if child.Tag == tag {
			found = append(found, child)
		}
	}
	return found
}

func xmlLeafToValue(xmlParent *patch.Node, schema patch.SchemaNode) attr.Value {
	matches := findChildrenByTag(xmlParent, schema.Name)
	if len(matches) == 0 {
		return types.StringNull()
	}
	// YANG empty type: presence element with no text
	return types.StringValue(matches[0].Text)
}

func xmlLeafListToValue(ctx context.Context, xmlParent *patch.Node, schema patch.SchemaNode) attr.Value {
	matches := findChildrenByTag(xmlParent, schema.Name)
	if len(matches) == 0 {
		return types.ListNull(types.StringType)
	}
	elems := make([]attr.Value, 0, len(matches))
	for _, m := range matches {
		elems = append(elems, types.StringValue(m.Text))
	}
	lv, _ := types.ListValue(types.StringType, elems)
	return lv
}

func xmlContainerToValue(ctx context.Context, xmlParent *patch.Node, schema patch.SchemaNode, diags *diag.Diagnostics) attr.Value {
	attrTypes := containerAttrTypes(schema)
	matches := findChildrenByTag(xmlParent, schema.Name)

	if len(schema.Children) == 0 {
		// Empty container presence → empty string marker
		if len(matches) == 0 {
			return types.StringNull()
		}
		return types.StringValue("")
	}

	objType := types.ObjectType{AttrTypes: attrTypes}

	if len(matches) == 0 {
		return types.ListNull(objType)
	}

	xmlNode := matches[0]

	objAttrs := make(map[string]attr.Value, len(schema.Children))
	for _, child := range schema.Children {
		childName := normalizeName(child.Name)
		objAttrs[childName] = xmlNodeToModelValue(ctx, xmlNode, child, diags)
		if diags.HasError() {
			return types.ListNull(objType)
		}
	}

	ov, d := types.ObjectValue(attrTypes, objAttrs)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType)
	}

	lv, d := types.ListValue(objType, []attr.Value{ov})
	diags.Append(d...)
	return lv
}

func xmlListToValue(ctx context.Context, xmlParent *patch.Node, schema patch.SchemaNode, diags *diag.Diagnostics) attr.Value {
	attrTypes := containerAttrTypes(schema)
	objType := types.ObjectType{AttrTypes: attrTypes}

	matches := findChildrenByTag(xmlParent, schema.Name)
	if len(matches) == 0 {
		if len(schema.Children) == 0 {
			return types.ListNull(types.StringType)
		}
		return types.ListNull(objType)
	}

	if len(schema.Children) == 0 {
		// List with no modeled children → list of strings (unlikely but safe)
		elems := make([]attr.Value, 0, len(matches))
		for _, m := range matches {
			elems = append(elems, types.StringValue(m.Text))
		}
		lv, _ := types.ListValue(types.StringType, elems)
		return lv
	}

	elems := make([]attr.Value, 0, len(matches))
	for _, xmlNode := range matches {
		objAttrs := make(map[string]attr.Value, len(schema.Children))
		for _, child := range schema.Children {
			childName := normalizeName(child.Name)
			objAttrs[childName] = xmlNodeToModelValue(ctx, xmlNode, child, diags)
			if diags.HasError() {
				return types.ListNull(objType)
			}
		}
		ov, d := types.ObjectValue(attrTypes, objAttrs)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType)
		}
		elems = append(elems, ov)
	}

	lv, d := types.ListValue(objType, elems)
	diags.Append(d...)
	return lv
}

// containerAttrTypes builds the attr.Type map for a container/list schema node.
func containerAttrTypes(schema patch.SchemaNode) map[string]attr.Type {
	if len(schema.Children) == 0 {
		return nil
	}
	attrTypes := make(map[string]attr.Type, len(schema.Children))
	for _, child := range schema.Children {
		name := normalizeName(child.Name)
		attrTypes[name] = schemaNodeToAttrType(child)
	}
	return attrTypes
}

// schemaNodeToAttrType returns the attr.Type for a schema node.
func schemaNodeToAttrType(node patch.SchemaNode) attr.Type {
	switch node.Type {
	case "leaf":
		return types.StringType
	case "leaf-list":
		return types.ListType{ElemType: types.StringType}
	case "container":
		childTypes := containerAttrTypes(node)
		if len(childTypes) == 0 {
			return types.StringType
		}
		return types.ListType{ElemType: types.ObjectType{AttrTypes: childTypes}}
	case "list":
		childTypes := containerAttrTypes(node)
		if len(childTypes) == 0 {
			return types.ListType{ElemType: types.StringType}
		}
		return types.ListType{ElemType: types.ObjectType{AttrTypes: childTypes}}
	default:
		return types.StringType
	}
}

// AttrTypesForSchema returns the top-level attr.Type map for the full schema.
// Used by the resource to construct typed objects.
func AttrTypesForSchema(idx *SchemaIndex) map[string]attr.Type {
	result := map[string]attr.Type{
		"resource_name": types.StringType,
	}
	for _, node := range idx.TopLevel {
		name := normalizeName(node.Name)
		result[name] = schemaNodeToAttrType(node)
	}
	return result
}
