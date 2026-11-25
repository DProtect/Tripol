#!/usr/bin/env python3
"""Reproducible v0.15/v0.16 comparison; no third-party packages needed."""
import argparse, json, pathlib, re, statistics, subprocess, tempfile, time, platform

def main():
    parser=argparse.ArgumentParser()
    parser.add_argument('--old',required=True,type=pathlib.Path)
    parser.add_argument('--new',required=True,type=pathlib.Path)
    parser.add_argument('--runs',type=int,default=11)
    parser.add_argument('--out',type=pathlib.Path)
    args=parser.parse_args()
    if not 3<=args.runs<=100: parser.error('runs must be 3..100')
    compilers=[args.old.resolve(),args.new.resolve()]
    root=pathlib.Path(__file__).resolve().parents[1]
    result={'host':platform.platform(),'runs':args.runs,'compilers':[subprocess.check_output([c,'--version'],text=True).strip() for c in compilers],'build':{},'runtime':{}}
    with tempfile.TemporaryDirectory() as tmp:
        d=pathlib.Path(tmp)
        source=d/'large.tri'
        source.write_text('func adjust(n) {return n*3+7}\ns=0\n'+''.join(f's+=adjust({i})\n' for i in range(5000))+'show s\n')
        # Native compiler's internal benchmark excludes startup and executable writing.
        for name,src in [('5000 calls',source),('game',root/'examples/nova_game.tri')]:
            values=[]
            for c in compilers:
                text=subprocess.check_output([c,'bench-build',src,'--runs',str(args.runs)],text=True)
                m=re.search(r'Build median: ([\d.]+) ms \| Min: ([\d.]+) ms \| Max: ([\d.]+) ms \| Output: (\d+) bytes',text)
                values.append(dict(zip(['median_ms','min_ms','max_ms','bytes'],[float(m[1]),float(m[2]),float(m[3]),int(m[4])])))
            result['build'][name]=values
        for name in ['bench_core_division.tri','bench_core_arithmetic.tri','fibonacci.tri']:
            src=root/'examples'/name
            if not src.exists(): continue
            paths=[d/f'run{i}' for i in range(2)]
            for c,p in zip(compilers,paths):subprocess.run([c,'build',src,'--target','linux','-o',p],check=True,stdout=subprocess.DEVNULL)
            outputs=[subprocess.check_output([p]) for p in paths]
            if outputs[0]!=outputs[1]:raise RuntimeError(f'Output mismatch {name}')
            samples=[[],[]]
            for i in range(args.runs):
                for j in ([0,1] if i%2==0 else [1,0]):
                    start=time.perf_counter_ns();got=subprocess.check_output([paths[j]]);elapsed=(time.perf_counter_ns()-start)/1e6
                    if got!=outputs[0]: raise RuntimeError('unexpected output')
                    samples[j].append(elapsed)
            result['runtime'][name]={'output':outputs[0].decode().strip(),'versions':[{'median_ms':statistics.median(x),'samples_ms':x,'bytes':p.stat().st_size} for p,x in zip(paths,samples)]}
    text=json.dumps(result,indent=2)
    print(text)
    if args.out: args.out.write_text(text+'\n')
if __name__=='__main__':main()
