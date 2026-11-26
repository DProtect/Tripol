package core

// WEB ENGINE R1: native frontend + an independent, bounded WebAssembly target.
// Numeric overflow and console I/O are checked by the browser host ABI.
import (
 "fmt"
)

type webCode []byte
func (b *webCode) u(n uint32){for {v:=byte(n&127);n>>=7;if n!=0{v|=128};*b=append(*b,v);if n==0{return}}}
func (b *webCode) signed(n int64){for {v:=byte(n&127);n>>=7;done:=(n==0&&v&64==0)||(n== -1&&v&64!=0);if !done{v|=128};*b=append(*b,v);if done{return}}}
func (b *webCode) name(s string){b.u(uint32(len(s)));*b=append(*b,[]byte(s)...)}
func (b *webCode) raw(v ...byte){*b=append(*b,v...)}
func (b *webCode) inst(op byte,n uint32){b.raw(op);b.u(n)}
func (b *webCode) i64(n int64){b.raw(0x42);b.signed(n)}
func (b *webCode) i32(n int64){b.raw(0x41);b.signed(n)}
func (b *webCode) section(id byte,p webCode){b.raw(id);b.u(uint32(len(p)));*b=append(*b,p...)}
type webLocal struct { index uint32; typ string }
type webSig struct { args []byte; results []byte }
type webFunction struct { f fn; index uint32; typ uint32 }
type webModule struct { funcs map[string]webFunction; order []webFunction; types []webSig; bodies []webCode; current int }
type webEmitter struct { module *webModule; code webCode; scopes []map[string]webLocal; locals,params uint32; line int; labels []string; result string }
func webReject(format string,args ...any){die("Web Engine R1: "+format,args...)}
func webScalar(t string)bool{return t=="int"||t=="bool"}
func (m *webModule) signature(a,r []byte)uint32{m.types=append(m.types,webSig{a,r});return uint32(len(m.types)-1)}
func (m *webModule) find(name string)webFunction{
 if v,ok:=m.funcs[name];ok{return v};webReject("unsupported function %s",name);return webFunction{}
}
func (m *webModule) declare(f fn)webFunction{
 if v,ok:=m.funcs[f.name];ok{return v}
 if f.record||!webScalar(f.ret){webReject("function %s requires a scalar return type",f.name)}
 args:=make([]byte,len(f.args));for i,a:=range f.args{if !webScalar(a.typ){webReject("parameter %s requires int or bool",a.name)};args[i]=0x7e}
 v:=webFunction{f:f,index:4+uint32(len(m.order)),typ:m.signature(args,[]byte{0x7e})};m.funcs[f.name]=v;m.order=append(m.order,v);return v
}
func (e *webEmitter) push(){e.scopes=append(e.scopes,map[string]webLocal{})}
func (e *webEmitter) pop(){e.scopes=e.scopes[:len(e.scopes)-1]}
func (e *webEmitter) get(n string)(webLocal,bool){for i:=len(e.scopes)-1;i>=0;i--{if v,ok:=e.scopes[i][n];ok{return v,true}};return webLocal{},false}
func (e *webEmitter) local(n,t string)webLocal{v:=webLocal{e.locals,t};e.locals++;e.scopes[len(e.scopes)-1][n]=v;return v}
func (e *webEmitter) tick(){e.code.i32(int64(e.line));e.code.inst(0x10,0)}
func (e *webEmitter) checked(op int64,a,b *expr)string{e.code.i32(op);ta:=e.expression(a);tb:=e.expression(b);if ta!="int"||tb!="int"{webReject("integer arithmetic required at line %d",e.line)};e.code.i32(int64(e.line));e.code.inst(0x10,1);return "int"}
func (e *webEmitter) expression(x *expr)string{
 if x==nil{webReject("missing expression")};switch x.op{
 case "int","bool":e.code.i64(x.n);return x.op
 case "var":v,ok:=e.get(x.s);if !ok{webReject("unknown variable %s",x.s)};e.code.inst(0x20,v.index);return v.typ
 case "un":
  if x.s=="-"{return e.checked(1,&expr{op:"int"},x.args[0])}
  if x.s=="!"{if e.expression(x.args[0])!="bool"{webReject("expected bool")};e.code.raw(0x50,0xad);return "bool"}
 case "bin":
  switch x.s{
  case "+","-","*","/","%":return e.checked(map[string]int64{"+":0,"-":1,"*":2,"/":3,"%":4}[x.s],x.args[0],x.args[1])
  case "==","!=","<",">","<=",">=":
   a:=e.expression(x.args[0]);b:=e.expression(x.args[1]);if a!=b{webReject("comparison type mismatch")};e.code.raw(map[string]byte{"==":0x51,"!=":0x52,"<":0x53,">":0x55,"<=":0x57,">=":0x59}[x.s],0xad);return "bool"
  case "&&","||":
   if e.expression(x.args[0])!="bool"{webReject("expected bool")};e.code.raw(0x50,0x04,0x7e)
   if x.s=="&&"{e.code.i64(0);e.code.raw(0x05);if e.expression(x.args[1])!="bool"{webReject("expected bool")}}else{if e.expression(x.args[1])!="bool"{webReject("expected bool")};e.code.raw(0x05);e.code.i64(1)}
   e.code.raw(0x0b);return "bool"
  }
 case "call":
  f:=e.module.find(x.s);if len(x.args)>len(f.f.args){webReject("wrong argument count")};for i,a:=range f.f.args{var arg *expr;if i<len(x.args){arg=x.args[i]}else if i<len(f.f.defaults){arg=f.f.defaults[i]};if arg==nil||e.expression(arg)!=a.typ{webReject("invalid argument for %s",x.s)}}
  e.code.inst(0x10,f.index);return f.f.ret
 }
 webReject("unsupported expression %s (%s) at line %d",x.op,x.s,e.line);return ""
}
func (e *webEmitter) begin(op byte,label string){e.code.raw(op,0x40);e.labels=append(e.labels,label)}
func (e *webEmitter) end(){e.code.raw(0x0b);e.labels=e.labels[:len(e.labels)-1]}
func (e *webEmitter) branch(label string){for i:=len(e.labels)-1;i>=0;i--{if e.labels[i]==label{e.code.inst(0x0c,uint32(len(e.labels)-1-i));return}};webReject("invalid %s",label)}
func (e *webEmitter) block(ss []stmt,nested bool){if nested{e.push();defer e.pop()};for _,s:=range ss{e.line=s.line;e.tick();switch s.op{
 case "decl","assign":
  t:=e.expression(s.e);v,ok:=e.get(s.name);if s.op=="decl"||!ok{v=e.local(s.name,t)};if v.typ!=t{webReject("assignment type mismatch")};e.code.inst(0x21,v.index)
 case "show","showmany":
  args:=[]*expr{s.e};if s.op=="showmany"{args=s.e.args};for i,a:=range args{t:=e.expression(a);kind:=int64(0);if t=="bool"{kind=1};e.code.i32(kind);last:=int64(0);if i==len(args)-1{last=1};e.code.i32(last);e.code.inst(0x10,2)}
 case "expr":e.expression(s.e);e.code.raw(0x1a)
 case "return":if e.result==""||e.expression(s.e)!=e.result{webReject("invalid return")};e.code.raw(0x0f)
 case "assert":if e.expression(s.e)!="bool"{webReject("expected bool")};e.code.raw(0x50);e.begin(0x04,"if");e.code.i32(int64(s.line));e.code.inst(0x10,3);e.end()
 case "if":if e.expression(s.e)!="bool"{webReject("expected bool")};e.code.raw(0x50,0x45);e.begin(0x04,"if");e.block(s.yes,true);if len(s.no)>0{e.code.raw(0x05);e.block(s.no,true)};e.end()
 case "while":
  e.begin(0x02,"break");e.begin(0x03,"continue");e.tick();if e.expression(s.e)!="bool"{webReject("expected bool")};e.code.raw(0x50);e.code.inst(0x0d,1);e.block(s.yes,true);e.code.inst(0x0c,0);e.end();e.end()
 case "for":
  e.push();if e.expression(s.e)!="int"{webReject("range requires integers")};v:=e.local(s.name,"int");e.code.inst(0x21,v.index);if e.expression(s.bound)!="int"{webReject("range requires integers")};bound:=e.local(fmt.Sprintf("@webbound%d",e.locals),"int");e.code.inst(0x21,bound.index)
  e.begin(0x02,"break");e.begin(0x03,"loop");e.tick();e.code.inst(0x20,v.index);e.code.inst(0x20,bound.index);e.code.raw(0x59);e.code.inst(0x0d,1)
  e.begin(0x02,"continue");e.block(s.yes,true);e.end()
  e.code.i32(0);e.code.inst(0x20,v.index);e.code.i64(1);e.code.i32(int64(s.line));e.code.inst(0x10,1);e.code.inst(0x21,v.index);e.code.inst(0x0c,0);e.end();e.end();e.pop()
 case "break","continue":e.branch(s.op)
 default:webReject("statement %s is not supported at line %d",s.op,s.line)
 }}}
func (e *webEmitter) body()webCode{var b webCode;n:=e.locals-e.params;if n==0{b.u(0)}else{b.u(1);b.u(n);b.raw(0x7e)};b=append(b,e.code...);b.raw(0x0b);return b}

// CompileWeb validates the program with the production frontend, then emits
// a standard WASM module. R1 accepts int/bool, functions, control flow and
// selected pure integer stdlib functions. Unsupported constructs are rejected.
func CompileWeb(source string)(out []byte,err error){
 defer func(){if v:=recover();v!=nil{out=nil;err=fmt.Errorf("%v",v)}}()
 if len(source)>32768{return nil,fmt.Errorf("Web Engine: 32 KiB source limit")}
 imports,_,_:=parseUnit(source);for _,n:=range imports{if n!="@std:math"{webReject("only tripol.math is supported by this target")}}
 fs,ss:=parse(source);validateCode(fs,ss,"linux")
 m:=&webModule{funcs:map[string]webFunction{}}
 m.signature([]byte{0x7f},nil) // tick(line)
 m.signature([]byte{0x7f,0x7e,0x7e,0x7f},[]byte{0x7e}) // checked arithmetic
 m.signature([]byte{0x7e,0x7f,0x7f},nil) // show
 m.signature([]byte{0x7f},nil) // assert failure
 // Register only reachable functions; imported float functions need not be emitted.
 byName:=map[string]fn{};for _,f:=range fs{byName[f.name]=f}
 seen:=map[string]bool{}
 var walkExpr func(*expr);var walkBlock func([]stmt)
 walkExpr=func(x *expr){if x==nil{return};if x.op=="call"&&!seen[x.s]{seen[x.s]=true;if f,ok:=byName[x.s];ok{m.declare(f);walkBlock(f.body);for _,d:=range f.defaults{walkExpr(d)}}};for _,a:=range x.args{walkExpr(a)}}
 walkBlock=func(ss []stmt){for _,s:=range ss{walkExpr(s.e);walkExpr(s.bound);walkBlock(s.yes);walkBlock(s.no);walkBlock(s.final)}}
 walkBlock(ss)
 for _,f:=range m.order{e:=&webEmitter{module:m,result:f.f.ret};e.push();for _,a:=range f.f.args{e.local(a.name,a.typ)};e.params=e.locals;e.block(f.f.body,false);e.code.raw(0x00);m.bodies=append(m.bodies,e.body())}
 mainType:=m.signature(nil,nil);mainIndex:=uint32(4+len(m.order));e:=&webEmitter{module:m};e.push();e.block(ss,false);m.bodies=append(m.bodies,e.body())
 var module webCode;module.raw(0,97,115,109,1,0,0,0)
 var ts webCode;ts.u(uint32(len(m.types)));for _,s:=range m.types{ts.raw(0x60);ts.u(uint32(len(s.args)));ts.raw(s.args...);ts.u(uint32(len(s.results)));ts.raw(s.results...)};module.section(1,ts)
 var im webCode;im.u(4);for i,n:=range []string{"tick","arithmetic","show","assertFail"}{im.name("tripol");im.name(n);im.raw(0);im.u(uint32(i))};module.section(2,im)
 var funcs webCode;funcs.u(uint32(len(m.order)+1));for _,f:=range m.order{funcs.u(f.typ)};funcs.u(mainType);module.section(3,funcs)
 var exports webCode;exports.u(1);exports.name("run");exports.raw(0);exports.u(mainIndex);module.section(7,exports)
 var code webCode;code.u(uint32(len(m.bodies)));for _,b:=range m.bodies{code.u(uint32(len(b)));code.raw(b...)};module.section(10,code)
 return []byte(module),nil
}
