package dfalearningtoolkit

import (
	"github.com/Cherrett/DFA-Learning-Toolkit/util"
	"math"
	"time"
)

// StatePairScore struct to store state pairs and their merge score.
type StatePairScore struct {
	State1 int     // StateID for first state.
	State2 int     // StateID for second state.
	Score  float64 // Score of merge for given states.
}

// ScoringFunction takes two stateIDs and two state partitions as input and returns a score as a float.
type ScoringFunction func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64

// RedBlueSearch deterministically merges all possible state pairs within red-blue framework.
// The first valid merge with respect to the rejecting examples is chosen.
// Returns the resultant state partition and merge data when no more valid merges are possible.
// Used by the regular positive and negative inference (RPNI) algorithm.
func RedBlueSearch(statePartition StatePartition) (StatePartition, MergeData) {
	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize merge data.
	mergeData := MergeData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// Start timer.
	start := time.Now()

	// Slice to store red states.
	redStates := []int{statePartition.StartingBlock()}
	// Generated slice of blue states from red states.
	blueStates := UpdateRedBlueSets(&statePartition, &redStates)

	// Iterate until blue set is empty.
	for len(blueStates) > 0 {
		// Remove and store first state in blue states queue.
		blueState := blueStates[0]
		blueStates = blueStates[1:]

		// Set merged flag to false.
		merged := false

		// Iterate over every red state.
		for _, redState := range redStates {
			// Increment merge count.
			mergeData.AttemptedMergesCount++

			// Check if states are mergeable.
			if copiedPartition.MergeStates(redState, blueState) {
				// Increment valid merge count.
				mergeData.ValidMergesCount++

				// Copy changes to original state partition.
				statePartition.CopyChangesFrom(&copiedPartition)

				// Set merged flag to true.
				merged = true
				break
			}

			// Undo merges from copied partition.
			copiedPartition.RollbackChangesFrom(statePartition)
		}

		// If merged flag is false, add current blue state
		// to red states set and ordered set.
		if !merged {
			redStates = append(redStates, blueState)
		}

		// Update red and blue states using UpdateOrderedRedBlueSets function.
		blueStates = UpdateRedBlueSets(&statePartition, &redStates)
	}

	// Add duration to merge data.
	mergeData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and merge data.
	return statePartition, mergeData
}

// UpdateRedBlueSets updates the red and blue sets given the state partition and the red set within the Red-Blue framework
// such as the RedBlueSearch function. It returns the blue set and modifies the red set via its pointer. This is
// used when the state partition is changed or when new states have been added to the red set.
func UpdateRedBlueSets(statePartition *StatePartition, redStates *[]int) []int {
	// Step 1 - Gather root of old red states and store in map declared below.

	// Initialize set of red root states (blocks) to empty set.
	redSet := make(map[int]util.Void, len(*redStates))

	// Iterate over every red state.
	for i := range *redStates {
		// Get red state.
		redState := &(*redStates)[i]
		// Replace red state with its parent using Find.
		*redState = statePartition.Find(*redState)
		// Add red state to red set.
		redSet[*redState] = util.Null
	}

	// Step 2 - Gather all blue states and store in slice declared below.

	// Initialize set of blue states to empty slice.
	var blueStates []int

	// Iterate over every red state.
	for _, redState := range *redStates {
		// Iterate over each symbol within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			// If transition is valid and resultant state is not red,
			// add the parent block of the resultant state to blue set.
			if resultantStateID := statePartition.Blocks[redState].Transitions[symbol]; resultantStateID > -1 {
				resultantStateID = statePartition.Find(resultantStateID)

				if _, exists := redSet[resultantStateID]; !exists {
					blueStates = append(blueStates, resultantStateID)
				}
			}
		}
	}

	// Return populated slice of blue states.
	return blueStates
}

// ExhaustiveSearchUsingScoringFunction deterministically merges all possible state pairs.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and merge data when no more valid merges are possible.
func ExhaustiveSearchUsingScoringFunction(statePartition StatePartition, scoringFunction ScoringFunction) (StatePartition, MergeData) {
	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize merge data.
	mergeData := MergeData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}
	// Start timer.
	start := time.Now()

	// Loop until no more deterministic merges are available.
	for {
		// Get root blocks within partition.
		blocks := statePartition.RootBlocks()

		// Get all valid merges and compute their score by
		// iterating over root blocks within partition.
		for i := 0; i < len(blocks); i++ {
			for j := i + 1; j < len(blocks); j++ {
				// Increment merge count.
				mergeData.AttemptedMergesCount++

				// Check if states are mergeable.
				if copiedPartition.MergeStates(blocks[i], blocks[j]) {
					// Increment valid merge count.
					mergeData.ValidMergesCount++

					// Calculate score.
					score := scoringFunction(blocks[i], blocks[j], statePartition, copiedPartition)

					// If score is bigger than state pair with the highest score,
					// set current state pair to state pair with the highest score.
					if score > highestScoringStatePair.Score {
						highestScoringStatePair = StatePairScore{
							State1: blocks[i],
							State2: blocks[j],
							Score:  score,
						}
					}
				}

				// Undo merges from copied partition.
				copiedPartition.RollbackChangesFrom(statePartition)
			}
		}

		// Check if any deterministic merges were found.
		if highestScoringStatePair.Score >= 0 {
			// Merge the state pairs with the highest score.
			copiedPartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)
			// Copy changes to original state partition.
			statePartition.CopyChangesFrom(&copiedPartition)

			// Add merged state pair with score to merge data.
			mergeData.Merges = append(mergeData.Merges, highestScoringStatePair)

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}
		} else {
			break
		}
	}

	// Add duration to merge data.
	mergeData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and merge data.
	return statePartition, mergeData
}

// WindowedSearchUsingScoringFunction deterministically merges state pairs within a given window.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and merge data when no more valid merges are possible.
func WindowedSearchUsingScoringFunction(statePartition StatePartition, windowSize int, windowGrow float64, scoringFunction ScoringFunction) (StatePartition, MergeData) {
	// Parameter Error Checking.
	if windowSize < 1 {
		panic("Window Size cannot be smaller than 1.")
	}
	if windowGrow <= 1.00 {
		panic("Window Grow cannot be smaller or equal to 1.")
	}

	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize merge data.
	mergeData := MergeData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// Get ordered blocks within partition.
	window := statePartition.OrderedBlocks()
	// Start timer.
	start := time.Now()

	// Iterate until stopped.
	for {
		// State pair with the highest score.
		highestScoringStatePair := StatePairScore{-1, -1, -1}

		// Set window size before to 0.
		previousWindowSize := 0

		// Set window size to window size parameter
		// or length of ordered blocks if smaller.
		windowSize = util.Min(windowSize, len(window))

		// Loop until no more deterministic merges are available within all possible windows.
		for {
			// Get all valid merges within window and compute their score.
			for i := 0; i < windowSize; i++ {
				for j := i + 1; j < windowSize; j++ {
					if j >= previousWindowSize {
						// Increment merge count.
						mergeData.AttemptedMergesCount++

						// Check if states are mergeable.
						if copiedPartition.MergeStates(window[i], window[j]) {
							// Increment valid merge count.
							mergeData.ValidMergesCount++

							// Calculate score.
							score := scoringFunction(window[i], window[j], statePartition, copiedPartition)

							// If score is bigger than state pair with the highest score,
							// set current state pair to state pair with the highest score.
							if score > highestScoringStatePair.Score {
								highestScoringStatePair = StatePairScore{
									State1: window[i],
									State2: window[j],
									Score:  score,
								}
							}
						}

						// Undo merges from copied partition.
						copiedPartition.RollbackChangesFrom(statePartition)
					}
				}
			}

			// Break loop if no deterministic merges were found or if the window size is the biggest possible.
			if highestScoringStatePair.Score >= 0 || windowSize >= len(window) {
				break
			} else {
				// No more possible merges were found so increase window size.
				// Set window size before to window size.
				previousWindowSize = windowSize
				// Get new window size which is the smallest from the window size multiplied
				// by window grow or the number of blocks within initial state partition.
				windowSize = util.Min(int(math.Round(float64(windowSize)*windowGrow)), len(window))
			}
		}

		if highestScoringStatePair.Score >= 0 {
			// Merge the state pairs with the highest score.
			copiedPartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)
			// Copy changes to original state partition.
			statePartition.CopyChangesFrom(&copiedPartition)

			// Add merged state pair with score to merge data.
			mergeData.Merges = append(mergeData.Merges, highestScoringStatePair)
		} else {
			break
		}

		// Update new window.
		window = UpdateWindow(window, statePartition)
	}

	// Add duration to merge data.
	mergeData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and merge data.
	return statePartition, mergeData
}

// UpdateWindow updates a window given the state partition within a Windowed framework such as the
// WindowedSearchUsingScoringFunction function. It returns the new window as a slice of integers.
// This works by gathering the root of each block within the previous window and assigns it to the
// position of the first index of any block which is part of that block. This is used to avoid attempting
// merges more than once within a windowed search.
func UpdateWindow(window []int, statePartition StatePartition) []int {
	// Gather root of ordered blocks and store in map and slice declared below.

	// Initialize set of root states (blocks) to empty set.
	rootBlocks := map[int]util.Void{}
	// Slice to store new blocks in canonical order (new window).
	var newWindow []int

	// Iterate over every state within window.
	for _, state := range window {
		// Get root block of current state.
		root := statePartition.Find(state)

		// If root is already visited, skip.
		// Else, add root to map and slice.
		if _, exists := rootBlocks[root]; !exists {
			rootBlocks[root] = util.Null
			newWindow = append(newWindow, root)
		}
	}

	// Return new red set and populated slice of red states in canonical order.
	return newWindow
}

// BlueFringeSearchUsingScoringFunction deterministically merges possible state pairs within red-blue sets.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and merge data when no more valid merges are possible.
func BlueFringeSearchUsingScoringFunction(statePartition StatePartition, scoringFunction ScoringFunction) (StatePartition, MergeData) {
	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize merge data.
	mergeData := MergeData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}
	// Start timer.
	start := time.Now()

	// Slice to store red states in canonical order.
	redStates := []int{statePartition.StartingBlock()}
	// Generated slice of ordered blue states from red states.
	blueStates := GenerateOrderedBlueSetFromRedSet(&statePartition, map[int]util.Void{redStates[0]: util.Null})

	// Initialize merged flag to false.
	merged := false

	// Iterate until blue set is empty.
	for len(blueStates) != 0 {
		// Iterate over every blue state in canonical order.
		for _, blueState := range blueStates {
			// Set merged flag to false.
			merged = false
			// Iterate over every red state in canonical order.
			for _, redState := range redStates {
				// Increment merge count.
				mergeData.AttemptedMergesCount++

				// Check if states are mergeable.
				if copiedPartition.MergeStates(blueState, redState) {
					// Increment valid merge count.
					mergeData.ValidMergesCount++

					// Calculate score.
					score := scoringFunction(blueState, redState, statePartition, copiedPartition)

					// If score is bigger than state pair with the highest score,
					// set current state pair to state pair with the highest score.
					if score > highestScoringStatePair.Score {
						highestScoringStatePair = StatePairScore{
							State1: blueState,
							State2: redState,
							Score:  score,
						}
					}

					// Set merged flag to true.
					merged = true
				}

				// Undo merge.
				copiedPartition.RollbackChangesFrom(statePartition)
			}

			// If merged flag is false, add current blue state
			// to red states set and ordered set and exit loop.
			if !merged {
				redStates = append(redStates, blueState)
				break
			}
		}

		// If merged flag is true, merge highest scoring merge.
		if merged {
			// Merge the state pairs with the highest score.
			copiedPartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)
			// Copy changes to original state partition.
			statePartition.CopyChangesFrom(&copiedPartition)

			// Add merged state pair with score to merge data.
			mergeData.Merges = append(mergeData.Merges, highestScoringStatePair)

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}
		}

		// Update red and blue states using UpdateOrderedRedBlueSets function.
		redStates, blueStates = UpdateOrderedRedBlueSets(&statePartition, redStates)
	}

	// Add duration to merge data.
	mergeData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and merge data.
	return statePartition, mergeData
}

// GenerateOrderedBlueSetFromRedSet generates the blue set given the state partition and the red set within the Red-Blue framework
// such as the BlueFringeSearchUsingScoringFunction function. It generates and returns the blue set in canonical order.
func GenerateOrderedBlueSetFromRedSet(statePartition *StatePartition, redSet map[int]util.Void) []int {
	// Step 1 - Gather all blue states and store in map declared below.

	// Initialize set of blue states to empty set.
	blue := map[int]util.Void{}

	// Iterate over every red state.
	for redState := range redSet {
		// Iterate over each symbol within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			if resultantStateID := statePartition.Blocks[redState].Transitions[symbol]; resultantStateID > -1 {
				// Store resultant stateID from red state.
				resultantStateID = statePartition.Find(resultantStateID)
				// If transition is valid and resultant state is not red,
				// add resultant state to blue set.
				if _, exists := redSet[resultantStateID]; !exists {
					blue[resultantStateID] = util.Null
				}
			}
		}
	}

	// Step 2 - Sort blue states by canonical order and store in slice declared below.

	// Slice to store blue states in canonical order.
	var orderedBlue []int

	// Get starting block ID.
	startingBlock := statePartition.StartingBlock()
	// Create a FIFO queue with starting block.
	queue := []int{statePartition.Find(startingBlock)}
	// Slice of boolean values to keep track of orders calculated.
	orderComputed := make([]bool, len(statePartition.Blocks))
	// Mark starting block as computed.
	orderComputed[startingBlock] = true

	// If starting block is in blue set, add to ordered slice.
	if _, exists := blue[startingBlock]; exists {
		orderedBlue = append(orderedBlue, startingBlock)
	}

	// Loop until queue is empty.
	for len(queue) > 0 {
		// Remove and store first state in queue.
		blockID := queue[0]
		queue = queue[1:]

		// Iterate over each symbol (alphabet) within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			// If transition from current state using current symbol is valid and is not a loop to the current state.
			if childStateID := statePartition.Blocks[blockID].Transitions[symbol]; childStateID >= 0 {
				// Get block ID of child state.
				childBlockID := statePartition.Find(childStateID)
				// If depth for child block has been computed, skip block.
				if !orderComputed[childBlockID] {
					// Add child block to queue.
					queue = append(queue, childBlockID)

					// If block is in blue set, add to ordered slice.
					if _, exists := blue[childBlockID]; exists {
						orderedBlue = append(orderedBlue, childBlockID)
					}

					// Mark block as computed.
					orderComputed[childBlockID] = true
				}
			}
		}
	}

	// Return populated slice of blue states in canonical order.
	return orderedBlue
}

// UpdateOrderedRedBlueSets updates the red and blue sets given the state partition and the red set within the Red-Blue framework
// such as the BlueFringeSearchUsingScoringFunction function. It returns the red and blue sets in canonical order. This is
// used when the state partition is changed or when new states have been added to the red set.
func UpdateOrderedRedBlueSets(statePartition *StatePartition, redStates []int) ([]int, []int) {
	// Step 1 - Gather root of old red states and store in map declared below.

	// Initialize set of red root states (blocks) to empty set.
	redSet := make(map[int]util.Void, len(redStates))

	// Iterate over every red state.
	for _, redState := range redStates {
		// Get root block of red state.
		root := statePartition.Find(redState)

		// Add root block to red set.
		redSet[root] = util.Null
	}

	// Step 2 - Gather all blue states and store in map declared below.

	// Initialize set of blue states to empty set.
	blueSet := map[int]util.Void{}

	// Iterate over every red state.
	for redState := range redSet {
		// Iterate over each symbol within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			if resultantStateID := statePartition.Blocks[redState].Transitions[symbol]; resultantStateID > -1 {
				// Store resultant stateID from red state.
				resultantStateID = statePartition.Find(resultantStateID)
				// If transition is valid and resultant state is not red,
				// add resultant state to blue set.
				if _, exists := redSet[resultantStateID]; !exists {
					blueSet[resultantStateID] = util.Null
				}
			}
		}
	}

	// Step 3 - Sort red and blue states by canonical order and store in slice declared below.

	// Slice to store red states in canonical order.
	var orderedRed []int
	// Slice to store red states in canonical order.
	var orderedBlue []int

	// Get starting block ID.
	startingBlock := statePartition.StartingBlock()
	// Create a FIFO queue with starting block.
	queue := []int{statePartition.Find(startingBlock)}
	// Slice of boolean values to keep track of orders calculated.
	orderComputed := make([]bool, len(statePartition.Blocks))
	// Mark starting block as computed.
	orderComputed[startingBlock] = true

	// If starting block is in red set, add to red ordered slice.
	if _, exists := redSet[startingBlock]; exists {
		orderedRed = append(orderedRed, startingBlock)
		// Else, if starting block is in blue set, add to blue ordered slice.
	} else if _, exists = blueSet[startingBlock]; exists {
		orderedBlue = append(orderedBlue, startingBlock)
	}

	// Loop until queue is empty.
	for len(queue) > 0 {
		// Remove and store first state in queue.
		blockID := queue[0]
		queue = queue[1:]

		// Iterate over each symbol (alphabet) within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			// If transition from current state using current symbol is valid and is not a loop to the current state.
			if childStateID := statePartition.Blocks[blockID].Transitions[symbol]; childStateID != -1 {
				// Get block ID of child state.
				childBlockID := statePartition.Find(childStateID)
				// If depth for child block has been computed, skip block.
				if !orderComputed[childBlockID] {
					// Add child block to queue.
					queue = append(queue, childBlockID)

					// If block is in red set, add to ordered red slice.
					if _, exists := redSet[childBlockID]; exists {
						orderedRed = append(orderedRed, childBlockID)
						// Else, if block is in blue set, add to ordered blue slice.
					} else if _, exists = blueSet[childBlockID]; exists {
						orderedBlue = append(orderedBlue, childBlockID)
					}

					// Mark block as computed.
					orderComputed[childBlockID] = true
				}
			}
		}
	}

	// Return populated slice of red states and  populated slice of blue states in canonical order.
	return orderedRed, orderedBlue
}
