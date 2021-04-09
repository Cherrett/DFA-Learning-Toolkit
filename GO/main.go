package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	// These are target DFA sizes we will test.
	dfaSizes := []int{16, 32, 64}
	// These are the training set sizes we will test.
	trainingSetSizes := []int{230, 607, 1521}

	// Benchmark over the problem instances.
	for i := range dfaSizes {
		targetSize := dfaSizes[i]
		trainingSetSize := trainingSetSizes[i]

		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("BENCHMARK %d (Target: %d states, Training: %d strings\n", i+1, targetSize, trainingSetSize)
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")

		// Create APTA.
		apta, valid := dfatoolkit.DFAFromJSON(fmt.Sprintf("TestingAPTAs/%d.json", targetSize))

		if !valid{
			panic("APTA read not valid.")
		}

		fmt.Printf("APTA size: %d\n", len(apta.States))

		// Perform all the merges.
		part := apta.ToStatePartition()
		snapshot := part.Copy()
		totalMerges := 0
		validMerges := 0
		start := time.Now()

		for i := 0; i < len(apta.States); i++ {
			for j := i + 1; j < len(apta.States); j++ {
				totalMerges++
				if snapshot.MergeStates(apta, i, j) {
					validMerges++
					//snapshot.LabelledBlocksCount(apta)
					//snapshot.BlocksCount()
					//print(temp, temp2)
				}

				snapshot.RollbackChanges(part)
			}
		}

		totalTime := (time.Now()).Sub(start).Seconds()
		fmt.Printf("Total merges: %d\n", totalMerges)
		fmt.Printf("Valid merges: %d\n", validMerges)
		fmt.Printf("Time: %.2fs\n", totalTime)
		fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)
	}
}
