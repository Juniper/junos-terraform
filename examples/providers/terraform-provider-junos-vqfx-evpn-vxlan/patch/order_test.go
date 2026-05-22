package patch

import (
	"strings"
	"testing"
)

func TestAlignXMLOrderToReference_ReordersKeyedListsAndLeafLists(t *testing.T) {
	idx := mustIdxFromSchema(t, testTrimmedSchema)
	current := []byte(`<configuration>
  <foo>
    <members>c</members>
    <members>a</members>
    <item><address>10.0.0.2</address><value>beta</value></item>
    <item><address>10.0.0.1</address><value>alpha</value></item>
  </foo>
</configuration>`)
	reference := []byte(`<configuration>
  <foo>
    <item><address>10.0.0.1</address><value>alpha</value></item>
    <item><address>10.0.0.2</address><value>beta</value></item>
    <members>a</members>
    <members>c</members>
  </foo>
</configuration>`)

	aligned, err := AlignXMLOrderToReference(current, reference, idx)
	if err != nil {
		t.Fatalf("AlignXMLOrderToReference() error: %v", err)
	}

	got := string(aligned)
	if strings.Index(got, "10.0.0.1") > strings.Index(got, "10.0.0.2") {
		t.Fatalf("expected keyed list entries to follow reference order, got %s", got)
	}
	if strings.Index(got, ">a</members>") > strings.Index(got, ">c</members>") {
		t.Fatalf("expected leaf-list values to follow reference order, got %s", got)
	}
}

func TestAlignXMLOrderToReference_UsesDeterministicFallbackWhenReferenceEmpty(t *testing.T) {
	idx := mustIdxFromSchema(t, testTrimmedSchema)
	current := []byte(`<configuration>
  <foo>
    <item><address>10.0.0.2</address><value>beta</value></item>
    <item><address>10.0.0.1</address><value>alpha</value></item>
  </foo>
</configuration>`)

	aligned, err := AlignXMLOrderToReference(current, nil, idx)
	if err != nil {
		t.Fatalf("AlignXMLOrderToReference() error: %v", err)
	}

	got := string(aligned)
	if strings.Index(got, "10.0.0.1") > strings.Index(got, "10.0.0.2") {
		t.Fatalf("expected deterministic fallback ordering by keyed identity, got %s", got)
	}
}

// ---------------------------------------------------------------------------
// Matrix Order Tests — O1-O6
// ---------------------------------------------------------------------------

// O1 — Top-level list reorder: [B,A,C] → reference [A,B,C]
func TestAlignOrder_TopLevelListReorder(t *testing.T) {
	idx := matrixIdx(t)

	current := []byte(`<configuration><interfaces>
  <interface><name>ge-0/0/2</name><description>c</description></interface>
  <interface><name>ge-0/0/0</name><description>a</description></interface>
  <interface><name>ge-0/0/1</name><description>b</description></interface>
</interfaces></configuration>`)

	reference := []byte(`<configuration><interfaces>
  <interface><name>ge-0/0/0</name><description>a</description></interface>
  <interface><name>ge-0/0/1</name><description>b</description></interface>
  <interface><name>ge-0/0/2</name><description>c</description></interface>
</interfaces></configuration>`)

	aligned, err := AlignXMLOrderToReference(current, reference, idx)
	if err != nil {
		t.Fatal(err)
	}
	got := string(aligned)

	idx0 := strings.Index(got, "ge-0/0/0")
	idx1 := strings.Index(got, "ge-0/0/1")
	idx2 := strings.Index(got, "ge-0/0/2")
	if idx0 > idx1 || idx1 > idx2 {
		t.Fatalf("expected order ge-0/0/0, ge-0/0/1, ge-0/0/2, got:\n%s", got)
	}
}

// O2 — Nested list reorder: each parent entry has independent child ordering
func TestAlignOrder_NestedListPerInstanceReorder(t *testing.T) {
	idx := mustIdxFromSchema(t, testStructuralKeyTrimmedSchema)

	current := []byte(`<configuration><system><syslog>
  <file><name>security</name>
    <contents><name>kernel</name><any/></contents>
    <contents><name>interactive-commands</name><any/></contents>
  </file>
  <file><name>messages</name>
    <contents><name>interactive-commands</name><any/></contents>
    <contents><name>kernel</name><any/></contents>
  </file>
</syslog></system></configuration>`)

	reference := []byte(`<configuration><system><syslog>
  <file><name>security</name>
    <contents><name>interactive-commands</name><any/></contents>
    <contents><name>kernel</name><any/></contents>
  </file>
  <file><name>messages</name>
    <contents><name>kernel</name><any/></contents>
    <contents><name>interactive-commands</name><any/></contents>
  </file>
</syslog></system></configuration>`)

	aligned, err := AlignXMLOrderToReference(current, reference, idx)
	if err != nil {
		t.Fatal(err)
	}
	got := string(aligned)

	// Split output at "messages" to get the two file sections
	msgIdx := strings.Index(got, "<name>messages</name>")
	if msgIdx < 0 {
		t.Fatalf("expected messages file in output:\n%s", got)
	}
	secSection := got[:msgIdx]
	msgSection := got[msgIdx:]

	// In file[security], interactive-commands should come before kernel
	icIdx := strings.Index(secSection, "interactive-commands")
	kerIdx := strings.Index(secSection, "kernel")
	if icIdx < 0 || kerIdx < 0 {
		t.Fatalf("file[security]: missing expected contents entries in:\n%s", secSection)
	}
	if icIdx > kerIdx {
		t.Errorf("file[security]: expected interactive-commands before kernel, got:\n%s", secSection)
	}

	// In file[messages], kernel should come before interactive-commands
	icIdx2 := strings.Index(msgSection, "interactive-commands")
	kerIdx2 := strings.Index(msgSection, "kernel")
	if icIdx2 < 0 || kerIdx2 < 0 {
		t.Fatalf("file[messages]: missing expected contents entries in:\n%s", msgSection)
	}
	if kerIdx2 > icIdx2 {
		t.Errorf("file[messages]: expected kernel before interactive-commands, got:\n%s", msgSection)
	}
}

// O3 — Extra entries in current (not in reference) sort after reference entries
func TestAlignOrder_ExtraEntriesInCurrent(t *testing.T) {
	idx := matrixIdx(t)

	current := []byte(`<configuration><interfaces>
  <interface><name>ge-0/0/3</name><description>extra-d</description></interface>
  <interface><name>ge-0/0/1</name><description>b</description></interface>
  <interface><name>ge-0/0/0</name><description>a</description></interface>
  <interface><name>ge-0/0/2</name><description>extra-c</description></interface>
</interfaces></configuration>`)

	reference := []byte(`<configuration><interfaces>
  <interface><name>ge-0/0/1</name><description>b</description></interface>
  <interface><name>ge-0/0/0</name><description>a</description></interface>
</interfaces></configuration>`)

	aligned, err := AlignXMLOrderToReference(current, reference, idx)
	if err != nil {
		t.Fatal(err)
	}
	got := string(aligned)

	// Reference entries first in ref order: ge-0/0/1, ge-0/0/0
	// Then extras sorted deterministically: ge-0/0/2, ge-0/0/3
	idx1 := strings.Index(got, "ge-0/0/1")
	idx0 := strings.Index(got, "ge-0/0/0")
	idx2 := strings.Index(got, "ge-0/0/2")
	idx3 := strings.Index(got, "ge-0/0/3")

	if idx1 > idx0 {
		t.Errorf("expected ge-0/0/1 before ge-0/0/0 (reference order), got:\n%s", got)
	}
	if idx0 > idx2 || idx0 > idx3 {
		t.Errorf("expected reference entries before extras, got:\n%s", got)
	}
}

// O4 — Missing entry in current (reference has entry current lacks)
func TestAlignOrder_MissingEntryInCurrent(t *testing.T) {
	idx := matrixIdx(t)

	current := []byte(`<configuration><interfaces>
  <interface><name>ge-0/0/2</name><description>c</description></interface>
  <interface><name>ge-0/0/0</name><description>a</description></interface>
</interfaces></configuration>`)

	reference := []byte(`<configuration><interfaces>
  <interface><name>ge-0/0/0</name><description>a</description></interface>
  <interface><name>ge-0/0/1</name><description>b</description></interface>
  <interface><name>ge-0/0/2</name><description>c</description></interface>
</interfaces></configuration>`)

	aligned, err := AlignXMLOrderToReference(current, reference, idx)
	if err != nil {
		t.Fatal(err)
	}
	got := string(aligned)

	// ge-0/0/0 should come before ge-0/0/2 (respecting reference order)
	idx0 := strings.Index(got, "ge-0/0/0")
	idx2 := strings.Index(got, "ge-0/0/2")
	if idx0 > idx2 {
		t.Errorf("expected ge-0/0/0 before ge-0/0/2, got:\n%s", got)
	}
	// Missing ge-0/0/1 should not crash or appear
	if strings.Contains(got, "ge-0/0/1") {
		t.Errorf("missing reference entry ge-0/0/1 should not appear in output")
	}
}

// O5 — Leaf-list reorder produces no diff (set semantics via LeafMapWithSchema)
func TestAlignOrder_LeafListSetSemantics(t *testing.T) {
	idx := matrixIdx(t)

	stateXML := `<configuration><policy-options><community><name>comm</name><members>b</members><members>a</members></community></policy-options></configuration>`
	planXML := `<configuration><policy-options><community><name>comm</name><members>a</members><members>b</members></community></policy-options></configuration>`

	stateMap := LeafMapWithSchema(mustTree(t, stateXML), idx)
	planMap := LeafMapWithSchema(mustTree(t, planXML), idx)
	diff := ComputeDiff(stateMap, planMap)

	if len(diff) != 0 {
		t.Fatalf("leaf-list reorder should produce empty diff (set semantics), got %d: %v", len(diff), diff)
	}
}

// O6 — Empty reference: deterministic fallback by identity
func TestAlignOrder_EmptyReference(t *testing.T) {
	idx := matrixIdx(t)

	current := []byte(`<configuration><interfaces>
  <interface><name>ge-0/0/2</name><description>c</description></interface>
  <interface><name>ge-0/0/0</name><description>a</description></interface>
  <interface><name>ge-0/0/1</name><description>b</description></interface>
</interfaces></configuration>`)

	aligned, err := AlignXMLOrderToReference(current, nil, idx)
	if err != nil {
		t.Fatal(err)
	}
	got := string(aligned)

	// With empty reference, should sort by key identity: ge-0/0/0, ge-0/0/1, ge-0/0/2
	idx0 := strings.Index(got, "ge-0/0/0")
	idx1 := strings.Index(got, "ge-0/0/1")
	idx2 := strings.Index(got, "ge-0/0/2")
	if idx0 > idx1 || idx1 > idx2 {
		t.Fatalf("expected deterministic fallback order ge-0/0/0, ge-0/0/1, ge-0/0/2, got:\n%s", got)
	}
}
