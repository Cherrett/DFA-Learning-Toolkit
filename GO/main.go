package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"time"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	winners := 0
	totalAccuracies := util.NewMinMaxAvg()
	totalNumberOfStates := util.NewMinMaxAvg()
	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)
		start := time.Now()

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := dfatoolkit.AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA := dfatoolkit.WindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		accuracy := resultantDFA.Accuracy(testingSet)

		totalAccuracies.Add(accuracy)
		totalNumberOfStates.Add(float64(len(resultantDFA.States)))

		if accuracy >= 0.99 {
			winners++
		}

		fmt.Printf("Duration: %.2fs\n\n", time.Since(start).Seconds())
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuracies.Min(), totalAccuracies.Max(), totalAccuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStates.Min(), totalNumberOfStates.Max(), totalNumberOfStates.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")
}
