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

	// Training set.
	training := DFA_Toolkit.GetDatasetFromAbbadingoFile("AbbadingoDatasets/dataset4/train.a")
	training.SortDatasetByLength()
	valid := training.ToJSON("AbbadingoDatasets/dataset4/train.json")

	if !valid{
		fmt.Println("Error")
	}

	//for i:=0; i < 5; i++{
	//	resultantDFA := DFA_Toolkit.GreedyEDSM(training, false)
	//	resultantDFA.ToJPG(fmt.Sprintf("../Temp/testing%d.jpg", i), true, true)
	//}
}