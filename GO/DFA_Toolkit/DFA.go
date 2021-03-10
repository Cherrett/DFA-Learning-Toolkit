package DFA_Toolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"reflect"
)

// DFA struct which represents a DFA.
type DFA struct {
	States          []State      // Slice of states within the DFA where the index is the State ID.
	StartingStateID int          // The ID of the starting state of the DFA.
	SymbolMap       map[rune]int // A map of each symbol within the DFA to its ID.

	Depth                 int  // The depth of the DFA.
	ComputedDepthAndOrder bool // Whether the Depth and Order were calculated DFA.
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
	dfa.States = append(dfa.States, State{stateStatus, transitions, -1, -1})
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
				if dfa.States[dfa.States[stateID].Transitions[symbolID]].Depth < dfa.States[stateID].Depth {
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
		dfa.States[i].Depth = -1
	}

	dfa.States[dfa.StartingStateID].Depth = 0

	currentOrder := 0
	queue := []int{dfa.StartingStateID}

	for len(queue) > 0{
		stateID := queue[0]
		queue = append(queue[:0], queue[1:]...)

		dfa.Depth = util.Max(dfa.States[stateID].Depth, dfa.Depth)
		dfa.States[stateID].Order = currentOrder
		currentOrder++

		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			if childStateID := dfa.States[stateID].Transitions[symbolID]; childStateID != -1 && childStateID != stateID{
				if dfa.States[childStateID].Depth == -1{
					dfa.States[childStateID].Depth = dfa.States[stateID].Depth + 1
					queue = append(queue, childStateID)
				}
			}
		}
	}
	dfa.ComputedDepthAndOrder = true
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
func (dfa DFA) Accuracy(dataset Dataset) float32 {
	// Correct classifications count.
	correctClassifications := float32(0)

	// Iterate over each string instance within dataset.
	for _, stringInstance := range dataset {
		// If the status of the string instance is equal to its state status
		// within the DFA, increment correct classifications count.
		if stringInstance.status == stringInstance.ParseToStateStatus(dfa) {
			correctClassifications++
		}
	}

	// Return the number of correct classifications divided by the length of
	// the dataset.
	return correctClassifications / float32(len(dataset))
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

		dfa.States[stateID1].Order = 0
		dfa2.States[stateID2].Order = 0

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
				if dfa.States[childStateID1].Depth == -1{
					dfa.States[childStateID1].Depth = dfa.States[stateID1].Depth + 1
					dfa2.States[childStateID2].Depth = dfa2.States[stateID2].Depth + 1
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

//func (dfa *DFA) MergeStates(state1 int, state2 int) bool{
//	state1Status := dfa.States[state1].StateStatus
//	state2Status := dfa.States[state2].StateStatus
//	var newStateStatus StateStatus = UNKNOWN
//	if (state1Status == ACCEPTING && state2Status == REJECTING) || (state1Status == REJECTING && state2Status == ACCEPTING){
//		return false
//	}else if state1Status != UNKNOWN{
//		newStateStatus = state1Status
//	}else if state2Status != UNKNOWN {
//		newStateStatus = state2Status
//	}
//
//	for i := 0; i < len(dfa.SymbolMap); i++{
//		state1transition := dfa.States[state1].Transitions[i]
//		state2transition := dfa.States[state2].Transitions[i]
//		if state1transition == -1 && state2transition != -1{
//			dfa.States[state1].Transitions[i] = state2transition
//		}else if (state2transition == -1 && state1transition != -1) || (state1transition == state1 && state2transition == state2){
//			dfa.States[state1].Transitions[i] = state1transition
//		}else if state1transition != state2transition{
//			if !dfa.MergeStates(state1transition, state2transition){
//				return false
//			}
//			dfa.States[state1].Transitions[i] = state1transition
//		}
//	}
//
//	// update new state status
//	dfa.States[state1].StateStatus = newStateStatus
//
//	// replace state2 with state1 (merged state)
//	dfa.ReplaceState(state2, state1)
//
//	return true
//}

// Replaces a state from DFA with the corresponding State ID.
//func (dfa *DFA) ReplaceState(stateID int, newStateID int) {
//	// If the state to be replaced is the starting state, replace it.
//	if dfa.StartingStateID == stateID {
//		if newStateID > stateID{
//			dfa.StartingStateID = newStateID - 1
//		}else{
//			dfa.StartingStateID = newStateID
//		}
//	// Panic if state ID is out of range.
//	}else if stateID > len(dfa.States)-1 || stateID < 0 || newStateID > len(dfa.States)-1 || newStateID < 0{
//		panic("stateID is out of range")
//	}
//	// Remove state from slice of states.
//	dfa.States = append(dfa.States[:stateID], dfa.States[stateID+1:]...)
//
//	if newStateID > stateID{
//		newStateID --
//	}
//
//	// Update transitions to account for changed State IDs and for removed state.
//	// Iterate over each state within the DFA.
//	for currentStateID := range dfa.States {
//		// Iterate over each symbol within the DFA.
//		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
//			// Store the ID of the resultant state.
//			resultantStateID := dfa.States[currentStateID].Transitions[symbolID]
//			// If the ID of the resultant state is equal to the ID of the removed state, set resultant state to new state.
//			if resultantStateID == stateID {
//				dfa.States[currentStateID].Transitions[symbolID] = newStateID
//				// Else, if the ID of the resultant state is bigger then the ID of the removed state, decrement starting state.
//			} else if resultantStateID > stateID {
//				dfa.States[currentStateID].Transitions[symbolID]--
//			}
//		}
//	}
//	// If the ID of the starting state is bigger then the ID of the removed state, decrement starting state.
//	if dfa.StartingStateID != newStateID && dfa.StartingStateID > stateID{
//		dfa.StartingStateID--
//	}
//}