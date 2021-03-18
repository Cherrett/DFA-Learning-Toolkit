package DFA_Toolkit

// GreedyEDSM is a greedy version of Evidence Driven State-Merging.
func GreedyEDSM(dataset Dataset, randomFromBest bool) DFA{
	// Store length of dataset.
	LengthOfDataset := len(dataset)
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// EDSM scoring function.
	EDSM := func (partition StatePartition) int {
		return LengthOfDataset - partition.NumberOfLabelledBlocks()
	}

	// Call GreedyMerge function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return GreedyMerge(APTA, EDSM, randomFromBest)
}

// WindowedEDSM is a windowed version of Evidence Driven State-Merging.
func WindowedEDSM(dataset Dataset, windowSize int, windowGrow float64, randomFromBest bool) DFA{
	// Store length of dataset.
	LengthOfDataset := len(dataset)
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// EDSM scoring function.
	EDSM := func (partition StatePartition) int {
		return LengthOfDataset - partition.NumberOfLabelledBlocks()
	}

	// Call WindowedMerge function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return WindowedMerge(APTA, windowSize, windowGrow, EDSM, randomFromBest)
}
