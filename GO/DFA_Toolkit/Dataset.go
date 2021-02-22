package DFA_Toolkit

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type StringInstance struct {
	value []rune
	status StateStatus
	length uint
}

type Dataset []StringInstance

func NewStringInstanceFromAbbadingoFile(text string, delimiter string) StringInstance {
	stringInstance := StringInstance{}
	splitString := strings.Split(text, delimiter)

	switch splitString[0] {
	case "0":
		stringInstance.status = REJECTING
		break
	case "1":
		stringInstance.status = ACCEPTING
		break
	case "-1":
		stringInstance.status = UNKNOWN
		break
	default:
		panic(fmt.Sprintf("Unknown string status - %s", splitString[0]))
	}

	i, err := strconv.Atoi(splitString[1])

	if err == nil {
		stringInstance.length = uint(i)
	} else {
		panic(fmt.Sprintf("Invalid string length - %s", splitString[1]))
	}

	stringInstance.value = []rune(strings.Join(splitString[2:], ""))

	return stringInstance
}

func (stringInstance StringInstance) ConsistentWithDFA(dfa DFA) bool{
	currentState := dfa.states[dfa.startingState]
	var count uint = 0
	for _, symbol := range stringInstance.value{
		count++
		if currentState.transitions[dfa.GetSymbolID(symbol)] != -1 {
			currentState = dfa.states[currentState.transitions[dfa.GetSymbolID(symbol)]]
			// last symbol in string check
			if count == stringInstance.length {
				if stringInstance.status == ACCEPTING {
					if currentState.stateStatus == REJECTING {
						return false
					}
				}else {
					if currentState.stateStatus == ACCEPTING {
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
	currentStateID := dfa.startingState
	count := uint(0)

	for _, symbol := range stringInstance.value {
		count++

		if dfa.states[currentStateID].transitions[dfa.GetSymbolID(symbol)] != -1 {
			currentStateID = dfa.states[currentStateID].transitions[dfa.GetSymbolID(symbol)]
			// last symbol in string check
			if count == stringInstance.length{
				return dfa.states[currentStateID].stateStatus
			}
		}else{
			return UNKNOWN
		}
	}
	return UNKNOWN
}

func (stringInstance StringInstance) ParseToState(dfa DFA) (bool, int){
	currentStateID := dfa.startingState
	count := uint(0)

	for _, symbol := range stringInstance.value {
		count++

		if dfa.states[currentStateID].transitions[dfa.GetSymbolID(symbol)] != -1 {
			currentStateID = dfa.states[currentStateID].transitions[dfa.GetSymbolID(symbol)]
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

func GetDatasetFromAbbadingoFile(fileName string) Dataset {
	dataset := Dataset{}

	file, err := os.Open(fileName)

	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan() // ignore first line
		for scanner.Scan() {
			dataset = append(dataset, NewStringInstanceFromAbbadingoFile(scanner.Text(), " "))
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	} else {
		panic("Invalid file name")
	}
	return dataset
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

func (dataset Dataset) GetAcceptingStringInstances() Dataset{
	var acceptingInstances Dataset

	for _, stringInstance := range dataset {
		if stringInstance.status == ACCEPTING {
			acceptingInstances = append(acceptingInstances, stringInstance)
		}
	}

	return acceptingInstances
}

func (dataset Dataset) GetRejectingStringInstances() Dataset{
	var rejectingInstances Dataset

	for _, stringInstance := range dataset {
		if stringInstance.status == REJECTING {
			rejectingInstances = append(rejectingInstances, stringInstance)
		}
	}

	return rejectingInstances
}
