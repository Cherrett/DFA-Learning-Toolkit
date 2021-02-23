package DFA_Toolkit

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetDatasetFromAbbadingoFile(fileName string) Dataset {
	dataset := Dataset{}

	file, err := os.Open(fileName)

	if err != nil {
		panic("Invalid file path/name")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // ignore first line
	for scanner.Scan() {
		dataset = append(dataset, NewStringInstanceFromAbbadingoFile(scanner.Text(), " "))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return dataset
}

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

func AbbadingoDFA(numberOfStates int, exact bool) DFA{
	dfaSize := int(math.Round((5.0 * float64(numberOfStates)) / 4.0))
	dfaDepth := uint(math.Round((2.0 * math.Log2(float64(numberOfStates))) - 2.0))
	// random seed
	rand.Seed(time.Now().UnixNano())
	for{
		dfa := NewDFA()
		dfa.AddSymbols([]rune{'a', 'b'})

		for i := 0; i < dfaSize; i++{
			if rand.Intn(2) == 0{
				dfa.AddState(ACCEPTING)
			}else{
				dfa.AddState(UNKNOWN)
			}
		}

		for stateID := range dfa.states{
			dfa.AddTransition(dfa.SymbolID('a'), stateID, rand.Intn(len(dfa.states)))
			dfa.AddTransition(dfa.SymbolID('b'), stateID, rand.Intn(len(dfa.states)))
		}
		dfa.startingState = rand.Intn(len(dfa.states))

		dfa = dfa.Minimise()
		currentDFADepth := dfa.Depth()

		if currentDFADepth == dfaDepth{
			if exact{
				if len(dfa.states) == numberOfStates{
					return dfa
				}
			}else{
				return dfa
			}
		}
	}
}

func AbbadingoDataset(dfa DFA, percentageFromSamplePool float64, testingRatio float64) (Dataset, Dataset){
	trainingDataset := Dataset{}
	testingDataset := Dataset{}
	maxLength := math.Round((2.0 * math.Log2(float64(len(dfa.states)))) + 3.0)
	maxDecimal := math.Pow(2, maxLength + 1) - 1
	totalSetSize := math.Round((percentageFromSamplePool / 100) * maxDecimal)
	trainingSetSize := int(math.Round((1 - testingRatio) * totalSetSize))

	// random seed
	rand.Seed(time.Now().UnixNano())
	// map to avoid duplicate values
	valueMap := map[int]bool{}

	for x := 0; x < (int(totalSetSize)); x++{
		// get random value from range [1, totalSetSize]
		value := rand.Intn(int(maxDecimal)) + 1
		// if value is duplicate decrement x and go to next loop
		// else write new value to map
		if valueMap[value]{
			x--
			continue
		}else{
			valueMap[value] = true
		}

		// convert value to binary string
		binaryString := strconv.FormatInt(int64(value), 2)
		// remove first '1'
		binaryString = binaryString[1:]

		if trainingDataset.AcceptingStringInstancesCount() +
			trainingDataset.RejectingStringInstancesCount() < trainingSetSize{
			trainingDataset = append(trainingDataset, BinaryStringToStringInstance(dfa, binaryString))
		}else{
			testingDataset = append(testingDataset, BinaryStringToStringInstance(dfa, binaryString))
		}
	}

	return trainingDataset, testingDataset
}

func (dataset Dataset) WriteToAbbadingoFile(filePath string){
	sortedDataset := dataset.SortDatasetByLength()
	file, err := os.Create(filePath)

	if err != nil {
		panic("Invalid file path/name")
	}

	defer file.Close()
	writer := bufio.NewWriter(file)

	_, _ = writer.WriteString(strconv.Itoa(len(dataset)) + " 2\n")

	for _, stringInstance := range sortedDataset{
		outputString := strconv.Itoa(int(stringInstance.status))+" "+strconv.Itoa(int(stringInstance.length))+" "
		for _, symbol := range stringInstance.value{
			if symbol == 'a'{
				outputString += "0 "
			}else{
				outputString += "1 "
			}
		}
		outputString = strings.TrimSuffix(outputString, " ") + "\n"
		_, _ = writer.WriteString(outputString)
	}

	_ = writer.Flush()
}