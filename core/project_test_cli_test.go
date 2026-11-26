package core
import(
 "context"
 "os"
 "path/filepath"
 "runtime"
 "testing"
)
func TestSimpleProjectRunner(t *testing.T){
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="windows"&&runtime.GOOS!="linux"){t.Skip("native target required")}
 dir:=t.TempDir()
 for name,src:=range map[string]string{"a.test.tri":"assert 2+2==4","b.test.tri":"assert false","main.tri":"assert false"}{
  if e:=os.WriteFile(filepath.Join(dir,name),[]byte(src),0600);e!=nil{t.Fatal(e)}
 }
 files,e:=projectTestFiles(dir);if e!=nil||len(files)!=2{t.Fatal(files,e)}
 if _,e=runProjectTest(context.Background(),files[0],runtime.GOOS);e!=nil{t.Fatal(e)}
 if _,e=runProjectTest(context.Background(),files[1],runtime.GOOS);e==nil{t.Fatal("failed assertion accepted")}
}
