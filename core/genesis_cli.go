package core

import (
 "bytes"
 "context"
 "encoding/hex"
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strings"
 "time"
)

// Host transport only: GENESIS produces every machine-code and ELF byte.
// No Python dependency; the stage is an explicitly selected native executable.
func genesisCLI(args []string) {
 if len(args)!=4||args[2]!="-o" {die("usage: tripol genesis STAGE FILE -o OUTPUT")}
 stage,e:=filepath.Abs(args[0]);if e!=nil{die("%v",e)}
 src,e:=filepath.Abs(args[1]);if e!=nil{die("%v",e)}
 if strings.ContainsAny(src,"\r\n"){die("invalid source path")}
 st,e:=os.Stat(src);if e!=nil||!st.Mode().IsRegular()||st.Size()>4096{die("GENESIS requires a regular source file of at most 4096 bytes")}
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second);defer cancel()
 cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(src+"\n")
 var out bytes.Buffer;cmd.Stdout=&out;cmd.Stderr=os.Stderr
 if e=cmd.Run();e!=nil{die("GENESIS stage failed: %v",e)}
 if out.Len()>140000{die("GENESIS output limit")}
 image,e:=hex.DecodeString(strings.TrimSpace(out.String()))
 if e!=nil||len(image)<120||!bytes.Equal(image[:4],[]byte{127,69,76,70}){die("invalid GENESIS ELF transport")}
 f,e:=os.OpenFile(args[3],os.O_WRONLY|os.O_CREATE|os.O_EXCL,0755);if e!=nil{die("output: %v",e)}
 _,e=f.Write(image);closeErr:=f.Close();if e!=nil||closeErr!=nil{_ = os.Remove(args[3]);die("failed writing GENESIS output")}
 fmt.Println("GENESIS ELF:",args[3],len(image),"bytes")
 if runtime.GOOS=="windows"{fmt.Println("Linux executable; run on Linux x86-64.")}
}
