package dfatoolkit

// GreedyEDSMFromDataset is a greedy version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
// The randomFromBest argument is a flag used within the GreedySearch function.
func GreedyEDSMFromDataset(dataset Dataset) DFA{
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call GreedyEDSM function using APTA constructed
	// above. Return resultant DFA.
	return GreedyEDSM(APTA)
}

// WindowedEDSMFromDataset is a windowed version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
// The randomFromBest argument is a flag used within the WindowedSearch function.
func WindowedEDSMFromDataset(dataset Dataset, windowSize int, windowGrow float64) DFA{
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call WindowedEDSM function using APTA constructed
	// above. Return resultant DFA.
	return WindowedEDSM(APTA, windowSize, windowGrow)
}

// BlueFringeEDSMFromDataset is a Blue Fringe version of Evidence Driven State-Merging.
// It takes a dataset as an argument which is used to generate an APTA.
// The randomFromBest argument is a flag used within the BlueFringeSearch function.
func BlueFringeEDSMFromDataset(dataset Dataset) DFA{
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call BlueFringeEDSM function using APTA constructed
	// above. Return resultant DFA.
	return BlueFringeEDSM(APTA)
}

// GreedyEDSM is a greedy version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the greedy search.
// The randomFromBest argument is a flag used within the GreedySearch function.
func GreedyEDSM(apta DFA) DFA{
	// Store length of dataset.
	LengthOfDataset := apta.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func (stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition, dfa DFA) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Call GreedySearch function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return GreedySearch(apta, EDSM)
}

// WindowedEDSM is a windowed version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the windowed search.
// The randomFromBest argument is a flag used within the WindowedSearch function.
func WindowedEDSM(apta DFA, windowSize int, windowGrow float64) DFA{
	// Store length of dataset.
	LengthOfDataset := apta.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func (stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition, dfa DFA) float64 {
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Call WindowedSearch function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return WindowedSearch(apta, windowSize, windowGrow, EDSM)
}

// BlueFringeEDSM is a Blue Fringe version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the blue-fringe search.
// The randomFromBest argument is a flag used within the BlueFringeSearch function.
func BlueFringeEDSM(apta DFA) DFA{
	// Store length of dataset.
	LengthOfDataset := apta.LabelledStatesCount()

	// EDSM scoring function.
	EDSM := func (stateID1, stateID2 int, partitionBefore, partitionAfter StatePartition, dfa DFA) float64 {
		//return (100 * (partitionBefore.NumberOfLabelledBlocks() - partitionAfter.NumberOfLabelledBlocks())) + 99 - dfa.States[stateID2].Depth()
		return float64(LengthOfDataset - partitionAfter.NumberOfLabelledBlocks())
	}

	// Call WindowedSearch function using APTA and EDSM scoring function
	// declared above. Return resultant DFA.
	return BlueFringeSearch(apta, EDSM)
}
