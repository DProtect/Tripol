package core

import "testing"

func TestQuantumTemporaryRegions(t *testing.T) {
	src := `import tripol.memory
prefix="id:"
before=memory.used()
total=0
repeat 100000 {
 temp=prefix+text(17)
 sub=slice(temp,0,3)
 total+=len(sub)
}
show total,before==memory.used()
repeat 100000 {
 data=list(0)
 alias=data
 data.add(5)
 alias.add(7)
 total+=alias[1]
}
show total,before==memory.used()
try {
 repeat 3 {
  temp=text(123)
  assert len(temp)==0
 }
} catch err {show err}
show before==memory.used()
outer="a"
repeat 2 {outer=outer+"b"}
show outer
`
	for _, p := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(p, func(t *testing.T) {
			got, bad := hyperRun(t, src, p)
			if bad || got != "300000 true\n1000000 true\n12\ntrue\nabb\n" {
				t.Fatalf("%q %v", got, bad)
			}
		})
	}
}

func TestNestedTemporaryListRegion(t *testing.T) {
	source := `import tripol.memory
 before=memory.used()
 total=0
 repeat 100000 {
  messages=string_list(0)
  messages.add(text(17))
  for j in range(3) {total+=len(messages[0])}
 }
 show total,before==memory.used()
 `
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		got, failed := hyperRun(t, source, profile)
		if failed || got != "600000 true\n" {
			t.Fatalf("%s: %q", profile, got)
		}
	}
}
