package patch

import (
	"encoding/json"
	"testing"
)

func TestNumRangeBounds(t *testing.T) {
	var rs []NumRange
	if err := json.Unmarshal([]byte(`[{"min": 1, "max": 65535}, {"min": 0, "max": "9223372036.854775807"}, {"min": -5}]`), &rs); err != nil {
		t.Fatal(err)
	}
	if *rs[0].Min != 1 || *rs[0].Max != 65535 {
		t.Fatalf("integer bounds: %v %v", *rs[0].Min, *rs[0].Max)
	}
	if *rs[1].Max != 9223372036.854775807 {
		t.Fatalf("decimal64 string bound: %v", *rs[1].Max)
	}
	if *rs[2].Min != -5 || rs[2].Max != nil {
		t.Fatalf("missing max: %v", rs[2].Max)
	}
	if err := json.Unmarshal([]byte(`[{"min": "max"}]`), &rs); err == nil {
		t.Fatal("expected an error for a non-numeric bound")
	}
}
