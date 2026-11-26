package core
import(
 "bytes"
 "runtime"
 "strings"
 "testing"
)
func TestSimpleShell(t *testing.T){
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="windows"&&runtime.GOOS!="linux"){t.Skip("native")}
 var out bytes.Buffer
 shellSession(strings.NewReader("name=\"Alex\"\nshow f\"Hola {name}\"\n25*4\nx=7\nx=1/0\nshow x\nwrite_text(\"x\",\"no\")\n:reset\nshow x\n:exit\n"),&out,runtime.GOOS)
 for _,v:=range []string{"Hola Alex","100","7","only pure scalar expressions","Session reset.","unknown variable x"}{if !strings.Contains(out.String(),v){t.Fatalf("missing %q: %s",v,out.String())}}
}
