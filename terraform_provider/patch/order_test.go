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
