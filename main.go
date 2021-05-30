package main

import (
	"fmt"
	dfalearningtoolkit "github.com/Cherrett/DFA-Learning-Toolkit/core"
	"math"
	"math/rand"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	// Random Seed.
	// rand.Seed(time.Now().UnixNano())
	rand.Seed(0)

	//dfa := dfatoolkit.NewDFA()
	//dfa.AddSymbol()
	//dfa.AddSymbol()
	//
	//for i := 0; i < 8; i++ {
	//	dfa.AddState(dfatoolkit.ACCEPTING)
	//}
	//
	//dfa.AddTransition()

	//_, trainingDataset, _ := dfatoolkit.DefaultStaminaInstance(2, 32, 10.0)

	trainingDataset := dfalearningtoolkit.GetDatasetFromAbbadingoFile("datasets/GI_learning_datasets/random_100_100_100.txt")

	fmt.Println("Positive Strings:", trainingDataset.AcceptingStringInstancesCount())
	fmt.Println("Negative Strings:", trainingDataset.RejectingStringInstancesCount())

	result, mergeData := dfalearningtoolkit.RPNIFromDataset(trainingDataset)

	fmt.Printf("Number of States: %d\n", len(result.States))
	fmt.Printf("Duration: %.5fs\n", mergeData.Duration.Seconds())
	fmt.Printf("Merges/s: %d\n", int(math.Round(mergeData.AttemptedMergesPerSecond())))
	fmt.Printf("Attempted Merges: %d\n", mergeData.AttemptedMergesCount)
	fmt.Printf("Valid Merges: %d\n", mergeData.ValidMergesCount)
}
