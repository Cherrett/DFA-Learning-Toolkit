package util

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

// Void struct represents an empty struct.
// Used to decrease memory overhead of maps.
type Void struct{}

// Null variable of type Void.
// Used to decrease memory overhead of maps.
var Null Void

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smallest of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// MaxSlice returns the largest value within a slice.
func MaxSlice(slice []int) int {
	maxValue := 0

	for element := range slice {
		if element > maxValue {
			maxValue = element
		}
	}

	return maxValue
}

// SumSlice returns the summed values within a slice.
func SumSlice(slice []int) int {
	count := 0

	for element := range slice {
		count += element
	}

	return count
}

// SumMap returns the summed values within a map.
func SumMap(currentMap map[int]int, key bool) int {
	count := 0

	if key {
		for key := range currentMap {
			count += key
		}
	} else {
		for _, element := range currentMap {
			count += element
		}
	}

	return count
}

// FileExists checks if a given.
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// MinMaxAvg struct is used to keep track of
// minimum, maximum and average values given
// a number of values.
type MinMaxAvg struct {
	min   float64 // Minimum of values.
	max   float64 // Maximum of values.
	sum   float64 // Sum of values.
	count uint    // Number of values.
}

// NewMinMaxAvg returns an empty MinMaxAvg struct.
func NewMinMaxAvg() MinMaxAvg {
	return MinMaxAvg{
		min:   math.Inf(1),
		max:   math.Inf(-1),
		sum:   0,
		count: 0,
	}
}

// Add adds a float value to the MinMaxAvg struct.
func (minMaxAvg *MinMaxAvg) Add(value float64) {
	// If value is smaller than the minimum value,
	// set minimum value within struct to value.
	if value < minMaxAvg.min {
		minMaxAvg.min = value
	}

	// If value is larger than the maximum value,
	// set maximum value within struct to value.
	if value > minMaxAvg.max {
		minMaxAvg.max = value
	}

	// Add value to sum.
	minMaxAvg.sum += value

	// Increment counter.
	minMaxAvg.count++
}

// AddInt adds an integer value to the MinMaxAvg struct.
func (minMaxAvg *MinMaxAvg) AddInt(intValue int) {
	value := float64(intValue)
	// If value is smaller than the minimum value,
	// set minimum value within struct to value.
	if value < minMaxAvg.min {
		minMaxAvg.min = value
	}

	// If value is larger than the maximum value,
	// set maximum value within struct to value.
	if value > minMaxAvg.max {
		minMaxAvg.max = value
	}

	// Add value to sum.
	minMaxAvg.sum += value

	// Increment counter.
	minMaxAvg.count++
}

// Min returns the minimum value within the MinMaxAvg struct.
func (minMaxAvg MinMaxAvg) Min() float64 {
	return minMaxAvg.min
}

// Max returns the maximum value within the MinMaxAvg struct.
func (minMaxAvg MinMaxAvg) Max() float64 {
	return minMaxAvg.max
}

// Avg returns the average value within the MinMaxAvg struct.
func (minMaxAvg MinMaxAvg) Avg() float64 {
	// Get average by dividing the sum of elements
	// by the number of elements within struct.
	return minMaxAvg.sum / float64(minMaxAvg.count)
}

// Factorial returns the factorial of n by recursively
// calling itself.
func Factorial(n int)(result int) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}

// DownloadAllStaminaDatasets downloads all of the stamina datasets to a given directory.
func DownloadAllStaminaDatasets(directory string){
	// Iterate from 1 to 100 (number of stamina datasets).
	for i := 1; i < 101; i++{
		// Get training and test sets from URL.
		resp, _ := http.Get(fmt.Sprintf("http://stamina.chefbe.net/downloads/grid/%d_training.txt", i))
		resp2, _ := http.Get(fmt.Sprintf("http://stamina.chefbe.net/downloads/grid/%d_test.txt", i))

		// Create training file.
		out, err := os.Create(fmt.Sprintf("%s/%d_training.txt", directory, i))
		if err != nil {
			panic("Training file failed to be created.")
		}

		// Create test file.
		out2, err := os.Create(fmt.Sprintf("%s/%d_test.txt", directory, i))
		if err != nil {
			panic("Testing file failed to be created.")
		}

		// Copy to files.
		_, _ = io.Copy(out, resp.Body)
		_, err = io.Copy(out2, resp2.Body)

		// Close io/file buffers.
		_ = resp.Body.Close()
		_ = resp2.Body.Close()
		out.Close()
		out2.Close()

		fmt.Printf("Downloaded dataset %d/100.\n", i)
	}
}
