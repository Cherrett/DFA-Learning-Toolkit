package DFA_Toolkit

import (
	"DFA_Toolkit/DFA_Toolkit/util"
	"fmt"
	"reflect"
)

type DFA struct {
	States          []State      // List of states within the DFA.
	StartingStateID int          // The ID of the starting state of the DFA.
	SymbolMap       map[rune]int // A map of each symbol within the DFA to its ID.

	Depth                 int  // The depth of the DFA.
	ComputedDepthAndOrder bool // Whether the Depth and Order were calculated DFA
}

func NewDFA() DFA {
	return DFA{States: make([]State, 0), StartingStateID: -1,
		SymbolMap: make(map[rune]int), Depth: -1, ComputedDepthAndOrder:false}
}

func (dfa *DFA) AddState(stateStatus StateStatus) int {
	transitions := make([]int, len(dfa.SymbolMap))
	for i := range transitions {
		transitions[i] = -1
	}
	dfa.States = append(dfa.States, State{stateStatus, transitions, -1, -1})
	return len(dfa.States) - 1
}

func (dfa *DFA) RemoveState(stateID int) {
	// panic if stateID to be removed is the starting state
	if dfa.StartingStateID == stateID {
		panic("Cannot remove starting state")
	}
	// remove state from slice of states
	dfa.States = append(dfa.States[:stateID], dfa.States[stateID+1:]...)
	// update transitions to account for new stateIDs and for removed state
	for stateIndex := range dfa.States {
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			if dfa.States[stateIndex].Transitions[symbolID] == stateID {
				dfa.States[stateIndex].Transitions[symbolID] = -1
			} else if dfa.States[stateIndex].Transitions[symbolID] > stateID {
				dfa.States[stateIndex].Transitions[symbolID]--
			}
		}
	}
	// update starting state
	if dfa.StartingStateID > stateID{
		dfa.StartingStateID--
	}
}

func (dfa DFA) SymbolID(symbol rune) int{
	return dfa.SymbolMap[symbol]
}

func (dfa DFA) Symbol(symbolID int) rune{
	for symbol := range dfa.SymbolMap {
		if dfa.SymbolMap[symbol] == symbolID{
			return symbol
		}
	}
	return -1
}

func (dfa *DFA) AddSymbol(symbol rune){
	dfa.SymbolMap[symbol] = len(dfa.SymbolMap)
	for stateIndex := range dfa.States {
		dfa.States[stateIndex].Transitions = append(dfa.States[stateIndex].Transitions, -1)
	}
}

func (dfa *DFA) AddSymbols(symbols []rune){
	for _, symbol := range symbols{
		dfa.AddSymbol(symbol)
	}
}

func (dfa *DFA) AddTransition(symbolID int, fromStateID int, toStateID int) {
	// error checking
	if fromStateID > len(dfa.States)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if toStateID > len(dfa.States)-1 || toStateID < 0 {
		panic("toStateID is out of range")
	} else if symbolID > len(dfa.SymbolMap)-1 || symbolID < 0 {
		panic("symbolID is out of range")
	}
	// add transition to fromState's transitions
	dfa.States[fromStateID].Transitions[symbolID] = toStateID
}

func (dfa *DFA) RemoveTransition(symbolID int, fromStateID int) {
	// error checking
	if fromStateID > len(dfa.States)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if symbolID > len(dfa.SymbolMap)-1 || symbolID < 0 {
		panic("symbolID is out of range")
	}
	// remove transition to fromState's transitions
	dfa.States[fromStateID].Transitions[symbolID] = -1
}

func (dfa DFA) AllStates() []int {
	var allStates []int

	for stateID := range dfa.States {
		allStates = append(allStates, stateID)
	}
	return allStates
}

func (dfa DFA) AcceptingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) RejectingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == REJECTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) UnknownStates() []int {
	var acceptingStates []int

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == UNKNOWN {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) AllStatesCount() int {
	return len(dfa.States)
}

func (dfa DFA) AcceptingStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == ACCEPTING {
			count++
		}
	}
	return count
}

func (dfa DFA) RejectingStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == REJECTING {
			count++
		}
	}
	return count
}

func (dfa DFA) UnknownStatesCount() int {
	count := 0

	for stateID := range dfa.States {
		if dfa.States[stateID].StateStatus == UNKNOWN {
			count++
		}
	}
	return count
}

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

func (dfa DFA) TransitionsCountForSymbol(symbol int) int{
	count := 0

	for stateIndex := range dfa.States {
		if dfa.States[stateIndex].Transitions[symbol] != -1 {
			count++
		}
	}
	return count
}

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
// the length of the shortest path from the root to x
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

func GetPTAFromDataset(dataset Dataset, APTA bool) DFA {
	dataset = dataset.SortDatasetByLength()
	alphabet := make(map[rune]bool)
	var count uint
	var currentStateID, newStateID int
	dfa := NewDFA()

	if dataset[0].length == 0 {
		if dataset[0].status == ACCEPTING {
			currentStateID = dfa.AddState(ACCEPTING)
		} else {
			currentStateID = dfa.AddState(REJECTING)
		}
	} else {
		currentStateID = dfa.AddState(UNKNOWN)
	}

	dfa.StartingStateID = currentStateID

	for _, stringInstance := range dataset {
		if !APTA && stringInstance.status != ACCEPTING {
			continue
		}
		currentStateID = dfa.StartingStateID
		count = 0
		for _, symbol := range stringInstance.value {
			count++
			// new alphabet check
			if !alphabet[symbol] {
				dfa.AddSymbol(symbol)
				alphabet[symbol] = true
			}

			symbolID := dfa.SymbolID(symbol)

			if dfa.States[currentStateID].Transitions[symbolID] != -1 {
				currentStateID = dfa.States[currentStateID].Transitions[symbolID]
				// last symbol in string check
				if count == stringInstance.length {
					if stringInstance.status == ACCEPTING {
						if dfa.States[currentStateID].StateStatus == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						} else {
							dfa.States[currentStateID].UpdateStateStatus(ACCEPTING)
						}
					} else {
						if dfa.States[currentStateID].StateStatus == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							dfa.States[currentStateID].UpdateStateStatus(REJECTING)
						}
					}
				}
			} else {
				// last symbol in string check
				if count == stringInstance.length {
					if stringInstance.status == ACCEPTING {
						newStateID = dfa.AddState(ACCEPTING)
					} else {
						newStateID = dfa.AddState(REJECTING)
					}
				} else {
					newStateID = dfa.AddState(UNKNOWN)
				}
				dfa.States[currentStateID].Transitions[symbolID] = newStateID
				currentStateID = newStateID
			}
		}
	}
	return dfa
}

func (dfa DFA) AccuracyOfDFA(dataset Dataset) float32 {
	correctClassifications := float32(0)

	for _, stringInstance := range dataset {
		if stringInstance.status == stringInstance.ParseToStateStatus(dfa) {
			correctClassifications++
		}
	}
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

type StateIDPair struct {
	state1 int
	state2 int
}

func (dfa *DFA) Mark() (map[StateIDPair]bool, map[StateIDPair]bool) {
	distinguishablePairs := map[StateIDPair]bool{}
	indistinguishablePairs := map[StateIDPair]bool{}
	allPairs := map[StateIDPair]bool{}

	dfa.RemoveUnreachableStates()
	for stateID := range dfa.States {
		for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
			allPairs[StateIDPair{stateID, stateID2}] = true
			if dfa.States[stateID].StateStatus != dfa.States[stateID2].StateStatus {
				distinguishablePairs[StateIDPair{stateID, stateID2}] = true
			}
		}
	}

	oldCount := 0
	for oldCount != len(distinguishablePairs) {
		oldCount = len(distinguishablePairs)

		for stateID := range dfa.States {
			for stateID2 := stateID + 1; stateID2 < len(dfa.States); stateID2++ {
				if distinguishablePairs[StateIDPair{stateID, stateID2}] {
					continue
				} else {
					for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
						if dfa.States[stateID].Transitions[symbolID] != -1 {
							if dfa.States[stateID2].Transitions[symbolID] != -1 {
								if distinguishablePairs[StateIDPair{dfa.States[stateID].Transitions[symbolID], dfa.States[stateID2].Transitions[symbolID]}] ||
									distinguishablePairs[StateIDPair{dfa.States[stateID2].Transitions[symbolID], dfa.States[stateID].Transitions[symbolID]}] {
									distinguishablePairs[StateIDPair{stateID, stateID2}] = true
								}
							}
						}
					}
				}
			}
		}
	}

	var distinguishablePairsList [][]int
	for stateIDPair := range distinguishablePairs {
		distinguishablePairsList = append(distinguishablePairsList, []int{stateIDPair.state1, stateIDPair.state2})
	}

	for pair := range allPairs{
		if !distinguishablePairs[StateIDPair{pair.state1, pair.state2}]{
			indistinguishablePairs[StateIDPair{pair.state1, pair.state2}] = true
		}
	}

	return distinguishablePairs, indistinguishablePairs
}

func (dfa DFA) Minimise() DFA{
	// get distinguishable and indistinguishable state pairs using Mark function
	_, indistinguishablePairs := dfa.Mark()

	// Partition states into blocks
	var currentPartition []map[int]bool
	for stateID := range dfa.States {
		exists := false
		for _, block := range currentPartition{
			if block[stateID]{
				exists = true
			}
		}
		if !exists{
			indistinguishable := false
			for indistinguishablePair := range indistinguishablePairs{
				if indistinguishablePair.state1 == stateID{
					for blockIndex, block := range currentPartition{
						if block[indistinguishablePair.state2]{
							currentPartition[blockIndex][stateID] = true
							indistinguishable = true
							break
						}
					}
					if indistinguishable{
						break
					}
				}else if indistinguishablePair.state2 == stateID{
					for blockIndex, block := range currentPartition{
						if block[indistinguishablePair.state1]{
							currentPartition[blockIndex][stateID] = true
							indistinguishable = true
							break
						}
					}
					if indistinguishable{
						break
					}
				}
			}
			if !indistinguishable{
				currentPartition = append(currentPartition, map[int]bool{stateID: true})
			}
		}
	}
	resultantDFA := DFA{SymbolMap: dfa.SymbolMap}

	// Create a new state for each block
	for blockIndex := range currentPartition{
		var stateStatus StateStatus = UNKNOWN
		for stateID := range currentPartition[blockIndex] {
			stateStatus = dfa.States[stateID].StateStatus
			break
		}
		resultantDFA.AddState(stateStatus)
	}

	// Initial State
	for blockIndex := range currentPartition{
		found := false
		for stateID := range currentPartition[blockIndex] {
			if stateID == dfa.StartingStateID {
				resultantDFA.StartingStateID = blockIndex
				found = true
				break
			}
		}
		if found{
			break
		}
	}

	// Transitions
	for stateID := range dfa.States {
		stateBlockIndex := 0
		found := false
		for blockIndex := range currentPartition{
			for stateID2 := range currentPartition[blockIndex] {
				if stateID == stateID2{
					stateBlockIndex = blockIndex
					found = true
					break
				}
			}
			if found{
				break
			}
		}

		for symbolID := range dfa.States[stateID].Transitions {
			resultantStateBlockIndex := 0
			if resultantDFA.States[stateBlockIndex].Transitions[symbolID] == -1 && dfa.States[stateID].Transitions[symbolID] != -1{
				found := false
				for blockIndex := range currentPartition{
					for stateID2 := range currentPartition[blockIndex] {
						if dfa.States[stateID].Transitions[symbolID] == stateID2{
							resultantStateBlockIndex = blockIndex
							found = true
							break
						}
					}
					if found{
						break
					}
				}
				resultantDFA.States[stateBlockIndex].Transitions[symbolID] = resultantStateBlockIndex
			}
		}
	}

	return resultantDFA
}

func (dfa DFA) Clone() DFA{
	return DFA{States: dfa.States, StartingStateID: dfa.StartingStateID, SymbolMap: dfa.SymbolMap}
}

func (dfa DFA) Equal(dfa2 DFA) bool{
	// need to confirm
	for stateID := range dfa.States {
		for symbolID := range dfa.States[stateID].Transitions {
			if dfa.States[stateID].Transitions[symbolID] != dfa2.States[stateID].Transitions[symbolID]{
				return false
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