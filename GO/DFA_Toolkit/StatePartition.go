package DFA_Toolkit

type StatePartition struct {
	root []int
	size []int
	changed []int
}

// Returns an initialized State Partition
func NewStatePartition(size int) *StatePartition {
	statePartition := new(StatePartition)
	statePartition.root = make([]int, size)
	statePartition.size = make([]int, size)
	statePartition.changed = []int{}

	for i := 0; i < size; i++ {
		statePartition.root[i] = i
		statePartition.size[i] = 1
	}

	return statePartition
}

// Union connects p and q by finding their roots and comparing their respective
// size arrays to keep the tree flat
func (statePartition *StatePartition) union(stateID1 int, stateID2 int) {
	stateID2Root := statePartition.Find(stateID2)
	stateID1Root := statePartition.Find(stateID1)

	if stateID1Root == stateID2Root{
		return
	}

	if statePartition.size[stateID2Root] < statePartition.size[stateID1Root] {
		statePartition.root[stateID2Root] = statePartition.root[stateID1Root]
		statePartition.size[stateID1Root] += statePartition.size[stateID2Root]
	} else {
		statePartition.root[stateID1Root] = statePartition.root[stateID2Root]
		statePartition.size[stateID2Root] += statePartition.size[stateID1Root]
	}

	statePartition.changed = append(statePartition.changed, stateID1)
	statePartition.changed = append(statePartition.changed, stateID2)
}

// Find traverses each parent element while compressing the
// levels to find the root element of p
// If we attempt to access an element outside the array it returns -1
func (statePartition *StatePartition) Find(p int) int {
	if p > len(statePartition.root)-1 {
		return -1
	}

	for statePartition.root[p] != p {
		statePartition.root[p] = statePartition.root[statePartition.root[p]]
		p = statePartition.root[p]
	}

	return p
}

// Check if items p,q are connected
func (statePartition *StatePartition) Connected(stateID1 int, stateID2 int) bool {
	return statePartition.Find(stateID1) == statePartition.Find(stateID2)
}

//func (statePartition StatePartition) GetSet(setID int) []int{
//	var set []int
//	root := setID
//	for set
//
//	return set
//}

// Convert a DFA to a State Partition
func (dfa DFA) ToStatePartition() *StatePartition {
	return NewStatePartition(len(dfa.States))
}

// Convert a State Partition to a DFA
func (statePartition StatePartition) ToDFA(dfa DFA) (bool, DFA){
	newMappings := map[int]int{}

	resultantDFA := DFA{
		States:                nil,
		StartingStateID:       -1,
		SymbolMap:             dfa.SymbolMap,
		Depth:                 -1,
		ComputedDepthAndOrder: false,
	}

	for stateID := range dfa.States {
		if newStateID, ok := newMappings[statePartition.Find(stateID)]; ok {
			if (resultantDFA.States[newStateID].StateStatus == ACCEPTING &&
				dfa.States[stateID].StateStatus == REJECTING) ||
				(resultantDFA.States[newStateID].StateStatus == REJECTING &&
					dfa.States[stateID].StateStatus == ACCEPTING){
				// not deterministic
				return false, DFA{}
			}else{
				resultantDFA.States[newStateID].StateStatus = dfa.States[stateID].StateStatus
			}
		}else{
			newMappings[statePartition.Find(stateID)] = resultantDFA.AddState(dfa.States[stateID].StateStatus)
		}
	}

	// update starting state via mappings
	resultantDFA.StartingStateID = newMappings[statePartition.Find(dfa.StartingStateID)]

	// update new transitions via mappings
	for stateID := range dfa.States{
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			if dfa.States[stateID].Transitions[symbolID] > -1{
				newStateID := newMappings[statePartition.Find(stateID)]
				resultantStateID := newMappings[statePartition.Find(dfa.States[stateID].Transitions[symbolID])]
				if resultantDFA.States[newStateID].Transitions[symbolID] > -1 &&
					resultantDFA.States[newStateID].Transitions[symbolID] != resultantStateID{
					// not deterministic
					return false, DFA{}
				}else{
					resultantDFA.States[newStateID].Transitions[symbolID] = resultantStateID
				}
			}
		}
	}
	return true, resultantDFA
}

// Recursively merge states to merge state1 and state2, returns false
// if the merge results in an NFA, or true if merge was successful
func (statePartition *StatePartition) MergeStates(dfa DFA, state1 int, state2 int) bool{
	state1Block := statePartition.Find(state1)
	state2Block := statePartition.Find(state2)
	// return same state partition if states are already in the same block
	if state1Block == state2Block{
		return true
	}

	var statesToBeMerged [][]int
	var block1Status StateStatus = UNKNOWN
	var block2Status StateStatus = UNKNOWN
	transitions := make([]int, len(dfa.SymbolMap))
	for i := range transitions {
		transitions[i] = -1
	}

	for stateID := range dfa.States{
		stateBlockID := statePartition.Find(stateID)
		if stateBlockID == state1Block{
			if block1Status == UNKNOWN{
				block1Status = dfa.States[stateID].StateStatus
			}else if (block1Status == ACCEPTING && dfa.States[stateID].StateStatus == REJECTING) ||
				(block1Status == REJECTING && dfa.States[stateID].StateStatus == ACCEPTING){
				// not deterministic
				return false
			}
			for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
				if transitions[symbolID] > -1 && dfa.States[stateID].Transitions[symbolID] > -1{
					resultantStateBlockID := statePartition.Find(dfa.States[stateID].Transitions[symbolID])
					if transitions[symbolID] != resultantStateBlockID{
						// not deterministic so merge
						statesToBeMerged = append(statesToBeMerged, []int{transitions[symbolID], resultantStateBlockID})
					}
				}else{
					if dfa.States[stateID].Transitions[symbolID] > -1{
						transitions[symbolID] = statePartition.Find(dfa.States[stateID].Transitions[symbolID])
					}
				}
			}
		}else if stateBlockID == state2Block{
			if block2Status == UNKNOWN{
				block2Status = dfa.States[stateID].StateStatus
			}else if (block2Status == ACCEPTING && dfa.States[stateID].StateStatus == REJECTING) ||
				(block2Status == REJECTING && dfa.States[stateID].StateStatus == ACCEPTING){
				return false
			}
			for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
				if transitions[symbolID] > -1 && dfa.States[stateID].Transitions[symbolID] > -1{
					resultantStateBlockID := statePartition.Find(dfa.States[stateID].Transitions[symbolID])
					if transitions[symbolID] != resultantStateBlockID{
						// not deterministic so merge
						statesToBeMerged = append(statesToBeMerged, []int{transitions[symbolID], resultantStateBlockID})
					}
				}else{
					if dfa.States[stateID].Transitions[symbolID] > -1{
						transitions[symbolID] = statePartition.Find(dfa.States[stateID].Transitions[symbolID])
					}
				}
			}
		}
	}

	// merge state within state partition
	statePartition.union(state1, state2)

	// merge conflicting states
	for pairID := range statesToBeMerged{
		if !statePartition.MergeStates(dfa, statesToBeMerged[pairID][0], statesToBeMerged[pairID][1]) {
			return false
		}
	}

	return true
}

func (statePartition StatePartition) Copy() *StatePartition{
	copiedStatePartition := NewStatePartition(len(statePartition.size))
	copiedStatePartition.size = append(copiedStatePartition.size, statePartition.size...)
	copiedStatePartition.root = append(copiedStatePartition.root, statePartition.root...)
	return copiedStatePartition
}

func (statePartition *StatePartition) RollbackChanges(originalStatePartition *StatePartition){
	for _, stateID := range statePartition.changed{
		statePartition.root[stateID] = originalStatePartition.root[stateID]
		statePartition.size[stateID] = originalStatePartition.size[stateID]
	}
	statePartition.changed = []int{}
}