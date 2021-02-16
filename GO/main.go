package main

import (
	"./DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	var timings []int64
	iterations := 5

	for i := 0; i < iterations; i++ {
		start := time.Now()

		listOfStrings := DFA_Toolkit.GetListOfStringInstancesFromFile("dataset4/train.a")
		APTA := DFA_Toolkit.GetPTAFromListOfStringInstances(listOfStrings, true)
		APTA.Describe(false)

		fmt.Println("DFA_Toolkit Depth:", APTA.Depth())

		//APTA.AddState(DFA_Toolkit.UNKNOWN)
		fmt.Println(DFA_Toolkit.ListOfStringInstancesConsistentWithDFA(listOfStrings, APTA))

		timings = append(timings, time.Since(start).Milliseconds())
	}
	var sum int64
	for _, timing := range timings{
		sum += timing
	}
	fmt.Printf("Average Time: %vms\n", sum/int64(iterations))
}
