package DFA

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type StateStatus uint8

const (
	REJECTING = iota // 0
	ACCEPTING        // 1
	UNKNOWN          // 2
)

type State struct {
	stateStatus StateStatus
	stateID     uint
	transitions map[int32]uint
}

type DFA struct {
	states        map[uint]State
	startingState State
	alphabet      map[int32]bool
}

func (dfa DFA) GetAcceptingStates() []State {
	var acceptingStates []State

	for _, v := range dfa.states {
		if v.stateStatus == ACCEPTING {
			acceptingStates = append(acceptingStates, v)
		}
	}

	return acceptingStates
}

func (dfa DFA) GetRejectingStates() []State {
	var rejectingStates []State

	for _, v := range dfa.states {
		if v.stateStatus == REJECTING {
			rejectingStates = append(rejectingStates, v)
		}
	}

	return rejectingStates
}

func (dfa *DFA) AddState(stateStatus StateStatus) {
	dfa.states[uint(len(dfa.states))] = State{stateStatus, uint(len(dfa.states)), map[int32]uint{}}
}

func (dfa DFA) Depth() uint {
	var stateMap = make(map[uint]uint)
	var maxValue uint

	stateMap = dfa.DepthUtil(dfa.startingState, 0, stateMap)

	for _, v := range stateMap {
		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue
}

func (dfa DFA) DepthUtil(state State, count uint, stateMap map[uint]uint) map[uint]uint {
	stateMap[state.stateID] = count

	for _, v := range state.transitions {
		if _, ok := stateMap[v]; !ok {
			stateMap = dfa.DepthUtil(dfa.states[v], count+1, stateMap)
		}
	}

	return stateMap
}

func (dfa DFA) Describe(detail bool) {
	fmt.Println("This DFA has", len(dfa.states), "states and", len(dfa.alphabet), "alphabet")
	if detail {
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
		fmt.Println("Accepting States:")
		for _, v := range dfa.GetAcceptingStates() {
			fmt.Println(v.stateID)
		}
		fmt.Println("Rejecting States:")
		for _, v := range dfa.GetRejectingStates() {
			fmt.Println(v.stateID)
		}
		fmt.Println("Starting State:", dfa.startingState.stateID)
		fmt.Println("Alphabet:")
		for key, _ := range dfa.alphabet {
			fmt.Println(string(key))
		}
	}
}

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

func GetPTAFromListOfStringInstances(strings []StringInstance, APTA bool) DFA {
	strings = SortListOfStringInstances(strings)
	var exists bool
	var count int
	var alphabet = make(map[int32]bool)
	var states = make(map[uint]State)
	var startingStateID uint = 0
	var currentStateID uint

	if strings[0].length == 0 {
		if strings[0].stringStatus == ACCEPTING {
			states[0] = State{ACCEPTING, 0, map[int32]uint{}}
		} else {
			states[0] = State{REJECTING, 0, map[int32]uint{}}
		}
	} else {
		states[0] = State{UNKNOWN, 0, map[int32]uint{}}
	}

	for _, stringInstance := range strings {
		if !APTA && stringInstance.stringStatus != ACCEPTING {
			continue
		}
		currentStateID = startingStateID
		count = 0
		for _, character := range stringInstance.stringValue {
			count++
			exists = false
			//alphabet check
			if !alphabet[character] {
				alphabet[character] = true
			}

			if value, ok := states[currentStateID].transitions[character]; ok {
				currentStateID = value
				exists = true
			}

			if !exists {
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						states[uint(len(states))] = State{ACCEPTING, uint(len(states)), map[int32]uint{}}
					} else {
						states[uint(len(states))] = State{REJECTING, uint(len(states)), map[int32]uint{}}
					}
				} else {
					states[uint(len(states))] = State{UNKNOWN, uint(len(states)), map[int32]uint{}}
				}
				states[currentStateID].transitions[character] = states[uint(len(states))-1].stateID
				currentStateID = states[uint(len(states))-1].stateID
			} else {
				// last symbol in string check
				if count == len(stringInstance.stringValue) {
					if stringInstance.stringStatus == ACCEPTING {
						if states[currentStateID].stateStatus == REJECTING {
							panic("State already set to rejecting, cannot set to accepting")
						} else {
							tempState := states[currentStateID]
							tempState.stateStatus = ACCEPTING
							states[currentStateID] = tempState
						}
					} else {
						if states[currentStateID].stateStatus == ACCEPTING {
							panic("State already set to accepting, cannot set to rejecting")
						} else {
							tempState := states[currentStateID]
							tempState.stateStatus = REJECTING
							states[currentStateID] = tempState
						}
					}
				}
			}
		}
	}
	return DFA{states, states[startingStateID], alphabet}
}
