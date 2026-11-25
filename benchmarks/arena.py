#!/usr/bin/env python3
"""Optional maintainer benchmark: identical naive recursive Fibonacci(30).
Process startup included. Native file size only; checked TRIPOL vs unchecked C/Go/Rust.
No dependency on this script for ordinary TRIPOL users.
"""
import argparse, json, os, pathlib, shutil, statistics, subprocess, sys, tempfile, time

SOURCES = {
    'TRIPOL': ('fib.tri', 'func fib(n) {if n<=1 {return n};return fib(n-1)+fib(n-2)}\nshow fib(30)\n'),
    'C': ('fib.c', '#include <stdio.h>\n#include <stdint.h>\nint64_t fib(int64_t n){return n<=1?n:fib(n-1)+fib(n-2);}\nint main(void){printf("%lld\\n",(long long)fib(30));}\n'),
    'C++': ('fib.cpp', '#include <cstdio>\n#include <cstdint>\nstd::int64_t fib(std::int64_t n){return n<=1?n:fib(n-1)+fib(n-2); }\nint main(){std::printf("%lld\\n",(long long)fib(30));}\n'),
    'Go': ('fib.go', 'package main\nimport "fmt"\nfunc fib(n int64) int64 {if n<=1{return n};return fib(n-1)+fib(n-2)}\nfunc main(){fmt.Println(fib(30))}\n'),
    'Rust': ('fib.rs', 'fn fib(n:i64)->i64 {if n<=1 {n} else {fib(n-1)+fib(n-2)}}\nfn main(){println!("{}",fib(30));}\n'),
    'Python': ('fib.py', 'def fib(n):\n    if n<=1: return n\n    return fib(n-1)+fib(n-2)\nprint(fib(30))\n'),
    'JavaScript': ('fib.js', 'function fib(n){return n<=1?n:fib(n-1)+fib(n-2)}\nconsole.log(fib(30))\n'),
}

def main():
    p=argparse.ArgumentParser();p.add_argument('--compiler',default='tripol');p.add_argument('--runs',type=int,default=5);a=p.parse_args()
    if not 1<=a.runs<=30: p.error('runs must be 1..30')
    compiler=shutil.which(a.compiler) or str(pathlib.Path(a.compiler).resolve())
    tools={'TRIPOL':compiler,'C':shutil.which('gcc'),'C++':shutil.which('g++'),'Go':shutil.which('go'),'Rust':shutil.which('rustc'),'Python':sys.executable,'JavaScript':shutil.which('node')}
    with tempfile.TemporaryDirectory() as d:
        root=pathlib.Path(d);commands={};results={};ext='.exe' if os.name=='nt' else ''
        for name,(filename,source) in SOURCES.items():
            tool=tools[name]
            if not tool:
                results[name]={'status':'not installed'};continue
            src=root/filename;src.write_text(source);binary=root/(name.lower()+ext)
            builds={'TRIPOL':[tool,'build',str(src),'-o',str(binary)],'C':[tool,'-O2',str(src),'-o',str(binary)],'C++':[tool,'-O2',str(src),'-o',str(binary)],'Go':[tool,'build','-trimpath','-ldflags=-s -w','-o',str(binary),str(src)],'Rust':[tool,'-O',str(src),'-o',str(binary)]}
            try:
                start=time.perf_counter_ns()
                if name in builds:
                    subprocess.run(builds[name],check=True,capture_output=True,timeout=120);cmd=[str(binary)];size=binary.stat().st_size;build=(time.perf_counter_ns()-start)/1e6
                else: cmd=[tool,str(src)];size=None;build=None
                if subprocess.check_output(cmd,timeout=30).strip()!=b'832040': raise ValueError('wrong output')
                commands[name]=cmd;results[name]={'build_ms_single':build,'native_bytes':size,'samples_ms':[]}
            except (OSError,subprocess.SubprocessError,ValueError) as e: results[name]={'status':str(e)}
        for i in range(a.runs):
            order=list(commands);order=order if i%2==0 else list(reversed(order))
            for name in order:
                start=time.perf_counter_ns();out=subprocess.check_output(commands[name],timeout=30);elapsed=(time.perf_counter_ns()-start)/1e6
                if out.strip()!=b'832040': raise RuntimeError(name+' wrong output')
                results[name]['samples_ms'].append(round(elapsed,4))
        for r in results.values():
            if 'samples_ms' in r:r['median_ms']=statistics.median(r['samples_ms'])
        print(json.dumps({'workload':'naive recursive Fibonacci(30), output 832040','startup_included':True,'runs':a.runs,'results':results},indent=2))

if __name__=='__main__':main()
