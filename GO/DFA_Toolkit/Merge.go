package dfatoolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
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
// Returns the resultant state partition when no more valid merges are possible.
func GreedySearch(statePartition StatePartition, scoringFunction ScoringFunction) StatePartition {
	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}
	start := time.Now()
	totalMerges := 0

	// Loop until no more deterministic merges are available.
	for {
		// Copy the state partition for undoing merging.
		copiedPartition := statePartition.Copy()

		// Get root blocks within partition.
		blocks := statePartition.RootBlocks()

		// Get all valid merges and compute their score by
		// iterating over root blocks within partition.
		for i := 0; i < len(blocks); i++ {
			for j := i + 1; j < len(blocks); j++ {
				totalMerges++
				// Check if states are mergeable.
				if copiedPartition.MergeStates(blocks[i], blocks[j]) {
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
				copiedPartition.RollbackChanges(statePartition)
			}
		}

		// Check if any deterministic merges were found.
		if highestScoringStatePair.Score != -1 {
			// Merge the state pairs with the highest score.
			statePartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}
		} else {
			break
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant state partition.
	return statePartition
}

// WindowedSearch deterministically merges state pairs within a given window.
// Returns the resultant state partition when no more valid merges are possible.
func WindowedSearch(statePartition StatePartition, windowSize int, windowGrow float64, scoringFunction ScoringFunction) StatePartition {
	start := time.Now()
	totalMerges := 0

	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}

	// Get ordered blocks within partition.
	orderedBlocks := statePartition.OrderedBlocks()

	// Iterate until stopped.
	for {
		// Window values.
		windowMin := 0
		windowMax := util.Min(windowSize, len(orderedBlocks))

		// Copy the state partition for undoing merging.
		copiedPartition := statePartition.Copy()

		// Loop until no more deterministic merges are available within all possible windows.
		for {
			// Get all valid merges within window and compute their score.
			for i := 0; i < windowMax; i++ {
				for j := windowMin; j < windowMax; j++ {
					if i < j {
						totalMerges++
						// Check if states are mergeable.
						if copiedPartition.MergeStates(orderedBlocks[i], orderedBlocks[j]) {
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
						copiedPartition.RollbackChanges(statePartition)
					}
				}
			}

			// Check if any deterministic merges were found.
			if highestScoringStatePair.Score != -1 {
				break
				// No more possible merges were found so increase window size.
			} else {
				windowMin += windowSize
				windowSize = int(math.Round(float64(windowSize) * windowGrow))
				windowMax = util.Min(windowMax + windowSize, len(orderedBlocks))

				// If the window size is out of bounds, break loop and return the
				// most recent DFA found.
				if windowMin >= len(orderedBlocks) {
					break
				}
			}
		}

		if highestScoringStatePair.Score != -1 {
			// Merge the state pairs with the highest score.
			statePartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Update ordered states within merged partition.
			orderedBlocks = statePartition.OrderedBlocks()

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}
		} else {
			break
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return statePartition
}

// BlueFringeSearch deterministically merges possible state pairs within red-blue sets.
// Returns the resultant state partition when no more valid merges are possible.
// TODO: Update using new StatePartition.
func BlueFringeSearch(statePartition StatePartition, scoringFunction ScoringFunction) StatePartition {
	start := time.Now()
	totalMerges := 0

	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}

	// Slice of state pairs to keep track of computed scores.
	scoresComputed := map[StateIDPair]util.Void{}

	// Initialize set of red states to starting state.
	red := map[int]util.Void{statePartition.StartingBlock(): util.Null}

	// Initialize merged flag to false.
	merged := false

	// Copy the state partition for undoing merging.
	copiedPartition := statePartition.Copy()

	// Iterate until stopped.
	for {
		// Initialize set of blue states to empty set.
		blue := map[int]util.Void{}

		// Iterate over every red state.
		for element := range red {
			// Iterate over each symbol within DFA.
			for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
				if resultantStateID := statePartition.Blocks[element].Transitions[symbol]; resultantStateID > -1{
					// Store resultant stateID from red state.
					resultantStateID = statePartition.Find(resultantStateID)
					// If transition is valid and resultant state is not red,
					// add resultant state to blue set.
					if _, exists := red[resultantStateID]; resultantStateID != -1 && !exists {
						blue[resultantStateID] = util.Null
					}
				}
			}
		}

		// If blue set is empty, break loop and return resultant DFA.
		if len(blue) == 0 {
			break
		}

		// Iterate over every blue state.
		for blueElement := range blue {
			// Set merged flag to false.
			merged = false
			// Iterate over every red state.
			for redElement := range red {
				// If scores for the current state pair has already been
				// computed, set merged flag to true and skip merge.
				if _, valid := scoresComputed[StateIDPair{blueElement, redElement}]; valid {
					merged = true
				} else {
					// If scores for the current state pair has not
					// been computed, attempt to merge state pair.
					totalMerges++
					// If states are mergeable, calculate score and add to detMerges.
					if copiedPartition.MergeStates(blueElement, redElement) {
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
					copiedPartition.RollbackChanges(statePartition)
				}
			}

			// If merged flag is false, add current blue state
			// to red states set and exit loop.
			if !merged {
				red[blueElement] = util.Null
				break
			}
		}

		// If merged flag is true, merge highest scoring merge.
		if merged {
			// Merge the state pairs with the highest score.
			statePartition.MergeStates(highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Copy the state partition for undoing merging.
			copiedPartition = statePartition.Copy()

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}

			// Slice of state pairs to keep track of computed scores.
			scoresComputed = map[StateIDPair]util.Void{}

			// Initialize set of red states to starting state.
			red = map[int]util.Void{statePartition.StartingBlock(): util.Null}
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return statePartition
}
