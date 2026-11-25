#!/usr/bin/env python3
"""Compare CORE and HYPER profiles with identical source and checked results."""
import argparse, json, pathlib, platform, random, statistics, subprocess, tempfile, time

def main():
    p=argparse.ArgumentParser();p.add_argument('--old',required=True,type=pathlib.Path);p.add_argument('--new',required=True,type=pathlib.Path);p.add_argument('--runs',type=int,default=11);p.add_argument('--out',type=pathlib.Path);a=p.parse_args()
    if not 3<=a.runs<=100:p.error('runs must be 3..100')
    old,new=a.old.resolve(),a.new.resolve();root=pathlib.Path(__file__).resolve().parents[1]
    profiles=[('CORE nano',old,'nano'),('CORE turbo',old,'turbo'),('HYPER nano',new,'nano'),('HYPER balanced',new,'balanced'),('HYPER turbo',new,'turbo')]
    result={'host':platform.platform(),'runs':a.runs,'versions':[subprocess.check_output([c,'--version'],text=True).strip() for c in [old,new]],'runtime':{}}
    rng=random.Random(17)
    with tempfile.TemporaryDirectory() as d:
        d=pathlib.Path(d)
        for name in ['bench_core_division.tri','bench_core_arithmetic.tri','bench_hyper_calls.tri','fibonacci.tri']:
            src=root/'examples'/name
            paths=[d/f'p{i}' for i in range(len(profiles))]
            for (_,c,prof),path in zip(profiles,paths):subprocess.run([c,'build',src,'--target','linux','--profile',prof,'-o',path],check=True,stdout=subprocess.DEVNULL)
            outputs=[subprocess.check_output([path]) for path in paths]
            if any(x!=outputs[0] for x in outputs):raise RuntimeError(f'{name}: output mismatch')
            samples=[[] for _ in paths]
            for i in range(a.runs):
                order=list(range(len(paths)));rng.shuffle(order)
                for j in order:
                    t=time.perf_counter_ns();out=subprocess.check_output([paths[j]]);elapsed=(time.perf_counter_ns()-t)/1e6
                    if out!=outputs[0]:raise RuntimeError('unexpected output')
                    samples[j].append(elapsed)
            result['runtime'][name]={'output':outputs[0].decode().strip(),'profiles':{k:{'median_ms':statistics.median(s),'samples_ms':s,'bytes':path.stat().st_size} for (k,_,_),s,path in zip(profiles,samples,paths)}}
        src=d/'large.tri';src.write_text('func adjust(n) {return n*3+7}\ns=0\n'+''.join(f's+=adjust({i})\n' for i in range(5000))+'show s\n')
        # Both compilers execute exactly the same external command. Includes startup/write.
        result['build_end_to_end']={}
        for name,src in [('5000 calls',src),('game',root/'examples/nova_game.tri')]:
            samples=[[],[]];out=d/'app.exe'
            for c in [old,new]:subprocess.run([c,'build',src,'--target','windows','-o',out],check=True,stdout=subprocess.DEVNULL)
            for i in range(a.runs):
                for j in ([0,1] if i%2==0 else [1,0]):
                    t=time.perf_counter_ns();subprocess.run([[old,new][j],'build',src,'--target','windows','-o',out],check=True,stdout=subprocess.DEVNULL);samples[j].append((time.perf_counter_ns()-t)/1e6)
            result['build_end_to_end'][name]=[{'median_ms':statistics.median(s),'samples_ms':s} for s in samples]
    text=json.dumps(result,indent=2);print(text)
    if a.out:a.out.write_text(text+'\n')
if __name__=='__main__':main()
