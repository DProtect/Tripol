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
		got:=propagateBlock(body,scalarFacts{},&cfgConstantStats{})
		if len(got)!=1||got[0].op!=op{t.Fatal("dead tail retained",op)}
	}
}
