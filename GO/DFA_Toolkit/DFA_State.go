package dfatoolkit

// StateLabel type as an 8-bit unsigned integer.
type StateLabel uint8

// Constants which represent the 3 possible state labels.
const (
	REJECTING = iota // 0
	ACCEPTING        // 1
	UNKNOWN          // 2
)

// State struct which represents a State within a DFA.
type State struct {
	Label       StateLabel // State Label (Rejecting, Accepting, or Unknown).
	Transitions []int      // Transition Table where each element corresponds to a transition for each symbol ID.
	depth       int        // Depth of State within DFA.
	order       int        // Order of State within DFA.
	dfa         *DFA       // Pointer to DFA containing state.
}

// IsAccepting returns true if state label is accepting, otherwise returns false.
func (state State) IsAccepting() bool {
	return state.Label == ACCEPTING
}

// IsRejecting returns true if state label is rejecting, otherwise returns false.
func (state State) IsRejecting() bool {
	return state.Label == REJECTING
}

// IsUnknown returns true if state label is unknown, otherwise returns false.
func (state State) IsUnknown() bool {
	return state.Label == UNKNOWN
}

// Depth returns the state's depth within DFA.
func (state *State) Depth() int {
	if state.depth == -1 {
		state.dfa.CalculateDepthAndOrder()
	}

	return state.depth
}

// Order returns the state's order within DFA.
func (state *State) Order() int {
	if state.order == -1 {
		state.dfa.CalculateDepthAndOrder()
	}

	return state.order
}

// DFA returns a pointer to the DFA which contains this State.
func (state State) DFA() *DFA {
	return state.dfa
}

// TransitionsCount returns the number of transitions to given state ID.
func (state State) TransitionsCount(stateID int) int {
	// Counter to store number of transitions.
	transitionsCount := 0

	// Iterate over each symbol within DFA.
	for alphabetID := range state.DFA().Alphabet {
		// If transition is to given state ID, increment transitions count.
		if state.Transitions[alphabetID] == stateID {
			transitionsCount++
		}
	}

	// Return transitions count.
	return transitionsCount
}
