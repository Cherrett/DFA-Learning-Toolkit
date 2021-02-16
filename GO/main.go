package main

import (
	"./DFA"
	"fmt"
	"time"
)

func main() {
	var timings []int64
	iterations := 5

	for i := 0; i < iterations; i++ {
		start := time.Now()

		listOfStrings := DFA.GetListOfStringInstancesFromFile("dataset4/train.a")
		APTA := DFA.GetPTAFromListOfStringInstances(listOfStrings, true)
		APTA.Describe(false)

		fmt.Println("DFA Depth:", APTA.Depth())

		//APTA.AddState(DFA.UNKNOWN)
		fmt.Println(DFA.ListOfStringInstancesConsistentWithDFA(listOfStrings, APTA))

		timings = append(timings, time.Since(start).Milliseconds())
	}
	var sum int64
	for _, timing := range timings{
		sum += timing
	}
	fmt.Printf("Average Time: %vms\n", sum/int64(iterations))
}
