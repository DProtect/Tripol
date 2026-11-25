#!/usr/bin/env python3
"""Compare HYPER and TITAN with same sources and expected output."""
import argparse,json,pathlib,platform,random,statistics,subprocess,tempfile,time

def main():
 p=argparse.ArgumentParser();p.add_argument('--old',required=True,type=pathlib.Path);p.add_argument('--new',required=True,type=pathlib.Path);p.add_argument('--runs',type=int,default=11);p.add_argument('--out',type=pathlib.Path);a=p.parse_args()
 if not 3<=a.runs<=100:p.error('runs must be 3..100')
 compilers=[a.old.resolve(),a.new.resolve()];root=pathlib.Path(__file__).resolve().parents[1];rng=random.Random(18);ext='.exe' if platform.system()=='Windows' else ''
 result={'host':platform.platform(),'versions':[subprocess.check_output([c,'--version'],text=True).strip() for c in compilers],'runs':a.runs,'workloads':{}}
 with tempfile.TemporaryDirectory() as d:
  d=pathlib.Path(d)
  for name,expected in [('bench_hyper_internal.tri','-750000'),('bench_titan_registers.tri','45000000 21000000 42000000 46500000')]:
   files=[d/f'p{i}{ext}' for i in range(2)];source=root/'examples'/name;build=[[],[]];internal=[[],[]];process=[[],[]]
   for i,c in enumerate(compilers):subprocess.run([c,'build',source,'--profile','turbo','-o',files[i]],check=True,stdout=subprocess.DEVNULL)
   for _ in range(a.runs):
    order=[0,1];rng.shuffle(order)
    for i in order:
     t=time.perf_counter_ns();subprocess.run([compilers[i],'build',source,'--profile','turbo','-o',files[i]],check=True,stdout=subprocess.DEVNULL);build[i].append((time.perf_counter_ns()-t)/1e6)
   def run(i,record):
    t=time.perf_counter_ns();lines=subprocess.check_output([files[i]],text=True).splitlines();ms=(time.perf_counter_ns()-t)/1e6
    if len(lines)!=2 or lines[0]!=expected or not lines[1].startswith('TRIPOL_NS:'):raise RuntimeError(lines)
    n=int(lines[1].split(':')[1]);assert n>=0
    if record:internal[i].append(n/1e6);process[i].append(ms)
   for i in range(2):run(i,False)
   for _ in range(a.runs):
    order=[0,1];rng.shuffle(order)
    for i in order:run(i,True)
   result['workloads'][name]={'output':expected,'profiles':[{ 'internal_median_ms':statistics.median(internal[i]),'process_median_ms':statistics.median(process[i]),'build_median_ms':statistics.median(build[i]),'bytes':files[i].stat().st_size,'internal_samples_ms':internal[i],'process_samples_ms':process[i],'build_samples_ms':build[i]} for i in range(2)]}
 text=json.dumps(result,indent=2);print(text)
 if a.out:a.out.write_text(text+'\n')
if __name__=='__main__':main()
