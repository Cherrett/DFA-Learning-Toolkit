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

func (dfa *DFA) RemoveState(state State) {
	delete(dfa.states, state.stateID)
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
		fmt.Println("Alphabet:")
		for key := range dfa.alphabet {
			fmt.Println(string(key))
		}
		fmt.Println("Starting State:", dfa.startingState.stateID)
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
		fmt.Println("Transitions:")
		for fromStateID, fromState := range dfa.states {
			for char, toStateID := range fromState.transitions{
				fmt.Println(fromStateID, "--", string(char), "->", toStateID)
			}
		}
	}
}

func (dfa *DFA) UpdateStateStatus(stateID uint, stateStatus StateStatus) {
	tempState := dfa.states[stateID]
	tempState.stateStatus = stateStatus
	dfa.states[stateID] = tempState
}

func GetPTAFromListOfStringInstances(strings []StringInstance, APTA bool) DFA {
	strings = SortListOfStringInstances(strings)
	var count int
	var currentStateID uint
	dfa := DFA{
		states:   make(map[uint]State),
		alphabet: make(map[int32]bool),
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
							dfa.UpdateStateStatus(currentStateID, ACCEPTING)
						}
					} else {
						if dfa.states[currentStateID].stateStatus == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							dfa.UpdateStateStatus(currentStateID, REJECTING)
						}
					}
				}
			} else {
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						dfa.AddState(ACCEPTING)
					} else {
						dfa.AddState(REJECTING)
					}
				} else {
					dfa.AddState(UNKNOWN)
				}
				dfa.states[currentStateID].transitions[character] = dfa.states[uint(len(dfa.states))-1].stateID
				currentStateID = dfa.states[uint(len(dfa.states))-1].stateID
			}
		}
	}
	return dfa
}

func (dfa DFA) AccuracyOfDFA(stringInstances []StringInstance) float32 {
	correctClassifications := float32(0)

	for _, stringInstance := range stringInstances {
		if stringInstance.stringStatus == GetStringStatusInRegardToDFA(stringInstance, dfa) {
			correctClassifications++
		}
	}
	return correctClassifications / float32(len(stringInstances))
}

func (dfa DFA) UnreachableStates() []State {
	reachableStates := map[uint]bool{dfa.startingState.stateID: true}
	currentStates := map[uint]bool{dfa.startingState.stateID: true}

	for len(currentStates) != 0 {
		nextStates := map[uint]bool{}
		for stateID := range currentStates {
			for character := range dfa.alphabet {
				if _, ok := dfa.states[stateID].transitions[character]; ok {
					nextStates[dfa.states[stateID].transitions[character]] = true
				}
			}
		}
		// Don’t visit states we know to be reachable
		currentStates = map[uint]bool{}
		for stateID := range nextStates {
			if !reachableStates[stateID] {
				currentStates[stateID] = true
			}
		}

		// States in Current are definitely reachable.
		for stateID := range currentStates {
			if !reachableStates[stateID] {
				reachableStates[stateID] = true
			}
		}
	}

	var unReachableStates []State
	for _, state := range dfa.states {
		if !reachableStates[state.stateID] {
			unReachableStates = append(unReachableStates, state)
		}
	}

	return unReachableStates
}

func (dfa *DFA) RemoveUnreachableStates() {
	unreachableStates := dfa.UnreachableStates()
	for _, unreachableState := range unreachableStates {
		dfa.RemoveState(unreachableState)
	}
}

type StateIDPair struct {
	state1 uint
	state2 uint
}

func (dfa DFA) Mark() [][]State {
	distinguishablePairs := map[StateIDPair]bool{}

	dfa.RemoveUnreachableStates()
	for stateID, state := range dfa.states {
		for stateID2, state2 := range dfa.states {
			if stateID != stateID2 && state.stateStatus != state2.stateStatus &&
				!distinguishablePairs[StateIDPair{state.stateID, state2.stateID}] &&
				!distinguishablePairs[StateIDPair{state2.stateID, state.stateID}] {
				distinguishablePairs[StateIDPair{state.stateID, state2.stateID}] = true
			}
		}
	}

	oldCount := 0
	for oldCount != len(distinguishablePairs) {
		oldCount = len(distinguishablePairs)

		for stateID, state := range dfa.states {
			for stateID2, state2 := range dfa.states {
				if stateID == stateID2 || distinguishablePairs[StateIDPair{state.stateID, state2.stateID}] ||
					distinguishablePairs[StateIDPair{state2.stateID, state.stateID}] {
					continue
				} else {
					for character := range dfa.alphabet {
						if resultantStateID1, ok := dfa.states[state.stateID].transitions[character]; ok {
							if resultantStateID2, ok := dfa.states[state2.stateID].transitions[character]; ok {
								if distinguishablePairs[StateIDPair{resultantStateID1, resultantStateID2}] ||
									distinguishablePairs[StateIDPair{resultantStateID2, resultantStateID1}] {
									distinguishablePairs[StateIDPair{state.stateID, state2.stateID}] = true
								}
							}
						}
					}
				}
			}
		}
	}

	var distinguishablePairsList [][]State
	for stateIDPair := range distinguishablePairs {
		distinguishablePairsList = append(distinguishablePairsList, []State{dfa.states[stateIDPair.state1], dfa.states[stateIDPair.state2]})
	}
	return distinguishablePairsList
}
