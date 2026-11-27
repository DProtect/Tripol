package core

import(
 "bytes"
 "fmt"
 "reflect"
 "strings"
 "testing"
)

func TestCFGParallelDeterminism(t *testing.T){
 var source strings.Builder
 for i:=0;i<96;i++{fmt.Fprintf(&source,"func f%d(flag:bool)->int {x=1;if flag {x=%d} else {x=%d};return x*2}\n",i,i,i)}
 source.WriteString("show f0(true),f95(false)")
 fs,ss:=parse(source.String())
 a,mainA,statsA:=optimizeCFGWorkers(fs,ss,1)
 for _,workers:=range []int{2,4}{
  b,mainB,statsB:=optimizeCFGWorkers(fs,ss,workers)
  if statsA!=statsB||!reflect.DeepEqual(a,b)||!reflect.DeepEqual(mainA,mainB){t.Fatal("worker-dependent CFG result",workers)}
  for _,target:=range []string{"linux","windows"}{
   serial:=generate(a,mainA,target,true,reachable(a,mainA),true,true)
   parallel:=generate(b,mainB,target,true,reachable(b,mainB),true,true)
   if !bytes.Equal(serial,parallel){t.Fatal("worker-dependent binary",target,workers)}
  }
 }
 t.Setenv("TRIPOL_JOBS","1");if cfgWorkerCount(96)!=1{t.Fatal("serial override")}
 t.Setenv("TRIPOL_JOBS","4");if cfgWorkerCount(96)!=4||cfgWorkerCount(12)!=1{t.Fatal("worker cap/threshold")}
}
func BenchmarkCFGWorkers(b *testing.B){
 var source strings.Builder
 for i:=0;i<300;i++{fmt.Fprintf(&source,"func f%d(flag:bool)->int {x=1;if flag {x=%d} else {x=%d};return x*2}\n",i,i,i)}
 source.WriteString("show f0(true)")
 fs,ss:=parse(source.String())
 for _,workers:=range []int{1,2,4}{b.Run(fmt.Sprint(workers),func(b *testing.B){b.ReportAllocs();for i:=0;i<b.N;i++{optimizeCFGWorkers(fs,ss,workers)}})}
}
