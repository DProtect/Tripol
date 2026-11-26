package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestLinkLiteralLayout(t *testing.T) {
	a := &asm{b:[]byte{0xaa}, labels:map[string]int{},data:map[string]string{"hello":"one","":"two","x\x00y":"three"},dataOrder:[]string{"hello","","x\x00y"}}
	appendLinkLiterals(a)
	if !bytes.Equal(a.b,[]byte("\xaahello\x00\x00x\x00y\x00")) { t.Fatal(a.b) }
	if a.labels["one"]!=1||a.labels["two"]!=7||a.labels["three"]!=8 { t.Fatal(a.labels) }
}

func TestLinkRelocationBounds(t *testing.T) {
	a:=&asm{b:make([]byte,8),labels:map[string]int{"back":0,"forward":12},fixes:[]fix{{0,"forward"},{4,"back"}}}
	a.finish()
	if binary.LittleEndian.Uint32(a.b)!=8||int32(binary.LittleEndian.Uint32(a.b[4:]))!=-8 { t.Fatal(a.b) }
	for _, f:=range []fix{{-1,"back"},{5,"back"},{0,"missing"}} {
		func(){
			defer func(){if recover()==nil{t.Error("invalid fix accepted",f)}}()
			b:=&asm{b:make([]byte,8),labels:a.labels,fixes:[]fix{f}}
			b.finish()
		}()
	}
}

func BenchmarkLinkLiterals(b *testing.B) {
	order:=make([]string,4096)
	data:=make(map[string]string,len(order))
	for i:=range order{order[i]=fmt.Sprintf("literal-%06d-abcdefghijklmnopqrstuvwxyz",i);data[order[i]]=fmt.Sprint(i)}
	for _, optimized:=range []bool{false,true}{
		b.Run(fmt.Sprint(optimized),func(b *testing.B){
			b.ReportAllocs()
			for i:=0;i<b.N;i++{
				a:=&asm{b:make([]byte,4096),labels:make(map[string]int,len(order)),data:data,dataOrder:order}
				if optimized {appendLinkLiterals(a)} else {for _,s:=range order{a.label(a.data[s]);a.b=append(a.b,s...);a.b=append(a.b,0)}}
			}
		})
	}
}
