package DFA_Toolkit

import (
	"math"
	"math/rand"
)

func AbbadingoDFA(numberOfStates int, exact bool) DFA{
	dfaSize := (5.0 * numberOfStates) / 4.0
	dfaDepth := uint((2.0 * math.Log2(float64(numberOfStates))) - 2.0)

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

		dfa.Minimise()

		if dfa.Depth() == dfaDepth{
			if exact{
				if len(dfa.states) == dfaSize{
					return dfa
				}
			}else{
				return dfa
			}
		}
	}
}