package dfatoolkit_test

import (
	"fmt"
	dfalearningtoolkit "github.com/Cherrett/DFA-Learning-Toolkit/core"
	"github.com/Cherrett/DFA-Learning-Toolkit/util"
	"math"
	"math/rand"
	"testing"
	"time"
)

// -------------------- BASIC TESTS --------------------

func TestAbbadingoDatasetFromFile(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	dataset := dfalearningtoolkit.GetDatasetFromAbbadingoFile("../datasets/Abbadingo/Problem S/train.a")

	if len(dataset) != 60000 {
		t.Errorf("Abbadingo Problem S length = %d, want 60000", len(dataset))
	}

	APTA := dataset.GetPTA(true)

	if len(APTA.Alphabet) != 2 {
		t.Errorf("APTA number of symbols = %d, want 2", APTA.Alphabet)
	}

	if len(APTA.States) != 322067 {
		t.Errorf("APTA number of states = %d, want 322067", len(APTA.States))
	}

	if APTA.Depth() != 21 {
		t.Errorf("APTA depth = %d, want 21", APTA.Depth())
	}
}

func TestAbbadingoDFAGeneration(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	numberOfStates := rand.Intn(499) + 1

	AbbadingoDFA := dfalearningtoolkit.AbbadingoDFA(numberOfStates, true)

	if len(AbbadingoDFA.Alphabet) != 2 {
		t.Errorf("AbbadingoDFA number of symbols = %d, want 2", AbbadingoDFA.Alphabet)
	}

	if len(AbbadingoDFA.States) != numberOfStates {
		t.Errorf("AbbadingoDFA number of states = %d, want %d", len(AbbadingoDFA.States), numberOfStates)
	}

	if AbbadingoDFA.Depth() != int(math.Round((2.0*math.Log2(float64(numberOfStates)))-2.0)) {
		t.Errorf("AbbadingoDFA depth = %d, want %d", AbbadingoDFA.Depth(), int(math.Round((2.0*math.Log2(float64(numberOfStates)))-2.0)))
	}
}

func TestAbbadingoDatasetGeneration(t *testing.T) {
	t.Parallel()
	// random seed
	rand.Seed(time.Now().UnixNano())

	numberOfStates := rand.Intn(99) + 1

	AbbadingoDFA, trainingDataset, testingDataset := dfalearningtoolkit.AbbadingoInstance(numberOfStates, false, 35, 0.25)

	trainingDatasetConsistentWithDFA := trainingDataset.ConsistentWithDFA(AbbadingoDFA)
	testingDatasetConsistentWithDFA := testingDataset.ConsistentWithDFA(AbbadingoDFA)

	if !trainingDatasetConsistentWithDFA || !testingDatasetConsistentWithDFA {
		t.Errorf("Expected both training and testing dataset to be consistent with AbbadingoDFA")
	}

	APTA := trainingDataset.GetPTA(true)

	trainingDatasetConsistentWithAPTA := trainingDataset.ConsistentWithDFA(APTA)

	if !trainingDatasetConsistentWithAPTA {
		t.Errorf("Expected training dataset to be consistent with APTA")
	}

	trainingDataset, _ = dfalearningtoolkit.AbbadingoDataset(AbbadingoDFA, 100, 0)

	if !trainingDataset.SymmetricallyStructurallyComplete(AbbadingoDFA) {
		t.Errorf("Expected training dataset to be symmetrically structurally complete with DFA")
	}
}

func TestStaminaDatasetFromFile(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	dataset := dfalearningtoolkit.GetDatasetFromStaminaFile("../datasets/Stamina/96/96_training.txt")

	if len(dataset) != 1093 {
		t.Errorf("Stamina Dataset 96 length = %d, want 1093", len(dataset))
	}

	APTA := dataset.GetPTA(true)

	if len(APTA.Alphabet) != 50 {
		t.Errorf("APTA number of symbols = %d, want 50", APTA.Alphabet)
	}

	if len(APTA.States) != 3503 {
		t.Errorf("APTA number of states = %d, want 3503", len(APTA.States))
	}

	if APTA.Depth() != 53 {
		t.Errorf("APTA depth = %d, want 53", APTA.Depth())
	}
}

func TestStaminaDFAGeneration(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	for _, alphabetSize := range []int{2, 5, 10, 20, 50} {
		StaminaDFA := dfalearningtoolkit.StaminaDFA(alphabetSize, 50)

		if len(StaminaDFA.Alphabet) != alphabetSize {
			t.Errorf("StaminaDFA number of symbols = %d, want %d", StaminaDFA.Alphabet, alphabetSize)
		}

		if len(StaminaDFA.States) < 48 {
			t.Errorf("StaminaDFA number of states = %d, want 48+", len(StaminaDFA.States))
		}
	}
}

func TestStaminaDatasetGeneration(t *testing.T) {
	t.Parallel()
	// random seed
	rand.Seed(time.Now().UnixNano())

	StaminaDFA, trainingDataset, testingDataset := dfalearningtoolkit.StaminaInstance(50, 50, 100, 25000, 2000)

	trainingDatasetConsistentWithDFA := trainingDataset.ConsistentWithDFA(StaminaDFA)
	testingDatasetConsistentWithDFA := testingDataset.ConsistentWithDFA(StaminaDFA)

	if !trainingDatasetConsistentWithDFA || !testingDatasetConsistentWithDFA {
		t.Errorf("Expected both training and testing dataset to be consistent with StaminaDFA")
	}

	APTA := trainingDataset.GetPTA(true)

	trainingDatasetConsistentWithAPTA := trainingDataset.ConsistentWithDFA(APTA)

	if !trainingDatasetConsistentWithAPTA {
		t.Errorf("Expected training dataset to be consistent with APTA")
	}

	// Cover all possible combinations used in the Stamina competition.
	for _, alphabetSize := range []int{2, 5, 10, 20, 50} {
		StaminaDFA = dfalearningtoolkit.StaminaDFA(alphabetSize, 50)

		for _, sparsityPercentage := range []float64{12.5, 25.0, 50.0, 100.0} {
			trainingDataset, testingDataset = dfalearningtoolkit.DefaultStaminaDataset(StaminaDFA, sparsityPercentage)

			if trainingDataset.AcceptingStringInstancesCount() == 0 || trainingDataset.RejectingStringInstancesCount() == 0 {
				t.Errorf("No accepting or rejecting string instances found within training dataset.")
			}

			if testingDataset.AcceptingStringInstancesCount() == 0 || testingDataset.RejectingStringInstancesCount() == 0 {
				t.Errorf("No accepting or rejecting string instances found within testing dataset.")
			}
		}
	}
}

func TestStateMergingAndDFAEquivalence(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	dataset := dfalearningtoolkit.GetDatasetFromAbbadingoFile("../datasets/Abbadingo/Simple/train.a")

	APTA := dataset.GetPTA(false)

	statePartition := APTA.ToStatePartition()
	statePartitionCopy := statePartition.Copy()

	if !statePartitionCopy.MergeStates(2, 4) {
		t.Errorf("Merge should be valid.")
	}

	mergedDFA1 := statePartitionCopy.ToQuotientDFA()

	if !mergedDFA1.IsValidSafe() {
		t.Errorf("State Partition should be valid.")
	}

	statePartitionCopy.RollbackChangesFrom(statePartition)

	if !statePartitionCopy.MergeStates(3, 5) {
		t.Errorf("Merge should be valid.")
	}

	if !statePartitionCopy.MergeStates(2, 4) {
		t.Errorf("Merge should be valid.")
	}

	mergedDFA2 := statePartitionCopy.ToQuotientDFA()

	if !mergedDFA2.IsValidSafe() {
		t.Errorf("State Partition should be valid.")
	}

	if !mergedDFA1.SameAs(mergedDFA2) {
		t.Errorf("Merged DFAs should be equal.")
	}
}

func TestDatasetJSON(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Training set from abbadingo file.
	training1 := dfalearningtoolkit.GetDatasetFromAbbadingoFile("../datasets/Abbadingo/Simple/train.a")
	// Training set from JSON file.
	training2, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Simple/train.json")

	if !valid {
		t.Errorf("Dataset was not read successfully from JSON.")
	}

	if !training1.SameAs(training2) {
		t.Errorf("datasets read not same as each other.")
	}
}

func TestDFAJSON(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	aptaNumberOfStates := []int{845, 2545, 7127}
	aptaDepths := []int{11, 13, 15}

	// Iterate over 3 different sizes of target DFA.
	for i, aptaName := range []int{16, 32, 64}{
		// Read DFA/APTA from JSON.
		APTA, valid := dfalearningtoolkit.DFAFromJSON(fmt.Sprintf("../datasets/TestingAPTAs/%d.json", aptaName))

		if !valid {
			t.Errorf(fmt.Sprintf("DFA was not read successfully from JSON (%d.json).", aptaName))
		}

		if !APTA.IsValidSafe() {
			t.Errorf(fmt.Sprintf("DFA read (%d.json) is not valid.", aptaName))
		}

		if len(APTA.States) != aptaNumberOfStates[i] || len(APTA.Alphabet) != 2 {
			t.Errorf(fmt.Sprintf("DFA read (%d.json) has an invalid number of states or alphabet size.", aptaName))
		}

		if APTA.Depth() != aptaDepths[i] {
			t.Errorf(fmt.Sprintf("DFA read (%d.json) has an invalid depth.", aptaName))
		}
	}
}

func TestVisualisation(t *testing.T) {
	t.Parallel()
	// Training set.
	training := dfalearningtoolkit.GetDatasetFromAbbadingoFile("../datasets/Abbadingo/Simple/train.a")

	test := training.GetPTA(true)

	// Visualisation Examples
	examplesFilenames := []string{"../datasets/Visualisation/test_leftright", "../datasets/Visualisation/test_leftright_ordered",
		"../datasets/Visualisation/test_topdown", "../datasets/Visualisation/test_topdown_ordered"}
	examplesRankByOrder := []bool{false, true, false, true}
	examplesTopDown := []bool{false, false, true, true}

	// To DOT scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".dot"
		test.ToDOT(filePath, dfalearningtoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex])
		if !util.FileExists(filePath) {
			t.Errorf("DFA toDOT failed, %s file not found.", filePath)
		}
	}

	// To PNG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".png"
		if !test.ToPNG(filePath, dfalearningtoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPNG failed, %s file not found.", filePath)
		}
	}

	// To JPG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".jpg"
		if !test.ToJPG(filePath, dfalearningtoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toJPG failed, %s file not found.", filePath)
		}
	}

	// To PDF scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".pdf"
		if !test.ToPDF(filePath, dfalearningtoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPDF failed, %s file not found.", filePath)
		}
	}

	// To SVG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".svg"
		if !test.ToSVG(filePath, dfalearningtoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toSVG failed, %s file not found.", filePath)
		}
	}
}

func TestRPNI(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Read training set from JSON file.
	training, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/train.json")

	if !valid {
		t.Errorf("Training dataset was not read successfully from JSON.")
	}

	// Read testing set from JSON file.
	test, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/test.json")

	if !valid {
		t.Errorf("Testing dataset was not read successfully from JSON.")
	}

	// Run RPNI version on training set.
	resultantDFA, mergeData := dfalearningtoolkit.RPNIFromDataset(training)

	// Get accuracy from resultant DFA on testing set.
	accuracy := resultantDFA.Accuracy(test)

	if mergeData.AttemptedMergesCount != 9831 ||
		mergeData.ValidMergesCount != 108 ||
		len(resultantDFA.States) != 107 ||
		accuracy != 0.5188888888888888 {
		t.Errorf("Discrepancies found in result of RPNI.")
	}
}

func TestExhaustiveEDSM(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Part 1 - Training and Testing Sets from file.

	// Read training set from JSON file.
	training, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/train.json")

	if !valid {
		t.Errorf("Training dataset was not read successfully from JSON.")
	}

	// Read testing set from JSON file.
	test, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/test.json")

	if !valid {
		t.Errorf("Testing dataset was not read successfully from JSON.")
	}

	// Run Exhaustive EDSM on training set.
	resultantDFA, mergeData := dfalearningtoolkit.ExhaustiveEDSMFromDataset(training)

	// Get accuracy from resultant DFA on testing set.
	accuracy := resultantDFA.Accuracy(test)

	// Confirm merge data and DFA values.
	if mergeData.AttemptedMergesCount != 15426663 ||
		mergeData.ValidMergesCount != 13591962 ||
		len(resultantDFA.States) != 31 ||
		accuracy != 0.9933333333333333 {
		t.Errorf("Discrepancies found in result of exhaustive EDSM (Part 1).")
	}

	// Part 2 - APTAs from files.

	// Slices of expected values.
	resultantNumberOfStates := []int{25, 31}
	resultantAttemptedMergesCount := []int{1060866, 15426663}
	resultantValidMergesCount := []int{868076, 13591962}

	// Iterate over 2 different sizes of target DFA.
	for i, dfaSize := range []int{16, 32} {
		// Read APTA from JSON file.
		APTA, valid := dfalearningtoolkit.DFAFromJSON(fmt.Sprintf("../datasets/TestingAPTAs/%d.json", dfaSize))

		if !valid {
			t.Errorf("APTA was not read successfully from JSON.")
		}

		// Run Exhaustive EDSM on APTA.
		resultantDFA, mergeData = dfalearningtoolkit.ExhaustiveEDSM(APTA)

		if len(resultantDFA.States) != resultantNumberOfStates[i] ||
			mergeData.AttemptedMergesCount != resultantAttemptedMergesCount[i] ||
			mergeData.ValidMergesCount != resultantValidMergesCount[i] {
			t.Errorf("Discrepancies found in result of exhaustive EDSM (Part 2).")
		}
	}
}

func TestWindowedEDSM(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Part 1 - Training and Testing Sets from file.

	// Read training set from JSON file.
	training, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/train.json")

	if !valid {
		t.Errorf("Training dataset was not read successfully from JSON.")
	}

	// Read testing set from JSON file.
	test, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/test.json")

	if !valid {
		t.Errorf("Testing dataset was not read successfully from JSON.")
	}

	// Run Windowed EDSM on training set.
	resultantDFA, mergeData := dfalearningtoolkit.WindowedEDSMFromDataset(training, 32*2, 2.0)

	// Get accuracy from resultant DFA on testing set.
	accuracy := resultantDFA.Accuracy(test)

	// Confirm merge data and DFA values.
	if mergeData.AttemptedMergesCount != 62786 ||
		mergeData.ValidMergesCount != 14018 ||
		len(resultantDFA.States) != 31 ||
		accuracy != 0.9933333333333333 {
		t.Errorf("Discrepancies found in result of windowed EDSM (Part 1).")
	}

	// Part 2 - APTAs from files.

	// Slices of expected values.
	resultantNumberOfStates := []int{25, 31, 208}
	resultantAttemptedMergesCount := []int{13078, 62786, 4819261}
	resultantValidMergesCount := []int{1861, 14018, 420270}

	// Iterate over 3 different sizes of target DFA.
	for i, dfaSize := range []int{16, 32, 64} {
		// Read APTA from JSON file.
		APTA, valid := dfalearningtoolkit.DFAFromJSON(fmt.Sprintf("../datasets/TestingAPTAs/%d.json", dfaSize))

		if !valid {
			t.Errorf("APTA was not read successfully from JSON.")
		}

		// Run Windowed EDSM on APTA.
		resultantDFA, mergeData = dfalearningtoolkit.WindowedEDSM(APTA, dfaSize*2, 2.0)

		if len(resultantDFA.States) != resultantNumberOfStates[i] ||
			mergeData.AttemptedMergesCount != resultantAttemptedMergesCount[i] ||
			mergeData.ValidMergesCount != resultantValidMergesCount[i] {
			t.Errorf("Discrepancies found in result of windowed EDSM (Part 2).")
		}
	}
}

func TestBlueFringeEDSM(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Part 1 - Training and Testing Sets from file.

	// Read training set from JSON file.
	training, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/train.json")

	if !valid {
		t.Errorf("Training dataset was not read successfully from JSON.")
	}

	// Read testing set from JSON file.
	test, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/test.json")

	if !valid {
		t.Errorf("Testing dataset was not read successfully from JSON.")
	}

	// Run Blue-Fringe EDSM on training set.
	resultantDFA, mergeData := dfalearningtoolkit.BlueFringeEDSMFromDataset(training)

	// Get accuracy from resultant DFA on testing set.
	accuracy := resultantDFA.Accuracy(test)

	// Confirm merge data and DFA values.
	if mergeData.AttemptedMergesCount != 12215 ||
		mergeData.ValidMergesCount != 789 ||
		len(resultantDFA.States) != 31 ||
		accuracy != 0.9933333333333333 {
		t.Errorf("Discrepancies found in result of blue-fringe EDSM (Part 1).")
	}

	// Part 2 - APTAs from files.

	// Slices of expected values.
	resultantNumberOfStates := []int{25, 31, 200}
	resultantAttemptedMergesCount := []int{5685, 12215, 1477532}
	resultantValidMergesCount := []int{764, 789, 73290}

	// Iterate over 3 different sizes of target DFA.
	for i, dfaSize := range []int{16, 32, 64} {
		// Read APTA from JSON file.
		APTA, valid := dfalearningtoolkit.DFAFromJSON(fmt.Sprintf("../datasets/TestingAPTAs/%d.json", dfaSize))

		if !valid {
			t.Errorf("APTA was not read successfully from JSON.")
		}

		// Run Blue-Fringe EDSM on APTA.
		resultantDFA, mergeData = dfalearningtoolkit.BlueFringeEDSM(APTA)

		if len(resultantDFA.States) != resultantNumberOfStates[i] ||
			mergeData.AttemptedMergesCount != resultantAttemptedMergesCount[i] ||
			mergeData.ValidMergesCount != resultantValidMergesCount[i] {
			t.Errorf("Discrepancies found in result of blue-fringe EDSM (Part 2).")
		}
	}
}

func TestAutomataTeams(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Read training set from JSON file.
	training, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/train.json")

	if !valid {
		t.Errorf("Training dataset was not read successfully from JSON.")
	}

	// Read testing set from JSON file.
	test, valid := dfalearningtoolkit.DatasetFromJSON("../datasets/Abbadingo/Generated/test.json")

	if !valid {
		t.Errorf("Testing dataset was not read successfully from JSON.")
	}

	// Run AutomataTeams version on training set.
	teamOfAutomata := dfalearningtoolkit.AutomataTeamsFromDataset(training, 81)
	fairVoteAccuracy := teamOfAutomata.FairVoteAccuracy(test)
	weightedVoteAccuracy := teamOfAutomata.WeightedVoteAccuracy(test)
	betterHalfWeightedVoteAccuracy := teamOfAutomata.BetterHalfWeightedVoteAccuracy(test)

	if fairVoteAccuracy > weightedVoteAccuracy || weightedVoteAccuracy > betterHalfWeightedVoteAccuracy {
		t.Errorf("Discrepancies found in result of AutomataTeams.")
	}
}
