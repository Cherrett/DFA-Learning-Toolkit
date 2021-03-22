package DFA_Toolkit

// State Status type as an 8-bit unsigned integer.
type StateStatus uint8

// Constants which represent the 3 possible state statuses.
const (
	REJECTING = iota // 0
	ACCEPTING        // 1
	UNKNOWN          // 2
)

// State struct which represents a State within a DFA.
type State struct {
	StateStatus StateStatus // State Status (Rejecting, Accepting, or Unknown).
	Transitions []int       // Transition Table where each element corresponds to a transition for each symbol ID.
	depth       int         // Depth of State within DFA.
	order       int         // Order of State within DFA.
	dfa         *DFA        // Pointer to DFA containing state.
}

// Returns true if state status is accepting, otherwise returns false.
func (state State) IsAccepting() bool{
	return state.StateStatus == ACCEPTING
}

// Returns true if state status is rejecting, otherwise returns false.
func (state State) IsRejecting() bool{
	return state.StateStatus == REJECTING
}

// Returns true if state status is unknown, otherwise returns false.
func (state State) IsUnknown() bool{
	return state.StateStatus == UNKNOWN
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