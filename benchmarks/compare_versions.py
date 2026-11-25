#!/usr/bin/env python3
"""Compare two native TRIPOL compilers using the same Arena workloads."""
import argparse, json, os, pathlib, subprocess, tempfile

def main():
    p=argparse.ArgumentParser()
    p.add_argument("--old",required=True);p.add_argument("--new",required=True)
    p.add_argument("--runs",type=int,default=7);p.add_argument("--out",required=True)
    args=p.parse_args()
    if not 1<=args.runs<=30: p.error("runs must be 1..30")
    env=dict(os.environ,TRIPOL_NO_UPDATE="1")
    with tempfile.TemporaryDirectory(prefix="tripol-compare-") as tmp:
        root=pathlib.Path(tmp)
        reports=[]
        for label,exe in [("old",args.old),("new",args.new)]:
            output=root/(label+".json")
            subprocess.run([str(pathlib.Path(exe).resolve()),"arena","--runs",str(args.runs),"--profile","turbo","--tripol-only","--json",str(output)],cwd=root,env=env,check=True,timeout=1200)
            reports.append(json.loads(output.read_text(encoding="utf8")))
        old,new=reports
        before={x["case"]:x for x in old["workloads"]}
        result=[]
        for case in new["workloads"]:
            previous=before.get(case["case"])
            if previous is None: continue
            if previous["input"]!=case["input"] or previous["verified_output"]!=case["verified_output"]:
                raise ValueError("workload differs: "+case["case"])
            a=previous["contestants"][0];b=case["contestants"][0]
            result.append({"case":case["case"],"old_median_ms":a["internal_median_ms"],"new_median_ms":b["internal_median_ms"],"old_p95_ms":a["internal_p95_ms"],"new_p95_ms":b["internal_p95_ms"],"old_bytes":a["bytes"],"new_bytes":b["bytes"],"ratio_old_over_new":a["internal_median_ms"]/b["internal_median_ms"] if b["internal_median_ms"]>0 else None})
        if not result: raise ValueError("no comparable cases")
        report={"method":"same host/profile/input/output; separate sequential Arena runs, not interleaved across versions","old":old,"new":new,"comparison":result}
        path=pathlib.Path(args.out);path.parent.mkdir(parents=True,exist_ok=True);path.write_text(json.dumps(report,indent=2)+"\n",encoding="utf8")
        for r in result: print(f"COMPARE {r['case']} old={r['old_median_ms']:.6f}ms new={r['new_median_ms']:.6f}ms ratio={r['ratio_old_over_new']}")
if __name__=="__main__":main()
