package generic

import (
	"context"

	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ctxBG is a package-level background context for type-system calls that
// require a context parameter but do not perform any I/O.
var ctxBG = context.Background()

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

// normalizeUnknowns converts any unknown descendants in a Terraform value to
// typed nulls so Create/Read never write unknown state back to Terraform.
func normalizeUnknowns(val attr.Value) attr.Value {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case basetypes.StringValue:
		if v.IsUnknown() {
			return types.StringNull()
		}
		return v
	case basetypes.ListValue:
		elemType := v.ElementType(ctxBG)
		if v.IsUnknown() {
			return types.ListNull(elemType)
		}
		if v.IsNull() {
			return v
		}

		elems := v.Elements()
		normalized := make([]attr.Value, 0, len(elems))
		for _, elem := range elems {
			normalized = append(normalized, normalizeUnknowns(elem))
		}

		lv, d := types.ListValue(elemType, normalized)
		if d.HasError() {
			return types.ListNull(elemType)
		}
		return lv
	case basetypes.ObjectValue:
		attrTypes := v.AttributeTypes(ctxBG)
		if v.IsUnknown() {
			return types.ObjectNull(attrTypes)
		}
		if v.IsNull() {
			return v
		}

		attrs := v.Attributes()
		normalized := make(map[string]attr.Value, len(attrs))
		for key, child := range attrs {
			normalized[key] = normalizeUnknowns(child)
		}

		ov, d := types.ObjectValue(attrTypes, normalized)
		if d.HasError() {
			return types.ObjectNull(attrTypes)
		}
		return ov
	default:
		return val
	}
}

// reconcileListOrder reorders list elements in observed to match the order in
// prior, matching elements by the "name" key field.  This handles the case
// where the device returns YANG list entries in a different order than the .tf
// file specified.  Terraform ListNestedAttribute is ordered, so mismatched
// ordering causes spurious diffs.
//
// The function recurses into container (single-element list) and list values
// to fix ordering at all nesting depths.
// Returns (result, changed) where changed indicates if the value was modified.
func reconcileListOrder(observed, prior attr.Value) (attr.Value, bool) {
	if prior == nil || prior.IsNull() || prior.IsUnknown() {
		return observed, false
	}
	if observed == nil || observed.IsNull() || observed.IsUnknown() {
		return observed, false
	}

	obsLV, obsOk := observed.(basetypes.ListValue)
	priorLV, priorOk := prior.(basetypes.ListValue)
	if !obsOk || !priorOk {
		return observed, false
	}

	obsElems := obsLV.Elements()
	priorElems := priorLV.Elements()

	if len(obsElems) == 0 || len(priorElems) == 0 {
		return observed, false
	}

	// Check if elements are objects (ListNestedAttribute) vs strings (leaf-list)
	if _, isObj := obsElems[0].(basetypes.ObjectValue); !isObj {
		return observed, false
	}

	// Single-element list (container): recurse into the object's attributes
	if len(obsElems) == 1 && len(priorElems) == 1 {
		obsObj := obsElems[0].(basetypes.ObjectValue)
		priorObj := priorElems[0].(basetypes.ObjectValue)
		reconciledObj := reconcileObjectAttrs(obsObj, priorObj)
		if reconciledObj != nil {
			lv, _ := types.ListValue(obsLV.ElementType(ctxBG), []attr.Value{*reconciledObj})
			return lv, true
		}
		return observed, false
	}

	// Multi-element list: reorder observed to match prior by "name" key
	obsMap := buildKeyMap(obsElems)
	if obsMap == nil {
		// No "name" keys found; cannot reconcile ordering, recurse into
		// positional elements instead.
		return reconcilePositional(obsLV, priorElems)
	}

	reordered := make([]attr.Value, 0, len(obsElems))
	used := make(map[string]bool)

	// First pass: emit elements in prior order
	for _, pe := range priorElems {
		po, ok := pe.(basetypes.ObjectValue)
		if !ok {
			continue
		}
		key := extractNameKey(po)
		if key == "" {
			continue
		}
		if ov, exists := obsMap[key]; exists {
			// Recurse into this element's nested lists
			reconciled := reconcileObjectAttrs(ov, po)
			if reconciled != nil {
				reordered = append(reordered, *reconciled)
			} else {
				reordered = append(reordered, ov)
			}
			used[key] = true
		}
	}

	// Second pass: append any observed elements not in prior (new entries)
	for _, oe := range obsElems {
		oo, ok := oe.(basetypes.ObjectValue)
		if !ok {
			continue
		}
		key := extractNameKey(oo)
		if key != "" && !used[key] {
			reordered = append(reordered, oe)
		}
	}

	if len(reordered) == 0 {
		return observed, false
	}

	lv, _ := types.ListValue(obsLV.ElementType(ctxBG), reordered)
	return lv, true
}

// reconcileObjectAttrs reconciles nested list ordering within object attributes.
// Returns a new ObjectValue if any nested lists were reordered, otherwise nil.
func reconcileObjectAttrs(obs, prior basetypes.ObjectValue) *basetypes.ObjectValue {
	obsAttrs := obs.Attributes()
	priorAttrs := prior.Attributes()
	changed := false
	result := make(map[string]attr.Value, len(obsAttrs))

	for k, ov := range obsAttrs {
		pv, exists := priorAttrs[k]
		if !exists {
			result[k] = ov
			continue
		}
		reconciled, didChange := reconcileListOrder(ov, pv)
		result[k] = reconciled
		if didChange {
			changed = true
		}
	}

	if !changed {
		return nil
	}
	ov, _ := types.ObjectValue(obs.AttributeTypes(ctxBG), result)
	return &ov
}

// reconcilePositional handles lists without name keys by recursing into
// positional elements to fix nested ordering.
func reconcilePositional(obsLV basetypes.ListValue, priorElems []attr.Value) (attr.Value, bool) {
	obsElems := obsLV.Elements()
	changed := false
	result := make([]attr.Value, len(obsElems))

	for i, oe := range obsElems {
		if i < len(priorElems) {
			obsObj, ok1 := oe.(basetypes.ObjectValue)
			priorObj, ok2 := priorElems[i].(basetypes.ObjectValue)
			if ok1 && ok2 {
				reconciled := reconcileObjectAttrs(obsObj, priorObj)
				if reconciled != nil {
					result[i] = *reconciled
					changed = true
					continue
				}
			}
		}
		result[i] = oe
	}

	if !changed {
		return obsLV, false
	}
	lv, _ := types.ListValue(obsLV.ElementType(ctxBG), result)
	return lv, true
}

// buildKeyMap indexes list elements by their "name" attribute value.
// Returns nil if no elements have a "name" key.
func buildKeyMap(elems []attr.Value) map[string]basetypes.ObjectValue {
	m := make(map[string]basetypes.ObjectValue, len(elems))
	hasKeys := false
	for _, e := range elems {
		ov, ok := e.(basetypes.ObjectValue)
		if !ok {
			continue
		}
		key := extractNameKey(ov)
		if key != "" {
			m[key] = ov
			hasKeys = true
		}
	}
	if !hasKeys {
		return nil
	}
	return m
}

// extractNameKey returns the "name" string attribute value from an object.
func extractNameKey(obj basetypes.ObjectValue) string {
	attrs := obj.Attributes()
	nameVal, exists := attrs["name"]
	if !exists {
		return ""
	}
	sv, ok := nameVal.(basetypes.StringValue)
	if !ok || sv.IsNull() || sv.IsUnknown() {
		return ""
	}
	return sv.ValueString()
}
