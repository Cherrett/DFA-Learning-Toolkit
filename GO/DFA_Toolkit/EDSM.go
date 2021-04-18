package dfatoolkit

// GreedyEDSMFromDataset is a greedy version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func GreedyEDSMFromDataset(dataset Dataset) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call GreedyEDSM function using APTA constructed
	// above. Return resultant DFA.
	return GreedyEDSM(APTA)
}

// WindowedEDSMFromDataset is a windowed version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func WindowedEDSMFromDataset(dataset Dataset, windowSize int, windowGrow float64) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call WindowedEDSM function using APTA constructed
	// above. Return resultant DFA.
	return WindowedEDSM(APTA, windowSize, windowGrow)
}

// WindowedEDSMFromDataset2 is a windowed version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func WindowedEDSMFromDataset2(dataset Dataset, windowSize int, windowGrow float64) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call WindowedEDSM function using APTA constructed
	// above. Return resultant DFA.
	return WindowedEDSM2(APTA, windowSize, windowGrow)
}

// BlueFringeEDSMFromDataset is a Blue Fringe version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func BlueFringeEDSMFromDataset(dataset Dataset) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call BlueFringeEDSM function using APTA constructed
	// above. Return resultant DFA.
	return BlueFringeEDSM(APTA)
}

// GreedyEDSM is a greedy version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the greedy search.
func GreedyEDSM(APTA DFA) (DFA, SearchData) {
	// Store length of dataset.
	LengthOfDataset := APTA.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call GreedySearch function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition.
	statePartition, searchData := GreedySearch(statePartition, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA.
	return resultantDFA, searchData
}

// WindowedEDSM is a windowed version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the windowed search.
func WindowedEDSM(APTA DFA, windowSize int, windowGrow float64) (DFA, SearchData) {
	// Store length of dataset.
	LengthOfDataset := APTA.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call WindowedSearch function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition.
	statePartition, searchData := WindowedSearch(statePartition, windowSize, windowGrow, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA.
	return resultantDFA, searchData
}

// WindowedEDSM2 is a windowed version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the windowed search.
func WindowedEDSM2(APTA DFA, windowSize int, windowGrow float64) (DFA, SearchData) {
	// Store length of dataset.
	LengthOfDataset := APTA.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call WindowedSearch function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition.
	statePartition, searchData := WindowedSearch2(statePartition, windowSize, windowGrow, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA.
	return resultantDFA, searchData
}

// BlueFringeEDSM is a Blue Fringe version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the blue-fringe search.
func BlueFringeEDSM(APTA DFA) (DFA, SearchData) {
	// Store length of dataset.
	LengthOfDataset := APTA.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call BlueFringeSearch function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition.
	statePartition, searchData := BlueFringeSearch(statePartition, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA.
	return resultantDFA, searchData
}
