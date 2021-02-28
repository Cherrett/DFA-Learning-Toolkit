package main

import (
	"DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	var timings []int64
	iterations := 5

	for i := 0; i < iterations; i++ {
		start := time.Now()

		AbbadingoDFA := DFA_Toolkit.AbbadingoDFA(20, true)
		AbbadingoDFA.Describe(false)
		fmt.Println("DFA Depth:", AbbadingoDFA.GetDepth())
		fmt.Println("DFA Loops:", AbbadingoDFA.LoopsCount())

		trainingDataset, testingDataset := DFA_Toolkit.AbbadingoDataset(AbbadingoDFA, 50, 0.2)

		trainingDatasetConsistentWithDFA := trainingDataset.ConsistentWithDFA(AbbadingoDFA)
		testingDatasetConsistentWithDFA := testingDataset.ConsistentWithDFA(AbbadingoDFA)

		if trainingDatasetConsistentWithDFA && testingDatasetConsistentWithDFA{
			fmt.Println("Both Consistent with AbbadingoDFA")
		}

		trainingDataset.WriteToAbbadingoFile("AbbadingoDatasets/customDataset1/training.a")
		testingDataset.WriteToAbbadingoFile("AbbadingoDatasets/customDataset1/testing.a")

		APTA := DFA_Toolkit.GetPTAFromDataset(trainingDataset, true)
		fmt.Println("APTA Depth:", APTA.GetDepth())
		fmt.Println("APTA Loops:", APTA.LoopsCount())

		timings = append(timings, time.Since(start).Milliseconds())
	}
	var sum int64
	for _, timing := range timings{
		sum += timing
	}
	fmt.Printf("Average Time: %vms\n", sum/int64(iterations))
}