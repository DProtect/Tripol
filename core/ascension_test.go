package core

import "testing"

func TestAscensionCopyInto(t *testing.T) {
	src := `import tripol.arrays
import tripol.memory
a=[1,-2,3];b=ints(3)
show arrays.copy_into(b,a),b
b[0]=8
show a,b
show arrays.copy_into(a,a),a
short=[7,9]
show arrays.copy_into(short,a),short
empty=ints(0)
show arrays.copy_into(empty,empty)
before=memory.used()
repeat 2000 {assert arrays.copy_into(b,a)}
show before==memory.used(),b`
	for _, p := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(p, func(t *testing.T) {
			got, bad := hyperRun(t, src, p)
			want := "true [1, -2, 3]\n[1, -2, 3] [8, -2, 3]\ntrue [1, -2, 3]\nfalse [7, 9]\ntrue\ntrue [1, -2, 3]\n"
			if bad || got != want {
				t.Fatalf("%q %v", got, bad)
			}
		})
	}
}
func TestAscensionPercentile(t *testing.T) {
	if arenaPercentile([]float64{5, 1, 3, 2, 4}, .95) != 5 {
		t.Fatal("p95")
	}
	if arenaPercentile([]float64{3}, .95) != 3 {
		t.Fatal("single")
	}
}
