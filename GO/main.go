package main

import (
	"./DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	var timings []int64
	iterations := 1

	for i := 0; i < iterations; i++ {
		start := time.Now()

		listOfStrings := DFA_Toolkit.GetListOfStringInstancesFromFile("dataset1/train.a")
		APTA := DFA_Toolkit.GetPTAFromListOfStringInstances(listOfStrings, true)
		APTA.Describe(false)

		fmt.Println("DFA_Toolkit Depth:", APTA.Depth())

		//APTA.AddState(DFA_Toolkit.UNKNOWN)
		fmt.Println(DFA_Toolkit.ListOfStringInstancesConsistentWithDFA(listOfStrings, APTA))

		result := DFA_Toolkit.RPNI(DFA_Toolkit.GetAcceptingStringInstances(listOfStrings),
			DFA_Toolkit.GetRejectingStringInstances(listOfStrings))

		result.Describe(false)
		timings = append(timings, time.Since(start).Milliseconds())
	}
	var sum int64
	for _, timing := range timings{
		sum += timing
	}
	fmt.Printf("Average Time: %vms\n", sum/int64(iterations))
}
