package DFA_Toolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// StatePairScore struct to store state pairs and their merge score.
type StatePairScore struct {
	State1 int // StateID for first state.
	State2 int // StateID for second state.
	Score  int // Score of merge for given states.
}

// ScoringFunction takes a state partition as input and returns an integer score.
type ScoringFunction func(partition StatePartition) int

// HighestScoringMerge returns the highest scoring state pair. If more than
// one state pair have the highest score, one is chosen randomly.
func HighestScoringMerge(statePairScores []StatePairScore, randomFromBest bool) StatePairScore{
	// Sort the state pairs by score.
	sort.Slice(statePairScores, func(i, j int) bool {
		return statePairScores[i].Score > statePairScores[j].Score
	})

	if !randomFromBest{
		return statePairScores[0]
	}

	// Declare the highest score from the first element within slice (since slice is sorted).
	highestScore := statePairScores[0].Score

	// Declare slice to store the highest scoring state pairs.
	highestScorePairs := []StatePairScore{}

	// Iterate over each state pair.
	for i := range statePairScores {
		// If the score of the state pair is equal to the highest score, add it to the highest scoring state pairs.
		if statePairScores[i].Score == highestScore{
			highestScorePairs = append(highestScorePairs, statePairScores[i])
			// Else, break out of loop (since slice is sorted, score cannot be bigger in other pairs).
		}else{
			break
		}
	}

	// If only one highest scoring state pair exists, return it.
	if len(highestScorePairs) == 1{
		return highestScorePairs[0]
		// Else, return a random state pair from highest scoring pairs.
	}else{
		return highestScorePairs[rand.Intn(len(highestScorePairs))]
	}
}

// Deterministically merges all possible state pairs.
// Returns the resultant DFA when no more valid merges are possible.
func GreedySearch(dfa DFA, scoringFunction ScoringFunction, randomFromBest bool) DFA{
	// Slice of state pairs with score.
	var detMerges []StatePairScore
	start := time.Now()
	totalMerges := 0

	// Loop until no more deterministic merges are available.
	for{
		// Convert DFA to StatePartition for state merging.
		partition := dfa.ToStatePartition()
		// Copy the state partition for undoing merging.
		snapshot := partition.Copy()
		// Get all valid merges and compute their score.
		for i := 0; i < len(dfa.States); i++ {
			for j := i + 1; j < len(dfa.States); j++ {
				totalMerges ++
				// If states are mergeable, calculate score and add to detMerges.
				if snapshot.MergeStates(dfa, i, j){
					detMerges = append(detMerges, StatePairScore{
						State1: i,
						State2: j,
						Score:  scoringFunction(snapshot),
					})
				}
				// Undo merges.
				snapshot.RollbackChanges(partition)
			}
		}

		// Check if any deterministic merges were found.
		if len(detMerges) > 0{
			highestScoringStatePair := HighestScoringMerge(detMerges, randomFromBest)

			// Merge the state pairs with the highest score.
			partition.MergeStates(dfa, highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Convert the state partition to a DFA.
			valid := false
			valid, dfa = partition.ToDFA(dfa)

			// Panic if state partition to DFA conversion was unsuccessful.
			if !valid{
				panic("Invalid merged DFA.")
			}

			// Remove previous state pairs' score.
			detMerges = nil
		}else{
			break
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return dfa
}

// Deterministically merges state pairs within a given window.
// Returns the resultant DFA when no more valid merges are possible.
func WindowedSearch(dfa DFA, windowSize int, windowGrow float64, scoringFunction ScoringFunction, randomFromBest bool) DFA{
	start := time.Now()
	totalMerges := 0
	// Slice of state pairs with score.
	var detMerges []StatePairScore
	// Store ordered states within DFA.
	orderedStates := dfa.OrderedStates()
	for{
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
							detMerges = append(detMerges, StatePairScore{
								State1: orderedStates[i],
								State2: orderedStates[j],
								Score:  scoringFunction(snapshot),
							})
						}
						// Undo merges.
						snapshot.RollbackChanges(partition)
					}
				}
			}

			// Check if any deterministic merges were found.
			if len(detMerges) > 0 {
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

		if len(detMerges) > 0{
			highestScoringStatePair := HighestScoringMerge(detMerges, randomFromBest)

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

			// Remove previous state pairs' score.
			detMerges = nil
		}else{
			break
		}
	}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return dfa
}

// Deterministically merges possible state pairs within red-blue sets.
// Returns the resultant DFA when no more valid merges are possible.
func BlueFringeSearch(dfa DFA, scoringFunction ScoringFunction, randomFromBest bool) DFA{
	dfa.GetDepth()

	// Slice of state pairs with score.
	var detMerges []StatePairScore
	start := time.Now()
	totalMerges := 0

	red := map[int]bool{dfa.StartingStateID: true}

	// Convert DFA to StatePartition for state merging.
	partition := dfa.ToStatePartition()
	// Copy the state partition for undoing merging.
	snapshot := partition.Copy()

	merged := false

	for{
		scoresComputed := make([][]bool, len(dfa.States))

		for i := range scoresComputed{
			scoresComputed[i] = make([]bool, len(dfa.States))
			for j := range dfa.States{
				scoresComputed[i][j] = false
			}
		}

		if merged{
			red = map[int]bool{dfa.StartingStateID: true}

			// Convert DFA to StatePartition for state merging.
			partition = dfa.ToStatePartition()
			// Copy the state partition for undoing merging.
			snapshot = partition.Copy()
		}

		blue := map[int]bool{}

		for element := range red{
			for symbolID:=0; symbolID < len(dfa.SymbolMap); symbolID++{
				resultantStateID := dfa.States[element].Transitions[symbolID]
				if resultantStateID != -1 && !red[resultantStateID]{
					blue[resultantStateID] = true
				}
			}
		}

		if len(blue) == 0{
			break
		}

		for blueElement := range blue{
			merged = false
			for redElement := range red{
				if scoresComputed[blueElement][redElement]{
					merged = true
				}else{
					totalMerges ++
					if blueElement != redElement && snapshot.MergeStates(dfa, blueElement, redElement){
						scoresComputed[blueElement][redElement] = true
						scoresComputed[redElement][blueElement] = true

						detMerges = append(detMerges, StatePairScore{
							State1: blueElement,
							State2: redElement,
							Score:  (100 * scoringFunction(snapshot)) + 99 - dfa.States[redElement].Depth,
						})

						merged = true
					}
					// Undo merges.
					snapshot.RollbackChanges(partition)
				}
			}

			if !merged{
				red[blueElement] = true
				break
			}
		}

		if merged{
			highestScoringStatePair := HighestScoringMerge(detMerges, randomFromBest)

			// Merge the state pairs with the highest score.
			partition.MergeStates(dfa, highestScoringStatePair.State1, highestScoringStatePair.State2)

			// Convert the state partition to a DFA.
			valid := false
			valid, dfa = partition.ToDFA(dfa)

			// Panic if state partition to DFA conversion was unsuccessful.
			if !valid{
				panic("Invalid merged DFA.")
			}

			dfa.GetDepth()

			// Remove previous state pairs' score.
			detMerges = nil
		}

		//blue = map[int]bool{}
		//
		//for element := range red{
		//	for symbolID:=0; symbolID < len(dfa.SymbolMap); symbolID++{
		//		resultantStateID := dfa.States[element].Transitions[symbolID]
		//		if resultantStateID != -1 && !red[resultantStateID]{
		//			blue[resultantStateID] = true
		//		}
		//	}
		//}
		//
		//if len(blue) == 0{
		//	break
		//}
	}


	//// Loop until no more deterministic merges are available.
	//for{
	//	// Convert DFA to StatePartition for state merging.
	//	partition := dfa.ToStatePartition()
	//	// Copy the state partition for undoing merging.
	//	snapshot := partition.Copy()
	//	// Get all valid merges and compute their score.
	//	for i := 0; i < len(dfa.States); i++ {
	//		for j := i + 1; j < len(dfa.States); j++ {
	//			totalMerges ++
	//			// If states are mergeable, calculate score and add to detMerges.
	//			if snapshot.MergeStates(dfa, i, j){
	//				detMerges = append(detMerges, StatePairScore{
	//					State1: i,
	//					State2: j,
	//					Score:  scoringFunction(snapshot),
	//				})
	//			}
	//			// Undo merges.
	//			snapshot.RollbackChanges(partition)
	//		}
	//	}
	//
	//	// Check if any deterministic merges were found.
	//	if len(detMerges) > 0{
	//		highestScoringStatePair := HighestScoringMerge(detMerges, randomFromBest)
	//
	//		// Merge the state pairs with the highest score.
	//		partition.MergeStates(dfa, highestScoringStatePair.State1, highestScoringStatePair.State2)
	//
	//		// Convert the state partition to a DFA.
	//		valid := false
	//		valid, dfa = partition.ToDFA(dfa)
	//
	//		// Panic if state partition to DFA conversion was unsuccessful.
	//		if !valid{
	//			panic("Invalid merged DFA.")
	//		}
	//
	//		// Remove previous state pairs' score.
	//		detMerges = nil
	//	}else{
	//		break
	//	}
	//}

	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Return the final resultant DFA.
	return dfa
}