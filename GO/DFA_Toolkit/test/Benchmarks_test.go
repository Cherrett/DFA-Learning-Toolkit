package dfatoolkit_test

import (
	"DFA_Toolkit/DFA_Toolkit"
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// -------------------- BENCHMARKS USING ABBADINGO PROTOCOL --------------------

// TestBenchmarkDetMerge benchmarks the performance of the MergeStates() function.
func TestBenchmarkDetMerge(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// These are target DFA sizes we will test.
	dfaSizes := []int{16, 32, 64}
	// These are the training set sizes we will test.
	trainingSetSizes := []int{230, 607, 1521}

	// Benchmark over the problem instances.
	for iterator := range dfaSizes {
		targetSize := dfaSizes[iterator]
		trainingSetSize := trainingSetSizes[iterator]

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training set.
		training, _ := dfatoolkit.AbbadingoDatasetExact(target, trainingSetSize, 0)

		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("BENCHMARK %d (Target: %d states, Training: %d strings\n", iterator+1, targetSize, len(training))
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")

		// Info about training set.
		fmt.Printf("Training proportion +ve: %.2f%%\n", training.AcceptingStringInstancesRatio()*100.0)
		fmt.Printf("Training proportion -ve: %.2f%%\n", training.RejectingStringInstancesRatio()*100.0)

		// Create APTA.
		apta := training.GetPTA(true)
		//apta := training.APTA(target.AlphabetSize)
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
		fmt.Printf("Time: %.2fs\n", totalTime)
		fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)
	}
}

// TestBenchmarkRPNI benchmarks the performance of the RPNIFromDataset() function.
func TestBenchmarkRPNI(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 607, 1800

	winners := 0
	accuracies := util.NewStatsTracker()
	numberOfStates := util.NewStatsTracker()
	durations := util.NewStatsTracker()
	mergesPerSec := util.NewStatsTracker()
	merges := util.NewStatsTracker()
	validMerges := util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		resultantDFA, searchData := dfatoolkit.RPNIFromDataset(trainingSet)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.AddInt(len(resultantDFA.States))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())
		merges.AddInt(searchData.AttemptedMergesCount)
		validMerges.AddInt(searchData.ValidMergesCount)

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := (float64(winners) / float64(n)) * 100
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Mean(), accuracies.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStates.Min()), int(numberOfStates.Max()), int(math.Round(numberOfStates.Mean())), int(numberOfStates.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durations.Min(), durations.Max(), durations.Mean(), durations.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Mean(), mergesPerSec.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(merges.Min()), int(merges.Max()), int(math.Round(merges.Mean())), int(merges.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMerges.Min()), int(validMerges.Max()), int(math.Round(validMerges.Mean())), int(validMerges.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")

	if successfulPercentage > 0 {
		t.Error("The percentage of successful DFAs is bigger than 0%.")
	}
}

// TestBenchmarkGreedyEDSM benchmarks the performance of the GreedyEDSMFromDataset() function.
func TestBenchmarkGreedyEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 607, 1800

	winners := 0
	accuracies := util.NewStatsTracker()
	numberOfStates := util.NewStatsTracker()
	durations := util.NewStatsTracker()
	mergesPerSec := util.NewStatsTracker()
	merges := util.NewStatsTracker()
	validMerges := util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		resultantDFA, searchData := dfatoolkit.GreedyEDSMFromDataset(trainingSet)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.AddInt(len(resultantDFA.States))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())
		merges.AddInt(searchData.AttemptedMergesCount)
		validMerges.AddInt(searchData.ValidMergesCount)

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := (float64(winners) / float64(n)) * 100
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Mean(), accuracies.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStates.Min()), int(numberOfStates.Max()), int(math.Round(numberOfStates.Mean())), int(numberOfStates.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durations.Min(), durations.Max(), durations.Mean(), durations.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Mean(), mergesPerSec.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(merges.Min()), int(merges.Max()), int(math.Round(merges.Mean())), int(merges.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMerges.Min()), int(validMerges.Max()), int(math.Round(validMerges.Mean())), int(validMerges.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentage < 9 || successfulPercentage > 15 {
			t.Error("The percentage of successful DFAs is less than 9% or bigger than 15%.")
		}
	}
}

// TestBenchmarkFastWindowedEDSM benchmarks the performance of the FastWindowedEDSMFromDataset() function.
func TestBenchmarkFastWindowedEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 607, 1800

	winners := 0
	accuracies := util.NewStatsTracker()
	numberOfStates := util.NewStatsTracker()
	durations := util.NewStatsTracker()
	mergesPerSec := util.NewStatsTracker()
	merges := util.NewStatsTracker()
	validMerges := util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		resultantDFA, searchData := dfatoolkit.FastWindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.AddInt(len(resultantDFA.States))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())
		merges.AddInt(searchData.AttemptedMergesCount)
		validMerges.AddInt(searchData.ValidMergesCount)

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := (float64(winners) / float64(n)) * 100
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Mean(), accuracies.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStates.Min()), int(numberOfStates.Max()), int(math.Round(numberOfStates.Mean())), int(numberOfStates.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durations.Min(), durations.Max(), durations.Mean(), durations.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Mean(), mergesPerSec.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(merges.Min()), int(merges.Max()), int(math.Round(merges.Mean())), int(merges.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMerges.Min()), int(validMerges.Max()), int(math.Round(validMerges.Mean())), int(validMerges.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentage < 7 || successfulPercentage > 15 {
			t.Error("The percentage of successful DFAs is less than 7% or bigger than 15%.")
		}
	}
}

// TestBenchmarkWindowedEDSM benchmarks the performance of the WindowedEDSMFromDataset() function.
func TestBenchmarkWindowedEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 607, 1800

	winners := 0
	accuracies := util.NewStatsTracker()
	numberOfStates := util.NewStatsTracker()
	durations := util.NewStatsTracker()
	mergesPerSec := util.NewStatsTracker()
	merges := util.NewStatsTracker()
	validMerges := util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		resultantDFA, searchData := dfatoolkit.WindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.AddInt(len(resultantDFA.States))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())
		merges.AddInt(searchData.AttemptedMergesCount)
		validMerges.AddInt(searchData.ValidMergesCount)

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := (float64(winners) / float64(n)) * 100
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Mean(), accuracies.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStates.Min()), int(numberOfStates.Max()), int(math.Round(numberOfStates.Mean())), int(numberOfStates.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durations.Min(), durations.Max(), durations.Mean(), durations.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Mean(), mergesPerSec.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(merges.Min()), int(merges.Max()), int(math.Round(merges.Mean())), int(merges.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMerges.Min()), int(validMerges.Max()), int(math.Round(validMerges.Mean())), int(validMerges.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentage < 7 || successfulPercentage > 15 {
			t.Error("The percentage of successful DFAs is less than 7% or bigger than 15%.")
		}
	}
}

// TestBenchmarkBlueFringeEDSM benchmarks the performance of the BlueFringeEDSMFromDataset() function.
func TestBenchmarkBlueFringeEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 607, 1800

	winners := 0
	accuracies := util.NewStatsTracker()
	numberOfStates := util.NewStatsTracker()
	durations := util.NewStatsTracker()
	mergesPerSec := util.NewStatsTracker()
	merges := util.NewStatsTracker()
	validMerges := util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		resultantDFA, searchData := dfatoolkit.BlueFringeEDSMFromDataset(trainingSet)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.AddInt(len(resultantDFA.States))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())
		merges.AddInt(searchData.AttemptedMergesCount)
		validMerges.AddInt(searchData.ValidMergesCount)

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := (float64(winners) / float64(n)) * 100
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Mean(), accuracies.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStates.Min()), int(numberOfStates.Max()), int(math.Round(numberOfStates.Mean())), int(numberOfStates.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durations.Min(), durations.Max(), durations.Mean(), durations.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Mean(), mergesPerSec.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(merges.Min()), int(merges.Max()), int(math.Round(merges.Mean())), int(merges.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMerges.Min()), int(validMerges.Max()), int(math.Round(validMerges.Mean())), int(validMerges.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentage < 7 || successfulPercentage > 15 {
			t.Error("The percentage of successful DFAs is less than 7% or bigger than 15%.")
		}
	}
}

// TestBenchmarkEDSM benchmarks the performance of the GreedyEDSMFromDataset(), FastWindowedEDSMFromDataset(),
// WindowedEDSMFromDataset() and BlueFringeEDSMFromDataset() functions while comparing their performance.
func TestBenchmarkEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 607, 1800

	// Initialize values.
	winnersGreedy, winnersFastWindowed, winnersWindowed, winnersBlueFringe := 0, 0, 0, 0
	accuraciesGreedy, accuraciesFastWindowed, accuraciesWindowed, accuraciesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	numberOfStatesGreedy, numberOfStatesFastWindowed, numberOfStatesWindowed, numberOfStatesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	durationGreedy, durationFastWindowed, durationWindowed, durationBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesPerSecGreedy, mergesPerSecFastWindowed, mergesPerSecWindowed, mergesPerSecBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesGreedy, mergesFastWindowed, mergesWindowed, mergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	validMergesGreedy, validMergesFastWindowed, validMergesWindowed, validMergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		// Greedy
		resultantDFA, searchData := dfatoolkit.GreedyEDSMFromDataset(trainingSet)
		durationGreedy.Add(searchData.Duration.Seconds())
		mergesPerSecGreedy.Add(searchData.AttemptedMergesPerSecond())
		accuracy := resultantDFA.Accuracy(testingSet)
		accuraciesGreedy.Add(accuracy)
		numberOfStatesGreedy.AddInt(len(resultantDFA.States))
		mergesGreedy.AddInt(searchData.AttemptedMergesCount)
		validMergesGreedy.AddInt(searchData.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersGreedy++
		}

		// Fast Windowed
		resultantDFA, searchData = dfatoolkit.FastWindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		durationFastWindowed.Add(searchData.Duration.Seconds())
		mergesPerSecFastWindowed.Add(searchData.AttemptedMergesPerSecond())
		accuracy = resultantDFA.Accuracy(testingSet)
		accuraciesFastWindowed.Add(accuracy)
		numberOfStatesFastWindowed.AddInt(len(resultantDFA.States))
		mergesFastWindowed.AddInt(searchData.AttemptedMergesCount)
		validMergesFastWindowed.AddInt(searchData.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersFastWindowed++
		}

		// Windowed
		resultantDFA, searchData = dfatoolkit.WindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		durationWindowed.Add(searchData.Duration.Seconds())
		mergesPerSecWindowed.Add(searchData.AttemptedMergesPerSecond())
		accuracy = resultantDFA.Accuracy(testingSet)
		accuraciesWindowed.Add(accuracy)
		numberOfStatesWindowed.AddInt(len(resultantDFA.States))
		mergesWindowed.AddInt(searchData.AttemptedMergesCount)
		validMergesWindowed.AddInt(searchData.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersWindowed++
		}

		// Blue-Fringe
		resultantDFA, searchData = dfatoolkit.BlueFringeEDSMFromDataset(trainingSet)
		durationBlueFringe.Add(searchData.Duration.Seconds())
		mergesPerSecBlueFringe.Add(searchData.AttemptedMergesPerSecond())
		accuracy = resultantDFA.Accuracy(testingSet)
		accuraciesBlueFringe.Add(accuracy)
		numberOfStatesBlueFringe.AddInt(len(resultantDFA.States))
		mergesBlueFringe.AddInt(searchData.AttemptedMergesCount)
		validMergesBlueFringe.AddInt(searchData.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersBlueFringe++
		}
	}

	successfulPercentageGreedy := (float64(winnersGreedy) / float64(n)) * 100
	successfulPercentageFastWindowed := (float64(winnersFastWindowed) / float64(n)) * 100
	successfulPercentageWindowed := (float64(winnersWindowed) / float64(n)) * 100
	successfulPercentageBlueFringe := (float64(winnersBlueFringe) / float64(n)) * 100

	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Greedy Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageGreedy)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesGreedy.Min(), accuraciesGreedy.Max(), accuraciesGreedy.Mean(), accuraciesGreedy.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesGreedy.Min()), int(numberOfStatesGreedy.Max()), int(math.Round(numberOfStatesGreedy.Mean())), int(numberOfStatesGreedy.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationGreedy.Min(), durationGreedy.Max(), durationGreedy.Mean(), durationGreedy.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecGreedy.Min(), mergesPerSecGreedy.Max(), mergesPerSecGreedy.Mean(), mergesPerSecGreedy.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesGreedy.Min()), int(mergesGreedy.Max()), int(math.Round(mergesGreedy.Mean())), int(mergesGreedy.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesGreedy.Min()), int(validMergesGreedy.Max()), int(math.Round(validMergesGreedy.Mean())), int(validMergesGreedy.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Fast Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageFastWindowed)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesFastWindowed.Min(), accuraciesFastWindowed.Max(), accuraciesFastWindowed.Mean(), accuraciesFastWindowed.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesFastWindowed.Min()), int(numberOfStatesFastWindowed.Max()), int(math.Round(numberOfStatesFastWindowed.Mean())), int(numberOfStatesFastWindowed.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationFastWindowed.Min(), durationFastWindowed.Max(), durationFastWindowed.Mean(), durationFastWindowed.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecFastWindowed.Min(), mergesPerSecFastWindowed.Max(), mergesPerSecFastWindowed.Mean(), mergesPerSecFastWindowed.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesFastWindowed.Min()), int(mergesFastWindowed.Max()), int(math.Round(mergesFastWindowed.Mean())), int(mergesFastWindowed.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesFastWindowed.Min()), int(validMergesFastWindowed.Max()), int(math.Round(validMergesFastWindowed.Mean())), int(validMergesFastWindowed.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageWindowed)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesWindowed.Min(), accuraciesWindowed.Max(), accuraciesWindowed.Mean(), accuraciesWindowed.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesWindowed.Min()), int(numberOfStatesWindowed.Max()), int(math.Round(numberOfStatesWindowed.Mean())), int(numberOfStatesWindowed.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationWindowed.Min(), durationWindowed.Max(), durationWindowed.Mean(), durationWindowed.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecWindowed.Min(), mergesPerSecWindowed.Max(), mergesPerSecWindowed.Mean(), mergesPerSecWindowed.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesWindowed.Min()), int(mergesWindowed.Max()), int(math.Round(mergesWindowed.Mean())), int(mergesWindowed.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesWindowed.Min()), int(validMergesWindowed.Max()), int(math.Round(validMergesWindowed.Mean())), int(validMergesWindowed.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Blue-Fringe Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageBlueFringe)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesBlueFringe.Min(), accuraciesBlueFringe.Max(), accuraciesBlueFringe.Mean(), accuraciesBlueFringe.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesBlueFringe.Min()), int(numberOfStatesBlueFringe.Max()), int(math.Round(numberOfStatesBlueFringe.Mean())), int(numberOfStatesBlueFringe.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationBlueFringe.Min(), durationBlueFringe.Max(), durationBlueFringe.Mean(), durationBlueFringe.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecBlueFringe.Min(), mergesPerSecBlueFringe.Max(), mergesPerSecBlueFringe.Mean(), mergesPerSecBlueFringe.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesBlueFringe.Min()), int(mergesBlueFringe.Max()), int(math.Round(mergesBlueFringe.Mean())), int(mergesBlueFringe.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesBlueFringe.Min()), int(validMergesBlueFringe.Max()), int(math.Round(validMergesBlueFringe.Mean())), int(validMergesBlueFringe.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentageGreedy < 9 || successfulPercentageGreedy > 15 {
			t.Error("The percentage of successful DFAs for Greedy EDSM is less than 9% or bigger than 15%.")
		}

		if successfulPercentageFastWindowed < 7 || successfulPercentageFastWindowed > 15 {
			t.Error("The percentage of successful DFAs for Fast Windowed EDSM is less than 7% or bigger than 15%.")
		}

		if successfulPercentageWindowed < 7 || successfulPercentageWindowed > 15 {
			t.Error("The percentage of successful DFAs for Windowed EDSM is less than 7% or bigger than 15%.")
		}

		if successfulPercentageBlueFringe < 7 || successfulPercentageBlueFringe > 15 {
			t.Error("The percentage of successful DFAs for Blue-Fringe EDSM is less than 7% or bigger than 15%.")
		}
	}
}

// TestBenchmarkEDSM concurrently benchmarks the performance of the GreedyEDSM(), FastWindowedEDSM(),
// WindowedEDSM() and BlueFringeEDSM() functions while comparing their performance.
func TestBenchmarkFastEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 607, 1800

	// Initialize values.
	winnersGreedy, winnersFastWindowed, winnersWindowed, winnersBlueFringe := 0, 0, 0, 0
	accuraciesGreedy, accuraciesFastWindowed, accuraciesWindowed, accuraciesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	numberOfStatesGreedy, numberOfStatesFastWindowed, numberOfStatesWindowed, numberOfStatesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	durationGreedy, durationFastWindowed, durationWindowed, durationBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesPerSecGreedy, mergesPerSecFastWindowed, mergesPerSecWindowed, mergesPerSecBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesGreedy, mergesFastWindowed, mergesWindowed, mergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	validMergesGreedy, validMergesFastWindowed, validMergesWindowed, validMergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		// Create APTA from training set.
		APTA := trainingSet.GetPTA(true)

		// Create wait group
		var wg sync.WaitGroup
		// Add 4 EDSM types to wait group.
		wg.Add(4)

		resultantDFAGreedy, resultantDFAFastWindowed, resultantDFAWindowed, resultantDFABlueFringe := dfatoolkit.DFA{}, dfatoolkit.DFA{}, dfatoolkit.DFA{}, dfatoolkit.DFA{}
		searchDataGreedy, searchDataFastWindowed, searchDataWindowed, searchDataBlueFringe := dfatoolkit.SearchData{}, dfatoolkit.SearchData{}, dfatoolkit.SearchData{}, dfatoolkit.SearchData{}

		// Greedy
		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAGreedy, searchDataGreedy = dfatoolkit.GreedyEDSM(APTA)
		}()

		// Fast Windowed
		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAFastWindowed, searchDataFastWindowed = dfatoolkit.FastWindowedEDSM(APTA, targetSize*2, 2.0)
		}()

		// Windowed
		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAWindowed, searchDataWindowed = dfatoolkit.WindowedEDSM(APTA, targetSize*2, 2.0)
		}()

		// Blue-Fringe
		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFABlueFringe, searchDataBlueFringe = dfatoolkit.BlueFringeEDSM(APTA)
		}()

		// Wait for all go routines within wait group to finish executing.
		wg.Wait()

		// Greedy
		durationGreedy.Add(searchDataGreedy.Duration.Seconds())
		mergesPerSecGreedy.Add(searchDataGreedy.AttemptedMergesPerSecond())
		accuracy := resultantDFAGreedy.Accuracy(testingSet)
		accuraciesGreedy.Add(accuracy)
		numberOfStatesGreedy.AddInt(len(resultantDFAGreedy.States))
		mergesGreedy.AddInt(searchDataGreedy.AttemptedMergesCount)
		validMergesGreedy.AddInt(searchDataGreedy.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersGreedy++
		}

		// Fast Windowed
		durationFastWindowed.Add(searchDataFastWindowed.Duration.Seconds())
		mergesPerSecFastWindowed.Add(searchDataFastWindowed.AttemptedMergesPerSecond())
		accuracy = resultantDFAFastWindowed.Accuracy(testingSet)
		accuraciesFastWindowed.Add(accuracy)
		numberOfStatesFastWindowed.AddInt(len(resultantDFAFastWindowed.States))
		mergesFastWindowed.AddInt(searchDataFastWindowed.AttemptedMergesCount)
		validMergesFastWindowed.AddInt(searchDataFastWindowed.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersFastWindowed++
		}

		// Windowed
		durationWindowed.Add(searchDataWindowed.Duration.Seconds())
		mergesPerSecWindowed.Add(searchDataWindowed.AttemptedMergesPerSecond())
		accuracy = resultantDFAWindowed.Accuracy(testingSet)
		accuraciesWindowed.Add(accuracy)
		numberOfStatesWindowed.AddInt(len(resultantDFAWindowed.States))
		mergesWindowed.AddInt(searchDataWindowed.AttemptedMergesCount)
		validMergesWindowed.AddInt(searchDataWindowed.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersWindowed++
		}

		// Blue-Fringe
		durationBlueFringe.Add(searchDataBlueFringe.Duration.Seconds())
		mergesPerSecBlueFringe.Add(searchDataBlueFringe.AttemptedMergesPerSecond())
		accuracy = resultantDFABlueFringe.Accuracy(testingSet)
		accuraciesBlueFringe.Add(accuracy)
		numberOfStatesBlueFringe.AddInt(len(resultantDFABlueFringe.States))
		mergesBlueFringe.AddInt(searchDataBlueFringe.AttemptedMergesCount)
		validMergesBlueFringe.AddInt(searchDataBlueFringe.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersBlueFringe++
		}
	}

	successfulPercentageGreedy := (float64(winnersGreedy) / float64(n)) * 100
	successfulPercentageFastWindowed := (float64(winnersFastWindowed) / float64(n)) * 100
	successfulPercentageWindowed := (float64(winnersWindowed) / float64(n)) * 100
	successfulPercentageBlueFringe := (float64(winnersBlueFringe) / float64(n)) * 100

	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Greedy Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageGreedy)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesGreedy.Min(), accuraciesGreedy.Max(), accuraciesGreedy.Mean(), accuraciesGreedy.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesGreedy.Min()), int(numberOfStatesGreedy.Max()), int(math.Round(numberOfStatesGreedy.Mean())), int(numberOfStatesGreedy.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationGreedy.Min(), durationGreedy.Max(), durationGreedy.Mean(), durationGreedy.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecGreedy.Min(), mergesPerSecGreedy.Max(), mergesPerSecGreedy.Mean(), mergesPerSecGreedy.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesGreedy.Min()), int(mergesGreedy.Max()), int(math.Round(mergesGreedy.Mean())), int(mergesGreedy.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesGreedy.Min()), int(validMergesGreedy.Max()), int(math.Round(validMergesGreedy.Mean())), int(validMergesGreedy.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Fast Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageFastWindowed)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesFastWindowed.Min(), accuraciesFastWindowed.Max(), accuraciesFastWindowed.Mean(), accuraciesFastWindowed.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesFastWindowed.Min()), int(numberOfStatesFastWindowed.Max()), int(math.Round(numberOfStatesFastWindowed.Mean())), int(numberOfStatesFastWindowed.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationFastWindowed.Min(), durationFastWindowed.Max(), durationFastWindowed.Mean(), durationFastWindowed.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecFastWindowed.Min(), mergesPerSecFastWindowed.Max(), mergesPerSecFastWindowed.Mean(), mergesPerSecFastWindowed.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesFastWindowed.Min()), int(mergesFastWindowed.Max()), int(math.Round(mergesFastWindowed.Mean())), int(mergesFastWindowed.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesFastWindowed.Min()), int(validMergesFastWindowed.Max()), int(math.Round(validMergesFastWindowed.Mean())), int(validMergesFastWindowed.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageWindowed)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesWindowed.Min(), accuraciesWindowed.Max(), accuraciesWindowed.Mean(), accuraciesWindowed.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesWindowed.Min()), int(numberOfStatesWindowed.Max()), int(math.Round(numberOfStatesWindowed.Mean())), int(numberOfStatesWindowed.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationWindowed.Min(), durationWindowed.Max(), durationWindowed.Mean(), durationWindowed.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecWindowed.Min(), mergesPerSecWindowed.Max(), mergesPerSecWindowed.Mean(), mergesPerSecWindowed.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesWindowed.Min()), int(mergesWindowed.Max()), int(math.Round(mergesWindowed.Mean())), int(mergesWindowed.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesWindowed.Min()), int(validMergesWindowed.Max()), int(math.Round(validMergesWindowed.Mean())), int(validMergesWindowed.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Blue-Fringe Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageBlueFringe)
	fmt.Printf("Accuracy - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", accuraciesBlueFringe.Min(), accuraciesBlueFringe.Max(), accuraciesBlueFringe.Mean(), accuraciesBlueFringe.PopulationStandardDev())
	fmt.Printf("Number of States - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(numberOfStatesBlueFringe.Min()), int(numberOfStatesBlueFringe.Max()), int(math.Round(numberOfStatesBlueFringe.Mean())), int(numberOfStatesBlueFringe.PopulationStandardDev()))
	fmt.Printf("Duration - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", durationBlueFringe.Min(), durationBlueFringe.Max(), durationBlueFringe.Mean(), durationBlueFringe.PopulationStandardDev())
	fmt.Printf("Merges/s - Min: %.2f, Max: %.2f, Avg: %.2f, Standard Dev: %.2f\n", mergesPerSecBlueFringe.Min(), mergesPerSecBlueFringe.Max(), mergesPerSecBlueFringe.Mean(), mergesPerSecBlueFringe.PopulationStandardDev())
	fmt.Printf("Attempted Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(mergesBlueFringe.Min()), int(mergesBlueFringe.Max()), int(math.Round(mergesBlueFringe.Mean())), int(mergesBlueFringe.PopulationStandardDev()))
	fmt.Printf("Valid Merges - Min: %d, Max: %d, Avg: %d, Standard Dev: %d\n", int(validMergesBlueFringe.Min()), int(validMergesBlueFringe.Max()), int(math.Round(validMergesBlueFringe.Mean())), int(validMergesBlueFringe.PopulationStandardDev()))
	fmt.Println("-----------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentageGreedy < 9 || successfulPercentageGreedy > 15 {
			t.Error("The percentage of successful DFAs for Greedy EDSM is less than 9% or bigger than 15%.")
		}

		if successfulPercentageFastWindowed < 7 || successfulPercentageFastWindowed > 15 {
			t.Error("The percentage of successful DFAs for Fast Windowed EDSM is less than 7% or bigger than 15%.")
		}

		if successfulPercentageWindowed < 7 || successfulPercentageWindowed > 15 {
			t.Error("The percentage of successful DFAs for Windowed EDSM is less than 7% or bigger than 15%.")
		}

		if successfulPercentageBlueFringe < 7 || successfulPercentageBlueFringe > 15 {
			t.Error("The percentage of successful DFAs for Blue-Fringe EDSM is less than 7% or bigger than 15%.")
		}
	}
}

// -------------------- BENCHMARKS USING STAMINA PROTOCOL --------------------
