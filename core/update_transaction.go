package core

import (
 "archive/tar"
 "bufio"
 "compress/gzip"
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strings"
 "time"
)

// Only the documented per-user installation locations are eligible. Portable
// binaries, symlinked executables and system installations remain manual.
func managedUpdateTarget()(string,error){
 var target string
 switch runtime.GOOS {
 case "windows":
  base:=os.Getenv("LOCALAPPDATA");if base==""||!filepath.IsAbs(base){return "",fmt.Errorf("LOCALAPPDATA unavailable")}
  target=filepath.Join(base,"Programs","TRIPOL","tripol.exe")
 case "linux":
  home,e:=os.UserHomeDir();if e!=nil{return "",e};target=filepath.Join(home,".local","bin","tripol")
 default:return "",fmt.Errorf("unsupported update platform")
 }
 st,e:=os.Lstat(target);if e!=nil{return "",e}
 if !st.Mode().IsRegular(){return "",fmt.Errorf("managed compiler must be a regular file")}
 return target,nil
}
func currentManagedTarget()(string,error){
 target,e:=managedUpdateTarget();if e!=nil{return "",e}
 self,e:=os.Executable();if e!=nil{return "",e}
 a,e:=filepath.EvalSymlinks(self);if e!=nil{return "",e}
 b,e:=filepath.EvalSymlinks(target);if e!=nil{return "",e}
 equal:=a==b;if runtime.GOOS=="windows"{equal=strings.EqualFold(a,b)}
 if !equal{return "",fmt.Errorf("automatic replacement is limited to the official per-user installation")}
 return target,nil
}
func fileUpdateDigest(path string)(string,error){
 st,e:=os.Lstat(path);if e!=nil{return "",e}
 if !st.Mode().IsRegular()||st.Size()<=0||st.Size()>updateLimit{return "",fmt.Errorf("invalid update file")}
 f,e:=os.Open(path);if e!=nil{return "",e};defer f.Close()
 h:=sha256.New();n,e:=io.Copy(h,io.LimitReader(f,updateLimit+1))
 if e!=nil{return "",e};if n!=st.Size(){return "",fmt.Errorf("update file changed while reading")}
 return "sha256:"+hex.EncodeToString(h.Sum(nil)),nil
}
func copyUpdateFile(src,dst,digest string)error{
 got,e:=fileUpdateDigest(src);if e!=nil{return e};if got!=strings.ToLower(digest){return fmt.Errorf("update hash mismatch")}
 in,e:=os.Open(src);if e!=nil{return e};defer in.Close()
 out,e:=os.OpenFile(dst,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0700);if e!=nil{return e}
 h:=sha256.New();n,e:=io.Copy(io.MultiWriter(out,h),io.LimitReader(in,updateLimit+1))
 if e==nil{e=out.Sync()};ce:=out.Close();if e==nil{e=ce}
 if e==nil&&(n<=0||n>updateLimit||"sha256:"+hex.EncodeToString(h.Sum(nil))!=strings.ToLower(digest)){e=fmt.Errorf("copied update hash mismatch")}
 if e!=nil{_ = os.Remove(dst)}
 return e
}
func probeUpdateCompiler(path string)(string,error){
 ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
 cmd:=exec.CommandContext(ctx,path,"--version");cmd.Env=append(os.Environ(),"TRIPOL_NO_UPDATE=1")
 out,e:=cmd.CombinedOutput();if e!=nil{return "",fmt.Errorf("compiler version probe: %w",e)}
 s:=strings.TrimSpace(string(out));if !strings.HasPrefix(s,"TRIPOL "){return "",fmt.Errorf("invalid compiler version response")}
 v:=strings.TrimPrefix(s,"TRIPOL ");if _,ok:=updateVersion(v);!ok{return "",fmt.Errorf("invalid compiler version")}
 return "v"+v,nil
}
// Extract only the known regular compiler member. No archive-supplied path is
// ever used for output. Reject duplicate, oversized, special and unknown entries.
func extractUpdateLinux(archive,digest,destination string)(string,error){
 actual,e:=fileUpdateDigest(archive);if e!=nil{return "",e};if actual!=strings.ToLower(digest){return "",fmt.Errorf("archive hash mismatch")}
 in,e:=os.Open(archive);if e!=nil{return "",e};defer in.Close()
 gz,e:=gzip.NewReader(in);if e!=nil{return "",e};defer gz.Close()
 tr:=tar.NewReader(io.LimitReader(gz,updateLimit+1048576))
 var found bool
 defer func(){if !found{_ = os.Remove(destination)}}()
 seen:=map[string]bool{};var total int64
 for {
  h,e:=tr.Next();if e==io.EOF{break};if e!=nil{found=false;return "",e}
  if (h.Name!="tripol"&&h.Name!="LICENSE")||seen[h.Name]||h.Typeflag!=tar.TypeReg||h.Size<0||h.Size>updateLimit||total>updateLimit-h.Size{found=false;return "",fmt.Errorf("unexpected update archive member")}
  total+=h.Size;seen[h.Name]=true
  if h.Name=="tripol" {
   if h.Size==0{found=false;return "",fmt.Errorf("empty compiler")}
   out,e:=os.OpenFile(destination,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0700);if e!=nil{return "",e}
   _,e=io.CopyN(out,tr,h.Size);if e==nil{e=out.Sync()};ce:=out.Close();if e==nil{e=ce}
   if e!=nil{return "",e};found=true
  }
 }
 if !found{return "",fmt.Errorf("archive missing compiler")}
 return fileUpdateDigest(destination)
}
type updateBackup struct {
 Tag string `json:"tag"`
 Digest string `json:"digest"`
 Created time.Time `json:"created"`
}
func backupExeName()string{if runtime.GOOS=="windows"{return "previous.exe"};return "previous"}
func applyUpdateTransaction(target,source,digest,tag,expectedCurrent string,probe func(string)(string,error))error{
 lockPath:=target+".update.lock"
 lock,e:=os.OpenFile(lockPath,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0600);if e!=nil{return fmt.Errorf("another update may be running: %w",e)}
 _ = lock.Close();defer os.Remove(lockPath)
 current,e:=fileUpdateDigest(target);if e!=nil{return e}
 if current!=expectedCurrent{return fmt.Errorf("installed compiler changed during update")}
 oldTag,e:=probe(target);if e!=nil{return e}
 if _,ok:=updateVersion(tag);!ok{return fmt.Errorf("invalid target version")}
 dir,e:=os.MkdirTemp(filepath.Dir(target),".tripol-rollback-");if e!=nil{return e}
 retain:=false;defer func(){if !retain{_ = os.RemoveAll(dir)}}()
 staged:=filepath.Join(dir,"candidate");if runtime.GOOS=="windows"{staged+=".exe"}
 if e=copyUpdateFile(source,staged,digest);e!=nil{return e}
 got,e:=probe(staged);if e!=nil||got!=tag{return fmt.Errorf("candidate version mismatch: %q %v",got,e)}
 backup:=filepath.Join(dir,backupExeName())
 manifest,e:=json.Marshal(updateBackup{oldTag,current,time.Now().UTC()});if e!=nil{return e}
 // Save recovery metadata before moving the original. A receipt is eligible
 // only if the previous binary exists and still matches its digest.
 if e=os.WriteFile(filepath.Join(dir,"receipt.json"),manifest,0600);e!=nil{return e}
 if e=os.Rename(target,backup);e!=nil{return fmt.Errorf("cannot move installed compiler; close other TRIPOL processes: %w",e)}
 retain=true
 restore:=func(cause error)error{
  if e:=os.Remove(target);e!=nil&&!os.IsNotExist(e){return fmt.Errorf("%v; restore blocked (%v); original preserved at %s",cause,e,backup)}
  if e:=os.Rename(backup,target);e!=nil{return fmt.Errorf("%v; restore failed (%v); original preserved at %s",cause,e,backup)}
  retain=false;return fmt.Errorf("%w; original restored",cause)
 }
 if e=os.Rename(staged,target);e!=nil{return restore(e)}
 got,e=probe(target)
 if e!=nil||got!=tag{return restore(fmt.Errorf("installed version verification failed: %q %v",got,e))}
 return nil
}
func latestUpdateBackup(target string)(string,updateBackup,error){
 entries,e:=os.ReadDir(filepath.Dir(target));if e!=nil{return "",updateBackup{},e}
 var best updateBackup;var path string
 for _,d:=range entries{
  if !d.IsDir()||!strings.HasPrefix(d.Name(),".tripol-rollback-"){continue}
  dir:=filepath.Join(filepath.Dir(target),d.Name())
  st,e:=os.Lstat(filepath.Join(dir,"receipt.json"));if e!=nil||!st.Mode().IsRegular()||st.Size()>4096{continue}
  b,e:=os.ReadFile(filepath.Join(dir,"receipt.json"));if e!=nil{continue}
  var receipt updateBackup;if json.Unmarshal(b,&receipt)!=nil{continue}
  if _,ok:=updateVersion(receipt.Tag);!ok{continue}
  p:=filepath.Join(dir,backupExeName());digest,e:=fileUpdateDigest(p)
  if e!=nil||digest!=receipt.Digest{continue}
  if path==""||receipt.Created.After(best.Created){best=receipt;path=p}
 }
 if path==""{return "",best,fmt.Errorf("no intact local rollback backup")}
 return path,best,nil
}
func updateConfirm(prompt string)bool{
 st,e:=os.Stdin.Stat();if e!=nil||st.Mode()&os.ModeCharDevice==0{return false}
 fmt.Print(prompt+" [y/N]: ")
 line,_:=bufio.NewReader(os.Stdin).ReadString('\n')
 line=strings.ToLower(strings.TrimSpace(line));return line=="y"||line=="yes"||line=="s"||line=="si"
}
func launchManagedUpdate(target,source,digest,tag string)error{
 old,e:=fileUpdateDigest(target);if e!=nil{return e}
 if runtime.GOOS=="windows"{return launchUpdateWorker(target,source,digest,tag,old)}
 if e=applyUpdateTransaction(target,source,digest,tag,old,probeUpdateCompiler);e!=nil{return e}
 fmt.Println("Installed",tag,"with a verified rollback backup.")
 return nil
}
func rollbackUpdate(){
 target,e:=currentManagedTarget();if e!=nil{die("%v",e)}
 source,backup,e:=latestUpdateBackup(target);if e!=nil{die("%v",e)}
 if !updateConfirm("Restore "+backup.Tag+" at "+target+"?"){fmt.Println("Rollback cancelled.");return}
 if e=launchManagedUpdate(target,source,backup.Digest,backup.Tag);e!=nil{die("%v",e)}
}
func withinUpdatePath(path,root string)bool{
 a,e:=filepath.Abs(path);if e!=nil{return false};b,e:=filepath.Abs(root);if e!=nil{return false}
 r,e:=filepath.Rel(b,a);return e==nil&&r!="."&&r!=".."&&!strings.HasPrefix(r,".."+string(filepath.Separator))&&!filepath.IsAbs(r)
}
func updateWorkerSourceAllowed(source,target string)bool{
 root,e:=updateRoot();if e==nil&&withinUpdatePath(source,root){return true}
 parent:=filepath.Dir(source)
 return filepath.Dir(parent)==filepath.Dir(target)&&strings.HasPrefix(filepath.Base(parent),".tripol-rollback-")&&filepath.Base(source)==backupExeName()
}
var errUpdateWorkerUnsupported=errors.New("update worker is only supported on Windows")
