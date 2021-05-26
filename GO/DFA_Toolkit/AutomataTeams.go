package dfatoolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"math/rand"
	"time"
)

// This file contains various functions to conduct automata teams as described in:
// P. García, M. Vázquez de Parga, D. López, and J. Ruiz, ‘Learning Automata Teams’,
// in Grammatical Inference: Theoretical Results and Applications, Berlin, Heidelberg,
// 2010, pp. 52–65. doi: 10.1007/978-3-642-15488-1_6.

// TeamOfAutomata represents a team of automata within the AutomataTeams algorithm.
type TeamOfAutomata struct {
	Team      []DFA
	MergeData MergeData
}

// TeamOfAutomataClassifierFunction takes a string instance as input and returns a state label.
// An example of how this function type should be used can be seen in BetterHalfWeightedVoteAccuracy.
type TeamOfAutomataClassifierFunction func(stringInstance StringInstance) StateLabel

// GRBM deterministically merges possible state pairs within red-blue sets.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and search data when no more valid merges are possible.
func GRBM(statePartition StatePartition) (StatePartition, MergeData) {
	// Clone StatePartition.
	statePartition = statePartition.Clone()
	// Copy the state partition for undoing and copying changed states.
	copiedPartition := statePartition.Copy()
	// Initialize search data.
	mergeData := MergeData{[]StatePairScore{}, 0, 0, time.Duration(0)}
	// Total merges and valid merges counter.
	totalMerges, totalValidMerges := 0, 0
	// Start timer.
	start := time.Now()

	// Slice to store red states in insertion order.
	redStates := []int{statePartition.StartingBlock()}

	// Generated slice of ordered blue states from red states.
	blueStates := GenerateBlueSetFromRedSetWithShuffle(&statePartition, map[int]util.Void{redStates[0]: util.Null})
	// Shuffle blue sets.
	rand.Shuffle(len(blueStates), func(i, j int) { blueStates[i], blueStates[j] = blueStates[j], blueStates[i] })

	// Initialize merged flag to false.
	merged := false

	// Iterate until blue set is empty.
	for len(blueStates) != 0 {
		// Iterate over every blue state in insertion order.
		for _, blueElement := range blueStates {
			// Set merged flag to false.
			merged = false
			// Iterate over every red state in insertion order.
			for _, redElement := range redStates {
				// Increment merge count.
				totalMerges++
				// If states are mergeable, calculate score and add to detMerges.
				if copiedPartition.MergeStates(blueElement, redElement) {
					// Increment valid merge count.
					totalValidMerges++

					// Set merged flag to true.
					merged = true

					// Copy changes to original state partition.
					statePartition.CopyChangesFrom(&copiedPartition)

					// Add merged state pair with score to search data.
					mergeData.Merges = append(mergeData.Merges, StatePairScore{blueElement, redElement, -1})

					break
				}

				// Undo merge.
				copiedPartition.RollbackChangesFrom(statePartition)
			}

			// If merged flag is false, add current blue state
			// to red states set and ordered set and exit loop.
			if !merged {
				redStates = append(redStates, blueElement)
			}

			// Update red and blue states using UpdateRedBlueSetsWithShuffle function.
			blueStates = UpdateRedBlueSetsWithShuffle(&statePartition, &redStates)
		}
	}

	// Add total and valid merges counts to search data.
	mergeData.AttemptedMergesCount = totalMerges
	mergeData.ValidMergesCount = totalValidMerges
	// Add duration to search data.
	mergeData.Duration = time.Now().Sub(start)

	// Return the final resultant state partition and search data.
	return statePartition, mergeData
}

// GenerateBlueSetFromRedSetWithShuffle generates the blue set given the state partition and the red set within the Red-Blue framework
// such as the GeneralizedRedBlueMerging function. It generates and returns the blue set in arbitrary order.
func GenerateBlueSetFromRedSetWithShuffle(statePartition *StatePartition, redSet map[int]util.Void) []int {
	// Step 1 - Gather all blue states and store in map declared below.

	// Initialize set of blue states to empty set.
	var blueStates []int

	// Iterate over every red state.
	for element := range redSet {
		// Iterate over each symbol within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			if resultantStateID := statePartition.Blocks[element].Transitions[symbol]; resultantStateID > -1 {
				// Store resultant stateID from red state.
				resultantStateID = statePartition.Find(resultantStateID)
				// If transition is valid and resultant state is not red,
				// add resultant state to blue set.
				if _, exists := redSet[resultantStateID]; !exists {
					blueStates = append(blueStates, resultantStateID)
				}
			}
		}
	}

	rand.Shuffle(len(blueStates), func(i, j int) { blueStates[i], blueStates[j] = blueStates[j], blueStates[i] })

	// Return populated slice of blue states in canonical order.
	return blueStates
}

// UpdateRedBlueSetsWithShuffle updates the red and blue sets given the state partition and the red set within the Red-Blue framework
// such as the GeneralizedRedBlueMerging function. It returns the blue set and modifies the red set via its pointer. Both sets are
// shuffled . This is used when the state partition is changed or when new states have been added to the red set.
func UpdateRedBlueSetsWithShuffle(statePartition *StatePartition, redStates *[]int) []int {
	// Step 1 - Gather root of old red states and store in map declared below.

	// Initialize set of red root states (blocks) to empty set.
	redSet := make(map[int]util.Void, len(*redStates))

	// Iterate over every red state.
	for i := range *redStates {
		redElement := &(*redStates)[i]
		*redElement = statePartition.Find(*redElement)

		redSet[*redElement] = util.Null
	}

	// Step 2 - Gather all blue states and store in map declared below.

	// Initialize set of blue states to empty set.
	var blueStates []int

	// Iterate over every red state.
	for _, element := range *redStates {
		// Iterate over each symbol within DFA.
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			if resultantStateID := statePartition.Blocks[element].Transitions[symbol]; resultantStateID > -1 {
				// Store resultant stateID from red state.
				resultantStateID = statePartition.Find(resultantStateID)
				// If transition is valid and resultant state is not red,
				// add resultant state to blue set.
				if _, exists := redSet[resultantStateID]; !exists {
					blueStates = append(blueStates, resultantStateID)
				}
			}
		}
	}

	// Shuffle both red and blue sets.
	rand.Shuffle(len(*redStates), func(i, j int) { (*redStates)[i], (*redStates)[j] = (*redStates)[j], (*redStates)[i] })
	rand.Shuffle(len(blueStates), func(i, j int) { blueStates[i], blueStates[j] = blueStates[j], blueStates[i] })

	// Return populated slice of red states and populated slice of blue states in canonical order.
	return blueStates
}

// GeneralizedRedBlueMergingFromDataset is a.
// It takes a dataset as an argument which is used to generate an APTA.
func GeneralizedRedBlueMergingFromDataset(dataset Dataset) (DFA, MergeData) {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call GeneralizedRedBlueMerging function using APTA constructed
	// above. Return resultant DFA and search data.
	return GeneralizedRedBlueMerging(APTA)
}

// GeneralizedRedBlueMerging is a.
// It takes a DFA (APTA) as an argument which is used within the greedy search.
func GeneralizedRedBlueMerging(APTA DFA) (DFA, MergeData) {
	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Call GRBM function using state partition declared above.
	// This function returns the resultant state partition and the search data.
	statePartition, mergeData := GRBM(statePartition)

	// Convert the state partition to a DFA.
	resultantDFA := statePartition.ToQuotientDFA()

	// Check if DFA generated is valid.
	resultantDFA.IsValidPanic()

	// Return resultant DFA and search data.
	return resultantDFA, mergeData
}

// AutomataTeamsFromDataset creates a team of automata of size teamSize. This creates a
// TeamOfAutomata instance which can then be used to classify testing examples.
// It takes a dataset as an argument which is used to generate an APTA.
func AutomataTeamsFromDataset(dataset Dataset, teamSize int) TeamOfAutomata {
	// Construct an APTA from dataset.
	APTA := dataset.GetPTA(true)

	// Call AutomataTeams function using APTA constructed
	// above. Return resultant team of automata.
	return AutomataTeams(APTA, teamSize)
}

// AutomataTeams creates a team of automata of size teamSize. This creates a
// TeamOfAutomata instance which can then be used to classify testing examples.
func AutomataTeams(APTA DFA, teamSize int) TeamOfAutomata {
	// Convert APTA to StatePartition for state merging.
	statePartition := APTA.ToStatePartition()

	// Initialize team.
	teamOfAutomata := TeamOfAutomata{
		Team:      make([]DFA, teamSize),
		MergeData: MergeData{},
	}

	for i := 0; i < teamSize; i++ {
		// Call GRBM function using state partition declared above.
		// This function returns the resultant state partition and the search data.
		resultantStatePartition, mergeData := GRBM(statePartition)

		// Convert the state partition to a DFA.
		resultantDFA := resultantStatePartition.ToQuotientDFA()

		// Check if DFA generated is valid.
		resultantDFA.IsValidPanic()

		// Add resultant DFA and merge data to team.
		teamOfAutomata.Team[i] = resultantDFA
		teamOfAutomata.MergeData.Duration += mergeData.Duration
		teamOfAutomata.MergeData.AttemptedMergesCount += mergeData.AttemptedMergesCount
		teamOfAutomata.MergeData.ValidMergesCount += mergeData.ValidMergesCount
	}

	// Return team of automata.
	return teamOfAutomata
}

// Accuracy returns the TeamOfAutomata's Accuracy with respect to a dataset.
func (teamOfAutomata TeamOfAutomata) Accuracy(dataset Dataset, teamOfAutomataClassifierFunction TeamOfAutomataClassifierFunction) float64 {
	// Correct classifications count.
	correctClassifications := float64(0)

	// Iterate over each string instance within dataset.
	for _, stringInstance := range dataset {
		// If the label of the string instance is equal to its state label
		// within the DFA, increment correct classifications count.
		if stringInstance.Accepting == (teamOfAutomataClassifierFunction(stringInstance) == ACCEPTING) {
			correctClassifications++
		}
	}

	// Return the number of correct classifications divided by the length of
	// the dataset.
	return correctClassifications / float64(len(dataset))
}

// FairVoteAccuracy returns the TeamOfAutomata's Accuracy with respect to a dataset
// using the fair vote scoring heuristic.
func (teamOfAutomata TeamOfAutomata) FairVoteAccuracy(dataset Dataset) float64 {
	fairVote := func(stringInstance StringInstance) StateLabel {
		accepting := 0
		unlabelledOrRejecting := 0

		for _, dfa := range teamOfAutomata.Team {
			if stringInstance.ParseToStateLabel(dfa) == ACCEPTING {
				accepting++
			} else {
				unlabelledOrRejecting++
			}
		}

		if accepting > unlabelledOrRejecting {
			return ACCEPTING
		} else if unlabelledOrRejecting > accepting {
			return REJECTING
		} else {
			if rand.Intn(2) == 0 {
				return ACCEPTING
			} else {
				return REJECTING
			}
		}
	}

	return teamOfAutomata.Accuracy(dataset, fairVote)
}

// WeightedVoteAccuracy returns the TeamOfAutomata's Accuracy with respect to a dataset
// using the weighted vote scoring heuristic.
func (teamOfAutomata TeamOfAutomata) WeightedVoteAccuracy(dataset Dataset) float64 {
	weightedVote := func(stringInstance StringInstance) StateLabel {
		accepting := 0.0
		unlabelledOrRejecting := 0.0

		for _, dfa := range teamOfAutomata.Team {
			dfaSize := float64(len(dfa.States))
			if stringInstance.ParseToStateLabel(dfa) == ACCEPTING {
				accepting += 1 / (dfaSize * dfaSize)
			} else {
				unlabelledOrRejecting += 1 / (dfaSize * dfaSize)
			}
		}

		if accepting > unlabelledOrRejecting {
			return ACCEPTING
		} else if unlabelledOrRejecting > accepting {
			return REJECTING
		} else {
			if rand.Intn(2) == 0 {
				return ACCEPTING
			} else {
				return REJECTING
			}
		}
	}

	return teamOfAutomata.Accuracy(dataset, weightedVote)
}

// BetterHalfWeightedVoteAccuracy returns the TeamOfAutomata's Accuracy with respect to a dataset
// using the better half weighted vote scoring heuristic.
func (teamOfAutomata TeamOfAutomata) BetterHalfWeightedVoteAccuracy(dataset Dataset) float64 {
	sizeCount := 0.0

	for _, dfa := range teamOfAutomata.Team {
		sizeCount += float64(len(dfa.States))
	}

	averageSize := sizeCount / float64(len(teamOfAutomata.Team))

	weightedVote := func(stringInstance StringInstance) StateLabel {
		accepting := 0.0
		unlabelledOrRejecting := 0.0

		for _, dfa := range teamOfAutomata.Team {
			dfaSize := float64(len(dfa.States))

			if dfaSize > averageSize {
				continue
			}

			if stringInstance.ParseToStateLabel(dfa) == ACCEPTING {
				accepting += 1 / (dfaSize * dfaSize)
			} else {
				unlabelledOrRejecting += 1 / (dfaSize * dfaSize)
			}
		}

		if accepting > unlabelledOrRejecting {
			return ACCEPTING
		} else if unlabelledOrRejecting > accepting {
			return REJECTING
		} else {
			if rand.Intn(2) == 0 {
				return ACCEPTING
			} else {
				return REJECTING
			}
		}
	}

	return teamOfAutomata.Accuracy(dataset, weightedVote)
}

// AverageNumberOfStates returns the average number of states of DFAs within team.
func (teamOfAutomata TeamOfAutomata) AverageNumberOfStates() int {
	averageNumberOfStates := util.StatsTracker{}

	for _, dfa := range teamOfAutomata.Team {
		averageNumberOfStates.AddInt(len(dfa.States))
	}

	return int(averageNumberOfStates.Mean())
}
