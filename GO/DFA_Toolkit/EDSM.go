package DFA_Toolkit

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type StatePairScore struct {
	state1 int
	state2 int
	score int
}

// HighestScoringMerge returns the highest scoring state pair. If more than
// one state pair have the highest score, one is chosen randomly.
func HighestScoringMerge(statePairScores []StatePairScore) StatePairScore{
	// Declare slice to store the highest scoring state pairs.
	highestScorePairs := []StatePairScore{}
	// Sort the state pairs by score.
	sort.Slice(statePairScores, func(i, j int) bool {
		return statePairScores[i].score > statePairScores[j].score
	})
	
	// Declare the highest score from the first element within slice (since slice is sorted).
	highestScore := statePairScores[0].score

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

// GreedyEDSM is a greedy version of Evidence Driven State-Merging
func GreedyEDSM(dataset Dataset) DFA{
	LengthOfDataset := len(dataset)
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)
	// Convert APTA to StatePartition for state merging.
	partition := APTA.ToStatePartition()
	// Copy the state partition for undoing merging.
	snapshot := partition.Copy()
	// Slice of state pairs with score.
	var detMerges []StatePairScore
	start := time.Now()
	totalMerges := 0

	// Get all valid merges and compute their score.
	for i := 0; i < len(APTA.States); i++ {
		for j := i + 1; j < len(APTA.States); j++ {
			totalMerges ++
			// If states are mergeable, calculate score and add to detMerges.
			if snapshot.MergeStates(APTA, i, j){
				detMerges = append(detMerges, StatePairScore{
					state1: i,
					state2: j,
					score: LengthOfDataset - snapshot.EDSMScore(),
				})
			}
			// Undo merges.
			snapshot.RollbackChanges(partition)
		}
	}
	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)

	// Set newDFA to APTA.
	newDFA := APTA

	// Loop until no more deterministic merges are available.
	for len(detMerges) != 0{

		highestScoringStatePair := HighestScoringMerge(detMerges)

		// Merge the state pairs with the highest score.
		partition.MergeStates(newDFA, highestScoringStatePair.state1, highestScoringStatePair.state2)

		// Convert the state partition to a DFA.
		valid := false
		valid, newDFA = partition.ToDFA(newDFA)

		// Panic if state partition to DFA conversion was unsuccessful.
		if !valid{
			panic("Invalid merged DFA.")
		}

		// Convert new DFA to StatePartition for state merging.
		partition = newDFA.ToStatePartition()
		// Copy the state partition for undoing merging.
		snapshot = partition.Copy()
		// Remove previous state pairs' score
		detMerges = nil

		totalMerges = 0

		// Get all valid merges and compute their score.
		for i := 0; i < len(newDFA.States); i++ {
			for j := i + 1; j < len(newDFA.States); j++ {
				totalMerges ++
				// If states are mergeable, calculate score and add to detMerges.
				if snapshot.MergeStates(newDFA, i, j){
					detMerges = append(detMerges, StatePairScore{
						state1: i,
						state2: j,
						score: LengthOfDataset - snapshot.EDSMScore(),
					})
				}
				// Undo merges.
				snapshot.RollbackChanges(partition)
			}
		}
	}

	// Convert the state partition to a DFA.
	valid, resultantDFA := partition.ToDFA(newDFA)

	// Panic if state partition to DFA conversion was unsuccessful.
	if !valid{
		panic("EDSM resultant DFA contains invalid merges.")
	}

	// Return the final resultant DFA.
	return resultantDFA
}