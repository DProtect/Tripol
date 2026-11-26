package core

import (
 "encoding/json"
 "strings"
)

type dictionaryField struct { Key string; Type string }
func isDictionary(t string)bool{return strings.HasPrefix(t,"@dict:")}
func dictionaryFields(t string)[]dictionaryField{
 var fields []dictionaryField
 if !isDictionary(t)||json.Unmarshal([]byte(strings.TrimPrefix(t,"@dict:")),&fields)!=nil{die("internal dictionary schema")}
 return fields
}
// A statically shaped, mixed-scalar dictionary. Key spelling and type form
// its schema. Unlike tripol.dictionary, its keys cannot be added at runtime.
func(c *compiler) dictionaryLiteral(e *expr)string{
 n:=len(e.args)/2
 if n==0||n>64{c.err("dictionary requires 1 to 64 entries")}
 fields:=make([]dictionaryField,0,n);seen:=map[string]bool{}
 a:=c.a
 for i:=0;i<len(e.args);i+=2{
  key:=e.args[i].s
  if key==""||len(key)>128||seen[key]{c.err("invalid or duplicate dictionary key "+key)}
  seen[key]=true
  typ:=c.expression(e.args[i+1])
  if listType(typ)==""{c.err("dictionary values must be int, string, bool or float")}
  fields=append(fields,dictionaryField{key,typ});a.emit(0x50)
 }
 a.imm(int64(n*8));a.call("alloc")
 for i:=n-1;i>=0;i--{a.emit(0x59,0x48,0x89,0x88);a.u32(i*8)}
 schema,_:=json.Marshal(fields);return "@dict:"+string(schema)
}
func(c *compiler) dictionaryField(v variable,key *expr)(dictionaryField,int){
 if key.op!="string"{c.err("dictionary literal keys must be known strings; use tripol.dictionary for dynamic keys")}
 for i,field:=range dictionaryFields(v.typ){if field.Key==key.s{return field,i*8}}
 c.err("unknown dictionary key "+key.s);return dictionaryField{},0
}
func(c *compiler) dictionaryGet(v variable,key *expr)string{
 field,off:=c.dictionaryField(v,key)
 c.a.load(v.off);c.a.emit(0x48,0x8b,0x80);c.a.u32(off)
 return field.Type
}
func(c *compiler) dictionaryAssign(v variable,s stmt){
 if !v.mut{c.err("cannot modify constant "+s.name)}
 field,off:=c.dictionaryField(v,s.bound);a:=c.a
 if s.typ!=""{
  if field.Type!="int"{c.err("compound dictionary assignment requires integer")}
  c.dictionaryGet(v,s.bound);a.emit(0x50)
 }
 c.require(c.expression(s.e),field.Type)
 if s.typ!=""{a.emit(0x48,0x89,0xc1,0x58);c.compoundInteger(s.typ)}
 a.loadRCX(v.off);a.emit(0x48,0x89,0x81);a.u32(off)
}
func(c *compiler) showDictionary(t string,newline bool){
 a:=c.a;a.emit(0x50)
 a.str("{");a.call("raw_string")
 for i,field:=range dictionaryFields(t){
  prefix:="";if i>0{prefix=", "}
  a.str(prefix+field.Key+": ");a.call("raw_string")
  a.emit(0x48,0x8b,0x04,0x24,0x48,0x8b,0x80);a.u32(i*8)
  a.call("raw_"+field.Type)
 }
 a.emit(0x58);end:="}";if newline{end+="\n"};a.str(end);a.call("raw_string")
}
