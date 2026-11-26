package core
import (
 "bytes"
 "fmt"
 "crypto/sha256"
 "testing"
)
func TestCompileReportMatchesBuild(t *testing.T){
 fs,ss:=parse(`func dead()->int {return 9}
 func used(x:int)->int {y=x+1;return y*2}
 show used(20)`)
 for _,target:=range []string{"windows","linux"}{
  for _,profile:=range []string{"nano","balanced","turbo","debug"}{
   got,r:=measuredCompilation(fs,ss,target,profile)
   want:=compileParsed(fs,ss,target,profile)
   if !bytes.Equal(got,want)||r.ImageBytes!=len(want)||r.SHA256!=fmt.Sprintf("%x",sha256.Sum256(want)){t.Fatal(target,profile,r)}
   if r.Functions!=2{t.Fatal(r)}
   if profile=="debug"{if r.EmittedFunctions!=2||r.ValidationBufferBytes!=0{t.Fatal(r)}}else{if r.EmittedFunctions!=1||r.ValidationBufferBytes==0{t.Fatal(r)}}
  }
 }
}
