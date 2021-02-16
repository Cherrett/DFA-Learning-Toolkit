package DFA_Toolkit

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
	dfa.states[uint(len(dfa.states))] = State{stateStatus, uint(len(dfa.states)), map[int32]uint{}}
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
	fmt.Println("This DFA_Toolkit has", len(dfa.states), "states and", len(dfa.alphabet), "alphabet")
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
	var count int
	var currentStateID uint
	dfa := DFA{
		states:        make(map[uint]State),
		alphabet:      make(map[int32]bool),
	}

	if strings[0].length == 0 {
		if strings[0].stringStatus == ACCEPTING {
			dfa.AddState(ACCEPTING)
		} else {
			dfa.AddState(REJECTING)
		}
	} else {
		dfa.AddState(UNKNOWN)
	}

	dfa.startingState = dfa.states[0]

	for _, stringInstance := range strings {
		if !APTA && stringInstance.stringStatus != ACCEPTING {
			continue
		}
		currentStateID = dfa.startingState.stateID
		count = 0
		for _, character := range stringInstance.stringValue {
			count++
			// new alphabet check
			if !dfa.alphabet[character] {
				dfa.alphabet[character] = true
			}

			if value, ok := dfa.states[currentStateID].transitions[character]; ok {
				currentStateID = value
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						if dfa.states[currentStateID].stateStatus == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						} else {
							tempState := dfa.states[currentStateID]
							tempState.stateStatus = ACCEPTING
							dfa.states[currentStateID] = tempState
						}
					} else {
						if dfa.states[currentStateID].stateStatus == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							tempState := dfa.states[currentStateID]
							tempState.stateStatus = REJECTING
							dfa.states[currentStateID] = tempState
						}
					}
				}
			}else{
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						dfa.AddState(ACCEPTING)
					} else {
						dfa.AddState(REJECTING)					}
				} else {
					dfa.AddState(UNKNOWN)				}
				dfa.states[currentStateID].transitions[character] = dfa.states[uint(len(dfa.states))-1].stateID
				currentStateID = dfa.states[uint(len(dfa.states))-1].stateID
			}
		}
	}
	return dfa
}
