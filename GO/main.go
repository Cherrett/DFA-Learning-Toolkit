package main

import (
	"DFA_Toolkit/DFA_Toolkit"
	"math/rand"
	"time"
)

func main() {
	// random seed
	rand.Seed(time.Now().UnixNano())

	// Training set.
	training := DFA_Toolkit.GetDatasetFromAbbadingoFile("AbbadingoDatasets/test.txt")
	resultantDFA := DFA_Toolkit.BlueFringeEDSM(training, false)
	resultantDFA.ToJPG("testing.jpg", true, false)
}