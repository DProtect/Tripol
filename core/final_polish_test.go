package core

import (
	"encoding/json"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func TestFinalPolishArgumentsOptionsAndTypedJSON(t *testing.T) {
	const source = `import tripol.json
struct Player {name:string;health:int=100;speed:float=1.5;score:int?=null}
func Player.net(self:Player,tax:int=1,bonus:int=5)->int {return self.health-tax+bonus}
func next(values:ints, n:int)->int {values[0]=values[0]*10+n;return n}
func points(first:int=1,second:int=2)->int {return first*10+second}
func pick(value:int, extra:int=3)->int {return value+extra}
func pick(text:string, suffix:string="?")->string {return text+suffix}
values=ints(1)
player=json.decode(Player,"{\"name\":\"Alex\"}")
show player.name,player.health,player.speed,player.score.has()
show json.encode(Player,player)
another=Player(speed: 2.5, name: "Eva", score: some(9))
show another.net(bonus:9),another.score.or_else(1)
show points(second:next(values,2),first:next(values,1)),values[0]
show pick(suffix:"!",text:"ok"),pick(extra:4,value:5)
players=list([list([player,another])])
show players[0][1].name
blank=json.decode(Player,"{\"name\":\"A\",\"score\":null,\"speed\":1.25e1}")
show blank.score.or_else(11),blank.speed
`
	want := "Alex 100 1.500000 false\n{\"name\":\"Alex\",\"health\":100,\"speed\":1.5,\"score\":null}\n108 9\n12 21\nok! 9\nEva\n11 12.500000\n"
	runCase(t, source, want, false)
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
}

func TestFinalPolishDecimalRoundTrip(t *testing.T) {
	decimals := []string{"1.015", "-333.25", "-0", "0.1", "9007199254740993",
		"2.2250738585072012e-308", "4.9406564584124654e-324", "2.4703282292062328e-324",
		"2.4703282292062327e-324", "1.7976931348623157e308", "1.00000000000000011102230246251565404236316680908203125",
		"1.00000000000000011102230246251565404236316680908203125000001"}
	var source, expected strings.Builder
	source.WriteString("import tripol.json\nstruct Number {value:float}\n")
	for _, decimal := range decimals {
		input := "{\"value\":" + decimal + "}"
		source.WriteString("show json.encode(Number,json.decode(Number," + strconv.Quote(input) + "))\n")
		v, err := strconv.ParseFloat(decimal, 64)
		if err != nil {
			t.Fatal(err)
		}
		var output string
		if v == 0 && math.Signbit(v) {
			output = "-0"
		} else {
			rat := new(big.Rat).SetFloat64(v)
			places := rat.Denom().BitLen() - 1
			output = rat.FloatString(places)
		}
		if !json.Valid([]byte(output)) {
			t.Fatal(output)
		}
		expected.WriteString("{\"value\":" + output + "}\n")
	}
	runCase(t, source.String(), expected.String(), false)
}

func TestFinalPolishNamedArgumentDiagnostics(t *testing.T) {
	for _, tc := range []struct{ call, message string }{
		{"g(vaule: 1)", "unknown named argument"},
		{"g(1, value: 2)", "duplicate argument"},
		{"g(extra: 1)", "missing required argument"},
		{"g(value: 1, 4)", "positional arguments before named"},
	} {
		t.Run(tc.call, func(t *testing.T) {
			_, err := CompileWeb("func g(value:int,extra:int=2)->int {return value+extra}\nshow " + tc.call)
			if err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatal(err)
			}
		})
	}
}

func TestFinalPolishClosureSnapshotsAndLifetime(t *testing.T) {
	const source = `func make(delta:int)->func(int)->int {
    return func(x:int) {x+delta}
}
func plus(x:int)->int {return x+1}
func apply(f:func(int)->int,n:int)->int {return f(n)}
number=2
first=func(n:int) {n+number}
number=9
second=make(number)
callbacks=list([first,second,&plus])
for cb in callbacks {show apply(cb,10)}
func nested(base:int)->int {
    outer=func(seed:int)->int {
        inner=func(extra:int) {base+seed+extra}
        return inner(1)
    }
    return outer(2)
}
show nested(3)
`
	runCase(t, source, "12\n19\n11\n6\n", false)
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
}

func TestFinalPolishClosureCannotEscapeScratch(t *testing.T) {
	_, err := CompileWeb(`func make(delta:int)->func(int)->int {return func(x:int) {x+delta}}
callback=make(1)
scratch {callback=make(2)}`)
	if err == nil || !strings.Contains(err.Error(), "scratch cannot assign an outer callback") {
		t.Fatal(err)
	}
}
