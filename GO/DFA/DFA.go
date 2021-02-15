package DFA

import (
	"fmt"
)

type StateStatus uint8

const (
	REJECTING = iota // 0
	ACCEPTING        // 1
	UNKNOWN          // 2
)

type State struct {
	stateStatus StateStatus
	stateID     uint
	transitions map[int32]uint
}

func NewState(stateStatus StateStatus, stateID uint, transitions map[int32]uint) *State{
	if stateStatus < 0 || stateStatus > 2{
		panic("State Status must be 0 (REJECTING), 1 (ACCEPTING) or 2 (UNKNOWN)")
	}

	if stateID < 0{
		panic("State ID must be 0 or bigger")
	}

	return &State{
		stateStatus: stateStatus,
		stateID:     stateID,
		transitions: transitions,
	}
}

type DFA struct {
	states        map[uint]State
	startingState State
	alphabet      map[int32]bool
}

func (dfa DFA) GetAcceptingStates() []State {
	var acceptingStates []State

	for _, v := range dfa.states {
		if v.stateStatus == ACCEPTING {
			acceptingStates = append(acceptingStates, v)
		}
	}

	return acceptingStates
}

func (dfa DFA) GetRejectingStates() []State {
	var rejectingStates []State

	for _, v := range dfa.states {
		if v.stateStatus == REJECTING {
			rejectingStates = append(rejectingStates, v)
		}
	}

	return rejectingStates
}

func (dfa *DFA) AddState(stateStatus StateStatus) {
	dfa.states[uint(len(dfa.states))] = *NewState(stateStatus, uint(len(dfa.states)), map[int32]uint{})
}

func (dfa DFA) Depth() uint {
	var stateMap = make(map[uint]uint)
	var maxValue uint

	stateMap = dfa.DepthUtil(dfa.startingState, 0, stateMap)

	for _, v := range stateMap {
		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue
}

func (dfa DFA) DepthUtil(state State, count uint, stateMap map[uint]uint) map[uint]uint {
	stateMap[state.stateID] = count

	for _, v := range state.transitions {
		if _, ok := stateMap[v]; !ok {
			stateMap = dfa.DepthUtil(dfa.states[v], count+1, stateMap)
		}
	}

	return stateMap
}

func (dfa DFA) Describe(detail bool) {
	fmt.Println("This DFA has", len(dfa.states), "states and", len(dfa.alphabet), "alphabet")
	if detail {
		fmt.Println("States:")
		for k, v := range dfa.states {
			switch v.stateStatus {
			case ACCEPTING:
				fmt.Println(k, "ACCEPTING")
				break
			case REJECTING:
				fmt.Println(k, "REJECTING")
				break
			case UNKNOWN:
				fmt.Println(k, "UNKNOWN")
				break
			}
		}
		fmt.Println("Accepting States:")
		for _, v := range dfa.GetAcceptingStates() {
			fmt.Println(v.stateID)
		}
		fmt.Println("Rejecting States:")
		for _, v := range dfa.GetRejectingStates() {
			fmt.Println(v.stateID)
		}
		fmt.Println("Starting State:", dfa.startingState.stateID)
		fmt.Println("Alphabet:")
		for key := range dfa.alphabet {
			fmt.Println(string(key))
		}
	}
}

func GetPTAFromListOfStringInstances(strings []StringInstance, APTA bool) DFA {
	strings = SortListOfStringInstances(strings)
	var exists bool
	var count int
	var alphabet = make(map[int32]bool)
	var states = make(map[uint]State)
	var startingStateID uint = 0
	var currentStateID uint

	if strings[0].length == 0 {
		if strings[0].stringStatus == ACCEPTING {
			states[0] = *NewState(ACCEPTING, 0, map[int32]uint{})
		} else {
			states[0] = *NewState(REJECTING, 0, map[int32]uint{})
		}
	} else {
		states[0] = *NewState(UNKNOWN, 0, map[int32]uint{})
	}

	for _, stringInstance := range strings {
		if !APTA && stringInstance.stringStatus != ACCEPTING {
			continue
		}
		currentStateID = startingStateID
		count = 0
		for _, character := range stringInstance.stringValue {
			count++
			exists = false
			// new alphabet check
			if !alphabet[character] {
				alphabet[character] = true
			}

			if value, ok := states[currentStateID].transitions[character]; ok {
				currentStateID = value
				exists = true
			}

			if !exists {
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						states[uint(len(states))] = *NewState(ACCEPTING, uint(len(states)), map[int32]uint{})
					} else {
						states[uint(len(states))] = *NewState(REJECTING, uint(len(states)), map[int32]uint{})
					}
				} else {
					states[uint(len(states))] = *NewState(UNKNOWN, uint(len(states)), map[int32]uint{})
				}
				states[currentStateID].transitions[character] = states[uint(len(states))-1].stateID
				currentStateID = states[uint(len(states))-1].stateID
			} else {
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						if states[currentStateID].stateStatus == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						} else {
							tempState := states[currentStateID]
							tempState.stateStatus = ACCEPTING
							states[currentStateID] = tempState
						}
					} else {
						if states[currentStateID].stateStatus == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							tempState := states[currentStateID]
							tempState.stateStatus = REJECTING
							states[currentStateID] = tempState
						}
					}
				}
			}
		}
	}
	return DFA{states, states[startingStateID], alphabet}
}
