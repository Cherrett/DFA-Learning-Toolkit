package DFA_Toolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

// DFA struct which represents a DFA.
type DFA struct {
	States                []State      // Slice of states within the DFA where the index is the State ID.
	StartingStateID       int          // The ID of the starting state of the DFA.
	SymbolMap             map[rune]int // A map of each symbol within the DFA to its ID.
	depth                 int          // The depth of the DFA.
	computedDepthAndOrder bool         // Whether the depth and Order were calculated.
}

// NewDFA initializes a new empty DFA.
func NewDFA() DFA {
	return DFA{States: make([]State, 0), StartingStateID: -1,
		SymbolMap: make(map[rune]int), depth: -1, computedDepthAndOrder:false}
}

// AddState adds a new state to the DFA with the corresponding State Status.
// Returns the new state's ID (index).
func (dfa *DFA) AddState(stateStatus StateStatus) int {
	// Create empty transition table with default values of -1 for each symbol within the DFA.
	transitions := make([]int, len(dfa.SymbolMap))
	for i := range transitions {
		transitions[i] = -1
	}
	// Initialize and add the new state to the slice of states within the DFA.
	dfa.States = append(dfa.States, State{stateStatus, transitions, -1, -1, dfa})
	// Return the ID of the newly created state.
	return len(dfa.States) - 1
}

// RemoveState removes a state from DFA with the corresponding State ID.
func (dfa *DFA) RemoveState(stateID int) {
	// Panic if the state to be removed is the starting state.
	if dfa.StartingStateID == stateID {
		panic("Cannot remove starting state")
	// Panic if state ID is out of range.
	}else if stateID > len(dfa.States)-1 || stateID < 0{
		panic("stateID is out of range")
	}
	// Remove state from slice of states.
	dfa.States = append(dfa.States[:stateID], dfa.States[stateID+1:]...)
	// Update transitions to account for changed State IDs and for removed state.
	// Iterate over each state within the DFA.
	for currentStateID := range dfa.States {
		// Iterate over each symbol within the DFA.
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			// Store the ID of the resultant state.
			resultantStateID := dfa.States[currentStateID].Transitions[symbolID]
			// If the ID of the resultant state is equal to the ID of the removed state, set resultant state to -1 (undefined).
			if resultantStateID == stateID {
				dfa.States[currentStateID].Transitions[symbolID] = -1
			// Else, if the ID of the resultant state is bigger then the ID of the removed state, decrement starting state.
			} else if resultantStateID > stateID {
				dfa.States[currentStateID].Transitions[symbolID]--
			}
		}
	}
	// If the ID of the starting state is bigger then the ID of the removed state, decrement starting state.
	if dfa.StartingStateID > stateID{
		dfa.StartingStateID--
	}
}

// SymbolID returns the symbol ID for the given symbol.
func (dfa DFA) SymbolID(symbol rune) int{
	return dfa.SymbolMap[symbol]
}

// Symbol returns the symbol for the given symbol ID.
func (dfa DFA) Symbol(symbolID int) rune{
	// Iterate over each symbol within the DFA.
	for symbol := range dfa.SymbolMap {
		// If symbol ID is equal to the current symbol, return it's ID.
		if dfa.SymbolMap[symbol] == symbolID{
			return symbol
		}
	}
	// Return -1 if symbol ID is not found.
	return -1
}

// AddSymbol adds a new symbol to the DFA.
func (dfa *DFA) AddSymbol(symbol rune){
	// Panic if symbol already exists.
	if _, ok := dfa.SymbolMap[symbol]; ok {
		panic("Symbol already exists.")
	}

	// Add symbol to symbol map within the DFA.
	dfa.SymbolMap[symbol] = len(dfa.SymbolMap)
	// Iterate over each state within the DFA and add an empty (-1) transition for the newly added state.
	for stateIndex := range dfa.States {
		dfa.States[stateIndex].Transitions = append(dfa.States[stateIndex].Transitions, -1)
	}
}

// AddSymbols adds multiple new symbols to the DFA.
func (dfa *DFA) AddSymbols(symbols []rune){
	// Iteratively add each symbol within slice to the DFA.
	for _, symbol := range symbols{
		dfa.AddSymbol(symbol)
	}
}

// AddTransition adds a new transition for a given symbol from one state to another.
func (dfa *DFA) AddTransition(symbolID int, fromStateID int, toStateID int) {
	// Panic if either state IDs are out of range.
	if fromStateID > len(dfa.States)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if toStateID > len(dfa.States)-1 || toStateID < 0 {
		panic("toStateID is out of range")
	} else if symbolID > len(dfa.SymbolMap)-1 || symbolID < 0 {
		panic("symbolID is out of range")
	}
	// Add transition to fromState's transitions.
	dfa.States[fromStateID].Transitions[symbolID] = toStateID
}

// RemoveTransition removes a transition for a given symbol from one state to another.
func (dfa *DFA) RemoveTransition(symbolID int, fromStateID int) {
	// Panic if either state IDs are out of range.
	if fromStateID > len(dfa.States)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if symbolID > len(dfa.SymbolMap)-1 || symbolID < 0 {
		panic("symbolID is out of range")
	}
	// Remove transition from fromState's transitions by assigning -1 to the transitions map.
	dfa.States[fromStateID].Transitions[symbolID] = -1
}

// AllStates returns all state IDs within DFA.
func (dfa DFA) AllStates() []int {
	var allStates []int

	for stateID := range dfa.States {
		allStates = append(allStates, stateID)
	}
	return allStates
}

// AcceptingStates returns state IDs of accepting states within DFA.
func (dfa DFA) AcceptingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

// RejectingStates returns state IDs of rejecting states within DFA.
func (dfa DFA) RejectingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == REJECTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

// UnknownStates returns state IDs of unknown states within DFA.
func (dfa DFA) UnknownStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == UNKNOWN {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

// AllStatesCount returns the number of states within DFA.
func (dfa DFA) AllStatesCount() int {
	return len(dfa.States)
}

// LabelledStatesCount returns the number of labelled states (accepting or rejecting) within DFA.
func (dfa DFA) LabelledStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING || dfa.States[stateID].StateStatus == REJECTING {
			count++
		}
	}
	return count
}

// AcceptingStatesCount returns the number of accepting states within DFA.
func (dfa DFA) AcceptingStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING {
			count++
		}
	}
	return count
}

// RejectingStatesCount returns the number of rejecting states within DFA.
func (dfa DFA) RejectingStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == REJECTING {
			count++
		}
	}
	return count
}

// UnknownStatesCount returns the number of unknown states within DFA.
func (dfa DFA) UnknownStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == UNKNOWN {
			count++
		}
	}
	return count
}

// TransitionsCount returns the number of transitions within DFA.
func (dfa DFA) TransitionsCount() int{
	count := 0

	for stateIndex := range dfa.States {
		for symbol := 0; symbol < len(dfa.SymbolMap); symbol++ {
			if dfa.States[stateIndex].Transitions[symbol] != -1 {
				count++
			}
		}
	}
	return count
}

// TransitionsCountForSymbol returns the number of transitions for a given symbol within DFA.
func (dfa DFA) TransitionsCountForSymbol(symbol int) int{
	count := 0

	for stateIndex := range dfa.States {
		if dfa.States[stateIndex].Transitions[symbol] != -1 {
			count++
		}
	}
	return count
}

// SymbolsCount returns the number of symbols (alphabet) within DFA.
func (dfa DFA) SymbolsCount() int {
	return len(dfa.SymbolMap)
}

// LeavesCount returns the number of leaves within DFA.
func (dfa DFA) LeavesCount() int{
	count := 0

	for stateIndex := range dfa.States {
		transitionsCount := 0
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			if dfa.States[stateIndex].Transitions[symbolID] != -1 {
				transitionsCount++
			}
		}
		if transitionsCount == 0{
			count++
		}
	}
	return count
}

// LoopsCount returns the number of loops within DFA.
func (dfa DFA) LoopsCount() int{
	var visitedStatesCount = make(map[int]int)

	for stateID := range dfa.States {
		for symbolID := range dfa.States[stateID].Transitions {
			if dfa.States[stateID].Transitions[symbolID] != -1 {
				if dfa.States[dfa.States[stateID].Transitions[symbolID]].depth < dfa.States[stateID].depth {
					if _, ok := visitedStatesCount[stateID]; ok {
						visitedStatesCount[stateID]++
					}else{
						visitedStatesCount[stateID] = 1
					}
				}
			}
		}
	}

	return util.SumMap(visitedStatesCount, false)
}

// IsTree returns true if DFA is a tree, false is returned otherwise.
func (dfa DFA) IsTree() bool{
	var visitedStates = make(map[int]bool)

	for stateID := range dfa.States {
		for symbolID := range dfa.States[stateID].Transitions {
			if dfa.States[stateID].Transitions[symbolID] != -1 && visitedStates[dfa.States[stateID].Transitions[symbolID]]{
				return false
			}else{
				visitedStates[dfa.States[stateID].Transitions[symbolID]] = true
			}
		}
	}

	return true
}

// IsComplete returns true if DFA is complete, false is returned otherwise.
func (dfa DFA) IsComplete() bool{
	for stateID := range dfa.States {
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++{
			if dfa.States[stateID].Transitions[symbolID] < 0{
				return false
			}
		}
	}

	return true
}

// Depth returns the DFA's depth which is defined as the maximum over all nodes x of
// the length of the shortest path from the root to x.
func (dfa *DFA) Depth() int {
	// If depth and order for DFA is not computed,
	// call CalculateDepthAndOrder.
	if !dfa.computedDepthAndOrder {
		dfa.CalculateDepthAndOrder()
	}

	// Return DFA's depth.
	return dfa.depth
}

// CalculateDepthAndOrder computes the depth and Order for each state within DFA.
// This is done by traversing the DFA in a breadth-first order.
func (dfa *DFA) CalculateDepthAndOrder(){
	// Checks if DFA is valid.
	// Panics otherwise.
	dfa.IsValid()

	// Set depth of DFA to -1.
	dfa.depth = -1

	// Iterate over each state within DFA and
	// set depth and order to -1.
	for i := range dfa.States {
		dfa.States[i].depth = -1
		dfa.States[i].order = -1
	}

	// Set the depth of the starting to 0.
	dfa.StartingState().depth = 0

	// Store the current order of states.
	// Set to 0 by default.
	currentOrder := 0
	// Create a FIFO queue with starting state.
	queue := []int{dfa.StartingStateID}

	// Loop until queue is empty.
	for len(queue) > 0{
		// Remove and store first state in queue.
		stateID := queue[0]
		queue = append(queue[:0], queue[1:]...)

		// If the depth of the current state is bigger than the depth of the
		// DFA, set the depth of the DFa to the the depth of the current state.
		dfa.depth = util.Max(dfa.States[stateID].depth, dfa.depth)
		// Set the order of the current state.
		dfa.States[stateID].order = currentOrder
		// Increment current state order.
		currentOrder++

		// Iterate over each symbol (alphabet) within DFA.
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			// If transition from current state using current symbol is valid and is not a loop to the current state.
			if childStateID := dfa.States[stateID].Transitions[symbolID]; childStateID != -1 && childStateID != stateID{
				// If depth for child state has been computed, skip state.
				if dfa.States[childStateID].depth == -1{
					// Set the depth of child state to current state's depth + 1.
					dfa.States[childStateID].depth = dfa.States[stateID].depth + 1
					// Add child state to queue.
					queue = append(queue, childStateID)
				}
			}
		}
	}

	// Set the computed depth and order flag to true.
	dfa.computedDepthAndOrder = true
}

// OrderedStates returns the state IDs in order.
func (dfa DFA) OrderedStates() []int{
	// If depth and order for DFA is not computed,
	// call CalculateDepthAndOrder.
	if !dfa.computedDepthAndOrder {
		dfa.CalculateDepthAndOrder()
	}

	// Slice of ordered states using the number of states.
	orderedStates := make([]int, dfa.AllStatesCount())

	// Iterate over each state.
	for stateID := range dfa.States{
		// Use the order as the index of the ordered states slice.
		orderedStates[dfa.States[stateID].order] = stateID
	}

	// Return ordered slice of state IDs.
	return orderedStates
}

// Describe prints the details of the DFA. If detail is set to true,
// each state and each transition will also be printed.
func (dfa DFA) Describe(detail bool) {
	// Print simple description.
	fmt.Println("This DFA has", dfa.AllStatesCount(), "states and", dfa.SymbolsCount(), "alphabet (symbols)")

	// If detail is set to true, print more details.
	if detail {
		// Print alphabet/symbols mapping.
		fmt.Println("Alphabet:")
		for symbol, symbolID := range dfa.SymbolMap {
			fmt.Println(symbolID,"-",string(symbol))
		}

		// Print starting state.
		fmt.Println("Starting State:", dfa.StartingStateID)

		// Print all states.
		fmt.Println("States:")
		for k, v := range dfa.States {
			switch v.StateStatus {
			case ACCEPTING:
				fmt.Println(k, "ACCEPTING")
				break
			case REJECTING:
				fmt.Println(k, "REJECTING")
				break
			case UNKNOWN:
				fmt.Println(k, "UNKNOWN")
				break
			}
		}

		// Print all transitions.
		fmt.Println("Transitions:")
		for fromStateID, fromState := range dfa.States {
			for symbolID, toStateID := range fromState.Transitions {
				fmt.Println(fromStateID, "--", symbolID, "->", toStateID)
			}
		}
	}
}

// Accuracy returns the DFA's Accuracy with respect to a dataset.
func (dfa DFA) Accuracy(dataset Dataset) float64 {
	// Correct classifications count.
	correctClassifications := float64(0)

	// Iterate over each string instance within dataset.
	for _, stringInstance := range dataset {
		// If the status of the string instance is equal to its state status
		// within the DFA, increment correct classifications count.
		if stringInstance.Accepting == (stringInstance.ParseToStateStatus(dfa) == ACCEPTING) {
			correctClassifications++
		}
	}

	// Return the number of correct classifications divided by the length of
	// the dataset.
	return correctClassifications / float64(len(dataset))
}

// UnreachableStates returns the state IDs of unreachable states. Extracted from:
// P. Linz, An Introduction to Formal Languages and Automata. Jones & Bartlett Publishers, 2011.
func (dfa DFA) UnreachableStates() []int {
	// Map of reachable states made up of starting state.
	reachableStates := map[int]bool{dfa.StartingStateID: true}
	// Map of current states made up of starting state.
	currentStates := map[int]bool{dfa.StartingStateID: true}

	// Iterate until current states is empty.
	for len(currentStates) != 0 {
		// Map of next states.
		nextStates := map[int]bool{}
		// Iterate over current states.
		for stateID := range currentStates {
			// Iterate over each symbol within DFA.
			for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
				// If transition from current state using current symbol
				// is valid, add resultant state to next states.
				if dfa.States[stateID].Transitions[symbolID] != -1{
					nextStates[dfa.States[stateID].Transitions[symbolID]] = true
				}
			}
		}

		// Remove all state IDs from current states.
		currentStates = map[int]bool{}
		// Iterate over next states.
		for stateID := range nextStates {
			// If state is not in reachable states map, add to
			// current states and to reachable states.
			// Else, ignore since state is already reachable.
			if !reachableStates[stateID] {
				currentStates[stateID] = true
				reachableStates[stateID] = true
			}
		}
	}

	// Slice of unreachable states.
	var unReachableStates []int
	// Iterate over each state within DFA.
	for stateID := range dfa.States {
		// If state ID is not in reachable states map,
		// add to unreachable states slice.
		if !reachableStates[stateID] {
			unReachableStates = append(unReachableStates, stateID)
		}
	}

	// Return state IDs of unreachable states.
	return unReachableStates
}

// RemoveUnreachableStates removes unreachable states from DFA.
func (dfa *DFA) RemoveUnreachableStates() {
	// Get unreachable states.
	unreachableStates := dfa.UnreachableStates()

	// Iterate over unreachable states.
	for index, stateID := range unreachableStates {
		// Remove unreachable state.
		dfa.RemoveState(stateID-index)
	}
}

// StartingState returns a pointer to the starting state within the DFA.
func (dfa DFA) StartingState() *State{
	return &dfa.States[dfa.StartingStateID]
}

// Clone returns a clone of DFA.
func (dfa DFA) Clone() DFA{
	return DFA{States: dfa.States, StartingStateID: dfa.StartingStateID, SymbolMap: dfa.SymbolMap}
}

// Equal checks whether DFA is equal to the given DFA.
func (dfa DFA) Equal(dfa2 DFA) bool{
	// Minimise both DFAs.
	dfa1 := dfa.Minimise()
	dfa2 = dfa2.Minimise()

	// If the number of states or the number of accepting states
	// or the number of symbols are not equal, return false.
	if (dfa.AllStatesCount() != dfa2.AllStatesCount()) ||
		(dfa1.SymbolsCount() != dfa2.SymbolsCount() ||
			len(dfa1.AcceptingStates()) != len(dfa2.AcceptingStates())) {
		return false
	}

	// Perform breadth first search to compare DFAs.
	queue1 := []int{dfa1.StartingStateID}
	queue2 := []int{dfa2.StartingStateID}

	for len(queue1) > 0{
		stateID1 := queue1[0]
		stateID2 := queue2[0]
		queue1 = append(queue1[:0], queue1[1:]...)
		queue2 = append(queue2[:0], queue2[1:]...)

		dfa1.States[stateID1].order = 0
		dfa2.States[stateID2].order = 0

		for symbolID := 0; symbolID < len(dfa1.SymbolMap); symbolID++ {
			childStateID1 := dfa1.States[stateID1].Transitions[symbolID]
			childStateID2 := dfa2.States[stateID2].Transitions[symbolID]
			if (childStateID1 == -1 && childStateID2 != -1) ||
				(childStateID1 != -1 && childStateID2 == -1) ||
				(dfa.States[childStateID1].StateStatus != dfa2.States[childStateID2].StateStatus){
				// If a transition exists for one DFA but does not exist
				// for another DFA, return false.
				return false
			}
			if childStateID1 != -1 && childStateID1 != stateID1 {
				if dfa1.States[childStateID1].depth == -1{
					dfa1.States[childStateID1].depth = dfa1.States[stateID1].depth + 1
					dfa2.States[childStateID2].depth = dfa2.States[stateID2].depth + 1
					queue1 = append(queue1, childStateID1)
					queue2 = append(queue2, childStateID2)
				}
			}
		}
	}

	// Return true if reached.
	return true
}

// SameAs checks whether DFA is the same as the given DFA.
// This function makes use of DeepEqual and if it
// returns false, it does not necessarily mean that
// the DFAs are not equal. Use Equal() for equivalence.
func (dfa DFA) SameAs(dfa2 DFA) bool {
	return reflect.DeepEqual(dfa, dfa2)
}

// IsValid Checks whether DFA is valid.
// Panics if not valid. Used for error checking.
func (dfa DFA) IsValid() bool{
	// Check if starting state is valid.
	if dfa.StartingStateID < 0 || dfa.StartingStateID >= dfa.AllStatesCount(){
		panic("Invalid starting state.")
	// Check if number of symbols is valid.
	}else if dfa.SymbolsCount() < 1{
		panic("DFA does not contain any symbols.")
	// Check if number of states is valid.
	}else if dfa.AllStatesCount() < 1{
		panic("DFA does not contain any states.")
	// Check if any unreachable states exist within DFA.
	}else if len(dfa.UnreachableStates()) > 0{
		panic("Unreachable State exist within DFA.")
	}

	// Return true if reached.
	return true
}

// ToJSON saves the DFA to a JSON file given a file path.
func (dfa DFA) ToJSON(filePath string) bool{
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

	// Convert DFA to JSON.
	resultantJSON, err := json.MarshalIndent(dfa, "", "\t")

	// If DFA was not converted successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Copy JSON to file created.
	_, err = io.Copy(file,  bytes.NewReader(resultantJSON))

	// If JSON was not copied successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Return true if reached.
	return true
}

// DFAFromJSON returns a DFA read from a JSON file
// given a file path. The boolean value returned is set to
// true if DFA was read successfully.
func DFAFromJSON(filePath string) (DFA, bool){
	// Open file from given a path/name.
	file, err := os.Open(filePath)

	// If file was not opened successfully,
	// return empty DFA and false.
	if err != nil {
		return DFA{}, false
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize empty DFA.
	resultantDFA := DFA{}

	// Convert JSON to DFA.
	err = json.NewDecoder(file).Decode(&resultantDFA)

	// If JSON was not converted successfully,
	// return empty DFA and false.
	if err != nil {
		return DFA{}, false
	}

	// Return populated DFA and true if reached.
	return resultantDFA, true
}
