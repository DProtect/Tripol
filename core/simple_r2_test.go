package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const simpleCallbacks = `import tripol.math

struct Box {prefix:string}

func Box.read(self:Box,value:int)->string {return self.prefix + text(value)}
func Box.read(self:Box,value:string)->string {return self.prefix + value}

func apply(action:func(int)->string) {
    action(7)
}

func render(value) {
    box=Box("n=")
    box.read(value)
}

func local() {
    const label = "inline "
    apply(func(n) {
        if n > 0 {
            label + text(n)
        } else {
            "zero"
        }
    })
}

square=math.square
box=Box("s=")
show apply(render),apply(&render),local(),box.read("ok"),Box("x=").read(5),square(4)
`

func TestSimpleR2ContextAndOverloads(t *testing.T) {
	const want = "n=7 n=7 inline 7 s=ok x=5 16\n"
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			got, failed := hyperRun(t, simpleCallbacks, profile)
			if failed || got != want {
				t.Fatalf("got %q (failed=%v)", got, failed)
			}
			compileProfile(simpleCallbacks, "windows", profile)
		})
	}
	if _, err := CompileWeb(simpleCallbacks); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleR2Diagnostics(t *testing.T) {
	for _, tc := range []struct {
		source, message string
		line            int
	}{
		{`func answer(x) {return x}`, `Cannot infer the type of parameter "x"`, 1},
		{`func answer(x) {return x}
show answer(1)
show answer("a")`, `Conflicting types for parameter "x"`, 3},
		{`func apply(cb:func(int)->int)->int {return cb(1)}
number=2
show apply(func(n) {n + number})`, `cannot capture "number"`, 3},
		{`func apply(cb:func(int)->int)->int {return cb(1)}
show apply(func(a,b) {a+b})`, `needs 1 parameters`, 2},
		{`struct X {value:int}
func X.get(self:X,x:int,y=1)->int {return x+y}
func X.get(self:X,x:int,y="a")->int {return x}
x=X(3)
show x.get(1)`, `Ambiguous overload`, 5},
		{`struct X {value:int}
func X.get(self:X,x:int)->int {return x}
func X.get(self:X,x:string)->int {return len(x)}
show X(3).get(true)`, `No overload matches`, 4},
	} {
		t.Run(tc.message, func(t *testing.T) {
			defer func() {
				v := recover()
				if v == nil {
					t.Fatal("invalid program compiled")
				}
				d := normalizeDiagnostic(v)
				if d.Line != tc.line || d.Column < 1 || !strings.Contains(d.Message, tc.message) {
					t.Fatalf("unexpected diagnostic %+v", d)
				}
			}()
			compile(tc.source, "linux")
		})
	}
}

func TestSimpleR2DiscordExamples(t *testing.T) {
	for _, tc := range []struct{ name, contains string }{
		{"embeds", `"embeds":[{"title":"TRIPOL","description":"Mi lenguaje de programacion","color":16777215`},
		{"components", `"type":17,"accent_color":16777215,"components"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Join("examples", "discord_simple_"+tc.name+".tri"))
			if err != nil {
				t.Fatal(err)
			}
			source := string(b)
			for _, target := range []string{"linux", "windows"} {
				compileProfile(source, target, "turbo")
			}
			preview := strings.Replace(source, "bot.start()", `show bot.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"data\":{\"type\":1,\"name\":\"info\"}}}")`, 1)
			if tc.name == "components" {
				preview += `show bot.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":3,\"data\":{\"component_type\":2,\"custom_id\":\"accept\"}}}")`
			}
			out, failed := hyperRun(t, preview, "turbo")
			if failed || !strings.Contains(out, tc.contains) || tc.name == "components" && !strings.Contains(out, `"content":"Gracias!","flags":64`) {
				t.Fatalf("unexpected preview %q", out)
			}
			if _, err := CompileWeb(preview); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSimpleR2PremiumComponent(t *testing.T) {
	const source = `import tripol.discord
panel=discord.container("#123ABC")
show panel.premium("123456")
show panel.payload()
try {panel.premium("invalid")} catch code {show code}
`
	runCase(t, source, "true\n{\"flags\":32768,\"components\":[{\"type\":17,\"accent_color\":1194684,\"components\":[{\"type\":1,\"components\":[{\"type\":2,\"style\":6,\"sku_id\":\"123456\"}]}]}],\"allowed_mentions\":{\"parse\":[]}}\n12\n", false)
}

func TestSimpleR2ImportedCallbackInference(t *testing.T) {
	dir := t.TempDir()
	entry := filepath.Join(dir, "main.tri")
	if err := os.WriteFile(filepath.Join(dir, "rules.tri"), []byte("func greeting(ctx) {\n ctx.reply(\"From a module\")\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entry, []byte("import tripol.discord\nimport \"rules.tri\"\nb=discord.bot()\nb.slash(\"hello\",\"Greets\",greeting)\nshow \"ready\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cache := &moduleCache{dir: filepath.Join(dir, "cache")}
	for i := 0; i < 2; i++ {
		fs, ss, _ := loadProjectCached(entry, cache)
		compileParsed(fs, ss, "linux", "turbo")
		for _, f := range fs {
			if f.name == "greeting" && (f.args[0].typ != "discord.Interaction" || f.ret != "string") {
				t.Fatal(f.args, f.ret)
			}
		}
	}
}
