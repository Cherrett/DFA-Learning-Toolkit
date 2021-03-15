package DFA_Toolkit

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// StatePairScore struct to store state pairs and their merge score.
type StatePairScore struct {
	state1 int	// StateID for first state.
	state2 int  // StateID for second state.
	score int	// Score of merge for given states.
}

// ScoringFunction takes a state partition as input and returns an integer score.
type ScoringFunction func(partition StatePartition) int

// HighestScoringMerge returns the highest scoring state pair. If more than
// one state pair have the highest score, one is chosen randomly.
func HighestScoringMerge(statePairScores []StatePairScore, randomFromBest bool) StatePairScore{
	// Sort the state pairs by score.
	sort.Slice(statePairScores, func(i, j int) bool {
		return statePairScores[i].score > statePairScores[j].score
	})

	if !randomFromBest{
		return statePairScores[0]
	}

	// Declare the highest score from the first element within slice (since slice is sorted).
	highestScore := statePairScores[0].score

	// Declare slice to store the highest scoring state pairs.
	highestScorePairs := []StatePairScore{}

	// Iterate over each state pair.
	for i := range statePairScores {
		// If the score of the state pair is equal to the highest score, add it to the highest scoring state pairs.
		if statePairScores[i].score == highestScore{
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

// Traverses a path within lattice of automaton given a DFA (generally an APTA) and a
// scoring function. Returns the resultant DFA when no more valid merges are possible.
func GreedyPath(dfa DFA, scoringFunction ScoringFunction, randomFromBest bool) DFA{
	// Convert DFA (APTA) to StatePartition for state merging.
	partition := dfa.ToStatePartition()
	// Copy the state partition for undoing merging.
	snapshot := partition.Copy()
	// Slice of state pairs with score.
	var detMerges []StatePairScore
	start := time.Now()
	totalMerges := 0

	// Loop until no more deterministic merges are available.
	for{
		// Get all valid merges and compute their score.
		for i := 0; i < len(dfa.States); i++ {
			for j := i + 1; j < len(dfa.States); j++ {
				totalMerges ++
				// If states are mergeable, calculate score and add to detMerges.
				if snapshot.MergeStates(dfa, i, j){
					detMerges = append(detMerges, StatePairScore{
						state1: i,
						state2: j,
						score: scoringFunction(snapshot),
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
			partition.MergeStates(dfa, highestScoringStatePair.state1, highestScoringStatePair.state2)

			// Convert the state partition to a DFA.
			valid := false
			valid, dfa = partition.ToDFA(dfa)

			// Panic if state partition to DFA conversion was unsuccessful.
			if !valid{
				panic("Invalid merged DFA.")
			}

			// Convert new DFA to StatePartition for state merging.
			partition = dfa.ToStatePartition()
			// Copy the state partition for undoing merging.
			snapshot = partition.Copy()
			// Remove previous state pairs' score.
			detMerges = nil
		}else{
			totalTime := (time.Now()).Sub(start).Seconds()
			fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

			// Return the final resultant DFA.
			return dfa
		}
	}

	//// Loop until no more deterministic merges are available.
	//for len(detMerges) != 0{
	//
	//	highestScoringStatePair := HighestScoringMerge(detMerges)
	//
	//	// Merge the state pairs with the highest score.
	//	partition.MergeStates(newDFA, highestScoringStatePair.state1, highestScoringStatePair.state2)
	//
	//	// Convert the state partition to a DFA.
	//	valid := false
	//	valid, newDFA = partition.ToDFA(newDFA)
	//
	//	// Panic if state partition to DFA conversion was unsuccessful.
	//	if !valid{
	//		panic("Invalid merged DFA.")
	//	}
	//
	//	// Convert new DFA to StatePartition for state merging.
	//	partition = newDFA.ToStatePartition()
	//	// Copy the state partition for undoing merging.
	//	snapshot = partition.Copy()
	//	// Remove previous state pairs' score
	//	detMerges = nil
	//
	//	totalMerges = 0
	//
	//	// Get all valid merges and compute their score.
	//	for i := 0; i < len(newDFA.States); i++ {
	//		for j := i + 1; j < len(newDFA.States); j++ {
	//			totalMerges ++
	//			// If states are mergeable, calculate score and add to detMerges.
	//			if snapshot.MergeStates(newDFA, i, j){
	//				detMerges = append(detMerges, StatePairScore{
	//					state1: i,
	//					state2: j,
	//					score: scoringFunction(snapshot),
	//				})
	//			}
	//			// Undo merges.
	//			snapshot.RollbackChanges(partition)
	//		}
	//	}
	//}
	//
	//totalTime := (time.Now()).Sub(start).Seconds()
	//fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)
	//
	//// Convert the state partition to a DFA.
	//valid, resultantDFA := partition.ToDFA(newDFA)
	//
	//// Panic if state partition to DFA conversion was unsuccessful.
	//if !valid{
	//	panic("EDSM resultant DFA contains invalid merges.")
	//}
	//
	//// Return the final resultant DFA.
	//return resultantDFA
}
