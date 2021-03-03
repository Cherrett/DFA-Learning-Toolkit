package DFA_Toolkit

import (
	"fmt"
	"time"
)

// BenchmarkDetMerge benchmarks the performance of the DetMerge() operation
// in a state partition.
func BenchmarkDetMerge() {
	// These are target DFA sizes we will test.
	//dfaSizes := []int{32, 64, 128}
	dfaSizes := []int{32, 64, 128}
	// These are the training set sizes we will test.
	//trainingSetSizes := []int{607, 1521, 4382}
	trainingSetSizes := []int{607, 1521, 4382}

	// Benchmark over the problem instances.
	for i := range dfaSizes {
		targetSize := dfaSizes[i]
		trainingSetSize := trainingSetSizes[i]

		// Create a target DFA.
		target := AbbadingoDFA(targetSize, true)

		// Training set.
		training, _ := AbbadingoDatasetExact(target, trainingSetSize, 0)

		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("BENCHMARK %d (Target: %d states, Training: %d strings\n", i+1, targetSize, len(training))
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")

		// Info about training set.
		fmt.Printf("Training proportion +ve: %.2f%%\n", training.AcceptingStringInstancesRatio()*100.0)
		fmt.Printf("Training proportion -ve: %.2f%%\n", training.RejectingStringInstancesRatio()*100.0)

		// Create APTA.
		apta := training.GetPTA(true)
		//apta := training.APTA(target.AlphabetSize)
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
				if snapshot.MergeStates(apta, i, j){
					validMerges++
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