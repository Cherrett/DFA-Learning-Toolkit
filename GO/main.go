package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"fmt"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	trainingSet := dfatoolkit.GetDatasetFromStaminaFile("Datasets/Stamina/31/31_training.txt")
	// Construct an APTA from the dataset.
	APTA := trainingSet.GetPTA(true)
	APTA = APTA.SetOrderAsID()
	resultantDFA, searchData := dfatoolkit.RPNI(APTA)
	
	fmt.Println(len(resultantDFA.States), resultantDFA.TransitionsCount())
	fmt.Println(searchData.Duration)
}
