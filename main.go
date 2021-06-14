package main

import (
	"fmt"
	dfalearningtoolkit "github.com/Cherrett/DFA-Learning-Toolkit/core"
	"github.com/Cherrett/DFA-Learning-Toolkit/util"
	"math"
	"math/rand"
	"os"
	"runtime"
	"text/tabwriter"
	"time"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	fmt.Println("BenchmarkMergeStates")

	// These are target DFA sizes we will test.
	dfaSizes := []int{16, 32, 64}
	// These are the training set sizes we will test.
	trainingSetSizes := []int{230, 607, 1521}

	// Benchmark over the problem instances.
	for iterator := range dfaSizes {
		targetSize := dfaSizes[iterator]
		trainingSetSize := trainingSetSizes[iterator]

		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("BENCHMARK %d (Target: %d states, Training: %d strings\n", iterator+1, targetSize, trainingSetSize)
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")

		// Read APTA.
		apta, _ := dfalearningtoolkit.DFAFromJSON(fmt.Sprintf("datasets/TestingAPTAs/%d.json", targetSize))

		fmt.Printf("APTA size: %d\n", len(apta.States))

		// Perform all the merges.
		part := apta.ToStatePartition()
		snapshot := part.Copy()
		totalMerges := 0
		validMerges := 0
		start := time.Now()

		for i := 0; i < len(apta.States); i++ {
			for j := i + 1; j < len(apta.States); j++ {
				totalMerges++
				if snapshot.MergeStates(i, j) {
					validMerges++
				}

				snapshot.RollbackChangesFrom(part)
			}
		}

		totalTime := (time.Now()).Sub(start).Seconds()
		fmt.Printf("Total merges: %d\n", totalMerges)
		fmt.Printf("Valid merges: %d\n", validMerges)
		fmt.Printf("Time: %.4fs\n", totalTime)
		fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)
	}

	fmt.Println("\nExhaustive EDSM")

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	numberOfStates := util.NewStatsTracker()
	durations := util.NewStatsTracker()
	mergesPerSec := util.NewStatsTracker()
	merges := util.NewStatsTracker()
	validMerges := util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Read APTA from file.
		apta, _ := dfalearningtoolkit.DFAFromJSON(fmt.Sprintf("datasets/Generated Abbadingo/%d/%d.json", targetSize, i))

		resultantDFA, mergeData := dfalearningtoolkit.ExhaustiveEDSM(*apta)

		numberOfStates.AddInt(len(resultantDFA.States))
		durations.Add(mergeData.Duration.Seconds())
		mergesPerSec.Add(mergeData.AttemptedMergesPerSecond())
		merges.AddInt(mergeData.AttemptedMergesCount)
		validMerges.AddInt(mergeData.ValidMergesCount)
	}

	fmt.Println("--------------------------------------------------------------------------------------------")
	PrintBenchmarkInformation(numberOfStates, durations, mergesPerSec, merges, validMerges)
	fmt.Println("--------------------------------------------------------------------------------------------")

	fmt.Println("\nRPNI")

	// Number of iterations.
	n = 128
	// Target size.
	targetSize = 32

	numberOfStates = util.NewStatsTracker()
	durations = util.NewStatsTracker()
	mergesPerSec = util.NewStatsTracker()
	merges = util.NewStatsTracker()
	validMerges = util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Read APTA from file.
		apta, _ := dfalearningtoolkit.DFAFromJSON(fmt.Sprintf("datasets/Generated Abbadingo/%d/%d.json", targetSize, i))

		resultantDFA, mergeData := dfalearningtoolkit.RPNI(*apta)

		numberOfStates.AddInt(len(resultantDFA.States))
		durations.Add(mergeData.Duration.Seconds())
		mergesPerSec.Add(mergeData.AttemptedMergesPerSecond())
		merges.AddInt(mergeData.AttemptedMergesCount)
		validMerges.AddInt(mergeData.ValidMergesCount)
	}

	fmt.Println("--------------------------------------------------------------------------------------------")
	PrintBenchmarkInformation(numberOfStates, durations, mergesPerSec, merges, validMerges)
	fmt.Println("--------------------------------------------------------------------------------------------")
}

func PrintBenchmarkInformation(numberOfStates, duration, mergesPerSec, merges, validMerges util.StatsTracker) {
	// Initialize tabwriter.
	w := new(tabwriter.Writer)

	// Determine OS tab width using runtime.GOOS.
	tabWidth := 4

	if runtime.GOOS != "windows" {
		tabWidth = 8
	}

	w.Init(os.Stdout, 17, tabWidth, 0, '\t', 0)

	_, _ = fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\t\n", "Minimum", "Maximum", "Average", "Standard Dev")
	_, _ = fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\t\n", "------------", "------------", "------------", "------------")
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Number of States", int(numberOfStates.Min()), int(numberOfStates.Max()), int(numberOfStates.Mean()), int(numberOfStates.PopulationStandardDev()))
	_, _ = fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t\n", "Duration", duration.Min(), duration.Max(), duration.Mean(), duration.PopulationStandardDev())
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Merges/s", int(math.Round(mergesPerSec.Min())), int(math.Round(mergesPerSec.Max())), int(math.Round(mergesPerSec.Mean())), int(math.Round(mergesPerSec.PopulationStandardDev())))
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Attempted Merges", int(merges.Min()), int(merges.Max()), int(merges.Mean()), int(merges.PopulationStandardDev()))
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Valid Merges", int(validMerges.Min()), int(validMerges.Max()), int(validMerges.Mean()), int(validMerges.PopulationStandardDev()))

	_ = w.Flush()
}
