package DFA_Toolkit

import (
	"sort"
	"strconv"
	"sync"
)

type StringInstance struct {
	value []rune
	status StateStatus
	length uint
}

type Dataset []StringInstance

func (dataset Dataset) GetPTA(APTA bool) DFA {
	dataset = dataset.SortDatasetByLength()
	alphabet := make(map[rune]bool)
	var count uint
	var currentStateID, newStateID int
	dfa := NewDFA()

	if dataset[0].length == 0 {
		if dataset[0].status == ACCEPTING {
			currentStateID = dfa.AddState(ACCEPTING)
		} else if APTA{
			currentStateID = dfa.AddState(REJECTING)
		}else{
			currentStateID = dfa.AddState(UNKNOWN)
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

func (stringInstance StringInstance) ConsistentWithDFA(dfa DFA) bool{
	currentState := dfa.States[dfa.StartingStateID]
	var count uint = 0
	for _, symbol := range stringInstance.value{
		count++
		if currentState.Transitions[dfa.SymbolID(symbol)] != -1 {
			currentState = dfa.States[currentState.Transitions[dfa.SymbolID(symbol)]]
			// last symbol in string check
			if count == stringInstance.length {
				if stringInstance.status == ACCEPTING {
					if currentState.StateStatus == REJECTING {
						return false
					}
				}else {
					if currentState.StateStatus == ACCEPTING {
						return false
					}
				}
			}
		}else{
			return !(stringInstance.status == ACCEPTING)
		}
	}
	return true
}

func (stringInstance StringInstance) ParseToStateStatus(dfa DFA) StateStatus{
	currentStateID := dfa.StartingStateID
	count := uint(0)

	for _, symbol := range stringInstance.value {
		count++

		if dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)] != -1 {
			currentStateID = dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)]
			// last symbol in string check
			if count == stringInstance.length{
				if dfa.States[currentStateID].StateStatus == UNKNOWN{
					return REJECTING
				}else{
					return ACCEPTING
				}
			}
		}else{
			return REJECTING
		}
	}
	return REJECTING
}

func (stringInstance StringInstance) ParseToState(dfa DFA) (bool, int){
	currentStateID := dfa.StartingStateID
	count := uint(0)

	for _, symbol := range stringInstance.value {
		count++

		if dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)] != -1 {
			currentStateID = dfa.States[currentStateID].Transitions[dfa.SymbolID(symbol)]
			// last symbol in string check
			if count == stringInstance.length{
				return true, currentStateID
			}
		}else{
			return false, -1
		}
	}
	return false, -1
}

func BinaryStringToStringInstance(dfa DFA, binaryString string) StringInstance{
	stringInstance := StringInstance{length: uint(len(binaryString))}

	for _, value := range binaryString{
		symbolID, err := strconv.Atoi(string(value))
		if err != nil || (symbolID != 0 && symbolID != 1){
			panic("Not a binary string")
		}
		stringInstance.value = append(stringInstance.value, dfa.Symbol(symbolID))
	}

	stringInstance.status = stringInstance.ParseToStateStatus(dfa)

	return stringInstance
}

func (dataset Dataset) SortDatasetByLength() Dataset{
	// sort all string instances by length
	sort.Slice(dataset[:], func(i, j int) bool {
		return dataset[i].length < dataset[j].length
	})
	return dataset
}

func (dataset Dataset) ConsistentWithDFA(dfa DFA) bool{
	consistent := true
	var wg sync.WaitGroup
	wg.Add(len(dataset))

	for _, stringInstance := range dataset {
		go func(stringInstance StringInstance, dfa DFA){
			defer wg.Done()
			if consistent {
				consistent = stringInstance.ConsistentWithDFA(dfa)
			}
		}(stringInstance, dfa)
	}

	wg.Wait()
	return consistent
}

func (dataset Dataset) AcceptingStringInstances() Dataset{
	var acceptingInstances Dataset

	for _, stringInstance := range dataset {
		if stringInstance.status == ACCEPTING {
			acceptingInstances = append(acceptingInstances, stringInstance)
		}
	}

	return acceptingInstances
}

func (dataset Dataset) RejectingStringInstances() Dataset{
	var rejectingInstances Dataset

	for _, stringInstance := range dataset {
		if stringInstance.status == REJECTING {
			rejectingInstances = append(rejectingInstances, stringInstance)
		}
	}

	return rejectingInstances
}

func (dataset Dataset) AcceptingStringInstancesCount() int{
	count := 0

	for _, stringInstance := range dataset {
		if stringInstance.status == ACCEPTING {
			count++
		}
	}

	return count
}

func (dataset Dataset) RejectingStringInstancesCount() int{
	count := 0

	for _, stringInstance := range dataset {
		if stringInstance.status == REJECTING {
			count++
		}
	}

	return count
}

func (dataset Dataset) AcceptingStringInstancesRatio() float64{
	return float64(dataset.AcceptingStringInstancesCount()) / float64(len(dataset))
}

func (dataset Dataset) RejectingStringInstancesRatio() float64{
	return float64(dataset.RejectingStringInstancesCount()) / float64(len(dataset))
}