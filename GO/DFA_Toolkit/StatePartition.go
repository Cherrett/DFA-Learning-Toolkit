package DFA_Toolkit

// StatePartition struct which represents a State Partition.
type StatePartition struct {
	root []int		// Parent block of each state.
	size []int		// Size (score) of each block.
	link []int		// Index of next state within the block.
	changed []int	// Slice of changed states/blocks.

	isCopy bool

	//temp
	labelledStates []int
}

// Returns an initialized State Partition.
func NewStatePartition(dfa DFA) *StatePartition {
	// Initialize new State Partition struct.
	statePartition := new(StatePartition)
	statePartition.isCopy = false
	// Initialize empty slices.
	statePartition.root = make([]int, len(dfa.States))
	statePartition.size = make([]int, len(dfa.States))
	statePartition.link = make([]int, len(dfa.States))
	statePartition.labelledStates = make([]int, len(dfa.States))
	// Set root and link as element, and size (score) as 1.
	for i := 0; i < len(dfa.States); i++ {
		statePartition.root[i] = i
		statePartition.size[i] = 1
		statePartition.link[i] = i
		if dfa.States[i].StateStatus == ACCEPTING || dfa.States[i].StateStatus == REJECTING{
			statePartition.labelledStates[i] = 1
		}else{
			statePartition.labelledStates[i] = 0
		}
	}

	return statePartition
}

// Connects two states by finding their roots and comparing their respective
// size (score) values to keep the tree flat.
func (statePartition *StatePartition) union(stateID1 int, stateID2 int, state1Status StateStatus, state2Status StateStatus){
	if (state1Status == ACCEPTING && state2Status == REJECTING) || (state1Status == REJECTING && state2Status == ACCEPTING){
		panic("Invalid merge.")
	}

	// Get root (block index) of both states.
	stateID1Root := statePartition.Find(stateID1)
	stateID2Root := statePartition.Find(stateID2)

	// If their root is not equal, the states are merged (union) using the
	// linkBlocks function. If their root is equal, the states are already
	// within the same block so the merge is not done.
	if stateID1Root != stateID2Root{
		// Add State IDs joined to changed struct so merge can be undone.
		if statePartition.isCopy{
			statePartition.changed = append(statePartition.changed, stateID1)
			statePartition.changed = append(statePartition.changed, stateID2)
		}
		statePartition.linkBlocks(stateID1Root, stateID2Root, state1Status, state2Status)
	}
}

func (statePartition *StatePartition) linkBlocks(blockID1 int, blockID2 int, state1Status StateStatus, state2Status StateStatus){
	statePartition.link[blockID1], statePartition.link[blockID2] =
	 	statePartition.link[blockID2], statePartition.link[blockID1]

	if statePartition.size[blockID1] > statePartition.size[blockID2] {
		statePartition.root[blockID2] = blockID1
		if state2Status == ACCEPTING || state2Status == REJECTING{
			statePartition.labelledStates[blockID1]++
		}
	}else{
		statePartition.root[blockID1] = blockID2
		if statePartition.size[blockID1] == statePartition.size[blockID2]{
			statePartition.size[blockID2]++
		}
		if state1Status == ACCEPTING || state1Status == REJECTING{
			statePartition.labelledStates[blockID2]++
		}
	}
}

// Find traverses each parent element while compressing the
// levels to find the root element of the stateID
// If we attempt to access an element outside the array it returns -1
func (statePartition *StatePartition) Find(stateID int) int {
	if stateID > len(statePartition.root)-1 {
		panic("StateID out of range.")
	}

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
	state1Status := dfa.States[state1].StateStatus
	state2Status := dfa.States[state2].StateStatus
	if (state1Status == ACCEPTING && state2Status == REJECTING) || (state1Status == REJECTING && state2Status == ACCEPTING){
		return false
	}

	// store block status, set as unknown by default
	var blockStatus StateStatus = UNKNOWN

	// store the block transitions and set to -1 by default
	transitions := make([]int, len(dfa.SymbolMap))
	for i := range transitions {
		transitions[i] = -1
	}

	// merge states within state partition
	statePartition.union(state1, state2, state1Status, state2Status)

	// iterate over each state within the block containing the merged states
	for _, stateID := range statePartition.ReturnSet(state1){
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
	// Initialize new State Partition struct.
	copiedStatePartition := new(StatePartition)
	copiedStatePartition.isCopy = true
	copiedStatePartition.root = []int{}
	copiedStatePartition.size = []int{}
	copiedStatePartition.link = []int{}
	copiedStatePartition.changed = []int{}
	copiedStatePartition.labelledStates = []int{}
	copiedStatePartition.root = append(copiedStatePartition.root, statePartition.root...)
	copiedStatePartition.size = append(copiedStatePartition.size, statePartition.size...)
	copiedStatePartition.link = append(copiedStatePartition.link, statePartition.link...)
	copiedStatePartition.labelledStates = append(copiedStatePartition.labelledStates, statePartition.labelledStates...)
	return copiedStatePartition
}

func (statePartition *StatePartition) RollbackChanges(originalStatePartition *StatePartition){
	if statePartition.isCopy{
		for _, stateID := range statePartition.changed{
			statePartition.root[stateID] = originalStatePartition.root[stateID]
			statePartition.size[stateID] = originalStatePartition.size[stateID]
			statePartition.link[stateID] = originalStatePartition.link[stateID]
			statePartition.labelledStates[stateID] = originalStatePartition.labelledStates[stateID]
		}
		statePartition.changed = []int{}
	}else{
		return
	}
}

func (statePartition StatePartition) EDSMScore() int{
	sum := 0

	for blockID := range statePartition.labelledStates{
		labelledStates := statePartition.labelledStates[blockID]
		if labelledStates > 1{
			sum += labelledStates - 1
		}
	}

	return sum
}