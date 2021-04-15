package dfatoolkit

// Block struct which represents a block within a partition.
type Block struct {
	Root    int        // Parent block of state.
	Size    int        // Size (score) of block.
	Link    int        // Index of next state within the block.
	Label   StateLabel // Label of block.
	Changed bool       // Whether block has been changed.
	Transitions []int  // Transition Table where each element corresponds to a transition for each symbol.
}

// StatePartition struct which represents a State Partition.
type StatePartition struct {
	Blocks               []Block // Slice of blocks.
	BlocksCount          int     // Number of blocks within partition.
	AcceptingBlocksCount int     // Number of accepting blocks within partition.
	RejectingBlocksCount int     // Number of rejecting blocks within partition.
	AlphabetSize int			 // Size of alphabet.
	StartingStateID int			 // The ID of the starting state.

	IsCopy             bool  // Whether state partition is a copy (for reverting merges).
	ChangedBlocks      []int // Slice of changed blocks.
	ChangedBlocksCount int   // Number of changed blocks within partition.
}

// NewStatePartition returns an initialized State Partition.
func NewStatePartition(referenceDFA DFA) StatePartition {
	// Initialize new State Partition struct with required values.
	statePartition := StatePartition{
		Blocks:          make([]Block, len(referenceDFA.States)),
		BlocksCount:     len(referenceDFA.States),
		AlphabetSize:    len(referenceDFA.Alphabet),
		StartingStateID: referenceDFA.StartingStateID,
	}

	// Set root and link as element, and size (score) as 1. Set block label
	// to state label and increment number of labelled states accordingly.
	// Transitions are then created using corresponding reference DFA.
	for i := 0; i < len(referenceDFA.States); i++ {
		statePartition.Blocks[i].Root = i
		statePartition.Blocks[i].Size = 1
		statePartition.Blocks[i].Link = i
		statePartition.Blocks[i].Label = referenceDFA.States[i].Label
		statePartition.Blocks[i].Changed = false
		if statePartition.Blocks[i].Label == ACCEPTING {
			statePartition.AcceptingBlocksCount++
		} else if statePartition.Blocks[i].Label == REJECTING {
			statePartition.RejectingBlocksCount++
		}

		// Initialize transitions.
		statePartition.Blocks[i].Transitions = make([]int, statePartition.AlphabetSize)
		for symbol := range referenceDFA.Alphabet{
			statePartition.Blocks[i].Transitions[symbol] = referenceDFA.States[i].Transitions[symbol]
		}
	}

	// Return initialized partition.
	return statePartition
}

// ChangedBlock updates the required fields to mark block as changed.
func (statePartition *StatePartition) ChangedBlock(blockID int) {
	// Check that state partition is a copy.
	if statePartition.IsCopy{
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
}

// Union connects two blocks by comparing their respective
// size (score) values to keep the tree flat.
func (statePartition *StatePartition) Union(blockID1 int, blockID2 int){
	// Mark both blocks as changed so merge can be undone.
	statePartition.ChangedBlock(blockID1)
	statePartition.ChangedBlock(blockID2)

	// Decrement blocks count.
	statePartition.BlocksCount--

	// If size of parent node is smaller than size of child node, switch
	// parent and child nodes.
	if statePartition.Blocks[blockID1].Size < statePartition.Blocks[blockID2].Size{
		blockID1, blockID2 = blockID2, blockID1
	}

	// Set pointer of block 1 to parent and pointer of block 2 to child.
	parent, child := &statePartition.Blocks[blockID1], &statePartition.Blocks[blockID2]

	// Link nodes by assigning the link of parent to link of child and vice versa.
	parent.Link, child.Link = child.Link, parent.Link

	// Set root of child node to parent node.
	child.Root = blockID1
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

	// Update transitions.
	for i := 0; i < statePartition.AlphabetSize; i++ {
		// If transition of parent does not exist,
		// set transition of child and set child
		// transition to -1 (remove transition).
		if parent.Transitions[i] == -1 {
			parent.Transitions[i] = child.Transitions[i]
			child.Transitions[i] = -1
		}
	}
}

// Find traverses each root element while compressing the
// levels to find the root element of the stateID.
func (statePartition *StatePartition) Find(stateID int) int {
	// Traverse each root block until state is reached.
	for statePartition.Blocks[stateID].Root != stateID {
		// Compress if necessary.
		if statePartition.Blocks[stateID].Root != statePartition.Blocks[statePartition.Blocks[stateID].Root].Root{
			// If compression is required, mark state as
			// changed and set root to root of parent.
			statePartition.ChangedBlock(stateID)
			statePartition.Blocks[stateID].Root = statePartition.Blocks[statePartition.Blocks[stateID].Root].Root
		}
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
func (statePartition StatePartition) ToDFA() DFA {
	// Map to store corresponding new state for
	// each root block within state partition.
	blockToStateMap := map[int]int{}

	// Initialize resultant DFA to be returned by function.
	resultantDFA := NewDFA()

	// Get root blocks within state partition.
	rootBlocks := statePartition.RootBlocks()

	// Create a new state within DFA for each root block and
	// set state label to block label.
	for _, stateID := range rootBlocks {
		blockToStateMap[stateID] = resultantDFA.AddState(statePartition.Blocks[stateID].Label)
	}

	// Create alphabet within DFA.
	for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
		resultantDFA.AddSymbol()
	}

	// Update transitions using transitions within blocks and block to state map.
	for _, stateID := range rootBlocks {
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			if resultantState := statePartition.Blocks[stateID].Transitions[symbol]; resultantState > -1{
				resultantDFA.States[blockToStateMap[stateID]].Transitions[symbol] = blockToStateMap[statePartition.Find(resultantState)]
			}
		}
	}

	// Set starting state using block to state map.
	resultantDFA.StartingStateID = blockToStateMap[statePartition.Find(statePartition.StartingStateID)]

	// Return populated resultant DFA.
	return resultantDFA
}

// MergeStates recursively merges states to merge state1 and state2. Returns false if merge
// results in a non-deterministic automaton. Returns true if merge was successful.
func (statePartition *StatePartition) MergeStates(state1 int, state2 int) bool {
	// If parent blocks (root) are the same as state ID, skip finding the root.
	// Else, find the parent block (root) using Find function.
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

	// Get pointer of both blocks.
	block1, block2 := &statePartition.Blocks[state1],  &statePartition.Blocks[state2]

	// If labels are contradicting, return false since this results
	// in a non-deterministic automaton so merge cannot be done.
	if (block1.Label == ACCEPTING && block2.Label == REJECTING) || (block1.Label == REJECTING && block2.Label == ACCEPTING) {
		return false
	}

	// Merge states within state partition.
	statePartition.Union(state1, state2)

	// Iterate over alphabet.
	for i := 0; i < statePartition.AlphabetSize; i++ {
		// If either block1 or block2 do not have a transition, continue
		// since no merge is required.
		if block1.Transitions[i] == -1 || block2.Transitions[i] == -1{
			continue
		}
		// Else, merge resultant blocks.
		if !statePartition.MergeStates(block1.Transitions[i], block2.Transitions[i]) {
			return false
		}
	}

	// Return true if this is reached (deterministic).
	return true
}

// Copy copies the state partition.
func (statePartition StatePartition) Copy() StatePartition {
	// Initialize new StatePartition struct using state partition.
	copiedStatePartition := StatePartition{
		Blocks:               make([]Block, len(statePartition.Blocks)),
		BlocksCount:          statePartition.BlocksCount,
		AcceptingBlocksCount: statePartition.AcceptingBlocksCount,
		RejectingBlocksCount: statePartition.RejectingBlocksCount,
		AlphabetSize:         statePartition.AlphabetSize,
		StartingStateID:      statePartition.StartingStateID,
		IsCopy:               true,
		ChangedBlocks:        make([]int, len(statePartition.Blocks)),
		ChangedBlocksCount:   0,
	}

	// Iterate over each block within blocks slice.
	for blockID := range statePartition.Blocks{
		// Get block pointer from copied state partition.
		block := &copiedStatePartition.Blocks[blockID]
		// Get block pointer from original state partition.
		originalBlock := &statePartition.Blocks[blockID]

		// Update root, size, link, and label.
		block.Root = originalBlock.Root
		block.Size = originalBlock.Size
		block.Link = originalBlock.Link
		block.Label = originalBlock.Label

		// Initialize and copy transitions.
		block.Transitions = make([]int, statePartition.AlphabetSize)
		copy(block.Transitions, originalBlock.Transitions)
	}

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

		// Iterate over each altered block.
		for _, blockID := range statePartition.ChangedBlocks[:statePartition.ChangedBlocksCount] {
			// Get block pointer from copied state partition.
			block := &statePartition.Blocks[blockID]
			// Get block pointer from original state partition.
			originalBlock := &originalStatePartition.Blocks[blockID]

			// Update root, size, link, and label.
			block.Root = originalBlock.Root
			block.Size = originalBlock.Size
			block.Link = originalBlock.Link
			block.Label = originalBlock.Label
			block.Changed = false

			// Copy transitions.
			copy(block.Transitions, originalBlock.Transitions)
		}

		// Set the changed blocks count to 0.
		statePartition.ChangedBlocksCount = 0
	}
}

// NumberOfLabelledBlocks returns the number of labelled blocks (states) within state partition.
func (statePartition StatePartition) NumberOfLabelledBlocks() int {
	// Return the sum of the accepting and rejecting blocks count.
	return statePartition.AcceptingBlocksCount + statePartition.RejectingBlocksCount
}

// RootBlocks returns the IDs of root blocks as a slice of integers.
func (statePartition StatePartition) RootBlocks() []int {
	// Initialize slice using blocks count value.
	rootBlocks := make([]int, statePartition.BlocksCount)
	// Index (count) of root blocks.
	index := 0

	// Iterate over each block within partition.
	for blockID := range statePartition.Blocks {
		// Check if root of current block is equal to the block ID
		if statePartition.Blocks[blockID].Root == blockID {
			// Add to rootBlocks slice using index.
			rootBlocks[index] = blockID
			// Increment index.
			index++
			// If index is equal to blocks count,
			// break since all root blocks have
			// been found.
			if index == statePartition.BlocksCount {
				break
			}
		}
	}

	// Return populated slice.
	return rootBlocks
}

// OrderedBlocks returns the IDs of root blocks in order as a slice of integers.
func (statePartition StatePartition) OrderedBlocks() []int{
	orderComputed := make([]bool, len(statePartition.Blocks))
	orderedBlocks := make([]int, statePartition.BlocksCount)
	index := 0

	// Create a FIFO queue with starting state.
	queue := []int{statePartition.Find(statePartition.Find(statePartition.StartingStateID))}

	// Loop until queue is empty.
	for len(queue) > 0 {
		// Remove and store first state in queue.
		blockID := queue[0]
		queue = queue[1:]

		if orderComputed[blockID]{
			continue
		}

		// Set the order of the current state.
		orderedBlocks[index] = blockID
		orderComputed[blockID] = true
		// Increment current state order.
		index++

		// Iterate over each symbol (alphabet) within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			// If transition from current state using current symbol is valid and is not a loop to the current state.
			if childStateID := statePartition.Blocks[blockID].Transitions[symbol]; childStateID != -1 {
				// If depth for child state has been computed, skip state.
				if childBlockID := statePartition.Find(childStateID); childBlockID != blockID {
					// Add child state to queue.
					queue = append(queue, childBlockID)
				}
			}
		}
	}

	return orderedBlocks
}