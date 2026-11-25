package core

import (
 "crypto/sha256"
 "encoding/json"
 "fmt"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "testing"
)
func TestQuantumXUpdateReceipt(t *testing.T){
 root:=t.TempDir();dir:=filepath.Join(root,"verified-example");if e:=os.Mkdir(dir,0700);e!=nil{t.Fatal(e)}
 a:=updateAsset{Name:"tripol.test",Digest:fmt.Sprintf("sha256:%x",sha256.Sum256([]byte("abc")))}
 path:=filepath.Join(dir,a.Name)
 if e:=recordUpdate(root,"v0.28.0",a,path);e!=nil{t.Fatal(e)}
 b,e:=os.ReadFile(filepath.Join(dir,"receipt.json"));if e!=nil{t.Fatal(e)}
 var r updateReceipt;if e=json.Unmarshal(b,&r);e!=nil||r.Tag!="v0.28.0"||r.Path!=path||r.Digest!=a.Digest{t.Fatal(r,e)}
}
func TestQuantumXUpdateTruncated(t *testing.T){
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){fmt.Fprint(w,"short")}));defer server.Close()
 a:=updateAsset{Name:"tripol.test",URL:server.URL,Size:10,Digest:fmt.Sprintf("sha256:%x",sha256.Sum256([]byte("0123456789")))}
 root:=t.TempDir();if _,e:=downloadUpdate(server.Client(),a,root);e==nil{t.Fatal("accepted truncated update")}
 entries,e:=os.ReadDir(root);if e!=nil||len(entries)!=0{t.Fatal("partial download remained",entries,e)}
}
