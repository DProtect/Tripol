package core

import (
 "crypto/sha256"
 "fmt"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
 "time"
)
func TestOverdriveVersions(t *testing.T){
 for _,tc:=range []struct{a,b string;want bool}{{"v0.26.0","0.25.0",true},{"v0.9.0","0.25.0",false},{"v0.26.0","0.26.0",false},{"v1.0.0","0.99.0",true},{"v0.26.1","0.26.0",true},{"v0.27.0-rc","0.26.0",false},{"v0.026.0","0.25.0",false}}{if newerVersion(tc.a,tc.b)!=tc.want{t.Fatalf("%s %s",tc.a,tc.b)}}
 r:=newestRelease([]updateRelease{{Tag:"v0.29.0",Draft:true},{Tag:"v0.25.0"},{Tag:"v0.26.0",Preview:true}})
 if r.Tag!="v0.26.0"{t.Fatal(r)}
}
func TestOverdriveUpdateHTTP(t *testing.T){
 srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){fmt.Fprint(w,`[{"tag_name":"v0.27.0","prerelease":true},{"tag_name":"v9.9.9","draft":true}]`)}));defer srv.Close()
 r,e:=getUpdateRelease(srv.Client(),srv.URL);if e!=nil||r.Tag!="v0.27.0"{t.Fatal(r,e)}
 root:=t.TempDir();writeUpdateCache(root,updateCache{time.Now(),r});if readUpdateCache(root).Release.Tag!=r.Tag{t.Fatal("cache")}
}
func TestOverdriveVerifiedDownload(t *testing.T){
 payload:="test asset"
 digest:=fmt.Sprintf("sha256:%x",sha256.Sum256([]byte(payload)))
 srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){fmt.Fprint(w,payload)}));defer srv.Close()
 a:=updateAsset{Name:"test.bin",URL:srv.URL,Digest:digest,Size:int64(len(payload))}
 root:=t.TempDir();p,e:=downloadUpdate(srv.Client(),a,root);if e!=nil{t.Fatal(e)}
 data,e:=os.ReadFile(p);if e!=nil||string(data)!=payload||filepath.Base(p)!=a.Name{t.Fatal(p,e)}
 a.Digest="sha256:"+strings.Repeat("0",64);if _,e=downloadUpdate(srv.Client(),a,root);e==nil{t.Fatal("accepted tampered asset")}
 a.Name="TRIPOL-v0.27-Setup-x64.exe";a.URL="https://example.com/evil";a.Digest=digest
 if _,e=updateAssetFor(updateRelease{Tag:"v0.27.0",Assets:[]updateAsset{a}},"windows");e==nil{t.Fatal("accepted external asset")}
}
