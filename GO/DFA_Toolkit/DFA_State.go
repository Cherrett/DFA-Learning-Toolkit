package dfatoolkit

// StateLabel type as an 8-bit unsigned integer.
type StateLabel uint8

// Constants which represent the 3 possible state labels.
const (
	REJECTING  = iota // 0
	ACCEPTING         // 1
	UNLABELLED        // 2
)

// State struct which represents a State within a DFA.
type State struct {
	Label       StateLabel // State Label (Rejecting, Accepting, or Unknown).
	Transitions []int      // Transition Table where each element corresponds to a transition for each symbol.
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
	return state.Label == UNLABELLED
}

// InDegree returns the in degree of the state..
func (state State) InDegree(stateID int) int {
	// Initialize in degree counter.
	count := 0

	// Iterate over each state within reference DFA.
	for _, state2 := range state.DFA().States {
		// Iterate over each transition from state.
		for _, toStateID := range state2.Transitions {
			// If a transition with a value equal to the stateID
			// is found, the in degree counter is incremented.
			if toStateID == stateID {
				count++
			}
		}
	}

	// Return in degree count.
	return count
}

// OutDegree returns the out degree of the state..
func (state State) OutDegree() int {
	// Initialize out degree counter.
	count := 0

	// Iterate over each transition from state.
	for _, toStateID := range state.Transitions {
		// If a transition with a value not equal to -1 is
		// found (valid), the counter is incremented.
		if toStateID != -1 {
			count++
		}
	}

	// Return out degree count.
	return count
}

// IsLeaf checks whether given state is a leaf within DFA.
// In other words, whether the state has any valid transitions.
func (state State) IsLeaf() bool {
	return state.OutDegree() == 0
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

// AllTransitionsExist checks whether all transitions from a given state
// exist (not -1). In other words, whether the state has a transition for
// each of the symbols within the alphabet.
func (state State) AllTransitionsExist() bool {
	// Iterate over each transition from state.
	for _, toStateID := range state.Transitions {
		// If a transition with -1 is found, false is
		// returned since this means that a transition
		// to that respective symbol does not exist.
		if toStateID == -1 {
			return false
		}
	}

	// Return true if reached since
	// all transitions exist.
	return true
}

// TransitionExists checks whether a transition exists from a given state
// to another state, regardless of the symbol.
func (state State) TransitionExists(stateID int) bool {
	// Iterate over each transition from state.
	for _, toStateID := range state.Transitions {
		// If a transition with stateID is found, true is
		// returned since this means that a transition
		// to that respective stateID exists.
		if toStateID == stateID {
			return true
		}
	}

	// Return false if reached since transition
	// to state does not exist.
	return false
}

// ValidTransitions returns all transitions from a given state
// that are valid (not -1). The symbolIDs of the corresponding
// valid transitions are returned in a slice of integers.
func (state State) ValidTransitions() []int {
	// Slice of symbolIDs.
	var validTransitions []int

	// Iterate over each transition from state.
	for symbolID, toStateID := range state.Transitions {
		// If a transition with a value not equal to
		// -1 is found (valid), the symbol is added
		// to the valid transitions slice.
		if toStateID != -1 {
			validTransitions = append(validTransitions, symbolID)
		}
	}

	// Return populated slice of symbolIDs.
	return validTransitions
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
	for symbol := range state.DFA().Alphabet {
		// If transition is to given state ID, increment transitions count.
		if state.Transitions[symbol] == stateID {
			transitionsCount++
		}
	}

	// Return transitions count.
	return transitionsCount
}

// TotalTransitionsCount returns the number of transitions from given state.
func (state State) TotalTransitionsCount() int {
	// Counter to store number of transitions.
	transitionsCount := 0

	// Iterate over each symbol within DFA.
	for symbol := range state.DFA().Alphabet {
		// If transition is valid (not -1), increment transitions count.
		if state.Transitions[symbol] != -1 {
			transitionsCount++
		}
	}

	// Return transitions count.
	return transitionsCount
}

// Clone returns a clone of the State.
func (state State) Clone() State {
	// Initialize cloned State.
	clonedState := State{
		Label:       state.Label,
		Transitions: make([]int, len(state.Transitions)),
		depth:       state.depth,
		order:       state.order,
		dfa:         state.dfa,
	}

	// Clone the transitions.
	copy(clonedState.Transitions, state.Transitions)

	// Return cloned State.
	return clonedState
}
