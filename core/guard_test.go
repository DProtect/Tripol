package core

import "testing"

func TestGuardNested(t *testing.T) {
	runCase(t, `func fail(n) {return 10/n}
try {
 show fail(0)
} catch code {
 show "caught",code
} finally {show "cleanup"}
try {show 7} catch e {show e} finally {show 8}
try {
 try {throw 101} catch e {show e;throw 102} finally {show 9}
} catch e {show e} finally {show 10}
show "done"`, "caught 10\ncleanup\n7\n8\n101\n9\n102\n10\ndone\n", false)
}
func TestGuardBounds(t *testing.T) {
	runCase(t, `a=[1,2]
try {show a[3]} catch e {show e} finally {show 99}
try {assert false,"bad"} catch e {show e}
repeat 5 {try {throw 123} catch e {show e}}`, "11\n99\n12\n123\n123\n123\n123\n123\n", false)
}
func TestGuardFatal(t *testing.T) {
	runCase(t, `try {throw 101} catch e {throw 102} finally {show 5}`, "5\nTRIPOL runtime error: unhandled exception\n", true)
}
