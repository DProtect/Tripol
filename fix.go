package main

import "strings"

// Guard Fix: conservative edit-distance suggestion, never changes source.
func (c *compiler) suggest(msg string) string {
	var target string
	candidates := map[string]bool{}
	if strings.HasPrefix(msg, "unknown variable ") {
		target = strings.TrimPrefix(msg, "unknown variable ")
		for _, scope := range c.scopes {
			for n := range scope {
				candidates[n] = true
			}
		}
	}
	if strings.HasPrefix(msg, "unknown function ") {
		target = strings.TrimPrefix(msg, "unknown function ")
		for n := range c.fs {
			candidates[n] = true
		}
		for _, n := range []string{"input", "input_int", "parse_int", "len", "sqrt", "sum", "min", "max", "abs", "text", "read_text", "write_text", "append_text"} {
			candidates[n] = true
		}
	}
	if len(target) < 3 || len(target) > 128 {
		return ""
	}
	best, dist, tie := "", 3, false
	for n := range candidates {
		if strings.HasPrefix(n, "@") || len(n) > 128 {
			continue
		}
		d := editDistance(target, n)
		if d < dist {
			best, dist, tie = n, d, false
		} else if d == dist {
			tie = true
		}
	}
	if tie || dist > 2 {
		return ""
	}
	return best
}
func editDistance(a, b string) int {
	row := make([]int, len(b)+1)
	for j := range row {
		row[j] = j
	}
	for i := range a {
		prev := row[0]
		row[0] = i + 1
		for j := range b {
			old := row[j+1]
			cost := 0
			if a[i] != b[j] {
				cost = 1
			}
			row[j+1] = min(prev+cost, row[j]+1, old+1)
			prev = old
		}
	}
	return row[len(b)]
}
