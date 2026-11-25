package main

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "os"
 "path/filepath"
 "runtime"
 "strconv"
 "strings"
 "time"
)

const updateAPI = "https://api.github.com/repos/Tripollas/Tripol/releases?per_page=30"
const updateTTL = 24*time.Hour
const updateLimit = 64<<20

type updateAsset struct {
 Name string `json:"name"`
 URL string `json:"browser_download_url"`
 Digest string `json:"digest"`
 Size int64 `json:"size"`
}
type updateRelease struct {
 Tag string `json:"tag_name"`
 Draft bool `json:"draft"`
 Preview bool `json:"prerelease"`
 Assets []updateAsset `json:"assets"`
}
type updateCache struct {
 Checked time.Time `json:"checked"`
 Release updateRelease `json:"release"`
}
func updateVersion(s string) ([3]int,bool) {
 var v [3]int
 p:=strings.Split(strings.TrimPrefix(s,"v"),".")
 if len(p)!=3 {return v,false}
 for i,x:=range p {
  if x=="" || (len(x)>1 && x[0]=='0') {return v,false}
  for _,c:=range x {if c<'0'||c>'9' {return v,false}}
  n,e:=strconv.Atoi(x);if e!=nil {return v,false};v[i]=n
 }
 return v,true
}
func newerVersion(a,b string) bool {
 x,ok:=updateVersion(a);y,valid:=updateVersion(b)
 if !ok||!valid {return false}
 for i:=range x {if x[i]!=y[i] {return x[i]>y[i]}}
 return false
}
func newestRelease(releases []updateRelease) updateRelease {
 var best updateRelease
 for _,r:=range releases {
  if _,ok:=updateVersion(r.Tag);!ok||r.Draft {continue}
  // TRIPOL is a public preview; numbered pre-releases are eligible.
  if best.Tag==""||newerVersion(r.Tag,best.Tag) {best=r}
 }
 return best
}
func updateRoot()(string,error){
 d,e:=os.UserCacheDir();if e!=nil{return "",e};return filepath.Join(d,"tripol","updates"),nil
}
func getUpdateRelease(client *http.Client,endpoint string)(updateRelease,error){
 req,e:=http.NewRequest("GET",endpoint,nil);if e!=nil{return updateRelease{},e}
 req.Header.Set("Accept","application/vnd.github+json")
 req.Header.Set("User-Agent","TRIPOL/"+version)
 resp,e:=client.Do(req);if e!=nil{return updateRelease{},e};defer resp.Body.Close()
 if resp.StatusCode!=200{return updateRelease{},fmt.Errorf("release server HTTP %d",resp.StatusCode)}
 var releases []updateRelease
 if e=json.NewDecoder(io.LimitReader(resp.Body,2<<20)).Decode(&releases);e!=nil{return updateRelease{},e}
 r:=newestRelease(releases);if r.Tag=="" {return r,fmt.Errorf("no published numbered release found")};return r,nil
}
func readUpdateCache(root string) updateCache {
 var c updateCache
 b,e:=os.ReadFile(filepath.Join(root,"check.json"));if e==nil&&len(b)<1<<20{_ = json.Unmarshal(b,&c)}
 return c
}
func writeUpdateCache(root string,c updateCache){
 if os.MkdirAll(root,0700)!=nil{return}
 b,e:=json.Marshal(c);if e!=nil{return}
 f,e:=os.CreateTemp(root,"check-");if e!=nil{return};name:=f.Name();defer os.Remove(name)
 _,e=f.Write(b);closeErr:=f.Close()
 if e==nil&&closeErr==nil{_ = os.Rename(name,filepath.Join(root,"check.json"))}
}
func updateLookup(force bool,timeout time.Duration)(updateRelease,error){
 root,e:=updateRoot();if e!=nil{return updateRelease{},e}
 c:=readUpdateCache(root);age:=time.Since(c.Checked)
 if !force && age>=0 && age<updateTTL {return c.Release,nil}
 r,e:=getUpdateRelease(&http.Client{Timeout:timeout},updateAPI)
 // Cache failed attempts too, so offline interactive invocations remain fast.
 if e!=nil {writeUpdateCache(root,updateCache{time.Now(),c.Release});return c.Release,e}
 writeUpdateCache(root,updateCache{time.Now(),r});return r,nil
}
func updateDisabled()bool{return os.Getenv("TRIPOL_NO_UPDATE")!=""||os.Getenv("CI")!=""}
func maybeUpdateNotice(args []string){
 if updateDisabled(){return}
 if len(args)>0 {
  switch args[0] {case "bench","bench-build","arena","check","ir","docs","libs","update","--version":return}
 }
 for _,a:=range args{if strings.HasPrefix(a,"--diagnostics"){return}}
 // Keep redirected output and automation free of notices.
 for _,f:=range []*os.File{os.Stdout,os.Stderr}{s,e:=f.Stat();if e!=nil||s.Mode()&os.ModeCharDevice==0{return}}
 r,_:=updateLookup(false,700*time.Millisecond)
 if newerVersion(r.Tag,version){fmt.Fprintf(os.Stderr,"TRIPOL update available: %s -> %s. Run: tripol update\n",version,r.Tag)}
}
func updateAssetFor(r updateRelease,platform string)(updateAsset,error){
 if _,ok:=updateVersion(r.Tag);!ok||!strings.HasPrefix(r.Tag,"v"){return updateAsset{},fmt.Errorf("invalid release tag")}
 p:=strings.Split(strings.TrimPrefix(r.Tag,"v"),".")
 // Current releases use major.minor in distribution names.
 short:=p[0]+"."+p[1]
 name:=""
 switch platform{case "windows":name="TRIPOL-v"+short+"-Setup-x64.exe";case "linux":name="TRIPOL-v"+short+"-Linux-x64.tar.gz";default:return updateAsset{},fmt.Errorf("unsupported platform")}
 for _,a:=range r.Assets{
  if a.Name!=name{continue}
  expected:="https://github.com/Tripollas/Tripol/releases/download/"+r.Tag+"/"+name
  d:=strings.TrimPrefix(a.Digest,"sha256:");raw,e:=hex.DecodeString(d)
  if a.URL!=expected||!strings.HasPrefix(a.Digest,"sha256:")||e!=nil||len(raw)!=32||a.Size<=0||a.Size>updateLimit{return updateAsset{},fmt.Errorf("invalid official asset metadata")}
  return a,nil
 }
 return updateAsset{},fmt.Errorf("official asset not found for %s",platform)
}
func downloadUpdate(client *http.Client,a updateAsset,root string)(string,error){
 if e:=os.MkdirAll(root,0700);e!=nil{return "",e}
 req,e:=http.NewRequest("GET",a.URL,nil);if e!=nil{return "",e}
 resp,e:=client.Do(req);if e!=nil{return "",e};defer resp.Body.Close()
 if resp.StatusCode!=200{return "",fmt.Errorf("download HTTP %d",resp.StatusCode)}
 f,e:=os.CreateTemp(root,"download-");if e!=nil{return "",e}
 name:=f.Name();defer os.Remove(name);defer f.Close()
 h:=sha256.New()
 n,e:=io.Copy(io.MultiWriter(f,h),io.LimitReader(resp.Body,a.Size+1))
 if e!=nil{return "",e}
 if n!=a.Size||"sha256:"+hex.EncodeToString(h.Sum(nil))!=strings.ToLower(a.Digest){return "",fmt.Errorf("SHA-256 or size mismatch; update rejected")}
 if e=f.Close();e!=nil{return "",e}
 // Use a fresh private directory, never overwrite an existing installer.
 dest,e:=os.MkdirTemp(root,"verified-");if e!=nil{return "",e}
 path:=filepath.Join(dest,a.Name)
 if e=os.Rename(name,path);e!=nil{_ = os.Remove(dest);return "",e}
 return path,nil
}
func updateCLI(args []string){
 if len(args)>1||(len(args)==1&&args[0]!="--check"){die("usage: tripol update [--check]")}
 r,e:=updateLookup(true,8*time.Second);if e!=nil{die("update check failed: %v",e)}
 if !newerVersion(r.Tag,version){fmt.Println("TRIPOL is up to date:",version);return}
 fmt.Printf("TRIPOL update available: %s -> %s\n",version,r.Tag)
 if len(args)==1{return}
 if runtime.GOARCH!="amd64"{die("updates require x86-64")}
 a,e:=updateAssetFor(r,runtime.GOOS);if e!=nil{die("%v",e)}
 root,e:=updateRoot();if e!=nil{die("%v",e)}
 client:=&http.Client{Timeout:2*time.Minute,CheckRedirect:func(req *http.Request,via []*http.Request)error{
  if len(via)>=8||req.URL.Scheme!="https"||(req.URL.Hostname()!="github.com"&&!strings.HasSuffix(req.URL.Hostname(),".githubusercontent.com")){return fmt.Errorf("untrusted release redirect")};return nil
 }}
 path,e:=downloadUpdate(client,a,root);if e!=nil{die("update download failed: %v",e)}
 fmt.Println("Verified SHA-256. Downloaded:",path)
 if runtime.GOOS=="windows"{fmt.Println("Close TRIPOL processes and run the unsigned setup EXE above to install the update.")}else{fmt.Println("Extract tripol from the archive and replace your installed compiler after closing running processes.")}
 fmt.Println("No files were installed automatically. Release:", "https://github.com/Tripollas/Tripol/releases/tag/"+r.Tag)
}
