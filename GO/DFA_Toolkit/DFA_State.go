package DFA_Toolkit

// State Label type as an 8-bit unsigned integer.
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

// Returns true if state label is accepting, otherwise returns false.
func (state State) IsAccepting() bool{
	return state.Label == ACCEPTING
}

// Returns true if state label is rejecting, otherwise returns false.
func (state State) IsRejecting() bool{
	return state.Label == REJECTING
}

// Returns true if state label is unknown, otherwise returns false.
func (state State) IsUnknown() bool{
	return state.Label == UNKNOWN
}

// Returns the state's depth within DFA.
func (state *State) Depth() int{
	if state.depth == -1{
		state.dfa.CalculateDepthAndOrder()
	}

	return state.depth
}

// Returns the state's order within DFA.
func (state *State) Order() int{
	if state.order == -1{
		state.dfa.CalculateDepthAndOrder()
	}

	return state.order
}

// Returns a pointer to the DFA which contains this State.
func (state State) DFA() *DFA{
	return state.dfa
}