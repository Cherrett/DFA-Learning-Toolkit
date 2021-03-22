package util

import (
	"math"
	"os"
)

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
func MaxSlice(slice []int) int{
	maxValue := 0

	for element := range slice{
		if element > maxValue{
			maxValue = element
		}
	}

	return maxValue
}

// SumSlice returns the summed values within a slice.
func SumSlice(slice []int) int{
	count := 0

	for element := range slice{
		count += element
	}

	return count
}

// SumMap returns the summed values within a map.
func SumMap(currentMap map[int]int, key bool) int{
	count := 0

	if key{
		for key := range currentMap{
			count += key
		}
	}else{
		for _, element := range currentMap{
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

type MinMaxAvg struct {
	min float64		// Minimum of elements.
	max float64		// Maximum of elements.
	sum float64		// Sum of elements.
	count uint		// Number of elements.
}

func NewMinMaxAvg() MinMaxAvg{
	return MinMaxAvg{
		min:   math.Inf(1),
		max:   math.Inf(-1),
		sum:   0,
		count: 0,
	}
}

func (minMaxAvg *MinMaxAvg) Add(element float64){
	if element < minMaxAvg.min{
		minMaxAvg.min = element
	}
	if element > minMaxAvg.max{
		minMaxAvg.max = element
	}
	minMaxAvg.sum += element
	minMaxAvg.count++
}

func (minMaxAvg MinMaxAvg) Min() float64{
	return minMaxAvg.min
}

func (minMaxAvg MinMaxAvg) Max() float64{
	return minMaxAvg.max
}

func (minMaxAvg MinMaxAvg) Avg() float64{
	return minMaxAvg.sum / float64(minMaxAvg.count)
}
