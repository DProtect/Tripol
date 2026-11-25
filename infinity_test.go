package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfinitySyntax(t *testing.T) {
	runCase(t, `func greet(name="Guest",score=10,active=true)->string {return f"{name}: {score} ({active}) {{ok}}"}
func num(x:int=2+3,y:int=-2)->int {return x+y}
func decimal(x=-1.5)->float {return x}
show greet()
show greet("Ana",21,false)
show num()
show num(9)
show decimal()
show "{not interpolated}"
show f""
show f"{{}}"`, "Guest: 10 (true) {ok}\nAna: 21 (false) {ok}\n3\n7\n-1.500000\n{not interpolated}\n\n{}\n", false)
}
func TestInfinityReject(t *testing.T) {
	for _, s := range []string{`func f(x=1,y){return x+y}`, `func f(x:int="a"){return x}`, `func f(x=1/0){return x}`, `func f(x=unknown){return x}`, `func f(x){return x};show f()`, `func f(x=1){return x};show f(1,2)`, `show f"{unknown}"`, `show f"{a+1}"`, `show f"{"`, `show f"}"`, `x=1.0;show f"{x}"`} {
		t.Run(s, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid code")
				}
			}()
			compile(s, "linux")
		})
	}
}
func TestInfinityCollections(t *testing.T) {
	runCase(t, `import tripol.collections
import tripol.compression
import tripol.sequences
import tripol.matrix
a=[1,1,1,4,4,9]
b=collections.append(a,9)
show b
show collections.insert(b,0,8)
show collections.remove_at(b,2)
show collections.part(b,2,4)
show collections.concat(a,b)
enc=compression.rle_encode(b)
show enc
out=compression.rle_decode(enc)
show collections.equal(b,out)
show compression.rle_valid(a)
empty=ints(0)
show compression.rle_encode(empty)
show compression.rle_decode(empty)
show matrix.identity(3)
show matrix.transpose(a,2,3)
show matrix.dot(a,a)
show matrix.add(a,a)
show sequences.arithmetic(5,2,3)
show sequences.fibonacci(10)
show sequences.factorial(20)
show sequences.triangular(4294967295)`, "[1, 1, 1, 4, 4, 9, 9]\n[8, 1, 1, 1, 4, 4, 9, 9]\n[1, 1, 4, 4, 9, 9]\n[1, 4]\n[1, 1, 1, 4, 4, 9, 1, 1, 1, 4, 4, 9, 9]\n[1, 3, 4, 2, 9, 2]\ntrue\ntrue\n[]\n[]\n[1, 0, 0, 0, 1, 0, 0, 0, 1]\n[1, 4, 1, 4, 1, 9]\n116\n[2, 2, 2, 8, 8, 18]\n[2, 5, 8, 11, 14]\n[0, 1, 1, 2, 3, 5, 8, 13, 21, 34]\n2432902008176640000\n9223372034707292160\n", false)
}
func TestInfinityNumbers(t *testing.T) {
	runCase(t, `import tripol.colors
import tripol.datetime
import tripol.bits
import tripol.checksum
import tripol.geometry
c=colors.rgb(255,128,0)
show colors.red(c),colors.green(c),colors.blue(c),colors.invert(c)
show colors.mix(0,16777215,50)
show colors.gray(128)
show datetime.valid(1900,2,29),datetime.valid(2000,2,29)
show datetime.weekday(2026,11,24)
show datetime.days_between(2024,2,28,2024,3,1)
show datetime.day_of_year(2024,12,31)
show bits.set(0,62),bits.clear(7,1),bits.toggle(7,1),bits.count(255)
show checksum.adler32("Wikipedia"),checksum.sum8("abc"),checksum.sum16("abc")
show geometry.distance(0.0,0.0,3.0,4.0),geometry.rectangle(2.0,3.0),geometry.lerp(0.0,10.0,0.5)`, "255 128 0 32767\n8355711\n8421504\nfalse true\n2\n2\n366\n4611686018427387904 5 5 8\n300286872 38 294\n5.000000 6.000000 5.000000\n", false)
}
func TestInfinityAllLibraries(t *testing.T) {
	var src strings.Builder
	for _, name := range standardNames {
		src.WriteString("import tripol." + name + "\n")
	}
	src.WriteString(`show "Hello, TRIPOL!"`)
	for _, target := range []string{"linux", "windows"} {
		if !bytes.Equal(compile(src.String(), target), compile(`show "Hello, TRIPOL!"`, target)) {
			t.Fatal("unused libraries changed binary", target)
		}
	}
	if len(standardNames) != 40 {
		t.Fatal(len(standardNames))
	}
}
func TestInfinityOptimizations(t *testing.T) {
	runCase(t, `x=1
repeat 100 {x=x*1+0-0}
show x,len("hello"),len("é")`, "1 5 2\n", false)
	runCase(t, `func f(x){return 10/x}
try {show f(0)*1+0} catch e {show e}`, "10\n", false)
}

func TestTIRInline(t *testing.T) {
	runCase(t, `func sq(n){return n*n}
func affine(x,y,z=2){return x*y+z}
func ignore(x){return 9}
func side(x){show x;return x}
show sq(9)
show affine(side(3),side(4))
show ignore(side(5))
try {show sq(9223372036854775807)} catch e {show e}
try {show ignore(1/0)} catch e {show e}
show affine(2,3,4)`, "81\n3\n4\n14\n5\n9\n10\n10\n10\n", false)
	fs, _ := parse(`func f(a,b){return a*b+1}`)
	if ir, ok := lowerTIR(fs[0]); !ok || ir.root.op != tirAdd {
		t.Fatal("TIR not lowered")
	}
	fs, _ = parse(`func f(a){return a/2}`)
	if _, ok := lowerTIR(fs[0]); ok {
		t.Fatal("unsupported expression lowered")
	}
}

func TestInfinityRecovery(t *testing.T) {
	runCase(t, `import tripol.colors
import tripol.compression
import tripol.collections
import tripol.datetime
import tripol.matrix
import tripol.bits
bad=[3,0]
show compression.rle_valid(bad)
try {show compression.rle_decode(bad)} catch e {show e}
try {show colors.rgb(256,0,0)} catch e {show e}
try {show collections.remove_at(bad,2)} catch e {show e}
try {show datetime.ordinal(1900,2,29)} catch e {show e}
try {show matrix.zeros(1000,1000)} catch e {show e}
try {show bits.mask(63)} catch e {show e}`, "false\n12\n12\n12\n12\n12\n12\n", false)
}
func TestInfinityPE(t *testing.T) {
	src := `import tripol.ui
show ui.warning("Test","Warning")
show ui.error("Test","Error")
show ui.retry("Test","Retry")
show ui.choose("Test","Choose")`
	bin := compile(src, "windows")
	if !bytes.Contains(bin, []byte("MessageBoxA")) || !bytes.Contains(bin, []byte("USER32.dll")) {
		t.Fatal("missing Win32 import")
	}
}
