package generic

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"terraform_provider/netconf"
	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Serve runs the generic provider for schema (the embedded, raw or gzipped
// JSON schema) as a Terraform plugin. providerType is the provider's type name
// (junos-srx); its one resource is terraform-provider-<providerType>.
//
// The provider speaks the plugin protocol directly, without
// terraform-plugin-framework: the resource's value is converted to and from
// configuration XML by walking the schema, and the framework's per-request
// walks of every attribute of every object are prohibitive with a full Junos
// model.
func Serve(address, providerType string, schema []byte, debug bool) error {
	var opts []tf6server.ServeOpt
	if debug {
		opts = append(opts, tf6server.WithManagedDebug())
	}
	srv := NewServer(providerType, schema)
	return tf6server.Serve(address, func() tfprotov6.ProviderServer { return srv }, opts...)
}

// Server is the generic provider's plugin protocol server.
type Server struct {
	resourceType string
	schemaRaw    []byte

	loadOnce sync.Once
	loaded   *loadedSchema
	loadErr  error

	mu     sync.Mutex
	client netconf.Client
}

type loadedSchema struct {
	idx    map[string]*patch.NodeInfo
	nodes  []patch.SchemaNode // under <configuration>
	schema *tfprotov6.Schema
	typ    tftypes.Object
}

var _ tfprotov6.ProviderServer = (*Server)(nil)

// NewServer returns a server for schema; the schema is loaded on first use.
func NewServer(providerType string, schema []byte) *Server {
	return &Server{resourceType: "terraform-provider-" + providerType, schemaRaw: schema}
}

func (s *Server) load() (*loadedSchema, error) {
	s.loadOnce.Do(func() {
		idx, nodes, err := LoadSchema(s.schemaRaw)
		if err == nil && len(nodes) == 0 {
			err = fmt.Errorf("schema has no configuration node")
		}
		if err == nil {
			err = ValidateNames(nodes[0].Children, nodes[0].Name)
		}
		if err != nil {
			s.loadErr = err
			return
		}
		configNodes := nodes[0].Children
		schema, typ := BuildSchema(configNodes)
		s.loaded = &loadedSchema{idx: idx, nodes: configNodes, schema: schema, typ: typ}
	})
	return s.loaded, s.loadErr
}

func (s *Server) device() (*device, []*tfprotov6.Diagnostic) {
	l, err := s.load()
	if err != nil {
		return nil, errorDiag("Invalid embedded schema", err)
	}
	s.mu.Lock()
	client := s.client
	s.mu.Unlock()
	if client == nil {
		return nil, errorDiag("Provider not configured", fmt.Errorf("no NETCONF client"))
	}
	return &device{client: client, idx: l.idx, nodes: l.nodes, typ: l.typ}, nil
}

func errorDiag(summary string, err error) []*tfprotov6.Diagnostic {
	return []*tfprotov6.Diagnostic{{
		Severity: tfprotov6.DiagnosticSeverityError,
		Summary:  summary,
		Detail:   err.Error(),
	}}
}

func unsupported(what string) []*tfprotov6.Diagnostic {
	return errorDiag("Unsupported", fmt.Errorf("%s is not supported by this provider", what))
}

// providerType is the provider configuration's type.
var providerType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"host":     tftypes.String,
	"username": tftypes.String,
	"password": tftypes.String,
	"port":     tftypes.Number,
	"sshkey":   tftypes.String,
}}

var providerSchema = &tfprotov6.Schema{Block: &tfprotov6.SchemaBlock{Attributes: []*tfprotov6.SchemaAttribute{
	{Name: "host", Type: tftypes.String, Required: true},
	{Name: "username", Type: tftypes.String, Required: true},
	{Name: "password", Type: tftypes.String, Optional: true, Sensitive: true},
	{Name: "port", Type: tftypes.Number, Required: true},
	{Name: "sshkey", Type: tftypes.String, Optional: true, Sensitive: true},
}}}

func (s *Server) GetMetadata(context.Context, *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	return &tfprotov6.GetMetadataResponse{
		ServerCapabilities: &tfprotov6.ServerCapabilities{},
		Resources:          []tfprotov6.ResourceMetadata{{TypeName: s.resourceType}},
	}, nil
}

func (s *Server) GetProviderSchema(context.Context, *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	resp := &tfprotov6.GetProviderSchemaResponse{
		ServerCapabilities: &tfprotov6.ServerCapabilities{},
		Provider:           providerSchema,
		ResourceSchemas:    map[string]*tfprotov6.Schema{},
	}
	l, err := s.load()
	if err != nil {
		resp.Diagnostics = errorDiag("Invalid embedded schema", err)
		return resp, nil
	}
	resp.ResourceSchemas[s.resourceType] = l.schema
	return resp, nil
}

func (s *Server) GetResourceIdentitySchemas(context.Context, *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	return &tfprotov6.GetResourceIdentitySchemasResponse{}, nil
}

func (s *Server) ValidateProviderConfig(_ context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	return &tfprotov6.ValidateProviderConfigResponse{PreparedConfig: req.Config}, nil
}

func (s *Server) ConfigureProvider(_ context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	resp := &tfprotov6.ConfigureProviderResponse{}
	v, err := req.Config.Unmarshal(providerType)
	if err != nil {
		resp.Diagnostics = errorDiag("Invalid provider configuration", err)
		return resp, nil
	}
	var attrs map[string]tftypes.Value
	if err := v.As(&attrs); err != nil {
		resp.Diagnostics = errorDiag("Invalid provider configuration", err)
		return resp, nil
	}
	str := func(name string) string {
		var out string
		if a := attrs[name]; a.IsKnown() && !a.IsNull() {
			_ = a.As(&out)
		}
		return out
	}
	port := 0
	if p := attrs["port"]; p.IsKnown() && !p.IsNull() {
		var n big.Float
		if err := p.As(&n); err == nil {
			i, _ := n.Int64()
			port = int(i)
		}
	}
	client, err := netconf.NewClient(str("username"), str("password"), str("sshkey"), str("host"), port)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to create NETCONF client", err)
		return resp, nil
	}
	s.mu.Lock()
	s.client = client
	s.mu.Unlock()
	return resp, nil
}

// StopProvider closes the NETCONF session.
func (s *Server) StopProvider(context.Context, *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	s.mu.Lock()
	client := s.client
	s.mu.Unlock()
	resp := &tfprotov6.StopProviderResponse{}
	if client != nil {
		if err := client.Close(); err != nil {
			resp.Error = err.Error()
		}
	}
	return resp, nil
}

func (s *Server) ValidateResourceConfig(_ context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateResourceConfigResponse{}
	if req.TypeName != s.resourceType {
		resp.Diagnostics = unsupported("resource type " + req.TypeName)
	}
	return resp, nil
}

// UpgradeResourceState decodes the stored state with the current schema,
// ignoring attributes the schema no longer has.
func (s *Server) UpgradeResourceState(_ context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	resp := &tfprotov6.UpgradeResourceStateResponse{}
	l, err := s.load()
	if err != nil {
		resp.Diagnostics = errorDiag("Invalid embedded schema", err)
		return resp, nil
	}
	v, err := req.RawState.UnmarshalWithOpts(l.typ, tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to decode state", err)
		return resp, nil
	}
	dv, err := tfprotov6.NewDynamicValue(l.typ, v)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to encode state", err)
		return resp, nil
	}
	resp.UpgradedState = &dv
	return resp, nil
}

func (s *Server) ReadResource(_ context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	resp := &tfprotov6.ReadResourceResponse{NewState: req.CurrentState, Private: req.Private}
	d, diags := s.device()
	if diags != nil {
		resp.Diagnostics = diags
		return resp, nil
	}
	current, err := req.CurrentState.Unmarshal(d.typ)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to decode state", err)
		return resp, nil
	}
	if current.IsNull() {
		return resp, nil
	}
	state, err := d.read(current)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to read configuration", err)
		return resp, nil
	}
	if resp.NewState, err = dynamicValue(d.typ, state); err != nil {
		resp.Diagnostics = errorDiag("Failed to encode state", err)
	}
	return resp, nil
}

// PlanResourceChange plans the proposed state: nothing is computed. A change
// of resource_name replaces the resource.
func (s *Server) PlanResourceChange(_ context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	resp := &tfprotov6.PlanResourceChangeResponse{PlannedState: req.ProposedNewState, PlannedPrivate: req.PriorPrivate}
	l, err := s.load()
	if err != nil {
		resp.Diagnostics = errorDiag("Invalid embedded schema", err)
		return resp, nil
	}
	prior, err := req.PriorState.Unmarshal(l.typ)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to decode prior state", err)
		return resp, nil
	}
	proposed, err := req.ProposedNewState.Unmarshal(l.typ)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to decode proposed state", err)
		return resp, nil
	}
	if !prior.IsNull() && !proposed.IsNull() && resourceName(prior) != resourceName(proposed) {
		resp.RequiresReplace = []*tftypes.AttributePath{tftypes.NewAttributePath().WithAttributeName(resourceNameAttribute)}
	}
	return resp, nil
}

func (s *Server) ApplyResourceChange(_ context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	resp := &tfprotov6.ApplyResourceChangeResponse{Private: req.PlannedPrivate}
	d, diags := s.device()
	if diags != nil {
		resp.Diagnostics = diags
		return resp, nil
	}
	prior, err := req.PriorState.Unmarshal(d.typ)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to decode prior state", err)
		return resp, nil
	}
	planned, err := req.PlannedState.Unmarshal(d.typ)
	if err != nil {
		resp.Diagnostics = errorDiag("Failed to decode planned state", err)
		return resp, nil
	}

	var state tftypes.Value
	switch {
	case planned.IsNull():
		if err := d.delete(prior); err != nil {
			resp.Diagnostics = errorDiag("Failed to delete configuration", err)
			return resp, nil
		}
		state = tftypes.NewValue(d.typ, nil)
	case prior.IsNull():
		if state, err = d.create(planned); err != nil {
			resp.Diagnostics = errorDiag("Failed to create configuration", err)
			return resp, nil
		}
	default:
		if state, err = d.update(planned); err != nil {
			resp.Diagnostics = errorDiag("Failed to update configuration", err)
			return resp, nil
		}
	}
	if resp.NewState, err = dynamicValue(d.typ, state); err != nil {
		resp.Diagnostics = errorDiag("Failed to encode state", err)
	}
	return resp, nil
}

func dynamicValue(typ tftypes.Type, v tftypes.Value) (*tfprotov6.DynamicValue, error) {
	dv, err := tfprotov6.NewDynamicValue(typ, v)
	if err != nil {
		return nil, err
	}
	return &dv, nil
}

func (s *Server) ImportResourceState(context.Context, *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	return &tfprotov6.ImportResourceStateResponse{Diagnostics: unsupported("import")}, nil
}

func (s *Server) MoveResourceState(context.Context, *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	return &tfprotov6.MoveResourceStateResponse{Diagnostics: unsupported("moving resource state")}, nil
}

func (s *Server) UpgradeResourceIdentity(context.Context, *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	return &tfprotov6.UpgradeResourceIdentityResponse{Diagnostics: unsupported("resource identity")}, nil
}

func (s *Server) GenerateResourceConfig(context.Context, *tfprotov6.GenerateResourceConfigRequest) (*tfprotov6.GenerateResourceConfigResponse, error) {
	return &tfprotov6.GenerateResourceConfigResponse{Diagnostics: unsupported("generating resource configuration")}, nil
}

func (s *Server) ValidateDataResourceConfig(_ context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	return &tfprotov6.ValidateDataResourceConfigResponse{Diagnostics: unsupported("data source " + req.TypeName)}, nil
}

func (s *Server) ReadDataSource(_ context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	return &tfprotov6.ReadDataSourceResponse{Diagnostics: unsupported("data source " + req.TypeName)}, nil
}

func (s *Server) CallFunction(_ context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	return &tfprotov6.CallFunctionResponse{Error: &tfprotov6.FunctionError{Text: "function " + req.Name + " is not supported by this provider"}}, nil
}

func (s *Server) GetFunctions(context.Context, *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	return &tfprotov6.GetFunctionsResponse{}, nil
}

func (s *Server) ValidateEphemeralResourceConfig(_ context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	return &tfprotov6.ValidateEphemeralResourceConfigResponse{Diagnostics: unsupported("ephemeral resource " + req.TypeName)}, nil
}

func (s *Server) OpenEphemeralResource(_ context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	return &tfprotov6.OpenEphemeralResourceResponse{Diagnostics: unsupported("ephemeral resource " + req.TypeName)}, nil
}

func (s *Server) RenewEphemeralResource(_ context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	return &tfprotov6.RenewEphemeralResourceResponse{Diagnostics: unsupported("ephemeral resource " + req.TypeName)}, nil
}

func (s *Server) CloseEphemeralResource(_ context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	return &tfprotov6.CloseEphemeralResourceResponse{Diagnostics: unsupported("ephemeral resource " + req.TypeName)}, nil
}

func (s *Server) ValidateListResourceConfig(_ context.Context, req *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	return &tfprotov6.ValidateListResourceConfigResponse{Diagnostics: unsupported("list resource " + req.TypeName)}, nil
}

func (s *Server) ListResource(_ context.Context, req *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	return nil, fmt.Errorf("list resource %s is not supported by this provider", req.TypeName)
}

func (s *Server) ValidateActionConfig(_ context.Context, req *tfprotov6.ValidateActionConfigRequest) (*tfprotov6.ValidateActionConfigResponse, error) {
	return &tfprotov6.ValidateActionConfigResponse{Diagnostics: unsupported("action " + req.ActionType)}, nil
}

func (s *Server) PlanAction(_ context.Context, req *tfprotov6.PlanActionRequest) (*tfprotov6.PlanActionResponse, error) {
	return &tfprotov6.PlanActionResponse{Diagnostics: unsupported("action " + req.ActionType)}, nil
}

func (s *Server) InvokeAction(_ context.Context, req *tfprotov6.InvokeActionRequest) (*tfprotov6.InvokeActionServerStream, error) {
	return nil, fmt.Errorf("action %s is not supported by this provider", req.ActionType)
}
