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
type webModule struct { definitions map[string]fn; funcs map[string]webFunction; order []webFunction; types []webSig; bodies []webCode; current int; strings []string; stringIDs map[string]int64; records map[string]fn; dictionaries map[string][]param; callbackSlots map[uint32]uint32; callbackElements []uint32; indirectTypes map[int]uint32 }
type webEmitter struct { module *webModule; code webCode; scopes []map[string]webLocal; locals,params uint32; line int; labels []string; result string }
func webReject(format string,args ...any){die("Web Engine R3: "+format,args...)}
func webScalar(t string)bool{return t=="int"||t=="bool"||t=="string"||t=="float"||t=="ints"||isList(t)||optionalElement(t)!=""}
func (m *webModule) signature(a,r []byte)uint32{m.types=append(m.types,webSig{a,r});return uint32(len(m.types)-1)}
func (m *webModule) supported(t string)bool{if webScalar(t)||strings.HasPrefix(t,callbackPrefix){return true};_,ok:=m.records[t];return ok}
func (m *webModule) fields(t string)[]param{if f,ok:=m.records[t];ok{return f.args};return m.dictionaries[t]}
func (m *webModule) field(t,n string)(int64,string){for i,f:=range m.fields(t){if f.name==n{return int64(i),f.typ}};webReject("unknown field %s on %s",n,t);return 0,""}
func (m *webModule) declare(f fn)webFunction{
 if v,ok:=m.funcs[f.name];ok{return v}
 if f.record||(f.ret!="void"&&!m.supported(f.ret)){webReject("function %s requires a scalar return type",f.name)}
 args:=make([]byte,len(f.args));for i,a:=range f.args{if !m.supported(a.typ){webReject("parameter %s requires a supported type",a.name)};args[i]=0x7e}
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
func webKind(t string)int64{switch t{case "int":return 0;case "bool":return 1;case "string":return 2;case "float":return 3};if t=="ints"||strings.HasPrefix(t,"array:"){return 4};if typedListElement(t)!=""{return 5};if isList(t){return 5+webKind(listElement(t))};webReject("unsupported value type %s",t);return 0}
func webElem(t string)string{if t=="ints"||strings.HasPrefix(t,"array:"){return "int"};return listElement(t)}
func (e *webEmitter) temp()uint32{v:=e.locals;e.locals++;return v}
func (e *webEmitter) save(x *expr)(uint32,string){t:=e.expression(x);v:=e.temp();e.code.inst(0x21,v);return v,t}
func (e *webEmitter) val(op int64,a,b,c uint32){e.code.i32(op);e.code.inst(0x20,a);e.code.inst(0x20,b);e.code.inst(0x20,c);e.code.i32(int64(e.line));e.code.inst(0x10,4)}
func (e *webEmitter) constant(n int64)uint32{v:=e.temp();e.code.i64(n);e.code.inst(0x21,v);return v}
func (e *webEmitter) arrayLiteral(x *expr,list bool,requested string)string{
 zero:=e.constant(0);size:=e.constant(int64(len(x.args)));kind:=int64(4);typ:="array:"+fmt.Sprint(len(x.args))
 if list{typ=requested;if len(x.args)>0&&requested=="list"{v,t:=e.save(x.args[0]);typ=listType(t);if typ==""{webReject("list requires matching concrete element types")};kind=webKind(typ);k:=e.constant(kind);e.val(10,size,k,zero);handle:=e.temp();e.code.inst(0x21,handle);e.val(14,handle,v,zero);e.code.raw(0x1a);for _,item:=range x.args[1:]{v,t:=e.save(item);if t!=listElement(typ){webReject("mixed list")};e.val(14,handle,v,zero);e.code.raw(0x1a)};e.code.inst(0x20,handle);return typ};kind=webKind(typ)}
 k:=e.constant(kind);e.val(10,size,k,zero);handle:=e.temp();e.code.inst(0x21,handle)
 for i,item:=range x.args{v,t:=e.save(item);elem:=webElem(typ);if t!=elem{webReject("invalid collection element")};if list{e.val(14,handle,v,zero)}else{idx:=e.constant(int64(i));e.val(12,handle,idx,v)};e.code.raw(0x1a)}
 e.code.inst(0x20,handle);return typ
}
func (e *webEmitter) aggregate(fields []param,values []*expr)uint32{
 z:=e.constant(0);e.val(10,e.constant(int64(len(fields))),e.constant(4),z);handle:=e.temp();e.code.inst(0x21,handle)
 for i,f:=range fields{v,t:=e.save(values[i]);if t!=f.typ{webReject("aggregate field type mismatch")};e.val(12,handle,e.constant(int64(i)),v);e.code.raw(0x1a)}
 return handle
}
func (e *webEmitter) expression(x *expr)string{
 if x==nil{webReject("missing expression")};switch x.op{
 case "int","bool","float":e.code.i64(x.n);return x.op
 case "string":e.code.i64(e.module.stringID(x.s));return "string"
 case "capture":
  env,ok:=e.get(closureEnvironment);if !ok{webReject("capture requires a closure")}
  e.val(11,env.index,e.constant(x.n+1),e.constant(0));return x.s
 case "closure":
  return e.webClosure(x)
 case "function-ref":
  f, ok := e.module.definitions[x.s]
  typ := functionCallbackType(f)
  if !ok || typ == "" { webReject("unknown callback target %s", x.s) }
  callee := e.module.declare(f)
  e.code.i64(int64(e.module.callbackSlot(callee.index)))
  return typ
 case "array":return e.arrayLiteral(x,false,"")
 case "none": e.code.i64(0);return x.s
 case "some":
  item,_:=e.save(x.args[0]);z:=e.constant(0)
  e.val(10,e.constant(1),e.constant(4),z);h:=e.temp();e.code.inst(0x21,h)
  e.val(12,h,z,item);e.code.raw(0x1a)
  e.code.inst(0x20,h);e.code.i64(1);e.code.raw(0x7c);return x.s
 case "bound-args":
  e.push();defer e.pop()
  for _, binding := range x.args[:len(x.args)-1] {
   saved,typ:=e.save(binding.args[0]);v:=e.local(binding.s,typ)
   e.code.inst(0x20,saved);e.code.inst(0x21,v.index)
  }
  return e.expression(x.args[len(x.args)-1])
 case "dictionary":
  fields:=[]param{};values:=[]uint32{};for i:=0;i<len(x.args);i+=2{if x.args[i].op!="string"{webReject("literal dictionary key required")};v,t:=e.save(x.args[i+1]);fields=append(fields,param{x.args[i].s,t});values=append(values,v)}
  // The native validator has already checked the shape. The web shape is
  // tracked independently and uses typed slots with a stable reference.
  shape:="@webdict:"+fmt.Sprint(len(e.module.dictionaries));for k,old:=range e.module.dictionaries{if fmt.Sprint(old)==fmt.Sprint(fields){shape=k;break}};e.module.dictionaries[shape]=fields
  z:=e.constant(0);e.val(10,e.constant(int64(len(fields))),e.constant(4),z);h:=e.temp();e.code.inst(0x21,h);for i,v:=range values{e.val(12,h,e.constant(int64(i)),v);e.code.raw(0x1a)};e.code.inst(0x20,h);return shape
 case "field":
  base,t:=e.save(x.args[0]);idx,kind:=e.module.field(t,x.s);z:=e.constant(0);e.val(11,base,e.constant(idx),z);return kind
 case "format":
  zero:=e.constant(0);acc:=e.constant(e.module.stringID(""));for _,part:=range x.args{v,t:=e.save(part);if t!="string"&&t!="int"&&t!="bool"{webReject("unsupported interpolation type")};k:=e.constant(webKind(t));e.val(21,v,k,zero);v=e.temp();e.code.inst(0x21,v);e.code.i32(5);e.code.inst(0x20,acc);e.code.inst(0x20,v);e.code.i32(int64(e.line));e.code.inst(0x10,1);e.code.inst(0x21,acc)};e.code.inst(0x20,acc);return "string"
 case "at":
  base,t:=e.save(x.args[0]);idx,i:=e.save(x.args[1])
  if i!="int"||webElem(t)==""{webReject("computed indexing requires a list and an integer")}
  e.val(11,base,idx,e.constant(0));return webElem(t)
 case "index":v,ok:=e.get(x.s);if !ok{webReject("unsupported index target %s",x.s)}
  if strings.HasPrefix(v.typ,"@webdict:"){if x.args[0].op!="string"{webReject("dictionary key must be literal")};idx,kind:=e.module.field(v.typ,x.args[0].s);e.val(11,v.index,e.constant(idx),e.constant(0));return kind}
  if webElem(v.typ)==""{webReject("unsupported index target %s",x.s)};idx,t:=e.save(x.args[0]);if t!="int"{webReject("index must be int")};e.val(11,v.index,idx,e.constant(0));return webElem(v.typ)
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
  // Target metadata is compiled for the web, never borrowed from Linux.
  if name=="system.architecture"{e.code.i64(e.module.stringID("wasm32"));return "string"}
  if name=="memory.capacity"{e.code.i64(6291456);return "int"}
  if name=="platform_name"||name=="compiler_version"{if len(x.args)!=0{webReject("metadata takes no arguments")};v:="web";if name=="compiler_version"{v=version};e.code.i64(e.module.stringID(v));return "string"}
  if name=="arena_used"||name=="arena_committed"{z:=e.constant(0);e.val(70,z,z,z);return "int"}
  if name=="read_text"||name=="write_text"||name=="append_text"||name=="read_bytes"||name=="write_bytes"{
   spec:=map[string]struct{op int64;args []string;ret string}{
    "read_text":{71,[]string{"string"},"string"},"write_text":{72,[]string{"string","string"},"int"},
    "append_text":{73,[]string{"string","string"},"int"},"read_bytes":{74,[]string{"string","int"},"ints"},
    "write_bytes":{75,[]string{"string","ints"},"int"},
   }[name];if len(x.args)!=len(spec.args){webReject("invalid file arguments")};z:=e.constant(0);v:=[3]uint32{z,z,z};for i,a:=range x.args{tmp,t:=e.save(a);if t!=spec.args[i]&&!(spec.args[i]=="ints"&&webElem(t)=="int"){webReject("invalid argument for %s",name)};v[i]=tmp};e.val(spec.op,v[0],v[1],v[2]);return spec.ret
  }
  if name=="atomic_load"||name=="atomic_exchange"||name=="atomic_add"||name=="atomic_compare_exchange"{
   n:=2;if name=="atomic_exchange"||name=="atomic_add"{n=3};if name=="atomic_compare_exchange"{n=4};if len(x.args)!=n{webReject("invalid atomic arguments")}
   a,t:=e.save(x.args[0]);b,u:=e.save(x.args[1]);if webElem(t)!="int"||u!="int"{webReject("invalid atomic array or index")};z:=e.constant(0)
   if n==2{e.val(76,a,b,z);return "int"}
   c,v:=e.save(x.args[2]);if v!="int"{webReject("atomic value requires int")}
   if n==4{d,w:=e.save(x.args[3]);if w!="int"{webReject("atomic value requires int")};e.val(79,a,b,c);e.code.raw(0x1a);e.val(80,a,b,d);return "bool"}
   op:=int64(77);if name=="atomic_add"{op=78};e.val(op,a,b,c);return "int"
  }
  if name=="sleep_ms"{if len(x.args)!=1{webReject("sleep requires milliseconds")};a,t:=e.save(x.args[0]);if t!="int"{webReject("sleep requires int")};z:=e.constant(0);e.val(81,a,z,z);return "int"}
  // Graphics host operations use the actual TRIPOL framebuffer or vertex arrays.
  if name=="game_key"||name=="native_mouse_axis"||name=="pad_read"{webReject("Live input is not yet supported by the browser graphics preview")}
  if name=="game_open"||name=="game_running"||name=="game_present"||name=="game_close"||name=="native_gl_frame"{
   spec:=map[string]struct{op int64;args []string}{
    "game_open":{82,[]string{"string","int","int"}},"game_running":{83,nil},
    "game_present":{84,[]string{"ints","int","int"}},"game_close":{85,nil},
    "native_gl_frame":{86,[]string{"ints","int"}},
   }[name];if len(x.args)!=len(spec.args){webReject("invalid graphics arguments for %s",name)}
   z:=e.constant(0);v:=[3]uint32{z,z,z};for i,a:=range x.args{tmp,t:=e.save(a);if t!=spec.args[i]&&!(spec.args[i]=="ints"&&webElem(t)=="int"){webReject("invalid graphics argument for %s",name)};v[i]=tmp};e.val(spec.op,v[0],v[1],v[2]);return "bool"
  }
  if name=="http_exchange"||name=="discord_request"||name=="discord_callback"||name=="discord_edit"||name=="discord_rest"||name=="gw_open"||strings.HasPrefix(name,"gw_host_")||name=="env_get"{
   webReject("%s requires a native target",name)
  }
  if name=="thread_sum"{webReject("native worker threads require a browser concurrency adapter")}
  if name=="audio_play_file"||name=="native_beep"{webReject("audio playback requires a browser user-gesture adapter; WAV generation is available")}

  if name == "input" || name == "input_int" {
   if len(x.args) > 1 { webReject("input expects at most one prompt") }
   prompt, silent := e.constant(e.module.stringID("")), e.constant(1)
   if len(x.args) == 1 {
    v, t := e.save(x.args[0])
    if t != "string" { webReject("input prompt requires string") }
    prompt, silent = v, e.constant(0)
   }
   if name == "input_int" { e.val(57, prompt, silent, e.constant(0)); return "int" }
   e.val(22, prompt, silent, e.constant(0)); return "string"
  }
  if name=="utf8_char" {
   if len(x.args)!=1 {webReject("utf8_char expects one int")}
   v,t:=e.save(x.args[0])
   if t!="int" {webReject("utf8_char expects an int")}
   z:=e.constant(0);e.val(88,v,z,z);return "string"
  }
  if name=="json_quote" {
   if len(x.args)!=1 {webReject("json_quote expects one string")}
   v,t:=e.save(x.args[0])
   if t!="string" {webReject("json_quote expects a string")}
   z:=e.constant(0);e.val(87,v,z,z);return "string"
  }
  if name=="text"||name=="parse_int"||name=="byte_at"||name=="slice"||name=="has_prefix"||name=="has_suffix"||name=="find_bytes"{
   spec:=map[string]struct{op int64;args []string;ret string}{
    "text":{50,[]string{"int"},"string"},"parse_int":{51,[]string{"string"},"int"},
    "byte_at":{52,[]string{"string","int"},"int"},"slice":{53,[]string{"string","int","int"},"string"},
    "has_prefix":{54,[]string{"string","string"},"bool"},"has_suffix":{55,[]string{"string","string"},"bool"},
    "find_bytes":{56,[]string{"string","string"},"int"},
   }[name]
   if len(x.args)!=len(spec.args){webReject("wrong argument count: %s",name)}
   z:=e.constant(0);v:=[3]uint32{z,z,z};for i,a:=range x.args{tmp,t:=e.save(a);if t!=spec.args[i]{webReject("invalid argument for %s",name)};v[i]=tmp};e.val(spec.op,v[0],v[1],v[2]);return spec.ret
  }
  if name=="u32_xor"||name=="u32_and"||name=="u32_rotr"||name=="u32_shr"{
   if len(x.args)!=2{webReject("wrong argument count")};a,t:=e.save(x.args[0]);b,u:=e.save(x.args[1]);if t!="int"||u!="int"{webReject("u32 requires integer operands")};e.val(map[string]int64{"u32_xor":60,"u32_and":61,"u32_rotr":62,"u32_shr":63}[name],a,b,e.constant(0));return "int"
  }
  if name=="monotonic_ms"||name=="monotonic_ns"{if len(x.args)!=0{webReject("clock takes no arguments")};z:=e.constant(0);op:=int64(58);if name=="monotonic_ns"{op=59};e.val(op,z,z,z);return "int"}
  if name=="array_equal"||name=="array_fill"||name=="array_fill_range"{
   a,t:=e.save(x.args[0]);if webElem(t)!="int"{webReject("integer array required")};b,bt:=e.save(x.args[1]);z:=e.constant(0)
   if name=="array_equal"{if webElem(bt)!="int"{webReject("integer array required")};e.val(43,a,b,z);return "bool"}
   if bt!="int"{webReject("integer argument required")}
   if name=="array_fill"{e.val(44,a,b,z);return "int"}
   n,nt:=e.save(x.args[2]);v,vt:=e.save(x.args[3]);if nt!="int"||vt!="int"{webReject("integer arguments required")};e.val(45,a,b,n);packed:=e.temp();e.code.inst(0x21,packed);e.val(46,a,packed,v);return "int"
  }
  if name=="array_find"||name=="array_copy"||name=="array_copy_into"{
   a,t:=e.save(x.args[0]);if webElem(t)!="int"{webReject("integer array required")};b:=e.constant(0);z:=e.constant(0)
   switch name{
   case "array_find":var bt string;b,bt=e.save(x.args[1]);if bt!="int"{webReject("integer search value required")};e.val(40,a,b,z);return "int"
   case "array_copy":e.val(41,a,b,z);return "ints"
   case "array_copy_into":var bt string;b,bt=e.save(x.args[1]);if webElem(bt)!="int"{webReject("integer source array required")};e.val(42,a,b,z);return "bool"
   }
  }
  if name=="len"||name=="@array_len"||name=="sum"{
   if len(x.args)!=1{webReject("wrong argument count")};v,t:=e.save(x.args[0]);zero:=e.constant(0)
   if name=="len"&&t=="string"{e.code.i32(6);e.code.inst(0x20,v);e.code.i64(0);e.code.i32(int64(e.line));e.code.inst(0x10,1);return "int"}
   if webElem(t)==""&&!(name=="len"&&strings.HasPrefix(t,"@webdict:")){webReject("expected collection")};op:=int64(13);if name=="sum"{if webElem(t)!="int"{webReject("sum requires integers")};op=20};e.val(op,v,zero,zero);return "int"
  }
  if name=="float_bits"{e.expression(x.args[0]);return "int"}
  if name=="float_from_bits"{e.expression(x.args[0]);return "float"}
  if name=="float"&&len(x.args)==1{v,t:=e.save(x.args[0]);if t!="int"{webReject("float requires int")};z:=e.constant(0);e.val(32,v,z,z);return "float"}
  if name=="sqrt"&&len(x.args)==1{v,t:=e.save(x.args[0]);if t!="float"{webReject("sqrt requires float")};z:=e.constant(0);e.val(34,v,z,z);return "float"}
  if name=="abs"||name=="min"||name=="max"{if len(x.args)<1||len(x.args)>2{webReject("invalid intrinsic args")};a,t:=e.save(x.args[0]);if t!="int"{webReject("expected int")};b:=e.constant(0);if len(x.args)==2{var bt string;b,bt=e.save(x.args[1]);if bt!="int"{webReject("expected int")}};e.val(map[string]int64{"abs":35,"min":36,"max":37}[name],a,b,e.constant(0));return "int"}
  if name=="ints"||isList(name){
   if len(x.args)==1&&x.args[0].op=="array"&&name!="ints"{return e.arrayLiteral(x.args[0],true,name)}
   n:=e.constant(4);if name=="ints"{n=e.constant(0)};if len(x.args)==1{var t string;n,t=e.save(x.args[0]);if t!="int"{webReject("capacity must be int")}}else if len(x.args)>1{webReject("invalid capacity")}
   e.val(10,n,e.constant(webKind(name)),e.constant(0));return name
  }
  if receiver,method,ok:=strings.Cut(name,".");ok{if v,exists:=e.get(receiver);exists&&isList(v.typ){
   return e.listMethod(v.index,v.typ,method,x.args)
  }}
  if f,ok:=e.module.records[x.s];ok{
   if len(x.args)>len(f.args){webReject("wrong record argument count")};values:=make([]*expr,len(f.args));for i:=range values{if i<len(x.args){values[i]=x.args[i]}else if i<len(f.defaults){values[i]=f.defaults[i]};if values[i]==nil{webReject("missing field value")}}
   h:=e.aggregate(f.args,values);e.code.inst(0x20,h);return f.name
  }
  if f, ok := e.module.definitions[x.s]; ok {
   return e.call(f, x.args, 0)
  }
  if v, ok := e.get(x.s); ok {
   if params, result, ok := callbackParts(v.typ); ok {
    return e.callCallback(v.index, params, result, x.args)
   }
  }
  if receiver, member, ok := strings.Cut(x.s, "."); ok {
   if v, exists := e.get(receiver); exists {
    return e.method(v.index, v.typ, member, x.args)
   }
  }
  webReject("unsupported function %s", x.s)
 case "method":
  receiver, typ := e.save(x.args[0])
  return e.method(receiver, typ, x.s, x.args[1:])
 }
 webReject("unsupported expression %s (%s) at line %d",x.op,x.s,e.line);return ""
}
func (e *webEmitter) listMethod(receiver uint32, typ, member string, args []*expr) string {
	ops := map[string]int64{"get":11,"set":12,"size":13,"add":14,"remove_at":15,"clear":16,"capacity":17,"contains":18,"remove":19,"sort":23}
	op, ok := ops[member]
	if !ok { webReject("unsupported list method %s", member) }
	z:=e.constant(0)
	a,b:=z,z
	if len(args)>0 { a,_=e.save(args[0]) }
	if len(args)>1 { b,_=e.save(args[1]) }
	e.val(op,receiver,a,b)
	switch member {
	case "get","set","remove_at":return listElement(typ)
	case "contains","remove":return "bool"
	default:return "int"
	}
}
func (e *webEmitter) method(receiver uint32, typ, member string, args []*expr) string {
	if element:=optionalElement(typ);element!="" {return e.webOptional(receiver,element,member,args)}
	if isList(typ) { return e.listMethod(receiver,typ,member,args) }
	f, ok := e.module.definitions[methodName(typ, member)]
	if !ok {
		webReject("unknown method %s.%s", typ, member)
	}
	e.code.inst(0x20, receiver)
	return e.call(f, args, 1)
}

func (m *webModule) callbackSlot(function uint32) uint32 {
	if slot, ok := m.callbackSlots[function]; ok {
		return slot
	}
	slot := uint32(len(m.callbackElements))
	m.callbackSlots[function] = slot
	m.callbackElements = append(m.callbackElements, function)
	return slot
}

func (m *webModule) indirectType(arity int) uint32 {
	if sig, ok := m.indirectTypes[arity]; ok {
		return sig
	}
	args := make([]byte, arity)
	for i := range args {
		args[i] = 0x7e
	}
	sig := m.signature(args, []byte{0x7e})
	m.indirectTypes[arity] = sig
	return sig
}

func (e *webEmitter) callCallback(receiver uint32, params []string, result string, args []*expr) string {
	if len(args) != len(params) {
		webReject("wrong callback argument count")
	}
	values:=make([]uint32,len(args))
	for i,arg:=range args {
		saved,typ:=e.save(arg)
		if typ!=params[i] {webReject("invalid callback argument %d",i+1)}
		values[i]=saved
	}
	emitCall:=func(slot uint32,arity int) {
		for _,value:=range values {e.code.inst(0x20,value)}
		e.code.inst(0x20,slot)
		e.code.raw(0xa7)
		e.code.inst(0x11,e.module.indirectType(arity))
		e.code.u(0)
	}
	// Nonnegative references are original function-table slots. The
	// negative range identifies boxed closures with a separate environment.
	e.code.inst(0x20,receiver)
	e.code.i64(0)
	e.code.raw(0x53,0x04,0x7e)
	e.code.i64(-1)
	e.code.inst(0x20,receiver)
	e.code.raw(0x7d)
	env:=e.temp()
	e.code.inst(0x21,env)
	z:=e.constant(0)
	e.val(11,env,z,z)
	slot:=e.temp()
	e.code.inst(0x21,slot)
	e.code.inst(0x20,env)
	emitCall(slot,len(params)+1)
	e.code.raw(0x05)
	emitCall(receiver,len(params))
	e.code.raw(0x0b)
	return result
}

func (e *webEmitter) call(f fn, args []*expr, implicit int) string {
	callee := e.module.declare(f)
	if len(args)+implicit > len(f.args) {
		webReject("wrong argument count for %s", f.name)
	}
	for i, a := range f.args[implicit:] {
		var arg *expr
		if i < len(args) {
			arg = args[i]
		} else if i+implicit < len(f.defaults) {
			arg = f.defaults[i+implicit]
		}
		if arg == nil {
			webReject("missing argument %s for %s at line %d", a.name, f.name, e.line)
		}
		actual := e.expression(arg)
		// Native read-only integer-array parameters can borrow fixed arrays.
		borrow := a.typ == "ints" && strings.HasPrefix(actual, "array:") && arg.op == "var" &&
			isStandardFunction(f.name) && (!mutableArrayFunction(f.name) || i > 0)
		if actual != a.typ && !borrow {
			webReject("invalid argument %s for %s: expected %s, got %s at line %d", a.name, f.name, a.typ, actual, e.line)
		}
	}
	e.code.inst(0x10, callee.index)
	return f.ret
}

func (e *webEmitter) begin(op byte,label string){e.code.raw(op,0x40);e.labels=append(e.labels,label)}
func (e *webEmitter) end(){e.code.raw(0x0b);e.labels=e.labels[:len(e.labels)-1]}
func (e *webEmitter) branch(label string){for i:=len(e.labels)-1;i>=0;i--{if e.labels[i]==label{e.code.inst(0x0c,uint32(len(e.labels)-1-i));return}};webReject("invalid %s",label)}
func (e *webEmitter) block(ss []stmt,nested bool){if nested{e.push();defer e.pop()};for _,s:=range ss{e.line=s.line;e.tick();switch s.op{
 case "decl","assign":
  t:=e.expression(s.e);v,ok:=e.get(s.name);if s.op=="decl"||!ok{v=e.local(s.name,t)};if v.typ!=t{webReject("assignment type mismatch")};e.code.inst(0x21,v.index)
 case "indexassign":
  v,ok:=e.get(s.name);if !ok{webReject("unsupported indexed assignment")}
  elem:=webElem(v.typ);var idx uint32
  if strings.HasPrefix(v.typ,"@webdict:"){if s.bound.op!="string"{webReject("dictionary key must be literal")};i,t:=e.module.field(v.typ,s.bound.s);idx=e.constant(i);elem=t}else{var t string;idx,t=e.save(s.bound);if t!="int"||elem==""{webReject("expected collection and integer index")}};zero:=e.constant(0);e.val(11,v.index,idx,zero);old:=e.temp();e.code.inst(0x21,old)
  rhs,rt:=e.save(s.e);if rt!=elem{webReject("indexed assignment type mismatch")};if s.typ!=""{op:=map[string]int64{"+":0,"-":1,"*":2,"/":3,"%":4}[s.typ];e.code.i32(op);e.code.inst(0x20,old);e.code.inst(0x20,rhs);e.code.i32(int64(e.line));e.code.inst(0x10,1);rhs=e.temp();e.code.inst(0x21,rhs)}
  e.val(12,v.index,idx,rhs);e.code.raw(0x1a)
 case "fieldassign":
  if s.bound.op!="field"||s.bound.args[0].op!="var"{webReject("field assignment requires a local record")}
  v,ok:=e.get(s.bound.args[0].s);if !ok{webReject("unknown record binding")};i,kind:=e.module.field(v.typ,s.bound.s);rhs,t:=e.save(s.e);if t!=kind{webReject("field assignment type mismatch")};e.val(12,v.index,e.constant(i),rhs);e.code.raw(0x1a)
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
 defer func(){
  if v:=recover();v!=nil {
   out=nil
   d:=normalizeDiagnostic(v)
   message:=d.Message
   if d.Suggestion!="" { message+=" Did you mean "+d.Suggestion+"?" }
   if d.Line>0 {
    location:=fmt.Sprintf("Line %d",d.Line)
    if d.Column>0{location+=fmt.Sprintf(":%d",d.Column)}
    message=location+": "+message
   }
   err=fmt.Errorf("%s: %s",d.Code,message)
  }
 }()
 if len(source)>32768{return nil,fmt.Errorf("Web Engine: 32 KiB source limit")}
 imports,fs,ss:=parseUnit(source);for _,n:=range imports{switch n{case "@std:discord","@std:math","@std:arrays","@std:collections","@std:datetime","@std:bits","@std:bitset","@std:binary","@std:buffer","@std:checksum","@std:colors","@std:compression","@std:console","@std:convert","@std:dictionary","@std:encoding","@std:geometry","@std:json","@std:matrix","@std:path","@std:physics","@std:random","@std:search","@std:sequences","@std:simd","@std:sort","@std:stats","@std:text","@std:validate","@std:vector","@std:stack","@std:queue","@std:animation","@std:camera","@std:particles","@std:canvas","@std:sprite","@std:testing","@std:errors","@std:network","@std:http","@std:image","@std:time","@std:profiler","@std:heap","@std:utf8","@std:csv","@std:stream","@std:regex","@std:cli","@std:crypto","@std:files","@std:io","@std:atomic","@std:memory","@std:system","@std:audio","@std:process","@std:thread","@std:game","@std:events":default:webReject("library %s requires a browser adapter or native target",n)}}
 fs = withStandardImports(imports, fs);fs=resolveSimpleProgram(fs,ss);validateCode(fs,ss,"linux")
 m:=&webModule{definitions:map[string]fn{},funcs:map[string]webFunction{},strings:[]string{},stringIDs:map[string]int64{},records:map[string]fn{},dictionaries:map[string][]param{},callbackSlots:map[uint32]uint32{},indirectTypes:map[int]uint32{}}
 m.signature([]byte{0x7f},nil) // tick(line)
 m.signature([]byte{0x7f,0x7e,0x7e,0x7f},[]byte{0x7e}) // checked arithmetic
 m.signature([]byte{0x7e,0x7f,0x7f},nil) // show
 m.signature([]byte{0x7f},nil) // assert failure
 m.signature([]byte{0x7f,0x7e,0x7e,0x7e,0x7f},[]byte{0x7e}) // value operations
 // Declare a function when a call is emitted. This resolves record methods
 // against the actual receiver and leaves unused native-only functions out.
 for _, f := range fs {
  if f.record { m.records[f.name] = f } else { m.definitions[f.name] = f }
 }
 main := &webEmitter{module:m}
 main.push()
 main.block(ss,false)
 for i := 0; i < len(m.order); i++ {
  f := m.order[i].f
  e := &webEmitter{module:m,result:f.ret}
  e.push()
  for _, a := range f.args { e.local(a.name,a.typ) }
  e.params=e.locals
  e.block(f.body,false)
  if f.ret=="void" {e.code.i64(0)} else {e.code.raw(0x00)}
  m.bodies=append(m.bodies,e.body())
 }
 mainType:=m.signature(nil,nil)
 mainIndex:=uint32(5+len(m.order))
 m.bodies=append(m.bodies,main.body())
 var module webCode;module.raw(0,97,115,109,1,0,0,0)
 var ts webCode;ts.u(uint32(len(m.types)));for _,s:=range m.types{ts.raw(0x60);ts.u(uint32(len(s.args)));ts.raw(s.args...);ts.u(uint32(len(s.results)));ts.raw(s.results...)};module.section(1,ts)
 var im webCode;im.u(5);for i,n:=range []string{"tick","arithmetic","show","assertFail","value"}{im.name("tripol");im.name(n);im.raw(0);im.u(uint32(i))};module.section(2,im)
 var funcs webCode;funcs.u(uint32(len(m.order)+1));for _,f:=range m.order{funcs.u(f.typ)};funcs.u(mainType);module.section(3,funcs)
 if len(m.callbackElements) != 0 {
  var table webCode
  table.u(1)
  table.raw(0x70, 0) // funcref, minimum-only limits
  table.u(uint32(len(m.callbackElements)))
  module.section(4, table)
 }
 var exports webCode;exports.u(1);exports.name("run");exports.raw(0);exports.u(mainIndex);module.section(7,exports)
 if len(m.callbackElements) != 0 {
  var elements webCode
  elements.u(1)
  elements.u(0) // active segment for table zero
  elements.i32(0)
  elements.raw(0x0b)
  elements.u(uint32(len(m.callbackElements)))
  for _, function := range m.callbackElements { elements.u(function) }
  module.section(9, elements)
 }
 var code webCode;code.u(uint32(len(m.bodies)));for _,b:=range m.bodies{code.u(uint32(len(b)));code.raw(b...)};module.section(10,code)
 var literals webCode;literals.name("tripol.strings");encoded,encodeErr:=json.Marshal(m.strings);if encodeErr!=nil{return nil,encodeErr};literals.raw(encoded...);module.section(0,literals)
 return []byte(module),nil
}
