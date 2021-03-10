package DFA_Toolkit

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

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
	if APTA.GetDepth() != 21{
		t.Errorf("APTA depth = %d, want 21", APTA.GetDepth())
	}
}

func TestAbbadingoDFAGeneration(t *testing.T) {
	t.Parallel()
	// random seed
	rand.Seed(time.Now().UnixNano())
	numberOfStates := rand.Intn(499) + 1

	AbbadingoDFA := AbbadingoDFA(numberOfStates, true)
	if len(AbbadingoDFA.SymbolMap) != 2{
		t.Errorf("AbbadingoDFA number of symbols = %d, want 2", len(AbbadingoDFA.SymbolMap))
	}
	if len(AbbadingoDFA.States) != numberOfStates{
		t.Errorf("AbbadingoDFA number of states = %d, want %d", len(AbbadingoDFA.States), numberOfStates)
	}
	if AbbadingoDFA.GetDepth() != int(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0)){
		t.Errorf("AbbadingoDFA depth = %d, want %d", AbbadingoDFA.GetDepth(), int(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0)))
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

func TestBenchmarkDetMerge(t *testing.T){
	// random seed
	//rand.Seed(time.Now().UnixNano())
	BenchmarkDetMerge()
}

func TestBenchmarkEDSM(t *testing.T){
	// random seed
	//rand.Seed(time.Now().UnixNano())
	BenchmarkEDSM(128)
}