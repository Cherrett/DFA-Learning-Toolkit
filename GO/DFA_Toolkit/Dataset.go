package dfatoolkit

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
	Value  []rune      // Slice of runes which represents the actual string value.
	Accepting bool 	   // Accepting label flag for StringInstance.
}

// Dataset which is a slice of string instances.
type Dataset []StringInstance

// GetPTA returns an APTA/PTA (DFA) given a dataset.
// If APTA is set to true, an APTA is returned.
func (dataset Dataset) GetPTA(APTA bool) DFA {
	// Sort dataset by length.
	sortedDataset := dataset.SortDatasetByLength()
	// Create a map to keep track of alphabet (symbols).
	alphabet := make(map[rune]bool)
	// Count to keep track of the number of states visited.
	var count int
	// Variables to store IDs of current and new states.
	var currentStateID, newStateID int
	// Initialize new DFA.
	dfa := NewDFA()

	// If the first string instance within dataset has a length of 0,
	// it represents the starting state so add a state to the DFA
	// using the label of the empty string instance.
	if sortedDataset[0].Length() == 0{
		if sortedDataset[0].Accepting {
			currentStateID = dfa.AddState(ACCEPTING)
		} else if APTA {
			currentStateID = dfa.AddState(REJECTING)
		} else {
			currentStateID = dfa.AddState(UNKNOWN)
		}
	} else {
		currentStateID = dfa.AddState(UNKNOWN)
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

			// If symbol is not within symbol map,
			// add symbol to DFA and to symbol map.
			if !alphabet[symbol] {
				dfa.AddSymbol(symbol)
				alphabet[symbol] = true
			}

			// Set symbol ID to current symbol.
			symbolID := dfa.SymbolID(symbol)

			// If a transition exists from the current state to any other state via
			// the current symbol, set resultant state ID to current state ID.
			if dfa.States[currentStateID].Transitions[symbolID] != -1 {
				currentStateID = dfa.States[currentStateID].Transitions[symbolID]
				// Check if last symbol in string.
				if count == stringInstance.Length() {
					if stringInstance.Accepting {
						// Panic if string instance is accepting and resultant state is rejecting.
						if dfa.States[currentStateID].Label == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						// If string instance is accepting and resultant state is not
						// rejecting, set state label to accepting.
						} else {
							dfa.States[currentStateID].Label = ACCEPTING
						}
					} else {
						// Panic if string instance is rejecting and resultant state is accepting.
						if dfa.States[currentStateID].Label == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						// If string instance is rejecting and resultant state is not
						// accepting, set state label to rejecting.
						} else {
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
					newStateID = dfa.AddState(UNKNOWN)
				}
				// Add a new transition from current state to newly created state.
				dfa.AddTransition(symbolID, currentStateID, newStateID)
				// Set current state ID to ID of newly created state.
				currentStateID = newStateID
			}
		}
	}

	// Return populated DFA.
	return dfa
}

// Length returns the length of the string.
func (stringInstance StringInstance) Length() int{
	return len(stringInstance.Value)
}

// ConsistentWithDFA returns whether a string instance is consistent
// within a given DFA.
func (stringInstance StringInstance) ConsistentWithDFA(dfa DFA) bool {
	// Set the current state ID to the starting state ID.
	currentStateID := dfa.StartingStateID
	// Set counter to 0.
	count := 0

	// Iterate over each symbol (alphabet) within value of string instance.
	for _, symbol := range stringInstance.Value {
		// Increment count.
		count++

		// If a transition exists from the current state to any other state via
		// the current symbol, set resultant state ID to current state ID.
		if dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)] != -1 {
			currentStateID = dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)]
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
		// If no transition exists and string instance is accepting, return false.
		// If string instance is rejecting, return true.
		} else {
			return !stringInstance.Accepting
		}
	}

	// Return true if reached.
	return true
}

// ParseToStateLabel returns the State Label of a given string instance
// within a given DFA. If state does not exist, REJECTING is returned.
func (stringInstance StringInstance) ParseToStateLabel(dfa DFA) StateLabel {
	// Set the current state ID to the starting state ID.
	currentStateID := dfa.StartingStateID
	// Set counter to 0.
	count := 0

	// Iterate over each symbol (alphabet) within value of string instance.
	for _, symbol := range stringInstance.Value {
		// Increment count.
		count++

		// If a transition exists from the current state to any other state via
		// the current symbol, set resultant state ID to current state ID.
		if dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)] != -1 {
			currentStateID = dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)]
			// Check if last symbol in string.
			if count == stringInstance.Length() {
				// If state is unknown or rejecting, return rejecting.
				if dfa.States[currentStateID].Label == UNKNOWN || dfa.States[currentStateID].Label == REJECTING {
					return REJECTING
				}
				// Else, if state is accepting, return accepting.
				return ACCEPTING
			}
		// Return rejecting if no transition exists.
		} else {
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

	// Iterate over each symbol (alphabet) within value of string instance.
	for _, symbol := range stringInstance.Value {
		// Increment count.
		count++

		// If a transition exists from the current state to any other state via
		// the current symbol, set resultant state ID to current state ID.
		if dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)] != -1 {
			currentStateID = dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)]
			// If last symbol in string, return the current true and the current state ID.
			if count == stringInstance.Length() {
				return true, currentStateID
			}
		// Return false if no transition exists.
		} else {
			return false, -1
		}
	}
	// Return false if reached.
	return false, -1
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
		stringInstance.Value = append(stringInstance.Value, dfa.Symbol(symbolID))
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

// ConsistentWithDFA checks whether dataset is consistent with a given DFA.
func (dataset Dataset) ConsistentWithDFA(dfa DFA) bool {
	// Set consistent flag to true.
	consistent := true
	// Create wait group
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
				consistent = stringInstance.ConsistentWithDFA(dfa)
			}
		}(stringInstance, dfa)
	}

	// Wait for all go routines within wait
	// group to finish executing.
	wg.Wait()

	// Return consistent flag.
	return consistent
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
		if stringInstance.Accepting{
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
