package core

import "slices"

// Reserve the complete literal tail once. Preserve the existing label order,
// exact bytes, and NUL terminators; no pointer identity or layout changes.
func appendLinkLiterals(a *asm) {
	total := 0
	maxInt := int(^uint(0) >> 1)
	for _, s := range a.dataOrder {
		if len(s) >= maxInt-total || total+len(s)+1 > maxInt-len(a.b) {
			die("internal literal section too large")
		}
		total += len(s)+1
	}
	a.b = slices.Grow(a.b,total)
	for _, s := range a.dataOrder {
		a.label(a.data[s])
		a.b = append(a.b,s...)
		a.b = append(a.b,0)
	}
}
