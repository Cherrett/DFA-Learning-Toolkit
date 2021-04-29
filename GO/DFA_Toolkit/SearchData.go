package dfatoolkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// SearchData struct to store merge search data.
type SearchData struct {
	Merges               []StatePairScore // Slice of state pairs and scores of merges done.
	AttemptedMergesCount int              // The number of attempted merges.
	ValidMergesCount	 int			  // The number of valid attempted merges.
	Duration             time.Duration    // The duration for the search function to finish.
}

// MergesCount returns the number of merges done before
// reaching the final StatePartition.
func (searchData SearchData) MergesCount() int{
	return len(searchData.Merges)
}

// AttemptedMergesPerSecond returns the number of attempted
// merges per second. Used for performance evaluation.
func (searchData SearchData) AttemptedMergesPerSecond() float64{
	return float64(searchData.AttemptedMergesCount) / searchData.Duration.Seconds()
}

// ToJSON saves the SearchData to a JSON file given a file path.
func (searchData SearchData) ToJSON(filePath string) bool {
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

	// Convert SearchData to JSON.
	resultantJSON, err := json.MarshalIndent(searchData, "", "\t")

	// If SearchData was not converted successfully,
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

// SearchDataFromJSON returns SearchData read from a JSON file
// given a file path. The boolean value returned is set to
// true if SearchData was read successfully.
func SearchDataFromJSON(filePath string) (SearchData, bool) {
	// Open file from given a path/name.
	file, err := os.Open(filePath)

	// If file was not opened successfully,
	// return empty SearchData and false.
	if err != nil {
		return SearchData{}, false
	}

	// Close file at end of function.
	defer file.Close()

	// Initialize empty SearchData.
	searchData := SearchData{}

	// Convert JSON to SearchData.
	err = json.NewDecoder(file).Decode(&searchData)

	// If JSON was not converted successfully,
	// return empty SearchData and false.
	if err != nil {
		return SearchData{}, false
	}

	// Return populated SearchData and true if reached.
	return searchData, true
}