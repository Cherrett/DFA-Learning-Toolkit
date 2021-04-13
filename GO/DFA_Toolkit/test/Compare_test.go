package dfatoolkit

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// TestBenchmarkDetMergeCompare benchmarks the performance of two MergeStates functions.
func TestBenchmarkDetMergeCompare(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// These are target DFA sizes we will test.
	dfaSizes := []int{16, 32, 64}
	// These are the training set sizes we will test.
	trainingSetSizes := []int{230, 607, 1521}

	// Benchmark over the problem instances.
	for iterator := range dfaSizes {
		targetSize := dfaSizes[iterator]
		trainingSetSize := trainingSetSizes[iterator]

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training set.
		training, _ := dfatoolkit.AbbadingoDatasetExact(target, trainingSetSize, 0)

		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("BENCHMARK %d (Target: %d states, Training: %d strings\n", iterator+1, targetSize, len(training))
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
				if snapshot.MergeStates(apta, i, j) {
					validMerges++
				}

				snapshot.RollbackChanges(part)
			}
		}

		totalTime := (time.Now()).Sub(start).Seconds()
		fmt.Printf("MergeStates Total merges: %d\n", totalMerges)
		fmt.Printf("MergeStates Valid merges: %d\n", validMerges)
		fmt.Printf("MergeStates Time: %.2fs\n", totalTime)
		fmt.Printf("MergeStates Merges per second: %.2f\n", float64(totalMerges)/totalTime)

		totalMerges = 0
		validMerges = 0
		start = time.Now()

		for i := 0; i < len(apta.States); i++ {
			for j := i + 1; j < len(apta.States); j++ {
				totalMerges++
				if snapshot.MergeStates(apta, i, j) {
					validMerges++
				}

				snapshot.RollbackChanges(part)
			}
		}

		totalTime = (time.Now()).Sub(start).Seconds()
		fmt.Printf("MergeStates2 Total merges: %d\n", totalMerges)
		fmt.Printf("MergeStates2 Valid merges: %d\n", validMerges)
		fmt.Printf("MergeStates2 Time: %.2fs\n", totalTime)
		fmt.Printf("MergeStates2 Merges per second: %.2f\n", float64(totalMerges)/totalTime)
	}
}