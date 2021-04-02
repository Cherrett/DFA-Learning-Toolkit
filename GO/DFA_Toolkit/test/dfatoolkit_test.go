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
	dataset := dfatoolkit.GetDatasetFromAbbadingoFile("../../AbbadingoDatasets/dataset4/train.a")
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

	AbbadingoDFA := dfatoolkit.AbbadingoDFA(numberOfStates, true)
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

	AbbadingoDFA := dfatoolkit.AbbadingoDFA(numberOfStates, false)

	trainingDataset, testingDataset := dfatoolkit.AbbadingoDataset(AbbadingoDFA, 35, 0.25)

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
	dataset := dfatoolkit.GetDatasetFromAbbadingoFile("../../AbbadingoDatasets/dataset1/train.a")
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

func TestDatasetJSON(t *testing.T){
	t.Parallel()

	// Training set from abbadingo file.
	training1 := dfatoolkit.GetDatasetFromAbbadingoFile("../../AbbadingoDatasets/dataset1/train.a")
	// Training set from JSON file.
	training2, valid := dfatoolkit.DatasetFromJSON("../../AbbadingoDatasets/dataset1/train.json")

	if !valid{
		t.Errorf("Dataset was not read successfuly from JSON.")
	}

	if !training1.SameAs(training2){
		t.Errorf("Datasets read not same as each other.")
	}
}

func TestVisualisation(t *testing.T){
	t.Parallel()
	// Training set.
	training := dfatoolkit.GetDatasetFromAbbadingoFile("../../AbbadingoDatasets/dataset1/train.a")

	test := training.GetPTA(true)

	// Visualisation Examples
	examplesFilenames := []string{"../../Visualisation/test_leftright", "../../Visualisation/test_leftright_ordered",
		"../../Visualisation/test_topdown", "../../Visualisation/test_topdown_ordered"}
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
		if !test.ToPNG(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]){
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPNG failed, %s file not found.", filePath)
		}
	}

	// To JPG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".jpg"
		if !test.ToJPG(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]){
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toJPG failed, %s file not found.", filePath)
		}
	}

	// To PDF scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".pdf"
		if !test.ToPDF(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]){
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toPDF failed, %s file not found.", filePath)
		}
	}

	// To SVG scenario
	for exampleIndex := range examplesFilenames {
		filePath := examplesFilenames[exampleIndex]+".svg"
		if !test.ToSVG(filePath, examplesRankByOrder[exampleIndex], examplesTopDown[exampleIndex]){
			t.Errorf("GraphViz Error. Probabbly not installed properly.")
		}
		if !util.FileExists(filePath) {
			t.Errorf("DFA toSVG failed, %s file not found.", filePath)
		}
	}
}