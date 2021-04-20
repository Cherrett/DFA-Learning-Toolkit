package dfatoolkit_test

import (
	"DFA_Toolkit/DFA_Toolkit"
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// -------------------- BENCHMARKS --------------------

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

				snapshot.RollbackChanges(part)
			}
		}

		totalTime := (time.Now()).Sub(start).Seconds()
		fmt.Printf("Total merges: %d\n", totalMerges)
		fmt.Printf("Valid merges: %d\n", validMerges)
		fmt.Printf("Time: %.2fs\n", totalTime)
		fmt.Printf("Merges per second: %.2f\n", float64(totalMerges)/totalTime)
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

	winners := 0
	accuracies := util.NewMinMaxAvg()
	numberOfStates := util.NewMinMaxAvg()
	durations := util.NewMinMaxAvg()
	mergesPerSec := util.NewMinMaxAvg()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := dfatoolkit.AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA, searchData := dfatoolkit.GreedyEDSMFromDataset(trainingSet)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.Add(float64(len(resultantDFA.States)))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStates.Min(), numberOfStates.Max(), numberOfStates.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durations.Min(), durations.Max(), durations.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentage < 0.10 || successfulPercentage > 0.15 {
		t.Error("The percentage of successful DFAs is less than 0.10 or bigger than 0.15.")
	}
}

// TestBenchmarkFastWindowedEDSM benchmarks the performance of the FastWindowedEDSMFromDataset() function.
func TestBenchmarkFastWindowedEDSM(t *testing.T) {
	// Random Seed.
	// rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 20
	// Target size.
	targetSize := 32

	winners := 0
	accuracies := util.NewMinMaxAvg()
	numberOfStates := util.NewMinMaxAvg()
	durations := util.NewMinMaxAvg()
	mergesPerSec := util.NewMinMaxAvg()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := dfatoolkit.AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA, searchData := dfatoolkit.FastWindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.Add(float64(len(resultantDFA.States)))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStates.Min(), numberOfStates.Max(), numberOfStates.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durations.Min(), durations.Max(), durations.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentage < 0.09 || successfulPercentage > 0.15 {
		t.Error("The percentage of successful DFAs is less than 0.09 or bigger than 0.15.")
	}
}

// TestBenchmarkWindowedEDSM benchmarks the performance of the WindowedEDSMFromDataset() function.
func TestBenchmarkWindowedEDSM(t *testing.T) {
	// Random Seed.
	// rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 20
	// Target size.
	targetSize := 32

	winners := 0
	accuracies := util.NewMinMaxAvg()
	numberOfStates := util.NewMinMaxAvg()
	durations := util.NewMinMaxAvg()
	mergesPerSec := util.NewMinMaxAvg()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := dfatoolkit.AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA, searchData := dfatoolkit.WindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.Add(float64(len(resultantDFA.States)))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStates.Min(), numberOfStates.Max(), numberOfStates.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durations.Min(), durations.Max(), durations.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentage < 0.09 || successfulPercentage > 0.15 {
		t.Error("The percentage of successful DFAs is less than 0.09 or bigger than 0.15.")
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

	winners := 0
	accuracies := util.NewMinMaxAvg()
	numberOfStates := util.NewMinMaxAvg()
	durations := util.NewMinMaxAvg()
	mergesPerSec := util.NewMinMaxAvg()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := dfatoolkit.AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA, searchData := dfatoolkit.BlueFringeEDSMFromDataset(trainingSet)
		accuracy := resultantDFA.Accuracy(testingSet)

		accuracies.Add(accuracy)
		numberOfStates.Add(float64(len(resultantDFA.States)))
		durations.Add(searchData.Duration.Seconds())
		mergesPerSec.Add(searchData.AttemptedMergesPerSecond())

		if accuracy >= 0.99 {
			winners++
		}
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuracies.Min(), accuracies.Max(), accuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStates.Min(), numberOfStates.Max(), numberOfStates.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durations.Min(), durations.Max(), durations.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSec.Min(), mergesPerSec.Max(), mergesPerSec.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentage < 0.07 || successfulPercentage > 0.15 {
		t.Error("The percentage of successful DFAs is less than 0.07 or bigger than 0.15.")
	}
}

// TestBenchmarkEDSM benchmarks the performance of the GreedyEDSMFromDataset(), FastWindowedEDSMFromDataset()
// and BlueFringeEDSMFromDataset() functions.
func TestBenchmarkEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	winnersGreedy, winnersWindowed, winnersBlueFringe := 0, 0, 0
	accuraciesGreedy, accuraciesWindowed, accuraciesBlueFringe := util.NewMinMaxAvg(), util.NewMinMaxAvg(), util.NewMinMaxAvg()
	numberOfStatesGreedy, numberOfStatesWindowed, numberOfStatesBlueFringe := util.NewMinMaxAvg(), util.NewMinMaxAvg(), util.NewMinMaxAvg()
	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := dfatoolkit.AbbadingoDatasetExact(target, 607, 1800)

		// Construct an APTA from training dataset.
		APTA := trainingSet.GetPTA(true)

		// Create wait group
		var wg sync.WaitGroup
		// Add 3 EDSM types to wait group.
		wg.Add(3)

		var resultantDFAGreedy dfatoolkit.DFA
		var resultantDFAWindowed dfatoolkit.DFA
		var resultantDFABlueFringe dfatoolkit.DFA

		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAGreedy, _ = dfatoolkit.GreedyEDSM(APTA)
		}()

		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAWindowed, _ = dfatoolkit.FastWindowedEDSM(APTA, targetSize*2, 2.0)
		}()

		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFABlueFringe, _ = dfatoolkit.BlueFringeEDSM(APTA)
		}()

		// Wait for all go routines within wait
		// group to finish executing.
		wg.Wait()

		accuracyGreedy := resultantDFAGreedy.Accuracy(testingSet)
		accuracyWindowed := resultantDFAWindowed.Accuracy(testingSet)
		accuracyBlueFringe := resultantDFABlueFringe.Accuracy(testingSet)

		accuraciesGreedy.Add(accuracyGreedy)
		accuraciesWindowed.Add(accuracyWindowed)
		accuraciesBlueFringe.Add(accuracyBlueFringe)

		numberOfStatesGreedy.Add(float64(len(resultantDFAGreedy.States)))
		numberOfStatesWindowed.Add(float64(len(resultantDFAWindowed.States)))
		numberOfStatesBlueFringe.Add(float64(len(resultantDFABlueFringe.States)))

		if accuracyGreedy >= 0.99 {
			winnersGreedy++
		}

		if accuracyWindowed >= 0.99 {
			winnersWindowed++
		}

		if accuracyBlueFringe >= 0.99 {
			winnersBlueFringe++
		}

		fmt.Printf("Greedy: Accuracy: %.3f, Number of States %d\n", accuracyGreedy, len(resultantDFAGreedy.States))
		fmt.Printf("Windowed: Accuracy: %.3f, Number of States %d\n", accuracyWindowed, len(resultantDFAWindowed.States))
		fmt.Printf("Blue-Fringe: Accuracy: %.3f, Number of States %d\n", accuracyBlueFringe, len(resultantDFABlueFringe.States))
		fmt.Print("-----------------------------------------------------------------------------\n")
	}

	successfulPercentageGreedy := float64(winnersGreedy) / float64(n)
	fmt.Print("Greedy Search:\n")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.3f%%\n", successfulPercentageGreedy)
	fmt.Printf("Minimum Accuracy: %.3f Maximum Accuracy: %.3f Average Accuracy: %.3f\n", accuraciesGreedy.Min(), accuraciesGreedy.Max(), accuraciesGreedy.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStatesGreedy.Min(), numberOfStatesGreedy.Max(), numberOfStatesGreedy.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	successfulPercentageWindowed := float64(winnersWindowed) / float64(n)
	fmt.Printf("Windowed Search:\n")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.3f%%\n", successfulPercentageWindowed)
	fmt.Printf("Minimum Accuracy: %.3f Maximum Accuracy: %.3f Average Accuracy: %.3f\n", accuraciesWindowed.Min(), accuraciesWindowed.Max(), accuraciesWindowed.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStatesWindowed.Min(), numberOfStatesWindowed.Max(), numberOfStatesWindowed.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	successfulPercentageBlueFringe := float64(winnersBlueFringe) / float64(n)
	fmt.Printf("Blue-Fringe Search:\n")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.3f%%\n", successfulPercentageBlueFringe)
	fmt.Printf("Minimum Accuracy: %.3f Maximum Accuracy: %.3f Average Accuracy: %.3f\n", accuraciesBlueFringe.Min(), accuraciesBlueFringe.Max(), accuraciesBlueFringe.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStatesBlueFringe.Min(), numberOfStatesBlueFringe.Max(), numberOfStatesBlueFringe.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentageGreedy < 0.09 || successfulPercentageGreedy > 0.15 {
		t.Error("The percentage of successful DFAs for Greedy EDSM is less than 0.09 or bigger than 0.15.")
	}

	if successfulPercentageWindowed < 0.09 || successfulPercentageWindowed > 0.15 {
		t.Error("The percentage of successful DFAs for Windowed EDSM is less than 0.09 or bigger than 0.15.")
	}

	if successfulPercentageBlueFringe < 0.07 || successfulPercentageBlueFringe > 0.15 {
		t.Error("The percentage of successful DFAs for Blue-Fringe EDSM is less than 0.07 or bigger than 0.15.")
	}
}

// TestBenchmarkEDSM2 benchmarks the performance of the GreedyEDSMFromDataset(), FastWindowedEDSMFromDataset()
// and BlueFringeEDSMFromDataset() functions while comparing performance (merges per second, duration, ..).
func TestBenchmarkEDSM2(t *testing.T){
	// Random Seed.
	// rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	// Greedy
	winnersGreedy := 0
	accuraciesGreedy := util.NewMinMaxAvg()
	numberOfStatesGreedy := util.NewMinMaxAvg()
	durationGreedy := util.NewMinMaxAvg()
	mergesPerSecGreedy := util.NewMinMaxAvg()

	// Fast Windowed
	winnersFastWindowed := 0
	accuraciesFastWindowed := util.NewMinMaxAvg()
	numberOfStatesFastWindowed := util.NewMinMaxAvg()
	durationFastWindowed := util.NewMinMaxAvg()
	mergesPerSecFastWindowed := util.NewMinMaxAvg()
	
	// Windowed
	winnersWindowed := 0
	accuraciesWindowed := util.NewMinMaxAvg()
	numberOfStatesWindowed := util.NewMinMaxAvg()
	durationWindowed := util.NewMinMaxAvg()
	mergesPerSecWindowed := util.NewMinMaxAvg()

	// Blue-Fringe
	winnersBlueFringe := 0
	accuraciesBlueFringe := util.NewMinMaxAvg()
	numberOfStatesBlueFringe := util.NewMinMaxAvg()
	durationBlueFringe := util.NewMinMaxAvg()
	mergesPerSecBlueFringe := util.NewMinMaxAvg()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)
		// Create a target DFA.
		target := dfatoolkit.AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := dfatoolkit.AbbadingoDatasetExact(target, 607, 1800)

		// Greedy
		resultantDFA, searchData := dfatoolkit.GreedyEDSMFromDataset(trainingSet)
		durationGreedy.Add(searchData.Duration.Seconds())
		mergesPerSecGreedy.Add(searchData.AttemptedMergesPerSecond())
		//fmt.Printf("Greedy Duration: %.2fs\n\n", searchData.Duration.Seconds())
		accuracy := resultantDFA.Accuracy(testingSet)
		accuraciesGreedy.Add(accuracy)
		numberOfStatesGreedy.Add(float64(len(resultantDFA.States)))
		if accuracy >= 0.99 {
			winnersGreedy++
		}
		
		// Fast Windowed
		resultantDFA, searchData = dfatoolkit.FastWindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		durationFastWindowed.Add(searchData.Duration.Seconds())
		mergesPerSecFastWindowed.Add(searchData.AttemptedMergesPerSecond())
		//fmt.Printf("Windowed Duration: %.2fs\n\n", searchData.Duration.Seconds())
		accuracy = resultantDFA.Accuracy(testingSet)
		accuraciesFastWindowed.Add(accuracy)
		numberOfStatesFastWindowed.Add(float64(len(resultantDFA.States)))
		if accuracy >= 0.99 {
			winnersFastWindowed++
		}
		
		// Windowed
		resultantDFA, searchData = dfatoolkit.WindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		durationWindowed.Add(searchData.Duration.Seconds())
		mergesPerSecWindowed.Add(searchData.AttemptedMergesPerSecond())
		//fmt.Printf("Windowed 2 Duration: %.2fs\n\n", searchData.Duration.Seconds())
		accuracy = resultantDFA.Accuracy(testingSet)
		accuraciesWindowed.Add(accuracy)
		numberOfStatesWindowed.Add(float64(len(resultantDFA.States)))
		if accuracy >= 0.99 {
			winnersWindowed++
		}

		// Blue-Fringe
		resultantDFA, searchData = dfatoolkit.BlueFringeEDSMFromDataset(trainingSet)
		durationBlueFringe.Add(searchData.Duration.Seconds())
		mergesPerSecBlueFringe.Add(searchData.AttemptedMergesPerSecond())
		//fmt.Printf("Blue-Fringe Duration: %.2fs\n\n", searchData.Duration.Seconds())
		accuracy = resultantDFA.Accuracy(testingSet)
		accuraciesBlueFringe.Add(accuracy)
		numberOfStatesBlueFringe.Add(float64(len(resultantDFA.States)))
		if accuracy >= 0.99 {
			winnersBlueFringe++
		}
	}

	successfulPercentageGreedy := float64(winnersGreedy) / float64(n)
	successfulPercentageWindowed := float64(winnersFastWindowed) / float64(n)
	successfulPercentageWindowed2 := float64(winnersWindowed) / float64(n)
	successfulPercentageBlueFringe := float64(winnersBlueFringe) / float64(n)

	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Greedy Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageGreedy)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuraciesGreedy.Min(), accuraciesGreedy.Max(), accuraciesGreedy.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStatesGreedy.Min(), numberOfStatesGreedy.Max(), numberOfStatesGreedy.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durationGreedy.Min(), durationGreedy.Max(), durationGreedy.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSecGreedy.Min(), mergesPerSecGreedy.Max(), mergesPerSecGreedy.Avg())
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Fast Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageWindowed)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuraciesFastWindowed.Min(), accuraciesFastWindowed.Max(), accuraciesFastWindowed.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStatesFastWindowed.Min(), numberOfStatesFastWindowed.Max(), numberOfStatesFastWindowed.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durationFastWindowed.Min(), durationFastWindowed.Max(), durationFastWindowed.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSecFastWindowed.Min(), mergesPerSecFastWindowed.Max(), mergesPerSecFastWindowed.Avg())
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageWindowed2)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuraciesWindowed.Min(), accuraciesWindowed.Max(), accuraciesWindowed.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStatesWindowed.Min(), numberOfStatesWindowed.Max(), numberOfStatesWindowed.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durationWindowed.Min(), durationWindowed.Max(), durationWindowed.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSecWindowed.Min(), mergesPerSecWindowed.Max(), mergesPerSecWindowed.Avg())
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println("Blue-Fringe")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageBlueFringe)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", accuraciesBlueFringe.Min(), accuraciesBlueFringe.Max(), accuraciesBlueFringe.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", numberOfStatesBlueFringe.Min(), numberOfStatesBlueFringe.Max(), numberOfStatesBlueFringe.Avg())
	fmt.Printf("Minimum Duration: %.2f Maximum Duration: %.2f Average Duration: %.2f\n", durationBlueFringe.Min(), durationBlueFringe.Max(), durationBlueFringe.Avg())
	fmt.Printf("Minimum Merges/s: %.2f Maximum Merges/s: %.2f Average Merges/s: %.2f\n", mergesPerSecBlueFringe.Min(), mergesPerSecBlueFringe.Max(), mergesPerSecBlueFringe.Avg())
	fmt.Println("-----------------------------------------------------------------------------")
}