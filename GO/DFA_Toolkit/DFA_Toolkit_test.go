package DFA_Toolkit

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestAbbadingoDFAFromFile(t *testing.T) {
	dataset := GetDatasetFromAbbadingoFile("../AbbadingoDatasets/dataset4/train.a")
	if len(dataset) != 60000{
		t.Errorf("Dataset4 length = %d, want 60000", len(dataset))
	}

	APTA := GetPTAFromDataset(dataset, true)
	if len(APTA.symbolMap) != 2{
		t.Errorf("APTA number of symbols = %d, want 2", len(APTA.symbolMap))
	}
	if len(APTA.states) != 322067{
		t.Errorf("APTA number of states = %d, want 322067", len(APTA.states))
	}
	if APTA.Depth() != 21{
		t.Errorf("APTA depth = %d, want 21", APTA.Depth())
	}
}

func TestAbbadingoDFAGeneration(t *testing.T) {
	// random seed
	rand.Seed(time.Now().UnixNano())
	numberOfStates := rand.Intn(499) + 1

	AbbadingoDFA := AbbadingoDFA(numberOfStates, true)
	if len(AbbadingoDFA.symbolMap) != 2{
		t.Errorf("AbbadingoDFA number of symbols = %d, want 2", len(AbbadingoDFA.symbolMap))
	}
	if len(AbbadingoDFA.states) != numberOfStates{
		t.Errorf("AbbadingoDFA number of states = %d, want %d", len(AbbadingoDFA.states), numberOfStates)
	}
	if AbbadingoDFA.Depth() != int(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0)){
		t.Errorf("AbbadingoDFA depth = %d, want %d", AbbadingoDFA.Depth(), int(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0)))
	}
}

func TestAbbadingoDatasetGeneration(t *testing.T){
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

	APTA := GetPTAFromDataset(trainingDataset, true)
	APTA.Describe(false)

	trainingDatasetConsistentWithAPTA := trainingDataset.ConsistentWithDFA(APTA)

	if !trainingDatasetConsistentWithAPTA{
		t.Errorf("Expected training dataset to be consistent with APTA")
	}
}