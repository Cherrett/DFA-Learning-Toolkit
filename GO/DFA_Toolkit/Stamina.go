package dfatoolkit

import (
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
		for i := 0; i < alphabetSize; i++{
			dfa.AddSymbol()
		}

		// Create first state.
		if rand.Intn(2) == 0 {
			dfa.AddState(ACCEPTING)
		} else {
			dfa.AddState(UNKNOWN)
		}

		// Iterate until required amount of states is reached.
		for len(dfa.States) < attemptedDFASize{
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
			for dfa.States[ambassadorNode].AllTransitionsExist(){
				// Randomly choose an ambassador node.
				ambassadorNode = rand.Intn(len(dfa.States))
			}

			// Select a random symbol within the alphabet.
			randomSymbol := rand.Intn(alphabetSize)

			// If transition from ambassador node using chosen symbol already exists,
			// choose another random symbol. Loop until a free symbol is found.
			for dfa.States[ambassadorNode].Transitions[randomSymbol] != -1{
				randomSymbol = rand.Intn(alphabetSize)
			}

			// Add edge (transition) from  ambassador node to new state using random valid symbol.
			dfa.AddTransition(randomSymbol, ambassadorNode, newStateID)

			// Self loop with probability l.
			if rand.Float64() < l{
				// Select a random symbol within the alphabet.
				randomSymbol = rand.Intn(alphabetSize)

				// If transition from current node using chosen symbol already exists,
				// choose another random symbol. Loop until a free symbol is found.
				for dfa.States[newStateID].Transitions[randomSymbol] != -1{
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
		dfa.Describe(false)
	}
}

// ModifiedForestFire is used within the StaminaDFA function to generate random DFAs using the Stamina protocol.
// This is equivalent to steps 2 an 3 within the Forest-Fire algorithm in
// Jure Leskovec, Jon Kleinberg, and Christos Faloutsos. 2007. Graph evolution:
// Densification and shrinking diameters. ACM Trans. Knowl. Discov. Data 1, 1
// (March 2007), 2–es. DOI:https://doi.org/10.1145/1217299.1217301
func ModifiedForestFire(currentNode int, ambassadorNode int, dfa *DFA, visitedStates *map[int]bool, alphabetSize int, f, b, p float64){
	// Step 1 within ForestFire algorithm.

	// If all transitions from current node are already
	// filled, return.
	if dfa.States[currentNode].AllTransitionsExist(){
		return
	}

	// Generate random number using a geometric distribution as per Stamina protocol.
	x := int(math.Ceil(math.Logb(1 - rand.Float64()) / (math.Logb(1 - (f / (1 - f))))))

	// Slice to store state IDs of nodes which have a transition to ambassador node.
	var fromNodes []int

	// Iterate over states within dfa.
	for stateID, state := range dfa.States{
		// If the amount of required nodes
		// is reached, break loop.
		if len(fromNodes) == x{
			break
		}

		// Check if state already visited and that not all transitions exist.
		if !(*visitedStates)[stateID] && !state.AllTransitionsExist(){
			// Iterate over each symbol within alphabet.
			for symbol := range dfa.Alphabet{
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
	y := int(math.Ceil(math.Logb(1 - rand.Float64()) / (math.Logb(1 - ((f * b) / (1 - (f * b)))))))

	// Map to store state IDs of nodes which have a transition from ambassador node.
	// A map is used rather than a slice to avoid duplicate state IDs.
	toNodes := map[int]bool{}

	// Iterate over each symbol within alphabet.
	for symbol := range dfa.Alphabet{
		// If the amount of required nodes
		// is reached, break loop.
		if len(toNodes) == y{
			break
		}

		// If resultant state is valid, not visited, not in to nodes and not all transitions exist, add to toNodes map.
		if resultantStateID := dfa.States[ambassadorNode].Transitions[symbol]; resultantStateID != -1 &&
			!(*visitedStates)[resultantStateID] && !toNodes[resultantStateID] && !dfa.States[resultantStateID].AllTransitionsExist(){
			toNodes[resultantStateID] = true
		}
	}

	// Step 2 within ForestFire algorithm.

	// Iterate over all from nodes selected.
	for _, stateID := range fromNodes{
		// Select a random symbol within the alphabet.
		randomSymbol := rand.Intn(alphabetSize)

		// If transition from state using chosen symbol already exists,
		// choose another random symbol. Loop until a free symbol is found.
		for dfa.States[stateID].Transitions[randomSymbol] != -1{
			randomSymbol = rand.Intn(alphabetSize)
		}

		// Add edge (transition) from current state to from state using random valid symbol.
		dfa.AddTransition(randomSymbol, currentNode, stateID)

		// Parallel edge label with probability l.
		if !dfa.States[currentNode].AllTransitionsExist() && rand.Float64() < p{
			// Select a random symbol within the alphabet.
			randomSymbol = rand.Intn(alphabetSize)

			// If transition from current node using chosen symbol already exists,
			// choose another random symbol. Loop until a free symbol is found.
			for dfa.States[currentNode].Transitions[randomSymbol] != -1{
				randomSymbol = rand.Intn(alphabetSize)
			}

			// Add parallel edge (transition) label using random valid symbol.
			dfa.AddTransition(randomSymbol, currentNode, stateID)
		}
	}

	// Iterate over all to nodes selected.
	for stateID := range toNodes{
		// Skip if all transitions already exist within
		// state. This happens since other recursive
		// calls can change these transitions.
		if dfa.States[stateID].AllTransitionsExist(){
			continue
		}

		// Select a random symbol within the alphabet.
		randomSymbol := rand.Intn(alphabetSize)

		// If transition from state using chosen symbol already exists,
		// choose another random symbol. Loop until a free symbol is found.
		for dfa.States[stateID].Transitions[randomSymbol] != -1{
			randomSymbol = rand.Intn(alphabetSize)
		}

		// Add edge (transition) from current state to from state using random valid symbol.
		dfa.AddTransition(randomSymbol, currentNode, stateID)

		// Parallel edge label with probability l.
		if !dfa.States[currentNode].AllTransitionsExist() && rand.Float64() < p{
			// Select a random symbol within the alphabet.
			randomSymbol = rand.Intn(alphabetSize)

			// If transition from current node using chosen symbol already exists,
			// choose another random symbol. Loop until a free symbol is found.
			for dfa.States[currentNode].Transitions[randomSymbol] != -1{
				randomSymbol = rand.Intn(alphabetSize)
			}

			// Add parallel edge (transition) label using random valid symbol.
			dfa.AddTransition(randomSymbol, currentNode, stateID)
		}
	}

	// Iterate over all from nodes selected.
	for _, stateID := range fromNodes{
		// Mark current node as visited.
		(*visitedStates)[stateID] = true
		// Recursively call ModifiedForestFire function.
		ModifiedForestFire(stateID, ambassadorNode, dfa, visitedStates, alphabetSize, f, b, p)
	}

	// Iterate over all to nodes selected.
	for stateID := range toNodes{
		// Mark current node as visited.
		(*visitedStates)[stateID] = true
		// Recursively call ModifiedForestFire function.
		ModifiedForestFire(stateID, ambassadorNode, dfa, visitedStates, alphabetSize, f, b, p)
	}
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
