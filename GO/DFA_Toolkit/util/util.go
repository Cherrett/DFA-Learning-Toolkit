package util

import (
	"fmt"
	"io"
	"math"
	"math/rand"
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

// StatsTracker struct is used to keep track of
// minimum, maximum, average, variance and
// standard deviation values given a sequence
// of values. Mean and variance calculation
// is done using Welford's online algorithm.
type StatsTracker struct {
	min   float64 // Minimum of values.
	max   float64 // Maximum of values.
	count uint    // Number of values.
	mean  float64 // Running average of values.
	m2    float64 // Running value used to calculate variance/standard dev.
}

// NewStatsTracker returns an empty StatsTracker struct.
func NewStatsTracker() StatsTracker {
	return StatsTracker{
		min:   math.Inf(1),
		max:   math.Inf(-1),
		count: 0,
		mean:  0.0,
		m2:    0.0,
	}
}

// Add adds a float value to the StatsTracker struct.
func (statsTracker *StatsTracker) Add(value float64) {
	// If value is smaller than the minimum value,
	// set minimum value within struct to value.
	if value < statsTracker.min {
		statsTracker.min = value
	}

	// If value is larger than the maximum value,
	// set maximum value within struct to value.
	if value > statsTracker.max {
		statsTracker.max = value
	}

	// Increment counter.
	statsTracker.count++

	muNew := statsTracker.mean + ((value - statsTracker.mean) / float64(statsTracker.count))

	statsTracker.m2 += (value - statsTracker.mean) * (value - muNew)

	statsTracker.mean = muNew
}

// AddInt adds an integer value to the StatsTracker struct.
func (statsTracker *StatsTracker) AddInt(intValue int) {
	// Cast to float64 and call Add function.
	statsTracker.Add(float64(intValue))
}

// Min returns the minimum value within the StatsTracker struct.
func (statsTracker StatsTracker) Min() float64 {
	return statsTracker.min
}

// Max returns the maximum value within the StatsTracker struct.
func (statsTracker StatsTracker) Max() float64 {
	return statsTracker.max
}

// Mean returns the average value within the StatsTracker struct.
func (statsTracker StatsTracker) Mean() float64 {
	// Get average by dividing the sum of elements
	// by the number of elements within struct.
	return statsTracker.mean
}

// PopulationVariance returns the population variance value within the StatsTracker struct.
func (statsTracker StatsTracker) PopulationVariance() float64 {
	if statsTracker.count > 1 {
		return statsTracker.m2 / float64(statsTracker.count)
	} else {
		return 0.0
	}
}

// SampleVariance returns the sample variance value within the StatsTracker struct.
func (statsTracker StatsTracker) SampleVariance() float64 {
	if statsTracker.count > 1 {
		return statsTracker.m2 / float64(statsTracker.count-1)
	} else {
		return 0.0
	}
}

// PopulationStandardDev returns the population standard deviation value within the StatsTracker struct.
func (statsTracker StatsTracker) PopulationStandardDev() float64 {
	return math.Sqrt(statsTracker.PopulationVariance())
}

// SampleStandardDev returns the sample standard deviation value within the StatsTracker struct.
func (statsTracker StatsTracker) SampleStandardDev() float64 {
	return math.Sqrt(statsTracker.SampleVariance())
}

// Factorial returns the factorial of n by recursively
// calling itself.
func Factorial(n int) int {
	if n > 0 {
		return n * Factorial(n-1)
	}
	return 1
}

// RandomGeometricProbability returns a random geometrically distributed integer using a mean value.
func RandomGeometricProbability(mean float64) int {
	return int(math.Log(rand.Float64()) / (math.Log(1 - mean)))
}

// DownloadAllStaminaDatasets downloads all of the stamina datasets to a given directory.
func DownloadAllStaminaDatasets(directory string) {
	// Iterate from 1 to 100 (number of stamina datasets).
	for i := 1; i < 101; i++ {
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
