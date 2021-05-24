package dfatoolkit

// ExhaustiveEDSMFromDataset is a greedy version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func ExhaustiveEDSMFromDataset(dataset Dataset) (DFA, MergeData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call ExhaustiveEDSM function using APTA constructed
	// above. Return resultant DFA and search data.
	return ExhaustiveEDSM(APTA)
}

// WindowedEDSMFromDataset is a windowed version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func WindowedEDSMFromDataset(dataset Dataset, windowSize int, windowGrow float64) (DFA, MergeData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call FastWindowedEDSM function using APTA constructed
	// above. Return resultant DFA.
	return WindowedEDSM(APTA, windowSize, windowGrow)
}

// BlueFringeEDSMFromDataset is a Blue Fringe version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func BlueFringeEDSMFromDataset(dataset Dataset) (DFA, MergeData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call BlueFringeEDSM function using APTA constructed
	// above. Return resultant DFA and search data.
	return BlueFringeEDSM(APTA)
}

// ExhaustiveEDSM is a greedy version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the greedy search.
func ExhaustiveEDSM(APTA DFA) (DFA, MergeData) {
	// Store length of dataset.
	LengthOfDataset := APTA.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call ExhaustiveSearchUsingScoringFunction function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition and the search data.
	statePartition, mergeData := ExhaustiveSearchUsingScoringFunction(statePartition, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
	return resultantDFA, mergeData
}

// WindowedEDSM is a windowed version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the windowed search.
func WindowedEDSM(APTA DFA, windowSize int, windowGrow float64) (DFA, MergeData) {
	// Store length of dataset.
	LengthOfDataset := APTA.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call FastWindowedSearchUsingScoringFunction function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition and the search data.
	statePartition, mergeData := WindowedSearchUsingScoringFunction(statePartition, windowSize, windowGrow, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
	return resultantDFA, mergeData
}

// BlueFringeEDSM is a Blue Fringe version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the blue-fringe search.
func BlueFringeEDSM(APTA DFA) (DFA, MergeData) {
	// Store length of dataset.
	LengthOfDataset := APTA.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func(stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call BlueFringeSearchUsingScoringFunction function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition and the search data.
	statePartition, mergeData := BlueFringeSearchUsingScoringFunction(statePartition, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
	return resultantDFA, mergeData
}
