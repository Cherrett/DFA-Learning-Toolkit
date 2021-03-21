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
	StateStatus StateStatus // State Status (Rejecting, Accepting, or Unknown)
	Transitions []int		// Transition Table where each element corresponds to a transition for each symbol ID.

	Depth int				// Depth of State within DFA
	Order int				// Order of State within DFA
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
