package dfatoolkit

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetDatasetFromStaminaFile returns a Dataset from a Stamina-Format File.
func GetDatasetFromStaminaFile(fileName string) Dataset {
	// Initialize new Dataset.
	dataset := Dataset{}

	// Open given file name.
	file, err := os.Open(fileName)

	// Panic if file does not exist.
	if err != nil {
		panic("Invalid file path/name")
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize a new scanner given the read file.
	scanner := bufio.NewScanner(file)
	// Iterate over each line.
	for scanner.Scan() {
		// Add StringInstance extracted from NewStringInstanceFromStaminaFile to dataset.
		dataset = append(dataset, NewStringInstanceFromStaminaFile(scanner.Text(), " "))
	}

	// Return read dataset.
	return dataset
}

// NewStringInstanceFromStaminaFile returns a StringInstance from a line within a Stamina-Format File.
func NewStringInstanceFromStaminaFile(text string, delimiter string) StringInstance {
	// Initialize new StringInstance.
	stringInstance := StringInstance{}

	// Split the string by delimiter inputted.
	splitString := strings.Split(text, delimiter)

	// Check whether string is accepting or rejecting.
	// Update label of StringInstance accordingly.
	// Panic if a value which is not +, or - is found.
	switch splitString[0] {
	case "+":
		stringInstance.Accepting = true
		break
	case "-":
		stringInstance.Accepting = false
		break
	default:
		panic(fmt.Sprintf("Unknown string label - %s", splitString[0]))
	}

	// Add the remaining split string values to value of StringInstance since
	// these contain the actual string value.
	for i := 1; i < len(splitString); i++{
		// Get integer value from split string.
		integerValue, err := strconv.Atoi(splitString[i])

		// Panic if error.
		if err != nil{
			panic(err)
		}

		// Add integer value to value slice within string instance.
		stringInstance.Value = append(stringInstance.Value, integerValue)
	}

	// Return populated string instance.
	return stringInstance
}

// ToStaminaFile writes a given Dataset to file in Stamina-Format.
func (dataset Dataset) ToStaminaFile(filePath string) {
	// Create file given a path/name.
	file, err := os.Create(filePath)

	// Panic if file was not created successfully.
	if err != nil {
		panic("Invalid file path/name")
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize a new writer given the created file.
	writer := bufio.NewWriter(file)

	// Write the length of the dataset and 2 (Since abbadingo DFAs
	// consist of  2 symbols (a and b).
	_, _ = writer.WriteString(strconv.Itoa(len(dataset)) + " 2\n")

	// Iterate over each string instance within sorted dataset.
	for _, stringInstance := range dataset {
		// Add string label and string length to output string.
		outputString := ""
		if stringInstance.Accepting {
			outputString = strconv.Itoa(1) + " " + strconv.Itoa(stringInstance.Length()) + " "
		} else {
			outputString = strconv.Itoa(0) + " " + strconv.Itoa(stringInstance.Length()) + " "
		}

		// Iterate over the value of string and add to output string.
		for _, symbol := range stringInstance.Value {
			if symbol == 'a' {
				outputString += "0 "
			} else {
				outputString += "1 "
			}
		}
		// Remove trailing space from output string and add new line char.
		outputString = strings.TrimSuffix(outputString, " ") + "\n"
		// Write output string to file.
		_, _ = writer.WriteString(outputString)
	}

	// Flush writer.
	_ = writer.Flush()
}

// GetDFAFromStaminaFile returns a DFA from a Stamina-Format File (adl).
func GetDFAFromStaminaFile(fileName string) DFA {
	// Initialize new DFA.
	dfa := NewDFA()

	// Open given file name.
	file, err := os.Open(fileName)

	// Panic if file does not exist.
	if err != nil {
		panic("Invalid file path/name")
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize a new scanner given the read file.
	scanner := bufio.NewScanner(file)
	// Read first line.
	scanner.Scan()

	// Get first value from first line by splitting it.
	numberOfStatesString := strings.Split(scanner.Text(), " ")[0]

	// Convert number of states (string) into an integer.
	numberOfStates, err := strconv.Atoi(numberOfStatesString)

	// Panic if error.
	if err != nil{
		panic(err)
	}

	// Iterate over each line which corresponds to a state.
	for stateID := 0; stateID < numberOfStates; stateID++ {
		// Scan next line.
		scanner.Scan()
		// Split read line.
		line := strings.Split(scanner.Text(), " ")

		// Check if state is starting state.
		initialState, err := strconv.ParseBool(line[1])

		// Panic if error.
		if err != nil{
			panic(err)
		}

		// If starting state is set to true, set starting state
		// within DFA to current state.
		if initialState{
			// Panic if starting state is already found.
			if dfa.StartingStateID > -1{
				panic("DFA cannot have more than 2 starting states.")
			}
			dfa.StartingStateID = stateID
		}

		// Check if state is accepting.
		acceptingState, err := strconv.ParseBool(line[1])

		// Panic if error.
		if err != nil{
			panic(err)
		}

		if acceptingState{
			dfa.AddState(ACCEPTING)
		}else{
			dfa.AddState(UNKNOWN)
		}
	}

	// Iterate over remaining lines (transitions).
	for scanner.Scan(){
		// Split read line.
		line := strings.Split(scanner.Text(), " ")

		// Get from state ID of transition.
		fromStateID, err := strconv.Atoi(line[0])

		// Panic if error.
		if err != nil{
			panic(err)
		}

		// Get to state ID of transition.
		toStateID, err := strconv.Atoi(line[1])

		// Panic if error.
		if err != nil{
			panic(err)
		}

		// Get symbol ID of transition.
		symbolID, err := strconv.Atoi(line[2])

		// Panic if error.
		if err != nil{
			panic(err)
		}

		// Add new symbol within alphabet until
		// symbolID is created.
		for symbolID + 1 > len(dfa.Alphabet){
			dfa.AddSymbol()
		}

		// Add transition to DFA.
		dfa.AddTransition(symbolID, fromStateID, toStateID)
	}

	// Return read dfa.
	return dfa
}

// ToStaminaFile writes a given DFA to file in Stamina-Format (adl).
func (dfa DFA) ToStaminaFile(filePath string) {
	// Create file given a path/name.
	file, err := os.Create(filePath)

	// Panic if file was not created successfully.
	if err != nil {
		panic("Invalid file path/name")
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize a new writer given the created file.
	writer := bufio.NewWriter(file)

	// Write the number of states and number of transitions.
	_, _ = writer.WriteString(strconv.Itoa(len(dfa.States)) + " " + strconv.Itoa(dfa.TransitionsCount()) + "\n")

	// Iterate over each state within dfa.
	for stateID, state := range dfa.States{
		// Write state information.
		_, _ = writer.WriteString(strconv.Itoa(stateID) + " " +
			strconv.FormatBool(dfa.StartingStateID == stateID) + " " +
			strconv.FormatBool(state.Label == ACCEPTING) + "\n")
	}

	// Iterate over each state within dfa.
	for stateID, state := range dfa.States{
		// Write transitions information.
		for symbolID := range dfa.Alphabet{
			// Write if transition exists.
			if resultantStateID := state.Transitions[symbolID]; resultantStateID > -1{
				_, _ = writer.WriteString(strconv.Itoa(stateID) + " " +
					strconv.Itoa(resultantStateID) + " " +
					strconv.Itoa(symbolID) + "\n")
			}
		}
	}

	// Flush writer.
	_ = writer.Flush()
}
