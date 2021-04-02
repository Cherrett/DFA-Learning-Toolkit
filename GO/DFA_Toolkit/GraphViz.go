package dfatoolkit

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// GraphViz.go consists of various functions which create a visualisation
// of a given DFA. Please note that GraphViz must be downloaded and
// installed before hand from https://graphviz.org/download/.

// ToDOT creates a .dot file using the DOT language. This DOT file
// contains a representation for the given DFA which can then be used
// to generate a visual representation of the DFA. If rankByOrder is set
// to true, the state IDs/labels are ordered by their canonical order
// within DFA. If topDown is set to true, the visual representation
// will be top down. A left to right representation is used otherwise.
// This function is also called from all of the functions below.
func (dfa DFA) ToDOT(filePath string, rankByOrder bool, topDown bool) {
	// If rank by order is set to true, make sure
	// that depth and order for DFA are computed.
	if rankByOrder {
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
	if topDown{
		_, _ = writer.WriteString("digraph g{\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	}else{
		_, _ = writer.WriteString("digraph g{\n\trankdir=LR;\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	}

	// Iterate over each state and write the correct format
	// depending on the state label.
	for stateID, state := range dfa.States{
		// If rankByOrder is set to true, the order
		// of the state is used as the label.
		if rankByOrder {
			if state.Label == UNKNOWN{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, state.order))
			}else if state.Label == ACCEPTING{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" peripheries=2 fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, state.order))
			}else{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"setlinewidth(3),filled\" fillcolor=\"white\"];\n", stateID, state.order))
			}
		}else{
			if state.Label == UNKNOWN{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, stateID))
			}else if state.Label == ACCEPTING{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" peripheries=2 fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, stateID))
			}else{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"setlinewidth(3),filled\" fillcolor=\"white\"];\n", stateID, stateID))
			}

		}

	}

	// Create the arrow of the starting state.
	_, _ = writer.WriteString(fmt.Sprintf("\tnull->q%d;\n", dfa.StartingStateID))

	// Iterate over each transition and write to file.
	for stateID, state := range dfa.States{
		for symbol, symbolID := range dfa.SymbolMap {
			resultantStateID := state.Transitions[symbolID]
			if resultantStateID > -1{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d->q%d [label=\"%s\" fontname=verdana fontsize=8];\n", stateID, resultantStateID, string(symbol)))
			}
		}
	}

	// Add closing curly bracket to file.
	_, _ = writer.WriteString("}")

	// Flush writer.
	_ = writer.Flush()
}

// ToPNG creates and saves a .png image to the given file path.
// Please note that GraphViz must be downloaded and installed before
// hand from https://graphviz.org/download/. If rankByOrder is set
// to true, the state IDs/labels are ordered by their canonical order
// within DFA. If topDown is set to true, the visual representation
// will be top down. A left to right representation is used otherwise.
func (dfa DFA) ToPNG(filePath string, rankByOrder bool, topDown bool){
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", rankByOrder, topDown)

	cmd := exec.Command("dot", "-Tpng", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
	}else{
		fmt.Printf("DFA outputted to: %s\n", filePath)
	}
}

// ToJPG creates and saves a .jpg image to the given file path.
// Please note that GraphViz must be downloaded and installed before
// hand from https://graphviz.org/download/. If rankByOrder is set
// to true, the state IDs/labels are ordered by their canonical order
// within DFA. If topDown is set to true, the visual representation
// will be top down. A left to right representation is used otherwise.
func (dfa DFA) ToJPG(filePath string, rankByOrder bool, topDown bool){
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", rankByOrder, topDown)

	cmd := exec.Command("dot", "-Tjpg", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
	}else{
		fmt.Printf("DFA outputted to: %s\n", filePath)
	}
}

// ToPDF creates and saves a .pdf file to the given file path.
// Please note that GraphViz must be downloaded and installed before
// hand from https://graphviz.org/download/. If rankByOrder is set
// to true, the state IDs/labels are ordered by their canonical order
// within DFA. If topDown is set to true, the visual representation
// will be top down. A left to right representation is used otherwise.
func (dfa DFA) ToPDF(filePath string, rankByOrder bool, topDown bool){
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", rankByOrder, topDown)

	cmd := exec.Command("dot", "-Tpdf", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
	}else{
		fmt.Printf("DFA outputted to: %s\n", filePath)
	}
}

// ToSVG creates and saves a .svg file to the given file path.
// Please note that GraphViz must be downloaded and installed before
// hand from https://graphviz.org/download/. If rankByOrder is set
// to true, the state IDs/labels are ordered by their canonical order
// within DFA. If topDown is set to true, the visual representation
// will be top down. A left to right representation is used otherwise.
func (dfa DFA) ToSVG(filePath string, rankByOrder bool, topDown bool){
	defer os.Remove("temp.dot")
	dfa.ToDOT("temp.dot", rankByOrder, topDown)

	cmd := exec.Command("dot", "-Tsvg", "temp.dot", "-o", filePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Kindly download Graphviz executable from https://graphviz.org/download/")
	}else{
		fmt.Printf("DFA outputted to: %s\n", filePath)
	}
}