package DFA_Toolkit

import (
	"fmt"
	"reflect"
)

type DFA struct {
	states        []State
	startingState int
	symbolMap     map[rune]int
}

func NewDFA() DFA {
	return DFA{states: make([]State, 0), startingState: -1, symbolMap: make(map[rune]int)}
}

func (dfa *DFA) AddState(stateStatus StateStatus) int {
	transitions := make([]int, len(dfa.symbolMap))
	for i := range transitions {
		transitions[i] = -1
	}
	dfa.states = append(dfa.states, State{stateStatus, transitions})
	return len(dfa.states) - 1
}

func (dfa *DFA) RemoveState(stateID int) {
	// panic if stateID to be removed is the starting state
	if dfa.startingState == stateID {
		panic("Cannot remove starting state")
	}
	// remove state from slice of states
	dfa.states = append(dfa.states[:stateID], dfa.states[stateID+1:]...)
	// update transitions to account for new stateIDs and for removed state
	for stateIndex := range dfa.states {
		for symbolID := 0; symbolID < len(dfa.symbolMap); symbolID++ {
			if dfa.states[stateIndex].transitions[symbolID] == stateID {
				dfa.states[stateIndex].transitions[symbolID] = -1
			} else if dfa.states[stateIndex].transitions[symbolID] > stateID {
				dfa.states[stateIndex].transitions[symbolID]--
			}
		}
	}
	// update starting state
	if dfa.startingState > stateID{
		dfa.startingState--
	}
}

func (dfa DFA) SymbolID(symbol rune) int{
	return dfa.symbolMap[symbol]
}

func (dfa DFA) Symbol(symbolID int) rune{
	for symbol := range dfa.symbolMap{
		if dfa.symbolMap[symbol] == symbolID{
			return symbol
		}
	}
	return -1
}

func (dfa *DFA) AddSymbol(symbol rune){
	dfa.symbolMap[symbol] = len(dfa.symbolMap)
	for stateIndex := range dfa.states {
		dfa.states[stateIndex].transitions = append(dfa.states[stateIndex].transitions, -1)
	}
}

func (dfa *DFA) AddSymbols(symbols []rune){
	for _, symbol := range symbols{
		dfa.AddSymbol(symbol)
	}
}

func (dfa *DFA) AddTransition(symbolID int, fromStateID int, toStateID int) {
	// error checking
	if fromStateID > len(dfa.states)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if toStateID > len(dfa.states)-1 || toStateID < 0 {
		panic("toStateID is out of range")
	} else if symbolID > len(dfa.symbolMap)-1 || symbolID < 0 {
		panic("symbolID is out of range")
	}
	// add transition to fromState's transitions
	dfa.states[fromStateID].transitions[symbolID] = toStateID
}

func (dfa *DFA) RemoveTransition(symbolID int, fromStateID int) {
	// error checking
	if fromStateID > len(dfa.states)-1 || fromStateID < 0 {
		panic("fromStateID is out of range")
	} else if symbolID > len(dfa.symbolMap)-1 || symbolID < 0 {
		panic("symbolID is out of range")
	}
	// remove transition to fromState's transitions
	dfa.states[fromStateID].transitions[symbolID] = -1
}

func (dfa DFA) AllStates() []int {
	var allStates []int

	for stateID := range dfa.states {
		allStates = append(allStates, stateID)
	}
	return allStates
}

func (dfa DFA) AcceptingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == ACCEPTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) RejectingStates() []int {
	var acceptingStates []int

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == REJECTING {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) UnknownStates() []int {
	var acceptingStates []int

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == UNKNOWN {
			acceptingStates = append(acceptingStates, stateID)
		}
	}
	return acceptingStates
}

func (dfa DFA) AllStatesCount() int {
	return len(dfa.states)
}

func (dfa DFA) AcceptingStatesCount() int {
	count := 0

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == ACCEPTING {
			count++
		}
	}
	return count
}

func (dfa DFA) RejectingStatesCount() int {
	count := 0

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == REJECTING {
			count++
		}
	}
	return count
}

func (dfa DFA) UnknownStatesCount() int {
	count := 0

	for stateID := range dfa.states {
		if dfa.states[stateID].stateStatus == UNKNOWN {
			count++
		}
	}
	return count
}

func (dfa DFA) TransitionsCount() int{
	count := 0

	for stateIndex := range dfa.states {
		for symbol := 0; symbol < len(dfa.symbolMap); symbol++ {
			if dfa.states[stateIndex].transitions[symbol] != -1 {
				count++
			}
		}
	}
	return count
}

func (dfa DFA) TransitionsCountForSymbol(symbol int) int{
	count := 0

	for stateIndex := range dfa.states {
		if dfa.states[stateIndex].transitions[symbol] != -1 {
			count++
		}
	}
	return count
}

func (dfa DFA) LeavesCount() int{
	count := 0

	for stateIndex := range dfa.states {
		transitionsCount := 0
		for symbolID := 0; symbolID < len(dfa.symbolMap); symbolID++ {
			if dfa.states[stateIndex].transitions[symbolID] != -1 {
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

	for stateID := range dfa.states{
		for symbolID := range dfa.states[stateID].transitions{
			if dfa.states[stateID].transitions[symbolID] != -1 {
				if _, ok := visitedStatesCount[stateID]; ok {
					visitedStatesCount[stateID]++
				}else{
					visitedStatesCount[stateID] = 1
				}
			}
		}
	}

	count := 0

	for visitedStateCount := range visitedStatesCount{
		if visitedStateCount > 1{
			count += visitedStateCount - 1
		}
	}

	return count
}

func (dfa DFA) IsTree() bool{
	var visitedStates = make(map[int]bool)

	for stateID := range dfa.states{
		for symbolID := range dfa.states[stateID].transitions{
			if dfa.states[stateID].transitions[symbolID] != -1 && visitedStates[dfa.states[stateID].transitions[symbolID]]{
				return false
			}else{
				visitedStates[dfa.states[stateID].transitions[symbolID]] = true
			}
		}
	}

	return true
}

// Returns the DFA's Depth which is defined as the maximum over all nodes x of
// the length of the shortest path from the root to x
func (dfa DFA) Depth() uint {
	// if the DFA is a tree, calculate the depth using a recursive function
	// which calculates the height of each node by assigning the maximum height
	// of its subtrees
	if dfa.IsTree(){
		return uint(dfa.DepthUtilTree(dfa.startingState))
		// if the DFA is not a tree, calculate the depth using Breadth First Traversal
		// while keeping track of the traversed nodes to avoid revisiting the same nodes
		// more than once in case of loops
	}else{
		var maxValue uint
		var stateMap = make(map[int]uint)

		for stateID := range dfa.states{
			stateMap[stateID] = dfa.DepthUtilNonTree(stateID)
		}

		for _, v := range stateMap {
			if v > maxValue {
				maxValue = v
			}
		}
		return maxValue
	}
}

// A recursive function which calculates the height of each node by assigning
// the maximum height of its subtrees. Used only if the DFA is a tree.
func (dfa DFA) DepthUtilTree(stateID int) int{
	if stateID == -1{
		return -1
	}
	var depths []int
	maxValue := -1

	for _, symbolID := range dfa.symbolMap{
		depths = append(depths, dfa.DepthUtilTree(dfa.states[stateID].transitions[symbolID]))
	}

	for _, depth := range depths {
		if depth > maxValue {
			maxValue = depth
		}
	}

	return maxValue + 1
}

// A function which calculates the DFA's depth using Breadth First Traversal
// while keeping track of the traversed nodes to avoid revisiting the same nodes
// more than once in case of loops
func (dfa DFA) DepthUtilNonTree(targetStateID int) uint{
	visitedStates := map[int]bool{dfa.startingState: true}
	queue := [][]int{{dfa.startingState, 0}}

	for len(queue) > 0{
		stateID := queue[0][0]
		depth := queue[0][1]
		queue = append(queue[:0], queue[1:]...)

		if stateID == targetStateID{
			return uint(depth)
		}

		for symbolID := range dfa.states[stateID].transitions {
			if dfa.states[stateID].transitions[symbolID] != -1 && dfa.states[stateID].transitions[symbolID] != stateID {
				if !visitedStates[dfa.states[stateID].transitions[symbolID]] {
					queue = append(queue, []int{dfa.states[stateID].transitions[symbolID], depth +1})
					visitedStates[dfa.states[stateID].transitions[symbolID]] = true
				}
			}
		}
	}

	return 0
}

func (dfa DFA) Describe(detail bool) {
	fmt.Println("This DFA has", len(dfa.states), "states and", len(dfa.symbolMap), "alphabet")
	if detail {
		fmt.Println("Alphabet:")
		for symbol, symbolID := range dfa.symbolMap{
			fmt.Println(symbolID,"-",string(symbol))
		}
		fmt.Println("Starting State:", dfa.startingState)
		fmt.Println("States:")
		for k, v := range dfa.states {
			switch v.stateStatus {
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
		for fromStateID, fromState := range dfa.states {
			for symbolID, toStateID := range fromState.transitions {
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

	dfa.startingState = currentStateID

	for _, stringInstance := range dataset {
		if !APTA && stringInstance.status != ACCEPTING {
			continue
		}
		currentStateID = dfa.startingState
		count = 0
		for _, symbol := range stringInstance.value {
			count++
			// new alphabet check
			if !alphabet[symbol] {
				dfa.AddSymbol(symbol)
				alphabet[symbol] = true
			}

			symbolID := dfa.SymbolID(symbol)

			if dfa.states[currentStateID].transitions[symbolID] != -1 {
				currentStateID = dfa.states[currentStateID].transitions[symbolID]
				// last symbol in string check
				if count == stringInstance.length {
					if stringInstance.status == ACCEPTING {
						if dfa.states[currentStateID].stateStatus == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						} else {
							dfa.states[currentStateID].UpdateStateStatus(ACCEPTING)
						}
					} else {
						if dfa.states[currentStateID].stateStatus == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							dfa.states[currentStateID].UpdateStateStatus(REJECTING)
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
				dfa.states[currentStateID].transitions[symbolID] = newStateID
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
	reachableStates := map[int]bool{dfa.startingState: true}
	currentStates := map[int]bool{dfa.startingState: true}

	for len(currentStates) != 0 {
		nextStates := map[int]bool{}
		for stateID := range currentStates {
			for symbolID := 0; symbolID < len(dfa.symbolMap); symbolID++ {
				if dfa.states[stateID].transitions[symbolID] != -1{
					nextStates[dfa.states[stateID].transitions[symbolID]] = true
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
	for stateID := range dfa.states {
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
	for stateID := range dfa.states {
		for stateID2 := stateID + 1; stateID2 < len(dfa.states); stateID2++ {
			allPairs[StateIDPair{stateID, stateID2}] = true
			if dfa.states[stateID].stateStatus != dfa.states[stateID2].stateStatus{
				distinguishablePairs[StateIDPair{stateID, stateID2}] = true
			}
		}
	}

	oldCount := 0
	for oldCount != len(distinguishablePairs) {
		oldCount = len(distinguishablePairs)

		for stateID := range dfa.states {
			for stateID2 := stateID + 1; stateID2 < len(dfa.states); stateID2++ {
				if distinguishablePairs[StateIDPair{stateID, stateID2}] {
					continue
				} else {
					for symbolID := 0; symbolID < len(dfa.symbolMap); symbolID++ {
						if dfa.states[stateID].transitions[symbolID] != -1 {
							if dfa.states[stateID2].transitions[symbolID] != -1 {
								if distinguishablePairs[StateIDPair{dfa.states[stateID].transitions[symbolID], dfa.states[stateID2].transitions[symbolID]}] ||
									distinguishablePairs[StateIDPair{dfa.states[stateID2].transitions[symbolID], dfa.states[stateID].transitions[symbolID]}] {
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
	for stateID := range dfa.states{
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
	resultantDFA := DFA{symbolMap: dfa.symbolMap}

	// Create a new state for each block
	for blockIndex := range currentPartition{
		var stateStatus StateStatus = UNKNOWN
		for stateID := range currentPartition[blockIndex] {
			stateStatus = dfa.states[stateID].stateStatus
			break
		}
		resultantDFA.AddState(stateStatus)
	}

	// Initial State
	for blockIndex := range currentPartition{
		found := false
		for stateID := range currentPartition[blockIndex] {
			if stateID == dfa.startingState{
				resultantDFA.startingState = blockIndex
				found = true
				break
			}
		}
		if found{
			break
		}
	}

	// Transitions
	for stateID := range dfa.states{
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

		for symbolID := range dfa.states[stateID].transitions{
			resultantStateBlockIndex := 0
			if resultantDFA.states[stateBlockIndex].transitions[symbolID] == -1 && dfa.states[stateID].transitions[symbolID] != -1{
				found := false
				for blockIndex := range currentPartition{
					for stateID2 := range currentPartition[blockIndex] {
						if dfa.states[stateID].transitions[symbolID] == stateID2{
							resultantStateBlockIndex = blockIndex
							found = true
							break
						}
					}
					if found{
						break
					}
				}
				resultantDFA.states[stateBlockIndex].transitions[symbolID] = resultantStateBlockIndex
			}
		}
	}

	return resultantDFA
}

func (dfa DFA) Clone() DFA{
	return DFA{states: dfa.states, startingState: dfa.startingState, symbolMap: dfa.symbolMap}
}

func (dfa DFA) Equal(dfa2 DFA) bool{
	// need to confirm
	for stateID := range dfa.states{
		for symbolID := range dfa.states[stateID].transitions{
			if dfa.states[stateID].transitions[symbolID] != dfa2.states[stateID].transitions[symbolID]{
				return false
			}
		}
	}

	return true
}

func (dfa DFA) SameAs(dfa2 DFA) bool {
	return reflect.DeepEqual(dfa, dfa2)
}