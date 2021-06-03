package main

import (
	"fmt"
	dfalearningtoolkit "github.com/Cherrett/DFA-Learning-Toolkit/core"
	"math"
	"math/rand"
	"time"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	// Random Seed.
	rand.Seed(time.Now().UnixNano())
	
	for _, fileName := range []string{"2_16_100.txt", "2_32_250.txt", "2_32_500.txt", "2_64_750.txt"}{
		trainingDataset := dfalearningtoolkit.GetDatasetFromStaminaFile(fmt.Sprintf("datasets/Comparison/%s", fileName))

		fmt.Println("Positive Strings:", trainingDataset.AcceptingStringInstancesCount())
		fmt.Println("Negative Strings:", trainingDataset.RejectingStringInstancesCount())
		fmt.Println("Total Strings:", len(trainingDataset))

		result, mergeData := dfalearningtoolkit.RPNIFromDataset(trainingDataset)

		fmt.Printf("-------------RPNI-------------\n")
		fmt.Printf("Number of States: %d\n", len(result.States))
		fmt.Printf("Duration: %.5fs\n", mergeData.Duration.Seconds())
		fmt.Printf("Merges/s: %d\n", int(math.Round(mergeData.AttemptedMergesPerSecond())))
		fmt.Printf("Attempted Merges: %d\n", mergeData.AttemptedMergesCount)
		fmt.Printf("Valid Merges: %d\n", mergeData.ValidMergesCount)

		result, mergeData = dfalearningtoolkit.BlueFringeEDSMFromDataset(trainingDataset)

		fmt.Printf("-------------EDSM-------------\n")
		fmt.Printf("Number of States: %d\n", len(result.States))
		fmt.Printf("Duration: %.5fs\n", mergeData.Duration.Seconds())
		fmt.Printf("Merges/s: %d\n", int(math.Round(mergeData.AttemptedMergesPerSecond())))
		fmt.Printf("Attempted Merges: %d\n", mergeData.AttemptedMergesCount)
		fmt.Printf("Valid Merges: %d\n", mergeData.ValidMergesCount)
	}
}