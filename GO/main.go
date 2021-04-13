package main

import (
	dfatoolkit "DFA_Toolkit/DFA_Toolkit"
	"fmt"
	"time"
)

func main() {
	// PROFILING
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// go tool pprof -http=:8081 cpu.pprof

	start := time.Now()

	for i := 0; i < 100; i++{
		_ = dfatoolkit.AbbadingoDFA(32, true)
	}

	fmt.Printf("Time in seconds: %.2fs\n\n", time.Since(start).Seconds())
}
