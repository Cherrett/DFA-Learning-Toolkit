package dfatoolkit_test

import (
	"DFA_Toolkit/DFA_Toolkit"
	"DFA_Toolkit/DFA_Toolkit/util"
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

	dataset := dfatoolkit.GetDatasetFromAbbadingoFile("../../Datasets/Abbadingo/Problem S/train.a")

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

	AbbadingoDFA := dfatoolkit.AbbadingoDFA(numberOfStates, true)

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

	AbbadingoDFA, trainingDataset, testingDataset := dfatoolkit.AbbadingoInstance(numberOfStates, false, 35, 0.25)

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

	trainingDataset, _ = dfatoolkit.AbbadingoDataset(AbbadingoDFA, 100, 0)

	if !trainingDataset.SymmetricallyStructurallyComplete(AbbadingoDFA) {
		t.Errorf("Expected training dataset to be symmetrically structurally complete with DFA")
	}
}

func TestStaminaDatasetFromFile(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	dataset := dfatoolkit.GetDatasetFromStaminaFile("../../Datasets/Stamina/96/96_training.txt")

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
		StaminaDFA := dfatoolkit.StaminaDFA(alphabetSize, 50)

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

	StaminaDFA, trainingDataset, testingDataset := dfatoolkit.StaminaInstance(50, 50, 100, 25000, 2000)

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

	trainingDataset, _ = dfatoolkit.DefaultStaminaDataset(StaminaDFA, 100)

	if !trainingDataset.SymmetricallyStructurallyComplete(StaminaDFA) {
		t.Errorf("Expected training dataset to be symmetrically structurally complete with DFA")
	}

	// Cover all possible combinations used in the Stamina competition.
	for _, alphabetSize := range []int{2, 5, 10, 20, 50} {
		StaminaDFA = dfatoolkit.StaminaDFA(alphabetSize, 50)

		for _, sparsityPercentage := range []float64{12.5, 25.0, 50.0, 100.0} {
			trainingDataset, testingDataset = dfatoolkit.DefaultStaminaDataset(StaminaDFA, sparsityPercentage)
		}

		if trainingDataset.AcceptingStringInstancesCount() == 0 || trainingDataset.RejectingStringInstancesCount() == 0 {
			t.Errorf("No accepting or rejecting string instances found within training dataset.")
		}

		if testingDataset.AcceptingStringInstancesCount() == 0 || testingDataset.RejectingStringInstancesCount() == 0 {
			t.Errorf("No accepting or rejecting string instances found within testing dataset.")
		}
	}
}

func TestStateMergingAndDFAEquivalence(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	dataset := dfatoolkit.GetDatasetFromAbbadingoFile("../../Datasets/Abbadingo/Simple/train.a")

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

	if !mergedDFA1.Equal(mergedDFA2) {
		t.Errorf("Merged DFAs should be equal.")
	}
}

func TestDatasetJSON(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Training set from abbadingo file.
	training1 := dfatoolkit.GetDatasetFromAbbadingoFile("../../Datasets/Abbadingo/Simple/train.a")
	// Training set from JSON file.
	training2, valid := dfatoolkit.DatasetFromJSON("../../Datasets/Abbadingo/Simple/train.json")

	if !valid {
		t.Errorf("Dataset was not read successfuly from JSON.")
	}

	if !training1.SameAs(training2) {
		t.Errorf("Datasets read not same as each other.")
	}
}

func TestDFAJSON(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Read DFAs from JSON.
	APTA16, valid := dfatoolkit.DFAFromJSON("../../TestingAPTAs/16.json")

	if !valid {
		t.Errorf("DFA was not read successfuly from JSON (16.json).")
	}

	if len(APTA16.States) != 845 || len(APTA16.Alphabet) != 2 {
		t.Errorf("DFA read (16.json) has an invalid number of states or alphabet size.")
	}

	APTA32, valid := dfatoolkit.DFAFromJSON("../../TestingAPTAs/32.json")

	if !valid {
		t.Errorf("DFA was not read successfuly from JSON (32.json).")
	}

	if len(APTA32.States) != 2545 || len(APTA32.Alphabet) != 2 {
		t.Errorf("DFA read (32.json) has an invalid number of states or alphabet size.")
	}

	APTA64, valid := dfatoolkit.DFAFromJSON("../../TestingAPTAs/64.json")

	if !valid {
		t.Errorf("DFA was not read successfuly from JSON (64.json).")
	}

	if len(APTA64.States) != 7127 || len(APTA64.Alphabet) != 2 {
		t.Errorf("DFA read (64.json) has an invalid number of states or alphabet size.")
	}

	if !APTA16.IsValidSafe() || !APTA32.IsValidSafe() || !APTA64.IsValidSafe() {
		t.Errorf("One of the read DFAs is not valid.")
	}
}

func TestVisualisation(t *testing.T) {
	t.Parallel()
	// Training set.
	training := dfatoolkit.GetDatasetFromAbbadingoFile("../../Datasets/Abbadingo/Simple/train.a")

	test := training.GetPTA(true)

	// Visualisation Examples
	examplesFilenames := []string{"../../Visualisation/test_leftright", "../../Visualisation/test_leftright_ordered",
		"../../Visualisation/test_topdown", "../../Visualisation/test_topdown_ordered"}
	examplesRankByOrder := []bool{false, true, false, true}
	examplesTopDown := []bool{false, false, true, true}

	// To DOT scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".dot"
		test.ToDOT(filePath, dfatoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex])
		if !util.FileExists(filePath) {
			t.Errorf("DFA toDOT failed, %s file not found.", filePath)
		}
	}

	// To PNG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".png"
		if !test.ToPNG(filePath, dfatoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPNG failed, %s file not found.", filePath)
		}
	}

	// To JPG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".jpg"
		if !test.ToJPG(filePath, dfatoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toJPG failed, %s file not found.", filePath)
		}
	}

	// To PDF scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".pdf"
		if !test.ToPDF(filePath, dfatoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPDF failed, %s file not found.", filePath)
		}
	}

	// To SVG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex] + ".svg"
		if !test.ToSVG(filePath, dfatoolkit.SymbolAlphabetMappingAbbadingo, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]) {
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
	training, valid := dfatoolkit.DatasetFromJSON("../../Datasets/Abbadingo/Generated/train.json")

	if !valid {
		t.Errorf("Training dataset was not read successfuly from JSON.")
	}

	// Read testing set from JSON file.
	test, valid := dfatoolkit.DatasetFromJSON("../../Datasets/Abbadingo/Generated/test.json")

	if !valid {
		t.Errorf("Testing dataset was not read successfuly from JSON.")
	}

	// Run RPNI version on training set.
	resultantDFA, mergeData := dfatoolkit.RPNIFromDataset(training)

	// Get accuracy from resultant DFA on testing set.
	accuracy := resultantDFA.Accuracy(test)

	if mergeData.AttemptedMergesCount != 3158232 ||
		mergeData.ValidMergesCount != 109 ||
		len(resultantDFA.States) != 108 ||
		accuracy != 0.5222222222222223 {
		t.Errorf("Discrepancies found in result of RPNI.")
	}
}

func TestEDSM(t *testing.T) {
	t.Parallel()
	// Random Seed.
	rand.Seed(time.Now().UnixNano())

	// Read training set from JSON file.
	training, valid := dfatoolkit.DatasetFromJSON("../../Datasets/Abbadingo/Generated/train.json")

	if !valid {
		t.Errorf("Training dataset was not read successfuly from JSON.")
	}

	// Read testing set from JSON file.
	test, valid := dfatoolkit.DatasetFromJSON("../../Datasets/Abbadingo/Generated/test.json")

	if !valid {
		t.Errorf("Testing dataset was not read successfuly from JSON.")
	}

	// Run all EDSM version on training set.
	exhaustiveResultantDFA, exhaustiveMergeData := dfatoolkit.ExhaustiveEDSMFromDataset(training)
	windowedResultantDFA, windowedMergeData := dfatoolkit.WindowedEDSMFromDataset(training, 64, 2.0)
	blueFringeResultantDFA, blueFringeMergeData := dfatoolkit.BlueFringeEDSMFromDataset(training)

	// Get accuracy from resultant DFAs on testing set.
	exhaustiveAccuracy := exhaustiveResultantDFA.Accuracy(test)
	windowedAccuracy := windowedResultantDFA.Accuracy(test)
	blueFringeAccuracy := blueFringeResultantDFA.Accuracy(test)

	if exhaustiveMergeData.AttemptedMergesCount != 15426663 ||
		exhaustiveMergeData.ValidMergesCount != 13591962 ||
		len(exhaustiveResultantDFA.States) != 31 ||
		exhaustiveAccuracy != 0.9933333333333333 {
		t.Errorf("Discrepancies found in result of exhaustive EDSM.")
	}

	if windowedMergeData.AttemptedMergesCount != 62786 ||
		windowedMergeData.ValidMergesCount != 14018 ||
		len(windowedResultantDFA.States) != 31 ||
		windowedAccuracy != 0.9933333333333333 {
		t.Errorf("Discrepancies found in result of windowed EDSM.")
	}

	if blueFringeMergeData.AttemptedMergesCount != 12215 ||
		blueFringeMergeData.ValidMergesCount != 789 ||
		len(blueFringeResultantDFA.States) != 31 ||
		blueFringeAccuracy != 0.9933333333333333 {
		t.Errorf("Discrepancies found in result of blue-fringe EDSM.")
	}
}
