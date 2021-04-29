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

func TestAbbadingoDFAFromFile(t *testing.T) {
	t.Parallel()
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

	if !trainingDataset.SymmetricallyStructurallyComplete(AbbadingoDFA){
		t.Errorf("Expected training dataset to be symmetrically structurally complete with DFA")
	}
}

func TestStaminaDFAFromFile(t *testing.T) {
	t.Parallel()
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

	StaminaDFA := dfatoolkit.StaminaDFA(50, 50)
	if len(StaminaDFA.Alphabet) != 50 {
		t.Errorf("StaminaDFA number of symbols = %d, want 50", StaminaDFA.Alphabet)
	}
	if len(StaminaDFA.States) != 50 {
		t.Errorf("StaminaDFA number of states = %d, want %d", len(StaminaDFA.States), 50)
	}
}

func TestStateMergingAndDFAEquivalence(t *testing.T) {
	t.Parallel()
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
