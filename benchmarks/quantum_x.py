#!/usr/bin/env python3
"""Linux checked Fibonacci comparison: independent tools, identical algorithm/output."""
import argparse,json,os,pathlib,platform,shutil,statistics,subprocess,sys,tempfile,time

TRI='''import tripol.time
func fib(n:int)->int {if n<=1 {return n};return fib(n-1)+fib(n-2)}
n=input_int()
assert n>=0 and n<=35
start=time.now_ns()
answer=fib(n)
ns=time.elapsed_ns(start)
show answer
show "TRIPOL_NS:",ns
'''
C=r'''#include <stdint.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
static int64_t fib(int64_t n){if(n<=1)return n;int64_t a=fib(n-1),b=fib(n-2),v;if(__builtin_add_overflow(a,b,&v))abort();return v;}
static int64_t ns(void){struct timespec t;if(clock_gettime(CLOCK_MONOTONIC,&t))abort();return (int64_t)t.tv_sec*1000000000+t.tv_nsec;}
int main(void){int64_t n;if(scanf("%" SCNd64,&n)!=1||n<0||n>35)return 1;int64_t s=ns(),v=fib(n),e=ns()-s;printf("%" PRId64 "\nTRIPOL_NS: %" PRId64 "\n",v,e);}
'''
RUST=r'''use std::io;
use std::time::Instant;
fn fib(n:i64)->i64{if n<=1{return n;}fib(n-1).checked_add(fib(n-2)).unwrap()}
fn main(){let mut s=String::new();io::stdin().read_line(&mut s).unwrap();let n:i64=s.trim().parse().unwrap();assert!((0..=35).contains(&n));let t=Instant::now();let v=fib(n);let ns=t.elapsed().as_nanos();println!("{}\nTRIPOL_NS: {}",v,ns);}
'''
GO=r'''package main
import("fmt";"time";"math")
func fib(n int64)int64{if n<=1{return n};a,b:=fib(n-1),fib(n-2);if a>math.MaxInt64-b{panic("overflow")};return a+b}
func main(){var n int64;if _,e:=fmt.Scan(&n);e!=nil||n<0||n>35{panic("input")};t:=time.Now();v:=fib(n);ns:=time.Since(t).Nanoseconds();fmt.Printf("%d\nTRIPOL_NS: %d\n",v,ns)}
'''
PYTHON='''import time
def fib(n):
    if n<=1:return n
    v=fib(n-1)+fib(n-2)
    if v>9223372036854775807:raise OverflowError()
    return v
n=int(input())
assert 0<=n<=35
start=time.perf_counter_ns()
v=fib(n)
elapsed=time.perf_counter_ns()-start
print(v)
print("TRIPOL_NS:",elapsed)
'''
ZIG=r'''extern "c" fn scanf([*:0]const u8, ...) c_int;
extern "c" fn printf([*:0]const u8, ...) c_int;
const Timespec=extern struct{sec:i64,nsec:i64};
extern "c" fn clock_gettime(c_int,*Timespec)c_int;
fn now()i64{var t:Timespec=undefined;if(clock_gettime(1,&t)!=0)@panic("clock");return t.sec*1000000000+t.nsec;}
fn fib(n:i64)i64{if(n<=1)return n;return fib(n-1)+fib(n-2);}
pub fn main()void{var n:i64=0;if(scanf("%lld",&n)!=1 or n<0 or n>35)@panic("input");const t=now();const v=fib(n);const ns=now()-t;_=printf("%lld\nTRIPOL_NS: %lld\n",v,ns);}
'''
def run(cmd,cwd,input_text=None,timeout=180):
    start=time.perf_counter_ns()
    proc=subprocess.run(cmd,cwd=cwd,input=input_text,text=True,capture_output=True,timeout=timeout,env=dict(os.environ,TRIPOL_NO_UPDATE="1"))
    ms=(time.perf_counter_ns()-start)/1e6
    if proc.returncode:raise RuntimeError(f"{cmd[0]}: {proc.stderr[-1000:]}")
    return proc.stdout,ms

def p95(values):
    import math
    return sorted(values)[max(0,math.ceil(len(values)*.95)-1)]
def main():
    p=argparse.ArgumentParser();p.add_argument("--compiler",required=True);p.add_argument("--runs",type=int,default=7);p.add_argument("--out",required=True);a=p.parse_args()
    if platform.system()!="Linux" or not 1<=a.runs<=30:p.error("Linux, 1..30 runs required")
    results=[];skipped=[]
    with tempfile.TemporaryDirectory(prefix="tripol-x-") as td:
        d=pathlib.Path(td)
        cases=[("TRIPOL",str(pathlib.Path(a.compiler).resolve()),".tri",TRI,lambda tool,s,b:[tool,"build",str(s),"--profile","turbo","-o",str(b)]),("C",shutil.which("gcc"),".c",C,lambda t,s,b:[t,"-O3",str(s),"-o",str(b)]),("C++",shutil.which("g++"),".cpp",C,lambda t,s,b:[t,"-O3",str(s),"-o",str(b)]),("Rust",shutil.which("rustc"),".rs",RUST,lambda t,s,b:[t,str(s),"-C","opt-level=3","-C","overflow-checks=yes","-C","strip=symbols","-o",str(b)]),("Go",shutil.which("go"),".go",GO,lambda t,s,b:[t,"build","-trimpath","-ldflags=-s -w","-o",str(b),str(s)]),("Zig",shutil.which("zig"),".zig",ZIG,lambda t,s,b:[t,"build-exe",str(s),"-O","ReleaseSafe","-lc","-femit-bin="+str(b)]),("Python",sys.executable,".py",PYTHON,None)]
        for i,(name,tool,ext,source,build) in enumerate(cases):
            if not tool:skipped.append({"name":name,"reason":"compiler unavailable"});continue
            s=d/(f"case{i}"+ext);b=d/f"bin{i}";s.write_text(source,encoding="utf8")
            builds=[]
            try:
                if build:
                    for _ in range(3):_,ms=run(build(tool,s,b),d);builds.append(ms)
                    cmd=[str(b)];size=b.stat().st_size
                else:cmd=[tool,str(s)];size=None
                version_args=["--version"] if name!="Go" else ["version"]
                ver,_=run([tool]+version_args,d)
                results.append({"name":name,"tool_version":ver.splitlines()[0],"cmd":cmd,"bytes":size,"build_samples_ms":builds,"internal_samples_ms":[],"process_samples_ms":[]})
            except Exception as e:
                if name=="TRIPOL":raise
                skipped.append({"name":name,"reason":str(e)})
        def sample(r,record):
            out,ms=run(r["cmd"],d,"30\n")
            lines=out.strip().splitlines()
            if len(lines)!=2 or lines[0]!="832040" or not lines[1].startswith("TRIPOL_NS:"):raise ValueError(r["name"]+" output mismatch "+repr(out))
            ns=int(lines[1].split(":")[1])
            if ns<0:raise ValueError("negative timer")
            if record:r["internal_samples_ms"].append(ns/1e6);r["process_samples_ms"].append(ms)
        for r in results:sample(r,False)
        for i in range(a.runs):
            for r in results[::1 if i%2==0 else -1]:sample(r,True)
        for r in results:
            r["internal_median_ms"]=statistics.median(r["internal_samples_ms"])
            r["internal_p95_ms"]=p95(r["internal_samples_ms"])
            r["process_median_ms"]=statistics.median(r["process_samples_ms"])
            r["build_median_ms"]=statistics.median(r["build_samples_ms"]) if r["build_samples_ms"] else None
            # Separate RSS probe to avoid timing wrapper overhead in measured samples.
            timer="/usr/bin/time"
            if pathlib.Path(timer).is_file():
                rss=d/"rss.txt"
                try:
                    subprocess.run([timer,"-f","%M","-o",str(rss),*r["cmd"]],input="30\n",text=True,capture_output=True,check=True,timeout=180)
                    r["peak_rss_kib_separate_probe"]=int(rss.read_text().strip())
                except Exception:r["peak_rss_kib_separate_probe"]=None
            del r["cmd"]
            print(f"LANG {r['name']} {r['internal_median_ms']:.6f}ms size={r['bytes']}")
        out=pathlib.Path(a.out);out.parent.mkdir(parents=True,exist_ok=True)
        out.write_text(json.dumps({"method":"checked signed Fibonacci(30), input before internal timer, 3 builds, alternating runs, warm-up, verified result 832040; Linux RSS separate probe; Python source is not a native binary","host":platform.platform(),"runs":a.runs,"results":results,"skipped":skipped},indent=2)+"\n",encoding="utf8")
if __name__=="__main__":main()
