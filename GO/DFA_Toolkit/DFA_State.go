package DFA_Toolkit

type StateStatus uint8

const (
	REJECTING = iota // 0
	ACCEPTING        // 1
	UNKNOWN          // 2
)

type State struct {
	StateStatus StateStatus
	Transitions []int

	Depth int
	Order int
}

func (state State) IsAccepting() bool{
	return state.StateStatus == ACCEPTING
}

func (state State) IsRejecting() bool{
	return state.StateStatus == REJECTING
}

func (state State) IsUnknown() bool{
	return state.StateStatus == UNKNOWN
}

func (state *State) UpdateStateStatus(stateStatus StateStatus){
	state.StateStatus = stateStatus
}