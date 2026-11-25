package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"time"
)

// Measures end-to-end execution, including OS process startup; stdout is discarded.
func benchmark(path string, runs int) {
	run := func() float64 {
		c := exec.Command(path)
		c.Stderr = os.Stderr
		start := time.Now()
		if err := c.Run(); err != nil {
			die("benchmark program failed: %v", err)
		}
		return float64(time.Since(start)) / float64(time.Millisecond)
	}
	run() // warmup
	times := make([]float64, runs)
	for i := range times {
		times[i] = run()
		fmt.Printf("Run %d: %.3f ms\n", i+1, times[i])
	}
	sort.Float64s(times)
	median := times[runs/2]
	if runs%2 == 0 {
		median = (times[runs/2-1] + median) / 2
	}
	fmt.Printf("Median: %.3f ms | Min: %.3f ms | Max: %.3f ms\n", median, times[0], times[runs-1])
	fmt.Println("Includes process startup; stdout discarded; stdin closed. One warm-up excluded.")
}
