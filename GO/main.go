package main

import (
	"DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// random seed
	rand.Seed(time.Now().UnixNano())

	//DFA_Toolkit.BenchmarkDetMerge()

	var timings []int64
	iterations := 1

	for i := 0; i < iterations; i++ {
		start := time.Now()

		//AbbadingoDFA := DFA_Toolkit.AbbadingoDFA(20, true)
		//AbbadingoDFA.Describe(false)
		//fmt.Println("DFA Depth:", AbbadingoDFA.GetDepth())
		//fmt.Println("DFA Loops:", AbbadingoDFA.LoopsCount())
		//
		//trainingDataset, testingDataset := DFA_Toolkit.AbbadingoDataset(AbbadingoDFA, 50, 0.2)
		//
		//trainingDatasetConsistentWithDFA := trainingDataset.ConsistentWithDFA(AbbadingoDFA)
		//testingDatasetConsistentWithDFA := testingDataset.ConsistentWithDFA(AbbadingoDFA)
		//
		//if trainingDatasetConsistentWithDFA && testingDatasetConsistentWithDFA{
		//	fmt.Println("Both Consistent with AbbadingoDFA")
		//}
		//
		//trainingDataset.WriteToAbbadingoFile("AbbadingoDatasets/customDataset1/training.a")
		//testingDataset.WriteToAbbadingoFile("AbbadingoDatasets/customDataset1/testing.a")

		dataset := DFA_Toolkit.GetDatasetFromAbbadingoFile("./AbbadingoDatasets/dataset1/train.a")

		APTA := dataset.GetPTA(false)
		fmt.Println("APTA Depth:", APTA.GetDepth())
		fmt.Println("APTA Loops:", APTA.LoopsCount())

		APTA.Describe(true)

		statePartition := APTA.ToStatePartition()

		valid := statePartition.MergeStates(APTA, 2, 4)
		fmt.Println(valid)
		valid, mergedDFA := statePartition.ToDFA(APTA)

		if !valid{
			fmt.Printf("Error!")
		}else{
			fmt.Printf("\nMerged DFA Below\n")
			mergedDFA.Describe(true)
		}

		timings = append(timings, time.Since(start).Milliseconds())
	}
	var sum int64
	for _, timing := range timings{
		sum += timing
	}
	fmt.Printf("Average Time: %vms\n", sum/int64(iterations))
}