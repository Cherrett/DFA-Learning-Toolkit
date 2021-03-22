package DFA_Toolkit

// GreedyEDSM is a greedy version of Evidence Driven State-Merging.
func GreedyEDSM(dataset Dataset, randomFromBest bool) DFA{
	// Store length of dataset.
	LengthOfDataset := len(dataset)
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// EDSM scoring function.
	EDSM := func (stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition, dfa DFA) int {
		return LengthOfDataset - partitionAfter.NumberOfLabelledBlocks()
	}

	// Call GreedySearch function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return GreedySearch(APTA, EDSM, randomFromBest)
}

// WindowedEDSM is a windowed version of Evidence Driven State-Merging.
func WindowedEDSM(dataset Dataset, windowSize int, windowGrow float64, randomFromBest bool) DFA{
	// Store length of dataset.
	LengthOfDataset := len(dataset)
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// EDSM scoring function.
	EDSM := func (stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition, dfa DFA) int {
		return LengthOfDataset - partitionAfter.NumberOfLabelledBlocks()
	}

	// Call WindowedSearch function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return WindowedSearch(APTA, windowSize, windowGrow, EDSM, randomFromBest)
}

// BlueFringeEDSM is a Blue Fringe version of Evidence Driven State-Merging.
func BlueFringeEDSM(dataset Dataset, randomFromBest bool) DFA{
	// Store length of dataset.
	LengthOfDataset := len(dataset)
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// EDSM scoring function.
	EDSM := func (stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition, dfa DFA) int {
		//return (100 * (partitionBefore.NumberOfLabelledBlocks() - partitionAfter.NumberOfLabelledBlocks())) + 99 - dfa.States[stateID2].Depth()
		return LengthOfDataset - partitionAfter.NumberOfLabelledBlocks()
	}

	// Call WindowedSearch function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return BlueFringeSearch(APTA, EDSM, randomFromBest)
}
