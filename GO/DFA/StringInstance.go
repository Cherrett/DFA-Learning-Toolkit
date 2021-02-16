package DFA

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
	stringValue  string
	stringStatus StateStatus
	length       uint
}

func NewStringInstance(text string, delimiter string) *StringInstance {
	stringInstance := StringInstance{}
	splitString := strings.Split(text, delimiter)

	switch splitString[0] {
	case "0":
		stringInstance.stringStatus = REJECTING
		break
	case "1":
		stringInstance.stringStatus = ACCEPTING
		break
	case "-1":
		stringInstance.stringStatus = UNKNOWN
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

	stringInstance.stringValue = strings.Join(splitString[2:], "")

	return &stringInstance
}

func GetListOfStringInstancesFromFile(fileName string) []StringInstance {
	var listOfStrings []StringInstance

	file, err := os.Open(fileName)

	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan() // ignore first line
		for scanner.Scan() {
			listOfStrings = append(listOfStrings, *NewStringInstance(scanner.Text(), " "))
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	} else {
		panic("Invalid file name")
	}
	return listOfStrings
}

func SortListOfStringInstances(strings []StringInstance) []StringInstance {
	sort.Slice(strings[:], func(i, j int) bool {
		return strings[i].length < strings[j].length
	})
	return strings
}

func StringInstanceConsistentWithDFA(stringInstance StringInstance, dfa DFA) bool{
	// Skip unknown strings (test data)
	if stringInstance.stringStatus == UNKNOWN{
		return true
	}

	exists := false
	currentState := dfa.startingState
	var count uint = 0
	for _, character := range stringInstance.stringValue{
		count++
		exists = false
		if value, ok := currentState.transitions[character]; ok {
			currentState = dfa.states[value]
			exists = true
		}

		if !exists {
			return false
		}

		// last symbol in string check
		if count == stringInstance.length {
			if stringInstance.stringStatus == ACCEPTING {
				if currentState.stateStatus == REJECTING {
					return false
				}
			}else {
				if currentState.stateStatus == ACCEPTING {
					return false
				}
			}
		}
	}
	return true
}

func ListOfStringInstancesConsistentWithDFA(stringInstances []StringInstance, dfa DFA) bool{
	consistent := true
	var wg sync.WaitGroup
	wg.Add(len(stringInstances))

	for _, stringInstance := range stringInstances {
		go func(stringInstance StringInstance, dfa DFA){
			defer wg.Done()
			if consistent {
				consistent = StringInstanceConsistentWithDFA(stringInstance, dfa)
			}
		}(stringInstance, dfa)
	}

	wg.Wait()
	return consistent
}