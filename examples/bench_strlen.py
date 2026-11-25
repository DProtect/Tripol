# Generate the same workload for both versions; Python is maintainer tooling only.
from pathlib import Path
import json
Path('bench_strlen.tri').write_text('s='+json.dumps('a'*32768)+'\ntotal=0\nrepeat 20000 {total+=len(s)}\nshow total\n',encoding='utf-8')
