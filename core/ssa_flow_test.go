package core

import (
	"fmt"
	"strings"
	"testing"
)

func TestFlowSSAProfiles(t *testing.T) {
	src := `func run(n:int,swap:bool)->int {
		a=1
		b=2
		c=3
		d=4
		e=5
		f=6
		i=0
		while i<n {
			i+=1
			if i==2 {continue}
			if i==5 {break}
			if swap {tmp=a;a=b;b=tmp} else {a+=i}
			for j in range(-1,2) {
				if j==0 {continue}
				c+=a+j
				if i>3 {break}
				d+=b-j
				e+=d
				f+=e
			}
		}
		return a+10*b+100*c+d+e+f
	}
	func safe(x:int,guard:bool)->bool {
		ok=guard or 4/x>0
		if ok and not guard {return true}
		return ok
	}
	show safe(0,true),safe(-1,false)
	try {safe(0,false)} catch code {show code}
	for i in range(-2,10) {show run(i,true),run(i,false)}
	`
	want, bad := hyperRun(t, src, "debug")
	if bad {
		t.Fatal(want)
	}
	for _, profile := range []string{"nano", "balanced", "turbo"} {
		got, failed := hyperRun(t, src, profile)
		if failed || got != want {
			t.Fatalf("%s got %q, want %q", profile, got, want)
		}
	}
	fs, ss := parse(src)
	for _, target := range []string{"linux", "windows"} {
		_, report := measuredCompilation(fs, ss, target, "turbo")
		if report.FlowSSAFunctions != 2 {
			t.Fatalf("%s lowered %d functions", target, report.FlowSSAFunctions)
		}
	}
}

func TestFlowSSARangeAndTraps(t *testing.T) {
	source := `func traverse(n:int,k:int)->int {
		total=0
		for i in range(n,k) {
			if i!=0 {total+=i} else {bad=1/i}
		}
		return total
	}
	show traverse(9223372036854775806,9223372036854775807)
	try {show traverse(-1,2)} catch code {show code}
	try {show traverse(9223372036854775805,9223372036854775807)} catch code {show code}
	`
	const want = "9223372036854775806\n10\n10\n"
	for _, profile := range []string{"debug", "turbo"} {
		got, bad := hyperRun(t, source, profile)
		if bad || got != want {
			t.Fatalf("%s %q", profile, got)
		}
	}
}

func TestFlowSSAFallback(t *testing.T) {
	fs, ss := parse("func dependent(x:int)->int {if x>0 {return x};return x+1};func caller(x:int)->int {if x==0 {return dependent(x)};return x};show caller(0)")
	_, report := measuredCompilation(fs, ss, "linux", "turbo")
	if report.FlowSSAFunctions != 1 {
		t.Fatal(report.FlowSSAFunctions)
	}
}

func TestFlowSSAParallelCopies(t *testing.T) {
	var source strings.Builder
	fmt.Fprintln(&source, "func rotate(n:int)->int {")
	for i := range 10 {
		fmt.Fprintf(&source, "v%d=%d\n", i, i)
	}
	fmt.Fprintln(&source, "repeat n {first=v0")
	for i := range 9 {
		fmt.Fprintf(&source, "v%d=v%d\n", i, i+1)
	}
	fmt.Fprintln(&source, "v9=first};return v0+10*v1+100*v2+1000*v3+10000*v4+100000*v5+1000000*v6+10000000*v7+100000000*v8+1000000000*v9}")
	fmt.Fprintln(&source, "for i in range(16){show rotate(i)}")
	want, failed := hyperRun(t, source.String(), "debug")
	got, bad := hyperRun(t, source.String(), "turbo")
	if failed || bad || got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
