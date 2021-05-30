package dfalearningtoolkit

import (
	"bufio"
	"fmt"
	"github.com/Cherrett/DFA-Learning-Toolkit/util"
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

	// Stop if no more characters exist in string. (Empty string case)
	if len(splitString) > 1 && splitString[1] != "" {
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
	// Forward-burning probability.
	f := 0.31
	// Backward-burning ratio.
	b := 0.385
	// Probability that a state loops to itself.
	l := 0.2
	// Probability of parallel transition (edge) labels.
	p := 0.2

	// Initialize a new DFA.
	dfa := NewDFA()

	// Add symbols to DFA until alphabet size is reached.
	for i := 0; i < alphabetSize; i++ {
		dfa.AddSymbol()
	}

	// Create and set starting state.
	if rand.Intn(2) == 0 {
		dfa.StartingStateID = dfa.AddState(ACCEPTING)
	} else {
		dfa.StartingStateID = dfa.AddState(UNLABELLED)
	}

	// Vertices added counter.
	numberOfStatesToAdd := targetDFASize
	minimisedDFAStates := 0

	// Iterate till a valid DFA is generated.
	for {
		// Iterate until required amount of states is reached.
		for len(dfa.States)-1 < numberOfStatesToAdd {
			// Map of visited states.
			visitedStates := map[int]bool{}

			// Newly created state ID placeholder.
			var newStateID int

			// Create new state and randomly choose whether
			// state is accepting or unlabelled.
			if rand.Intn(2) == 0 {
				newStateID = dfa.AddState(ACCEPTING)
			} else {
				newStateID = dfa.AddState(UNLABELLED)
			}

			// Ambassador state ID placeholder.
			ambassadorState := -1
			// Map of tried states to be ambassador state.
			tried := map[int]bool{newStateID: true}

			// Iterate until a valid ambassador state is found.
			for {
				// Slice of state pool for ambassador state selection.
				var statePool []int

				// Iterate over states within dfa.
				for stateID := range dfa.States {
					// Add to state pool if not already tried.
					if _, alreadyTried := tried[stateID]; !alreadyTried {
						statePool = append(statePool, stateID)
					}
				}

				// Recursively call StaminaDFA if state pool is empty.
				if len(statePool) == 0 {
					return StaminaDFA(alphabetSize, targetDFASize)
				}

				// Randomly select an ambassador state from state pool.
				ambassadorState = statePool[rand.Intn(len(statePool))]

				// If ambassador state already covers all transitions,
				// choose another ambassador state. Else break since a
				// valid ambassador state is found.
				if !dfa.States[ambassadorState].AllTransitionsExist() {
					break
				}

				// Add newly selected ambassador state to tried map.
				tried[ambassadorState] = true
			}

			// Add an edge from the ambassador state to the newly created state.
			addTransition(ambassadorState, newStateID, &dfa, p)

			// Self loop with probability l.
			if util.RandomGeometricProbability(1-l) > 0 && !dfa.States[newStateID].AllTransitionsExist() {
				// Add a self-loop edge within the newly created state.
				addTransition(newStateID, newStateID, &dfa, p)
			}

			// Mark both ambassador state and new state as visited.
			visitedStates[ambassadorState] = true
			visitedStates[newStateID] = true

			// modifiedForestFire is called using the new state, the ambassador state, a pointer
			// to the DFA and to the visited states map, alongside other probability values.
			modifiedForestFire(newStateID, ambassadorState, &dfa, &visitedStates, f, b, p)
		}

		// Minimise DFA created and remove sink state.
		minimisedDFA := dfa.minimiseAndRemoveSinkState()

		// Store number of states within minimised DFA.
		minimisedDFAStates = len(minimisedDFA.States)
		// Increment number of states to add by decrementing the number of
		// states within minimised dfa from the size of the target DFA.
		numberOfStatesToAdd += targetDFASize - minimisedDFAStates

		// If the number of states within minimised DFA is equal (or within +2)
		// to the target DFA size, the DFA is returned. Else, try again until
		// required target is found.
		if minimisedDFAStates+2 >= targetDFASize {
			// Return the minimised DFA since it
			// meets all of the requirements.
			return minimisedDFA
		}
	}
}

// modifiedForestFire is used within the StaminaDFA function to generate random DFAs using the Stamina protocol.
// This is equivalent to steps 2 an 3 within the Forest-Fire algorithm in
// Jure Leskovec, Jon Kleinberg, and Christos Faloutsos. 2007. Graph evolution:
// Densification and shrinking diameters. ACM Trans. Knowl. Discov. Data 1, 1
// (March 2007), 2–es. DOI:https://doi.org/10.1145/1217299.1217301
func modifiedForestFire(currentState int, ambassadorState int, dfa *DFA, visitedStates *map[int]bool, f, b, p float64) {
	// Generate random number using a geometric distribution as per Stamina protocol.
	x := util.RandomGeometricProbability(1 - f)

	// Generate random number using a geometric distribution as per Stamina protocol.
	y := util.RandomGeometricProbability(1 - f*b)

	// Slice to store IDs of states which have a transition from ambassador state.
	// A map is also created to avoid duplicate state IDs.
	toStatesSet := map[int]bool{}
	var toStates []int

	// Iterate over each symbol within alphabet.
	for symbol := range dfa.Alphabet {
		// If resultant state is valid, not visited, not in to states, and not
		// equal to current state, add to toStatesSet map and to toStates slice.
		if resultantStateID := dfa.States[ambassadorState].Transitions[symbol]; resultantStateID >= 0 &&
			!(*visitedStates)[resultantStateID] && !toStatesSet[resultantStateID] {
			toStatesSet[resultantStateID] = true
			toStates = append(toStates, resultantStateID)
		}
	}

	// Shuffle to states selected and limit number of states to y.
	if len(toStates) > y {
		rand.Shuffle(len(toStates), func(i, j int) { toStates[i], toStates[j] = toStates[j], toStates[i] })
		toStates = toStates[:y]
	}

	// Slice to store IDs of states which have a transition to the ambassador state.
	var fromStates []int

	// Iterate over states within dfa.
	for stateID, state := range dfa.States {
		// Check if state already visited.
		if !(*visitedStates)[stateID] {
			// Iterate over each symbol within alphabet.
			for symbol := range dfa.Alphabet {
				// If resultant state is equal to ID of ambassador state and not equal
				// to current state, add to fromStates slice and break.
				if resultantStateID := state.Transitions[symbol]; resultantStateID == ambassadorState && !toStatesSet[stateID] {
					fromStates = append(fromStates, stateID)
					break
				}
			}
		}
	}

	// Shuffle from states selected and limit number of states to x.
	if len(fromStates) > x {
		rand.Shuffle(len(fromStates), func(i, j int) { fromStates[i], fromStates[j] = fromStates[j], fromStates[i] })
		fromStates = fromStates[:x]
	}

	// Iterate over all from states selected.
	for _, fromStateID := range fromStates {
		// Randomly choose direction of transition.
		if rand.Intn(2) == 0 {
			// Check whether any transitions are available from current state.
			if !dfa.States[currentState].AllTransitionsExist() {
				// Add edge from current state to 'from' state.
				addTransition(currentState, fromStateID, dfa, p)
			}
		} else {
			// Check whether any transitions are available from 'from' state.
			if !dfa.States[fromStateID].AllTransitionsExist() {
				// Add edge from 'from' state to current state.
				addTransition(fromStateID, currentState, dfa, p)
			}
		}
		// Mark state as visited.
		(*visitedStates)[fromStateID] = true
	}

	// Iterate over all to states selected.
	for _, toStateID := range toStates {
		// Randomly choose direction of transition.
		if rand.Intn(2) == 0 {
			// Check whether any transitions are available from current state.
			if !dfa.States[currentState].AllTransitionsExist() {
				// Add edge from current state to 'to' state.
				addTransition(currentState, toStateID, dfa, p)
			}
		} else {
			// Check whether any transitions are available from 'to' state.
			if !dfa.States[toStateID].AllTransitionsExist() {
				// Add edge from 'to' state to current state.
				addTransition(toStateID, currentState, dfa, p)
			}
		}
		// Mark state as visited.
		(*visitedStates)[toStateID] = true
	}

	// Iterate over all from states selected.
	for _, fromStateID := range fromStates {
		// Recursively call modifiedForestFire function on 'from' state.
		modifiedForestFire(currentState, fromStateID, dfa, visitedStates, f, b, p)
	}

	// Iterate over all to statesselected.
	for _, toStateID := range toStates {
		// Recursively call modifiedForestFire function on 'to' state.
		modifiedForestFire(currentState, toStateID, dfa, visitedStates, f, b, p)
	}
}

// addTransition adds a transition from one state to another (can be the same state) while
// possibly adding a number of parallel edges using mean p.
// Used within StaminaDFA and modifiedForestFire functions.
func addTransition(fromState, toState int, dfa *DFA, p float64) {
	// Call addTransitionInternal function to add a transition
	// from 'from' state to 'to' state.
	addTransitionInternal(fromState, toState, dfa)

	// Parallel edge label with mean p.
	parallelEdges := util.RandomGeometricProbability(1 - p)

	// Randomly choose direction of transition.
	if rand.Intn(2) == 0 {
		// Check whether number of parallel edges is reached and if any transitions are available from 'from' state.
		for i := 0; i < parallelEdges && !dfa.States[fromState].AllTransitionsExist(); i++ {
			// Call addTransitionInternal function to add a transition
			// from 'from' state to 'to' state.
			addTransitionInternal(fromState, toState, dfa)
		}
	} else {
		// Check whether number of parallel edges is reached and if any transitions are available from 'to' state.
		for i := 0; i < parallelEdges && !dfa.States[toState].AllTransitionsExist(); i++ {
			// Call addTransitionInternal function to add a transition
			// from 'to' state to 'from' state.
			addTransitionInternal(toState, fromState, dfa)
		}
	}
}

// addTransitionInternal adds a transition from one state to another (can be the same state)
// Used within addTransition function which is used within StaminaDFA and modifiedForestFire functions.
func addTransitionInternal(fromState, toState int, dfa *DFA) {
	// Slice to store alphabet pool of 'from' state.
	var alphabetPool []int

	// Iterate over each symbol within alphabet.
	for symbolID := range dfa.Alphabet {
		// If resultant state ID is smaller than 0 (does not exist), add to alphabet pool.
		if resultantStateID := dfa.States[fromState].Transitions[symbolID]; resultantStateID < 0 {
			alphabetPool = append(alphabetPool, symbolID)
		}
	}

	// Panic if alphabet pool is empty.
	if len(alphabetPool) == 0 {
		panic("Alphabet pool cannot be empty.")
	}

	// Randomly select a symbol from alphabet pool.
	randomSymbol := alphabetPool[rand.Intn(len(alphabetPool))]

	// Add transition from 'from' state to 'to' state using selected symbol.
	dfa.AddTransition(randomSymbol, fromState, toState)
}

// minimiseAndRemoveSinkState is a modified version of the Minimise function. If a non-complete
// DFA is inputted, a sink state is added and it is minimised. After this, the sink state is removed.
// Used within StaminaDFA function.
func (dfa DFA) minimiseAndRemoveSinkState() DFA {
	// Clone DFA to work on.
	temporaryDFA := dfa.Clone()

	// Remove unreachable states within DFA.
	temporaryDFA.RemoveUnreachableStates()

	// If DFA is complete, use normal minimisation function.
	if temporaryDFA.IsComplete() {
		return temporaryDFA.Minimise()
	}

	// Add sink state and store its ID.
	sinkStateID := temporaryDFA.AddSinkState()

	// Get indistinguishable state pairs from IndistinguishableStatePairs function.
	indistinguishablePairs := temporaryDFA.IndistinguishableStatePairs()

	// Convert DFA to state partition.
	statePartition := temporaryDFA.ToStatePartition()

	// Merge indistinguishable pairs.
	for _, indistinguishablePair := range indistinguishablePairs {
		block1 := statePartition.Find(indistinguishablePair.state1)
		block2 := statePartition.Find(indistinguishablePair.state2)
		if block1 != block2 {
			statePartition.Union(block1, block2)
		}
	}

	// Get the block ID which contains the sink state after merging states.
	sinkBlockID := statePartition.Find(sinkStateID)

	// Convert state partition to Quotient DFA while getting the mappings used.
	resultantDFA, blockToStateMap := statePartition.ToQuotientDFAWithMapping()

	// Remove state which contains the sink state.
	resultantDFA.RemoveState(blockToStateMap[sinkBlockID])

	// Return resultant minimised DFA.
	return resultantDFA
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
				if rand.Float64() < 1.0/float64(1+(2*currentState.OutDegree())) {
					break
				}
			}

			// Break if out degree of current state is equal to 0.
			if currentState.OutDegree() == 0 {
				break
			}

			validTransitions := currentState.ValidTransitions()

			// Randomly choose a symbol with a valid transitions.
			validSymbol := validTransitions[rand.Intn(len(validTransitions))]

			currentString.Value = append(currentString.Value, validSymbol)

			currentState = &dfa.States[currentState.Transitions[validSymbol]]
		}

		if currentState.IsAccepting() {
			// Add string to first sample.
			firstSample = append(firstSample, currentString)
		}
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

	// Iterate over each string instance within sorted dataset.
	for _, stringInstance := range dataset {
		// Add string label and string length to output string.
		outputString := ""
		if stringInstance.Accepting {
			outputString = "+ "
		} else {
			outputString = "- "
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
			dfa.AddState(UNLABELLED)
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
