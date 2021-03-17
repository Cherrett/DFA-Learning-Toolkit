package DFA_Toolkit

// StatePartition struct which represents a State Partition.
type StatePartition struct {
	root []int					// Parent block of each state.
	size []int					// Size (score) of each block.
	link []int					// Index of next state within the block.
	changed []int				// Slice of changed states/blocks.

	isCopy bool					// Copied flag for reverting merges.
	//labelledStatesCount int		// Number of labelled states within partition.
	acceptingBlocksCount int	// Number of accepting blocks (states) within partition.
	rejectingBlocksCount int	// Number of rejecting blocks (states) within partition.
	blockStatus []StateStatus	// Status of each block.
}

// Returns an initialized State Partition.
func NewStatePartition(dfa DFA) *StatePartition {
	// Initialize new State Partition struct.
	statePartition := new(StatePartition)
	statePartition.isCopy = false
	// Initialize empty slices.
	statePartition.acceptingBlocksCount = 0
	statePartition.rejectingBlocksCount = 0
	statePartition.root = make([]int, len(dfa.States))
	statePartition.size = make([]int, len(dfa.States))
	statePartition.link = make([]int, len(dfa.States))
	statePartition.blockStatus = make([]StateStatus, len(dfa.States))
	// Set root and link as element, and size (score) as 1. Set block status
	// to state status and increment number of labelled states accordingly.
	for i := 0; i < len(dfa.States); i++ {
		statePartition.root[i] = i
		statePartition.size[i] = 1
		statePartition.link[i] = i
		statePartition.blockStatus[i] = dfa.States[i].StateStatus
		if statePartition.blockStatus[i] == ACCEPTING{
			statePartition.acceptingBlocksCount++
		}else if statePartition.blockStatus[i] == REJECTING{
			statePartition.rejectingBlocksCount++
		}
	}

	// Return initialized partition.
	return statePartition
}

// Connects two blocks by comparing their respective
// size (score) values to keep the tree flat.
func (statePartition *StatePartition) union(blockID1 int, blockID2 int){
	// Add State IDs joined to changed struct so merge can be undone.
	if statePartition.isCopy{
		statePartition.changed = append(statePartition.changed, blockID1)
		statePartition.changed = append(statePartition.changed, blockID2)
	}

	statePartition.link[blockID1], statePartition.link[blockID2] =
		statePartition.link[blockID2], statePartition.link[blockID1]

	block1Status := statePartition.blockStatus[blockID1]
	block2Status := statePartition.blockStatus[blockID2]

	if statePartition.size[blockID1] > statePartition.size[blockID2] {
		statePartition.root[blockID2] = blockID1
		statePartition.size[blockID1] ++

		if block1Status == UNKNOWN && block2Status != UNKNOWN{
			statePartition.blockStatus[blockID1] = block2Status
		}else if block1Status == ACCEPTING && block2Status == ACCEPTING{
			statePartition.acceptingBlocksCount--
		}else if block1Status == REJECTING && block2Status == REJECTING{
			statePartition.rejectingBlocksCount--
		}
	}else{
		statePartition.root[blockID1] = blockID2
		statePartition.size[blockID2] ++

		if block2Status == UNKNOWN && block1Status != UNKNOWN{
			statePartition.blockStatus[blockID2] = block1Status
		}else if block1Status == ACCEPTING && block2Status == ACCEPTING{
			statePartition.acceptingBlocksCount--
		}else if block1Status == REJECTING && block2Status == REJECTING{
			statePartition.rejectingBlocksCount--
		}
	}
}

// Find traverses each parent element while compressing the
// levels to find the root element of the stateID.
// If we attempt to access an element outside the array it returns -1.
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

// Checks if states are within the same block.
func (statePartition *StatePartition) WithinSameBlock(stateID1 int, stateID2 int) bool {
	return statePartition.Find(stateID1) == statePartition.Find(stateID2)
}

// Converts a DFA to a State Partition.
func (dfa DFA) ToStatePartition() *StatePartition {
	// Return
	return NewStatePartition(dfa)
}

// Converts a State Partition to a DFA.
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

// Recursively merges states to merge state1 and state2. Returns false if merge
// results in a non-deterministic automaton. Returns true if merge was successful.
func (statePartition *StatePartition) MergeStates(dfa DFA, state1 int, state2 int) bool{
	// If parent blocks (root) are the same as state ID, skip finding the root.
	// Else, find the parent block (root) using Find.
	if statePartition.root[state1] != state1{
		state1 = statePartition.Find(state1)
	}
	if statePartition.root[state2] != state2{
		state2 = statePartition.Find(state2)
	}

	// Return true if states are already in the same block
	// since merge is not required.
	if state1 == state2 {
		return true
	}

	// Get status of each block.
	state1Status := statePartition.blockStatus[state1]
	state2Status := statePartition.blockStatus[state2]
	// If statuses are contradicting, return false since this results
	// in a non-deterministic automaton so merge cannot be done.
	if (state1Status == ACCEPTING && state2Status == REJECTING) || (state1Status == REJECTING && state2Status == ACCEPTING){
		return false
	}

	// Initialize a slice which will be used to evaluate the validity of
	// the transitions for the merged blocks.
	transitions := make([]int, len(dfa.SymbolMap))
	// Set all values to -1 by default.
	for i := range transitions {
		transitions[i] = -1
	}

	// Get the states within each block.
	block1Set := statePartition.ReturnSet(state1)
	block2Set := statePartition.ReturnSet(state2)

	// Merge states within state partition.
	statePartition.union(state1, state2)

	// Iterate over each symbol within DFA.
	for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
		// Iterate over each state within first block.
		for _, stateID := range block1Set{
			// Store resultant state from state transition of current state.
			currentResultantStateID := dfa.States[stateID].Transitions[symbolID]

			// If resultant state ID is bigger than -1 (valid transition), get
			// the block containing state and store in transitions. The loop is
			// then broken since the transition for the current symbol was found.
			if currentResultantStateID > -1{
				// Set resultant state to state transition for current symbol.
				transitions[symbolID] = currentResultantStateID
				// Break loop since the transition for the current symbol is found.
				break
			}
		}
	}

	// Iterate over each symbol within DFA.
	for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
		// If no transitions exist from block for the given symbol, skip
		// current loop iteration.
		if transitions[symbolID] == -1{
			continue
		}
		// Iterate over each state within second block.
		for _, stateID := range block2Set{
			// Store resultant state from state transition of current state.
			currentResultantStateID := dfa.States[stateID].Transitions[symbolID]
			// If resultant state ID is bigger than -1 (valid transition), get the
			// block containing state and compare it to the transitions found above.
			// If they are not equal, merge blocks to eliminate non-determinism.
			if currentResultantStateID > -1{
				// If resultant block is not equal to the block containing the state within transitions
				// found above, merge the two states to eliminate non-determinism.
				if statePartition.Find(currentResultantStateID) != statePartition.Find(transitions[symbolID]){
					// Not deterministic so merge, if states cannot be merged, return false.
					if !statePartition.MergeStates(dfa, transitions[symbolID], currentResultantStateID) {
						return false
					}
				}
			}
		}
	}

	// Return true if this is reached (deterministic).
	return true
}

// Copies the state partition.
func (statePartition StatePartition) Copy() StatePartition{
	// Initialize new State Partition struct.
	copiedStatePartition := StatePartition{
		root:                make([]int, len(statePartition.root)),
		size:                make([]int, len(statePartition.size)),
		link:                make([]int, len(statePartition.link)),
		changed:             []int{},
		isCopy:              true,
		acceptingBlocksCount: statePartition.acceptingBlocksCount,
		rejectingBlocksCount: statePartition.rejectingBlocksCount,
		blockStatus:         make([]StateStatus, len(statePartition.blockStatus)),
	}

	// Copy root, size, link and blockStatus slices.
	copy(copiedStatePartition.root, statePartition.root)
	copy(copiedStatePartition.size, statePartition.size)
	copy(copiedStatePartition.link, statePartition.link)
	copy(copiedStatePartition.blockStatus, statePartition.blockStatus)

	// Return copied state partition.
	return copiedStatePartition
}

// Reverts any changes made within state partition given the original state partition.
func (statePartition *StatePartition) RollbackChanges(originalStatePartition *StatePartition){
	// If the state partition is a copy, copy values of changed blocks from original
	// state partition. Else, do nothing.
	if statePartition.isCopy{
		// Set accepting and rejecting blocks count values to the original values.
		statePartition.acceptingBlocksCount = originalStatePartition.acceptingBlocksCount
		statePartition.rejectingBlocksCount = originalStatePartition.rejectingBlocksCount
		// Iterate over each altered block (state).
		for _, stateID := range statePartition.changed{
			// Update root, size, link and blockStatus values.
			statePartition.root[stateID] = originalStatePartition.root[stateID]
			statePartition.size[stateID] = originalStatePartition.size[stateID]
			statePartition.link[stateID] = originalStatePartition.link[stateID]
			statePartition.blockStatus[stateID] = originalStatePartition.blockStatus[stateID]
		}
		// Empty the changed blocks slice.
		statePartition.changed = []int{}
	}
}

// Returns the number of labelled blocks (states) within state partition.
func (statePartition StatePartition) NumberOfLabelledBlocks() int{
	// Return the sum of the accepting and rejecting blocks count.
	return statePartition.acceptingBlocksCount + statePartition.rejectingBlocksCount
}