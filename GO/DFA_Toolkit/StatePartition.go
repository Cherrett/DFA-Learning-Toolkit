package DFA_Toolkit

// StatePartition struct which represents a State Partition.
type StatePartition struct {
	root []int		// Parent block of each state.
	size []int		// Size (score) of each block.
	link []int		// Index of next state within the block.
	changed []int	// Slice of changed states/blocks.
}

// Returns an initialized State Partition.
func NewStatePartition(size int) *StatePartition {
	// Initialize new State Partition struct.
	statePartition := new(StatePartition)
	// Initialize empty slices.
	statePartition.root = make([]int, size)
	statePartition.size = make([]int, size)
	statePartition.link = make([]int, size)
	statePartition.changed = []int{}

	// Set root and link as element, and size (score) as 1.
	for i := 0; i < size; i++ {
		statePartition.root[i] = i
		statePartition.size[i] = 1
		statePartition.link[i] = i
	}

	return statePartition
}

// Connects two states by finding their roots and comparing their respective
// size (score) values to keep the tree flat. Returns the parent element.
func (statePartition *StatePartition) union(stateID1 int, stateID2 int) int{
	// Get root (block index) of both states.
	stateID1Root := statePartition.Find(stateID1)
	stateID2Root := statePartition.Find(stateID2)

	// If their root is not equal, the states are merged (union) using
	// linkBlocks function.
	if stateID1Root != stateID2Root{
		// Add State IDs joined to changed struct.
		statePartition.changed = append(statePartition.changed, stateID1)
		statePartition.changed = append(statePartition.changed, stateID2)
		return statePartition.linkBlocks(stateID1Root, stateID2Root)
	// If their root is equal, the states are already within
	// the same block so the root for state 1 is returned.
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
// levels to find the root element of the stateID
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

// Checks if states are within the same block
func (statePartition *StatePartition) WithinSameBlock(stateID1 int, stateID2 int) bool {
	return statePartition.Find(stateID1) == statePartition.Find(stateID2)
}

// Converts a DFA to a State Partition
func (dfa DFA) ToStatePartition() *StatePartition {
	// Return
	return NewStatePartition(len(dfa.States))
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

// Recursively merges states to merge state1 and state2, returns false
// if the merge results in an NFA, or true if merge was successful
func (statePartition *StatePartition) MergeStates(dfa DFA, state1 int, state2 int) bool{
	// return same state partition if states are already in the same block
	if statePartition.Find(state1) == statePartition.Find(state2){
		return true
	}

	// store block status, set as unknown by default
	var blockStatus StateStatus = UNKNOWN

	// store the block transitions and set to -1 by default
	transitions := make([]int, len(dfa.SymbolMap))
	for i := range transitions {
		transitions[i] = -1
	}

	// merge states within state partition
	mergedStateBlock := statePartition.union(state1, state2)

	// iterate over each state within the block containing the merged states
	for _, stateID := range statePartition.ReturnSet(mergedStateBlock){
		// store the current state status
		currentStateStatus := dfa.States[stateID].StateStatus
		// if the block status is unknown and the current state status
		// is not unknown, set the block status to the current state status
		if blockStatus == UNKNOWN && currentStateStatus != UNKNOWN{
			blockStatus = currentStateStatus
		// else check if the block status and the current state status are
		// non deterministic (one is accepting and one is rejecting)
		}else if (blockStatus == ACCEPTING && currentStateStatus == REJECTING) ||
			(blockStatus == REJECTING && currentStateStatus == ACCEPTING){
			// return false since merge is non-deterministic
			return false
		}
		// iterate over each symbol within DFA
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			// store resultant state from state transition of current state
			currentResultantStateID := dfa.States[stateID].Transitions[symbolID]
			// store resultant block from state transition of current block
			resultantBlockID := transitions[symbolID]
			// if both resultant block and current resultant state are bigger than
			// -1 (valid transition)
			if resultantBlockID > -1 && currentResultantStateID > -1{
				// get block which contains current current resultant state
				currentResultantBlockID := statePartition.Find(currentResultantStateID)
				// if the resultant block is not equal to the current resultant block,
				// it means that we have non-determinism
				if resultantBlockID != currentResultantBlockID {
					// not deterministic so merge, if states cannot be merged, return false
					if !statePartition.MergeStates(dfa, resultantBlockID, currentResultantBlockID) {
						return false
					}
				}
			}else{
				// if the current resultant state is initialized, set to block transition
				// for the current symbol being iterated
				if currentResultantStateID > -1{
					transitions[symbolID] = statePartition.Find(currentResultantStateID)
				}
			}
		}
	}

	// return true if this is reached (deterministic)
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