package patch

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
