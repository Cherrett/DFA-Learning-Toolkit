package DFA_Toolkit

import (
	"fmt"
	"sort"
	"time"
)

type StatePairScore struct {
	state1 int
	state2 int
	score int
}

// Greedy Version of Evidence Driven State-Merging
func GreedyEDSM(dataset Dataset) DFA{
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
					score: snapshot.EDSMScore(),
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
		// Sort the deterministic merges by score.
		sort.Slice(detMerges, func(i, j int) bool {
			return detMerges[i].score > detMerges[j].score
		})
		// Merge the state pairs with the highest score.
		partition.MergeStates(newDFA, detMerges[0].state1, detMerges[0].state2)

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

		start = time.Now()
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
						score: snapshot.EDSMScore(),
					})
				}
				// Undo merges.
				snapshot.RollbackChanges(partition)
			}
		}
		totalTime := (time.Now()).Sub(start).Seconds()
		fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)
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