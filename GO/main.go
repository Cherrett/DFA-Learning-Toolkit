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

	// random seed
	rand.Seed(time.Now().UnixNano())

	// Training set.
	training := DFA_Toolkit.GetDatasetFromAbbadingoFile("../AbbadingoDatasets/test.txt")

	for i:=0; i < 5; i++{
		resultantDFA := DFA_Toolkit.GreedyEDSM(training, false)
		resultantDFA.ToJPG(fmt.Sprintf("../Temp/testing%d.jpg", i), true, true)
	}
}