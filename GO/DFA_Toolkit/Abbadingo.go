package DFA_Toolkit

import (
	"math"
	"math/rand"
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
			dfa.AddTransition(dfa.GetSymbolID('a'), stateID, rand.Intn(len(dfa.states)))
			dfa.AddTransition(dfa.GetSymbolID('b'), stateID, rand.Intn(len(dfa.states)))
		}
		dfa.startingState = rand.Intn(len(dfa.states))

		dfa = dfa.Minimise()

		if dfa.Depth() == dfaDepth{
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