package util

import (
	"math"
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

// Add adds an element to the MinMaxAvg struct.
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
