package DFA_Toolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// -------------------- BASIC TESTS --------------------

func TestAbbadingoDFAFromFile(t *testing.T) {
	t.Parallel()
	dataset := GetDatasetFromAbbadingoFile("../AbbadingoDatasets/dataset4/train.a")
	if len(dataset) != 60000{
		t.Errorf("Dataset4 length = %d, want 60000", len(dataset))
	}

	APTA := dataset.GetPTA(true)
	if len(APTA.SymbolMap) != 2{
		t.Errorf("APTA number of symbols = %d, want 2", len(APTA.SymbolMap))
	}
	if len(APTA.States) != 322067{
		t.Errorf("APTA number of states = %d, want 322067", len(APTA.States))
	}
	if APTA.Depth() != 21{
		t.Errorf("APTA depth = %d, want 21", APTA.Depth())
	}
}

func TestAbbadingoDFAGeneration(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())
	numberOfStates := rand.Intn(499) + 1

	AbbadingoDFA := AbbadingoDFA(numberOfStates, true)
	if len(AbbadingoDFA.SymbolMap) != 2{
		t.Errorf("AbbadingoDFA number of symbols = %d, want 2", len(AbbadingoDFA.SymbolMap))
	}
	if len(AbbadingoDFA.States) != numberOfStates{
		t.Errorf("AbbadingoDFA number of states = %d, want %d", len(AbbadingoDFA.States), numberOfStates)
	}
	if AbbadingoDFA.Depth() != int(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0)){
		t.Errorf("AbbadingoDFA depth = %d, want %d", AbbadingoDFA.Depth(), int(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0)))
	}
}

func TestAbbadingoDatasetGeneration(t *testing.T){
	t.Parallel()
	// random seed
	rand.Seed(time.Now().UnixNano())
	numberOfStates := rand.Intn(99) + 1

	AbbadingoDFA := AbbadingoDFA(numberOfStates, false)

	trainingDataset, testingDataset := AbbadingoDataset(AbbadingoDFA, 35, 0.25)

	trainingDatasetConsistentWithDFA := trainingDataset.ConsistentWithDFA(AbbadingoDFA)
	testingDatasetConsistentWithDFA := testingDataset.ConsistentWithDFA(AbbadingoDFA)

	if !trainingDatasetConsistentWithDFA || !testingDatasetConsistentWithDFA{
		t.Errorf("Expected both training and testing dataset to be consistent with AbbadingoDFA")
	}

	APTA := trainingDataset.GetPTA(true)

	trainingDatasetConsistentWithAPTA := trainingDataset.ConsistentWithDFA(APTA)

	if !trainingDatasetConsistentWithAPTA{
		t.Errorf("Expected training dataset to be consistent with APTA")
	}
}

func TestStateMergingAndDFAEquivalence(t *testing.T){
	t.Parallel()
	dataset := GetDatasetFromAbbadingoFile("../AbbadingoDatasets/dataset1/train.a")
	APTA := dataset.GetPTA(false)

	statePartition := APTA.ToStatePartition()
	statePartitionCopy := statePartition.Copy()

	if !statePartitionCopy.MergeStates(APTA, 2, 4){
		t.Errorf("Merge should be valid.")
	}
	valid1, mergedDFA1 := statePartitionCopy.ToDFA(APTA)
	if !valid1{
		t.Errorf("State Partition should be valid.")
	}
	statePartitionCopy.RollbackChanges(statePartition)

	if !statePartitionCopy.MergeStates(APTA, 3, 5){
		t.Errorf("Merge should be valid.")
	}
	if !statePartitionCopy.MergeStates(APTA, 2, 4){
		t.Errorf("Merge should be valid.")
	}
	valid2, mergedDFA2 := statePartitionCopy.ToDFA(APTA)
	if !valid2{
		t.Errorf("State Partition should be valid.")
	}

	if !mergedDFA1.Equal(mergedDFA2){
		t.Errorf("Merged DFAs should be equal.")
	}
}

func TestVisualisation(t *testing.T){
	t.Parallel()
	// Training set.
	training := GetDatasetFromAbbadingoFile("../AbbadingoDatasets/test.txt")

	test := training.GetPTA(true)

	// Visualisation Examples
	examplesFilenames := []string{"../Visualisation/test_leftright", "../Visualisation/test_leftright_ordered",
		"../Visualisation/test_topdown", "../Visualisation/test_topdown_ordered"}
	examplesRankByOrder := []bool{false, true, false, true}
	examplesTopDown := []bool{false, false, true, true}

	// To DOT scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".dot"
		test.ToDOT(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex])
		if !util.FileExists(filePath) {
			t.Errorf("DFA toDOT failed, %s file not found.", filePath)
		}
	}

	// To PNG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".png"
		test.ToPNG(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex])
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPNG failed, %s file not found.", filePath)
		}
	}

	// To JPG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".jpg"
		test.ToJPG(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex])
		if !util.FileExists(filePath) {
			t.Errorf("DFA toJPG failed, %s file not found.", filePath)
		}
	}

	// To PDF scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".pdf"
		test.ToPDF(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex])
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPDF failed, %s file not found.", filePath)
		}
	}

	// To SVG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".svg"
		test.ToSVG(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex])
		if !util.FileExists(filePath) {
			t.Errorf("DFA toSVG failed, %s file not found.", filePath)
		}
	}
}

// -------------------- BENCHMARKS --------------------

// TestBenchmarkDetMerge benchmarks the performance of the MergeStates() function.
func TestBenchmarkDetMerge(t *testing.T){
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// These are target DFA sizes we will test.
	dfaSizes := []int{16, 32, 64, 128}
	// These are the training set sizes we will test.
	trainingSetSizes := []int{230, 607, 1521, 4382}

	// Benchmark over the problem instances.
	for i := range dfaSizes {
		targetSize := dfaSizes[i]
		trainingSetSize := trainingSetSizes[i]

		// Create a target DFA.
		target := AbbadingoDFA(targetSize, true)

		// Training set.
		training, _ := AbbadingoDatasetExact(target, trainingSetSize, 0)

		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("BENCHMARK %d (Target: %d states, Training: %d strings\n", i+1, targetSize, len(training))
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
				if snapshot.MergeStates(apta, i, j){
					validMerges++
					//snapshot.LabelledBlocksCount(apta)
					//snapshot.BlocksCount()
					//print(temp, temp2)
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
func TestBenchmarkGreedyEDSM(t *testing.T){
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	winners := 0
	totalAccuracies := util.NewMinMaxAvg()
	totalNumberOfStates := util.NewMinMaxAvg()
	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)
		start := time.Now()

		// Create a target DFA.
		target := AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA := GreedyEDSMFromDataset(trainingSet)
		accuracy := resultantDFA.Accuracy(testingSet)

		totalAccuracies.Add(accuracy)
		totalNumberOfStates.Add(float64(resultantDFA.AllStatesCount()))

		if accuracy >= 0.99{
			winners++
		}

		fmt.Printf("Duration: %.2fs\n\n", time.Since(start).Seconds())
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuracies.Min(), totalAccuracies.Max(), totalAccuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStates.Min(), totalNumberOfStates.Max(), totalNumberOfStates.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentage < 0.10 || successfulPercentage > 0.15{
		t.Error("The percentage of successful DFAs is less than 0.10 or bigger than 0.15.")
	}
}

// TestBenchmarkWindowedEDSM benchmarks the performance of the WindowedEDSMFromDataset() function.
func TestBenchmarkWindowedEDSM(t *testing.T){
	// Random Seed.
	// rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	winners := 0
	totalAccuracies := util.NewMinMaxAvg()
	totalNumberOfStates := util.NewMinMaxAvg()
	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)
		start := time.Now()

		// Create a target DFA.
		target := AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA := WindowedEDSMFromDataset(trainingSet, targetSize*2, 2.0)
		accuracy := resultantDFA.Accuracy(testingSet)

		totalAccuracies.Add(accuracy)
		totalNumberOfStates.Add(float64(resultantDFA.AllStatesCount()))

		if accuracy >= 0.99{
			winners++
		}

		fmt.Printf("Duration: %.2fs\n\n", time.Since(start).Seconds())
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuracies.Min(), totalAccuracies.Max(), totalAccuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStates.Min(), totalNumberOfStates.Max(), totalNumberOfStates.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentage < 0.09 || successfulPercentage > 0.15{
		t.Error("The percentage of successful DFAs is less than 0.09 or bigger than 0.15.")
	}
}

// TestBenchmarkBlueFringeEDSM benchmarks the performance of the BlueFringeEDSMFromDataset() function.
func TestBenchmarkBlueFringeEDSM(t *testing.T){
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	winners := 0
	totalAccuracies := util.NewMinMaxAvg()
	totalNumberOfStates := util.NewMinMaxAvg()
	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)
		start := time.Now()

		// Create a target DFA.
		target := AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := AbbadingoDatasetExact(target, 607, 1800)

		resultantDFA := BlueFringeEDSMFromDataset(trainingSet)
		accuracy := resultantDFA.Accuracy(testingSet)

		totalAccuracies.Add(accuracy)
		totalNumberOfStates.Add(float64(resultantDFA.AllStatesCount()))

		if accuracy >= 0.99{
			winners++
		}

		fmt.Printf("Duration: %.2fs\n\n", time.Since(start).Seconds())
	}

	successfulPercentage := float64(winners) / float64(n)
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentage)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuracies.Min(), totalAccuracies.Max(), totalAccuracies.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStates.Min(), totalNumberOfStates.Max(), totalNumberOfStates.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentage < 0.07 || successfulPercentage > 0.15{
		t.Error("The percentage of successful DFAs is less than 0.07 or bigger than 0.15.")
	}
}

func TestBenchmarkEDSM(t *testing.T) {
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Number of iterations.
	n := 128
	// Target size.
	targetSize := 32

	winnersGreedy, winnersWindowed, winnersBlueFringe := 0, 0, 0
	totalAccuraciesGreedy, totalAccuraciesWindowed, totalAccuraciesBlueFringe := util.NewMinMaxAvg(), util.NewMinMaxAvg(), util.NewMinMaxAvg()
	totalNumberOfStatesGreedy, totalNumberOfStatesWindowed, totalNumberOfStatesBlueFringe := util.NewMinMaxAvg(), util.NewMinMaxAvg(), util.NewMinMaxAvg()
	for i := 0; i < n; i++ {
		fmt.Printf("BENCHMARK %d/%d\n", i+1, n)

		// Create a target DFA.
		target := AbbadingoDFA(targetSize, true)

		// Training testing sets.
		trainingSet, testingSet := AbbadingoDatasetExact(target, 607, 1800)

		// Construct an APTA from training dataset.
		APTA := trainingSet.GetPTA(true)

		// Create wait group
		var wg sync.WaitGroup
		// Add 3 EDSM types to wait group.
		wg.Add(3)

		var resultantDFAGreedy DFA
		var resultantDFAWindowed DFA
		var resultantDFABlueFringe DFA

		go func(){
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAGreedy = GreedyEDSM(APTA)
		}()

		go func(){
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFAWindowed = WindowedEDSM(APTA, targetSize*2, 2.0)
		}()

		go func(){
			// Decrement 1 from wait group.
			defer wg.Done()
			resultantDFABlueFringe = BlueFringeEDSM(APTA)
		}()

		// Wait for all go routines within wait
		// group to finish executing.
		wg.Wait()

		accuracyGreedy := resultantDFAGreedy.Accuracy(testingSet)
		accuracyWindowed := resultantDFAWindowed.Accuracy(testingSet)
		accuracyBlueFringe := resultantDFABlueFringe.Accuracy(testingSet)

		totalAccuraciesGreedy.Add(accuracyGreedy)
		totalAccuraciesWindowed.Add(accuracyWindowed)
		totalAccuraciesBlueFringe.Add(accuracyBlueFringe)

		totalNumberOfStatesGreedy.Add(float64(resultantDFAGreedy.AllStatesCount()))
		totalNumberOfStatesWindowed.Add(float64(resultantDFAWindowed.AllStatesCount()))
		totalNumberOfStatesBlueFringe.Add(float64(resultantDFABlueFringe.AllStatesCount()))

		if accuracyGreedy >= 0.99{
			winnersGreedy++
		}

		if accuracyWindowed >= 0.99{
			winnersWindowed++
		}

		if accuracyBlueFringe >= 0.99{
			winnersBlueFringe++
		}

		fmt.Printf("Greedy: Accuracy: %.2f, Number of States %d\n", accuracyGreedy, resultantDFAGreedy.AllStatesCount())
		fmt.Printf("Windowed: Accuracy: %.2f, Number of States %d\n", accuracyWindowed, resultantDFAWindowed.AllStatesCount())
		fmt.Printf("Blue-Fringe: Accuracy: %.2f, Number of States %d\n", accuracyBlueFringe, resultantDFABlueFringe.AllStatesCount())
		fmt.Print("-----------------------------------------------------------------------------\n")
	}

	successfulPercentageGreedy := float64(winnersGreedy) / float64(n)
	fmt.Print("Greedy Search:\n")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageGreedy)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuraciesGreedy.Min(), totalAccuraciesGreedy.Max(), totalAccuraciesGreedy.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStatesGreedy.Min(), totalNumberOfStatesGreedy.Max(), totalNumberOfStatesGreedy.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	successfulPercentageWindowed := float64(winnersWindowed) / float64(n)
	fmt.Printf("Windowed Search:\n")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageWindowed)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuraciesWindowed.Min(), totalAccuraciesWindowed.Max(), totalAccuraciesWindowed.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStatesWindowed.Min(), totalNumberOfStatesWindowed.Max(), totalNumberOfStatesWindowed.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	successfulPercentageBlueFringe := float64(winnersBlueFringe) / float64(n)
	fmt.Printf("Blue-Fringe Search:\n")
	fmt.Printf("Percentage of 0.99 <= Accuracy: %.2f%%\n", successfulPercentageBlueFringe)
	fmt.Printf("Minimum Accuracy: %.2f Maximum Accuracy: %.2f Average Accuracy: %.2f\n", totalAccuraciesBlueFringe.Min(), totalAccuraciesBlueFringe.Max(), totalAccuraciesBlueFringe.Avg())
	fmt.Printf("Minimum States: %.2f Maximum States: %.2f Average States: %.2f\n", totalNumberOfStatesBlueFringe.Min(), totalNumberOfStatesBlueFringe.Max(), totalNumberOfStatesBlueFringe.Avg())
	fmt.Print("-----------------------------------------------------------------------------\n\n")

	if successfulPercentageGreedy < 0.09 || successfulPercentageGreedy > 0.15{
		t.Error("The percentage of successful DFAs for Greedy EDSM is less than 0.09 or bigger than 0.15.")
	}

	if successfulPercentageWindowed < 0.09 || successfulPercentageWindowed > 0.15{
		t.Error("The percentage of successful DFAs for Windowed EDSM is less than 0.09 or bigger than 0.15.")
	}

	if successfulPercentageBlueFringe < 0.07 || successfulPercentageBlueFringe > 0.15{
		t.Error("The percentage of successful DFAs for Blue-Fringe EDSM is less than 0.07 or bigger than 0.15.")
	}
}