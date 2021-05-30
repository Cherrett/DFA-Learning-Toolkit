package dfalearningtoolkit

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// GetDatasetFromAbbadingoFile returns a Dataset from an Abbadingo-Format File.
func GetDatasetFromAbbadingoFile(fileName string) Dataset {
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
	// Ignore first line.
	scanner.Scan()
	// Iterate over each line.
	for scanner.Scan() {
		// Add StringInstance extracted from NewStringInstanceFromAbbadingoFile to dataset.
		dataset = append(dataset, NewStringInstanceFromAbbadingoFile(scanner.Text(), " "))
	}

	// Return read dataset.
	return dataset
}

// NewStringInstanceFromAbbadingoFile returns a StringInstance from a line within an Abbadingo-Format File.
func NewStringInstanceFromAbbadingoFile(text string, delimiter string) StringInstance {
	// Initialize new StringInstance.
	stringInstance := StringInstance{}

	// Split the string by delimiter inputted.
	splitString := strings.Split(text, delimiter)

	// Check whether string is rejecting or accepting.
	// Update label of StringInstance accordingly.
	// Panic if a value which is not 0 or 1 is found.
	switch splitString[0] {
	case "0":
		stringInstance.Accepting = false
		break
	case "1":
		stringInstance.Accepting = true
		break
	default:
		panic(fmt.Sprintf("Unknown string label - %s", splitString[0]))
	}

	// Add the remaining split string values to value of StringInstance since
	// these contain the actual string value.
	for i := 2; i < len(splitString); i++ {
		intValue, err := strconv.Atoi(splitString[i])

		if err != nil {
			panic("Not an integer.")
		}

		stringInstance.Value = append(stringInstance.Value, intValue)
	}

	// Return populated string instance.
	return stringInstance
}

// AbbadingoDFA returns a random DFA using the Abbadingo protocol given a number of states.
// If exact is set to true, the resultant DFA will have the required depth as per Abbadingo protocol.
func AbbadingoDFA(numberOfStates int, exact bool) DFA {
	// The size of the DFA to be created.
	dfaSize := int(math.Round((5.0 * float64(numberOfStates)) / 4.0))

	// The depth of the DFA to be created.
	dfaDepth := int(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0))

	// Iterate till a valid DFA is generated.
	for {
		// Initialize a new DFA.
		dfa := NewDFA()
		// Add symbols 'a' and 'b' since these are
		// used in Abbadingo DFAs.
		dfa.AddSymbol()
		dfa.AddSymbol()

		// Create new states and assign either
		// an accepting or unknown label.
		for i := 0; i < dfaSize; i++ {
			if rand.Intn(2) == 0 {
				dfa.AddState(ACCEPTING)
			} else {
				dfa.AddState(UNLABELLED)
			}
		}

		// Iterate over created states.
		for stateID := range dfa.States {
			// Randomly add transitions for both symbols.
			dfa.AddTransition(0, stateID, rand.Intn(len(dfa.States)))
			dfa.AddTransition(1, stateID, rand.Intn(len(dfa.States)))
		}

		// Randomly choose a starting state.
		dfa.StartingStateID = rand.Intn(len(dfa.States))

		// Minimise DFA created.
		dfa = dfa.Minimise()

		// Get depth of DFA created.
		currentDFADepth := dfa.Depth()

		// If depth of created DFA is not equal to the required
		// depth, restart the process of creating a random DFA.
		if currentDFADepth == dfaDepth {
			// If exact is set to true, the number of states
			// within created DFA must be equal to the amount
			// of required states.
			if exact {
				if len(dfa.States) == numberOfStates {
					// Return the created DFA since it
					// meets all of the requirements.
					return dfa
				}
			} else {
				// If exact is set to false, return the created DFA
				// since it meets the requirements.
				return dfa
			}
		}
	}
}

// AbbadingoDataset returns a training and testing Dataset using the
// Abbadingo protocol given a DFA, a sparsity percentage and a training:testing ratio.
func AbbadingoDataset(dfa DFA, sparsityPercentage float64, testingRatio float64) (Dataset, Dataset) {
	// Calculate the length of the longest string.
	maxLength := math.Round((2.0 * math.Log2(float64(len(dfa.States)))) + 3.0)
	// Calculate the number which represents the longest string.
	maxDecimal := math.Pow(2, maxLength+1) - 1
	// Calculate the total size of the dataset.
	totalSetSize := math.Round((sparsityPercentage / 100) * maxDecimal)
	// Calculate the size of the training dataset.
	trainingSetSize := int(math.Round((1 - testingRatio) * totalSetSize))

	// Return the datasets generated by AbbadingoDatasetExact using the training set and
	// testing set sizes found above.
	return AbbadingoDatasetExact(dfa, trainingSetSize, int(totalSetSize)-trainingSetSize)
}

// AbbadingoDatasetExact returns a training and testing Dataset using the
// Abbadingo protocol given a DFA and a set size for each.
func AbbadingoDatasetExact(dfa DFA, trainingSetSize int, testingSetSize int) (Dataset, Dataset) {
	// Initialize training dataset.
	trainingDataset := Dataset{}
	// Initialize testing dataset.
	testingDataset := Dataset{}
	// Calculate the length of the longest string.
	maxLength := math.Round((2.0 * math.Log2(float64(len(dfa.States)))) + 3.0)
	// Calculate the number which represents the longest string.
	maxDecimal := math.Pow(2, maxLength+1) - 1

	// Value map to avoid duplicate values.
	valueMap := map[int]bool{}

	// Iterate until both training and testing datasets are filled.
	for x := 0; x < (trainingSetSize + testingSetSize); x++ {
		// Get random value from range [1, totalSetSize].
		value := rand.Intn(int(maxDecimal)) + 1

		// If value is duplicate, decrement x and go to next loop
		// else write new value to map.
		if valueMap[value] {
			x--
			continue
		} else {
			valueMap[value] = true
		}

		// Convert value to binary string.
		binaryString := strconv.FormatInt(int64(value), 2)
		// Remove first '1' within binary string.
		binaryString = binaryString[1:]

		// If training dataset is filled, add string to testing dataset.
		// BinaryStringToStringInstance is used to convert binary string to
		// correct string instance within given DFA.
		if trainingDataset.Count() < trainingSetSize {
			trainingDataset = append(trainingDataset, BinaryStringToStringInstance(dfa, binaryString))
		} else {
			testingDataset = append(testingDataset, BinaryStringToStringInstance(dfa, binaryString))
		}
	}

	// Return populated training and testing datasets.
	return trainingDataset, testingDataset
}

// AbbadingoInstance returns a random DFA using the Abbadingo protocol given a number of states while
// returning a training and testing dataset built on the generated DFA given a sparsity percentage and
// a training:testing ratio. If exact is set to true, the resultant DFA will have the required depth as
// per Abbadingo protocol.
func AbbadingoInstance(numberOfStates int, exact bool, sparsityPercentage float64, testingRatio float64) (DFA, Dataset, Dataset) {
	dfa := AbbadingoDFA(numberOfStates, exact)
	trainingSet, testingSet := AbbadingoDataset(dfa, sparsityPercentage, testingRatio)

	return dfa, trainingSet, testingSet
}

// AbbadingoInstanceExact returns a random DFA using the Abbadingo protocol given a number of states while
// returning a training and testing dataset built on the generated DFA given a set size for each.
// If exact is set to true, the resultant DFA will have the required depth as per Abbadingo protocol.
func AbbadingoInstanceExact(numberOfStates int, exact bool, trainingSetSize int, testingSetSize int) (DFA, Dataset, Dataset) {
	dfa := AbbadingoDFA(numberOfStates, exact)
	trainingSet, testingSet := AbbadingoDatasetExact(dfa, trainingSetSize, testingSetSize)

	return dfa, trainingSet, testingSet
}

// ToAbbadingoFile writes a given Dataset to file in Abbadingo-Format.
func (dataset Dataset) ToAbbadingoFile(filePath string) {
	// Sort the dataset by length.
	sortedDataset := dataset.SortDatasetByLength()

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
	for _, stringInstance := range sortedDataset {
		// Add string label and string length to output string.
		outputString := ""
		if stringInstance.Accepting {
			outputString = strconv.Itoa(1) + " " + strconv.Itoa(stringInstance.Length()) + " "
		} else {
			outputString = strconv.Itoa(0) + " " + strconv.Itoa(stringInstance.Length()) + " "
		}

		// Iterate over the value of string and add to output string.
		for _, symbol := range stringInstance.Value {
			outputString += strconv.Itoa(symbol) + " "
		}
		// Remove trailing space from output string and add new line char.
		outputString = strings.TrimSuffix(outputString, " ") + "\n"
		// Write output string to file.
		_, _ = writer.WriteString(outputString)
	}

	// Flush writer.
	_ = writer.Flush()
}
