package DFA_Toolkit

import "fmt"

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
		for _, states := range state.transitions {
			if len(states) > 1{
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
				currentState.transitions[char] = map[uint]bool{stateID2: true}
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

func RPNI(acceptingStrings []StringInstance, rejectingStrings[]StringInstance) DFA{
	PTA := GetPTAFromListOfStringInstances(acceptingStrings, false)
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

			for _, m := range tempPartition{
				fmt.Print("(")
				for stateID := range m{
					fmt.Print(stateID, " ")
				}
				fmt.Print(")")
			}
			fmt.Println()

			// get quotient automaton
			tempHypothesisNFA := RPNIDerive(PTA, tempPartition)

			if tempHypothesisNFA.Deterministic(){
				tempHypothesis = tempHypothesisNFA.ToDFA()
			}else{
				panic("not deterministic, all ok")
				// TODO: determine the quotient automaton (if necessary) by state merging
			}

			if ListOfStringInstancesConsistentWithDFA(rejectingStrings, tempHypothesis) {
				// Treat tempHypothesis as the current hypothesis
				currentHypothesis = tempHypothesis
				currentPartition = tempPartition
			}
		}
	}

	return currentHypothesis
}
