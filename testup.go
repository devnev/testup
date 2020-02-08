package testup

import "testing"

type Register func(caseName string, caseImpl func())

// Suite re-runs the callback for every check registered during the first callback run.
// Only the check test case callback for the current run is actually invoked.
// This allows common setup/teardown in the suite to be re-executed for every test case.
func Suite(t *testing.T, suite func(t *testing.T, check Register)) {
	t.Helper()

	// Run the suite once without calling into any cases to discover all the cases
	var names []string
	{
		seen := map[string]struct{}{}
		suite(t, func(name string, _ func()) {
			if _, ok := seen[name]; ok {
				t.Fatalf("duplicate test case %q", name)
			}
			seen[name] = struct{}{}
			names = append(names, name)
		})
	}

	// Run the suite for each registered case, only executing the callback for that case
	for targetIdx, targetName := range names {
		// Use the case name to create a subtest
		t.Run(targetName, func(t *testing.T) {
			caseIdx := 0
			suite(t, func(caseName string, tc func()) {
				// Be strict about extra cases and unexpected names to avoid tests silently passing
				if caseIdx >= len(names) {
					t.Fatalf("unexpected extra case %q", caseName)
				}
				if caseName != names[caseIdx] {
					t.Fatalf("case name at index %d changed; first %q then %q", caseIdx, names[caseIdx], caseName)
				}

				// Keep track of how many cases we've seen and run the current case if it's our target
				runCase := caseIdx == targetIdx
				caseIdx++
				if runCase {
					tc()
				}
			})
			// Be strict about the number of callbacks seen in any case. This also ensures we actually called the target case.
			if caseIdx != len(names) {
				t.Fatalf("missing test case callbacks; expected %d but got %d", len(names), caseIdx)
			}
		})
	}
}
