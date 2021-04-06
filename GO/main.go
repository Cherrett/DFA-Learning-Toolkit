package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	// random seed
	// rand.Seed(time.Now().UnixNano())

	structurallyCompleteCount := 0
	iterations := 100
	start := time.Now()

	for i := 0; i < iterations; i++ {
		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(32, true)

		//target.ToJPG("temp.jpg", false, true)

		// Training set.
		training, _ := dfatoolkit.AbbadingoDataset(target, 100, 0)

		if training.SymmetricallyStructurallyComplete(target) {
			structurallyCompleteCount++
		}

		fmt.Printf("Iteration %d/%d\n", i+1, iterations)
	}

	fmt.Printf("Percentage which were Structurally Complete: %.4f\n", float64(structurallyCompleteCount)/float64(iterations))
	totalTime := (time.Now()).Sub(start).Seconds()
	fmt.Printf("Total Time: %.2f seconds.\n", totalTime)
}
