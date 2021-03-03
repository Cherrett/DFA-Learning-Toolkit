package DFA_Toolkit

type StatePartition struct {
	root []int
	size []int
	link []int
	changed []int
}

// Returns an initialized State Partition
func NewStatePartition(size int) *StatePartition {
	statePartition := new(StatePartition)
	statePartition.root = make([]int, size)
	statePartition.size = make([]int, size)
	statePartition.link = make([]int, size)
	statePartition.changed = []int{}

	for i := 0; i < size; i++ {
		statePartition.root[i] = i
		statePartition.size[i] = 1
		statePartition.link[i] = i
	}

	return statePartition
}

// Union connects p and q by finding their roots and comparing their respective
// size arrays to keep the tree flat
func (statePartition *StatePartition) union(stateID1 int, stateID2 int) int{
	stateID1Root := statePartition.Find(stateID1)
	stateID2Root := statePartition.Find(stateID2)

	if stateID1Root != stateID2Root{
		statePartition.changed = append(statePartition.changed, stateID1)
		statePartition.changed = append(statePartition.changed, stateID2)
		return statePartition.linkBlocks(stateID1Root, stateID2Root)
	}else{
		return stateID1Root
	}
}

func (statePartition *StatePartition) linkBlocks(blockID1 int, blockID2 int) int{
	statePartition.link[blockID1], statePartition.link[blockID2] =
	 	statePartition.link[blockID2], statePartition.link[blockID1]

	if statePartition.size[blockID1] > statePartition.size[blockID2] {
		statePartition.root[blockID2] = blockID1
		return blockID1
	}else{
		statePartition.root[blockID1] = blockID2
		if statePartition.size[blockID1] == statePartition.size[blockID2]{
			statePartition.size[blockID2]++
		}
		return blockID2
	}
}

// Find traverses each parent element while compressing the
// levels to find the root element of p
// If we attempt to access an element outside the array it returns -1
func (statePartition *StatePartition) Find(stateID int) int {
	if stateID > len(statePartition.root)-1 {
		return -1
	}

	for statePartition.root[stateID] != stateID {
		statePartition.root[stateID] = statePartition.root[statePartition.root[stateID]]
		stateID = statePartition.root[stateID]
	}

	return stateID
}

func (statePartition StatePartition) ReturnSet(blockID int) []int{
	blockElements := []int{blockID}
	root := blockID
	for statePartition.link[blockID] != root{
		blockID = statePartition.link[blockID]
		blockElements = append(blockElements, blockID)
	}
	return blockElements
}

// Check if items p,q are connected
func (statePartition *StatePartition) Connected(stateID1 int, stateID2 int) bool {
	return statePartition.Find(stateID1) == statePartition.Find(stateID2)
}

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
	// return same state partition if states are already in the same block
	if statePartition.Find(state1) == statePartition.Find(state2){
		return true
	}

	//
	var blockStatus StateStatus = UNKNOWN
	//var block2Status StateStatus = UNKNOWN
	transitions := make([]int, len(dfa.SymbolMap))
	for i := range transitions {
		transitions[i] = -1
	}

	// merge state within state partition
	mergedStateBlock := statePartition.union(state1, state2)

	for _, stateID := range statePartition.ReturnSet(mergedStateBlock){
		currentStateStatus := dfa.States[stateID].StateStatus
		if blockStatus == UNKNOWN && currentStateStatus != UNKNOWN{
			blockStatus = currentStateStatus
		}else if (blockStatus == ACCEPTING && currentStateStatus == REJECTING) ||
			(blockStatus == REJECTING && currentStateStatus == ACCEPTING){
			// not deterministic
			return false
		}
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			currentResultantStateID := dfa.States[stateID].Transitions[symbolID]
			resultantStateID := transitions[symbolID]
			if resultantStateID > -1 && currentResultantStateID > -1{
				resultantStateBlockID := statePartition.Find(currentResultantStateID)
				if resultantStateID != resultantStateBlockID{
					// not deterministic so merge, if states cannot be merged, return false
					if !statePartition.MergeStates(dfa, resultantStateID, resultantStateBlockID) {
						return false
					}
				}
			}else{
				if currentResultantStateID > -1{
					transitions[symbolID] = statePartition.Find(currentResultantStateID)
				}
			}
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
		statePartition.link[stateID] = originalStatePartition.link[stateID]
	}
	statePartition.changed = []int{}
}