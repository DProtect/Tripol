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

func TestCFGSignedImmediateProfiles(t *testing.T) {
	source := `func signed(x:int)->int {
  short=-128
  medium=-129
  min32=-2147483648
  wide=-4294967296
  show x+short,x+medium,x-medium,x*medium
  show x+min32,x+wide,x-wide,x*wide,x==wide,x>wide
  return x
 }
 func checked(x:int)->int {
  k=-129
  show 7
  return x+k
 }
 signed(2)
 signed(0)
 try {show checked(-9223372036854775807)} catch code {show code}
 `
	const want = "-126 -127 131 -258\n-2147483646 -4294967294 4294967298 -8589934592 false true\n-128 -129 129 0\n-2147483648 -4294967296 4294967296 0 false true\n7\n10\n"
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			got, failed := hyperRun(t, source, profile)
			if failed || got != want {
				t.Fatalf("got %q (failed=%v), want %q", got, failed, want)
			}
		})
	}
}
