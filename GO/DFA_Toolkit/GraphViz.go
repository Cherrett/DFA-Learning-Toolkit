package DFA_Toolkit

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func (dfa DFA) ToDOT(filePath string, rankByOrder bool, topDown bool) {
	if rankByOrder {
		// If depth and order for DFA is not computed,
		// call CalculateDepthAndOrder.
		if !dfa.ComputedDepthAndOrder {
			dfa.CalculateDepthAndOrder()
		}
	}

	file, err := os.Create(filePath)

	if err != nil {
		panic("Invalid file path/name")
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	if topDown{
		_, _ = writer.WriteString("digraph g{\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	}else{
		_, _ = writer.WriteString("digraph g{\n\trankdir=LR;\n\tgraph [dpi=300 ordering=\"out\"];\n\tmargin=0;\n\tnull [style=invis];\n")
	}

	for stateID, state := range dfa.States{
		if rankByOrder {
			if state.StateStatus == UNKNOWN{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, state.order))
			}else if state.StateStatus == ACCEPTING{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" peripheries=2 fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, state.order))
			}else{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"setlinewidth(3),filled\" fillcolor=\"white\"];\n", stateID, state.order))
			}
		}else{
			if state.StateStatus == UNKNOWN{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, stateID))
			}else if state.StateStatus == ACCEPTING{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" peripheries=2 fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"filled\" fillcolor=\"white\"];\n", stateID, stateID))
			}else{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d [label=\"q%d\" shape=\"circle\" fontname=verdana fontsize=8 color=\"black\" fontcolor=\"black\" style=\"setlinewidth(3),filled\" fillcolor=\"white\"];\n", stateID, stateID))
			}

		}

	}

	_, _ = writer.WriteString(fmt.Sprintf("\tnull->q%d;\n", dfa.StartingStateID))

	for stateID, state := range dfa.States{
		for symbol, symbolID := range dfa.SymbolMap {
			resultantStateID := state.Transitions[symbolID]
			if resultantStateID > -1{
				_, _ = writer.WriteString(fmt.Sprintf("\tq%d->q%d [label=\"%s\" fontname=verdana fontsize=8];\n", stateID, resultantStateID, string(symbol)))
			}
		}
	}

	_, _ = writer.WriteString("}")
	_ = writer.Flush()
}

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