package testup

import "testing"

type Register func(caseName string, caseImpl func())

// Suite re-runs the callback for every check registered during the first callback run.
// Only the check test case callback for the current run is actually invoked.
// This allows common setup/teardown in the suite to be re-executed for every test case.
func Suite(t *testing.T, suite func(t *testing.T, check Register)) {
	runTargetAndRecurse(t, []*stackFrame{}, suite)
}

type suiteFunc func(t *testing.T, check Register)

type stackFrame struct {
	names  []string
	target int
}

func runTargetAndRecurse(t *testing.T, stack []*stackFrame, suite suiteFunc) {
	newNames := runStackTarget(t, stack, suite)
	if len(newNames) > 0 {
		runLastFrame(t, append(stack, &stackFrame{names: newNames}), suite)
	}
}

func runStackTarget(t *testing.T, stack []*stackFrame, suite suiteFunc) (subNames []string) {
	seenNewNames := map[string]struct{}{}

	currentCase := make([]int, 0, len(stack)+1)
	currentCase = append(currentCase, 0)

	suite(t, func(name string, cb func()) {
		currentDepth := len(currentCase)

		// If we have a longer index than we have stack, this callback is being executed from
		// within the target test case. Record the name of sub-checks without executing them.
		if currentDepth > len(stack) {
			if _, ok := seenNewNames[name]; ok {
				t.Fatalf("duplicate test case %q", name)
			}
			seenNewNames[name] = struct{}{}
			subNames = append(subNames, name)
			return
		}

		// Find the frame for the current case and check that the case is valid.
		currIdx := currentCase[currentDepth-1]
		currFrame := stack[currentDepth-1]
		if currIdx >= len(currFrame.names) {
			t.Fatalf("unexpected extra case %q", name)
		}
		if recordedName := currFrame.names[currIdx]; name != recordedName {
			// Although not necessary, we're strict about the test case names staying the same to help
			// debug test code.
			t.Fatalf("case name at index %d changed; first %q then %q", currIdx, recordedName, name)
		}

		// Determine if we need to do anything, then record that we've seen the current case by
		// incrementing the case index.
		targetIdx := currFrame.target
		runCase := currIdx == targetIdx
		currentCase[currentDepth-1]++
		if !runCase {
			return
		}

		// We know that the current case is in the path to the target. Add a new frame of indexes
		// for the sub-cases of the current case.
		currentCase = append(currentCase, 0)
		defer func() {
			currentCase = currentCase[:currentDepth]
		}()

		// Execute check callback
		cb()

		// The check callback should have called back to us the same number of times as previously
		// recorded, unless we were recording a new frame. We verify that these callbacks actually
		// happened as the strict enforcement should help debug test code and also ensures that the
		// target case was actually executed.
		if len(stack) >= len(currentCase) {
			called := currentCase[currentDepth]
			expectedCalls := len(stack[currentDepth].names)
			if called < expectedCalls {
				t.Fatalf("missing test case callbacks; expected %d but got %d", expectedCalls, called)
			}
		}
	})

	return subNames
}

func runLastFrame(t *testing.T, stack []*stackFrame, suite suiteFunc) {
	newFrame := stack[len(stack)-1]
	for newFrame.target = 0; newFrame.target < len(newFrame.names); newFrame.target++ {
		caseName := newFrame.names[newFrame.target]
		t.Run(caseName, func(t *testing.T) {
			runTargetAndRecurse(t, stack, suite)
		})
	}
}
