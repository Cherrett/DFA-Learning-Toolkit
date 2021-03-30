package DFA_Toolkit

type Node struct{
	root int					// Parent block of state.
	size int					// Size (score) of block.
	link int					// Index of next state within the block.
	status StateStatus			// Status of block.
	changed bool
}

// StatePartition struct which represents a State Partition.
type StatePartitionV2 struct {
	nodes []Node				// Slice of nodes.
	changed []int				// Slice of changed states/blocks.

	isCopy bool					// Copied flag for reverting merges.
	//labelledStatesCount int		// Number of labelled states within partition.
	acceptingBlocksCount int	// Number of accepting blocks (states) within partition.
	rejectingBlocksCount int	// Number of rejecting blocks (states) within partition.
	changedBlocks int
}

// Returns an initialized State Partition.
func NewStatePartitionV2(dfa DFA) StatePartitionV2 {
	// Initialize new State Partition struct and
	// initialize empty slices.
	statePartition := StatePartitionV2{
		nodes:                make([]Node, len(dfa.States)),
		changed:              make([]int, len(dfa.States)),
		isCopy:               false,
		acceptingBlocksCount: 0,
		rejectingBlocksCount: 0,
		changedBlocks:		  0,
	}

	// Set root and link as element, and size (score) as 1. Set block status
	// to state status and increment number of labelled states accordingly.
	for i := 0; i < len(dfa.States); i++ {
		statePartition.nodes[i].root = i
		statePartition.nodes[i].size = 1
		statePartition.nodes[i].link = i
		statePartition.nodes[i].status = dfa.States[i].StateStatus
		statePartition.nodes[i].changed = false
		statePartition.changed[i] = -1
		if statePartition.nodes[i].status == ACCEPTING{
			statePartition.acceptingBlocksCount++
		}else if statePartition.nodes[i].status == REJECTING{
			statePartition.rejectingBlocksCount++
		}
	}

	// Return initialized partition.
	return statePartition
}

func (statePartition *StatePartitionV2) modifiedBlock(blockID int){
	if !statePartition.nodes[blockID].changed{
		statePartition.changed[statePartition.changedBlocks] = blockID
		statePartition.changedBlocks++
		statePartition.nodes[blockID].changed = true

		if statePartition.changedBlocks == len(statePartition.changed){
			panic("Test")
		}
	}
}

// Connects two blocks by comparing their respective
// size (score) values to keep the tree flat.
func (statePartition *StatePartitionV2) union(blockID1 int, blockID2 int){
	// Add State IDs joined to changed struct so merge can be undone.
	if statePartition.isCopy{
		statePartition.modifiedBlock(blockID1)
		statePartition.modifiedBlock(blockID2)
	}

	statePartition.nodes[blockID1].link, statePartition.nodes[blockID2].link =
		statePartition.nodes[blockID2].link, statePartition.nodes[blockID1].link

	block1Status := statePartition.nodes[blockID1].status
	block2Status := statePartition.nodes[blockID2].status

	if statePartition.nodes[blockID1].size > statePartition.nodes[blockID2].size {
		statePartition.nodes[blockID2].root = blockID1
		statePartition.nodes[blockID1].size ++

		if block1Status == UNKNOWN && block2Status != UNKNOWN{
			statePartition.nodes[blockID1].status = block2Status
		}else if block1Status == ACCEPTING && block2Status == ACCEPTING{
			statePartition.acceptingBlocksCount--
		}else if block1Status == REJECTING && block2Status == REJECTING{
			statePartition.rejectingBlocksCount--
		}
	}else{
		statePartition.nodes[blockID1].root = blockID2
		statePartition.nodes[blockID2].size ++

		if block2Status == UNKNOWN && block1Status != UNKNOWN{
			statePartition.nodes[blockID2].status = block1Status
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
func (statePartition *StatePartitionV2) Find(stateID int) int {
	if stateID > len(statePartition.nodes)-1 {
		panic("StateID out of range.")
	}

	for statePartition.nodes[stateID].root != stateID {
		statePartition.nodes[stateID].root = statePartition.nodes[statePartition.nodes[stateID].root].root
		stateID = statePartition.nodes[stateID].root
	}

	return stateID
}

func (statePartition StatePartitionV2) ReturnSet(stateID int) []int{
	blockElements := []int{stateID}
	root := stateID
	for statePartition.nodes[stateID].link != root{
		stateID = statePartition.nodes[stateID].link
		blockElements = append(blockElements, stateID)
		if len(blockElements) > len(statePartition.nodes){
			panic("Error in state linking.")
		}
	}
	return blockElements
}

// Checks if states are within the same block.
func (statePartition *StatePartitionV2) WithinSameBlock(stateID1 int, stateID2 int) bool {
	return statePartition.Find(stateID1) == statePartition.Find(stateID2)
}

// Converts a DFA to a State Partition.
func (dfa DFA) ToStatePartitionV2() StatePartitionV2 {
	// Return NewStatePartition.
	return NewStatePartitionV2(dfa)
}

// Converts a State Partition to a DFA.
func (statePartition StatePartitionV2) ToDFA(dfa DFA) (bool, DFA){
	newMappings := map[int]int{}

	resultantDFA := DFA{
		States:                nil,
		StartingStateID:       -1,
		SymbolMap:             dfa.SymbolMap,
		depth:                 -1,
		computedDepthAndOrder: false,
	}

	for stateID := range dfa.States {
		currentBlockID := statePartition.Find(stateID)
		if _, ok := newMappings[currentBlockID]; !ok {
			newMappings[currentBlockID] = resultantDFA.AddState(statePartition.nodes[currentBlockID].status)
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
func (statePartition *StatePartitionV2) MergeStates(dfa DFA, state1 int, state2 int) bool{
	// If parent blocks (root) are the same as state ID, skip finding the root.
	// Else, find the parent block (root) using Find.
	if statePartition.nodes[state1].root != state1{
		state1 = statePartition.Find(state1)
	}
	if statePartition.nodes[state2].root != state2{
		state2 = statePartition.Find(state2)
	}

	// Return true if states are already in the same block
	// since merge is not required.
	if state1 == state2 {
		return true
	}

	// Get status of each block.
	state1Status := statePartition.nodes[state1].status
	state2Status := statePartition.nodes[state2].status
	// If statuses are contradicting, return false since this results
	// in a non-deterministic automaton so merge cannot be done.
	if (state1Status == ACCEPTING && state2Status == REJECTING) || (state1Status == REJECTING && state2Status == ACCEPTING){
		return false
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
			// the block containing state and store in transitionResult. The
			// states within the second block are then iterated and checked
			// for non-deterministic properties.
			if currentResultantStateID > -1{
				// Set resultant state to state transition for current symbol.
				transitionResult := currentResultantStateID

				// Iterate over each state within second block.
				for _, stateID2 := range block2Set{
					// Store resultant state from state transition of current state.
					currentResultantStateID = dfa.States[stateID2].Transitions[symbolID]
					// If resultant state ID is bigger than -1 (valid transition), get the
					// block containing state and compare it to the transition found above.
					// If they are not equal, merge blocks to eliminate non-determinism.
					if currentResultantStateID > -1{
						// If resultant block is not equal to the block containing the state within transition
						// found above, merge the two states to eliminate non-determinism.
						if statePartition.Find(currentResultantStateID) != statePartition.Find(transitionResult){
							// Not deterministic so merge, if states cannot be merged, return false.
							if !statePartition.MergeStates(dfa, transitionResult, currentResultantStateID) {
								return false
							}
						}
					}
				}

				// The loop is broken since the transition for the current symbol was found.
				break
			}
		}
	}

	// Return true if this is reached (deterministic).
	return true
}

// Copies the state partition.
func (statePartition StatePartitionV2) Copy() StatePartitionV2{
	// Initialize new State Partition struct.
	copiedStatePartition := StatePartitionV2{
		nodes:                make([]Node, len(statePartition.nodes)),
		changed:              make([]int, len(statePartition.nodes)),
		isCopy:               true,
		acceptingBlocksCount: statePartition.acceptingBlocksCount,
		rejectingBlocksCount: statePartition.rejectingBlocksCount,
		changedBlocks:        0,
	}

	// Copy root, size, link and blockStatus slices.
	copy(copiedStatePartition.nodes, statePartition.nodes)
	copy(copiedStatePartition.changed, statePartition.changed)

	// Return copied state partition.
	return copiedStatePartition
}

// Reverts any changes made within state partition given the original state partition.
func (statePartition *StatePartitionV2) RollbackChanges(originalStatePartition StatePartitionV2){
	// If the state partition is a copy, copy values of changed blocks from original
	// state partition. Else, do nothing.
	if statePartition.isCopy{
		// Set accepting and rejecting blocks count values to the original values.
		statePartition.acceptingBlocksCount = originalStatePartition.acceptingBlocksCount
		statePartition.rejectingBlocksCount = originalStatePartition.rejectingBlocksCount

		// Iterate over each altered block (state).
		for i := 0; i < statePartition.changedBlocks; i++ {
			// Update root, size, link and blockStatus values.
			stateID := statePartition.changed[i]
			statePartition.nodes[stateID] = originalStatePartition.nodes[stateID]
			statePartition.changed[i] = -1
		}
		// Empty the changed blocks slice.
		statePartition.changedBlocks = 0
	}
}

// Returns the number of labelled blocks (states) within state partition.
func (statePartition StatePartitionV2) NumberOfLabelledBlocks() int{
	// Return the sum of the accepting and rejecting blocks count.
	return statePartition.acceptingBlocksCount + statePartition.rejectingBlocksCount
}