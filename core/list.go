package core

import "strings"

func (c *compiler) listCreate(e *expr) string {
 a:=c.a
 if len(e.args)>1{c.err("list expects optional capacity or an integer literal array")}
 if len(e.args)==1&&e.args[0].op=="array" {
  items:=e.args[0].args
  if len(items)>400{c.err("list initializer exceeds 400 elements")}
  a.imm(int64(len(items)));a.call("list_new")
  for _,item:=range items{
   a.emit(0x50);c.require(c.expression(item),"int")
   a.emit(0x48,0x89,0xc1,0x58,0x50);a.call("list_add");a.emit(0x58)
  }
  return "list"
 }
 if len(e.args)==1 {c.require(c.expression(e.args[0]),"int")} else {a.imm(4)}
 a.call("list_new")
 return "list"
}
func (c *compiler) listMethod(e *expr)(string,bool){
 name,method,ok:=strings.Cut(e.s,".")
 if !ok{return "",false}
 v,ok:=c.lookup(name)
 if !ok||v.typ!="list"{return "",false}
 a:=c.a
 count,known:=map[string]int{"add":1,"remove_at":1,"get":1,"set":2,"size":0,"capacity":0,"clear":0}[method]
 if !known{c.err("unknown list method "+method)}
 if len(e.args)!=count{c.err("wrong list argument count")}
 if (method=="add"||method=="remove_at"||method=="set"||method=="clear")&&!v.mut{c.err("cannot modify constant "+name)}
 a.load(v.off)
 if count==0 {
  if method=="capacity"{a.emit(0x48,0x8b,0x00);return "int",true}
  a.emit(0x48,0x8b,0x40,0x08)
  if method=="clear"{a.emit(0x48,0xc7,0x00,0,0,0,0);a.imm(0)}else{a.emit(0x48,0x8b,0x00)}
  return "int",true
 }
 a.emit(0x50)
 c.require(c.expression(e.args[0]),"int")
 if count==2{
  a.emit(0x50);c.require(c.expression(e.args[1]),"int")
  // RAX=value, RCX=index, RDX=stable handle; reload backing after argument calls.
  a.emit(0x59,0x5a,0x48,0x8b,0x52,0x08,0x48,0x3b,0x0a);a.j(0x83,"bounds")
  a.emit(0x48,0x89,0x44,0xca,0x08);return "int",true
 }
 a.emit(0x48,0x89,0xc1,0x58)
 switch method{
 case "add":a.call("list_add")
 case "remove_at":a.call("list_remove")
 case "get":a.emit(0x48,0x8b,0x50,0x08,0x48,0x3b,0x0a);a.j(0x83,"bounds");a.emit(0x48,0x8b,0x44,0xca,0x08)
 }
 return "int",true
}
