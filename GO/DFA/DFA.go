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
}

type DFA struct {
	states          map[uint]State
	startingState   State
	alphabet        map[int32]bool
	transitionTable map[uint]map[int32]uint
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
	newID := uint(len(dfa.states) + 1)
	dfa.states[newID] = State{stateID: newID, stateStatus: stateStatus}
	dfa.transitionTable[newID] = make(map[int32]uint)
	for k := range dfa.alphabet {
		dfa.transitionTable[newID][k] = 0
	}
}

func (dfa *DFA) AddToAlphabet(letter int32) {
	dfa.alphabet[letter] = true
	for k := range dfa.transitionTable {
		dfa.transitionTable[k][letter] = 0
	}
}

func (dfa DFA) Depth() uint {
	var stateMap = make(map[uint]uint)
	var maxValue uint

	stateMap = dfa.DepthUtil(dfa.startingState.stateID, 0, stateMap)

	for _, v := range stateMap {
		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue
}

func (dfa DFA) DepthUtil(stateID uint, count uint, stateMap map[uint]uint) map[uint]uint {
	stateMap[stateID] = count

	for _, v := range dfa.transitionTable[stateID] {
		for key := range dfa.alphabet{
			if dfa.transitionTable[v][key] != 0{
				if _, ok := stateMap[v]; !ok {
					tempStateMap := dfa.DepthUtil(dfa.transitionTable[v][key], count+1, stateMap)
					if tempStateMap[dfa.transitionTable[v][key]] > stateMap[dfa.transitionTable[v][key]]{
						stateMap = tempStateMap
					}
				}
			}
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
	var count uint
	var currentStateID uint
	dfa := DFA{
		states: make(map[uint]State),
		alphabet:        make(map[int32]bool),
		transitionTable: make(map[uint]map[int32]uint),
	}

	dfa.AddState(UNKNOWN)
	dfa.startingState = dfa.states[1]

	if strings[0].length == 0 {
		if strings[0].stringStatus == ACCEPTING {
			dfa.startingState.stateStatus = ACCEPTING
		} else {
			dfa.startingState.stateStatus = REJECTING
		}
	}

	for _, stringInstance := range strings {
		if !APTA && stringInstance.stringStatus != ACCEPTING {
			continue
		}
		currentStateID = dfa.startingState.stateID
		count = 0
		for _, character := range stringInstance.stringValue {
			count++
			exists = false
			// new alphabet check
			if !dfa.alphabet[character] {
				dfa.AddToAlphabet(character)
			}

			if dfa.transitionTable[currentStateID][character] != 0 {
				currentStateID = dfa.transitionTable[currentStateID][character]
				exists = true
			}

			if !exists {
				// last symbol in string check
				if count == stringInstance.length {
					if stringInstance.stringStatus == ACCEPTING {
						dfa.AddState(ACCEPTING)
					} else {
						dfa.AddState(REJECTING)
					}
				} else {
					dfa.AddState(UNKNOWN)
				}
				dfa.transitionTable[currentStateID][character] = uint(len(dfa.states))
				currentStateID = uint(len(dfa.states))
			} else {
				// last symbol in string check
				if count == stringInstance.length {
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
			}
		}
	}
	return dfa
}
