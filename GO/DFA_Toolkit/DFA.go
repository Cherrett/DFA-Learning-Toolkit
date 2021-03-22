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
	StartingStateID 	  int          // The ID of the starting state of the DFA.
	SymbolMap       	  map[rune]int // A map of each symbol within the DFA to its ID.
	Depth                 int  	       // The depth of the DFA.
	ComputedDepthAndOrder bool 	       // Whether the Depth and Order were calculated.
}

// Initializes a new empty DFA.
func NewDFA() DFA {
	return DFA{States: make([]State, 0), StartingStateID: -1,
		SymbolMap: make(map[rune]int), Depth: -1, ComputedDepthAndOrder:false}
}

// Adds a new state to the DFA with the corresponding State Status.
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

// Removes a state from DFA with the corresponding State ID.
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

// Returns the symbol ID for the given symbol.
func (dfa DFA) SymbolID(symbol rune) int{
	return dfa.SymbolMap[symbol]
}

// Returns the symbol for the given symbol ID.
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

// Adds a new symbol to the DFA.
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

// Adds multiple new symbols to the DFA.
func (dfa *DFA) AddSymbols(symbols []rune){
	// Iteratively add each symbol within slice to the DFA.
	for _, symbol := range symbols{
		dfa.AddSymbol(symbol)
	}
}

// Adds a new transition for a given symbol from one state to another.
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

// Removes a transition for a given symbol from one state to another.
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

// Returns all state IDs within DFA.
func (dfa DFA) AllStates() []int {
	var allStates []int

	for stateID := range dfa.States {
		allStates = append(allStates, stateID)
	}
	return allStates
}

// Returns state IDs of accepting states within DFA.
func (dfa DFA) AcceptingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

// Returns state IDs of rejecting states within DFA.
func (dfa DFA) RejectingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == REJECTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

// Returns state IDs of unknown states within DFA.
func (dfa DFA) UnknownStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == UNKNOWN {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

// Returns the number of all states within DFA.
func (dfa DFA) AllStatesCount() int {
	return len(dfa.States)
}

// Returns the number of labelled states (accepting or rejecting) within DFA.
func (dfa DFA) LabelledStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING || dfa.States[stateID].StateStatus == REJECTING {
			count++
		}
	}
	return count
}

// Returns the number of accepting states within DFA.
func (dfa DFA) AcceptingStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING {
			count++
		}
	}
	return count
}

// Returns the number of rejecting states within DFA.
func (dfa DFA) RejectingStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == REJECTING {
			count++
		}
	}
	return count
}

// Returns the number of unknown states within DFA.
func (dfa DFA) UnknownStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == UNKNOWN {
			count++
		}
	}
	return count
}

// Returns the number of transitions within DFA.
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

// Returns the number of transitions for a given symbol within DFA.
func (dfa DFA) TransitionsCountForSymbol(symbol int) int{
	count := 0

	for stateIndex := range dfa.States {
		if dfa.States[stateIndex].Transitions[symbol] != -1 {
			count++
		}
	}
	return count
}

// Returns the number of leaves within DFA.
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

// Returns the number of loops within DFA.
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

// Returns true if DFA is a tree, false is returned otherwise.
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

// Returns true if DFA is complete, false is returned otherwise.
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

// Returns the DFA's Depth which is defined as the maximum over all nodes x of
// the length of the shortest path from the root to x.
func (dfa *DFA) GetDepth() int {
	if !dfa.ComputedDepthAndOrder {
		dfa.CalculateDepthAndOrder()
	}
	return dfa.Depth
}

func (dfa *DFA) CalculateDepthAndOrder(){
	dfa.IsValid()

	dfa.Depth = -1

	for i := range dfa.States {
		dfa.States[i].depth = -1
	}

	dfa.States[dfa.StartingStateID].depth = 0

	currentOrder := 0
	queue := []int{dfa.StartingStateID}

	for len(queue) > 0{
		stateID := queue[0]
		queue = append(queue[:0], queue[1:]...)

		dfa.Depth = util.Max(dfa.States[stateID].depth, dfa.Depth)
		dfa.States[stateID].order = currentOrder
		currentOrder++

		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			if childStateID := dfa.States[stateID].Transitions[symbolID]; childStateID != -1 && childStateID != stateID{
				if dfa.States[childStateID].depth == -1{
					dfa.States[childStateID].depth = dfa.States[stateID].depth + 1
					queue = append(queue, childStateID)
				}
			}
		}
	}
	dfa.ComputedDepthAndOrder = true
}

// Returns the state IDs in order,
func (dfa DFA) OrderedStates() []int{
	dfa.GetDepth()
	orderedStates := make([]int, len(dfa.States))

	for stateID := range dfa.States{
		orderedStates[dfa.States[stateID].order] = stateID
	}

	return orderedStates
}

func (dfa DFA) Describe(detail bool) {
	fmt.Println("This DFA has", len(dfa.States), "states and", len(dfa.SymbolMap), "alphabet")
	if detail {
		fmt.Println("Alphabet:")
		for symbol, symbolID := range dfa.SymbolMap {
			fmt.Println(symbolID,"-",string(symbol))
		}
		fmt.Println("Starting State:", dfa.StartingStateID)
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
		fmt.Println("Transitions:")
		for fromStateID, fromState := range dfa.States {
			for symbolID, toStateID := range fromState.Transitions {
				fmt.Println(fromStateID, "--", symbolID, "->", toStateID)
			}
		}
	}
}

// Returns the DFA's Accuracy with respect to a dataset.
func (dfa DFA) Accuracy(dataset Dataset) float64 {
	// Correct classifications count.
	correctClassifications := float64(0)

	// Iterate over each string instance within dataset.
	for _, stringInstance := range dataset {
		// If the status of the string instance is equal to its state status
		// within the DFA, increment correct classifications count.
		if stringInstance.Status == stringInstance.ParseToStateStatus(dfa) {
			correctClassifications++
		}
	}

	// Return the number of correct classifications divided by the length of
	// the dataset.
	return correctClassifications / float64(len(dataset))
}

func (dfa DFA) UnreachableStates() []int {
	reachableStates := map[int]bool{dfa.StartingStateID: true}
	currentStates := map[int]bool{dfa.StartingStateID: true}

	for len(currentStates) != 0 {
		nextStates := map[int]bool{}
		for stateID := range currentStates {
			for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
				if dfa.States[stateID].Transitions[symbolID] != -1{
					nextStates[dfa.States[stateID].Transitions[symbolID]] = true
				}
			}
		}
		// Don’t visit states we know to be reachable
		currentStates = map[int]bool{}
		for stateID := range nextStates {
			if !reachableStates[stateID] {
				currentStates[stateID] = true
			}
		}

		// States in Current are definitely reachable.
		for stateID := range currentStates {
			if !reachableStates[stateID] {
				reachableStates[stateID] = true
			}
		}
	}

	var unReachableStates []int
	for stateID := range dfa.States {
		if !reachableStates[stateID] {
			unReachableStates = append(unReachableStates, stateID)
		}
	}

	return unReachableStates
}

func (dfa *DFA) RemoveUnreachableStates() {
	unreachableStates := dfa.UnreachableStates()
	for index, stateID := range unreachableStates {
		dfa.RemoveState(stateID-index)
	}
}

// Returns a pointer to the starting state within the DFA.
func (dfa DFA) StartingState() *State{
	return &dfa.States[dfa.StartingStateID]
}

func (dfa DFA) Clone() DFA{
	return DFA{States: dfa.States, StartingStateID: dfa.StartingStateID, SymbolMap: dfa.SymbolMap}
}

func (dfa DFA) Equal(dfa2 DFA) bool{
	dfa = dfa.Minimise()
	dfa2 = dfa2.Minimise()

	if (len(dfa.States) != len(dfa2.States)) ||
		(len(dfa.SymbolMap) != len(dfa2.SymbolMap) ||
			len(dfa.AcceptingStates()) != len(dfa2.AcceptingStates())) {
		return false
	}
	//perform breadth first search to compare DFAs
	queue1 := []int{dfa.StartingStateID}
	queue2 := []int{dfa2.StartingStateID}

	for len(queue1) > 0{
		stateID1 := queue1[0]
		stateID2 := queue2[0]
		queue1 = append(queue1[:0], queue1[1:]...)
		queue2 = append(queue2[:0], queue2[1:]...)

		dfa.States[stateID1].order = 0
		dfa2.States[stateID2].order = 0

		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			childStateID1 := dfa.States[stateID1].Transitions[symbolID]
			childStateID2 := dfa2.States[stateID2].Transitions[symbolID]
			if (childStateID1 == -1 && childStateID2 != -1) ||
				(childStateID1 != -1 && childStateID2 == -1) ||
				(childStateID1 == childStateID1 && childStateID2 != childStateID2) ||
				(childStateID1 != childStateID1 && childStateID2 == childStateID2){
				return false
			}
			if childStateID1 != -1 && childStateID1 != stateID1 {
				if dfa.States[childStateID1].depth == -1{
					dfa.States[childStateID1].depth = dfa.States[stateID1].depth + 1
					dfa2.States[childStateID2].depth = dfa2.States[stateID2].depth + 1
					queue1 = append(queue1, childStateID1)
					queue2 = append(queue2, childStateID2)
				}
			}
		}
	}

	return true
}

func (dfa DFA) SameAs(dfa2 DFA) bool {
	return reflect.DeepEqual(dfa, dfa2)
}

func (dfa DFA) IsValid() bool{
	if dfa.StartingStateID < 0 || dfa.StartingStateID >= len(dfa.States){
		panic("Invalid starting state.")
	}else if len(dfa.SymbolMap) < 1{
		panic("DFA does not contain any symbols.")
	}else if len(dfa.States) < 1{
		panic("DFA does not contain any states.")
	}else if len(dfa.UnreachableStates()) > 0{
		panic("Unreachable State exist within DFA.")
	}
	return true
}

func (dfa DFA) ToJSON(filePath string) bool{
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()
	resultantJSON, err := json.MarshalIndent(dfa, "", "\t")
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = io.Copy(file,  bytes.NewReader(resultantJSON))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func DFAFromJSON(filePath string) (DFA, bool){
	file, err := os.Open(filePath)
	if err != nil {
		return DFA{}, false
	}
	defer file.Close()

	resultantDFA := DFA{}
	err = json.NewDecoder(file).Decode(&resultantDFA)

	if err != nil {
		return DFA{}, false
	}

	return resultantDFA, true
}
