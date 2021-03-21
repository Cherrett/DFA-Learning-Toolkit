package DFA_Toolkit

type StateIDPair struct {
	state1 int
	state2 int
}

func (dfa *DFA) Mark() (map[StateIDPair]bool, map[StateIDPair]bool) {
	distinguishablePairs := map[StateIDPair]bool{}
	indistinguishablePairs := map[StateIDPair]bool{}
	allPairs := map[StateIDPair]bool{}

	dfa.RemoveUnreachableStates()
	for stateID := range dfa.States {
		for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
			allPairs[StateIDPair{stateID, stateID2}] = true
			if dfa.States[stateID].StateStatus != dfa.States[stateID2].StateStatus {
				distinguishablePairs[StateIDPair{stateID, stateID2}] = true
			}
		}
	}

	oldCount := 0
	for oldCount != len(distinguishablePairs) {
		oldCount = len(distinguishablePairs)

		for stateID := range dfa.States {
			for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
				if distinguishablePairs[StateIDPair{stateID, stateID2}] {
					continue
				} else {
					for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
						if dfa.States[stateID].Transitions[symbolID] != -1 {
							if dfa.States[stateID2].Transitions[symbolID] != -1 {
								if distinguishablePairs[StateIDPair{dfa.States[stateID].Transitions[symbolID], dfa.States[stateID2].Transitions[symbolID]}] ||
									distinguishablePairs[StateIDPair{dfa.States[stateID2].Transitions[symbolID], dfa.States[stateID].Transitions[symbolID]}] {
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

	for pair := range allPairs{
		if !distinguishablePairs[StateIDPair{pair.state1, pair.state2}]{
			indistinguishablePairs[StateIDPair{pair.state1, pair.state2}] = true
		}
	}

	return distinguishablePairs, indistinguishablePairs
}

func (dfa DFA) Minimise() DFA{
	// get distinguishable and indistinguishable state pairs using Mark function
	_, indistinguishablePairs := dfa.Mark()

	// Partition states into blocks
	var currentPartition []map[int]bool
	for stateID := range dfa.States {
		exists := false
		for _, block := range currentPartition{
			if block[stateID]{
				exists = true
			}
		}
		if !exists{
			indistinguishable := false
			for indistinguishablePair := range indistinguishablePairs{
				if indistinguishablePair.state1 == stateID{
					for blockIndex, block := range currentPartition{
						if block[indistinguishablePair.state2]{
							currentPartition[blockIndex][stateID] = true
							indistinguishable = true
							break
						}
					}
					if indistinguishable{
						break
					}
				}else if indistinguishablePair.state2 == stateID{
					for blockIndex, block := range currentPartition{
						if block[indistinguishablePair.state1]{
							currentPartition[blockIndex][stateID] = true
							indistinguishable = true
							break
						}
					}
					if indistinguishable{
						break
					}
				}
			}
			if !indistinguishable{
				currentPartition = append(currentPartition, map[int]bool{stateID: true})
			}
		}
	}
	resultantDFA := DFA{SymbolMap: dfa.SymbolMap}

	// Create a new state for each block
	for blockIndex := range currentPartition{
		var stateStatus StateStatus = UNKNOWN
		for stateID := range currentPartition[blockIndex] {
			stateStatus = dfa.States[stateID].StateStatus
			break
		}
		resultantDFA.AddState(stateStatus)
	}

	// Initial State
	for blockIndex := range currentPartition{
		found := false
		for stateID := range currentPartition[blockIndex] {
			if stateID == dfa.StartingStateID {
				resultantDFA.StartingStateID = blockIndex
				found = true
				break
			}
		}
		if found{
			break
		}
	}

	// Transitions
	for stateID := range dfa.States {
		stateBlockIndex := 0
		found := false
		for blockIndex := range currentPartition{
			for stateID2 := range currentPartition[blockIndex] {
				if stateID == stateID2{
					stateBlockIndex = blockIndex
					found = true
					break
				}
			}
			if found{
				break
			}
		}

		for symbolID := range dfa.States[stateID].Transitions {
			resultantStateBlockIndex := 0
			if resultantDFA.States[stateBlockIndex].Transitions[symbolID] == -1 && dfa.States[stateID].Transitions[symbolID] != -1{
				found := false
				for blockIndex := range currentPartition{
					for stateID2 := range currentPartition[blockIndex] {
						if dfa.States[stateID].Transitions[symbolID] == stateID2{
							resultantStateBlockIndex = blockIndex
							found = true
							break
						}
					}
					if found{
						break
					}
				}
				resultantDFA.States[stateBlockIndex].Transitions[symbolID] = resultantStateBlockIndex
			}
		}
	}

	return resultantDFA
}