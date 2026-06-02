package generic

import (
	"context"
	"fmt"
	"os"

	"terraform_provider/netconf"
	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProviderConfig holds the NETCONF client and host info passed from the provider.
type ProviderConfig struct {
	netconf.Client
	Host string
}

// ConfigResource implements a Terraform resource using schema-driven generic logic.
type ConfigResource struct {
	client    ProviderConfig
	idx       *SchemaIndex
	tfSchema  schema.Schema
	rawSchema string
}

// NewConfigResource creates a new generic config resource backed by the given schema JSON.
func NewConfigResource(schemaJSON string) func() resource.Resource {
	return func() resource.Resource {
		idx, err := LoadSchema([]byte(schemaJSON))
		if err != nil {
			panic(fmt.Sprintf("generic: failed to load schema: %v", err))
		}
		s := BuildTerraformSchema(idx)
		return &ConfigResource{
			idx:       idx,
			tfSchema:  s,
			rawSchema: schemaJSON,
		}
	}
}

func (r *ConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(ProviderConfig)
}

func (r *ConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "terraform-provider-" + req.ProviderTypeName
}

func (r *ConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.tfSchema
}

func (r *ConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	attrs := r.getTopLevelAttrs(ctx, &req.Plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	xmlBytes := ModelToXMLBytes(ctx, attrs, r.idx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.SendDirectTransactionRaw(string(xmlBytes), false); err != nil {
		resp.Diagnostics.AddError("Failed while applying configuration", err.Error())
		return
	}
	if err := r.client.SendCommit(); err != nil {
		resp.Diagnostics.AddError("Failed while committing configuration", err.Error())
		return
	}

	state := r.readAndBuildState(ctx, attrs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.setTopLevelAttrs(ctx, &resp.State, state, &resp.Diagnostics)
}

func (r *ConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	priorAttrs := r.getTopLevelAttrs(ctx, &req.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	state := r.readAndBuildState(ctx, priorAttrs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	r.setTopLevelAttrs(ctx, &resp.State, state, &resp.Diagnostics)
}

func (r *ConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planAttrs := r.getTopLevelAttrs(ctx, &req.Plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	planXML := ModelToXMLBytes(ctx, planAttrs, r.idx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	stateXML, err := r.client.GetConfigXML()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read current configuration", err.Error())
		return
	}

	patchIdx, err := patch.UnmarshalTrimmedSchemaIndex(r.rawSchema)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse schema", err.Error())
		return
	}

	planTree, err := patch.BuildTree(planXML)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse plan XML", err.Error())
		return
	}
	stateTree, err := patch.BuildTree(stateXML)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse state XML", err.Error())
		return
	}

	planMap := patch.LeafMapWithSchema(planTree, patchIdx)
	stateMap := patch.LeafMapWithSchema(stateTree, patchIdx)
	diffMap := patch.ComputeDiff(stateMap, planMap)

	if len(diffMap) == 0 {
		r.setTopLevelAttrs(ctx, &resp.State, planAttrs, &resp.Diagnostics)
		return
	}

	resourceName := r.getResourceName(planAttrs)
	patchXMLBytes, err := patch.CreateDiffPatch(diffMap, resourceName)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build NETCONF patch", err.Error())
		return
	}

	debugPatchUpdate(resourceName, planXML, stateXML, diffMap, string(patchXMLBytes))

	if err := r.client.SendUpdate("", string(patchXMLBytes), false); err != nil {
		resp.Diagnostics.AddError("Failed while sending diff patch", err.Error())
		return
	}
	if err := r.client.SendCommit(); err != nil {
		resp.Diagnostics.AddError("Failed while committing configuration", err.Error())
		return
	}

	// Verify patch took effect
	verifiedXML, err := r.client.GetConfigXML()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read patched configuration", err.Error())
		return
	}
	verifiedTree, err := patch.BuildTree(verifiedXML)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse verified XML", err.Error())
		return
	}
	verifiedMap := patch.LeafMapWithSchema(verifiedTree, patchIdx)
	remainingDiff := patch.ComputeDiff(verifiedMap, planMap)
	if len(remainingDiff) > 0 {
		if err := r.client.SendDirectTransactionRaw(string(planXML), false); err != nil {
			resp.Diagnostics.AddError("Patch had no effect and fallback update failed", err.Error())
			return
		}
		if err := r.client.SendCommit(); err != nil {
			resp.Diagnostics.AddError("Fallback update commit failed", err.Error())
			return
		}
	}

	// Use plan values as state: the verification+fallback above ensure the
	// device has the correct config.  Reading back from device can return
	// extra list elements (e.g. from merge) that were not in the plan, which
	// Terraform rejects as "inconsistent result after apply".
	r.setTopLevelAttrs(ctx, &resp.State, planAttrs, &resp.Diagnostics)
}

func (r *ConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateAttrs := r.getTopLevelAttrs(ctx, &req.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	stateXML := ModelToXMLBytes(ctx, stateAttrs, r.idx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	emptyXML := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<configuration></configuration>")

	patchIdx, err := patch.UnmarshalTrimmedSchemaIndex(r.rawSchema)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse schema", err.Error())
		return
	}

	stateTree, err := patch.BuildTree(stateXML)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse state XML", err.Error())
		return
	}
	emptyTree, err := patch.BuildTree(emptyXML)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse empty XML", err.Error())
		return
	}

	stateMap := patch.LeafMapWithSchema(stateTree, patchIdx)
	emptyMap := patch.LeafMapWithSchema(emptyTree, patchIdx)
	diffMap := patch.ComputeDiff(stateMap, emptyMap)
	if len(diffMap) == 0 {
		return
	}

	resourceName := r.getResourceName(stateAttrs)
	patchXMLBytes, err := patch.CreateDiffPatch(diffMap, resourceName)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build delete patch", err.Error())
		return
	}

	if err := r.client.SendUpdate("", string(patchXMLBytes), false); err != nil {
		resp.Diagnostics.AddError("Failed while deleting configuration", err.Error())
		return
	}
	if err := r.client.SendCommit(); err != nil {
		resp.Diagnostics.AddError("Failed while committing configuration", err.Error())
		return
	}
}

// --- Helper methods ---

// attrAccessor abstracts Terraform Plan and State for attribute access.
type attrAccessor interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

// attrWriter abstracts Terraform State for attribute writes.
type attrWriter interface {
	SetAttribute(ctx context.Context, p path.Path, val interface{}) diag.Diagnostics
}

func (r *ConfigResource) getTopLevelAttrs(ctx context.Context, src attrAccessor, diags *diag.Diagnostics) map[string]attr.Value {
	attrs := make(map[string]attr.Value)

	var resourceName types.String
	diags.Append(src.GetAttribute(ctx, path.Root("resource_name"), &resourceName)...)
	if diags.HasError() {
		return nil
	}
	attrs["resource_name"] = resourceName

	for _, node := range r.idx.TopLevel {
		name := normalizeName(node.Name)
		val := r.getAttr(ctx, src, name, node, diags)
		if diags.HasError() {
			return nil
		}
		attrs[name] = val
	}
	return attrs
}

func (r *ConfigResource) getAttr(ctx context.Context, src attrAccessor, name string, node patch.SchemaNode, diags *diag.Diagnostics) attr.Value {
	switch node.Type {
	case "leaf":
		var val types.String
		diags.Append(src.GetAttribute(ctx, path.Root(name), &val)...)
		return val
	case "leaf-list":
		var val types.List
		diags.Append(src.GetAttribute(ctx, path.Root(name), &val)...)
		return val
	case "container":
		var val types.List
		diags.Append(src.GetAttribute(ctx, path.Root(name), &val)...)
		return val
	case "list":
		var val types.List
		diags.Append(src.GetAttribute(ctx, path.Root(name), &val)...)
		return val
	default:
		var val types.String
		diags.Append(src.GetAttribute(ctx, path.Root(name), &val)...)
		return val
	}
}

func (r *ConfigResource) setTopLevelAttrs(ctx context.Context, dst attrWriter, attrs map[string]attr.Value, diags *diag.Diagnostics) {
	for name, val := range attrs {
		diags.Append(dst.SetAttribute(ctx, path.Root(name), val)...)
		if diags.HasError() {
			return
		}
	}
}

func (r *ConfigResource) readAndBuildState(ctx context.Context, referenceAttrs map[string]attr.Value, diags *diag.Diagnostics) map[string]attr.Value {
	xmlBytes, err := r.client.GetConfigXML()
	if err != nil {
		diags.AddError("Failed to read configuration", err.Error())
		return nil
	}

	observed := XMLBytesToModel(ctx, xmlBytes, r.idx, diags)
	if diags.HasError() {
		return nil
	}

	// Preserve resource_name from reference (not in device XML)
	if rn, ok := referenceAttrs["resource_name"]; ok {
		observed["resource_name"] = rn
	}

	// For null observed attrs where reference has a known value, preserve reference
	// (handles empty containers whose XML is not emitted by device).
	// Skip unknown references — state must not contain unknowns.
	for key, refVal := range referenceAttrs {
		if key == "resource_name" {
			continue
		}
		obsVal, exists := observed[key]
		if !exists || obsVal == nil || obsVal.IsNull() {
			if refVal != nil && !refVal.IsNull() && !refVal.IsUnknown() {
				observed[key] = refVal
			}
		}
	}

	// Reconcile list ordering: device may return list elements in different
	// order than the prior state. Terraform ListNestedAttribute is ordered,
	// so we must reorder observed elements to match prior state order by key.
	for key, refVal := range referenceAttrs {
		if key == "resource_name" {
			continue
		}
		obsVal, exists := observed[key]
		if !exists || obsVal == nil || obsVal.IsNull() {
			continue
		}
		reconciled, _ := reconcileListOrder(obsVal, refVal)
		observed[key] = reconciled
	}

	for key, val := range observed {
		observed[key] = normalizeUnknowns(val)
	}

	return observed
}

func (r *ConfigResource) getResourceName(attrs map[string]attr.Value) string {
	if rn, ok := attrs["resource_name"]; ok {
		if sv, ok := rn.(types.String); ok {
			return sv.ValueString()
		}
	}
	return ""
}

func debugPatchUpdate(resourceName string, planXML []byte, stateXML []byte, diffMap map[string]patch.Change, patchPayload string) {
	if os.Getenv("JUNOS_TF_DEBUG_PATCH") == "" {
		return
	}
	fmt.Printf("\n=== terraform diff patch debug: %s ===\n", resourceName)
	fmt.Printf("--- state xml ---\n%s\n", string(stateXML))
	fmt.Printf("--- plan xml ---\n%s\n", string(planXML))
	fmt.Printf("--- diff map ---\n")
	for _, entry := range patch.DebugSortedChanges(diffMap) {
		fmt.Printf("%v | %s | old=%q | new=%q\n", entry.Op, entry.Path, entry.OldVal, entry.NewVal)
	}
	fmt.Printf("--- patch payload ---\n%s\n", patchPayload)
}
