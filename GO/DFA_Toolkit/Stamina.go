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

// WriteToStaminaFile writes a given Dataset to file in Stamina-Format.
func (dataset Dataset) WriteToStaminaFile(filePath string) {
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