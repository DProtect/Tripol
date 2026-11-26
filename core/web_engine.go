package core

// WEB ENGINE R3: native frontend + an independent, bounded WebAssembly target.
// Numeric overflow and console I/O are checked by the browser host ABI.
import (
 "fmt"
 "encoding/json"
 "strings"
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
type webModule struct { funcs map[string]webFunction; order []webFunction; types []webSig; bodies []webCode; current int; strings []string; stringIDs map[string]int64 }
type webEmitter struct { module *webModule; code webCode; scopes []map[string]webLocal; locals,params uint32; line int; labels []string; result string }
func webReject(format string,args ...any){die("Web Engine R3: "+format,args...)}
func webScalar(t string)bool{return t=="int"||t=="bool"||t=="string"||t=="float"||t=="ints"||isList(t)}
func (m *webModule) signature(a,r []byte)uint32{m.types=append(m.types,webSig{a,r});return uint32(len(m.types)-1)}
func (m *webModule) find(name string)webFunction{
 if v,ok:=m.funcs[name];ok{return v};webReject("unsupported function %s",name);return webFunction{}
}
func (m *webModule) declare(f fn)webFunction{
 if v,ok:=m.funcs[f.name];ok{return v}
 if f.record||!webScalar(f.ret){webReject("function %s requires a scalar return type",f.name)}
 args:=make([]byte,len(f.args));for i,a:=range f.args{if !webScalar(a.typ){webReject("parameter %s requires a supported type",a.name)};args[i]=0x7e}
 v:=webFunction{f:f,index:5+uint32(len(m.order)),typ:m.signature(args,[]byte{0x7e})};m.funcs[f.name]=v;m.order=append(m.order,v);return v
}
func (e *webEmitter) push(){e.scopes=append(e.scopes,map[string]webLocal{})}
func (e *webEmitter) pop(){e.scopes=e.scopes[:len(e.scopes)-1]}
func (e *webEmitter) get(n string)(webLocal,bool){for i:=len(e.scopes)-1;i>=0;i--{if v,ok:=e.scopes[i][n];ok{return v,true}};return webLocal{},false}
func (e *webEmitter) local(n,t string)webLocal{v:=webLocal{e.locals,t};e.locals++;e.scopes[len(e.scopes)-1][n]=v;return v}
func (e *webEmitter) tick(){e.code.i32(int64(e.line));e.code.inst(0x10,0)}
func (e *webEmitter) checked(op int64,a,b *expr)string{e.code.i32(op);ta:=e.expression(a);tb:=e.expression(b);if ta!="int"||tb!="int"{webReject("integer arithmetic required at line %d",e.line)};e.code.i32(int64(e.line));e.code.inst(0x10,1);return "int"}
func (m *webModule) stringID(s string)int64{
 if id,ok:=m.stringIDs[s];ok{return id}
 id:=int64(len(m.strings));m.strings=append(m.strings,s);m.stringIDs[s]=id;return id
}
func webKind(t string)int64{switch t{case "int":return 0;case "bool":return 1;case "string":return 2;case "float":return 3};if t=="ints"||strings.HasPrefix(t,"array:"){return 4};if isList(t){return 5+webKind(listElement(t))};webReject("unsupported value type %s",t);return 0}
func webElem(t string)string{if t=="ints"||strings.HasPrefix(t,"array:"){return "int"};return listElement(t)}
func (e *webEmitter) temp()uint32{v:=e.locals;e.locals++;return v}
func (e *webEmitter) save(x *expr)(uint32,string){t:=e.expression(x);v:=e.temp();e.code.inst(0x21,v);return v,t}
func (e *webEmitter) val(op int64,a,b,c uint32){e.code.i32(op);e.code.inst(0x20,a);e.code.inst(0x20,b);e.code.inst(0x20,c);e.code.i32(int64(e.line));e.code.inst(0x10,4)}
func (e *webEmitter) constant(n int64)uint32{v:=e.temp();e.code.i64(n);e.code.inst(0x21,v);return v}
func (e *webEmitter) arrayLiteral(x *expr,list bool,requested string)string{
 zero:=e.constant(0);size:=e.constant(int64(len(x.args)));kind:=int64(4);typ:="array:"+fmt.Sprint(len(x.args))
 if list{typ=requested;if len(x.args)>0&&requested=="list"{v,t:=e.save(x.args[0]);typ=listType(t);if typ==""{webReject("list requires scalar")};kind=webKind(typ);k:=e.constant(kind);e.val(10,size,k,zero);handle:=e.temp();e.code.inst(0x21,handle);e.val(14,handle,v,zero);e.code.raw(0x1a);for _,item:=range x.args[1:]{v,t:=e.save(item);if t!=listElement(typ){webReject("mixed list")};e.val(14,handle,v,zero);e.code.raw(0x1a)};e.code.inst(0x20,handle);return typ};kind=webKind(typ)}
 k:=e.constant(kind);e.val(10,size,k,zero);handle:=e.temp();e.code.inst(0x21,handle)
 for i,item:=range x.args{v,t:=e.save(item);elem:=webElem(typ);if t!=elem{webReject("invalid collection element")};if list{e.val(14,handle,v,zero)}else{idx:=e.constant(int64(i));e.val(12,handle,idx,v)};e.code.raw(0x1a)}
 e.code.inst(0x20,handle);return typ
}
func (e *webEmitter) expression(x *expr)string{
 if x==nil{webReject("missing expression")};switch x.op{
 case "int","bool","float":e.code.i64(x.n);return x.op
 case "string":e.code.i64(e.module.stringID(x.s));return "string"
 case "array":return e.arrayLiteral(x,false,"")
 case "format":
  zero:=e.constant(0);acc:=e.constant(e.module.stringID(""));for _,part:=range x.args{v,t:=e.save(part);if t!="string"&&t!="int"&&t!="bool"{webReject("unsupported interpolation type")};k:=e.constant(webKind(t));e.val(21,v,k,zero);v=e.temp();e.code.inst(0x21,v);e.code.i32(5);e.code.inst(0x20,acc);e.code.inst(0x20,v);e.code.i32(int64(e.line));e.code.inst(0x10,1);e.code.inst(0x21,acc)};e.code.inst(0x20,acc);return "string"
 case "index":v,ok:=e.get(x.s);if !ok||webElem(v.typ)==""{webReject("unsupported index target %s",x.s)};idx,t:=e.save(x.args[0]);if t!="int"{webReject("index must be int")};e.val(11,v.index,idx,e.constant(0));return webElem(v.typ)
 case "var":v,ok:=e.get(x.s);if !ok{webReject("unknown variable %s",x.s)};e.code.inst(0x20,v.index);return v.typ
 case "un":
  if x.s=="-"{v,t:=e.save(x.args[0]);zero:=e.constant(0);if t=="float"{e.val(33,v,zero,zero);return "float"};if t!="int"{webReject("numeric operand required")};e.code.i32(1);e.code.i64(0);e.code.inst(0x20,v);e.code.i32(int64(e.line));e.code.inst(0x10,1);return "int"}
  if x.s=="!"{if e.expression(x.args[0])!="bool"{webReject("expected bool")};e.code.raw(0x50,0xad);return "bool"}
 case "bin":
  switch x.s{
  case "+":
   // Evaluate operands once, left to right; keep them in unnamed temporary locals.
   a:=e.expression(x.args[0]);left:=e.locals;e.locals++;e.code.inst(0x21,left)
   b:=e.expression(x.args[1]);right:=e.locals;e.locals++;e.code.inst(0x21,right)
   if a=="float"&&b=="float"{e.val(30,left,right,e.constant(0));return "float"};if a!=b||(a!="int"&&a!="string"){webReject("unsupported addition at line %d",e.line)}
   operation:=int64(0);if a=="string"{operation=5}
   e.code.i32(operation);e.code.inst(0x20,left);e.code.inst(0x20,right);e.code.i32(int64(e.line));e.code.inst(0x10,1);return a
  case "-","*","/","%":
   l,a:=e.save(x.args[0]);r,b:=e.save(x.args[1]);if a!=b{webReject("numeric type mismatch")};op:=map[string]int64{"-":1,"*":2,"/":3,"%":4}[x.s];if a=="float"{if op==4{webReject("float remainder unsupported")};e.val(30,l,r,e.constant(op));return "float"};if a!="int"{webReject("integer arithmetic required")};e.code.i32(op);e.code.inst(0x20,l);e.code.inst(0x20,r);e.code.i32(int64(e.line));e.code.inst(0x10,1);return "int"
  case "==","!=","<",">","<=",">=":
   l,a:=e.save(x.args[0]);r,b:=e.save(x.args[1]);if a!=b{webReject("comparison type mismatch")};if a=="float"{e.val(31,l,r,e.constant(map[string]int64{"==":0,"!=":1,"<":2,">":3,"<=":4,">=":5}[x.s]));return "bool"};e.code.inst(0x20,l);e.code.inst(0x20,r);if a=="string"&&x.s!="=="&&x.s!="!="{webReject("string ordering not supported")};e.code.raw(map[string]byte{"==":0x51,"!=":0x52,"<":0x53,">":0x55,"<=":0x57,">=":0x59}[x.s],0xad);return "bool"
  case "&&","||":
   if e.expression(x.args[0])!="bool"{webReject("expected bool")};e.code.raw(0x50,0x04,0x7e)
   if x.s=="&&"{e.code.i64(0);e.code.raw(0x05);if e.expression(x.args[1])!="bool"{webReject("expected bool")}}else{if e.expression(x.args[1])!="bool"{webReject("expected bool")};e.code.raw(0x05);e.code.i64(1)}
   e.code.raw(0x0b);return "bool"
  }
 case "call":
  name:=strings.TrimPrefix(x.s,"@builtin:")
  if name=="len"||name=="@array_len"||name=="sum"{
   if len(x.args)!=1{webReject("wrong argument count")};v,t:=e.save(x.args[0]);zero:=e.constant(0)
   if name=="len"&&t=="string"{e.code.i32(6);e.code.inst(0x20,v);e.code.i64(0);e.code.i32(int64(e.line));e.code.inst(0x10,1);return "int"}
   if webElem(t)==""{webReject("expected collection")};op:=int64(13);if name=="sum"{if webElem(t)!="int"{webReject("sum requires integers")};op=20};e.val(op,v,zero,zero);return "int"
  }
  if name=="input"&&len(x.args)==1{v,t:=e.save(x.args[0]);if t!="string"{webReject("input prompt requires string")};z:=e.constant(0);e.val(22,v,z,z);return "string"}
  if name=="float"&&len(x.args)==1{v,t:=e.save(x.args[0]);if t!="int"{webReject("float requires int")};z:=e.constant(0);e.val(32,v,z,z);return "float"}
  if name=="sqrt"&&len(x.args)==1{v,t:=e.save(x.args[0]);if t!="float"{webReject("sqrt requires float")};z:=e.constant(0);e.val(34,v,z,z);return "float"}
  if name=="abs"||name=="min"||name=="max"{if len(x.args)<1||len(x.args)>2{webReject("invalid intrinsic args")};a,t:=e.save(x.args[0]);if t!="int"{webReject("expected int")};b:=e.constant(0);if len(x.args)==2{var bt string;b,bt=e.save(x.args[1]);if bt!="int"{webReject("expected int")}};e.val(map[string]int64{"abs":35,"min":36,"max":37}[name],a,b,e.constant(0));return "int"}
  if name=="ints"||isList(name){
   if len(x.args)==1&&x.args[0].op=="array"&&name!="ints"{return e.arrayLiteral(x.args[0],true,name)}
   n:=e.constant(4);if name=="ints"{n=e.constant(0)};if len(x.args)==1{var t string;n,t=e.save(x.args[0]);if t!="int"{webReject("capacity must be int")}}else if len(x.args)>1{webReject("invalid capacity")}
   e.val(10,n,e.constant(webKind(name)),e.constant(0));return name
  }
  if receiver,method,ok:=strings.Cut(name,".");ok{if v,exists:=e.get(receiver);exists&&isList(v.typ){
   elem:=listElement(v.typ);zero:=e.constant(0);a,b:=zero,zero
   if len(x.args)>0{a,_=e.save(x.args[0])};if len(x.args)>1{b,_=e.save(x.args[1])}
   op,known:=map[string]int64{"get":11,"set":12,"size":13,"add":14,"remove_at":15,"clear":16,"capacity":17,"contains":18,"remove":19,"sort":23}[method];if !known{webReject("unsupported list method")}
   e.val(op,v.index,a,b);if method=="get"||method=="set"||method=="remove_at"{return elem};if method=="contains"||method=="remove"{return "bool"};return "int"
  }}
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
 case "indexassign":
  v,ok:=e.get(s.name);if !ok||webElem(v.typ)==""{webReject("unsupported indexed assignment")}
  idx,t:=e.save(s.bound);if t!="int"{webReject("expected int index")};zero:=e.constant(0);e.val(11,v.index,idx,zero);old:=e.temp();e.code.inst(0x21,old)
  rhs,rt:=e.save(s.e);if rt!=webElem(v.typ){webReject("indexed assignment type mismatch")};if s.typ!=""{op:=map[string]int64{"+":0,"-":1,"*":2,"/":3,"%":4}[s.typ];e.code.i32(op);e.code.inst(0x20,old);e.code.inst(0x20,rhs);e.code.i32(int64(e.line));e.code.inst(0x10,1);rhs=e.temp();e.code.inst(0x21,rhs)}
  e.val(12,v.index,idx,rhs);e.code.raw(0x1a)
 case "show","showmany":
  args:=[]*expr{s.e};if s.op=="showmany"{args=s.e.args};for i,a:=range args{t:=e.expression(a);kind:=webKind(t);e.code.i32(kind);last:=int64(0);if i==len(args)-1{last=1};e.code.i32(last);e.code.inst(0x10,2)}
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
// a standard WASM module. R3 adds floats, arrays, lists, formatting, input and
// selected pure integer stdlib functions. Unsupported constructs are rejected.
func CompileWeb(source string)(out []byte,err error){
 defer func(){if v:=recover();v!=nil{out=nil;err=fmt.Errorf("%v",v)}}()
 if len(source)>32768{return nil,fmt.Errorf("Web Engine: 32 KiB source limit")}
 imports,_,_:=parseUnit(source);for _,n:=range imports{if n!="@std:math"{webReject("only tripol.math is currently enabled by this target")}}
 fs,ss:=parse(source);validateCode(fs,ss,"linux")
 m:=&webModule{funcs:map[string]webFunction{},strings:[]string{},stringIDs:map[string]int64{}}
 m.signature([]byte{0x7f},nil) // tick(line)
 m.signature([]byte{0x7f,0x7e,0x7e,0x7f},[]byte{0x7e}) // checked arithmetic
 m.signature([]byte{0x7e,0x7f,0x7f},nil) // show
 m.signature([]byte{0x7f},nil) // assert failure
 m.signature([]byte{0x7f,0x7e,0x7e,0x7e,0x7f},[]byte{0x7e}) // value operations
 // Register only reachable functions; imported float functions need not be emitted.
 byName:=map[string]fn{};for _,f:=range fs{byName[f.name]=f}
 seen:=map[string]bool{}
 var walkExpr func(*expr);var walkBlock func([]stmt)
 walkExpr=func(x *expr){if x==nil{return};if x.op=="call"&&!seen[x.s]{seen[x.s]=true;if f,ok:=byName[x.s];ok{m.declare(f);walkBlock(f.body);for _,d:=range f.defaults{walkExpr(d)}}};for _,a:=range x.args{walkExpr(a)}}
 walkBlock=func(ss []stmt){for _,s:=range ss{walkExpr(s.e);walkExpr(s.bound);walkBlock(s.yes);walkBlock(s.no);walkBlock(s.final)}}
 walkBlock(ss)
 for _,f:=range m.order{e:=&webEmitter{module:m,result:f.f.ret};e.push();for _,a:=range f.f.args{e.local(a.name,a.typ)};e.params=e.locals;e.block(f.f.body,false);e.code.raw(0x00);m.bodies=append(m.bodies,e.body())}
 mainType:=m.signature(nil,nil);mainIndex:=uint32(5+len(m.order));e:=&webEmitter{module:m};e.push();e.block(ss,false);m.bodies=append(m.bodies,e.body())
 var module webCode;module.raw(0,97,115,109,1,0,0,0)
 var ts webCode;ts.u(uint32(len(m.types)));for _,s:=range m.types{ts.raw(0x60);ts.u(uint32(len(s.args)));ts.raw(s.args...);ts.u(uint32(len(s.results)));ts.raw(s.results...)};module.section(1,ts)
 var im webCode;im.u(5);for i,n:=range []string{"tick","arithmetic","show","assertFail","value"}{im.name("tripol");im.name(n);im.raw(0);im.u(uint32(i))};module.section(2,im)
 var funcs webCode;funcs.u(uint32(len(m.order)+1));for _,f:=range m.order{funcs.u(f.typ)};funcs.u(mainType);module.section(3,funcs)
 var exports webCode;exports.u(1);exports.name("run");exports.raw(0);exports.u(mainIndex);module.section(7,exports)
 var code webCode;code.u(uint32(len(m.bodies)));for _,b:=range m.bodies{code.u(uint32(len(b)));code.raw(b...)};module.section(10,code)
 var literals webCode;literals.name("tripol.strings");encoded,encodeErr:=json.Marshal(m.strings);if encodeErr!=nil{return nil,encodeErr};literals.raw(encoded...);module.section(0,literals)
 return []byte(module),nil
}
