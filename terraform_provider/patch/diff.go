package patch

import (
	"unicode/utf8"
)

// NormalizeLeafMapUTF8 creates a copy of the leaf map with all string values
// sanitized: double-encoded UTF-8 sequences (where UTF-8 bytes were
// misinterpreted as Latin-1 and re-encoded) are repaired back to their
// original form. This prevents false diffs caused by encoding round-trip
// issues between Junos NETCONF responses and Go's xml.Marshal/Unmarshal.
func NormalizeLeafMapUTF8(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = repairDoubleEncodedUTF8(v)
	}
	return result
}

// repairDoubleEncodedUTF8 detects and repairs strings where UTF-8 bytes were
// misinterpreted as Latin-1 (ISO-8859-1) and then re-encoded to UTF-8.
// For example, em-dash U+2014 (UTF-8: E2 80 94) becomes "â\x80\x94"
// when double-encoded. This function reverses that transformation.
func repairDoubleEncodedUTF8(s string) string {
	// Quick check: if the string contains any rune in the C2-F4 range
	// (UTF-8 lead bytes when misread as Latin-1 code points), it might
	// be double-encoded.  Also check for control chars (0x80-0x9F) which
	// appear as raw runes when UTF-8 continuation bytes are misread.
	hasDoubleEncodeSignal := false
	for _, r := range s {
		if (r >= 0x80 && r <= 0x9F) || (r >= 0xC0 && r <= 0xF4) {
			hasDoubleEncodeSignal = true
			break
		}
	}
	if !hasDoubleEncodeSignal {
		return s
	}

	// Try to decode: treat each rune as a byte value (Latin-1 → byte)
	// and see if the resulting byte sequence is valid UTF-8
	bytes := make([]byte, 0, len(s))
	for _, r := range s {
		if r > 0xFF {
			// Rune above Latin-1 range — not double-encoded
			return s
		}
		bytes = append(bytes, byte(r))
	}

	if utf8.Valid(bytes) {
		repaired := string(bytes)
		// Sanity check: repaired string should be shorter (fewer bytes)
		if len(repaired) < len(s) {
			return repaired
		}
	}

	return s
}

// ComputeDiff compares stateMap (what is currently on the device) with
// planMap (what Terraform wants it to be) and returns a map of leaf paths
// to their required CRUD operation.
//
// Rules:
//   - Path in state only               → Delete  (remove it from the device)
//   - Path in both, values differ      → Replace (update the existing value)
//   - Path in plan only                → Create  (add new leaf to the device)
//   - Path in both, values identical   → omitted (no change needed)
func ComputeDiff(stateMap, planMap map[string]string) map[string]Change {
	diff := make(map[string]Change)

	// First pass: iterate state — find deletions and replacements
	for path, stateVal := range stateMap {
		if planVal, exists := planMap[path]; !exists {
			diff[path] = Change{Op: Delete, OldVal: stateVal, NewVal: ""}
		} else if planVal != stateVal {
			diff[path] = Change{Op: Replace, OldVal: stateVal, NewVal: planVal}
		}
		// Values match — no change, do not add to diff
	}

	// Second pass: iterate plan — find creations
	for path, planVal := range planMap {
		if _, exists := stateMap[path]; !exists {
			diff[path] = Change{Op: Create, OldVal: "", NewVal: planVal}
		}
	}

	return diff
}

type DebugChange struct {
	Path   string
	Op     ChangeType
	OldVal string
	NewVal string
}

func DebugSortedChanges(diffMap map[string]Change) []DebugChange {
	ordered := orderedChanges(diffMap)
	result := make([]DebugChange, 0, len(ordered))
	for _, entry := range ordered {
		result = append(result, DebugChange{
			Path:   entry.path,
			Op:     entry.change.Op,
			OldVal: entry.change.OldVal,
			NewVal: entry.change.NewVal,
		})
	}
	return result
}
