//go:build js && wasm

package main

import (
 "syscall/js"
 "tripol/core"
)

func main(){
 compile:=js.FuncOf(func(_ js.Value,args []js.Value) any {
  if len(args)!=1 || args[0].Type()!=js.TypeString{return map[string]any{"ok":false,"error":"Expected TRIPOL source string"}}
  bytes,err:=core.CompileWeb(args[0].String())
  if err!=nil{return map[string]any{"ok":false,"error":err.Error()}}
  arr:=js.Global().Get("Uint8Array").New(len(bytes));js.CopyBytesToJS(arr,bytes)
  return map[string]any{"ok":true,"wasm":arr}
 })
 js.Global().Set("tripolCompileWeb",compile)
 select{}
}
