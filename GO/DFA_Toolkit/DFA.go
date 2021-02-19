package DFA_Toolkit

import (
	"fmt"
)

type DFA struct {
	states        []State
	startingState int
	symbolMap     map[int32]int
}

func NewDFA() DFA {
	return DFA{states: make([]State, 0), startingState: -1, symbolMap: make(map[int32]int)}
}

func (dfa *DFA) AddState(stateStatus StateStatus) int {
	transitions := make([]int, len(dfa.symbolMap))
	for i := range transitions {
		transitions[i] = -1
	}
	dfa.states = append(dfa.states, State{stateStatus, transitions})
	return len(dfa.states) - 1
}

func (dfa *DFA) RemoveState(stateID int) {
	// panic if stateID to be removed is the starting state
	if dfa.startingState == stateID {
		panic("Cannot remove starting state")
	}
	// remove state from slice of states
	dfa.states = append(dfa.states[:stateID], dfa.states[stateID+1:]...)
	// update transitions to account for new stateIDs and for removed state
	for stateIndex := range dfa.states {
		for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
			if dfa.states[stateIndex].transitions[symbol] == stateID {
				dfa.states[stateIndex].transitions[symbol] = -1
			} else if dfa.states[stateIndex].transitions[symbol] > stateID {
				dfa.states[stateIndex].transitions[symbol] -= 1
			}
		}
	}
}

func (dfa DFA) GetSymbolID(symbol int32) int{
	return dfa.symbolMap[symbol]
}

func (dfa *DFA) AddSymbol(symbol int32){
	dfa.symbolMap[symbol] = len(dfa.symbolMap)
	for stateIndex := range dfa.states {
		dfa.states[stateIndex].transitions = append(dfa.states[stateIndex].transitions, -1)
	}
}

func (dfa *DFA) RemoveSymbol(symbol int32){
	// TODO: CONTINUE THIS FUNCTION
	//symbolID := dfa.symbolMap[symbol]
	// remove symbol from symbolMap
	//delete(dfa.symbolMap, symbol)
	// update transitions to account for removed symbol
	//for stateIndex := range dfa.states {
	//	for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
	//		if dfa.states[stateIndex].transitions[symbol] == stateID {
	//			dfa.states[stateIndex].transitions[symbol] = -1
	//		} else if dfa.states[stateIndex].transitions[symbol] > stateID {
	//			dfa.states[stateIndex].transitions[symbol] -= 1
	//		}
	//	}
	//}
}

func (dfa *DFA) AddTransition(symbol int, fromStateID int, toStateID int) {
	// error checking
	if fromStateID > len(dfa.states)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if toStateID > len(dfa.states)-1 || toStateID < 0 {
		panic("toStateID is out of range")
	} else if symbol > len(dfa.symbolMap)-1 || symbol < 0 {
		panic("symbol is out of range")
	}
	// add transition to fromState's transitions
	dfa.states[fromStateID].transitions[symbol] = toStateID
}

func (dfa *DFA) RemoveTransition(symbol int, fromStateID int) {
	// error checking
	if fromStateID > len(dfa.states)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if symbol > len(dfa.symbolMap)-1 || symbol < 0 {
		panic("symbol is out of range")
	}
	// remove transition to fromState's transitions
	dfa.states[fromStateID].transitions[symbol] = -1
}

func (dfa DFA) AllStates() []int {
	var allStates []int

	for stateID := range dfa.states {
		allStates = append(allStates, stateID)
	}
	return allStates
}

func (dfa DFA) AcceptingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == ACCEPTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) RejectingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == REJECTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) UnknownStates() []int {
	var acceptingStates []int

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == UNKNOWN {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) AllStatesCount() int {
	return len(dfa.states)
}

func (dfa DFA) AcceptingStatesCount() int {
	count := 0

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == ACCEPTING {
			count++
		}
	}
	return count
}

func (dfa DFA) RejectingStatesCount() int {
	count := 0

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == REJECTING {
			count++
		}
	}
	return count
}

func (dfa DFA) UnknownStatesCount() int {
	count := 0

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == UNKNOWN {
			count++
		}
	}
	return count
}

func (dfa DFA) TransitionsCount() int{
	count := 0

	for stateIndex := range dfa.states {
		for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
			if dfa.states[stateIndex].transitions[symbol] != -1 {
				count++
			}
		}
	}
	return count
}

func (dfa DFA) TransitionsCountForSymbol(symbol int) int{
	count := 0

	for stateIndex := range dfa.states {
		if dfa.states[stateIndex].transitions[symbol] != -1 {
			count++
		}
	}
	return count
}

func (dfa DFA) LeavesCount() int{
	count := 0

	for stateIndex := range dfa.states {
		transitionsCount := 0
		for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
			if dfa.states[stateIndex].transitions[symbol] != -1 {
				transitionsCount++
			}
		}
		if transitionsCount == 0{
			count++
		}
	}
	return count
}

func (dfa DFA) LoopsCount() int{
	count := 0

	for stateIndex := range dfa.states {
		transitionsCount := 0
		for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
			if dfa.states[stateIndex].transitions[symbol] != -1 {
				transitionsCount++
			}
		}
		if transitionsCount == 0{
			count++
		}
	}
	return count
}

func (dfa DFA) Depth() uint {
	var stateMap = make(map[int]uint)
	var maxValue uint

	stateMap = dfa.DepthUtil(dfa.startingState, 0, stateMap)

	for _, v := range stateMap {
		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue
}

func (dfa DFA) DepthUtil(stateID int, count uint, stateMap map[int]uint) map[int]uint {
	stateMap[stateID] = count

	for symbol := range dfa.states[stateID].transitions {
		if dfa.states[stateID].transitions[symbol] != -1 {
			stateMap = dfa.DepthUtil(dfa.states[stateID].transitions[symbol], count+1, stateMap)
		}
	}

	return stateMap
}

func (dfa DFA) Describe(detail bool) {
	fmt.Println("This DFA has", len(dfa.states), "states and", len(dfa.symbolMap), "alphabet")
	if detail {
		fmt.Println("Alphabet:")
		for symbol, index := range dfa.symbolMap{
			fmt.Println(index,"-",string(symbol))
		}
		fmt.Println("Starting State:", dfa.startingState)
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
			for symbol, toStateID := range fromState.transitions {
				fmt.Println(fromStateID, "--", symbol, "->", toStateID)
			}
		}
	}
}

func GetPTAFromListOfStringInstances(strings []StringInstance, APTA bool) DFA {
	strings = SortListOfStringInstances(strings)
	alphabet := make(map[int32]bool)
	var count int
	var currentStateID, newStateID int
	dfa := NewDFA()
	//dfa := DFA{
	//	states:   make(map[uint]State),
	//	alphabet: make(map[int32]bool),
	//}

	if strings[0].length == 0 {
		if strings[0].stringStatus == ACCEPTING {
			currentStateID = dfa.AddState(ACCEPTING)
		} else {
			currentStateID = dfa.AddState(REJECTING)
		}
	} else {
		currentStateID = dfa.AddState(UNKNOWN)
	}

	dfa.startingState = currentStateID

	for _, stringInstance := range strings {
		if !APTA && stringInstance.stringStatus != ACCEPTING {
			continue
		}
		currentStateID = dfa.startingState
		count = 0
		for _, symbol := range stringInstance.stringValue {
			count++
			// new alphabet check
			if !alphabet[symbol] {
				dfa.AddSymbol(symbol)
				alphabet[symbol] = true
			}

			symbolID := dfa.GetSymbolID(symbol)

			if dfa.states[currentStateID].transitions[symbolID] != -1 {
				currentStateID = dfa.states[currentStateID].transitions[symbolID]
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						if dfa.states[currentStateID].stateStatus == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						} else {
							dfa.states[currentStateID].UpdateStateStatus(ACCEPTING)
						}
					} else {
						if dfa.states[currentStateID].stateStatus == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							dfa.states[currentStateID].UpdateStateStatus(REJECTING)
						}
					}
				}
			} else {
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						newStateID = dfa.AddState(ACCEPTING)
					} else {
						newStateID = dfa.AddState(REJECTING)
					}
				} else {
					newStateID = dfa.AddState(UNKNOWN)
				}
				dfa.states[currentStateID].transitions[symbolID] = newStateID
				currentStateID = newStateID
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

func (dfa DFA) UnreachableStates() []int {
	reachableStates := map[int]bool{dfa.startingState: true}
	currentStates := map[int]bool{dfa.startingState: true}

	for len(currentStates) != 0 {
		nextStates := map[int]bool{}
		for stateID := range currentStates {
			for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
				if dfa.states[stateID].transitions[symbol] != -1{
					nextStates[dfa.states[stateID].transitions[symbol]] = true
				}
			}
		}
		// Don’t visit states we know to be reachable
		currentStates = map[int]bool{}
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

	var unReachableStates []int
	for stateID := range dfa.states {
		if !reachableStates[stateID] {
			unReachableStates = append(unReachableStates, stateID)
		}
	}

	return unReachableStates
}

func (dfa *DFA) RemoveUnreachableStates() {
	unreachableStates := dfa.UnreachableStates()
	for stateID := range unreachableStates {
		dfa.RemoveState(stateID)
	}
}

type StateIDPair struct {
	state1 int
	state2 int
}

func (dfa DFA) Mark() [][]int {
	distinguishablePairs := map[StateIDPair]bool{}

	dfa.RemoveUnreachableStates()
	for stateID, state := range dfa.states {
		for stateID2, state2 := range dfa.states {
			if stateID != stateID2 && state.stateStatus != state2.stateStatus &&
				!distinguishablePairs[StateIDPair{stateID, stateID2}] &&
				!distinguishablePairs[StateIDPair{stateID2, stateID}] {
				distinguishablePairs[StateIDPair{stateID, stateID2}] = true
			}
		}
	}

	oldCount := 0
	for oldCount != len(distinguishablePairs) {
		oldCount = len(distinguishablePairs)

		for stateID, state := range dfa.states {
			for stateID2, state2 := range dfa.states {
				if stateID == stateID2 || distinguishablePairs[StateIDPair{stateID, stateID2}] ||
					distinguishablePairs[StateIDPair{stateID2, stateID}] {
					continue
				} else {
					for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
						if state.transitions[symbol] != -1 {
							if state2.transitions[symbol] != -1 {
								if distinguishablePairs[StateIDPair{state.transitions[symbol], state2.transitions[symbol]}] ||
									distinguishablePairs[StateIDPair{state2.transitions[symbol], state.transitions[symbol]}] {
									distinguishablePairs[StateIDPair{stateID, stateID2}] = true
								}
							}
						}
					}
				}
			}
		}
	}

	var distinguishablePairsList [][]int
	for stateIDPair := range distinguishablePairs {
		distinguishablePairsList = append(distinguishablePairsList, []int{stateIDPair.state1, stateIDPair.state2})
	}
	return distinguishablePairsList
}
