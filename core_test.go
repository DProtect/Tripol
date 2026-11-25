package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// Frozen v0.15 scanner is a test-only grammar oracle.
var oldLexRE = regexp.MustCompile(`^(?:([ \t\r]+)|(//[^\n]*|#[^\n]*)|(\n)|("(?:[^"\\\n]|\\.)*")|([0-9]+(?:\.[0-9]+)?)|([A-Za-z_][A-Za-z_0-9]*)|(\+=|-=|\*=|/=|%=|==|!=|<=|>=|&&|\|\||->|=>|[+*/%<>=!(){}\[\],.:;\-]))`)

func oldLex(s string) []token {
	var ts []token
	line, col := 1, 1
	for len(s) > 0 {
		m := oldLexRE.FindStringSubmatch(s)
		if m == nil {
			lexError(s[0], line, col)
		}
		k := 0
		for i := 1; i < len(m); i++ {
			if m[i] != "" {
				k = i
				break
			}
		}
		if k > 2 {
			ts = append(ts, token{m[0], []string{"", "", "", "nl", "str", "int", "id", "op"}[k], line, col})
		}
		if k == 3 {
			line++
			col = 1
		} else {
			col += len(m[0])
		}
		s = s[len(m[0]):]
	}
	return append(ts, token{"<eof>", "eof", line, col})
}
func captureLex(f func(string) []token, s string) (tokens []token, err string) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Sprintf("%#v", e)
		}
	}()
	tokens = f(s)
	return
}
func TestCoreLexerParity(t *testing.T) {
	r := rand.New(rand.NewSource(1616))
	alphabet := []byte("aZ_109.\t\r\n /#\"\\+-*%!=&|><:;()[]{},\x00\xffé你好")
	cases := []string{"", "// hi\nshow f\"Hola {x}\"", `"escaped \" and \\ and é"`, "\"a\\é\"", "1.2.3 1. 99 => -> /= // end", "\"a\\\n\"", "\"\xff\"", "show\r\n3"}
	for i := 0; i < 20000; i++ {
		b := make([]byte, r.Intn(100))
		for j := range b {
			b[j] = alphabet[r.Intn(len(alphabet))]
		}
		cases = append(cases, string(b))
	}
	for _, s := range cases {
		a, ae := captureLex(oldLex, s)
		b, be := captureLex(lex, s)
		if ae != be || !reflect.DeepEqual(a, b) {
			t.Fatalf("%q: old %v %v new %v %v", s, a, ae, b, be)
		}
	}
}
func TestCoreSignedDivision(t *testing.T) {
	var src, w strings.Builder
	for _, x := range []int64{math.MinInt64, math.MinInt64 + 1, -1025, -17, -8, -7, -1, 0, 1, 7, 8, 17, 1025, math.MaxInt64} {
		lit := fmt.Sprint(x)
		if x == math.MinInt64 {
			lit = "(-9223372036854775807-1)"
		}
		fmt.Fprintf(&src, "x=%s\n", lit)
		for _, n := range []int64{1, 2, 4, 8, 16, 256, 1 << 30, 3, 0x7fffffff} {
			fmt.Fprintf(&src, "show x/%d,x%%%d\n", n, n)
			fmt.Fprintf(&w, "%d %d\n", x/n, x%n)
		}
	}
	runCase(t, src.String(), w.String(), false)
}
func TestCoreMemoryOperands(t *testing.T) {
	runCase(t, `func calc(a,b){return (a+b)*(a-b)}
show calc(9,4)
func wide(a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q){return q-a}
show wide(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17)
x=2;y=3;z=4
show x+y,x-y,y*z
sum=0
for i in range(9223372036854775805,9223372036854775807) {sum+=1}
show sum
for i in range(5,2){show i}
for i in range(0,5){if i==2{continue};if i==4{break};show i}`, "65\n16\n5 -1 12\n2\n0\n1\n3\n", false)
	runCase(t, `x=9223372036854775807;y=1
try {x=x+y} catch code {show code}
show x`, "10\n9223372036854775807\n", false)
}
func TestCoreDivisionEvaluation(t *testing.T) {
	runCase(t, `func side(){show "called";return 19}
show side()%1
show side()/1`, "called\n0\ncalled\n19\n", false)
}
func TestCoreFixIndex(t *testing.T) {
	a := asm{}
	a.addFix(fix{0, "one"})
	a.addFix(fix{4, "two"})
	if !a.needed("one") || !a.needed("two") || a.needed("three") {
		t.Fatal("incorrect relocation index")
	}
}
