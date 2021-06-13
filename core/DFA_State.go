package dfalearningtoolkit

// StateLabel type as an 8-bit unsigned integer.
type StateLabel uint8

// Constants which represent the 3 possible state labels.
const (
	REJECTING  = iota // 0
	ACCEPTING         // 1
	UNLABELLED        // 2
)

// State struct which represents a State within a DFA. Please note that while
// the Transitions are made to be accessible from outside of this toolkit use
// AddTransition, UpdateTransition and GetTransitionValue where required. It
// was left public to be accessible to json.Marshal for import/export.
type State struct {
	Label       StateLabel // State Label (Rejecting, Accepting, or Unlabelled).
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

// IsUnlabelled returns true if state label is unlabelled, otherwise returns false.
func (state State) IsUnlabelled() bool {
	return state.Label == UNLABELLED
}

// GetTransitions returns the transitions from the current state.
func (state State) GetTransitions() []int {
	// Return transitions.
	return state.Transitions
}

// AddTransition adds a transition given a stateID. The stateID can be set to -1
// which indicates that the transition for this symbol does not exist (not valid).
func (state *State) AddTransition(stateID int) {
	// Panic if stateID is out of range.
	if stateID > len(state.dfa.States)-1 {
		panic("stateID out of range.")
	}

	// Add new transition using stateID value.
	state.Transitions = append(state.Transitions, stateID)

	// Set computedDepthAndOrder flag to false since DFA was modified.
	state.dfa.computedDepthAndOrder = false
}

// UpdateTransition updates a transition given a symbolID and stateID. This
// is recommended to be used when changing transitions since it updates the
// computedDepthAndOrder flag within DFA. The stateID can be set to -1 to
// 'remove' the transition or any other stateID to add or modify a transition.
func (state *State) UpdateTransition(symbolID, stateID int) {
	// Panic if symbolID is out of range.
	if symbolID > len(state.Transitions)-1 {
		panic("symbolID out of range.")
	}

	// Panic if stateID is out of range.
	if stateID > len(state.dfa.States)-1 {
		panic("stateID out of range.")
	}

	// Update transition.
	state.Transitions[symbolID] = stateID

	// Set computedDepthAndOrder flag to false since DFA was modified.
	state.dfa.computedDepthAndOrder = false
}

// GetTransitionValue returns the transition value given a symbolID.
func (state State) GetTransitionValue(symbolID int) int {
	// Panic if symbolID is out of range.
	if symbolID > len(state.Transitions)-1 {
		panic("symbolID out of range.")
	}

	// Return transition value given symbolID.
	return state.Transitions[symbolID]
}

// InDegree returns the in degree of the state.
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
		if toStateID >= 0 {
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
		if toStateID >= 0 {
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
		if state.GetTransitionValue(symbol) == stateID {
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
