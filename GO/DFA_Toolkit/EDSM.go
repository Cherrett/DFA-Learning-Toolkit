package dfatoolkit

// ExhaustiveEDSMFromDataset is a greedy version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func ExhaustiveEDSMFromDataset(dataset Dataset) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call ExhaustiveEDSM function using APTA constructed
	// above. Return resultant DFA and search data.
	return ExhaustiveEDSM(APTA)
}

// FastWindowedEDSMFromDataset is a fast windowed version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func FastWindowedEDSMFromDataset(dataset Dataset, windowSize int, windowGrow float64) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call FastWindowedEDSM function using APTA constructed
	// above. Return resultant DFA and search data.
	return FastWindowedEDSM(APTA, windowSize, windowGrow)
}

// WindowedEDSMFromDataset is a windowed version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func WindowedEDSMFromDataset(dataset Dataset, windowSize int, windowGrow float64) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call FastWindowedEDSM function using APTA constructed
	// above. Return resultant DFA.
	return WindowedEDSM(APTA, windowSize, windowGrow)
}

// BlueFringeEDSMFromDataset is a Blue Fringe version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
func BlueFringeEDSMFromDataset(dataset Dataset) (DFA, SearchData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call BlueFringeEDSM function using APTA constructed
	// above. Return resultant DFA and search data.
	return BlueFringeEDSM(APTA)
}

// ExhaustiveEDSM is a greedy version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the greedy search.
func ExhaustiveEDSM(APTA DFA) (DFA, SearchData) {
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

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
	return resultantDFA, searchData
}

// FastWindowedEDSM is a fast windowed version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the windowed search.
func FastWindowedEDSM(APTA DFA, windowSize int, windowGrow float64) (DFA, SearchData) {
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
	statePartition, searchData := FastWindowedSearchUsingScoringFunction(statePartition, windowSize, windowGrow, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
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

	// Call FastWindowedSearchUsingScoringFunction function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition and the search data.
	statePartition, searchData := WindowedSearchUsingScoringFunction(statePartition, windowSize, windowGrow, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
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

	// Call BlueFringeSearchUsingScoringFunction function using state partition and EDSM scoring function
	// declared above. This function returns the resultant state partition and the search data.
	statePartition, searchData := BlueFringeSearchUsingScoringFunction(statePartition, EDSM)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
	return resultantDFA, searchData
}
