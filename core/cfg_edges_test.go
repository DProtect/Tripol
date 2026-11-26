package core

import "testing"

func TestCFGComparisonEdgeFacts(t *testing.T){
 for _,op:=range []string{"==","!="}{
  yes,no:=scalarFacts{},scalarFacts{}
  cond:=&expr{op:"bin",s:op,args:[]*expr{{op:"var",s:"x"},{op:"int",n:7}}}
  refineCondition(cond,yes,no)
  correct,other:=yes,no
  if op=="!="{correct,other=no,yes}
  if correct["x"]!=(scalarFact{"int",7})||len(other)!=0{t.Fatal(op,yes,no)}
 }
 yes,no:=scalarFacts{},scalarFacts{}
 refineCondition(&expr{op:"var",s:"flag"},yes,no)
 if yes["flag"]!=(scalarFact{"bool",1})||no["flag"]!=(scalarFact{"bool",0}){t.Fatal(yes,no)}
}
