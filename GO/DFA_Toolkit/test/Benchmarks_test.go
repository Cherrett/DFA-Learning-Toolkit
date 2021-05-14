package dfatoolkit_test

import (
	"DFA_Toolkit/DFA_Toolkit"
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
	"testing"
	"text/tabwriter"
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
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentage)
	PrintBenchmarkInformation(accuracies, numberOfStates, durations, mergesPerSec, merges, validMerges)
	fmt.Println("--------------------------------------------------------------------------------------------")

	if successfulPercentage > 0 {
		t.Error("The percentage of successful DFAs is bigger than 0%.")
	}
}

// TestBenchmarkExhaustiveEDSM benchmarks the performance of the ExhaustiveEDSMFromDataset() function.
func TestBenchmarkExhaustiveEDSM(t *testing.T) {
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

		resultantDFA, searchData := dfatoolkit.ExhaustiveEDSMFromDataset(trainingSet)
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
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentage)
	PrintBenchmarkInformation(accuracies, numberOfStates, durations, mergesPerSec, merges, validMerges)
	fmt.Println("--------------------------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentage < 9 || successfulPercentage > 15 {
			t.Error("The percentage of successful DFAs is less than 9% or bigger than 15%.")
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
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentage)
	PrintBenchmarkInformation(accuracies, numberOfStates, durations, mergesPerSec, merges, validMerges)
	fmt.Println("--------------------------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentage < 7 || successfulPercentage > 15 {
			t.Error("The percentage of successful DFAs is less than 7% or bigger than 15%.")
		}
	}
}

// TestBenchmarkBlueFringeEDSM benchmarks the performance of the BlueFringeEDSMFromDataset() function.
func TestBenchmarkBlueFringeEDSM(t *testing.T) {
	// Random Seed.
	//rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 256
	// Target size.
	targetSize := 64
	// Training and testing set sizes.
	trainingSetSize, testingSetSize := 1521, 1800

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
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentage)
	PrintBenchmarkInformation(accuracies, numberOfStates, durations, mergesPerSec, merges, validMerges)
	fmt.Println("--------------------------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentage < 7 || successfulPercentage > 15 {
			t.Error("The percentage of successful DFAs is less than 7% or bigger than 15%.")
		}
	}
}

// TestBenchmarkEDSM benchmarks the performance of the ExhaustiveEDSMFromDataset(), FastWindowedEDSMFromDataset(),
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
	winnersExhaustive, winnersWindowed, winnersBlueFringe := 0, 0, 0
	accuraciesExhaustive, accuraciesWindowed, accuraciesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	numberOfStatesExhaustive, numberOfStatesWindowed, numberOfStatesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	durationExhaustive, durationWindowed, durationBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesPerSecExhaustive, mergesPerSecWindowed, mergesPerSecBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesExhaustive, mergesWindowed, mergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	validMergesExhaustive, validMergesWindowed, validMergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()

	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA, training set, and testing set.
		_, trainingSet, testingSet := dfatoolkit.AbbadingoInstanceExact(targetSize, true, trainingSetSize, testingSetSize)

		// Exhaustive
		resultantDFA, searchData := dfatoolkit.ExhaustiveEDSMFromDataset(trainingSet)
		durationExhaustive.Add(searchData.Duration.Seconds())
		mergesPerSecExhaustive.Add(searchData.AttemptedMergesPerSecond())
		accuracy := resultantDFA.Accuracy(testingSet)
		accuraciesExhaustive.Add(accuracy)
		numberOfStatesExhaustive.AddInt(len(resultantDFA.States))
		mergesExhaustive.AddInt(searchData.AttemptedMergesCount)
		validMergesExhaustive.AddInt(searchData.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersExhaustive++
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

	successfulPercentageExhaustive := (float64(winnersExhaustive) / float64(n)) * 100
	successfulPercentageWindowed := (float64(winnersWindowed) / float64(n)) * 100
	successfulPercentageBlueFringe := (float64(winnersBlueFringe) / float64(n)) * 100

	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Exhaustive Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentageExhaustive)
	PrintBenchmarkInformation(accuraciesExhaustive, numberOfStatesExhaustive, durationExhaustive, mergesPerSecExhaustive, mergesExhaustive, validMergesExhaustive)
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentageWindowed)
	PrintBenchmarkInformation(accuraciesWindowed, numberOfStatesWindowed, durationWindowed, mergesPerSecWindowed, mergesWindowed, validMergesWindowed)
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Blue-Fringe Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentageBlueFringe)
	PrintBenchmarkInformation(accuraciesBlueFringe, numberOfStatesBlueFringe, durationBlueFringe, mergesPerSecBlueFringe, mergesBlueFringe, validMergesBlueFringe)
	fmt.Println("--------------------------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentageExhaustive < 9 || successfulPercentageExhaustive > 15 {
			t.Error("The percentage of successful DFAs for Exhaustive EDSM is less than 9% or bigger than 15%.")
		}

		if successfulPercentageWindowed < 7 || successfulPercentageWindowed > 15 {
			t.Error("The percentage of successful DFAs for Windowed EDSM is less than 7% or bigger than 15%.")
		}

		if successfulPercentageBlueFringe < 7 || successfulPercentageBlueFringe > 15 {
			t.Error("The percentage of successful DFAs for Blue-Fringe EDSM is less than 7% or bigger than 15%.")
		}
	}
}

// TestBenchmarkEDSM concurrently benchmarks the performance of the ExhaustiveEDSM(), FastWindowedEDSM(),
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
	winnersExhaustive, winnersWindowed, winnersBlueFringe := 0, 0, 0
	accuraciesExhaustive, accuraciesWindowed, accuraciesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	numberOfStatesExhaustive, numberOfStatesWindowed, numberOfStatesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	durationExhaustive, durationWindowed, durationBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesPerSecExhaustive, mergesPerSecWindowed, mergesPerSecBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	mergesExhaustive, mergesWindowed, mergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()
	validMergesExhaustive, validMergesWindowed, validMergesBlueFringe := util.NewStatsTracker(), util.NewStatsTracker(), util.NewStatsTracker()

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

		resultantDFAExhaustive, resultantDFAWindowed, resultantDFABlueFringe := dfatoolkit.DFA{}, dfatoolkit.DFA{}, dfatoolkit.DFA{}
		searchDataExhaustive, searchDataWindowed, searchDataBlueFringe := dfatoolkit.SearchData{}, dfatoolkit.SearchData{}, dfatoolkit.SearchData{}

		// Exhaustive
		go func() {
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAExhaustive, searchDataExhaustive = dfatoolkit.ExhaustiveEDSM(APTA)
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

		// Exhaustive
		durationExhaustive.Add(searchDataExhaustive.Duration.Seconds())
		mergesPerSecExhaustive.Add(searchDataExhaustive.AttemptedMergesPerSecond())
		accuracy := resultantDFAExhaustive.Accuracy(testingSet)
		accuraciesExhaustive.Add(accuracy)
		numberOfStatesExhaustive.AddInt(len(resultantDFAExhaustive.States))
		mergesExhaustive.AddInt(searchDataExhaustive.AttemptedMergesCount)
		validMergesExhaustive.AddInt(searchDataExhaustive.ValidMergesCount)
		if accuracy >= 0.99 {
			winnersExhaustive++
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

	successfulPercentageExhaustive := (float64(winnersExhaustive) / float64(n)) * 100
	successfulPercentageWindowed := (float64(winnersWindowed) / float64(n)) * 100
	successfulPercentageBlueFringe := (float64(winnersBlueFringe) / float64(n)) * 100

	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Exhaustive Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentageExhaustive)
	PrintBenchmarkInformation(accuraciesExhaustive, numberOfStatesExhaustive, durationExhaustive, mergesPerSecExhaustive, mergesExhaustive, validMergesExhaustive)
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Windowed Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentageWindowed)
	PrintBenchmarkInformation(accuraciesWindowed, numberOfStatesWindowed, durationWindowed, mergesPerSecWindowed, mergesWindowed, validMergesWindowed)
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Blue-Fringe Search")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n\n", successfulPercentageBlueFringe)
	PrintBenchmarkInformation(accuraciesBlueFringe, numberOfStatesBlueFringe, durationBlueFringe, mergesPerSecBlueFringe, mergesBlueFringe, validMergesBlueFringe)
	fmt.Println("--------------------------------------------------------------------------------------------")

	if targetSize == 32 {
		if successfulPercentageExhaustive < 9 || successfulPercentageExhaustive > 15 {
			t.Error("The percentage of successful DFAs for Exhaustive EDSM is less than 9% or bigger than 15%.")
		}

		if successfulPercentageWindowed < 7 || successfulPercentageWindowed > 15 {
			t.Error("The percentage of successful DFAs for Windowed EDSM is less than 7% or bigger than 15%.")
		}

		if successfulPercentageBlueFringe < 7 || successfulPercentageBlueFringe > 15 {
			t.Error("The percentage of successful DFAs for Blue-Fringe EDSM is less than 7% or bigger than 15%.")
		}
	}
}

func PrintBenchmarkInformation(accuracies, numberOfStates, duration, mergesPerSec, merges, validMerges util.StatsTracker){
	// Initialize tabwriter.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 17, 4, 0, '\t', 0)

	_, _ = fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\t\n",  "Minimum", "Maximum", "Average", "Standard Dev")
	_, _ = fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\t\n", "------------", "------------", "------------", "------------")
	_, _ = fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%.2f\t\n", "Accuracy", accuracies.Min(), accuracies.Max(), accuracies.Mean(), accuracies.PopulationStandardDev())
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Number of States", int(numberOfStates.Min()), int(numberOfStates.Max()), int(numberOfStates.Mean()), int(numberOfStates.PopulationStandardDev()))
	_, _ = fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%.2f\t\n", "Duration", duration.Min(), duration.Max(), duration.Mean(), duration.PopulationStandardDev())
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Merges/s", int(math.Round(mergesPerSec.Min())), int(math.Round(mergesPerSec.Max())), int(math.Round(mergesPerSec.Mean())), int(math.Round(mergesPerSec.PopulationStandardDev())))
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Attempted Merges", int(merges.Min()), int(merges.Max()), int(merges.Mean()), int(merges.PopulationStandardDev()))
	_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", "Valid Merges", int(validMerges.Min()), int(validMerges.Max()), int(validMerges.Mean()), int(validMerges.PopulationStandardDev()))

	_ = w.Flush()
}

// -------------------- BENCHMARKS USING STAMINA PROTOCOL --------------------
