package dfa_learning_toolkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// MergeData struct to store merge merge data.
type MergeData struct {
	Merges               []StatePairScore // Slice of state pairs and scores of merges done.
	AttemptedMergesCount int              // The number of attempted merges.
	ValidMergesCount     int              // The number of valid attempted merges.
	Duration             time.Duration    // The duration for the search function to finish.
}

// MergesCount returns the number of merges done before
// reaching the final StatePartition.
func (mergeData MergeData) MergesCount() int {
	return len(mergeData.Merges)
}

// AttemptedMergesPerSecond returns the number of attempted
// merges per second. Used for performance evaluation.
func (mergeData MergeData) AttemptedMergesPerSecond() float64 {
	return float64(mergeData.AttemptedMergesCount) / mergeData.Duration.Seconds()
}

// ToJSON saves the MergeData to a JSON file given a file path.
func (mergeData MergeData) ToJSON(filePath string) bool {
	// Create file given a path/name.
	file, err := os.Create(filePath)

	// If file was not created successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Close file at end of function.
	defer file.Close()

	// Convert MergeData to JSON.
	resultantJSON, err := json.MarshalIndent(mergeData, "", "\t")

	// If MergeData was not converted successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Copy JSON to file created.
	_, err = io.Copy(file, bytes.NewReader(resultantJSON))

	// If JSON was not copied successfully,
	// print error and return false.
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Return true if reached.
	return true
}

// mergeDataFromJSON returns MergeData read from a JSON file
// given a file path. The boolean value returned is set to
// true if MergeData was read successfully.
func mergeDataFromJSON(filePath string) (MergeData, bool) {
	// Open file from given a path/name.
	file, err := os.Open(filePath)

	// If file was not opened successfully,
	// return empty MergeData and false.
	if err != nil {
		return MergeData{}, false
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize empty MergeData.
	mergeData := MergeData{}

	// Convert JSON to MergeData.
	err = json.NewDecoder(file).Decode(&mergeData)

	// If JSON was not converted successfully,
	// return empty MergeData and false.
	if err != nil {
		return MergeData{}, false
	}

	// Return populated MergeData and true if reached.
	return mergeData, true
}
