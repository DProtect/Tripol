//go:build windows

package core

import (
 "bufio"
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strconv"
 "syscall"
 "time"
)

func launchUpdateWorker(target,source,digest,tag,old string)error{
 root,e:=updateRoot();if e!=nil{return e};if e=os.MkdirAll(root,0700);e!=nil{return e}
 dir,e:=os.MkdirTemp(root,"worker-");if e!=nil{return e}
 helper:=filepath.Join(dir,"worker.exe")
 if e=copyUpdateFile(target,helper,old);e!=nil{_ = os.RemoveAll(dir);return e}
 logPath:=filepath.Join(dir,"result.txt");log,e:=os.OpenFile(logPath,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0600);if e!=nil{return e}
 defer log.Close()
 cmd:=exec.Command(helper,"update","--worker",strconv.Itoa(os.Getpid()),source,digest,tag,old)
 cmd.Stderr=log
 pipe,e:=cmd.StdoutPipe();if e!=nil{return e}
 if e=cmd.Start();e!=nil{return e}
 ready:=make(chan string,1)
 go func(){scanner:=bufio.NewScanner(pipe);if scanner.Scan(){ready<-scanner.Text()}else{ready<-""}}()
 select{
 case line:=<-ready:
  if line!="READY"{_ = cmd.Process.Kill();_ = cmd.Wait();return fmt.Errorf("update worker failed; see %s",logPath)}
 case <-time.After(8*time.Second):
  _ = cmd.Process.Kill();_ = cmd.Wait();return fmt.Errorf("update worker did not initialize; see %s",logPath)
 }
 _ = pipe.Close();_ = cmd.Process.Release()
 fmt.Println("Verified update worker is ready; it will replace the compiler after this process exits.")
 fmt.Println("Update result:",logPath)
 return nil
}
func updateWorkerCLI(args []string){
 if len(args)!=5{die("invalid update worker arguments")}
 pid,e:=strconv.ParseUint(args[0],10,32);if e!=nil||pid==0||int(pid)==os.Getpid(){die("invalid parent process")}
 target,e:=managedUpdateTarget();if e!=nil{die("%v",e)}
 root,e:=updateRoot();if e!=nil{die("%v",e)}
 self,e:=os.Executable();if e!=nil||!withinUpdatePath(self,root)||filepath.Base(self)!="worker.exe"{die("invalid update worker location")}
 if !updateWorkerSourceAllowed(args[1],target){die("invalid update source")}
 handle,e:=syscall.OpenProcess(syscall.SYNCHRONIZE,false,uint32(pid));if e!=nil{die("cannot open update parent: %v",e)}
 defer syscall.CloseHandle(handle)
 fmt.Println("READY")
 event,e:=syscall.WaitForSingleObject(handle,30000)
 if e!=nil||event!=syscall.WAIT_OBJECT_0{die("parent did not exit; compiler unchanged")}
 if e=applyUpdateTransaction(target,args[1],args[2],args[3],args[4],probeUpdateCompiler);e!=nil{die("%v",e)}
 fmt.Fprintln(os.Stderr,"[OK] Installed",args[3],"with a verified rollback backup.")
}
