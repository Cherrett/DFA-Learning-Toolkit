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

	// Iterate over all unique state pairs within DFA where stateID != stateID2.
	for stateID := range dfa.States {
		for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
			// If the state pair have different types, add to distinguishable pairs map.
			if dfa.States[stateID].Label != dfa.States[stateID2].Label {
				// Make sure that the smaller stateID is in first position within pair
				if stateID <= stateID2{
					distinguishablePairs[StateIDPair{stateID, stateID2}] = true
				}else{
					distinguishablePairs[StateIDPair{stateID2, stateID}] = true
				}
			}
		}
	}

	// Set counter for length of distinguishable pairs to 0.
	oldCount := -1

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
					for symbol := range dfa.Alphabet {
						resultantStateID1 := dfa.States[stateID].Transitions[symbol]
						resultantStateID2 := dfa.States[stateID2].Transitions[symbol]
						// If both states have a valid transition using current symbol.
						if resultantStateID1 != -1 && resultantStateID2 != -1 {
							// If pair containing both resultant state IDs is marked as distinguishable,
							// mark current state pair as distinguishable.
							if resultantStateID1 <= resultantStateID2{
								if distinguishablePairs[StateIDPair{resultantStateID1, resultantStateID2}]{
									distinguishablePairs[StateIDPair{stateID, stateID2}] = true
								}
							}else{
								if distinguishablePairs[StateIDPair{resultantStateID2, resultantStateID1}]{
									distinguishablePairs[StateIDPair{stateID, stateID2}] = true
								}
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
	// Check that DFA does not have any rejecting states (only accepting and unlabelled).
	if dfa.RejectingStatesCount() != 0 {
		panic("The DFA contains rejecting states so cannot be minimised.")
	}

	// Clone DFA to work on.
	temporaryDFA := dfa.Clone()

	// Remove unreachable states within DFA.
	temporaryDFA.RemoveUnreachableStates()

	// If DFA is not complete, add sink state to make it complete.
	if !temporaryDFA.IsComplete() {
		temporaryDFA.AddSinkState()
	}

	// Get indistinguishable state pairs from IndistinguishableStatePairs function.
	indistinguishablePairs := temporaryDFA.IndistinguishableStatePairs()

	// Convert DFA to state partition.
	statePartition := temporaryDFA.ToStatePartition()

	// Merge indistinguishable pairs.
	for _, indistinguishablePair := range indistinguishablePairs{
		block1 := statePartition.Find(indistinguishablePair.state1)
		block2 := statePartition.Find(indistinguishablePair.state2)
		if block1 != block2{
			statePartition.Union(block1, block2)
		}
	}

	// Convert state partition to Quotient DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Return resultant minimised DFA.
	return resultantDFA
}

// AddSinkState adds a sink state within DFA. This is used to convert
// a non-complete DFA to a complete DFA by adding a sink state.
func (dfa *DFA) AddSinkState() int {
	// Return if dfa is complete since sink state is not required.
	if dfa.IsComplete(){
		return -1
	}

	// Add sink state and store its ID.
	sinkStateID := dfa.AddState(UNLABELLED)

	// Iterate over alphabet within DFA.
	for symbolID := range dfa.Alphabet {
		// Add a self-transition using symbol.
		dfa.AddTransition(symbolID, sinkStateID, sinkStateID)
	}

	// Iterate over each state within DFA.
	for stateID := range dfa.States {
		// Iterate over alphabet within DFA.
		for symbolID := range dfa.Alphabet {
			// If a transition does not exist with symbol.
			if dfa.States[stateID].Transitions[symbolID] == -1 {
				// Add a transition to sink state using symbol.
				dfa.States[stateID].Transitions[symbolID] = sinkStateID
			}
		}
	}

	// Return ID of sink state.
	return sinkStateID
}
