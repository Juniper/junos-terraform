package generic

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"terraform_provider/netconf"
	"terraform_provider/patch"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// device applies the resource's configuration to a Junos device over NETCONF
// and reads it back, using the schema to convert between the resource's value
// and configuration XML.
type device struct {
	client netconf.Client
	idx    map[string]*patch.NodeInfo
	nodes  []patch.SchemaNode // under <configuration>
	typ    tftypes.Object
}

// create loads the planned configuration (merge) and commits.
func (d *device) create(plan tftypes.Value) (tftypes.Value, error) {
	root, err := ValueToConfig(plan, d.nodes)
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("building configuration: %w", err)
	}
	if err := d.load(root); err != nil {
		return tftypes.Value{}, fmt.Errorf("applying configuration: %w", err)
	}
	if err := d.client.SendCommit(); err != nil {
		return tftypes.Value{}, fmt.Errorf("committing configuration: %w", err)
	}
	return d.read(plan)
}

// update applies the difference between the device's configuration and the
// plan as one patch and commits it. If the device still differs afterwards, it
// loads the whole planned configuration (merge) and commits again.
func (d *device) update(plan tftypes.Value) (tftypes.Value, error) {
	planRoot, err := ValueToConfig(plan, d.nodes)
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("building configuration: %w", err)
	}
	planMap := patch.LeafMapWithSchema(planRoot, d.idx)

	current, err := d.config()
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("reading current configuration: %w", err)
	}
	diff := patch.ComputeDiff(patch.LeafMapWithSchema(current, d.idx), planMap)

	if len(diff) > 0 {
		name := resourceName(plan)
		patchXML, err := patch.CreateDiffPatch(diff, name)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("building NETCONF patch: %w", err)
		}
		debugPatch(name, current, planRoot, diff, string(patchXML))
		if err := d.client.SendUpdate("", string(patchXML), false); err != nil {
			return tftypes.Value{}, fmt.Errorf("sending diff patch: %w", err)
		}
		if err := d.client.SendCommit(); err != nil {
			return tftypes.Value{}, fmt.Errorf("committing configuration: %w", err)
		}

		verified, err := d.config()
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("reading patched configuration: %w", err)
		}
		if len(patch.ComputeDiff(patch.LeafMapWithSchema(verified, d.idx), planMap)) == 0 {
			// The configuration just read is the new state.
			return d.state(verified, plan)
		}
		if err := d.load(planRoot); err != nil {
			return tftypes.Value{}, fmt.Errorf("patch had no effect and fallback update failed: %w", err)
		}
		if err := d.client.SendCommit(); err != nil {
			return tftypes.Value{}, fmt.Errorf("committing fallback update: %w", err)
		}
		return d.read(plan)
	}

	return d.state(current, plan)
}

// delete removes everything in the state from the device and commits.
func (d *device) delete(state tftypes.Value) error {
	root, err := ValueToConfig(state, d.nodes)
	if err != nil {
		return fmt.Errorf("building configuration: %w", err)
	}
	diff := patch.ComputeDiff(patch.LeafMapWithSchema(root, d.idx),
		patch.LeafMapWithSchema(&patch.Node{Tag: "configuration"}, d.idx))
	if len(diff) == 0 {
		return nil
	}
	patchXML, err := patch.CreateDiffPatch(diff, resourceName(state))
	if err != nil {
		return fmt.Errorf("building delete patch: %w", err)
	}
	if err := d.client.SendUpdate("", string(patchXML), false); err != nil {
		return fmt.Errorf("deleting configuration: %w", err)
	}
	if err := d.client.SendCommit(); err != nil {
		return fmt.Errorf("committing configuration: %w", err)
	}
	return nil
}

// read returns the device's configuration as the resource's state, with list
// entries in the order they have in reference (the plan or prior state) where
// the device's order is not significant, and reference's resource_name.
func (d *device) read(reference tftypes.Value) (tftypes.Value, error) {
	current, err := d.config()
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("reading configuration: %w", err)
	}
	return d.state(current, reference)
}

// state returns a configuration read from the device (config) as the
// resource's state, ordered like reference and with its resource_name.
func (d *device) state(current *patch.Node, reference tftypes.Value) (tftypes.Value, error) {
	refRoot, err := ValueToConfig(reference, d.nodes)
	if err != nil {
		return tftypes.Value{}, err
	}
	currentXML, err := patch.MarshalTree(current)
	if err != nil {
		return tftypes.Value{}, err
	}
	refXML, err := patch.MarshalTree(refRoot)
	if err != nil {
		return tftypes.Value{}, err
	}
	aligned, err := patch.AlignXMLOrderToReference(currentXML, refXML, d.idx)
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("aligning configuration order: %w", err)
	}
	tree, err := patch.BuildTree(aligned)
	if err != nil {
		return tftypes.Value{}, err
	}
	observed, err := ConfigToValue(tree, d.nodes, d.typ)
	if err != nil {
		return tftypes.Value{}, err
	}
	return withResourceName(observed, reference, d.typ)
}

// rawConfiguration carries a <configuration> element's content as raw XML.
type rawConfiguration struct {
	XMLName xml.Name `xml:"configuration"`
	Inner   []byte   `xml:",innerxml"`
}

// config reads the device's configuration, keeping only what the schema
// models, in schema order: the form the plan is compared with.
func (d *device) config() (*patch.Node, error) {
	var raw rawConfiguration
	if err := d.client.MarshalConfig(&raw); err != nil {
		return nil, err
	}
	var doc bytes.Buffer
	doc.WriteString("<configuration>")
	doc.Write(raw.Inner)
	doc.WriteString("</configuration>")
	tree, err := patch.BuildTree(doc.Bytes())
	if err != nil {
		return nil, err
	}
	v, err := ConfigToValue(tree, d.nodes, d.typ)
	if err != nil {
		return nil, err
	}
	return ValueToConfig(v, d.nodes)
}

// load merges a configuration into the device's candidate configuration.
func (d *device) load(root *patch.Node) error {
	var inner bytes.Buffer
	for _, c := range root.Children {
		b, err := patch.MarshalTree(c)
		if err != nil {
			return err
		}
		inner.Write(b)
	}
	return d.client.SendDirectTransaction(rawConfiguration{Inner: inner.Bytes()}, false)
}

// withResourceName sets v's resource_name to reference's.
func withResourceName(v, reference tftypes.Value, typ tftypes.Object) (tftypes.Value, error) {
	var attrs map[string]tftypes.Value
	if err := v.As(&attrs); err != nil {
		return tftypes.Value{}, err
	}
	attrs[resourceNameAttribute] = tftypes.NewValue(tftypes.String, nil)
	if !reference.IsNull() {
		var refAttrs map[string]tftypes.Value
		if err := reference.As(&refAttrs); err != nil {
			return tftypes.Value{}, err
		}
		if name, ok := refAttrs[resourceNameAttribute]; ok {
			attrs[resourceNameAttribute] = name
		}
	}
	return tftypes.NewValue(typ, attrs), nil
}

func resourceName(v tftypes.Value) string {
	var attrs map[string]tftypes.Value
	if v.IsNull() || v.As(&attrs) != nil {
		return ""
	}
	var name string
	if n, ok := attrs[resourceNameAttribute]; ok && n.IsKnown() && !n.IsNull() {
		_ = n.As(&name)
	}
	return name
}

func debugPatch(name string, current, plan *patch.Node, diff map[string]patch.Change, payload string) {
	if os.Getenv("JUNOS_TF_DEBUG_PATCH") == "" {
		return
	}
	currentXML, _ := patch.MarshalTree(current)
	planXML, _ := patch.MarshalTree(plan)
	fmt.Fprintf(os.Stderr, "\n=== terraform diff patch debug: %s ===\n", name)
	fmt.Fprintf(os.Stderr, "--- state xml ---\n%s\n", currentXML)
	fmt.Fprintf(os.Stderr, "--- plan xml ---\n%s\n", planXML)
	fmt.Fprintf(os.Stderr, "--- diff map ---\n")
	for _, entry := range patch.DebugSortedChanges(diff) {
		fmt.Fprintf(os.Stderr, "%v | %s | old=%q | new=%q\n", entry.Op, entry.Path, entry.OldVal, entry.NewVal)
	}
	fmt.Fprintf(os.Stderr, "--- patch payload ---\n%s\n", payload)
}
