package DFA_Toolkit

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

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

	for x := 0; x < (int(totalSetSize)); x++{
		// get random value from range [1, totalSetSize]
		value := rand.Intn(int(maxDecimal)) + 1
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