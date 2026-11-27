package core

import (
	"strings"
	"testing"
)

const callbackExample = `import tripol.math

struct Context {username:string}
struct Handler {reply:func(Context)->string}

func greet(ctx:Context)->string {return "Hola, "+ctx.username}
func alternate(ctx:Context)->string {return "Bienvenido, "+ctx.username}
func dispatch(handler:Handler,ctx:Context)->string {
 reply=handler.reply
 return reply(ctx)
}
func twice(n:int)->int {return n*2}
func identity(handler:func(int)->int)->func(int)->int {return handler}
func fib(n:int)->int {
 if n<2 {return n}
 return fib(n-1)+fib(n-2)
}
ctx=Context("Ada")
h=Handler(&greet)
show dispatch(h,ctx)
h.reply=&alternate
show dispatch(h,ctx)
cb=identity(&twice)
show cb(21),twice(3)
square=&math.square
show square(7)
seq=&fib
show seq(10),fib(7)`

func TestTypedCallbacks(t *testing.T) {
	const expected = "Hola, Ada\nBienvenido, Ada\n42 6\n49\n55 13\n"
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			got, failed := hyperRun(t, callbackExample, profile)
			if failed || got != expected {
				t.Fatalf("got %q (failed=%v), want %q", got, failed, expected)
			}
		})
	}
	if _, err := CompileWeb(callbackExample); err != nil {
		t.Fatal(err)
	}
}

func TestTypedCallbacksRejectTypeMismatch(t *testing.T) {
	for _, tc := range []struct{ source, diagnostic string }{
		{`func square(n:int)->int {return n*n}
let cb:func(string)->int=&square`, "expected func(string)->int"},
		{`func greet(n:string)->string {return n}
cb=&greet
show cb(42)`, "expected string"},
		{`func greet(n:string)->string {return n}
cb=&greet
show cb()`, "wrong callback argument count"},
		{`let forged:func(int)->int=42`, "expected func(int)->int"},
		{`struct Record {value:int}
cb=&Record`, "callback target must be a function"},
	} {
		if _, err := CompileWeb(tc.source); err == nil || !strings.Contains(err.Error(), tc.diagnostic) {
			t.Fatalf("%s: got %v, want %q", tc.source, err, tc.diagnostic)
		}
	}
}
