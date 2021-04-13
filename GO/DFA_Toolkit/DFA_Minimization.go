package dfatoolkit

// StateIDPair which represents a pair of states.
type StateIDPair struct {
	state1 int
	state2 int
}

// IndistinguishableStatePairs returns a slice of indistinguishable state pairs within DFA.
// P. Linz, An Introduction to Formal Languages and Automata. Jones & Bartlett Publishers, 2011.
func (dfa *DFA) IndistinguishableStatePairs() []StateIDPair {
	// Map to store distinguishable state pairs.
	distinguishablePairs := map[StateIDPair]bool{}

	// Remove unreachable states within DFA.
	dfa.RemoveUnreachableStates()

	// Iterate over all unique state pairs within DFA where stateID != stateID2.
	for stateID := range dfa.States {
		for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
			// If the state pair have different types, add to distinguishable pairs map.
			if dfa.States[stateID].Label != dfa.States[stateID2].Label {
				distinguishablePairs[StateIDPair{stateID, stateID2}] = true
			}
		}
	}

	// Set counter for length of distinguishable pairs to 0.
	oldCount := 0

	// Iterate until no new pairs have been added to distinguishable pairs.
	for oldCount != len(distinguishablePairs) {
		// Set counter to length of distinguishable pairs.
		oldCount = len(distinguishablePairs)

		// Iterate over all unique state pairs within DFA where stateID != stateID2.
		for stateID := range dfa.States {
			for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
				// If state pair is already marked as distinguishable, skip.
				if distinguishablePairs[StateIDPair{stateID, stateID2}] {
					continue
				} else {
					// Iterate over each symbol within DFA.
					for alphabetID := range dfa.Alphabet {
						// If both states have a valid transition using current symbol.
						if dfa.States[stateID].Transitions[alphabetID] != -1 &&
							dfa.States[stateID2].Transitions[alphabetID] != -1 {
							// If pair containing both resultant state IDs is marked as distinguishable,
							// mark current state pair as distinguishable.
							if distinguishablePairs[StateIDPair{dfa.States[stateID].Transitions[alphabetID],
								dfa.States[stateID2].Transitions[alphabetID]}] ||
								distinguishablePairs[StateIDPair{dfa.States[stateID2].Transitions[alphabetID],
									dfa.States[stateID].Transitions[alphabetID]}] {
								distinguishablePairs[StateIDPair{stateID, stateID2}] = true
							}
						}
					}
				}
			}
		}
	}

	// Slice to store indistinguishable pairs. The size of this slice is all
	// state pairs (n(n-1)/2 is used via the triangle number method) minus the
	// number of distinguishable state pairs.
	indistinguishablePairs := make([]StateIDPair, ((len(dfa.States)*(len(dfa.States)-1))/2)-len(distinguishablePairs))

	// Set counter to 0.
	count := 0

	// Iterate over all unique state pairs within DFA where stateID != stateID2.
	for stateID := range dfa.States {
		for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
			// If state pair is not marked as distinguishable, add to indistinguishable pairs.
			if !distinguishablePairs[StateIDPair{stateID, stateID2}] {
				indistinguishablePairs[count] = StateIDPair{stateID, stateID2}
				// Increment count.
				count++
			}
		}
	}

	// Return slice of indistinguishable state pairs.
	return indistinguishablePairs
}

// Minimise returns a minimised version of the DFA.
// P. Linz, An Introduction to Formal Languages and Automata. Jones & Bartlett Publishers, 2011.
func (dfa DFA) Minimise() DFA {
	// Get indistinguishable state pairs from Mark function.
	indistinguishablePairs := dfa.IndistinguishableStatePairs()

	// Partition states into blocks.
	var currentPartition []map[int]bool
	for stateID := range dfa.States {
		exists := false
		for _, block := range currentPartition {
			if block[stateID] {
				exists = true
			}
		}
		if !exists {
			indistinguishable := false
			for _, indistinguishablePair := range indistinguishablePairs {
				if indistinguishablePair.state1 == stateID {
					for blockIndex, block := range currentPartition {
						if block[indistinguishablePair.state2] {
							currentPartition[blockIndex][stateID] = true
							indistinguishable = true
							break
						}
					}
					if indistinguishable {
						break
					}
				} else if indistinguishablePair.state2 == stateID {
					for blockIndex, block := range currentPartition {
						if block[indistinguishablePair.state1] {
							currentPartition[blockIndex][stateID] = true
							indistinguishable = true
							break
						}
					}
					if indistinguishable {
						break
					}
				}
			}
			if !indistinguishable {
				currentPartition = append(currentPartition, map[int]bool{stateID: true})
			}
		}
	}
	resultantDFA := DFA{Alphabet: dfa.Alphabet}

	// Create a new state for each block.
	for blockIndex := range currentPartition {
		var stateLabel StateLabel = UNKNOWN
		for stateID := range currentPartition[blockIndex] {
			stateLabel = dfa.States[stateID].Label
			break
		}
		resultantDFA.AddState(stateLabel)
	}

	// Set Initial State.
	for blockIndex := range currentPartition {
		found := false
		for stateID := range currentPartition[blockIndex] {
			if stateID == dfa.StartingStateID {
				resultantDFA.StartingStateID = blockIndex
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// Create Transitions.
	for stateID := range dfa.States {
		stateBlockIndex := 0
		found := false
		for blockIndex := range currentPartition {
			for stateID2 := range currentPartition[blockIndex] {
				if stateID == stateID2 {
					stateBlockIndex = blockIndex
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		for alphabetID := range dfa.States[stateID].Transitions {
			resultantStateBlockIndex := 0
			if resultantDFA.States[stateBlockIndex].Transitions[alphabetID] == -1 && dfa.States[stateID].Transitions[alphabetID] != -1 {
				found := false
				for blockIndex := range currentPartition {
					for stateID2 := range currentPartition[blockIndex] {
						if dfa.States[stateID].Transitions[alphabetID] == stateID2 {
							resultantStateBlockIndex = blockIndex
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				resultantDFA.States[stateBlockIndex].Transitions[alphabetID] = resultantStateBlockIndex
			}
		}
	}

	// Return resultant minimised DFA.
	return resultantDFA
}
