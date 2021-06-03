package dfalearningtoolkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Block struct which represents a block within a partition.
type Block struct {
	Root        int        // Parent block of state.
	Size        int        // Size (score) of block.
	Link        int        // Index of next state within the block.
	Label       StateLabel // Label of block.
	Changed     bool       // Whether block has been changed.
	Transitions []int      // Transition Table where each element corresponds to a transition for each symbol.
}

// StatePartition struct which represents a State Partition.
type StatePartition struct {
	Blocks               []Block // Slice of blocks.
	BlocksCount          int     // Number of blocks within partition.
	AcceptingBlocksCount int     // Number of accepting blocks within partition.
	RejectingBlocksCount int     // Number of rejecting blocks within partition.
	AlphabetSize         int     // Size of alphabet.
	StartingStateID      int     // The ID of the starting state.

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
		for symbol := range referenceDFA.Alphabet {
			statePartition.Blocks[i].Transitions[symbol] = referenceDFA.States[i].Transitions[symbol]
		}
	}

	// Return initialized partition.
	return statePartition
}

// ChangedBlock updates the required fields to mark block as changed.
func (statePartition *StatePartition) ChangedBlock(blockID int) {
	// Check that state partition is a copy and that block is not modified.
	if statePartition.IsCopy && !statePartition.Blocks[blockID].Changed {
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
func (statePartition *StatePartition) Union(blockID1 int, blockID2 int) {
	// Mark both blocks as changed so merge can be undone.
	statePartition.ChangedBlock(blockID1)
	statePartition.ChangedBlock(blockID2)

	// Decrement blocks count.
	statePartition.BlocksCount--

	// If size of parent node is smaller than size of child node, switch
	// parent and child nodes.
	if statePartition.Blocks[blockID1].Size < statePartition.Blocks[blockID2].Size {
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

	// If label of parent is unlabelled and label of child is
	// labelled, set label of parent to label of child.
	if parent.Label == UNLABELLED && child.Label != UNLABELLED {
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
		if statePartition.Blocks[stateID].Root != statePartition.Blocks[statePartition.Blocks[stateID].Root].Root {
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

// ToQuotientDFA converts a State Partition to a quotient DFA and returns it.
func (statePartition *StatePartition) ToQuotientDFA() DFA {
	// Map to store corresponding new state for
	// each root block within state partition.
	blockToStateMap := map[int]int{}

	// Initialize resultant DFA to be returned by function.
	resultantDFA := NewDFA()

	// Get root blocks within state partition.
	rootBlocks := statePartition.RootBlocks()

	// Create alphabet within DFA.
	for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
		resultantDFA.AddSymbol()
	}

	// Create a new state within DFA for each root block and
	// set state label to block label.
	for _, stateID := range rootBlocks {
		blockToStateMap[stateID] = resultantDFA.AddState(statePartition.Blocks[stateID].Label)
	}

	// Update transitions using transitions within blocks and block to state map.
	for _, stateID := range rootBlocks {
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			if resultantState := statePartition.Blocks[stateID].Transitions[symbol]; resultantState > -1 {
				resultantDFA.States[blockToStateMap[stateID]].Transitions[symbol] = blockToStateMap[statePartition.Find(resultantState)]
			}
		}
	}

	// Set starting state using block to state map.
	resultantDFA.StartingStateID = blockToStateMap[statePartition.StartingBlock()]

	// Return populated resultant DFA.
	return resultantDFA
}

// ToQuotientDFAWithMapping converts a State Partition to a quotient DFA and returns it. This function
// also returns the state partition's blocks to state mapping.
func (statePartition *StatePartition) ToQuotientDFAWithMapping() (DFA, map[int]int) {
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
			if resultantState := statePartition.Blocks[stateID].Transitions[symbol]; resultantState > -1 {
				resultantDFA.States[blockToStateMap[stateID]].Transitions[symbol] = blockToStateMap[statePartition.Find(resultantState)]
			}
		}
	}

	// Set starting state using block to state map.
	resultantDFA.StartingStateID = blockToStateMap[statePartition.StartingBlock()]

	// Return populated resultant DFA and blocks to state map.
	return resultantDFA, blockToStateMap
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
	block1, block2 := &statePartition.Blocks[state1], &statePartition.Blocks[state2]

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
		if block1.Transitions[i] == -1 || block2.Transitions[i] == -1 {
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

// Clone returns a clone of the state partition.
func (statePartition StatePartition) Clone() StatePartition {
	// Initialize new StatePartition struct using state partition.
	clonedStatePartition := StatePartition{
		Blocks:               make([]Block, len(statePartition.Blocks)),
		BlocksCount:          statePartition.BlocksCount,
		AcceptingBlocksCount: statePartition.AcceptingBlocksCount,
		RejectingBlocksCount: statePartition.RejectingBlocksCount,
		AlphabetSize:         statePartition.AlphabetSize,
		StartingStateID:      statePartition.StartingStateID,
		IsCopy:               statePartition.IsCopy,
		ChangedBlocksCount:   statePartition.ChangedBlocksCount,
	}

	// Copy blocks.
	copy(clonedStatePartition.Blocks, statePartition.Blocks)

	// Iterate over each block within blocks slice.
	for blockID := range statePartition.Blocks {
		// Copy transitions.
		clonedStatePartition.Blocks[blockID].Transitions = make([]int, statePartition.AlphabetSize)
		copy(clonedStatePartition.Blocks[blockID].Transitions, statePartition.Blocks[blockID].Transitions)
	}

	// If state partition is already a copy.
	if statePartition.IsCopy {
		// Allocate slice for changed blocks.
		clonedStatePartition.ChangedBlocks = make([]int, len(statePartition.ChangedBlocks))
		// Copy changed blocks slice.
		copy(clonedStatePartition.ChangedBlocks, statePartition.ChangedBlocks)
	}

	// Return cloned state partition.
	return clonedStatePartition
}

// Copy returns a copy of the state partition in 'copy' (undo) mode.
func (statePartition StatePartition) Copy() StatePartition {
	// Panic if state partition is already a copy.
	if statePartition.IsCopy {
		panic("This state partition is already a copy.")
	}

	// Create a clone of the state partition.
	copiedStatePartition := statePartition.Clone()

	// Mark copied state partition as a copy.
	copiedStatePartition.IsCopy = true
	copiedStatePartition.ChangedBlocksCount = 0
	copiedStatePartition.ChangedBlocks = make([]int, len(statePartition.Blocks))

	// Return copied state partition.
	return copiedStatePartition
}

// RollbackChangesFrom reverts any changes made within state partition given the original state partition.
func (statePartition *StatePartition) RollbackChangesFrom(originalStatePartition StatePartition) {
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

// CopyChangesFrom copies the changes from one state partition to another and resets
// the changed values within the copied state partition.
func (statePartition *StatePartition) CopyChangesFrom(copiedStatePartition *StatePartition) {
	// If the state partition is a copy, copy values of changed blocks to original
	// state partition. Else, do nothing.
	if copiedStatePartition.IsCopy {
		// Set blocks count values to the new values.
		statePartition.BlocksCount = copiedStatePartition.BlocksCount
		statePartition.AcceptingBlocksCount = copiedStatePartition.AcceptingBlocksCount
		statePartition.RejectingBlocksCount = copiedStatePartition.RejectingBlocksCount

		// Iterate over each altered block.
		for _, blockID := range copiedStatePartition.ChangedBlocks[:copiedStatePartition.ChangedBlocksCount] {
			// Get changed block pointer from copied state partition.
			changedBlock := &copiedStatePartition.Blocks[blockID]
			// Get block pointer from original state partition.
			originalBlock := &statePartition.Blocks[blockID]

			// Update root, size, link, and label.
			originalBlock.Root = changedBlock.Root
			originalBlock.Size = changedBlock.Size
			originalBlock.Link = changedBlock.Link
			originalBlock.Label = changedBlock.Label

			// Set changed property of changed block to false.
			changedBlock.Changed = false

			// Copy transitions from changed to original.
			copy(originalBlock.Transitions, changedBlock.Transitions)
		}

		// Set the changed blocks count within the copied partition to 0.
		copiedStatePartition.ChangedBlocksCount = 0
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
func (statePartition *StatePartition) OrderedBlocks() []int {
	// Slice of boolean values to keep track of orders calculated.
	orderComputed := make([]bool, statePartition.BlocksCount)
	// Slice of integer values to keep track of ordered blocks.
	orderedBlocks := make([]int, statePartition.BlocksCount)

	// Get starting block ID.
	startingBlock := statePartition.StartingBlock()
	// Create a FIFO queue with starting block.
	queue := []int{startingBlock}
	// Add starting block to ordered blocks.
	orderedBlocks[0] = startingBlock
	// Mark starting block as computed.
	orderComputed[startingBlock] = true
	// Set index to 1.
	index := 1

	// Loop until queue is empty.
	for len(queue) > 0 {
		// Remove and store first state in queue.
		blockID := queue[0]
		queue = queue[1:]

		// Iterate over each symbol (alphabet) within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			// If transition from current state using current symbol is valid.
			if childStateID := statePartition.Blocks[blockID].Transitions[symbol]; childStateID >= 0 {
				// Get block ID of child state.
				childBlockID := statePartition.Find(childStateID)
				// If depth for child block has been computed, skip block.
				if !orderComputed[childBlockID] {
					// Add child block to queue.
					queue = append(queue, childBlockID)
					// Set the order of the current block.
					orderedBlocks[index] = childBlockID
					// Mark block as computed.
					orderComputed[childBlockID] = true
					// Increment current block order.
					index++
				}
			}
		}
	}

	return orderedBlocks
}

// StartingBlock returns the ID of the block which contains the starting state.
func (statePartition *StatePartition) StartingBlock() int {
	return statePartition.Find(statePartition.StartingStateID)
}

// DepthOfBlocks returns the depth of each block.
func (statePartition *StatePartition) DepthOfBlocks() map[int]int {
	// Create a FIFO queue with starting state.
	start := statePartition.StartingBlock()
	result := map[int]int{start: 0}
	queue := []int{start}

	for len(queue) > 0 {
		// Remove and store first state in queue.
		blockID := queue[0]
		queue = queue[1:]
		depth := result[blockID]

		for symbolID := 0; symbolID < statePartition.AlphabetSize; symbolID++ {
			if childStateID := statePartition.Blocks[blockID].Transitions[symbolID]; childStateID >= 0 {
				childBlockID := statePartition.Find(childStateID)
				if _, exists := result[childBlockID]; !exists {
					result[childBlockID] = depth + 1
					queue = append(queue, childBlockID)
				}
			}
		}
	}

	return result
}

// OrderOfBlocks returns the order of each block.
func (statePartition *StatePartition) OrderOfBlocks() map[int]int {
	// Map of integer values to keep track of ordered blocks.
	orderedBlocks := make(map[int]int, statePartition.BlocksCount)

	// Get starting block ID.
	startingBlock := statePartition.StartingBlock()
	// Create a FIFO queue with starting block.
	queue := []int{startingBlock}
	// Add starting block to ordered blocks.
	orderedBlocks[startingBlock] = 0
	// Set index to 1.
	index := 1

	// Loop until queue is empty.
	for len(queue) > 0 {
		// Remove and store first state in queue.
		blockID := queue[0]
		queue = queue[1:]

		// Iterate over each symbol (alphabet) within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			// If transition from current state using current symbol is valid.
			if childStateID := statePartition.Blocks[blockID].Transitions[symbol]; childStateID >= 0 {
				// Get block ID of child state.
				childBlockID := statePartition.Find(childStateID)
				// If depth for child block has been computed, skip block.
				if _, exists := orderedBlocks[childBlockID]; !exists {
					// Add child block to queue.
					queue = append(queue, childBlockID)
					// Set the order of the current block.
					orderedBlocks[childBlockID] = index
					// Increment current block order.
					index++
				}
			}
		}
	}

	return orderedBlocks
}

// ToJSON saves the StatePartition to a JSON file given a file path.
func (statePartition StatePartition) ToJSON(filePath string) bool {
	// Create file given a path/name.
	file, err := os.Create(filePath)

	// If file was not created successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Close file at end of function.
	defer file.Close()

	// Convert StatePartition to JSON.
	resultantJSON, err := json.MarshalIndent(statePartition, "", "\t")

	// If StatePartition was not converted successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Copy JSON to file created.
	_, err = io.Copy(file, bytes.NewReader(resultantJSON))

	// If JSON was not copied successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Return true if reached.
	return true
}

// StatePartitionFromJSON returns a StatePartition read from a JSON file
// given a file path. The boolean value returned is set to
// true if DFA was read successfully.
func StatePartitionFromJSON(filePath string) (StatePartition, bool) {
	// Open file from given a path/name.
	file, err := os.Open(filePath)

	// If file was not opened successfully,
	// return empty DFA and false.
	if err != nil {
		return StatePartition{}, false
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize empty StatePartition.
	resultantStatePartition := StatePartition{}

	// Convert JSON to StatePartition.
	err = json.NewDecoder(file).Decode(&resultantStatePartition)

	// If JSON was not converted successfully,
	// return empty StatePartition and false.
	if err != nil {
		return StatePartition{}, false
	}

	// Return populated StatePartition and true if reached.
	return resultantStatePartition, true
}
