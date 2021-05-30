package dfalearningtoolkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

// StringInstance struct which represents
// a string instance within a dataset.
type StringInstance struct {
	Value     []int // Slice of integers which represents the symbol values.
	Accepting bool  // Accepting label flag for StringInstance.
}

// Dataset which is a slice of string instances.
type Dataset []StringInstance

// GetPTA returns an APTA/PTA (DFA) given a dataset.
// If APTA is set to true, an APTA is returned.
func (dataset Dataset) GetPTA(APTA bool) DFA {
	// Sort dataset by length.
	sortedDataset := dataset.SortDatasetByLength()
	// Count to keep track of the number of states visited.
	var count int
	// Variables to store IDs of current and new states.
	var currentStateID, newStateID int
	// Initialize new DFA.
	dfa := NewDFA()

	// If the first string instance within dataset has a length of 0,
	// it represents the starting state so add a state to the DFA
	// using the label of the empty string instance.
	if sortedDataset[0].Length() == 0 {
		if sortedDataset[0].Accepting {
			currentStateID = dfa.AddState(ACCEPTING)
		} else if APTA {
			currentStateID = dfa.AddState(REJECTING)
		} else {
			currentStateID = dfa.AddState(UNLABELLED)
		}
	} else {
		currentStateID = dfa.AddState(UNLABELLED)
	}

	// Set starting state ID within DFA to current state ID.
	dfa.StartingStateID = currentStateID

	// Iterate over each string instance within sorted dataset.
	for _, stringInstance := range sortedDataset {
		// If APTA is set to false and string instance is
		// rejecting, skip string instance.
		if !APTA && !stringInstance.Accepting {
			continue
		}
		// Set current state ID to starting state ID.
		currentStateID = dfa.StartingStateID
		// Set count to 0.
		count = 0
		// Iterate over each symbol (alphabet) within
		// value of string instance.
		for _, symbol := range stringInstance.Value {
			// Increment count.
			count++

			// While symbol is not within DFA's symbols,
			// add new symbol to DFA.
			for symbol > len(dfa.Alphabet)-1 {
				dfa.AddSymbol()
			}

			// If a transition exists from the current state to any other state via
			// the current symbol, set resultant state ID to current state ID.
			if dfa.States[currentStateID].Transitions[symbol] >= 0 {
				currentStateID = dfa.States[currentStateID].Transitions[symbol]
				// Check if last symbol in string.
				if count == stringInstance.Length() {
					if stringInstance.Accepting {
						// Panic if string instance is accepting and resultant state is rejecting.
						if dfa.States[currentStateID].Label == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						} else {
							// If string instance is accepting and resultant state is not
							// rejecting, set state label to accepting.
							dfa.States[currentStateID].Label = ACCEPTING
						}
					} else {
						// Panic if string instance is rejecting and resultant state is accepting.
						if dfa.States[currentStateID].Label == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							// If string instance is rejecting and resultant state is not
							// accepting, set state label to rejecting.
							dfa.States[currentStateID].Label = REJECTING
						}
					}
				}
				// If no transition exists, add new state.
			} else {
				// Check if last symbol in string.
				if count == stringInstance.Length() {
					// If string instance is accepting, add an accepting state.
					// Otherwise, add a rejecting state.
					if stringInstance.Accepting {
						newStateID = dfa.AddState(ACCEPTING)
					} else {
						newStateID = dfa.AddState(REJECTING)
					}
					// If not last symbol in string, add an unknown state.
				} else {
					newStateID = dfa.AddState(UNLABELLED)
				}
				// Add a new transition from current state to newly created state.
				dfa.AddTransition(symbol, currentStateID, newStateID)
				// Set current state ID to ID of newly created state.
				currentStateID = newStateID
			}
		}
	}

	// Return populated DFA.
	return dfa
}

// Length returns the length of the string.
func (stringInstance StringInstance) Length() int {
	return len(stringInstance.Value)
}

// ConsistentWithDFA returns whether a string instance is consistent
// within a given DFA.
func (stringInstance StringInstance) ConsistentWithDFA(dfa DFA) bool {
	// Set the current state ID to the starting state ID.
	currentStateID := dfa.StartingStateID
	// Set counter to 0.
	count := 0

	// If string instance is the empty string, compare label
	// with starting state within DFA.
	if len(stringInstance.Value) == 0 {
		if stringInstance.Accepting {
			// If string instance is accepting and starting state is rejecting, return false.
			if dfa.StartingState().Label == REJECTING {
				return false
			}
		} else {
			// If string instance is rejecting and starting state is accepting, return false.
			if dfa.StartingState().Label == ACCEPTING {
				return false
			}
		}

		// Return true since labels match.
		return true
	}

	// Iterate over each symbol (alphabet) within value of string instance.
	for _, symbol := range stringInstance.Value {
		// Increment count.
		count++

		// If a transition exists from the current state to any other state via
		// the current symbol, set resultant state ID to current state ID.
		if dfa.States[currentStateID].Transitions[symbol] >= 0 {
			currentStateID = dfa.States[currentStateID].Transitions[symbol]
			// Check if last symbol in string.
			if count == stringInstance.Length() {
				if stringInstance.Accepting {
					// If string instance is accepting and state is rejecting, return false.
					if dfa.States[currentStateID].Label == REJECTING {
						return false
					}
				} else {
					// If string instance is rejecting and state is accepting, return false.
					if dfa.States[currentStateID].Label == ACCEPTING {
						return false
					}
				}
			}
		} else {
			// If no transition exists and string instance is accepting, return false.
			// If string instance is rejecting, return true.
			return !stringInstance.Accepting
		}
	}

	// Return true if reached.
	return true
}

// ConsistentWithDFA checks whether dataset is consistent with a given DFA.
func (dataset Dataset) ConsistentWithDFA(dfa DFA) bool {
	// Set consistent flag to true.
	consistent := true
	// Create wait group.
	var wg sync.WaitGroup
	// Add length of dataset to wait group.
	wg.Add(dataset.Count())

	// Iterate over each string instance within dataset.
	for _, stringInstance := range dataset {
		go func(stringInstance StringInstance, dfa DFA) {
			// Decrement 1 from wait group.
			defer wg.Done()
			// If consistent flag is true, check if current
			// string instance is consistent with DFA. If
			// consistent flag is already set to false, skip
			// checking the remaining string instances.
			if consistent {
				if !stringInstance.ConsistentWithDFA(dfa) {
					consistent = false
				}
			}
		}(stringInstance, dfa)
	}

	// Wait for all go routines within wait
	// group to finish executing.
	wg.Wait()

	// Return consistent flag.
	return consistent
}

// ConsistentWithStatePartition returns whether a string instance is consistent
// within a given StatePartition.
func (stringInstance StringInstance) ConsistentWithStatePartition(statePartition StatePartition) bool {
	// Set the current block ID to the starting block ID.
	currentBlockID := statePartition.StartingBlock()
	// Set counter to 0.
	count := 0

	// If string instance is the empty string, compare label
	// with starting state within DFA.
	if len(stringInstance.Value) == 0 {
		if stringInstance.Accepting {
			// If string instance is accepting and starting state is rejecting, return false.
			if statePartition.Blocks[currentBlockID].Label == REJECTING {
				return false
			}
		} else {
			// If string instance is rejecting and starting state is accepting, return false.
			if statePartition.Blocks[currentBlockID].Label == ACCEPTING {
				return false
			}
		}

		// Return true since labels match.
		return true
	}

	// Iterate over each symbol (alphabet) within value of string instance.
	for _, symbol := range stringInstance.Value {
		// Increment count.
		count++

		// If a transition exists from the current state to any other state via
		// the current symbol, set resultant state ID to current state ID.
		if statePartition.Blocks[currentBlockID].Transitions[symbol] >= 0 {
			currentStateID := statePartition.Blocks[currentBlockID].Transitions[symbol]
			currentBlockID = statePartition.Find(currentStateID)
			// Check if last symbol in string.
			if count == stringInstance.Length() {
				if stringInstance.Accepting {
					// If string instance is accepting and state is rejecting, return false.
					if statePartition.Blocks[currentBlockID].Label == REJECTING {
						return false
					}
				} else {
					// If string instance is rejecting and state is accepting, return false.
					if statePartition.Blocks[currentBlockID].Label == ACCEPTING {
						return false
					}
				}
			}
		} else {
			// If no transition exists and string instance is accepting, return false.
			// If string instance is rejecting, return true.
			return !stringInstance.Accepting
		}
	}

	// Return true if reached.
	return true
}

// ConsistentWithStatePartition checks whether dataset is consistent with a given StatePartition.
func (dataset Dataset) ConsistentWithStatePartition(statePartition StatePartition) bool {
	// Set consistent flag to true.
	consistent := true
	// Create wait group.
	var wg sync.WaitGroup
	// Add length of dataset to wait group.
	wg.Add(dataset.Count())

	// Iterate over each string instance within dataset.
	for _, stringInstance := range dataset {
		go func(stringInstance StringInstance, statePartition StatePartition) {
			// Decrement 1 from wait group.
			defer wg.Done()
			// If consistent flag is true, check if current
			// string instance is consistent with DFA. If
			// consistent flag is already set to false, skip
			// checking the remaining string instances.
			if consistent {
				if !stringInstance.ConsistentWithStatePartition(statePartition) {
					consistent = false
				}
			}
		}(stringInstance, statePartition)
	}

	// Wait for all go routines within wait
	// group to finish executing.
	wg.Wait()

	// Return consistent flag.
	return consistent
}

// ParseToStateLabel returns the State Label of a given string instance
// within a given DFA. If state does not exist, REJECTING is returned.
func (stringInstance StringInstance) ParseToStateLabel(dfa DFA) StateLabel {
	// Set the current state ID to the starting state ID.
	currentStateID := dfa.StartingStateID
	// Set counter to 0.
	count := 0

	// If string instance is the empty string and the starting
	// state is accepting, return accepting. Else return rejecting.
	if len(stringInstance.Value) == 0 {
		if dfa.StartingState().Label == ACCEPTING {
			return ACCEPTING
		}
	}

	// Iterate over each symbol (alphabet) within value of string instance.
	for _, symbol := range stringInstance.Value {
		// Increment count.
		count++

		// If a transition exists from the current state to any other state via
		// the current symbol, set resultant state ID to current state ID.
		if dfa.States[currentStateID].Transitions[symbol] >= 0 {
			currentStateID = dfa.States[currentStateID].Transitions[symbol]
			// Check if last symbol in string.
			if count == stringInstance.Length() {
				// If state is unknown or rejecting, return rejecting.
				if dfa.States[currentStateID].Label == UNLABELLED || dfa.States[currentStateID].Label == REJECTING {
					return REJECTING
				}
				// Else, if state is accepting, return accepting.
				return ACCEPTING
			}
		} else {
			// Return rejecting if no transition exists.
			return REJECTING
		}
	}

	// Return rejecting if reached.
	return REJECTING
}

// ParseToState returns the State ID of a given string instance within
// a given DFA. A boolean is also returned which returns false if string
// instance does not correspond to any state within DFA.
func (stringInstance StringInstance) ParseToState(dfa DFA) (bool, int) {
	// Set the current state ID to the starting state ID.
	currentStateID := dfa.StartingStateID
	// Set counter to 0.
	count := 0

	// If string instance is the empty string, return
	// starting state ID within DFA.
	if len(stringInstance.Value) == 0 {
		return true, dfa.StartingStateID
	}

	// Iterate over each symbol (alphabet) within value of string instance.
	for _, symbol := range stringInstance.Value {
		// Increment count.
		count++

		// If a transition exists from the current state to any other state via
		// the current symbol, set resultant state ID to current state ID.
		if dfa.States[currentStateID].Transitions[symbol] >= 0 {
			currentStateID = dfa.States[currentStateID].Transitions[symbol]
			// If last symbol in string, return the current true and the current state ID.
			if count == stringInstance.Length() {
				return true, currentStateID
			}
		} else {
			// Return false if no transition exists.
			return false, -1
		}
	}

	// Return false if reached.
	return false, -1
}

// WithinDataset checks whether a StringInstance is within given Dataset.
func (stringInstance StringInstance) WithinDataset(dataset Dataset) bool {
	// Sort dataset by length.
	dataset = dataset.SortDatasetByLength()

	// Iterate over sorted dataset.
	for _, instance := range dataset {
		// Skip if length of string instance is smaller
		// than that of the stringInstance being searched.
		if len(instance.Value) < len(stringInstance.Value) {
			continue
		} else if len(instance.Value) > len(stringInstance.Value) {
			// Break if length of string instance is bigger
			// than that of the stringInstance being searched since
			// all strings which have the same size as the target
			// string have already been checked.
			break
		}

		// Else (same length as string being searched),
		// compare using deep equal. Return true if
		// they are equal.
		if reflect.DeepEqual(instance.Value, stringInstance.Value) {
			return true
		}
	}

	// Return false if reached since
	// loop was broken.
	return false
}

// BinaryStringToStringInstance returns a StringInstances from a binary string.
func BinaryStringToStringInstance(dfa DFA, binaryString string) StringInstance {
	// Initialize string instance.
	stringInstance := StringInstance{}

	// Iterate over binary string.
	for _, value := range binaryString {
		// Convert character to integer (symbol ID).
		symbolID, err := strconv.Atoi(string(value))
		// Panic if character is not an integer or integer
		// is not 0 or 1 (not binary).
		if err != nil || (symbolID != 0 && symbolID != 1) {
			panic("Not a binary string")
		}

		// Add symbol to value of string instance.
		stringInstance.Value = append(stringInstance.Value, symbolID)
	}

	// Set string instance accepting to true if string instance is accepting.
	// It is set to false if string instance is rejecting.
	stringInstance.Accepting = stringInstance.ParseToStateLabel(dfa) == ACCEPTING

	// Return populated string instance.
	return stringInstance
}

// SortDatasetByLength sorts the given dataset by length and returns it.
func (dataset Dataset) SortDatasetByLength() Dataset {
	// Sort all string instances by length in ascending order.
	sort.Slice(dataset[:], func(i, j int) bool {
		return len(dataset[i].Value) < len(dataset[j].Value)
	})

	// Return sorted dataset.
	return dataset
}

// AcceptingStringInstances returns the accepting string instances.
func (dataset Dataset) AcceptingStringInstances() Dataset {
	// Initialize empty dataset.
	var acceptingInstances Dataset

	// Populate dataset initialized with accepting string instances.
	for _, stringInstance := range dataset {
		if stringInstance.Accepting {
			acceptingInstances = append(acceptingInstances, stringInstance)
		}
	}

	// Return populated dataset.
	return acceptingInstances
}

// RejectingStringInstances returns the rejecting string instances.
func (dataset Dataset) RejectingStringInstances() Dataset {
	// Initialize empty dataset.
	var rejectingInstances Dataset

	// Populate dataset initialized with rejecting string instances.
	for _, stringInstance := range dataset {
		if !stringInstance.Accepting {
			rejectingInstances = append(rejectingInstances, stringInstance)
		}
	}

	// Return populated dataset.
	return rejectingInstances
}

// Count returns the number of string instances within dataset.
func (dataset Dataset) Count() int {
	return len(dataset)
}

// AcceptingStringInstancesCount returns the number of accepting
// string instances within dataset.
func (dataset Dataset) AcceptingStringInstancesCount() int {
	count := 0

	for _, stringInstance := range dataset {
		if stringInstance.Accepting {
			count++
		}
	}

	return count
}

// RejectingStringInstancesCount returns the number of accepting
// string instances within dataset.
func (dataset Dataset) RejectingStringInstancesCount() int {
	count := 0

	for _, stringInstance := range dataset {
		if !stringInstance.Accepting {
			count++
		}
	}

	return count
}

// AcceptingStringInstancesRatio returns the ratio of accepting string
// instances to the rejecting string instances within dataset.
func (dataset Dataset) AcceptingStringInstancesRatio() float64 {
	return float64(dataset.AcceptingStringInstancesCount()) / float64(len(dataset))
}

// RejectingStringInstancesRatio returns the ratio of rejecting string
// instances to the accepting string instances within dataset.
func (dataset Dataset) RejectingStringInstancesRatio() float64 {
	return float64(dataset.RejectingStringInstancesCount()) / float64(len(dataset))
}

// SameAs checks whether Dataset is the same as the given Dataset.
// Both datasets are sorted before checking with DeepEqual.
func (dataset Dataset) SameAs(dataset2 Dataset) bool {
	dataset1 := dataset.SortDatasetByLength()
	dataset2 = dataset2.SortDatasetByLength()
	return reflect.DeepEqual(dataset1, dataset2)
}

// StructurallyComplete checks if Dataset is structurally
// complete with respect to a DFA.
func (dataset Dataset) StructurallyComplete(dfa DFA) bool {
	return dfa.StructurallyComplete(dataset)
}

// SymmetricallyStructurallyComplete checks if Dataset is symmetrically
// structurally complete with respect to a DFA.
func (dataset Dataset) SymmetricallyStructurallyComplete(dfa DFA) bool {
	return dfa.SymmetricallyStructurallyComplete(dataset)
}

// Accuracy returns the DFA's Accuracy with respect to the dataset.
func (dataset Dataset) Accuracy(dfa DFA) float64 {
	return dfa.Accuracy(dataset)
}

// ToJSON saves the dataset to a JSON file given a file path.
func (dataset Dataset) ToJSON(filePath string) bool {
	// Sort the dataset by length.
	dataset.SortDatasetByLength()

	// Create file given a path/name.
	file, err := os.Create(filePath)

	// If file was not created successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Close file at end of function.
	defer file.Close()

	// Convert dataset to JSON.
	resultantJSON, err := json.MarshalIndent(dataset, "", "\t")

	// If dataset was not converted successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Copy JSON to file created.
	_, err = io.Copy(file, bytes.NewReader(resultantJSON))

	// If JSON was not copied successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Return true if reached.
	return true
}

// DatasetFromJSON returns a dataset read from a JSON file
// given a file path. The boolean value returned is set to
// true if Dataset was read successfully.
func DatasetFromJSON(filePath string) (Dataset, bool) {
	// Open file from given a path/name.
	file, err := os.Open(filePath)

	// If file was not opened successfully,
	// return empty dataset and false.
	if err != nil {
		return Dataset{}, false
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize empty Dataset.
	resultantDataset := Dataset{}

	// Convert JSON to Dataset.
	err = json.NewDecoder(file).Decode(&resultantDataset)

	// If JSON was not converted successfully,
	// return empty dataset and false.
	if err != nil {
		return Dataset{}, false
	}

	// Return populated Dataset and true if reached.
	return resultantDataset, true
}
