package DFA_Toolkit

import (
	"fmt"
	"time"
)

// BenchmarkEDSM benchmarks the performance of the GreedyEDSM() function.
func BenchmarkEDSM(n int) {
	winners := 0
	for i := 0; i < n; i++ {
		start := time.Now()

		// Create a target DFA.
		target := AbbadingoDFA(32, true)

		// Training set.
		training, testing := AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA := GreedyEDSM(training)
		accuracy := resultantDFA.Accuracy(testing)

		if accuracy >= 0.99{
			winners++
		}

		fmt.Printf("BENCHMARK %d/%d. Duration: %.2fs. Resultant DFA: %d states, Accuracy: %.2f\n", i+1, n, time.Since(start).Seconds(), len(resultantDFA.States), accuracy)
	}
	overallAccuracy := float64(winners) / float64(n)
	fmt.Printf("Overall Accuracy: %.2f%%\n", overallAccuracy)
}
