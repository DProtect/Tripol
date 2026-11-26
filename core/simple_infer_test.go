package core
import "testing"
func TestSimpleInferredFunctions(t *testing.T){
 src:=`func greet(name) {return f"Hola, {name}"}
 func welcome(name) => "Welcome "+name
 func twice(n) => n*2
 func opposite(flag) => not flag
 func identity(x) => x
 names=list(["Ana","Carlos"])
 for name in names {show greet(name)}
 show welcome("Eva"),twice(7),opposite(true),identity("hello")
 `
 want:="Hola, Ana\nHola, Carlos\nWelcome Eva 14 false hello\n"
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);if bad||got!=want{t.Fatalf("%q %v",got,bad)}})}
}

func TestSimpleExplicitTypeStillChecked(t *testing.T){
 for _,src:=range []string{`func f()->int {return true};show f()`,`func f(x:int) => x;show f("no")`,`func f(x) => x;show f("a"),f(2)`}{
  t.Run(src,func(t *testing.T){defer func(){if recover()==nil{t.Fatal("accepted conflicting explicit or inferred signature")}}();compile(src,"linux")})
 }
}
