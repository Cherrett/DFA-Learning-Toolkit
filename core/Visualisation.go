package dfalearningtoolkit

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// SymbolAlphabetMappingAbbadingo represents the symbol-alphabet mapping for
// the Abbadingo competition standard (using only a's and b's / 0s and 1s).
var SymbolAlphabetMappingAbbadingo = map[int]string{0: "a", 1: "b"}

// Visualisation.go consists of various functions which create a visualisation of a
// given DFA or State Partition. Please note that Graphviz must be downloaded
// and installed before hand from https://graphviz.org/download.

// ToDOT creates a .dot file using the DOT language. This DOT file
// contains a representation for the given DFA which can then be used
// to generate a visual representation of the DFA. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a
// string. This can be set to nil, which will map each symbol with a
// number starting from 0. If showOrder is set to true, the canonical
// order of the states is shown inside the node within DFA. If topDown
// is set to true, the visual representation will be top down. A left
// to right representation is used otherwise. This function is also
// called from all of the functions below.
func (dfa *DFA) ToDOT(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) {
	// If show order is set to true, make sure
	// that depth and order for DFA are computed.
	if showOrder {
		// If depth and order for DFA are not computed,
		// call CalculateDepthAndOrder.
		if !dfa.computedDepthAndOrder {
			dfa.CalculateDepthAndOrder()
		}
	}

	// Create file given a path/name.
	file, err := os.Create(filePath)

	// If file was not created successfully,
	// print error and panic.
	if err != nil {
		panic("Invalid file path/name")
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize a new writer given the created file.
	writer := bufio.NewWriter(file)

	// If topDown is set to true, do not set rankdir field
	// within first line of DOT file. Otherwise, rankdir
	// is set to LR (left-to-right).
	if topDown {
		_, _ = writer.WriteString("digraph g{\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	} else {
		_, _ = writer.WriteString("digraph g{\n\trankdir=LR;\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	}

	// Iterate over each state and write the correct format
	// depending on the state label.
	for stateID, state := range dfa.States {
		// If showOrder is set to true, the order
		// of the state is shown within the label.
		label := ""
		if showOrder {
			label = fmt.Sprintf("\"q%d\\n(%d)\"", stateID, state.order)
		} else {
			label = fmt.Sprintf("\"q%d\"", stateID)
		}

		if state.Label == UNLABELLED {
			_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=%s shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, label))
		} else if state.Label == ACCEPTING {
			_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=%s shape=\"circle\" peripheries=2 fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, label))
		} else {
			_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=%s shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"setlinewidth(3),filled\" fillcolor=\"white\"];\n", stateID, label))
		}
	}

	// Create the arrow of the starting state.
	_, _ = writer.WriteString(fmt.Sprintf("\tnull->q%d;\n", dfa.StartingStateID))

	// Iterate over each transition and write to file.
	for stateID, state := range dfa.States {
		for symbol := range dfa.Alphabet {
			if symbolMapping == nil {
				resultantStateID := state.GetTransitionValue(symbol)
				if resultantStateID > -1 {
					_, _ = writer.WriteString(fmt.Sprintf("\tq%d->q%d [label=\"%d\" fontname=verdana fontsize=8];\n", stateID, resultantStateID, symbol))
				}
			} else {
				if value, exists := symbolMapping[symbol]; exists {
					resultantStateID := state.GetTransitionValue(symbol)
					if resultantStateID > -1 {
						_, _ = writer.WriteString(fmt.Sprintf("\tq%d->q%d [label=\"%s\" fontname=verdana fontsize=8];\n", stateID, resultantStateID, value))
					}
				} else {
					panic(fmt.Sprintf("Symbol ID %d not in symbolMapping map.", symbol))
				}
			}
		}
	}

	// Add closing curly bracket to file.
	_, _ = writer.WriteString("}")

	// Flush writer.
	_ = writer.Flush()
}

// ToPNG creates and saves a .png image which represents the DFA to the
// given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within DFA. If topDown is set to true, the visual
// representation will be top down. A left to right representation is used
// otherwise. Returns true if successful or false if an error occurs.
func (dfa DFA) ToPNG(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tpng", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("DFA outputted to: %s\n", filePath)
	return true
}

// ToJPG creates and saves a .jpg image which represents the DFA to the
// given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within DFA. If topDown is set to true, the visual
// representation will be top down. A left to right representation is used
// otherwise. Returns true if successful or false if an error occurs.
func (dfa DFA) ToJPG(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tjpg", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("DFA outputted to: %s\n", filePath)
	return true
}

// ToPDF creates and saves a .pdf file which represents the DFA to the
// given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within DFA. If topDown is set to true, the visual
// representation will be top down. A left to right representation is used
// otherwise. Returns true if successful or false if an error occurs.
func (dfa DFA) ToPDF(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tpdf", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("DFA outputted to: %s\n", filePath)
	return true
}

// ToSVG creates and saves a .svg file which represents the DFA to the
// given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within DFA. If topDown is set to true, the visual
// representation will be top down. A left to right representation is used
// otherwise. Returns true if successful or false if an error occurs.
func (dfa DFA) ToSVG(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tsvg", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("DFA outputted to: %s\n", filePath)
	return true
}

// ToDOT creates a .dot file using the DOT language. This DOT file
// contains a representation for the given StatePartition which can then
// be used to generate a visual representation of the DFA represented by the
// State Partition. The symbolMapping parameter is a map which maps each symbol
// within the alphabet to a string. This can be set to nil, which will map each
// symbol with a number starting from 0. If showOrder is set to true, the canonical
// order of the states is shown inside the node within partition. If topDown is set
// to true, the visual representation will be top down. A left to right representation
// is used otherwise. This function is also called from all of the functions below.
func (statePartition StatePartition) ToDOT(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) {
	var rootBlocks []int

	// If show order is set to true, get root blocks in order.
	// Else, just get root blocks.
	if showOrder {
		rootBlocks = statePartition.OrderedBlocks()
	} else {
		rootBlocks = statePartition.RootBlocks()
	}

	// Create file given a path/name.
	file, err := os.Create(filePath)

	// If file was not created successfully,
	// print error and panic.
	if err != nil {
		panic("Invalid file path/name")
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize a new writer given the created file.
	writer := bufio.NewWriter(file)

	// If topDown is set to true, do not set rankdir field
	// within first line of DOT file. Otherwise, rankdir
	// is set to LR (left-to-right).
	if topDown {
		_, _ = writer.WriteString("digraph g{\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	} else {
		_, _ = writer.WriteString("digraph g{\n\trankdir=LR;\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	}

	// Iterate over each blockID and write the correct format
	// depending on the blockID label.
	for index, blockID := range rootBlocks {
		statesWithinBlock := statePartition.ReturnSet(blockID)
		label := fmt.Sprintf("\"q%d\\n", blockID)
		for index2, stateID := range statesWithinBlock {
			label += fmt.Sprintf("%d", stateID)
			if index2 < len(statesWithinBlock)-1 {
				label += "\\n"
			}
		}
		// If showOrder is set to true, the order
		// of the blockID is shown within the label.
		if showOrder {
			label += fmt.Sprintf("\\n(%d)", index)
		}
		label += "\""

		if statePartition.Blocks[blockID].Label == UNLABELLED {
			_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=%s shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", blockID, label))
		} else if statePartition.Blocks[blockID].Label == ACCEPTING {
			_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=%s shape=\"circle\" peripheries=2 fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", blockID, label))
		} else {
			_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=%s shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"setlinewidth(3),filled\" fillcolor=\"white\"];\n", blockID, label))
		}
	}

	// Create the arrow of the starting blockID.
	_, _ = writer.WriteString(fmt.Sprintf("\tnull->q%d;\n", statePartition.StartingBlock()))

	// Iterate over each transition and write to file.
	for _, blockID := range rootBlocks {
		for symbol := 0; symbol < statePartition.AlphabetSize; symbol++ {
			if symbolMapping == nil {
				if resultantBlockID := statePartition.Blocks[blockID].Transitions[symbol]; resultantBlockID > -1 {
					resultantBlockRootID := statePartition.Find(resultantBlockID)
					_, _ = writer.WriteString(fmt.Sprintf("\tq%d->q%d [label=\"%d\" fontname=verdana fontsize=8];\n", blockID, resultantBlockRootID, symbol))
				}
			} else {
				if value, exists := symbolMapping[symbol]; exists {
					if resultantBlockID := statePartition.Blocks[blockID].Transitions[symbol]; resultantBlockID > -1 {
						resultantBlockRootID := statePartition.Find(resultantBlockID)
						_, _ = writer.WriteString(fmt.Sprintf("\tq%d->q%d [label=\"%s\" fontname=verdana fontsize=8];\n", blockID, resultantBlockRootID, value))
					}
				} else {
					panic(fmt.Sprintf("Symbol ID %d not in symbolMapping map.", symbol))
				}
			}
		}
	}

	// Add closing curly bracket to file.
	_, _ = writer.WriteString("}")

	// Flush writer.
	_ = writer.Flush()
}

// ToPNG creates and saves a .png image which represents the state partition
// to the given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within partition. If topDown is set to true, the
// visual representation will be top down. A left to right representation is
// used otherwise. Returns true if successful or false if an error occurs.
func (statePartition StatePartition) ToPNG(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	statePartition.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tpng", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("State Partition outputted to: %s\n", filePath)
	return true
}

// ToJPG creates and saves a .jpg image which represents the state partition
// to the given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within partition. If topDown is set to true, the
// visual representation will be top down. A left to right representation is
// used otherwise. Returns true if successful or false if an error occurs.
func (statePartition StatePartition) ToJPG(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	statePartition.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tjpg", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("State Partition outputted to: %s\n", filePath)
	return true
}

// ToPDF creates and saves a .pdf file which represents the state partition
// to the given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within partition. If topDown is set to true, the
// visual representation will be top down. A left to right representation is
// used otherwise. Returns true if successful or false if an error occurs.
func (statePartition StatePartition) ToPDF(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	statePartition.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tpdf", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("State Partition outputted to: %s\n", filePath)
	return true
}

// ToSVG creates and saves a .svg file which represents the state partition
// to the given file path. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/. The symbolMapping
// parameter is a map which maps each symbol within the alphabet to a string.
// This can be set to nil, which will map each symbol with a number starting
// from 0. If showOrder is set to true, the canonical order of the states is
// shown inside the node within partition. If topDown is set to true, the
// visual representation will be top down. A left to right representation is
// used otherwise. Returns true if successful or false if an error occurs.
func (statePartition StatePartition) ToSVG(filePath string, symbolMapping map[int]string, showOrder bool, topDown bool) bool {
	defer os.Remove("temp.dot")
	statePartition.ToDOT("temp.dot", symbolMapping, showOrder, topDown)

	cmd := exec.Command("dot", "-Tsvg", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
		return false
	}

	fmt.Printf("State Partition outputted to: %s\n", filePath)
	return true
}
