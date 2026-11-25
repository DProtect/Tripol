package core

import (
 "bytes"
 "crypto/sha256"
 pecheck "debug/pe"
 "fmt"
 "runtime"
 "strings"
 "testing"
)

func TestQuantumCrypto(t *testing.T){
 for _,s:=range []string{"","abc",strings.Repeat("a",55),strings.Repeat("a",56),strings.Repeat("a",64),strings.Repeat("a",1000)}{
  src:=fmt.Sprintf("import tripol.crypto\nshow crypto.sha256_text(%q)",s)
  want:=fmt.Sprintf("%x\n",sha256.Sum256([]byte(s)))
  runCase(t,src,want,false)
 }
 runCase(t,`import tripol.crypto
 data=[0,255,17,0]
 show crypto.sha256(data)
 try {bad=[256];crypto.sha256(bad)} catch e {show e}`,fmt.Sprintf("%x\n12\n",sha256.Sum256([]byte{0,255,17,0})),false)
}
func TestQuantumAudioCLI(t *testing.T){
 runCase(t,`import tripol.audio
import tripol.binary
import tripol.cli
samples=[0,-32768,32767]
wav=audio.wav_mono(samples,8000)
show len(wav),binary.read_u32_le(wav,4),binary.read_u32_le(wav,24),binary.read_u32_le(wav,40)
show wav[46],wav[47],wav[48],wav[49]
show audio.square_wave(2000,1,8000,100)
line="--verbose --port=8080 --name=x --port=9000"
show cli.has(line,"verbose"),cli.has(line,"verb"),cli.integer(line,"port"),cli.value(line,"name"),cli.value(line,"none","fallback")
`,"50 42 8000 6\n0 128 255 127\n[100, 100, -100, -100, 100, 100, -100, -100]\ntrue false 9000 x fallback\n",false)
}
func TestQuantumProcessThread(t *testing.T){
 runCase(t,`import tripol.process
 show process.id()>0
 show process.platform()==platform_name()`,"true\ntrue\n",false)
 src:=`import tripol.thread
 a=[1,2,-4,100]
 show thread.sum(a)
 b=ints(0)
 show thread.sum(b)
 c=[9223372036854775807,1]
 try {thread.sum(c)} catch e {show e}`
 if runtime.GOOS=="windows"{runCase(t,src,"99\n0\n10\n",false)}
 if runtime.GOOS=="linux"{runCase(t,`import tripol.thread
 show thread.supported()`,"false\n",false)}
 f,e:=pecheck.NewFile(bytes.NewReader(compile(src+"\nimport_invalid=0","windows")));if e!=nil{t.Fatal(e)};defer f.Close()
 syms,e:=f.ImportedSymbols();if e!=nil{t.Fatal(e)}
 for _,name:=range []string{"CreateThread","WaitForSingleObject","CloseHandle"}{if !strings.Contains(strings.Join(syms," "),name){t.Fatal(name,syms)}}
 audio,e:=pecheck.NewFile(bytes.NewReader(compile("import tripol.audio\nshow audio.play_file(\"test.wav\")","windows")));if e!=nil{t.Fatal(e)};defer audio.Close()
 symbols,e:=audio.ImportedSymbols();if e!=nil||!strings.Contains(strings.Join(symbols," "),"PlaySoundA"){t.Fatal(symbols,e)}
}
