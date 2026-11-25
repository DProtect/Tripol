package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Embedded workloads: no Python, Cargo, network or package dependencies.
// Rust receives the same runtime input, i64 arithmetic and overflow checks.
type arenaCase struct{ name, tri, rust, input, expected string }

const rustInput = `use std::io;
use std::time::Instant;
fn main() {
 let mut line=String::new(); io::stdin().read_line(&mut line).unwrap();
 let n:i64=line.trim().parse().unwrap();
`
const triInput = "import tripol.time\nn=input_int()\n"

func arenaCases() []arenaCase {
	var a, b, c, d int64
	for i := int64(0); i < 6000000; i++ {
		a += i % 16
		b += a % 8
		if i%2 == 0 {
			c += b % 32
		} else {
			d += a % 32
		}
	}
	return []arenaCase{
		{"division", triInput + `start=time.now_ns()
s=0
for i in range(n) { x=i-n/2; s+=x/8; s+=x%8 }
ns=time.elapsed_ns(start)
show s
show "TRIPOL_NS:",ns
`, rustInput + `let start=Instant::now(); let mut s=0i64;
for i in 0..n {let x=i-n/2; s+=x/8; s+=x%8;}
let ns=start.elapsed().as_nanos(); println!("{}",s); println!("TRIPOL_NS: {}",ns);
}
`, "12000000", "-750000"},
		{"registers", triInput + `start=time.now_ns()
a=0;b=0;c=0;d=0
for i in range(n) {
 a+=i%16; b+=a%8
 if i%2==0 {c+=b%32} else {d+=a%32}
}
ns=time.elapsed_ns(start)
show a,b,c,d
show "TRIPOL_NS:",ns
`, rustInput + `let start=Instant::now(); let (mut a,mut b,mut c,mut d)=(0i64,0i64,0i64,0i64);
for i in 0..n {a+=i%16; b+=a%8; if i%2==0 {c+=b%32} else {d+=a%32}}
let ns=start.elapsed().as_nanos(); println!("{} {} {} {}",a,b,c,d); println!("TRIPOL_NS: {}",ns);
}
`, "6000000", fmt.Sprintf("%d %d %d %d", a, b, c, d)},
		{"fibonacci", `import tripol.time
func fib(n:int)->int {if n<=1 {return n};return fib(n-1)+fib(n-2)}
n=input_int()
start=time.now_ns()
s=fib(n)
ns=time.elapsed_ns(start)
show s
show "TRIPOL_NS:",ns
`, `fn fib(n:i64)->i64 {if n<=1 {n} else {fib(n-1)+fib(n-2)}}
` + rustInput + `let start=Instant::now(); let s=fib(n); let ns=start.elapsed().as_nanos(); println!("{}",s); println!("TRIPOL_NS: {}",ns);
}
`, "30", "832040"},

		{"arrays", triInput + `data=ints(n)
for i in range(n){data[i]=i%1024}
start=time.now_ns()
s=0
for i in range(n){s+=data[i]}
ns=time.elapsed_ns(start)
show s
show "TRIPOL_NS:",ns
`, rustInput + `let mut data=vec![0i64;n as usize];for i in 0..n {data[i as usize]=i%1024;}
let start=Instant::now();let mut s=0i64;for i in 0..n {s+=data[i as usize];}
let ns=start.elapsed().as_nanos();println!("{}",s);println!("TRIPOL_NS: {}",ns);}
`, "1000000", "511370976"},
		{"copy", "import tripol.arrays\n" + triInput + `data=ints(n)
for i in range(n){data[i]=i%1024}
start=time.now_ns()
b=arrays.copy(data)
ns=time.elapsed_ns(start)
s=0
for i in range(n){s+=b[i]}
show s
show "TRIPOL_NS:",ns
`, rustInput + `let mut data=vec![0i64;n as usize];for i in 0..n {data[i as usize]=i%1024;}
let start=Instant::now();let b=std::hint::black_box(&data).to_vec();let ns=start.elapsed().as_nanos();
let mut s=0i64;for x in b {s+=x;} println!("{}",s);println!("TRIPOL_NS: {}",ns);}
`, "1000000", "511370976"},
		{"search", "import tripol.search\n" + triInput + `data=ints(n)
for i in range(n){data[i]=i%1024}
start=time.now_ns()
s=search.linear(data,n)
ns=time.elapsed_ns(start)
show s
show "TRIPOL_NS:",ns
`, rustInput + `let mut data=vec![0i64;n as usize];for i in 0..n {data[i as usize]=i%1024;}
let start=Instant::now();let s=std::hint::black_box(&data).iter().position(|&x|x==n).map(|i|i as i64).unwrap_or(-1);let ns=start.elapsed().as_nanos();println!("{}",s);println!("TRIPOL_NS: {}",ns);}
`, "1000000", "-1"},
		{"copy-reuse", "import tripol.arrays\n" + triInput + `data=ints(n);b=ints(n)
for i in range(n){data[i]=i%1024}
start=time.now_ns()
assert arrays.copy_into(b,data)
ns=time.elapsed_ns(start)
s=0
for i in range(n){s+=b[i]}
show s
show "TRIPOL_NS:",ns
`, rustInput + `let mut data=vec![0i64;n as usize];let mut b=vec![0i64;n as usize];for i in 0..n {data[i as usize]=i%1024;}
let start=Instant::now();b.copy_from_slice(std::hint::black_box(&data));let ns=start.elapsed().as_nanos();
let mut s=0i64;for x in b {s+=x;} println!("{}",s);println!("TRIPOL_NS: {}",ns);}
`, "1000000", "511370976"},
	}
}
func arenaOutput(out, expected string) (int64, bool) {
	ns, ok := parseInternal(out)
	if !ok {
		return 0, false
	}
	var lines []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "TRIPOL_NS:") {
			lines = append(lines, line)
		}
	}
	return ns, len(lines) == 1 && lines[0] == expected
}
func arenaMedian(v []float64) float64 {
	s := append([]float64(nil), v...)
	sort.Float64s(s)
	return (s[(len(s)-1)/2] + s[len(s)/2]) / 2
}

// Nearest-rank percentile; no samples discarded.
func arenaPercentile(v []float64, p float64) float64 {
	s := append([]float64(nil), v...)
	sort.Float64s(s)
	i := int(p * float64(len(s)))
	if float64(i) < p*float64(len(s)) {
		i++
	}
	if i < 1 {
		i = 1
	}
	if i > len(s) {
		i = len(s)
	}
	return s[i-1]
}
func arenaCommand(path string, args []string, input string, limit time.Duration) (string, float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdin = strings.NewReader(input)
	start := time.Now()
	out, e := cmd.CombinedOutput()
	elapsed := float64(time.Since(start)) / 1e6
	if ctx.Err() != nil {
		return string(out), elapsed, ctx.Err()
	}
	return string(out), elapsed, e
}
func arenaCLI(args []string) {
	runs := 10
	profile := "turbo"
	selected := "all"
	jsonPath := ""
	only := false
	timeout := 120 * time.Second
	for len(args) > 0 {
		if args[0] == "--tripol-only" {
			only = true
			args = args[1:]
			continue
		}
		if len(args) < 2 {
			die("missing arena option value")
		}
		switch args[0] {
		case "--runs":
			n, e := strconv.Atoi(args[1])
			if e != nil || n < 1 || n > 30 {
				die("arena runs must be 1..30")
			}
			runs = n
		case "--profile":
			profile = args[1]
			if profile != "nano" && profile != "balanced" && profile != "turbo" && profile != "debug" {
				die("invalid profile")
			}
		case "--timeout":
			n, e := strconv.Atoi(args[1])
			if e != nil || n < 1 || n > 900 {
				die("arena timeout must be 1..900 seconds")
			}
			timeout = time.Duration(n) * time.Second
		case "--json":
			jsonPath = args[1]
		case "--case":
			selected = args[1]
			if selected != "all" && selected != "division" && selected != "registers" && selected != "fibonacci" && selected != "arrays" && selected != "copy" && selected != "search" && selected != "copy-reuse" {
				die("unknown arena case")
			}
		default:
			die("unknown arena option %s", args[0])
		}
		args = args[2:]
	}
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "windows" && runtime.GOOS != "linux") {
		die("arena requires Windows or Linux x86-64")
	}
	self, e := os.Executable()
	if e != nil {
		die("%v", e)
	}
	rust := ""
	rv := "not requested"
	if !only {
		rust, e = exec.LookPath("rustc")
		if e != nil {
			die("rustc not found. Install Rust from https://rustup.rs/ and reopen CMD; or use --tripol-only")
		}
		v, _, err := arenaCommand(rust, []string{"--version"}, "", timeout)
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Fprintln(os.Stderr, "Rust version probe timed out; retrying once (check disk activity).")
			v, _, err = arenaCommand(rust, []string{"--version"}, "", timeout)
		}
		if err != nil {
			die("rustc probe failed: %s %v; try rustc --version, --timeout 300 or --tripol-only", v, err)
		}
		rv = strings.TrimSpace(v)
	}
	root, e := os.MkdirTemp("", "tripol-arena-")
	if e != nil {
		die("%v", e)
	}
	defer os.RemoveAll(root)
	fmt.Printf("TRIPOL QUANTUM X Arena | %s/%s | %d runs\nTRIPOL %s (%s) | %s\n", runtime.GOOS, runtime.GOARCH, runs, version, profile, rv)
	fmt.Println("Checked i64 arithmetic. Runtime input before timer. Rust: -C opt-level=3 -C overflow-checks=yes -C debuginfo=0 -C strip=symbols.")
	fmt.Println("Build median: 3 full CLI builds incl. startup/write. Runtime: interleaved after warm-up; every result checked.")
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	report := map[string]any{"timeout_seconds": int(timeout / time.Second), "version": version, "rust": rv, "os": runtime.GOOS, "arch": runtime.GOARCH, "profile": profile, "runs": runs}
	workloads := []map[string]any{}
	for _, w := range arenaCases() {
		if selected != "all" && selected != w.name {
			continue
		}
		tri := filepath.Join(root, w.name+".tri")
		rs := filepath.Join(root, w.name+".rs")
		if err := os.WriteFile(tri, []byte(w.tri), 0600); err != nil {
			die("%v", err)
		}
		if err := os.WriteFile(rs, []byte(w.rust), 0600); err != nil {
			die("%v", err)
		}
		type contestant struct {
			name, tool, bin          string
			args                     []string
			build, internal, process []float64
			size                     int64
		}
		tbin := filepath.Join(root, w.name+"-tripol"+ext)
		rbin := filepath.Join(root, w.name+"-rust"+ext)
		cs := []contestant{{name: "TRIPOL", tool: self, bin: tbin, args: []string{"build", tri, "--profile", profile, "-o", tbin}}}
		if !only {
			cs = append(cs, contestant{name: "Rust", tool: rust, bin: rbin, args: []string{rs, "-C", "opt-level=3", "-C", "overflow-checks=yes", "-C", "debuginfo=0", "-C", "strip=symbols", "-o", rbin}})
		}
		fmt.Printf("\n=== %s | input %s ===\n", w.name, w.input)
		for i := 0; i < 3; i++ {
			for j := range cs {
				k := j
				if i%2 == 1 {
					k = len(cs) - 1 - j
				}
				c := &cs[k]
				out, ms, err := arenaCommand(c.tool, c.args, "", timeout)
				if err != nil {
					die("%s build failed: %v\n%s", c.name, err, out)
				}
				c.build = append(c.build, ms)
			}
		}
		run := func(c *contestant, record bool) {
			out, ms, err := arenaCommand(c.bin, nil, w.input+"\n", timeout)
			ns, ok := arenaOutput(out, w.expected)
			if err != nil || !ok {
				die("%s %s result mismatch: expected %s; got %q (%v)", w.name, c.name, w.expected, out, err)
			}
			if record {
				c.internal = append(c.internal, float64(ns)/1e6)
				c.process = append(c.process, ms)
			}
		}
		for i := range cs {
			st, err := os.Stat(cs[i].bin)
			if err != nil {
				die("%v", err)
			}
			cs[i].size = st.Size()
			run(&cs[i], false)
		}
		for i := 0; i < runs; i++ {
			for j := range cs {
				k := j
				if i%2 == 1 {
					k = len(cs) - 1 - j
				}
				run(&cs[k], true)
			}
		}
		fmt.Printf("Verified output: %s\n", w.expected)
		entries := []map[string]any{}
		for _, c := range cs {
			entries = append(entries, map[string]any{"name": c.name, "build_median_ms": arenaMedian(c.build), "internal_median_ms": arenaMedian(c.internal), "process_median_ms": arenaMedian(c.process), "bytes": c.size, "build_samples_ms": c.build, "internal_samples_ms": c.internal, "process_samples_ms": c.process, "internal_p95_ms": arenaPercentile(c.internal, 0.95)})
			fmt.Printf("%-7s | build %.3f ms | internal %.6f ms | process %.3f ms | %d bytes\n", c.name, arenaMedian(c.build), arenaMedian(c.internal), arenaMedian(c.process), c.size)
		}
		workloads = append(workloads, map[string]any{"case": w.name, "input": w.input, "verified_output": w.expected, "contestants": entries})
		if len(cs) == 2 {
			t, r := arenaMedian(cs[0].internal), arenaMedian(cs[1].internal)
			if t > 0 && r > 0 {
				fmt.Printf("Internal Rust/TRIPOL ratio: %.3fx (above 1 favors TRIPOL; below 1 favors Rust)\n", r/t)
			}
		}
	}
	if jsonPath != "" {
		report["workloads"] = workloads
		b, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			die("%v", err)
		}
		if err = os.WriteFile(jsonPath, append(b, '\n'), 0644); err != nil {
			die("%v", err)
		}
		fmt.Println("JSON report:", jsonPath)
	}
	fmt.Println("Results are workload-specific; Rust includes its standard runtime and uses its own optimizer. File size is not RAM usage.")
}
