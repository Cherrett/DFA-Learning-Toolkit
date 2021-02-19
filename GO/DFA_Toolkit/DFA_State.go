package DFA_Toolkit

type StateStatus uint8

const (
	REJECTING = iota // 0
	ACCEPTING        // 1
	UNKNOWN          // 2
)

type State struct {
	stateStatus StateStatus
	transitions []int
}

func (state State) IsAccepting() bool{
	return state.stateStatus == ACCEPTING
}

func (state State) IsRejecting() bool{
	return state.stateStatus == REJECTING
}

func (state State) IsUnknown() bool{
	return state.stateStatus == UNKNOWN
}

func (state *State) UpdateStateStatus(stateStatus StateStatus){
	state.stateStatus = stateStatus
}