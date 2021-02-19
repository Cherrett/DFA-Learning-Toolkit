package DFA_Toolkit

import "sync"

type NFAState struct {
	stateStatus StateStatus
	stateID     uint
	transitions map[int32]map[uint]bool
}

type NFA struct {
	states        map[uint]NFAState
	startingState NFAState
	alphabet      map[int32]bool
}

func (nfa NFA) Deterministic() bool{
	for _, state := range nfa.states{
		for _, transitions := range state.transitions{
			if len(transitions) > 1{
				return false
			}
		}
	}
	return true
}

func (nfa NFA) ToDFA() DFA{
	dfa := DFA{
		states:   make(map[uint]State),
		alphabet: nfa.alphabet,
	}
	for _, state := range nfa.states{
		dfa.states[state.stateID] = State{
			stateStatus: state.stateStatus,
			stateID:     state.stateID,
			transitions: make(map[int32]uint),
		}
		for char, states := range state.transitions {
			if len(states) > 1{
				panic("Cannot convert NFA to DFA since it is non-deterministic")
			}else{
				for state2 := range states{
					dfa.states[state.stateID].transitions[char] = state2
				}
			}
		}
	}

	dfa.startingState = dfa.states[nfa.startingState.stateID]
	return dfa
}

func (nfa *NFA) UpdateStateStatus(stateID uint, stateStatus StateStatus) {
	tempState := nfa.states[stateID]
	tempState.stateStatus = stateStatus
	nfa.states[stateID] = tempState
}

func RPNIDerive(dfa DFA, partition []map[uint]bool) NFA{
	newMappings := map[uint]uint{}
	nfa := NFA{
		states:        make(map[uint]NFAState),
		alphabet:      dfa.alphabet,
	}

	for _, currentBlock := range partition{
		currentState := NFAState{
			stateStatus: UNKNOWN,
			stateID:     uint(len(nfa.states)),
			transitions: make(map[int32]map[uint]bool),
		}

		for stateID := range currentBlock{
			if dfa.states[stateID].stateStatus == ACCEPTING{
				currentState.stateStatus = ACCEPTING
			}
			for char, stateID2 := range dfa.states[stateID].transitions{
				if len(currentState.transitions[char]) == 0{
					currentState.transitions[char] = map[uint]bool{stateID2: true}
				}else{
					currentState.transitions[char][stateID2] = true
				}
			}
			newMappings[stateID] = currentState.stateID
		}
		nfa.states[currentState.stateID] = currentState
	}
	// update new states via mappings
	for _, state := range nfa.states{
		for char, states := range state.transitions{
			for stateID2 := range states{
				if stateID2 != newMappings[stateID2]{
					state.transitions[char][newMappings[stateID2]] = true
					delete(state.transitions[char], stateID2)
				}
			}
		}
	}

	return nfa
}

func RPNIMerge(nfa NFA, state1 uint, state2 uint) NFA{
	var stateStatus StateStatus = UNKNOWN
	for stateID, state := range nfa.states{
		if stateID == state1{
			if state.stateStatus == ACCEPTING{
				stateStatus = ACCEPTING
			}
		}else if stateID == state2{
			if state.stateStatus == ACCEPTING{
				stateStatus = ACCEPTING
			}

			for char, transitions := range state.transitions{
				for transition := range transitions{
					if _, ok := nfa.states[state1].transitions[char]; ok{
						nfa.states[state1].transitions[char][transition] = true
					}else{
						// temporary solution
						tempState := nfa.states[state1]
						tempState.transitions = map[int32]map[uint]bool{char: {transition: true}}
						nfa.states[state1] = tempState
					}
				}
			}
			continue
		}

		for char, transitions := range state.transitions{
			for transition := range transitions{
				if transition == state2{
					if _, ok := state.transitions[char][state1]; !ok {
						nfa.states[stateID].transitions[char][state1] = true
					}
					delete(nfa.states[stateID].transitions[char], state2)
				}
			}
		}

	}
	delete(nfa.states, state2)
	nfa.UpdateStateStatus(state1, stateStatus)

	if nfa.startingState.stateID == state2{
		nfa.startingState.stateID = state1
	}

	return nfa
}

func RPNIDeterministicMerge(nfa NFA, partition []map[uint]bool) (DFA, []map[uint]bool){
	for !nfa.Deterministic(){
		exitLoop := false
		for _, state := range nfa.states{
			for _, transitions := range state.transitions {
				if len(transitions) > 1 {

					var stateIDsToMerge = []uint{}
					for stateID := range transitions{
						if len(stateIDsToMerge) < 2{
							stateIDsToMerge = append(stateIDsToMerge, stateID)
						}else{
							break
						}
					}
					nfa = RPNIMerge(nfa, stateIDsToMerge[0], stateIDsToMerge[1])

					for blockIndex, block := range partition{
						if _, ok := block[stateIDsToMerge[1]]; ok{
							partition[blockIndex][stateIDsToMerge[1]] = true
							break
						}
					}

					for blockIndex, block := range partition{
						if _, ok := block[stateIDsToMerge[0]]; ok{
							// replace unwanted block with last block in partition
							partition[blockIndex] = partition[len(partition)-1]
							// remove last block
							partition = partition[:len(partition)-1]
							break
						}
					}

					// need to create normalization function after merge
					exitLoop = true
					break
				}
			}
			if exitLoop{
				break
			}
		}
	}

	return nfa.ToDFA(), partition
}

func RPNIStringInstanceConsistentWithDFA(stringInstance StringInstance, dfa DFA) bool{
	if stringInstance.length == 0{
		return !(dfa.startingState.stateStatus == ACCEPTING)
	}

	currentState := dfa.startingState
	var count uint = 0
	for _, character := range stringInstance.stringValue{
		count++
		if value, ok := currentState.transitions[character]; ok {
			currentState = dfa.states[value]
			// last symbol in string check
			if count == stringInstance.length {
				if currentState.stateStatus == ACCEPTING {
					return false
				}
			}
		}
	}

	return true
}

func RPNIListOfStringInstancesConsistentWithDFA(stringInstances []StringInstance, dfa DFA) bool{
	consistent := true
	var wg sync.WaitGroup
	wg.Add(len(stringInstances))

	for _, stringInstance := range stringInstances {
		go func(stringInstance StringInstance, dfa DFA){
			defer wg.Done()
			if consistent {
				consistent = RPNIStringInstanceConsistentWithDFA(stringInstance, dfa)
			}
		}(stringInstance, dfa)
	}

	wg.Wait()
	return consistent
}

func RPNI(acceptingStrings []StringInstance, rejectingStrings[]StringInstance) DFA{
	PTA := GetPTAFromListOfStringInstances(acceptingStrings, false)
	PTA.Describe(false)

	currentHypothesis := PTA
	tempHypothesis := PTA
	var currentPartition []map[uint]bool

	for stateID := range currentHypothesis.states{
		currentPartition = append(currentPartition, map[uint]bool{stateID: true})
	}

	for i := uint(1); i < uint(len(PTA.states)); i++ {
		for j := uint(0); j < i; j++ {
			// merge the block which contains state i with the block which contains state j
			var tempPartition []map[uint]bool
			tempBlock := map[uint]bool{}

			for _, block := range currentPartition{
				if block[i]{
					for stateID := range block{
						tempBlock[stateID] = true
					}
				}else if block[j]{
					for stateID := range block{
						tempBlock[stateID] = true
					}
				}else{
					tempPartition = append(tempPartition, block)
				}
			}
			tempPartition = append(tempPartition, tempBlock)

			//for _, m := range tempPartition{
			//	fmt.Print("(")
			//	for stateID := range m{
			//		fmt.Print(stateID, " ")
			//	}
			//	fmt.Print(")")
			//}
			//fmt.Println()

			// get quotient automaton
			tempHypothesisNFA := RPNIDerive(PTA, tempPartition)

			if tempHypothesisNFA.Deterministic(){
				tempHypothesis = tempHypothesisNFA.ToDFA()
			}else{
				// determine the quotient automaton (if necessary) by state merging
				tempHypothesis, tempPartition = RPNIDeterministicMerge(tempHypothesisNFA, tempPartition)
			}

			if RPNIListOfStringInstancesConsistentWithDFA(rejectingStrings, tempHypothesis) {
				// Treat tempHypothesis as the current hypothesis
				currentHypothesis = tempHypothesis
				currentPartition = tempPartition
				break
			}
		}
	}

	return currentHypothesis
}
