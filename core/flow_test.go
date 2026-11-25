package core

import "testing"

func TestQuantumXFlow(t *testing.T){
 g,ok:=inspectFlow([]stmt{{op:"decl"},{op:"while",yes:[]stmt{{op:"if",yes:[]stmt{{op:"break"}},no:[]stmt{{op:"continue"}}}}},{op:"return"}})
 if !ok||g.entry==g.exit{t.Fatal(g)}
 found:=map[string]bool{}
 for i,b:=range g.blocks{found[b.op]=true;for _,n:=range b.next{if n<0||n>=len(g.blocks){t.Fatal(i,n)}}}
 for _,op:=range []string{"while.condition","if","break","continue","return"}{if !found[op]{t.Fatal(op)}}
 if _,ok:=inspectFlow([]stmt{{op:"try"}});ok{t.Fatal("protected flow must remain on existing backend")}
}
