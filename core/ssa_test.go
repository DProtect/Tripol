package core

import "testing"

func TestQuantumXSSA(t *testing.T){
 fs,_:=parse("func f(n:int)->int {x=n+1;x=x*2;dead=x/0;return x}\nshow 0")
 p,ok:=buildScalarSSA(fs[0]);if !ok||p.result<0{t.Fatal(p)}
 found:=false
 for i,v:=range p.values{for _,a:=range v.args{if a>=i{t.Fatal("non-dominating definition")}};if v.op=="div.checked.i64"{found=true;if v.uses!=0{t.Fatal("dead operation unexpectedly used")}}}
 if !found{t.Fatal("checked dead division removed")}
}
