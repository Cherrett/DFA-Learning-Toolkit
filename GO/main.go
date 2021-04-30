package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"fmt"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	// Random Seed.
	// rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 50
	// Alphabet size.
	alphabetSize := 2
	// Training sparsity percentage.
	sparsityPercentage := 12.5

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, _ := dfatoolkit.DefaultStaminaInstance(alphabetSize, targetSize, sparsityPercentage)

		trainingSet.GetPTA(true)
	}
}
