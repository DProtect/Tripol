package core

import "testing"

func TestSimpleExpressionsAndLists(t *testing.T) {
 src:=`import tripol.errors
 func twice(n) => n*2
 func side(n) {show n;return n}
 struct Player {score:int=3}
 p=Player()
 a=list([side(1),twice(2),p.score+1])
 show a
 a.add(9)
 show f"Score: {p.score*2}; first: {a[0]}; {twice(4)}"
 show f"{{ok}} {2<3} {1+2*3}"
 show f"{side(5)+side(6)}"
 empty=list([])
 show len(empty)
 try {show f"{1/0}"} catch code {show errors.message(code)}
 try {show a[10]} catch code {show errors.message_es(code)}
 `
 want:="1\n[1, 4, 4]\nScore: 6; first: 1; 8\n{ok} true 7\n5\n6\n11\n0\nInteger overflow or division by zero\nIndice fuera de rango\n"
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);if bad||got!=want{t.Fatalf("%q %v",got,bad)}})
 }
}
func TestSimpleReject(t *testing.T){
 for _,src:=range []string{`a=list(["text"])`,`show f"{}"`,`show f"{1;2}"`,`show f"{1 2}"`,`show f"{1+}"`,`show f"{1.5}"`}{
  t.Run(src,func(t *testing.T){defer func(){if recover()==nil{t.Fatal("accepted invalid program")}}();compile(src,"linux")})
 }
}
