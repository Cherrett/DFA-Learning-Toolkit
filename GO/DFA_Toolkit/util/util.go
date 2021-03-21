package util

// Max returns the larger of x or y
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// MaxSlice returns the largest value within a slice
func MaxSlice(slice []int) int{
	maxValue := 0

	for element := range slice{
		if element > maxValue{
			maxValue = element
		}
	}

	return maxValue
}

// SumSlice returns the summed values within a slice
func SumSlice(slice []int) int{
	count := 0

	for element := range slice{
		count += element
	}

	return count
}

// SumMap returns the summed values within a map
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