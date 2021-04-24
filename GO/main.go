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
	resultantDFA, _ := dfatoolkit.RPNIFromDataset(trainingSet)

	fmt.Println(len(resultantDFA.States), resultantDFA.TransitionsCount())
}
