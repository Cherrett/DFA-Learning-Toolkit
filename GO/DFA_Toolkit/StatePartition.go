package dfatoolkit

// Block struct which represents a block within a partition.
type Block struct {
	Root    int        // Parent block of state.
	Size    int        // Size (score) of block.
	Link    int        // Index of next state within the block.
	Label   StateLabel // Label of block.
	Changed bool       // Whether block has been changed.
}

// StatePartition struct which represents a State Partition.
type StatePartition struct {
	Blocks               []Block // Slice of blocks.
	BlocksCount          int     // Number of blocks within partition.
	AcceptingBlocksCount int     // Number of accepting blocks within partition.
	RejectingBlocksCount int     // Number of rejecting blocks within partition.

	IsCopy             bool  // Whether state partition is a copy (for reverting merges).
	ChangedBlocks      []int // Slice of changed blocks.
	ChangedBlocksCount int   // Number of changed blocks within partition.
}

// NewStatePartition returns an initialized State Partition.
func NewStatePartition(dfa DFA) StatePartition {
	// Initialize new State Partition struct and initialize
	// empty slices of nodes and changed blocks.
	statePartition := StatePartition{
		Blocks:               make([]Block, len(dfa.States)),
		IsCopy:               false,
		BlocksCount:          len(dfa.States),
		AcceptingBlocksCount: 0,
		RejectingBlocksCount: 0,
		ChangedBlocksCount:   0,
	}

	// Set root and link as element, and size (score) as 1. Set block label
	// to state label and increment number of labelled states accordingly.
	for i := 0; i < len(dfa.States); i++ {
		statePartition.Blocks[i].Root = i
		statePartition.Blocks[i].Size = 1
		statePartition.Blocks[i].Link = i
		statePartition.Blocks[i].Label = dfa.States[i].Label
		statePartition.Blocks[i].Changed = false
		if statePartition.Blocks[i].Label == ACCEPTING {
			statePartition.AcceptingBlocksCount++
		} else if statePartition.Blocks[i].Label == REJECTING {
			statePartition.RejectingBlocksCount++
		}
	}

	// Return initialized partition.
	return statePartition
}

// ChangedBlock updates the required fields to mark block as changed.
func (statePartition *StatePartition) ChangedBlock(blockID int) {
	// If block is not already changed.
	if !statePartition.Blocks[blockID].Changed {
		// Update changed slice to include changed block ID.
		statePartition.ChangedBlocks[statePartition.ChangedBlocksCount] = blockID
		// Increment the changed blocks counter.
		statePartition.ChangedBlocksCount++
		// Set changed flag within block to true.
		statePartition.Blocks[blockID].Changed = true
	}
}

// Union connects two blocks by comparing their respective
// size (score) values to keep the tree flat.
func (statePartition *StatePartition) Union(blockID1 int, blockID2 int) int{
	// If state partition is a copy, call ChangedBlock for
	// both blocks so merge can be undone if necessary.
	if statePartition.IsCopy {
		statePartition.ChangedBlock(blockID1)
		statePartition.ChangedBlock(blockID2)
	}

	// Decrement blocks count.
	statePartition.BlocksCount--

	// If size of parent node is smaller than size of child node, switch
	// parent and child nodes.
	if statePartition.Blocks[blockID1].Size < statePartition.Blocks[blockID2].Size{
		blockID1, blockID2 = blockID2, blockID1
	}

	// Set block 1 to parent and block 2 to child.
	parent, child := &statePartition.Blocks[blockID1], &statePartition.Blocks[blockID2]

	// Link nodes by assigning the link of parent to link of child and vice versa.
	parent.Link, child.Link = child.Link, parent.Link

	// Set root of child node to parent node.
	child.Root = parent.Root
	// Increment size (score) of parent node by size of child node.
	parent.Size += child.Size

	// If label of parent is unknown and label of child is
	// not unknown, set label of parent to label of child.
	if parent.Label == UNKNOWN && child.Label != UNKNOWN {
		parent.Label = child.Label
	} else if parent.Label == ACCEPTING && child.Label == ACCEPTING {
		// Else, if both blocks are accepting, decrement accepting blocks count.
		statePartition.AcceptingBlocksCount--
	} else if parent.Label == REJECTING && child.Label == REJECTING {
		// Else, if both blocks are rejecting, decrement rejecting blocks count.
		statePartition.RejectingBlocksCount--
	}

	return blockID1
}

// Find traverses each root element while compressing the
// levels to find the root element of the stateID.
func (statePartition *StatePartition) Find(stateID int) int {
	// Panic if out of range.
	//if stateID > len(statePartition.Blocks)-1 {
	//	panic("StateID out of range.")
	//}

	// Traverse each root block until state is reached.
	for statePartition.Blocks[stateID].Root != stateID {
		// Compress root.
		statePartition.Blocks[stateID].Root = statePartition.Blocks[statePartition.Blocks[stateID].Root].Root
		stateID = statePartition.Blocks[stateID].Root
	}

	return stateID
}

// ReturnSet returns the state IDs within given block.
func (statePartition StatePartition) ReturnSet(blockID int) []int {
	// Slice of state IDs.
	blockElements := []int{blockID}
	// Set root to block ID.
	root := blockID

	// Iterate until link of current block ID is
	// not equal to the root block.
	for statePartition.Blocks[blockID].Link != root {
		// Set block ID to link of current block.
		blockID = statePartition.Blocks[blockID].Link
		// Add block ID to block elements slice.
		blockElements = append(blockElements, blockID)

		// Panic if length of elements slice is bigger than number of blocks.
		//if len(blockElements) > len(statePartition.Blocks) {
		//	panic("Error in state linking.")
		//}
	}

	// Return state IDs within block.
	return blockElements
}

// WithinSameBlock checks whether two states are within the same block.
func (statePartition *StatePartition) WithinSameBlock(stateID1 int, stateID2 int) bool {
	return statePartition.Find(stateID1) == statePartition.Find(stateID2)
}

// ToStatePartition converts a DFA to a State Partition.
func (dfa DFA) ToStatePartition() StatePartition {
	// Return NewStatePartition.
	return NewStatePartition(dfa)
}

// ToDFA converts a State Partition to a DFA. Returns true and
// the corresponding DFA if state partition is valid. Else,
// false and an empty DFA are returned.
func (statePartition StatePartition) ToDFA(dfa DFA) (bool, DFA) {
	newMappings := map[int]int{}

	resultantDFA := DFA{
		States:                nil,
		StartingStateID:       -1,
		Alphabet:              dfa.Alphabet,
		depth:                 -1,
		computedDepthAndOrder: false,
	}

	for _, stateID := range statePartition.RootBlocks() {
		newMappings[stateID] = resultantDFA.AddState(statePartition.Blocks[stateID].Label)
	}

	// update starting state via mappings
	resultantDFA.StartingStateID = newMappings[statePartition.Find(dfa.StartingStateID)]

	// update new transitions via mappings
	for stateID := range dfa.States {
		for alphabetID := range dfa.Alphabet {
			oldResultantStateID := dfa.States[stateID].Transitions[alphabetID]
			if oldResultantStateID > -1 {
				newStateID := newMappings[statePartition.Find(stateID)]
				resultantStateID := newMappings[statePartition.Find(oldResultantStateID)]
				if resultantDFA.States[newStateID].Transitions[alphabetID] > -1 &&
					resultantDFA.States[newStateID].Transitions[alphabetID] != resultantStateID {
					// not deterministic
					return false, DFA{}
				}
				resultantDFA.States[newStateID].Transitions[alphabetID] = resultantStateID
			}
		}
	}
	return true, resultantDFA
}

// MergeStates recursively merges states to merge state1 and state2. Returns false if merge
// results in a non-deterministic automaton. Returns true if merge was successful.
func (statePartition *StatePartition) MergeStates(dfa DFA, state1 int, state2 int) bool {
	// If parent blocks (root) are the same as state ID, skip finding the root.
	// Else, find the parent block (root) using Find.
	if statePartition.Blocks[state1].Root != state1 {
		state1 = statePartition.Find(state1)
	}
	if statePartition.Blocks[state2].Root != state2 {
		state2 = statePartition.Find(state2)
	}

	// Return true if states are already in the same block
	// since merge is not required.
	if state1 == state2 {
		return true
	}

	// Get label of each block.
	state1Label := statePartition.Blocks[state1].Label
	state2Label := statePartition.Blocks[state2].Label
	// If labels are contradicting, return false since this results
	// in a non-deterministic automaton so merge cannot be done.
	if (state1Label == ACCEPTING && state2Label == REJECTING) || (state1Label == REJECTING && state2Label == ACCEPTING) {
		return false
	}

	// Merge states within state partition.
	root := statePartition.Union(state1, state2)

	// Iterate over each symbol within DFA.
	for alphabetID := range dfa.Alphabet {
		// Set block ID to root.
		blockID := root
		// Set resultant state to -1.
		resultantState := -1

		// Iterate until link of current block ID is
		// not equal to the root block.
		for{
			// Store current resultant state from state transition of current state and check if transition is valid.
			if currentResultantState := dfa.States[blockID].Transitions[alphabetID]; currentResultantState >= 0{
				// If transition is valid and resultant state is already found, merge
				// current resultant state and resultant state. Else, if resultant
				// state is not found, set it to current resultant state.
				if resultantState >= 0{
					// Return false if merge is nondeterministic.
					if !statePartition.MergeStates(dfa, currentResultantState, resultantState) {
						return false
					}
				}else{
					resultantState = currentResultantState
				}
			}

			// Set block ID to link of current block ID and break if
			// new block ID is equal to the root block.
			if blockID = statePartition.Blocks[blockID].Link; blockID == root{
				break
			}
		}
	}

	// Return true if this is reached (deterministic).
	return true
}

// Copy copies the state partition.
func (statePartition StatePartition) Copy() StatePartition {
	// Initialize new State Partition struct.
	copiedStatePartition := StatePartition{
		Blocks:               make([]Block, len(statePartition.Blocks)),
		ChangedBlocks:        make([]int, len(statePartition.Blocks)),
		IsCopy:               true,
		BlocksCount:          statePartition.BlocksCount,
		AcceptingBlocksCount: statePartition.AcceptingBlocksCount,
		RejectingBlocksCount: statePartition.RejectingBlocksCount,
		ChangedBlocksCount:   0,
	}

	// Copy blocks slice.
	copy(copiedStatePartition.Blocks, statePartition.Blocks)

	// Return copied state partition.
	return copiedStatePartition
}

// RollbackChanges reverts any changes made within state partition given the original state partition.
func (statePartition *StatePartition) RollbackChanges(originalStatePartition StatePartition) {
	// If the state partition is a copy, copy values of changed blocks from original
	// state partition. Else, do nothing.
	if statePartition.IsCopy {
		// Set blocks count values to the original values.
		statePartition.BlocksCount = originalStatePartition.BlocksCount
		statePartition.AcceptingBlocksCount = originalStatePartition.AcceptingBlocksCount
		statePartition.RejectingBlocksCount = originalStatePartition.RejectingBlocksCount

		// Iterate over each altered block (state).
		for i := 0; i < statePartition.ChangedBlocksCount; i++ {
			// Update root, size, link and blockLabel values.
			stateID := statePartition.ChangedBlocks[i]
			statePartition.Blocks[stateID] = originalStatePartition.Blocks[stateID]
		}
		// Empty the changed blocks slice.
		statePartition.ChangedBlocksCount = 0
	}
}

// NumberOfLabelledBlocks returns the number of labelled blocks (states) within state partition.
func (statePartition StatePartition) NumberOfLabelledBlocks() int {
	// Return the sum of the accepting and rejecting blocks count.
	return statePartition.AcceptingBlocksCount + statePartition.RejectingBlocksCount
}

// RootBlocks returns
func (statePartition StatePartition) RootBlocks() []int {
	rootBlocks := make([]int, statePartition.BlocksCount)
	index := 0

	for blockID := range statePartition.Blocks {
		if statePartition.Blocks[blockID].Root == blockID {
			rootBlocks[index] = blockID
			index++
			if index == statePartition.BlocksCount {
				break
			}
		}
	}

	return rootBlocks
}
