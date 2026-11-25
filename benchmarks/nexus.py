"""Run on Linux x64: python benchmarks/nexus.py /path/omega /path/eclipse."""
import json, pathlib, statistics, subprocess, sys, tempfile
root = pathlib.Path(__file__).resolve().parents[1]
compilers = [str(pathlib.Path(p).resolve()) for p in sys.argv[1:]]
assert len(compilers) == 2, 'pass ECLIPSE and NEXUS compilers'
report = {'platform': 'Linux x86-64', 'profile': 'turbo', 'runs':21, 'cases':{}}
with tempfile.TemporaryDirectory() as d:
    for case, expected in [("search","-1")]:
        bins=[]
        for i, compiler in enumerate(compilers):
            binary=pathlib.Path(d)/f'{case}{i}'
            subprocess.run([compiler,'build',str(root/'examples'/f'bench_eclipse_{case}.tri'),'--profile','turbo','-o',str(binary)],check=True,capture_output=True)
            bins.append(binary)
        def run(i):
            out=subprocess.run([str(bins[i])],input='1000000\n',text=True,capture_output=True,check=True).stdout.splitlines()
            assert out[0]==expected and len(out)==2 and out[1].startswith('TRIPOL_NS: '),out
            return int(out[1].split(':')[1])/1e6
        run(0);run(1)
        samples=[[],[]]
        for j in range(21):
            for i in ([0,1] if j%2==0 else [1,0]):samples[i].append(run(i))
        report['cases'][case]={'verified':expected,'samples_ms':samples,'medians_ms':[statistics.median(s) for s in samples]}
print(json.dumps(report,indent=2))
