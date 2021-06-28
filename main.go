package main

import (
	"bufio"
	"fmt"
	dfalearningtoolkit "github.com/Cherrett/DFA-Learning-Toolkit/core"
	"os"
	"strconv"
	"strings"
)

func main() {
	var dfa dfalearningtoolkit.DFA
	var dataset dfalearningtoolkit.Dataset = nil
	reader := bufio.NewReader(os.Stdin)
	exit := false
	restart := false

	for !exit {
		restart = false
		fmt.Println("\nSelect from the following options\n1. Generate DFA and Dataset\n2. Read DFA/APTA from file\n3. Read Dataset from file\n4. Exit")

		choice := readIntInput("Choice: ", reader)

		// First main choice.
		switch choice {
		case 1:
			// Generate DFA and Dataset
			valid := false

			for !valid {
				fmt.Println("\nSelect from the following options\n1. Generate DFA and Dataset using the Abbadingo Protocol\n2. Generate DFA and Dataset using the Stamina Protocol")
				choice = readIntInput("Choice: ", reader)

				switch choice {
				case 1:
					valid = true
					targetDFASize := readIntInput("Input size of target DFA: ", reader)
					numberOfTrainingExamples := readIntInput("Input number of training examples: ", reader)
					_, dataset, _ = dfalearningtoolkit.AbbadingoInstanceExact(targetDFASize, true, numberOfTrainingExamples, 0)
				case 2:
					valid = true
					alphabetSize := readIntInput("Input alphabet size: ", reader)
					targetDFASize := readIntInput("Input size of target DFA: ", reader)
					sparsityPercentage := readFloatInput("Input sparsity percentage of training set (in floating point form): ", reader)

					_, dataset, _ = dfalearningtoolkit.DefaultStaminaInstance(alphabetSize, targetDFASize, sparsityPercentage)
				default:
					fmt.Println("Invalid choice. Please choose from options 1-2.")
				}
			}

		case 2:
			// Read DFA/APTA from file
			valid := false

			for !valid {
				fmt.Println("\nSelect from the following options\n1. Read DFA/APTA from JSON\n2. Read DFA/APTA from Stamina file")
				choice = readIntInput("Choice: ", reader)

				switch choice {
				case 1:
					valid = true
					filePath := readInput("Input file path of JSON DFA file: ", reader)
					dfa, valid = dfalearningtoolkit.DFAFromJSON(filePath)
				case 2:
					filePath := readInput("Input file path of Stamina DFA file: ", reader)
					dfa = dfalearningtoolkit.GetDFAFromStaminaFile(filePath)
					valid = true
				default:
					fmt.Println("Invalid choice. Please choose from options 1-2.")
				}
			}
		case 3:
			// Read Dataset from file
			valid := false

			for !valid {
				fmt.Println("\nSelect from the following options\n1. Read Dataset from JSON\n2. Read Dataset from Abbadingo file\n3. Read Dataset from Stamina file")
				choice = readIntInput("Choice: ", reader)

				switch choice {
				case 1:
					valid = true
					filePath := readInput("Input file path of JSON Dataset file: ", reader)
					dataset, valid = dfalearningtoolkit.DatasetFromJSON(filePath)
				case 2:
					filePath := readInput("Input file path of Abbadingo Dataset file: ", reader)
					dataset = dfalearningtoolkit.GetDatasetFromAbbadingoFile(filePath)
					valid = true
				case 3:
					filePath := readInput("Input file path of Stamina Dataset file: ", reader)
					dataset = dfalearningtoolkit.GetDatasetFromStaminaFile(filePath)
					valid = true
				default:
					fmt.Println("Invalid choice. Please choose from options 1-2.")
				}
			}
		case 4:
			fmt.Println("Exiting.")
			exit = true
		default:
			fmt.Println("Invalid choice. Please choose from options 1-4.")
			restart = true
		}

		for !exit && !restart {
			// Second main choice.
			fmt.Println("\nSelect from the following options\n1. Run EDSM on DFA/APTA or Dataset generated/read\n2. Run RPNI on DFA/APTA or Dataset generated/read\n3. Run AutomataTeams on DFA or Dataset generated/read\n4. Visualisation of DFA/APTA\n5. Exit")
			choice = readIntInput("Choice: ", reader)

			switch choice {
			case 1:

			case 2:

			case 3:

			case 4:

			case 5:
				fmt.Println("Exiting.")
				exit = true
			default:
				fmt.Println("Invalid choice. Please choose from options 1-4.")
			}

			dfa.ToJSON("")
			dataset.ToJSON("")
		}
	}
}

func readInput(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	value, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Invalid input.")
		fmt.Print(prompt)
		value, err = reader.ReadString('\n')
	}

	return strings.TrimRight(value, "\r\n")
}

func readIntInput(prompt string, reader *bufio.Reader) int {
	fmt.Print(prompt)
	value, err := reader.ReadString('\n')
	intValue, err2 := strconv.Atoi(strings.TrimRight(value, "\r\n"))

	if err != nil || err2 != nil {
		fmt.Println("Invalid input.")
		fmt.Print(prompt)
		value, err = reader.ReadString('\n')
		intValue, err2 = strconv.Atoi(strings.TrimRight(value, "\r\n"))
	}

	return intValue
}

func readFloatInput(prompt string, reader *bufio.Reader) float64 {
	fmt.Print(prompt)
	value, err := reader.ReadString('\n')
	floatValue, err2 := strconv.ParseFloat(strings.TrimRight(value, "\r\n"), 64)

	if err != nil || err2 != nil {
		fmt.Println("Invalid input.")
		fmt.Print(prompt)
		value, err = reader.ReadString('\n')
		floatValue, err2 = strconv.ParseFloat(strings.TrimRight(value, "\r\n"), 64)
	}

	return floatValue
}
