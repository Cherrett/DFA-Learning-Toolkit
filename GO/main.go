package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	dfa := dfatoolkit.StaminaDFA(5, 50)
	dfa.Describe(false)
	//dfa.ToJPG("temp.jpg", nil, false, false)
	// Construct an APTA from the dataset.
	//APTA := trainingSet.GetPTA(true)
	//APTA = APTA.SetOrderAsID()
	//resultantDFA, searchData := dfatoolkit.RPNI(APTA)

	//fmt.Println(len(resultantDFA.States), resultantDFA.TransitionsCount())
	//fmt.Println(searchData.Duration)
}
