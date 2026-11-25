package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Nexus conveniences lower to existing native operations. No interpreter or extra runtime.
func (c *compiler) inputPrompt(e *expr) {
	if len(e.args) > 1 {
		c.err(e.s + " expects zero arguments or one text prompt")
	}
	if len(e.args) == 1 {
		c.require(c.expression(e.args[0]), "string")
		c.a.call("print:string")
	}
}
func (c *compiler) nexusBuiltin(e *expr) (string, bool) {
	a := c.a
	switch e.s {
	case "input_int":
		c.inputPrompt(e)
		a.call("input")
		a.call("parseint")
		return "int", true
	case "@array_len":
		v := c.arrayVar(e.args[0].s)
		if v.typ == "ints" {
			a.load(v.off)
			a.emit(0x48, 0x8b, 0x00)
		} else {
			a.imm(int64(arrayLen(v)))
		}
		return "int", true
	case "abs":
		if len(e.args) != 1 {
			c.err("abs expects one integer")
		}
		c.require(c.expression(e.args[0]), "int")
		done := a.fresh()
		a.emit(0x48, 0x85, 0xc0)
		a.j(0x89, done)
		a.emit(0x48, 0xf7, 0xd8)
		a.j(0x80, "err")
		a.label(done)
		return "int", true
	case "min", "max":
		if len(e.args) != 2 {
			c.err(e.s + " expects two integers")
		}
		c.require(c.expression(e.args[0]), "int")
		a.emit(0x50)
		c.require(c.expression(e.args[1]), "int")
		a.emit(0x48, 0x89, 0xc1, 0x58, 0x48, 0x39, 0xc8)
		cc := byte(0x4f)
		if e.s == "max" {
			cc = 0x4c
		}
		a.emit(0x48, 0x0f, cc, 0xc1)
		return "int", true
	case "sum":
		if len(e.args) != 1 || e.args[0].op != "var" {
			c.err("sum expects a named integer array")
		}
		v := c.arrayVar(e.args[0].s)
		if v.typ == "ints" {
			a.load(v.off)
			a.emit(0x48, 0x8b, 0x08, 0x48, 0x8d, 0x50, 0x08) // rcx=len, rdx=data
		} else {
			a.emit(0xb9)
			a.u32(arrayLen(v))
			a.emit(0x48, 0x8d, 0x95)
			a.u32(v.off)
		}
		a.emit(0x31, 0xc0)
		loop, end := a.fresh(), a.fresh()
		a.label(loop)
		a.emit(0x48, 0x85, 0xc9)
		a.j(0x84, end)
		a.emit(0x48, 0x03, 0x02)
		a.j(0x80, "err")
		a.emit(0x48, 0x83, 0xc2, 8, 0x48, 0xff, 0xc9)
		a.jmp(loop)
		a.label(end)
		return "int", true
	}
	return "", false
}
func learn() {
	fmt.Print(`TRIPOL 0.11 — Learn in five minutes

1. Show values:
   show "Hello, TRIPOL!"
   show "Total:", 10 + 20

2. Read input:
   name = input("Your name?")
   age = input_int("Your age?")
   show "Hello", name

3. Decide:
   if age >= 18 { show "Adult" } else { show "Under 18" }

4. Repeat:
   repeat 3 { show "Hi" }
   for i in range(5) { show i }
   numbers = [10, 20, 30]
   for n in numbers { show n }
   show sum(numbers)
   show numbers

5. Save data:
   write_text("note.txt", "Hello")
   show read_text("note.txt")

6. Group data:
   struct Player { name: string; score: int }
   p = Player("Ana", 10)
   p.score += 5
   show p.name, p.score

7. Libraries: import tripol.math
   show math.square(7)
   tripol libs lists available modules.
   tripol docs sort shows sorting routines.
   tripol packages lists local Vault packages.

Save code as main.tri, then: tripol run main.tri
Build a standalone executable: tripol build main.tri
`)
}
func doctor() {
	exe, err := os.Executable()
	if err != nil {
		die("%v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		die("%v", err)
	}
	fmt.Printf("TRIPOL %s\nCompiler: %s\nDirectory: %s\nHost: %s/%s\nTargets: windows/amd64, linux/amd64\nExternal compiler required: no\n", version, exe, cwd, runtime.GOOS, runtime.GOARCH)
	if runtime.GOARCH != "amd64" || !strings.Contains("|linux|windows|", "|"+runtime.GOOS+"|") {
		fmt.Println("Native run: unsupported host; cross-build is available")
	} else {
		fmt.Println("Native run: supported host")
	}
}

func newProjectTemplate(dir, template string) {
	if template == "basic" {
		newProject(dir)
		return
	}
	templates := map[string]string{
		"ecosystem": `import tripol.sort
import tripol.time
import tripol.system
nums=[5,9,1,3]
start=time.now_ms()
show "Sorted:",sort.ascending(nums)
show "Platform:",system.platform()
show "Elapsed ms:",time.elapsed(start)
`,
		"libraries": `import tripol.math
import tripol.stats
numbers = [10, 20, 30]
show "Square:", math.square(7)
show "Mean:", stats.mean(numbers)
`,
		"calculator": `a = input_int("First number?")
b = input_int("Second number?")
show "Sum:", a + b
show "Product:", a * b
if b != 0 { show "Division:", float(a) / float(b) }
`,
		"records": `struct Player { name: string; score: int }
p = Player(input("Your name?"), 0)
p.score += 10
show "Player:", p.name, "Score:", p.score
`,
		"files": `name = input("Your name?")
write_text("profile.txt", "Name: " + name + "\n")
show "Saved:", read_text("profile.txt")
`,
	}
	source, ok := templates[template]
	if !ok {
		die("unknown project template %s; choose basic, calculator, files, records libraries or ecosystem", template)
	}
	if err := os.Mkdir(dir, 0755); err != nil {
		die("cannot create project: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.tri"), []byte(source), 0644); err != nil {
		die("%v", err)
	}
	fmt.Printf("Created %s (%s)\ncd %s\ntripol run\n", dir, template, dir)
}
