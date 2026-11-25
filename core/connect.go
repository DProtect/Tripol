package core

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func (c *compiler) networkBuiltin(e *expr) string {
	if len(e.args) != 3 {
		c.err("tcp_exchange expects IPv4 integer, port and request text")
	}
	a := c.a
	c.require(c.expression(e.args[0]), "int")
	a.emit(0x50)
	c.require(c.expression(e.args[1]), "int")
	a.emit(0x50)
	c.require(c.expression(e.args[2]), "string")
	a.emit(0x48, 0x89, 0xc2, 0x59, 0x58) // rax address, rcx port, rdx request
	a.call("tcp_exchange")
	return "string"
}

// PE import table with a separate descriptor per DLL; old non-network layout is retained.
func networkImports(a *asm, kernel []string, rva, textRVA int) ([]byte, int, int, int) {
	type group struct {
		dll      string
		names    []string
		first    int
		ilt, iat int
	}
	groups := []group{{dll: "KERNEL32.dll", names: kernel}}
	if a.needed("iat:10") {
		groups = append(groups, group{dll: "WS2_32.dll", names: []string{"WSAStartup", "socket", "connect", "setsockopt", "send", "shutdown", "recv", "closesocket", "WSACleanup"}, first: 10})
	}
	if a.needed("iat:19") || a.needed("iat:20") || a.needed("game_state") {
		names := []string{"MessageBoxA", "DialogBoxIndirectParamW", "SetWindowLongPtrA", "GetWindowLongPtrA", "SetWindowTextA", "SetDlgItemTextA", "GetDlgItemTextA", "EndDialog"}
		if a.needed("game_state") {
			for len(names) < 31 {
				names = append(names, "")
			}
			names = append(names, "CreateWindowExA", "IsWindow", "PeekMessageA", "TranslateMessage", "DispatchMessageA", "DestroyWindow", "GetDC", "ReleaseDC", "GetClientRect", "GetAsyncKeyState", "GetForegroundWindow")
		}
		if a.needed("iat:91") || a.needed("iat:92"){for len(names)<72{names=append(names,"")};names=append(names,"GetCursorPos","ScreenToClient")}
		for i := range names {
			if !a.needed(fmt.Sprintf("iat:%d", 19+i)) {
				names[i] = ""
			}
		}
		groups = append(groups, group{dll: "USER32.dll", names: names, first: 19})
	}
	if a.needed("iat:30") || a.needed("gw_state") {
		names := []string{"WinHttpOpen", "WinHttpConnect", "WinHttpOpenRequest", "WinHttpSendRequest", "WinHttpReceiveResponse", "WinHttpReadData", "WinHttpCloseHandle", "WinHttpSetTimeouts", "WinHttpSetOption", "WinHttpQueryHeaders"}
		if a.needed("iat:66") || a.needed("iat:67") || a.needed("iat:68") {
			for len(names) < 36 {
				names = append(names, "")
			}
			names = append(names, "WinHttpWebSocketCompleteUpgrade", "WinHttpWebSocketSend", "WinHttpWebSocketReceive")
		}
		groups = append(groups, group{dll: "WINHTTP.dll", names: names, first: 30})
	}

	if a.needed("iat:61") || a.needed("gl_state") {
        names:=make([]string,19)
        if a.needed("iat:61"){names[0]="StretchDIBits"}
        if a.needed("gl_state"){copy(names[15:],[]string{"GetPixelFormat","ChoosePixelFormat","SetPixelFormat","SwapBuffers"})}
        groups=append(groups,group{dll:"GDI32.dll",names:names,first:61})
    }
    if a.needed("iat:75"){groups=append(groups,group{dll:"XINPUT1_4.dll",names:[]string{"XInputGetState"},first:75})}
    if a.needed("gl_state"){groups=append(groups,group{dll:"OPENGL32.dll",names:[]string{"wglCreateContext","wglMakeCurrent","wglDeleteContext","glViewport","glClearColor","glClear","glBegin","glColor3ub","glVertex2f","glEnd"},first:80})}
	if a.needed("iat:74") {groups=append(groups,group{dll:"WINMM.dll",names:[]string{"PlaySoundA"},first:74})}
	size := 1024
	if len(groups) > 2 {
		size = 1536
	}
	if a.needed("iat:20") || a.needed("iat:30") || a.needed("gw_state") {
		size = 3072
	}
	if a.needed("game_state") {
		size = 4096
	}
	id := make([]byte, size)
	cursor := align((len(groups)+1)*20, 16)
	for i := range groups {
		g := &groups[i]
		g.ilt = cursor
		for _, n := range g.names {
			if n != "" {
				cursor += 8
			}
		}
		cursor += 8
	}
	firstIAT := cursor
	for i := range groups {
		g := &groups[i]
		g.iat = cursor
		for _, n := range g.names {
			if n != "" {
				cursor += 8
			}
		}
		cursor += 8
	}
	iatSize := cursor - firstIAT
	for i, g := range groups {
		dll := cursor
		copy(id[cursor:], g.dll+"\x00")
		cursor = align(cursor+len(g.dll)+1, 2)
		put32(id, i*20, rva+g.ilt)
		put32(id, i*20+12, rva+dll)
		put32(id, i*20+16, rva+g.iat)
		j := 0
		for k, n := range g.names {
			if n == "" {
				continue
			}
			put64(id, g.ilt+j*8, uint64(rva+cursor))
			put64(id, g.iat+j*8, uint64(rva+cursor))
			copy(id[cursor+2:], n+"\x00")
			cursor = align(cursor+2+len(n)+1, 2)
			a.labels[fmt.Sprintf("iat:%d", g.first+k)] = rva - textRVA + g.iat + j*8
			j++
		}
	}
	reserve := 40
	if a.needed("game_state") {
		reserve += 16
		a.labels["game_state"] = rva - textRVA + len(id) - reserve
		a.labels["game_buffer"] = a.labels["game_state"] + 8
	}
	if a.needed("gl_state"){
        reserve+=16
        a.labels["gl_state"]=rva-textRVA+len(id)-reserve
        a.labels["gl_context"]=a.labels["gl_state"]+8
    }
    if cursor > len(id)-reserve {
		die("network import table overflow")
	}
	if a.needed("arena_state") || a.needed("guard_state") {
		a.labels["arena_state"] = rva - textRVA + len(id) - 40
		a.labels["guard_state"] = rva - textRVA + len(id) - 8
	}
	return id, firstIAT, iatSize, (len(groups) + 1) * 20
}

// One-shot localhost server for a safe, dependency-free installation smoke test.
func netDemo(args []string) {
	port := 8765
	if len(args) > 1 {
		die("usage: tripol net-demo [PORT]")
	}
	if len(args) == 1 {
		var e error
		port, e = strconv.Atoi(args[0])
		if e != nil || port < 1024 || port > 65535 {
			die("demo port must be 1024..65535")
		}
	}
	listener, e := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
	if e != nil {
		die("listen: %v", e)
	}
	defer listener.Close()
	fmt.Printf("TRIPOL demo listening on 127.0.0.1:%d (one request)\n", port)
	conn, e := listener.Accept()
	if e != nil {
		die("accept: %v", e)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	b, e := io.ReadAll(io.LimitReader(conn, 65536))
	if e != nil || len(b) > 65535 {
		die("demo request failed or too large")
	}
	body := "Hello from TRIPOL!"
	response := fmt.Sprintf("HTTP/1.0 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", len(body), body)
	if _, e = io.WriteString(conn, response); e != nil {
		die("demo response: %v", e)
	}
	fmt.Println("Request handled. Demo stopped.")
}
