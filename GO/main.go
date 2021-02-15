package main

import (
	"./DFA"
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	listOfStrings := DFA.GetListOfStringInstancesFromFile("dataset4/train.a")
	APTA := DFA.GetPTAFromListOfStringInstances(listOfStrings, true)
	APTA.Describe(false)

	fmt.Println("DFA Depth:", APTA.Depth())

	//APTA.AddState(DFA.UNKNOWN)
	fmt.Println(DFA.ListOfStringInstancesConsistentWithDFA(listOfStrings, APTA))

	elapsed := time.Since(start)
	fmt.Printf("Time: %s\n", elapsed)
}
