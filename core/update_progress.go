package core

import (
 "encoding/json"
 "fmt"
 "os"
 "path/filepath"
 "time"
)

type updateProgress struct{ total,done int64; last int64; enabled bool }
func (p *updateProgress) Write(b []byte)(int,error){
 p.done+=int64(len(b))
 if p.enabled&&p.total>0 {
  pct:=p.done*100/p.total
  if pct>=p.last+10||p.done==p.total {fmt.Fprintf(os.Stderr,"\rDownloading: %d/%d bytes (%d%%)",p.done,p.total,pct);p.last=pct}
 }
 return len(b),nil
}
func interactiveUpdate()bool{
 s,e:=os.Stderr.Stat();return e==nil&&s.Mode()&os.ModeCharDevice!=0&&!updateDisabled()
}
type updateReceipt struct{Tag string `json:"tag"`; Name string `json:"asset"`; Digest string `json:"digest"`; Path string `json:"path"`; Downloaded time.Time `json:"downloaded"`}
func recordUpdate(root,tag string,a updateAsset,path string)error{
 // Each verified download has its own directory, avoiding concurrent history writes.
 b,e:=json.MarshalIndent(updateReceipt{tag,a.Name,a.Digest,path,time.Now().UTC()},"","  ");if e!=nil{return e}
 return os.WriteFile(filepath.Join(filepath.Dir(path),"receipt.json"),b,0600)
}
func updateHistoryCLI(){
 root,e:=updateRoot();if e!=nil{die("%v",e)}
 entries,e:=os.ReadDir(root);if os.IsNotExist(e){fmt.Println("No verified downloads.");return};if e!=nil{die("%v",e)}
 n:=0
 for _,d:=range entries{
  if !d.IsDir(){continue}
  b,e:=os.ReadFile(filepath.Join(root,d.Name(),"receipt.json"));if e!=nil||len(b)>16384{continue}
  var r updateReceipt;if json.Unmarshal(b,&r)!=nil{continue}
  fmt.Printf("%s %s %s\n",r.Downloaded.Format(time.RFC3339),r.Tag,r.Path);n++
 }
 if n==0{fmt.Println("No verified downloads.")}
}
