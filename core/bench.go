package core

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
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

// Internal timing uses an explicit marker emitted by the program after its
// timed section. End-to-end time is also reported; neither is subtracted.
func parseInternal(s string) (int64, bool) {
	var ns int64
	found := false
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "TRIPOL_NS:") {
			continue
		}
		n, e := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "TRIPOL_NS:")), 10, 64)
		if found || e != nil || n < 0 {
			return 0, false
		}
		ns = n
		found = true
	}
	return ns, found
}
func benchmarkInternal(path string, runs int) {
	run := func() (float64, float64) {
		c := exec.Command(path)
		c.Stderr = os.Stderr
		start := time.Now()
		out, e := c.Output()
		elapsed := float64(time.Since(start)) / float64(time.Millisecond)
		if e != nil {
			die("benchmark program failed: %v", e)
		}
		ns, ok := parseInternal(string(out))
		if !ok {
			die("--internal requires exactly one TRIPOL_NS: <nonnegative integer> output line")
		}
		return float64(ns) / 1e6, elapsed
	}
	run()
	compute := make([]float64, runs)
	total := make([]float64, runs)
	for i := range compute {
		compute[i], total[i] = run()
		fmt.Printf("Run %d: internal %.6f ms | process %.3f ms\n", i+1, compute[i], total[i])
	}
	sort.Float64s(compute)
	sort.Float64s(total)
	fmt.Printf("Internal median: %.6f ms | Process median: %.3f ms\n", (compute[(runs-1)/2]+compute[runs/2])/2, (total[(runs-1)/2]+total[runs/2])/2)
	fmt.Println("Internal = user-marked timed region, process = startup + execution + captured stdout. One warm-up excluded. Use a noninteractive program with small output.")
}
