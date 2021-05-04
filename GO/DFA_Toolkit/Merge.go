package dfatoolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
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

// GreedySearch deterministically merges all possible state pairs.
// The first valid merge with respect to the rejecting examples is chosen.
// Returns the resultant state partition and search data when no more valid merges are possible.
// Used by the regular positive and negative inference (RPNI) algorithm
func GreedySearch(statePartition StatePartition) (StatePartition, SearchData) {
	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize search data.
	searchData := SearchData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// Total merges and valid merges counter.
	totalMerges, totalValidMerges := 0, 0
	// Start timer.
	start := time.Now()

	// Get ordered blocks within partition.
	orderedBlocks := statePartition.OrderedBlocks()

	// Deterministically merge all valid merges by
	// iterating over root blocks within partition.
	for i := 1; i < len(orderedBlocks); i++ {
		//if copiedPartition.Blocks[i].Root != i{
		//	continue
		//}
		for j := 0; j < i; j++ {
			if copiedPartition.Blocks[j].Root != j {
				continue
			}
			// Increment merge count.
			totalMerges++
			// Check if states are mergeable.
			if copiedPartition.MergeStates(orderedBlocks[i], orderedBlocks[j]) {
				// Do not merge if states are within same block.
				if !statePartition.WithinSameBlock(orderedBlocks[i], orderedBlocks[j]) {
					// Increment valid merge count.
					totalValidMerges++

					// Copy changes to original state partition.
					statePartition.CopyChangesFrom(&copiedPartition)

					// Add merged state pair with score to search data.
					searchData.Merges = append(searchData.Merges, StatePairScore{orderedBlocks[i], orderedBlocks[j], 0})
					break
				}
			}

			// Undo merges from copied partition.
			copiedPartition.RollbackChangesFrom(statePartition)
		}
	}

	// Add total and valid merges counts to search data.
	searchData.AttemptedMergesCount = totalMerges
	searchData.ValidMergesCount = totalValidMerges
	// Add duration to search data.
	searchData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and search data.
	return statePartition, searchData
}

// GreedySearchUsingScoringFunction deterministically merges all possible state pairs.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and search data when no more valid merges are possible.
func GreedySearchUsingScoringFunction(statePartition StatePartition, scoringFunction ScoringFunction) (StatePartition, SearchData) {
	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize search data.
	searchData := SearchData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}
	// Total merges and valid merges counter.
	totalMerges, totalValidMerges := 0, 0
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
				totalMerges++
				// Check if states are mergeable.
				if copiedPartition.MergeStates(blocks[i], blocks[j]) {
					// Increment valid merge count.
					totalValidMerges++

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
		if highestScoringStatePair.Score != -1 {
			// Merge the state pairs with the highest score.
			copiedPartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)
			// Copy changes to original state partition.
			statePartition.CopyChangesFrom(&copiedPartition)

			// Add merged state pair with score to search data.
			searchData.Merges = append(searchData.Merges, highestScoringStatePair)

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}
		} else {
			break
		}
	}

	// Add total and valid merges counts to search data.
	searchData.AttemptedMergesCount = totalMerges
	searchData.ValidMergesCount = totalValidMerges
	// Add duration to search data.
	searchData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and search data.
	return statePartition, searchData
}

// FastWindowedSearchUsingScoringFunction deterministically merges state pairs within a given window.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and search data when no more valid merges are possible.
func FastWindowedSearchUsingScoringFunction(statePartition StatePartition, windowSize int, windowGrow float64, scoringFunction ScoringFunction) (StatePartition, SearchData) {
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
	// Initialize search data.
	searchData := SearchData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// Total merges and valid merges counter.
	totalMerges, totalValidMerges := 0, 0

	// Start timer.
	start := time.Now()

	// Iterate until stopped.
	for {
		// State pair with the highest score.
		highestScoringStatePair := StatePairScore{-1, -1, -1}
		// Get ordered blocks within partition.
		orderedBlocks := statePartition.OrderedBlocks()
		// Set previous window size to 0.
		previousWindowSize := 0
		// Set window size to window size parameter
		// or length of ordered blocks if smaller.
		windowSize = util.Min(windowSize, len(orderedBlocks))

		// Loop until no more deterministic merges are available within all possible windows.
		for {
			// Get all valid merges within window and compute their score.
			for i := 0; i < windowSize; i++ {
				for j := i + 1; j < windowSize; j++ {
					if j >= previousWindowSize {
						// Increment merge count.
						totalMerges++
						// Check if states are mergeable.
						if copiedPartition.MergeStates(orderedBlocks[i], orderedBlocks[j]) {
							// Increment valid merge count.
							totalValidMerges++

							// Calculate score.
							score := scoringFunction(orderedBlocks[i], orderedBlocks[j], statePartition, copiedPartition)

							// If score is bigger than state pair with the highest score,
							// set current state pair to state pair with the highest score.
							if score > highestScoringStatePair.Score {
								highestScoringStatePair = StatePairScore{
									State1: orderedBlocks[i],
									State2: orderedBlocks[j],
									Score:  score,
								}
							}
						}

						// Undo merges from copied partition.
						copiedPartition.RollbackChangesFrom(statePartition)
					}
				}
			}

			// Check if any deterministic merges were found.
			if highestScoringStatePair.Score != -1 {
				break
				// No more possible merges were found so increase window size.
			} else {
				// If the window size is biggest possible, break loop
				// and return the most recent State Partition.
				if windowSize >= len(orderedBlocks) {
					break
				}

				previousWindowSize = windowSize
				windowSize = util.Min(int(math.Round(float64(windowSize)*windowGrow)), len(orderedBlocks))
			}
		}

		if highestScoringStatePair.Score != -1 {
			// Merge the state pairs with the highest score.
			copiedPartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)
			// Copy changes to original state partition.
			statePartition.CopyChangesFrom(&copiedPartition)

			// Add merged state pair with score to search data.
			searchData.Merges = append(searchData.Merges, highestScoringStatePair)
		} else {
			break
		}
	}

	// Add total and valid merges counts to search data.
	searchData.AttemptedMergesCount = totalMerges
	searchData.ValidMergesCount = totalValidMerges
	// Add duration to search data.
	searchData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and search data.
	return statePartition, searchData
}

// WindowedSearchUsingScoringFunction deterministically merges state pairs within a given window.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and search data when no more valid merges are possible.
func WindowedSearchUsingScoringFunction(statePartition StatePartition, windowSize int, windowGrow float64, scoringFunction ScoringFunction) (StatePartition, SearchData) {
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
	// Initialize search data.
	searchData := SearchData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// Total merges and valid merges counter.
	totalMerges, totalValidMerges := 0, 0
	// Get ordered blocks within partition.
	orderedBlocks := statePartition.OrderedBlocks()
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
		windowSize = util.Min(windowSize, len(orderedBlocks))

		// Loop until no more deterministic merges are available within all possible windows.
		for {
			// Get all valid merges within window and compute their score.
			for i := 0; i < windowSize; i++ {
				for j := i + 1; j < windowSize; j++ {
					if j >= previousWindowSize {
						// Increment merge count.
						totalMerges++

						// Check if states are mergeable.
						if copiedPartition.MergeStates(orderedBlocks[i], orderedBlocks[j]) {
							// Do not compute score if states are within same block.
							if !statePartition.WithinSameBlock(orderedBlocks[i], orderedBlocks[j]) {
								// Increment valid merge count.
								totalValidMerges++

								// Calculate score.
								score := scoringFunction(orderedBlocks[i], orderedBlocks[j], statePartition, copiedPartition)

								// If score is bigger than state pair with the highest score,
								// set current state pair to state pair with the highest score.
								if score > highestScoringStatePair.Score {
									highestScoringStatePair = StatePairScore{
										State1: orderedBlocks[i],
										State2: orderedBlocks[j],
										Score:  score,
									}
								}
							}
						}

						// Undo merges from copied partition.
						copiedPartition.RollbackChangesFrom(statePartition)
					}
				}
			}

			// Check if any deterministic merges were found.
			if highestScoringStatePair.Score != -1 {
				break
				// No more possible merges were found so increase window size.
			} else {
				// If the window size is biggest possible, break loop
				// and return the most recent State Partition.
				if windowSize >= len(orderedBlocks) {
					break
				}
				// Set window size before to window size.
				previousWindowSize = windowSize
				// Get new window size which is the smallest from the window size multiplied
				// by window grow or the number of blocks within initial state partition.
				windowSize = util.Min(int(math.Round(float64(windowSize)*windowGrow)), len(orderedBlocks))
			}
		}

		if highestScoringStatePair.Score != -1 {
			// Merge the state pairs with the highest score.
			copiedPartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)
			// Copy changes to original state partition.
			statePartition.CopyChangesFrom(&copiedPartition)

			// Add merged state pair with score to search data.
			searchData.Merges = append(searchData.Merges, highestScoringStatePair)
		} else {
			break
		}
	}

	// Add total and valid merges counts to search data.
	searchData.AttemptedMergesCount = totalMerges
	searchData.ValidMergesCount = totalValidMerges
	// Add duration to search data.
	searchData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and search data.
	return statePartition, searchData
}

// BlueFringeSearchUsingScoringFunction deterministically merges possible state pairs within red-blue sets.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and search data when no more valid merges are possible.
func BlueFringeSearchUsingScoringFunction(statePartition StatePartition, scoringFunction ScoringFunction) (StatePartition, SearchData) {
	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize search data.
	searchData := SearchData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}
	// Total merges and valid merges counter.
	totalMerges, totalValidMerges := 0, 0
	// Start timer.
	start := time.Now()

	// Slice of state pairs to keep track of computed scores.
	scoresComputed := map[StateIDPair]util.Void{}

	// Initialize set of red states to starting state.
	red := map[int]util.Void{statePartition.StartingBlock(): util.Null}
	// Slice to store red states in insertion order.
	orderedRed := []int{statePartition.StartingBlock()}

	// Initialize merged flag to false.
	merged := false

	// Iterate until stopped.
	for {
		// Initialize set of blue states to empty set.
		blue := map[int]util.Void{}
		// Slice to store blue states in insertion order.
		var orderedBlue []int

		// Iterate over every red state.
		for _, element := range orderedRed {
			// Iterate over each symbol within DFA.
			for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
				if resultantStateID := statePartition.Blocks[element].Transitions[symbol]; resultantStateID > -1 {
					// Store resultant stateID from red state.
					resultantStateID = statePartition.Find(resultantStateID)
					// If transition is valid and resultant state is not red,
					// add resultant state to blue set.
					if _, exists := red[resultantStateID]; resultantStateID != -1 && !exists {
						// If resultant state is not already in blue set, add to
						// both blue set and ordered blue set.
						if _, exists := blue[resultantStateID]; !exists {
							orderedBlue = append(orderedBlue, resultantStateID)
							blue[resultantStateID] = util.Null
						}
					}
				}
			}
		}

		// If blue set is empty, break loop and return resultant DFA.
		if len(blue) == 0 {
			break
		}

		// Iterate over every blue state in insertion order.
		for _, blueElement := range orderedBlue {
			// Set merged flag to false.
			merged = false
			// Iterate over every red state in insertion order.
			for _, redElement := range orderedRed {
				// If scores for the current state pair has already been
				// computed, set merged flag to true and skip merge.
				if _, valid := scoresComputed[StateIDPair{blueElement, redElement}]; valid {
					merged = true
				} else {
					// If scores for the current state pair has not
					// been computed, attempt to merge state pair.

					// Increment merge count.
					totalMerges++
					// If states are mergeable, calculate score and add to detMerges.
					if copiedPartition.MergeStates(blueElement, redElement) {
						// Increment valid merge count.
						totalValidMerges++

						// Set the state pairs score as computed.
						scoresComputed[StateIDPair{blueElement, redElement}] = util.Null
						scoresComputed[StateIDPair{redElement, blueElement}] = util.Null

						// Calculate score.
						score := scoringFunction(blueElement, redElement, statePartition, copiedPartition)

						// If score is bigger than state pair with the highest score,
						// set current state pair to state pair with the highest score.
						if score > highestScoringStatePair.Score {
							highestScoringStatePair = StatePairScore{
								State1: blueElement,
								State2: redElement,
								Score:  score,
							}
						}

						// Set merged flag to true.
						merged = true
					}
					// Undo merge.
					copiedPartition.RollbackChangesFrom(statePartition)
				}
			}

			// If merged flag is false, add current blue state
			// to red states set and ordered set and exit loop.
			if !merged {
				red[blueElement] = util.Null
				orderedRed = append(orderedRed, blueElement)
				break
			}
		}

		// If merged flag is true, merge highest scoring merge.
		if merged {
			// Merge the state pairs with the highest score.
			copiedPartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)
			// Copy changes to original state partition.
			statePartition.CopyChangesFrom(&copiedPartition)

			// Add merged state pair with score to search data.
			searchData.Merges = append(searchData.Merges, highestScoringStatePair)

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}

			// Slice of state pairs to keep track of computed scores.
			scoresComputed = map[StateIDPair]util.Void{}

			// Set set of red states to starting state.
			red = map[int]util.Void{statePartition.StartingBlock(): util.Null}
			// Set slice to store red states in insertion order.
			orderedRed = []int{statePartition.StartingBlock()}
		}
	}

	// Add total and valid merges counts to search data.
	searchData.AttemptedMergesCount = totalMerges
	searchData.ValidMergesCount = totalValidMerges
	// Add duration to search data.
	searchData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and search data.
	return statePartition, searchData
}
