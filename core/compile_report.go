package core

import(
 "crypto/sha256"
 "encoding/json"
 "fmt"
 "runtime"
 "time"
)

type compilationReport struct{
 Version string `json:"version"`
 Target string `json:"target"`
 Profile string `json:"profile"`
 LoadNS int64 `json:"load_ns"`
 ValidationNS int64 `json:"validation_ns"`
 OptimizationNS int64 `json:"optimization_ns"`
 CFGSubstitutions int `json:"cfg_substitutions"`
 CFGLoopInvariants int `json:"cfg_loop_invariants"`
 CFGJoins int `json:"cfg_joins"`
 CFGConstantBranches int `json:"cfg_constant_branches"`
 ReachabilityNS int64 `json:"reachability_ns"`
 EmissionNS int64 `json:"emission_ns"`
 LinkNS int64 `json:"link_ns"`
 TotalNS int64 `json:"total_ns"`
 SourceFiles int `json:"source_files"`
 Functions int `json:"functions"`
 EmittedFunctions int `json:"emitted_functions"`
 ValidationBufferBytes int `json:"validation_buffer_capacity_bytes"`
 ImageBytes int `json:"image_bytes"`
 SHA256 string `json:"sha256"`
}

func measuredCompilation(fs []fn,ss []stmt,target,profile string)([]byte,compilationReport){
 r:=compilationReport{Version:version,Target:target,Profile:profile}
 start:=time.Now()
 if profile!="debug"{
  phase:=time.Now();validation:=validateCode(fs,ss,target)
  r.ValidationNS=time.Since(phase).Nanoseconds()
  r.ValidationBufferBytes=cap(validation.b)
 }
 phase:=time.Now()
 if profile=="balanced"||profile=="turbo"{
  var stats cfgConstantStats
  fs,ss,stats=optimizeCFG(fs,ss)
  r.CFGSubstitutions=stats.Substitutions;r.CFGJoins=stats.Joins;r.CFGLoopInvariants=stats.LoopInvariants;r.CFGConstantBranches=stats.ConstantBranches
 }
 r.OptimizationNS=time.Since(phase).Nanoseconds()
 phase=time.Now()
 var live map[string]bool
 if profile!="debug"{live=reachable(fs,ss)}
 for _,f:=range fs{if !f.record{r.Functions++;if profile=="debug"||live[f.name]{r.EmittedFunctions++}}}
 r.ReachabilityNS=time.Since(phase).Nanoseconds()
 phase=time.Now()
 a:=generateCode(fs,ss,target,profile!="debug",live,profile=="turbo",profile=="balanced"||profile=="turbo")
 r.EmissionNS=time.Since(phase).Nanoseconds()
 phase=time.Now()
 bin:=linkCode(a)
 r.LinkNS=time.Since(phase).Nanoseconds()
 r.TotalNS=time.Since(start).Nanoseconds()
 r.ImageBytes=len(bin);r.SHA256=fmt.Sprintf("%x",sha256.Sum256(bin))
 return bin,r
}

func compileReportCLI(args []string){
 file:="main.tri";profile:="balanced";target:=runtime.GOOS
 if target!="windows"{target="linux"}
 jsonOutput:=false;chosen:=false
 for i:=0;i<len(args);i++{
  switch args[i]{
  case "--json":jsonOutput=true
  case "--profile","--target":
   flag:=args[i];i++;if i>=len(args){die("missing value for %s",flag)}
   if flag=="--profile"{profile=args[i]}else{target=args[i]}
  default:
   if chosen||len(args[i])==0||args[i][0]=='-'{die("usage: tripol compile-report [FILE.tri] [--profile PROFILE] [--target TARGET] [--json]")}
   file=args[i];chosen=true
  }
 }
 switch profile{case "nano","balanced","turbo","debug":default:die("unknown profile %s",profile)}
 if target!="linux"&&target!="windows"{die("unsupported target %s",target)}
 entryFile=file
 begin:=time.Now()
 fs,ss,inputs:=loadProject(file)
 loadNS:=time.Since(begin).Nanoseconds()
 _,r:=measuredCompilation(fs,ss,target,profile)
 r.LoadNS=loadNS;r.TotalNS+=loadNS;r.SourceFiles=len(inputs)
 if jsonOutput{b,e:=json.MarshalIndent(r,"","  ");if e!=nil{die("%v",e)};fmt.Println(string(b));return}
 fmt.Printf("TRIPOL %s | %s | %s\n",version,target,profile)
 for _,p:=range []struct{name string;ns int64}{
 {"Load and parse",r.LoadNS},{"Validation",r.ValidationNS},{"CFG optimization",r.OptimizationNS},{"Reachability",r.ReachabilityNS},{"Native emission",r.EmissionNS},{"Runtime and linking",r.LinkNS},{"Total",r.TotalNS},
 }{fmt.Printf("%-22s %10.3f ms\n",p.name,float64(p.ns)/1e6)}
 fmt.Printf("CFG substitutions: %d; agreed joins: %d; constant branches: %d; loop invariants: %d\n",r.CFGSubstitutions,r.CFGJoins,r.CFGConstantBranches,r.CFGLoopInvariants)
 fmt.Printf("Functions selected: %d / %d\nValidation code buffer: %d bytes (capacity, not RSS)\nImage: %d bytes\nSHA256: %s\n",r.EmittedFunctions,r.Functions,r.ValidationBufferBytes,r.ImageBytes,r.SHA256)
}
