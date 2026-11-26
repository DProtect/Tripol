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

func TestCFGWideSignedComparisonImmediate(t *testing.T){
 // -4294967296 would become zero if truncated to a signed imm32.
 fs,ss:=parse("func f(x:int)->int {k=-4294967296;if x==k {return 1};return 0};show f(0),f(-4294967296)")
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  for _,target:=range []string{"linux","windows"}{
   if len(compileParsed(fs,ss,target,profile))==0{t.Fatal("empty image",target,profile)}
  }
 }
}
