package core

import (
	"reflect"
	"testing"
)

func TestCFGPruneAfterTerminator(t *testing.T) {
	fs, ss := parse("func f(x:int)->int {return x+1;show 999};show f(2)")
	validateCode(fs,ss,"linux")
	before := append([]stmt(nil),fs[len(fs)-1].body...)
	out, _, _ := optimizeCFGWorkers(fs,ss,1)
	if len(out[len(out)-1].body)!=1 || out[len(out)-1].body[0].op!="return" {t.Fatal("dead tail retained")}
	if !reflect.DeepEqual(fs[len(fs)-1].body,before){t.Fatal("original AST changed")}
	if out[len(out)-1].body[0].e.op!="bin"{t.Fatal("checked return expression lost")}
}

func TestCFGPruneLoopControl(t *testing.T) {
	for _, op:=range []string{"break","continue"}{
		body:=[]stmt{{op:op},{op:"show",e:&expr{op:"int",n:999}}}
		got,fallsThrough:=propagateBlock(body,scalarFacts{},&cfgConstantStats{})
		if len(got)!=1||got[0].op!=op||fallsThrough{t.Fatal("loop control has a fallthrough edge",op)}
	}
}

func TestCFGJoinIgnoresTerminatingBranch(t *testing.T) {
	const source = `func choose(n:int)->int {
		value=3
		if n<0 {
			value=99
			if n==-1 {return value}
			return 44
		} else {value=5}
		return value+1
	}
	show choose(-1),choose(-2),choose(1)`
	fs, ss := parse(source)
	validateCode(fs, ss, "linux")
	optimized, _, _ := optimizeCFGWorkers(fs, ss, 1)
	body := optimized[len(optimized)-1].body
	if typ, n, ok := fold(body[len(body)-1].e); !ok || typ != "int" || n != 6 {
		t.Fatal("reachable predecessor's value did not reach the join")
	}
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		got, failed := hyperRun(t, source, profile)
		if failed || got != "99 44 6\n" {
			t.Fatalf("%s: got %q (failed=%v)", profile, got, failed)
		}
	}
}
