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

// ScoringFunction takes a state partition as input and returns an integer score.
type ScoringFunction func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition, dfa DFA) float64

// GreedySearch deterministically merges all possible state pairs.
// Returns the resultant DFA when no more valid merges are possible.
func GreedySearch(dfa DFA, scoringFunction ScoringFunction) DFA {
	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}
	start := time.Now()
	totalMerges := 0

	// Loop until no more deterministic merges are available.
	for {
		// Convert DFA to StatePartition for state merging.
		partition := dfa.ToStatePartition()
		// Copy the state partition for undoing merging.
		snapshot := partition.Copy()
		// Get all valid merges and compute their score.
		for i := 0; i < len(dfa.States); i++ {
			for j := i + 1; j < len(dfa.States); j++ {
				totalMerges++
				// If states are mergeable, calculate score and add to detMerges.
				if snapshot.MergeStates(dfa, i, j) {
					// Calculate score.
					score := scoringFunction(i, j, partition, snapshot, dfa)

					// If score is bigger than state pair with the highest score,
					// set current state pair to state pair with the highest score.
					if score > highestScoringStatePair.Score {
						highestScoringStatePair = StatePairScore{
							State1: i,
							State2: j,
							Score:  score,
						}
					}
				}
				// Undo merges.
				snapshot.RollbackChanges(partition)
			}
		}

		// Check if any deterministic merges were found.
		if highestScoringStatePair.Score != -1 {
			// Merge the state pairs with the highest score.
			partition.MergeStates(dfa, highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Convert the state partition to a DFA.
			valid := false
			valid, dfa = partition.ToDFA(dfa)

			// Panic if state partition to DFA conversion was unsuccessful.
			if !valid {
				panic("Invalid merged DFA.")
			}

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}
		} else {
			break
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return dfa
}

// WindowedSearch deterministically merges state pairs within a given window.
// Returns the resultant DFA when no more valid merges are possible.
func WindowedSearch(dfa DFA, windowSize int, windowGrow float64, scoringFunction ScoringFunction) DFA {
	start := time.Now()
	totalMerges := 0
	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}
	// Store ordered states within DFA.
	orderedStates := dfa.OrderedStates()

	// Iterate until stopped.
	for {
		// Window values.
		windowMin := 0
		windowMax := util.Min(windowSize, len(dfa.States))

		// Convert DFA to StatePartition for state merging.
		partition := dfa.ToStatePartition()
		// Copy the state partition for undoing merging.
		snapshot := partition.Copy()

		// Loop until no more deterministic merges are available within all possible windows.
		for {
			// Get all valid merges within window and compute their score.
			for i := 0; i < windowMax; i++ {
				for j := windowMin; j < windowMax; j++ {
					if i < j {
						totalMerges++
						// If states are mergeable, calculate score and add to detMerges.
						if snapshot.MergeStates(dfa, orderedStates[i], orderedStates[j]) {
							// Calculate score.
							score := scoringFunction(orderedStates[i], orderedStates[j], partition, snapshot, dfa)

							// If score is bigger than state pair with the highest score,
							// set current state pair to state pair with the highest score.
							if score > highestScoringStatePair.Score {
								highestScoringStatePair = StatePairScore{
									State1: orderedStates[i],
									State2: orderedStates[j],
									Score:  score,
								}
							}
						}
						// Undo merges.
						snapshot.RollbackChanges(partition)
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
				windowMax = util.Min(windowMax+windowSize, len(dfa.States))

				// If the window size is out of bounds, break loop and return the
				// most recent DFA found.
				if windowMin >= len(dfa.States) {
					break
				}
			}
		}

		if highestScoringStatePair.Score != -1 {
			// Merge the state pairs with the highest score.
			partition.MergeStates(dfa, highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Convert the state partition to a DFA.
			valid := false
			valid, dfa = partition.ToDFA(dfa)

			// Panic if state partition to DFA conversion was unsuccessful.
			if !valid {
				panic("Invalid merged DFA.")
			}

			// Update ordered states within new DFA.
			orderedStates = dfa.OrderedStates()

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}
		} else {
			break
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return dfa
}

// BlueFringeSearch deterministically merges possible state pairs within red-blue sets.
// Returns the resultant DFA when no more valid merges are possible.
func BlueFringeSearch(dfa DFA, scoringFunction ScoringFunction) DFA {
	start := time.Now()
	totalMerges := 0

	// State pair with the highest score.
	highestScoringStatePair := StatePairScore{-1, -1, -1}

	// Slice of state pairs to keep track of computed scores.
	scoresComputed := map[StateIDPair]util.Void{}

	// Initialize set of red states to starting state.
	red := map[int]util.Void{dfa.StartingStateID: util.Null}

	// Convert DFA to StatePartition for state merging.
	partition := dfa.ToStatePartition()
	// Copy the state partition for undoing merging.
	snapshot := partition.Copy()

	// Initialize merged flag to false.
	merged := false

	// Iterate until stopped.
	for {
		// Initialize set of blue states to empty set.
		blue := map[int]util.Void{}

		// Iterate over every red state.
		for element := range red {
			// Iterate over each symbol within DFA.
			for alphabetID := range dfa.Alphabet {
				// Store resultant stateID from red state.
				resultantStateID := dfa.States[element].Transitions[alphabetID]
				// If transition is valid and resultant state is not red,
				// add resultant state to blue set.
				if _, exists := red[resultantStateID]; resultantStateID != -1 && !exists {
					blue[resultantStateID] = util.Null
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
					if snapshot.MergeStates(dfa, blueElement, redElement) {
						// Set the state pairs score as computed.
						scoresComputed[StateIDPair{blueElement, redElement}] = util.Null
						scoresComputed[StateIDPair{redElement, blueElement}] = util.Null

						// Calculate score.
						score := scoringFunction(blueElement, redElement, partition, snapshot, dfa)

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
					snapshot.RollbackChanges(partition)
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
			partition.MergeStates(dfa, highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Convert the state partition to a DFA.
			valid := false
			valid, dfa = partition.ToDFA(dfa)

			// Convert DFA to StatePartition for state merging.
			partition = dfa.ToStatePartition()
			// Copy the state partition for undoing merging.
			snapshot = partition.Copy()

			// Panic if state partition to DFA conversion was unsuccessful.
			if !valid {
				panic("Invalid merged DFA.")
			}

			// Remove previous state pair with the highest score.
			highestScoringStatePair = StatePairScore{-1, -1, -1}

			// Slice of state pairs to keep track of computed scores.
			scoresComputed = map[StateIDPair]util.Void{}

			// Initialize set of red states to starting state.
			red = map[int]util.Void{dfa.StartingStateID: util.Null}
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return dfa
}
