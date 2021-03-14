package DFA_Toolkit

// StatePartition struct which represents a State Partition.
type StatePartition struct {
	root []int		// Parent block of each state.
	size []int		// Size (score) of each block.
	link []int		// Index of next state within the block.
	changed []int	// Slice of changed states/blocks.

	isCopy bool

	//temp
	labelledStatesCount int
	blockStatus []StateStatus
}

// Returns an initialized State Partition.
func NewStatePartition(dfa DFA) *StatePartition {
	// Initialize new State Partition struct.
	statePartition := new(StatePartition)
	statePartition.isCopy = false
	// Initialize empty slices.
	statePartition.labelledStatesCount = 0
	statePartition.root = make([]int, len(dfa.States))
	statePartition.size = make([]int, len(dfa.States))
	statePartition.link = make([]int, len(dfa.States))
	statePartition.blockStatus = make([]StateStatus, len(dfa.States))
	// Set root and link as element, and size (score) as 1.
	for i := 0; i < len(dfa.States); i++ {
		statePartition.root[i] = i
		statePartition.size[i] = 1
		statePartition.link[i] = i
		statePartition.blockStatus[i] = dfa.States[i].StateStatus
		if statePartition.blockStatus[i] == ACCEPTING || statePartition.blockStatus[i] == REJECTING{
			statePartition.labelledStatesCount++
		}
	}

	return statePartition
}

// Connects two states by finding their roots and comparing their respective
// size (score) values to keep the tree flat.
func (statePartition *StatePartition) union(stateID1Root int, stateID2Root int){
	// If their root is not equal, the states are merged (union) using the
	// linkBlocks function. If their root is equal, the states are already
	// within the same block so the merge is not done.
	if stateID1Root != stateID2Root{
		// Add State IDs joined to changed struct so merge can be undone.
		if statePartition.isCopy{
			statePartition.changed = append(statePartition.changed, stateID1Root)
			statePartition.changed = append(statePartition.changed, stateID2Root)
		}
		statePartition.linkBlocks(stateID1Root, stateID2Root)
	}
}

func (statePartition *StatePartition) linkBlocks(blockID1 int, blockID2 int){
	statePartition.link[blockID1], statePartition.link[blockID2] =
	 	statePartition.link[blockID2], statePartition.link[blockID1]

	block1Status := statePartition.blockStatus[blockID1]
	block2Status := statePartition.blockStatus[blockID2]

	if statePartition.size[blockID1] > statePartition.size[blockID2] {
		statePartition.root[blockID2] = blockID1
		statePartition.size[blockID1] ++

		if block1Status == UNKNOWN && block2Status != UNKNOWN{
			statePartition.blockStatus[blockID1] = block2Status
		}else if (block1Status == ACCEPTING || block1Status == REJECTING) &&
			(block2Status == ACCEPTING || block2Status == REJECTING){
			statePartition.labelledStatesCount--
		}
	}else{
		statePartition.root[blockID1] = blockID2
		statePartition.size[blockID2] ++

		if block2Status == UNKNOWN && block1Status != UNKNOWN{
			statePartition.blockStatus[blockID2] = block1Status
		}else if (block1Status == ACCEPTING || block1Status == REJECTING) &&
			(block2Status == ACCEPTING || block2Status == REJECTING){
			statePartition.labelledStatesCount--
		}
	}
}

// Find traverses each parent element while compressing the
// levels to find the root element of the stateID
// If we attempt to access an element outside the array it returns -1
func (statePartition *StatePartition) Find(stateID int) int {
	//if stateID > len(statePartition.root)-1 {
	//	panic("StateID out of range.")
	//}

	for statePartition.root[stateID] != stateID {
		statePartition.root[stateID] = statePartition.root[statePartition.root[stateID]]
		stateID = statePartition.root[stateID]
	}

	return stateID
}

func (statePartition StatePartition) ReturnSet(stateID int) []int{
	blockElements := []int{stateID}
	root := stateID
	for statePartition.link[stateID] != root{
		stateID = statePartition.link[stateID]
		blockElements = append(blockElements, stateID)
		if len(blockElements) > len(statePartition.root){
			panic("Error in state linking.")
		}
	}
	return blockElements
}

// Checks if states are within the same block
func (statePartition *StatePartition) WithinSameBlock(stateID1 int, stateID2 int) bool {
	return statePartition.Find(stateID1) == statePartition.Find(stateID2)
}

// Converts a DFA to a State Partition
func (dfa DFA) ToStatePartition() *StatePartition {
	// Return
	return NewStatePartition(dfa)
}

// Converts a State Partition to a DFA
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
		currentBlockID := statePartition.Find(stateID)
		if _, ok := newMappings[currentBlockID]; !ok {
			newMappings[currentBlockID] = resultantDFA.AddState(statePartition.blockStatus[currentBlockID])
		}
	}

	// update starting state via mappings
	resultantDFA.StartingStateID = newMappings[statePartition.Find(dfa.StartingStateID)]

	// update new transitions via mappings
	for stateID := range dfa.States{
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			oldResultantStateID := dfa.States[stateID].Transitions[symbolID]
			if oldResultantStateID > -1{
				newStateID := newMappings[statePartition.Find(stateID)]
				resultantStateID := newMappings[statePartition.Find(oldResultantStateID)]
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

// Recursively merges states to merge state1 and state2, returns false
// if the merge results in an NFA, or true if merge was successful
func (statePartition *StatePartition) MergeStates(dfa DFA, state1 int, state2 int) bool{
	state1Block := 0
	state2Block := 0

	if statePartition.root[state1] != state1{
		state1Block = statePartition.Find(state1)
	}else{
		state1Block = state1
	}

	if statePartition.root[state2] != state2{
		state2Block = statePartition.Find(state2)
	}else{
		state2Block = state2
	}

	// return true if states are already in the same block
	if state1Block == state2Block{
		return true
	}

	state1Status := statePartition.blockStatus[state1Block]
	state2Status := statePartition.blockStatus[state2Block]
	if (state1Status == ACCEPTING && state2Status == REJECTING) || (state1Status == REJECTING && state2Status == ACCEPTING){
		return false
	}

	// store the block transitions and set to -1 by default
	transitions := make([]int, len(dfa.SymbolMap))
	for i := range transitions {
		transitions[i] = -1
	}

	block1Set := statePartition.ReturnSet(state1Block)
	block2Set := statePartition.ReturnSet(state2Block)

	// merge states within state partition
	statePartition.union(state1Block, state2Block)

	// iterate over each state within the block containing the merged states
	for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
		for _, stateID := range block1Set{
			// store resultant state from state transition of current state
			currentResultantStateID := dfa.States[stateID].Transitions[symbolID]

			if currentResultantStateID > -1{
				// store resultant block from state transition of current block
				transitions[symbolID] = currentResultantStateID
				break
			}
		}
	}

	// iterate over each symbol within DFA
	for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
		if transitions[symbolID] == -1{
			continue
		}
		// iterate over each state within the block containing the merged states
		for _, stateID := range block2Set{
			// store resultant state from state transition of current state
			currentResultantStateID := dfa.States[stateID].Transitions[symbolID]
			if currentResultantStateID > -1{
				currentResultantBlockID := statePartition.Find(currentResultantStateID)
				if currentResultantBlockID != statePartition.Find(transitions[symbolID]){
					// not deterministic so merge, if states cannot be merged, return false
					if !statePartition.MergeStates(dfa, transitions[symbolID], currentResultantStateID) {
						return false
					}
				}
			}
		}
	}

	// return true if this is reached (deterministic)
	return true
}

func (statePartition StatePartition) Copy() *StatePartition{
	// Initialize new State Partition struct.
	copiedStatePartition := new(StatePartition)
	copiedStatePartition.isCopy = true
	copiedStatePartition.labelledStatesCount = statePartition.labelledStatesCount
	copiedStatePartition.root = make([]int, len(statePartition.root))
	copiedStatePartition.size = make([]int, len(statePartition.size))
	copiedStatePartition.link = make([]int, len(statePartition.link))
	copiedStatePartition.blockStatus = make([]StateStatus, len(statePartition.blockStatus))
	copiedStatePartition.changed = []int{}

	copy(copiedStatePartition.root, statePartition.root)
	copy(copiedStatePartition.size, statePartition.size)
	copy(copiedStatePartition.link, statePartition.link)
	copy(copiedStatePartition.blockStatus, statePartition.blockStatus)

	return copiedStatePartition
}

func (statePartition *StatePartition) RollbackChanges(originalStatePartition *StatePartition){
	if statePartition.isCopy{
		statePartition.labelledStatesCount = originalStatePartition.labelledStatesCount
		for _, stateID := range statePartition.changed{
			statePartition.root[stateID] = originalStatePartition.root[stateID]
			statePartition.size[stateID] = originalStatePartition.size[stateID]
			statePartition.link[stateID] = originalStatePartition.link[stateID]
			statePartition.blockStatus[stateID] = originalStatePartition.blockStatus[stateID]
		}
		statePartition.changed = []int{}
	}else{
		return
	}
}

func (statePartition StatePartition) EDSMScore() int{
	return statePartition.labelledStatesCount
}