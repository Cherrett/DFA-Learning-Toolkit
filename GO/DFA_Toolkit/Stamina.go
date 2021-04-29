package dfatoolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"bufio"
	"fmt"
	"math"
	"math/rand"
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
	for i := 1; i < len(splitString); i++ {
		// Get integer value from split string.
		integerValue, err := strconv.Atoi(splitString[i])

		// Panic if error.
		if err != nil {
			panic(err)
		}

		// Add integer value to value slice within string instance.
		stringInstance.Value = append(stringInstance.Value, integerValue)
	}

	// Return populated string instance.
	return stringInstance
}

// StaminaDFA returns a random DFA using the Stamina protocol given a number of states.
// This process is described in http://stamina.chefbe.net/machines and in
// Walkinshaw, N., Lambeau, B., Damas, C., Bogdanov, K., & Dupont, P. (2012). STAMINA:
// a competition to encourage the development and assessment of software model inference
// techniques. Empirical Software Engineering, 18(4), 791–824. doi:10.1007/s10664-012-9210-3
func StaminaDFA(alphabetSize int, targetDFASize int) DFA {
	// Get attempted DFA size by increasing the target DFA size by 25%.
	// This is done since targetDFASize cannot be reached using original
	// value since DFA is minimized after the process is finished.
	attemptedDFASize := int(math.Ceil(float64(targetDFASize) * 1.25))

	// Forward-burning probability.
	f := 0.31
	// Backward-burning ratio.
	b := 0.385
	// Probability that a state loops to itself.
	l := 0.2
	// Probability of parallel transition (edge) labels.
	p := 0.2

	// Iterate till a valid DFA is generated.
	for {
		// Initialize a new DFA.
		dfa := NewDFA()

		// Add symbols to DFA until alphabet size is reached.
		for i := 0; i < alphabetSize; i++ {
			dfa.AddSymbol()
		}

		// Create first state.
		if rand.Intn(2) == 0 {
			dfa.AddState(ACCEPTING)
		} else {
			dfa.AddState(UNKNOWN)
		}

		// Iterate until required amount of states is reached.
		for len(dfa.States) < attemptedDFASize {
			// Step 1 within ForestFire algorithm.

			// Map of visited states.
			visitedStates := map[int]bool{}

			newStateID := -1
			// Create new state.
			if rand.Intn(2) == 0 {
				newStateID = dfa.AddState(ACCEPTING)
			} else {
				newStateID = dfa.AddState(UNKNOWN)
			}

			// Randomly choose an ambassador node.
			ambassadorNode := rand.Intn(len(dfa.States))

			// If ambassador node already covers all transitions,
			// choose another ambassador node. Loop until a valid
			// ambassador node is found.
			for dfa.States[ambassadorNode].AllTransitionsExist() {
				// Randomly choose an ambassador node.
				ambassadorNode = rand.Intn(len(dfa.States))
			}

			// Select a random symbol within the alphabet.
			randomSymbol := rand.Intn(alphabetSize)

			// If transition from ambassador node using chosen symbol already exists,
			// choose another random symbol. Loop until a free symbol is found.
			for dfa.States[ambassadorNode].Transitions[randomSymbol] != -1 {
				randomSymbol = rand.Intn(alphabetSize)
			}

			// Add edge (transition) from  ambassador node to new state using random valid symbol.
			dfa.AddTransition(randomSymbol, ambassadorNode, newStateID)

			// Self loop with probability l.
			if rand.Float64() < l {
				// Select a random symbol within the alphabet.
				randomSymbol = rand.Intn(alphabetSize)

				// If transition from current node using chosen symbol already exists,
				// choose another random symbol. Loop until a free symbol is found.
				for dfa.States[newStateID].Transitions[randomSymbol] != -1 {
					randomSymbol = rand.Intn(alphabetSize)
				}

				// Add self-looping edge (transition) using random valid symbol.
				dfa.AddTransition(randomSymbol, newStateID, newStateID)
			}

			// ModifiedForestFire is called using the new state and the ambassador node.
			ModifiedForestFire(newStateID, ambassadorNode, &dfa, &visitedStates, alphabetSize, f, b, p)
		}

		// Randomly choose a starting state.
		dfa.StartingStateID = rand.Intn(len(dfa.States))

		// Minimise DFA created.
		dfa = dfa.Minimise()

		// If the number of states within created DFA is equal
		// to the target DFA size, the DFA is returned.
		// Else, try again until required target is found.
		if len(dfa.States) == targetDFASize {
			// Return the created DFA since it
			// meets all of the requirements.
			return dfa
		}
	}
}

// ModifiedForestFire is used within the StaminaDFA function to generate random DFAs using the Stamina protocol.
// This is equivalent to steps 2 an 3 within the Forest-Fire algorithm in
// Jure Leskovec, Jon Kleinberg, and Christos Faloutsos. 2007. Graph evolution:
// Densification and shrinking diameters. ACM Trans. Knowl. Discov. Data 1, 1
// (March 2007), 2–es. DOI:https://doi.org/10.1145/1217299.1217301
func ModifiedForestFire(currentNode int, ambassadorNode int, dfa *DFA, visitedStates *map[int]bool, alphabetSize int, f, b, p float64) {
	// Step 1 within ForestFire algorithm.

	// If all transitions from current node are already
	// filled, return.
	if dfa.States[currentNode].AllTransitionsExist() {
		return
	}

	// Generate random number using a geometric distribution as per Stamina protocol.
	x := int(math.Ceil(math.Logb(1-rand.Float64()) / (math.Logb(1 - (f / (1 - f))))))

	// Slice to store state IDs of nodes which have a transition to ambassador node.
	var fromNodes []int

	// Iterate over states within dfa.
	for stateID, state := range dfa.States {
		// If the amount of required nodes
		// is reached, break loop.
		if len(fromNodes) == x {
			break
		}

		// Check if state already visited and that not all transitions exist.
		if !(*visitedStates)[stateID] && !state.AllTransitionsExist() {
			// Iterate over each symbol within alphabet.
			for symbol := range dfa.Alphabet {
				// If resultant state is equal to ID of ambassador node, add to fromNodes slice and
				// break.
				if resultantStateID := state.Transitions[symbol]; resultantStateID == ambassadorNode {
					fromNodes = append(fromNodes, stateID)
					break
				}
			}
		}
	}

	// Generate random number using a geometric distribution as per Stamina protocol.
	y := int(math.Ceil(math.Logb(1-rand.Float64()) / (math.Logb(1 - ((f * b) / (1 - (f * b)))))))

	// Map to store state IDs of nodes which have a transition from ambassador node.
	// A map is used rather than a slice to avoid duplicate state IDs.
	toNodes := map[int]bool{}

	// Iterate over each symbol within alphabet.
	for symbol := range dfa.Alphabet {
		// If the amount of required nodes
		// is reached, break loop.
		if len(toNodes) == y {
			break
		}

		// If resultant state is valid, not visited, not in to nodes and not all transitions exist, add to toNodes map.
		if resultantStateID := dfa.States[ambassadorNode].Transitions[symbol]; resultantStateID != -1 &&
			!(*visitedStates)[resultantStateID] && !toNodes[resultantStateID] && !dfa.States[resultantStateID].AllTransitionsExist() {
			toNodes[resultantStateID] = true
		}
	}

	// Step 2 within ForestFire algorithm.

	// Iterate over all from nodes selected.
	for _, stateID := range fromNodes {
		// Select a random symbol within the alphabet.
		randomSymbol := rand.Intn(alphabetSize)

		// If transition from state using chosen symbol already exists,
		// choose another random symbol. Loop until a free symbol is found.
		for dfa.States[stateID].Transitions[randomSymbol] != -1 {
			randomSymbol = rand.Intn(alphabetSize)
		}

		// Add edge (transition) from current state to from state using random valid symbol.
		dfa.AddTransition(randomSymbol, currentNode, stateID)

		// Parallel edge label with probability l.
		if !dfa.States[currentNode].AllTransitionsExist() && rand.Float64() < p {
			// Select a random symbol within the alphabet.
			randomSymbol = rand.Intn(alphabetSize)

			// If transition from current node using chosen symbol already exists,
			// choose another random symbol. Loop until a free symbol is found.
			for dfa.States[currentNode].Transitions[randomSymbol] != -1 {
				randomSymbol = rand.Intn(alphabetSize)
			}

			// Add parallel edge (transition) label using random valid symbol.
			dfa.AddTransition(randomSymbol, currentNode, stateID)
		}
	}

	// Iterate over all to nodes selected.
	for stateID := range toNodes {
		// Skip if all transitions already exist within
		// state. This happens since other recursive
		// calls can change these transitions.
		if dfa.States[stateID].AllTransitionsExist() {
			continue
		}

		// Select a random symbol within the alphabet.
		randomSymbol := rand.Intn(alphabetSize)

		// If transition from state using chosen symbol already exists,
		// choose another random symbol. Loop until a free symbol is found.
		for dfa.States[stateID].Transitions[randomSymbol] != -1 {
			randomSymbol = rand.Intn(alphabetSize)
		}

		// Add edge (transition) from current state to from state using random valid symbol.
		dfa.AddTransition(randomSymbol, currentNode, stateID)

		// Parallel edge label with probability l.
		if !dfa.States[currentNode].AllTransitionsExist() && rand.Float64() < p {
			// Select a random symbol within the alphabet.
			randomSymbol = rand.Intn(alphabetSize)

			// If transition from current node using chosen symbol already exists,
			// choose another random symbol. Loop until a free symbol is found.
			for dfa.States[currentNode].Transitions[randomSymbol] != -1 {
				randomSymbol = rand.Intn(alphabetSize)
			}

			// Add parallel edge (transition) label using random valid symbol.
			dfa.AddTransition(randomSymbol, currentNode, stateID)
		}
	}

	// Iterate over all from nodes selected.
	for _, stateID := range fromNodes {
		// Mark current node as visited.
		(*visitedStates)[stateID] = true
		// Recursively call ModifiedForestFire function.
		ModifiedForestFire(stateID, ambassadorNode, dfa, visitedStates, alphabetSize, f, b, p)
	}

	// Iterate over all to nodes selected.
	for stateID := range toNodes {
		// Mark current node as visited.
		(*visitedStates)[stateID] = true
		// Recursively call ModifiedForestFire function.
		ModifiedForestFire(stateID, ambassadorNode, dfa, visitedStates, alphabetSize, f, b, p)
	}
}

// StaminaDataset returns a training and testing Dataset using the
// Stamina protocol given a DFA and a sparsity percentage.
// This process is described in http://stamina.chefbe.net/samples and in
// Walkinshaw, N., Lambeau, B., Damas, C., Bogdanov, K., & Dupont, P. (2012). STAMINA:
// a competition to encourage the development and assessment of software model inference
// techniques. Empirical Software Engineering, 18(4), 791–824. doi:10.1007/s10664-012-9210-3
func StaminaDataset(dfa DFA, sparsityPercentage float64, initialStringsGenerated, maximumTestingStringInstances int) (Dataset, Dataset) {
	// Initialize training dataset.
	trainingDataset := Dataset{}
	// Initialize testing dataset.
	testingDataset := Dataset{}

	// Step 1 - First Sample
	// Initialize first sample.
	firstSample := Dataset{}

	// Positive strings into first sample.
	for len(firstSample) < initialStringsGenerated {
		currentString := StringInstance{make([]int, 0), true}
		currentState := dfa.StartingState()
		for {
			if currentState.IsAccepting() {
				// End generation with probability 1.0/(1 + 2*outdegree(v).
				if rand.Float64() < 1.0/float64(1+(2*currentState.TotalTransitionsCount())) {
					break
				}
			}

			validTransitions := currentState.ValidTransitions()

			// Randomly choose a symbol with a valid transitions.
			validSymbol := validTransitions[rand.Intn(len(validTransitions))]

			currentString.Value = append(currentString.Value, validSymbol)

			currentState = &dfa.States[currentState.Transitions[validSymbol]]
		}

		// Add string to first sample.
		firstSample = append(firstSample, currentString)
	}

	// Negative strings into first sample.
	for i := 0; i < initialStringsGenerated/2; i++ {
		stringToBeChanged := &firstSample[i]

		// Poisson distribution (with a mean of 3).
		randomValue := rand.Float64()
		sum := 0.0
		numberOfEditingOperations := 0

		for {
			sum += (math.Pow(3, float64(numberOfEditingOperations)) * math.Pow(math.E, -3)) / float64(util.Factorial(numberOfEditingOperations))

			if randomValue < sum {
				break
			}

			numberOfEditingOperations++
		}

		for j := 0; j < numberOfEditingOperations; j++ {
			if len(stringToBeChanged.Value) == 0 {
				break
			}
			switch rand.Intn(3) {
			// Substitution
			case 0:
				substitutionPosition := rand.Intn(len(stringToBeChanged.Value))
				substitutionSymbol := rand.Intn(len(dfa.Alphabet))
				stringToBeChanged.Value[substitutionPosition] = substitutionSymbol
				break
			// Insertion
			case 1:
				insertionPosition := rand.Intn(len(stringToBeChanged.Value))
				insertionSymbol := rand.Intn(len(dfa.Alphabet))
				stringToBeChanged.Value = append(stringToBeChanged.Value[:insertionPosition+1], stringToBeChanged.Value[insertionPosition:]...)
				stringToBeChanged.Value[insertionPosition] = insertionSymbol
				break
			// Deletion
			case 2:
				deletionPosition := rand.Intn(len(stringToBeChanged.Value))
				stringToBeChanged.Value = append(stringToBeChanged.Value[:deletionPosition], stringToBeChanged.Value[deletionPosition+1:]...)
				break
			}
		}

		if stringToBeChanged.ParseToStateLabel(dfa) != ACCEPTING {
			stringToBeChanged.Accepting = false
		}
	}

	// Step 2 - Split first sample.
	// Initialize 2 sample sets.
	firstSet := Dataset{}
	secondSet := Dataset{}
	acceptingStringInstances := firstSample.AcceptingStringInstances()
	rejectingStringInstances := firstSample.RejectingStringInstances()

	for _, stringInstance := range acceptingStringInstances {
		if len(firstSet) < len(acceptingStringInstances)/2 {
			firstSet = append(firstSet, stringInstance)
		} else {
			secondSet = append(secondSet, stringInstance)
		}
	}

	for _, stringInstance := range rejectingStringInstances {
		if len(firstSet) < (len(acceptingStringInstances)+len(rejectingStringInstances))/2 {
			firstSet = append(firstSet, stringInstance)
		} else {
			secondSet = append(secondSet, stringInstance)
		}
	}

	// Step 4 - Populate Training Sample.
	requiredStrings := float64(len(secondSet)) * (sparsityPercentage / 100)
	acceptingStringInstances = secondSet.AcceptingStringInstances()
	rejectingStringInstances = secondSet.RejectingStringInstances()
	counter := 0

	for float64(len(trainingDataset)) < requiredStrings {
		if counter < len(acceptingStringInstances) {
			trainingDataset = append(trainingDataset, acceptingStringInstances[counter])
		}

		if counter < len(rejectingStringInstances) {
			trainingDataset = append(trainingDataset, rejectingStringInstances[counter])
		}

		counter++
	}

	// Step 3 - Populate Test Sample.
	for len(testingDataset) < maximumTestingStringInstances && len(firstSet) > 0 {
		randomIndex := rand.Intn(len(firstSet))
		randomString := firstSet[randomIndex]
		firstSet = append(firstSet[:randomIndex], firstSet[randomIndex+1:]...)

		if randomString.WithinDataset(trainingDataset) || randomString.WithinDataset(testingDataset) {
			continue
		}

		testingDataset = append(testingDataset, randomString)
	}

	// Return populated training and testing datasets.
	return trainingDataset, testingDataset
}

// DefaultStaminaDataset returns a training and testing Dataset using the
// Stamina protocol given a DFA and a sparsity percentage.
// The default values for the initial strings generated (20000) and the maximum testing string
// instances (1500) which were used within the Stamina competition are used within this function.
func DefaultStaminaDataset(dfa DFA, sparsityPercentage float64) (Dataset, Dataset) {
	return StaminaDataset(dfa, sparsityPercentage, 20000, 1500)
}

// StaminaInstance returns a random DFA using the Stamina protocol given a target size while
// returning a training and testing dataset built on the generated DFA given a sparsity percentage,
// an initial string generated value, and a maximum testing string instances value.
func StaminaInstance(alphabetSize, targetDFASize int, sparsityPercentage float64, initialStringsGenerated, maximumTestingStringInstances int) (DFA, Dataset, Dataset) {
	dfa := StaminaDFA(alphabetSize, targetDFASize)
	trainingSet, testingSet := StaminaDataset(dfa, sparsityPercentage, initialStringsGenerated, maximumTestingStringInstances)

	return dfa, trainingSet, testingSet
}

// DefaultStaminaInstance returns a random DFA using the Stamina protocol given a target size while
// returning a training and testing dataset built on the generated DFA given a sparsity percentage.
// The default values for the initial strings generated (20000) and the maximum testing string
// instances (1500) which were used within the Stamina competition are used within this function.
func DefaultStaminaInstance(alphabetSize, targetDFASize int, sparsityPercentage float64) (DFA, Dataset, Dataset) {
	dfa := StaminaDFA(alphabetSize, targetDFASize)
	trainingSet, testingSet := DefaultStaminaDataset(dfa, sparsityPercentage)

	return dfa, trainingSet, testingSet
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
	if err != nil {
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
		if err != nil {
			panic(err)
		}

		// If starting state is set to true, set starting state
		// within DFA to current state.
		if initialState {
			// Panic if starting state is already found.
			if dfa.StartingStateID > -1 {
				panic("DFA cannot have more than 2 starting states.")
			}
			dfa.StartingStateID = stateID
		}

		// Check if state is accepting.
		acceptingState, err := strconv.ParseBool(line[1])

		// Panic if error.
		if err != nil {
			panic(err)
		}

		if acceptingState {
			dfa.AddState(ACCEPTING)
		} else {
			dfa.AddState(UNKNOWN)
		}
	}

	// Iterate over remaining lines (transitions).
	for scanner.Scan() {
		// Split read line.
		line := strings.Split(scanner.Text(), " ")

		// Get from state ID of transition.
		fromStateID, err := strconv.Atoi(line[0])

		// Panic if error.
		if err != nil {
			panic(err)
		}

		// Get to state ID of transition.
		toStateID, err := strconv.Atoi(line[1])

		// Panic if error.
		if err != nil {
			panic(err)
		}

		// Get symbol ID of transition.
		symbolID, err := strconv.Atoi(line[2])

		// Panic if error.
		if err != nil {
			panic(err)
		}

		// Add new symbol within alphabet until
		// symbolID is created.
		for symbolID+1 > len(dfa.Alphabet) {
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
	for stateID, state := range dfa.States {
		// Write state information.
		_, _ = writer.WriteString(strconv.Itoa(stateID) + " " +
			strconv.FormatBool(dfa.StartingStateID == stateID) + " " +
			strconv.FormatBool(state.Label == ACCEPTING) + "\n")
	}

	// Iterate over each state within dfa.
	for stateID, state := range dfa.States {
		// Write transitions information.
		for symbolID := range dfa.Alphabet {
			// Write if transition exists.
			if resultantStateID := state.Transitions[symbolID]; resultantStateID > -1 {
				_, _ = writer.WriteString(strconv.Itoa(stateID) + " " +
					strconv.Itoa(resultantStateID) + " " +
					strconv.Itoa(symbolID) + "\n")
			}
		}
	}

	// Flush writer.
	_ = writer.Flush()
}
