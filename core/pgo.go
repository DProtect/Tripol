package core

import (
 "bytes"
 "context"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "io"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "sort"
 "strings"
 "time"
)

const pgoPrefix = "TRIPOL_PGO_R1:"
const pgoMaxLog = 8 << 20
type pgoProfile struct {
 Key string `json:"key"`
 Target string `json:"target"`
 Counts map[string]uint64 `json:"function_entries"`
}

// Training uses the ordinary checked stack ABI (no register leaves/inlining).
// Emit an observation before the first body instruction, without modifying
// parameters. Self tail jumps do not count as new physical function entries.
func (a *asm) emitPGOEntry(name string){
 marker:=pgoPrefix+hex.EncodeToString([]byte(name))+"\n"
 a.str(marker)
 a.emit(0x48,0x89,0xc6,0xba)
 a.u32(len(marker))
 a.call("writeerr")
}

func pgoOrder(fs []fn,counts map[string]uint64)[]fn{
 if len(counts)==0{return fs}
 out:=append([]fn(nil),fs...)
 sort.SliceStable(out,func(i,j int)bool{return counts[out[i].name]>counts[out[j].name]})
 return out
}

func pgoKey(fs []fn,ss []stmt,target string)string{
 digest,ok:=cacheCompilerDigest();if !ok{die("cannot identify compiler for PGO")}
 key,ok:=buildCacheKey(digest,fs,ss,target,"pgo-r1");if !ok{die("unsupported project shape for PGO")}
 return key
}

func readPGO(path string,fs []fn,ss []stmt,target string)map[string]uint64{
 f,e:=os.Open(path);if e!=nil{die("%v",e)};defer f.Close()
 b,e:=io.ReadAll(io.LimitReader(f,pgoMaxLog+1));if e!=nil||len(b)>pgoMaxLog{die("invalid or oversized PGO profile")}
 var p pgoProfile
 if json.Unmarshal(b,&p)!=nil||p.Target!=target||p.Key!=pgoKey(fs,ss,target)||p.Counts==nil{die("PGO profile does not match this project, compiler or target; retrain it")}
 allowed:=map[string]bool{};for _,f:=range fs{if !f.record{allowed[f.name]=true}}
 for name:=range p.Counts{if !allowed[name]{die("unknown function in PGO profile: %s",name)}}
 return p.Counts
}

type pgoLog struct{ buf bytes.Buffer }
func (b *pgoLog) Write(p []byte)(int,error){
 if len(p)>pgoMaxLog-b.buf.Len(){return 0,fmt.Errorf("PGO observation limit exceeded")}
 return b.buf.Write(p)
}
func pgoCounts(log string,fs []fn)map[string]uint64{
 known:=map[string]string{}
 for _,f:=range fs{if !f.record{known[hex.EncodeToString([]byte(f.name))]=f.name}}
 counts:=map[string]uint64{}
 for _,line:=range strings.Split(log,"\n"){
  if !strings.HasPrefix(line,pgoPrefix){continue}
  if name,ok:=known[strings.TrimPrefix(line,pgoPrefix)];ok{counts[name]++}
 }
 return counts
}
func pgoTrainCLI(args []string){
 if len(args)!=3||args[1]!="-o"||args[2]==""{die("usage: tripol pgo-train FILE.tri -o PROFILE.json")}
 if diagnosticTrace||runtime.GOARCH!="amd64"||(runtime.GOOS!="windows"&&runtime.GOOS!="linux"){die("PGO training requires an untraced local x86-64 Windows or Linux host")}
 entryFile=args[0]
 fs,ss,_:=loadProject(args[0])
 key:=pgoKey(fs,ss,runtime.GOOS)
 // Debug emitter checks all functions; instrumented output is never cached.
 image:=linkCode(generateCodePGO(fs,ss,runtime.GOOS,false,nil,false,nil,true))
 dir,e:=os.MkdirTemp("","tripol-pgo-");if e!=nil{die("%v",e)};defer os.RemoveAll(dir)
 exe:=filepath.Join(dir,"training");if runtime.GOOS=="windows"{exe+=".exe"}
 write(exe,image)
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 cmd:=exec.CommandContext(ctx,exe)
 cmd.Stdout=io.Discard
 var log pgoLog
 cmd.Stderr=&log
 cmd.WaitDelay=2*time.Second
 if e=cmd.Run();e!=nil{die("PGO training failed (30s / 8MiB log limit): %v",e)}
 counts:=pgoCounts(log.buf.String(),fs)
 data,e:=json.MarshalIndent(pgoProfile{key,runtime.GOOS,counts},"","  ");if e!=nil{die("%v",e)}
 out,e:=os.OpenFile(args[2],os.O_WRONLY|os.O_CREATE|os.O_EXCL,0600);if e!=nil{die("cannot create PGO profile: %v",e)}
 _,werr:=out.Write(append(data,'\n'));cerr:=out.Close()
 if werr!=nil||cerr!=nil{_ = os.Remove(args[2]);die("cannot write PGO profile")}
 fmt.Printf("PGO profile saved: %s (%d observed functions)\n",args[2],len(counts))
}
