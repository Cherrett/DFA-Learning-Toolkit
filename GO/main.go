package main

import (
	"DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	var timings []int64
	iterations := 1

	for i := 0; i < iterations; i++ {
		start := time.Now()

		AbbadingoDFA := DFA_Toolkit.AbbadingoDFA(5, true)
		AbbadingoDFA.Describe(false)
		fmt.Println("DFA Depth:", AbbadingoDFA.Depth())
		fmt.Println("DFA Loops:", AbbadingoDFA.LoopsCount())

		trainingDataset, testingDataset := DFA_Toolkit.AbbadingoDataset(AbbadingoDFA, 75, 0.2)

		trainingDatasetConsistentWithDFA := trainingDataset.ConsistentWithDFA(AbbadingoDFA)
		testingDatasetConsistentWithDFA := testingDataset.ConsistentWithDFA(AbbadingoDFA)

		if trainingDatasetConsistentWithDFA && testingDatasetConsistentWithDFA{
			fmt.Println("Both Consistent with AbbadingoDFA")
		}

		timings = append(timings, time.Since(start).Milliseconds())
	}
	var sum int64
	for _, timing := range timings{
		sum += timing
	}
	fmt.Printf("Average Time: %vms\n", sum/int64(iterations))
}