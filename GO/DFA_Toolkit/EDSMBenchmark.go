package DFA_Toolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"time"
)

// BenchmarkEDSM benchmarks the performance of the GreedyEDSM() function.
func BenchmarkEDSM(n int) {
	winners := 0
	totalAccuracies := util.NewMinMaxAvg()
	totalNumberOfStates := util.NewMinMaxAvg()
	for i := 0; i < n; i++ {
		start := time.Now()

		// Create a target DFA.
		target := AbbadingoDFA(32, true)

		// Training set.
		training, testing := AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA := GreedyEDSM(training)
		accuracy := resultantDFA.Accuracy(testing)

		totalAccuracies.Add(accuracy)
		totalNumberOfStates.Add(float64(resultantDFA.AllStatesCount()))

		if accuracy >= 0.99{
			winners++
		}

		fmt.Printf("BENCHMARK %d/%d. Duration: %.2fs.\n", i+1, n, time.Since(start).Seconds())
	}
	fmt.Printf("Percentage of 0.99+ Accuracy: %.2f%%\n", float64(winners) / float64(n))
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuracies.Min(), totalAccuracies.Max(), totalAccuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStates.Min(), totalNumberOfStates.Max(), totalNumberOfStates.Avg())
}
