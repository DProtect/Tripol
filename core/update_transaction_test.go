package core

import (
 "archive/tar"
 "compress/gzip"
 "errors"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestUpdateTransactionAndRollback(t *testing.T){
 dir:=t.TempDir();target:=filepath.Join(dir,"tripol");candidate:=filepath.Join(dir,"candidate")
 if e:=os.WriteFile(target,[]byte("v0.27.0"),0700);e!=nil{t.Fatal(e)}
 if e:=os.WriteFile(candidate,[]byte("v0.27.1"),0700);e!=nil{t.Fatal(e)}
 probe:=func(p string)(string,error){b,e:=os.ReadFile(p);return string(b),e}
 old,_:=fileUpdateDigest(target);digest,_:=fileUpdateDigest(candidate)
 if e:=applyUpdateTransaction(target,candidate,digest,"v0.27.1",old,probe);e!=nil{t.Fatal(e)}
 if got,_:=probe(target);got!="v0.27.1"{t.Fatal(got)}
 backup,receipt,e:=latestUpdateBackup(target);if e!=nil||receipt.Tag!="v0.27.0"{t.Fatal(receipt,e)}
 current,_:=fileUpdateDigest(target)
 if e=applyUpdateTransaction(target,backup,receipt.Digest,receipt.Tag,current,probe);e!=nil{t.Fatal(e)}
 if got,_:=probe(target);got!="v0.27.0"{t.Fatal(got)}
 // Simulate a failure after rename: old binary must be restored.
 old,_=fileUpdateDigest(target)
 fail:=func(p string)(string,error){v,e:=probe(p);if p==target&&v=="v0.27.1"{return "",errors.New("simulated probe failure")};return v,e}
 if e=applyUpdateTransaction(target,candidate,digest,"v0.27.1",old,fail);e==nil||!strings.Contains(e.Error(),"original restored"){t.Fatal(e)}
 if got,_:=probe(target);got!="v0.27.0"{t.Fatal(got)}
 if e=applyUpdateTransaction(target,candidate,"sha256:"+strings.Repeat("0",64),"v0.27.1",old,probe);e==nil{t.Fatal("tampered candidate accepted")}
}
func TestUpdateArchivePolicy(t *testing.T){
 for _,name:=range []string{"tripol","../tripol","/tripol","other"}{
  t.Run(name,func(t *testing.T){
   dir:=t.TempDir();archive:=filepath.Join(dir,"archive.tar.gz");f,e:=os.Create(archive);if e!=nil{t.Fatal(e)}
   gz:=gzip.NewWriter(f);tw:=tar.NewWriter(gz)
   payload:=[]byte("test compiler")
   if e=tw.WriteHeader(&tar.Header{Name:name,Mode:0755,Size:int64(len(payload)),Typeflag:tar.TypeReg});e!=nil{t.Fatal(e)}
   if _,e=tw.Write(payload);e!=nil{t.Fatal(e)};if e=tw.Close();e!=nil{t.Fatal(e)};if e=gz.Close();e!=nil{t.Fatal(e)};if e=f.Close();e!=nil{t.Fatal(e)}
   digest,_:=fileUpdateDigest(archive);dest:=filepath.Join(dir,"candidate")
   _,e=extractUpdateLinux(archive,digest,dest)
   if name=="tripol"{if e!=nil{t.Fatal(e)};b,e:=os.ReadFile(dest);if e!=nil||string(b)!=string(payload){t.Fatal(e)}}else{if e==nil{t.Fatal("unsafe archive accepted")};if _,e=os.Stat(dest);!os.IsNotExist(e){t.Fatal("unexpected output")}}
  })
 }
}
