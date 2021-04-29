package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"fmt"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	dfa := dfatoolkit.StaminaDFA(10, 50)
	dfa.Describe(false)
	training, testing := dfatoolkit.StaminaDataset(dfa, 12.5, 20000, 1500)

	fmt.Println(training.Count())
	fmt.Println(training.AcceptingStringInstancesCount())
	fmt.Println(training.AcceptingStringInstancesRatio())
	fmt.Println(training.RejectingStringInstancesCount())
	fmt.Println(training.RejectingStringInstancesRatio())
	fmt.Println()
	fmt.Println(testing.Count())
	fmt.Println(testing.AcceptingStringInstancesCount())
	fmt.Println(testing.AcceptingStringInstancesRatio())
	fmt.Println(testing.RejectingStringInstancesCount())
	fmt.Println(testing.RejectingStringInstancesRatio())
}
