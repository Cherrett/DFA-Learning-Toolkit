package main

import (
	"DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	var timings []int64
	iterations := 1

	for i := 0; i < iterations; i++ {
		start := time.Now()

		AbbadingoDFA := DFA_Toolkit.AbbadingoDFA(5, false)
		AbbadingoDFA.Describe(false)
		fmt.Println("DFA Depth:", AbbadingoDFA.Depth())

		timings = append(timings, time.Since(start).Milliseconds())
	}
	var sum int64
	for _, timing := range timings{
		sum += timing
	}
	fmt.Printf("Average Time: %vms\n", sum/int64(iterations))
}