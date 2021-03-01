//https://github.com/theodesp/unionfind
package DFA_Toolkit

type StatePartition struct {
	root []int
	size []int
}

// NewStatePartition returns an initialized list of size
func NewStatePartition(size int) *StatePartition {
	return new(StatePartition).init(size)
}

// Constructor initializes root and size arrays
func (statePartition *StatePartition) init(size int) *StatePartition {
	statePartition = new(StatePartition)
	statePartition.root = make([]int, size)
	statePartition.size = make([]int, size)

	for i := 0; i < size; i++ {
		statePartition.root[i] = i
		statePartition.size[i] = 1
	}

	return statePartition
}

// Union connects p and q by finding their roots and comparing their respective
// size arrays to keep the tree flat
func (statePartition *StatePartition) Union(p int, q int) {
	qRoot := statePartition.Root(q)
	pRoot := statePartition.Root(p)

	if statePartition.size[qRoot] < statePartition.size[pRoot] {
		statePartition.root[qRoot] = statePartition.root[pRoot]
		statePartition.size[pRoot] += statePartition.size[qRoot]
	} else {
		statePartition.root[pRoot] = statePartition.root[qRoot]
		statePartition.size[qRoot] += statePartition.size[pRoot]
	}
}

// Root or Find traverses each parent element while compressing the
// levels to find the root element of p
// If we attempt to access an element outside the array it returns -1
func (statePartition *StatePartition) Root(p int) int {
	if p > len(statePartition.root)-1 {
		return -1
	}

	for statePartition.root[p] != p {
		statePartition.root[p] = statePartition.root[statePartition.root[p]]
		p = statePartition.root[p]
	}

	return p
}

// Root or Find
func (statePartition *StatePartition) Find(p int) int {
	return statePartition.Root(p)
}

// Check if items p,q are connected
func (statePartition *StatePartition) Connected(p int, q int) bool {
	return statePartition.Root(p) == statePartition.Root(q)
}

func (dfa DFA) ToStatePartition() *StatePartition {
	return NewStatePartition(len(dfa.States))
}

func (statePartition StatePartition) ToDFA(dfa DFA) (bool, DFA){
	newMappings := map[int]int{}

	resultantDFA := DFA{
		States:                nil,
		StartingStateID:       -1,
		SymbolMap:             dfa.SymbolMap,
		Depth:                 -1,
		ComputedDepthAndOrder: false,
	}

	for _, stateID := range statePartition.root {
		if newStateID, ok := newMappings[stateID]; ok {
			if (resultantDFA.States[newStateID].StateStatus == ACCEPTING &&
				dfa.States[stateID].StateStatus == REJECTING) ||
				(resultantDFA.States[newStateID].StateStatus == REJECTING &&
					dfa.States[stateID].StateStatus == ACCEPTING){
				return false, DFA{}
			}else{
				resultantDFA.States[newStateID].StateStatus = dfa.States[stateID].StateStatus
			}
		}else{
			newMappings[stateID] = resultantDFA.AddState(dfa.States[stateID].StateStatus)
		}
	}

	// update new states via mappings
	for stateID := range dfa.States{
		for symbolID := 0; symbolID < len(dfa.SymbolMap); symbolID++ {
			if dfa.States[stateID].Transitions[symbolID] > -1{
				newStateID := newMappings[statePartition.root[stateID]]
				resultantStateID := newMappings[statePartition.root[dfa.States[stateID].Transitions[symbolID]]]
				if resultantDFA.States[newStateID].Transitions[symbolID] != -1 &&
					resultantDFA.States[newStateID].Transitions[symbolID] != resultantStateID{
					panic("Not deterministic")
				}else{
					resultantDFA.States[newStateID].Transitions[symbolID] = resultantStateID
				}
			}
		}
	}
	// update starting state via mappings
	resultantDFA.StartingStateID = newMappings[statePartition.root[dfa.StartingStateID]]

	return true, resultantDFA
}

func (statePartition *StatePartition) mergeStates(dfa DFA, state1 int, state2 int) bool{
	if (dfa.States[state1].StateStatus == ACCEPTING &&
		dfa.States[state2].StateStatus== REJECTING) ||
		(dfa.States[state2].StateStatus == ACCEPTING &&
			dfa.States[state1].StateStatus== REJECTING){
		return false
	}

	return true
}