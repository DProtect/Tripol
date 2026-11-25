.intel_syntax noprefix
.section .text.print_float,"ax"
# rax = IEEE754 bits; format rounded fixed six decimals, trap outside int64 scaled range.
 movq xmm0,rax
 mov rdx,0x412e848000000000
 movq xmm1,rdx
 mulsd xmm0,xmm1
 cvtsd2si rax,xmm0
 mov rdx,0x8000000000000000
 cmp rax,rdx
 je float_error
 push rbp
 mov rbp,rsp
 sub rsp,48
 lea rsi,[rbp-1]
 mov byte ptr [rsi],10
 mov r8d,1
 xor r9d,r9d
 test rax,rax
 jns 1f
 mov r9b,1
 neg rax
1: mov ecx,10
 mov r10d,6
2: xor edx,edx
 div rcx
 add dl,48
 dec rsi
 mov [rsi],dl
 inc r8
 dec r10d
 jnz 2b
 dec rsi
 mov byte ptr [rsi],46
 inc r8
3: xor edx,edx
 div rcx
 add dl,48
 dec rsi
 mov [rsi],dl
 inc r8
 test rax,rax
 jnz 3b
 test r9b,r9b
 jz 4f
 dec rsi
 mov byte ptr [rsi],45
 inc r8
4: mov rdx,r8
 call write
 leave
 ret
.section .text.strlen,"ax"
 # SSE2 aligned scan: scalar prefix avoids crossing an unmapped page.
 mov rcx,rax
 xor eax,eax
 push r10
 push r11
 mov r11,rcx
1: test r11b,15
 jz 3f
 cmp byte ptr [r11],0
 je 5f
 inc r11
 inc rax
 jmp 1b
3: pxor xmm1,xmm1
4: movdqa xmm0,[r11]
 pcmpeqb xmm0,xmm1
 pmovmskb r10d,xmm0
 test r10d,r10d
 jnz 6f
 add r11,16
 add rax,16
 jmp 4b
6: bsf r10d,r10d
 add rax,r10
5: pop r11
 pop r10
 ret
.section .text.concat,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,40
 mov [rbp-8],rax
 mov [rbp-16],rcx
 call strlen
 mov [rbp-24],rax
 mov rax,[rbp-16]
 call strlen
 mov [rbp-32],rax
 add rax,[rbp-24]
 cmp rax,65535
 ja text_error
 inc rax
 call alloc
 mov [rbp-40],rax
 mov rdi,rax
 mov rsi,[rbp-8]
 mov rcx,[rbp-24]
 cld
 rep movsb
 mov rsi,[rbp-16]
 mov rcx,[rbp-32]
 rep movsb
 mov byte ptr [rdi],0
 mov rax,[rbp-40]
 leave
 ret
.section .text.input,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,16
 lea rdx,[rip+arena_state]
 mov rax,[rdx+16]
 test rax,rax
 jnz 5f
 mov eax,4098
 call alloc
 lea rdx,[rip+arena_state]
 mov [rdx+16],rax
5: mov [rbp-8],rax
 mov qword ptr [rbp-16],0
1: mov rax,[rbp-8]
 add rax,[rbp-16]
 call readbyte
 test rax,rax
 js text_error
 jz 3f
 mov rcx,[rbp-16]
 mov rdx,[rbp-8]
 mov al,[rdx+rcx]
 cmp al,10
 je 3f
 test al,al
 jz text_error
 cmp rcx,4096
 jb 6f
 jne text_error
 cmp al,13
 jne text_error
6: inc qword ptr [rbp-16]
 jmp 1b
3: mov rcx,[rbp-16]
 mov rax,[rbp-8]
 test rcx,rcx
 jz 4f
 cmp byte ptr [rax+rcx-1],13
 jne 4f
 dec rcx
4: mov byte ptr [rax+rcx],0
 inc rcx
 mov [rbp-16],rcx
 mov rax,rcx
 call alloc
 mov rdi,rax
 mov rsi,[rbp-8]
 mov rcx,[rbp-16]
 cld
 rep movsb
 leave
 ret
.section .text.parseint,"ax"
 mov rsi,rax
 xor eax,eax
 xor r8d,r8d
 cmp byte ptr [rsi],45
 jne 1f
 mov r8b,1
 inc rsi
 jmp 2f
1: cmp byte ptr [rsi],43
 jne 2f
 inc rsi
2: cmp byte ptr [rsi],0
 je text_error
3: movzx ecx,byte ptr [rsi]
 test cl,cl
 jz 4f
 sub ecx,48
 cmp ecx,9
 ja text_error
 imul rax,rax,10
 jo text_error
 sub rax,rcx
 jo text_error
 inc rsi
 jmp 3b
4: test r8b,r8b
 jnz 5f
 neg rax
 jo text_error
5: ret
 .section .text.alloc,"ax"
 # 256 MiB contiguous virtual reservation. Commit pages only as needed.
 # [cursor, remaining reservation, input scratch, committed end].
 push rbp
 mov rbp,rsp
 sub rsp,24
 mov [rbp-8],rax
 lea rdx,[rip+arena_state]
 cmp qword ptr [rdx],0
 jne 1f
 mov eax,268435456
 call osalloc
 lea rdx,[rip+arena_state]
 mov [rdx],rax
 mov qword ptr [rdx+8],268435456
 mov [rdx+24],rax
1: mov rcx,[rbp-8]
 cmp rcx,[rdx+8]
 ja text_error
 mov rax,[rdx]
 mov [rbp-16],rax
 add rax,rcx
 cmp rax,[rdx+24]
 jbe 2f
 add rax,4095
 and rax,-4096
 mov [rbp-24],rax
 mov rcx,rax
 mov rax,[rdx+24]
 sub rcx,rax
 call oscommit
 lea rdx,[rip+arena_state]
 mov rax,[rbp-24]
 mov [rdx+24],rax
2: mov rcx,[rbp-8]
 add [rdx],rcx
 sub [rdx+8],rcx
 mov rax,[rbp-16]
 leave
 ret
.section .text.osalloc_linux,"ax"
 mov rsi,rax
 xor edi,edi
 xor edx,edx
 mov r10d,0x22
 mov r8,-1
 xor r9d,r9d
 mov eax,9
 syscall
 cmp rax,-4095
 jae text_error
 ret
.section .text.osalloc_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rdx,rax
 xor ecx,ecx
 mov r8d,0x2000
 mov r9d,4
 call qword ptr [rip+iat4]
 test rax,rax
 jz text_error
 leave
 ret
.section .text.oscommit_linux,"ax"
 mov rdi,rax
 mov rsi,rcx
 mov edx,3
 mov eax,10
 syscall
 test rax,rax
 jnz text_error
 ret
.section .text.oscommit_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rdx,rcx
 mov rcx,rax
 mov r8d,0x1000
 mov r9d,4
 call qword ptr [rip+iat4]
 test rax,rax
 jz text_error
 leave
 ret
.section .text.arena_used,"ax"
 lea rdx,[rip+arena_state]
 xor eax,eax
 cmp qword ptr [rdx],0
 je 1f
 mov eax,268435456
 sub rax,[rdx+8]
1: ret
.section .text.arena_committed,"ax"
 lea rdx,[rip+arena_state]
 xor eax,eax
 cmp qword ptr [rdx],0
 je 1f
 mov rax,[rdx+24]
 sub rax,[rdx]
 sub rax,[rdx+8]
 add rax,268435456
1: ret
.section .text.array_fill,"ax"
 # rax dynamic array pointer, rcx fill value
 mov rdi,rax
 mov rdx,[rdi]
 add rdi,8
 mov rax,rcx
 mov rcx,rdx
 cld
 rep stosq
 mov rax,rdx
 ret
.section .text.readbyte_linux,"ax"
 mov rsi,rax
 xor edi,edi
 mov edx,1
1: xor eax,eax
 syscall
 cmp rax,-4
 je 1b
 ret
.section .text.readbyte_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,64
 mov [rsp+48],rax
 mov ecx,-10
 call qword ptr [rip+iat0]
 mov rcx,rax
 mov rdx,[rsp+48]
 mov r8d,1
 lea r9,[rsp+56]
 mov qword ptr [rsp+32],0
 call qword ptr [rip+iat3]
 test eax,eax
 jz 1f
 mov eax,[rsp+56]
 leave
 ret
1: call qword ptr [rip+iat5]
 cmp eax,109
 je 2f
 mov rax,-1
 leave
 ret
2: xor eax,eax
 leave
 ret
.section .text.text_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,13
 jmp exit
1: .ascii "TRIPOL runtime error: invalid input, text too long or allocation failed\n"
2:
.section .text.float_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,16
 jmp exit
1: .ascii "TRIPOL runtime error: float outside printable range\n"
2:
.section .text.file_read,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,32
 xor edx,edx
 call file_open
 mov [rbp-8],rax
 mov eax,65536
 call alloc
 mov [rbp-16],rax
 mov qword ptr [rbp-24],0
1: mov rax,[rbp-8]
 mov rsi,[rbp-16]
 add rsi,[rbp-24]
 mov edx,65536
 sub rdx,[rbp-24]
 call file_chunk_read
 test rax,rax
 js file_error
 jz 2f
 add [rbp-24],rax
 cmp qword ptr [rbp-24],65536
 jae file_error
 jmp 1b
2: mov rax,[rbp-8]
 call file_close
 mov rax,[rbp-16]
 mov rcx,[rbp-24]
 mov byte ptr [rax+rcx],0
 # Reject embedded NUL; string API is NUL-terminated.
 xor edx,edx
3: cmp rdx,rcx
 je 4f
 cmp byte ptr [rax+rdx],0
 je file_error
 inc rdx
 jmp 3b
4: # Return unused arena space; no allocation happened after the buffer.
 mov edx,65535
 sub rdx,rcx
 lea r8,[rip+arena_state]
 sub [r8],rdx
 add [r8+8],rdx
 leave
 ret
.section .text.file_write,"ax"
 mov edx,1
 jmp file_write_common
.section .text.file_append,"ax"
 mov edx,2
 jmp file_write_common
.section .text.file_write_common,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,40
 mov [rbp-8],rcx
 call file_open
 mov [rbp-16],rax
 mov rax,[rbp-8]
 call strlen
 mov [rbp-24],rax
 mov qword ptr [rbp-32],0
1: mov rdx,[rbp-24]
 sub rdx,[rbp-32]
 jz 2f
 mov rax,[rbp-16]
 mov rsi,[rbp-8]
 add rsi,[rbp-32]
 call file_chunk_write
 test rax,rax
 jle file_error
 add [rbp-32],rax
 jmp 1b
2: mov rax,[rbp-16]
 call file_close
 mov rax,[rbp-24]
 leave
 ret
.section .text.file_open_linux,"ax"
 mov rdi,rax
 xor esi,esi
 cmp edx,0
 je 1f
 mov esi,577 # O_WRONLY|O_CREAT|O_TRUNC
 cmp edx,1
 je 1f
 mov esi,1089 # O_WRONLY|O_CREAT|O_APPEND
1: mov edx,438 # 0666, filtered by umask
 mov eax,2
 syscall
 test rax,rax
 js file_error
 ret
.section .text.file_chunk_read_linux,"ax"
 mov rdi,rax
1: xor eax,eax
 syscall
 cmp rax,-4
 je 1b
 ret
.section .text.file_chunk_write_linux,"ax"
 mov rdi,rax
1: mov eax,1
 syscall
 cmp rax,-4
 je 1b
 ret
.section .text.file_close_linux,"ax"
 mov rdi,rax
 mov eax,3
 syscall
 test rax,rax
 js file_error
 ret
.section .text.file_open_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,64
 mov rcx,rax
 mov r8d,1 # FILE_SHARE_READ
 xor r9d,r9d
 mov qword ptr [rsp+40],128 # FILE_ATTRIBUTE_NORMAL
 mov qword ptr [rsp+48],0
 cmp edx,0
 je 1f
 cmp edx,1
 je 2f
 mov edx,4 # FILE_APPEND_DATA
 mov qword ptr [rsp+32],4 # OPEN_ALWAYS
 jmp 3f
1: mov edx,0x80000000 # GENERIC_READ
 mov qword ptr [rsp+32],3 # OPEN_EXISTING
 jmp 3f
2: mov edx,0x40000000 # GENERIC_WRITE
 mov qword ptr [rsp+32],2 # CREATE_ALWAYS
3: call qword ptr [rip+iat6]
 cmp rax,-1
 je file_error
 leave
 ret
.section .text.file_chunk_read_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,64
 mov rcx,rax
 mov r8d,edx
 mov rdx,rsi
 lea r9,[rsp+48]
 mov qword ptr [rsp+32],0
 call qword ptr [rip+iat3]
 test eax,eax
 jz file_error
 mov eax,[rsp+48]
 leave
 ret
.section .text.file_chunk_write_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,64
 mov rcx,rax
 mov r8d,edx
 mov rdx,rsi
 lea r9,[rsp+48]
 mov qword ptr [rsp+32],0
 call qword ptr [rip+iat1]
 test eax,eax
 jz file_error
 mov eax,[rsp+48]
 leave
 ret
.section .text.file_close_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rcx,rax
 call qword ptr [rip+iat7]
 test eax,eax
 jz file_error
 leave
 ret
.section .text.file_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,14
 jmp exit
1: .ascii "TRIPOL runtime error: file operation failed or text exceeds limit\n"
2:
.section .text.int_text,"ax"
 # signed int64 -> immutable decimal text; unsigned division handles INT64_MIN.
 push rbp
 mov rbp,rsp
 sub rsp,48
 lea rsi,[rbp-1]
 mov byte ptr [rsi],0
 mov ecx,10
 xor r8d,r8d
 test rax,rax
 jns 1f
 mov r8b,1
 neg rax
1: xor edx,edx
 div rcx
 add dl,48
 dec rsi
 mov [rsi],dl
 test rax,rax
 jnz 1b
 test r8b,r8b
 jz 2f
 dec rsi
 mov byte ptr [rsi],45
2: lea rax,[rbp]
 sub rax,rsi
 mov [rbp-40],rsi
 mov [rbp-48],rax
 call alloc
 mov rdi,rax
 mov rsi,[rbp-40]
 mov rcx,[rbp-48]
 cld
 rep movsb
 leave
 ret
.section .text.text_slice,"ax"
 # rax=string, rcx=start, rdx=end (exclusive), byte offsets
 push rbp
 mov rbp,rsp
 sub rsp,32
 mov [rbp-8],rax
 mov [rbp-16],rcx
 mov [rbp-24],rdx
 test rcx,rcx
 js text_error
 cmp rdx,rcx
 jl text_error
 call strlen
 cmp [rbp-24],rax
 ja text_error
 mov rax,[rbp-24]
 sub rax,[rbp-16]
 inc rax
 mov [rbp-32],rax
 call alloc
 mov rdi,rax
 mov rsi,[rbp-8]
 add rsi,[rbp-16]
 mov rcx,[rbp-32]
 dec rcx
 cld
 rep movsb
 mov byte ptr [rdi],0
 leave
 ret
.section .text.ints_alloc,"ax"
 # Runtime-sized int64 array: [length][elements...], bounded arena allocation.
 cmp rax,65535
 ja text_error
 push rbp
 mov rbp,rsp
 sub rsp,8
 mov [rbp-8],rax
 lea rax,[rax*8+8]
 call alloc
 mov rdx,rax
 mov rcx,[rbp-8]
 mov [rax],rcx
 lea rdi,[rax+8]
 xor eax,eax
 cld
 rep stosq
 mov rax,rdx
 leave
 ret

.section .text.raw_float,"ax"
# rax = IEEE754 bits; format rounded fixed six decimals, trap outside int64 scaled range.
 movq xmm0,rax
 mov rdx,0x412e848000000000
 movq xmm1,rdx
 mulsd xmm0,xmm1
 cvtsd2si rax,xmm0
 mov rdx,0x8000000000000000
 cmp rax,rdx
 je float_error
 push rbp
 mov rbp,rsp
 sub rsp,48
 lea rsi,[rbp-1]
 xor r8d,r8d
 xor r9d,r9d
 test rax,rax
 jns 1f
 mov r9b,1
 neg rax
1: mov ecx,10
 mov r10d,6
2: xor edx,edx
 div rcx
 add dl,48
 dec rsi
 mov [rsi],dl
 inc r8
 dec r10d
 jnz 2b
 dec rsi
 mov byte ptr [rsi],46
 inc r8
3: xor edx,edx
 div rcx
 add dl,48
 dec rsi
 mov [rsi],dl
 inc r8
 test rax,rax
 jnz 3b
 test r9b,r9b
 jz 4f
 dec rsi
 mov byte ptr [rsi],45
 inc r8
4: mov rdx,r8
 call write
 leave
 ret

.section .text.raw_string,"ax"
 mov rsi,rax
 xor edx,edx
1: cmp byte ptr [rsi+rdx],0
 je 2f
 inc rdx
 jmp 1b
2: jmp write
.section .text.raw_bool,"ax"
 test rax,rax
 jz 1f
 lea rax,[rip+2f]
 jmp raw_string
1: lea rax,[rip+3f]
 jmp raw_string
2: .asciz "true"
3: .asciz "false"
.section .text.raw_int,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,48
 lea rsi,[rbp-1]
 xor r8d,r8d
 xor r9d,r9d
 test rax,rax
 jns 1f
 mov r9b,1
 neg rax
1: mov ecx,10
2: xor edx,edx
 div rcx
 add dl,48
 dec rsi
 mov [rsi],dl
 inc r8
 test rax,rax
 jnz 2b
 test r9b,r9b
 jz 3f
 dec rsi
 mov byte ptr [rsi],45
 inc r8
3: mov rdx,r8
 call write
 leave
 ret

# Monotonic time; x86-64 Linux clock_gettime(CLOCK_MONOTONIC).
.section .text.clock_ms_linux,"ax"
 sub rsp,16
 mov eax,228
 mov edi,1
 mov rsi,rsp
 syscall
 test rax,rax
 js text_error
 mov rax,[rsp+8]
 xor edx,edx
 mov ecx,1000000
 div rcx
 imul rcx,[rsp],1000
 add rax,rcx
 add rsp,16
 ret
.section .text.clock_ms_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 call [rip+iat8]
 leave
 ret

# Sleep milliseconds, 0..INT32_MAX, retry Linux nanosleep on EINTR.
.section .text.sleep_ms_linux,"ax"
 cmp rax,2147483647
 ja text_error
 xor edx,edx
 mov ecx,1000
 div rcx
 imul rdx,rdx,1000000
 sub rsp,32
 mov [rsp],rax
 mov [rsp+8],rdx
1: mov eax,35
 mov rdi,rsp
 lea rsi,[rsp+16]
 syscall
 test rax,rax
 jz 2f
 cmp rax,-4
 jne text_error
 mov rax,[rsp+16]
 mov rdx,[rsp+24]
 mov [rsp],rax
 mov [rsp+8],rdx
 jmp 1b
2: add rsp,32
 xor eax,eax
 ret
.section .text.sleep_ms_windows,"ax"
 cmp rax,2147483647
 ja text_error
 mov ecx,eax
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 call [rip+iat9]
 leave
 xor eax,eax
 ret

# Allocation-free bytewise prefix and suffix checks. Inputs RAX, RCX.
.section .text.has_prefix,"ax"
1: mov dl,[rcx]
 test dl,dl
 jz 2f
 cmp dl,[rax]
 jne 3f
 inc rax
 inc rcx
 jmp 1b
2: mov eax,1
 ret
3: xor eax,eax
 ret
.section .text.has_suffix,"ax"
 mov rdx,rax
 mov r8,rcx
1: cmp byte ptr [rdx],0
 je 2f
 inc rdx
 jmp 1b
2: cmp byte ptr [r8],0
 je 3f
 inc r8
 jmp 2b
3: sub r8,rcx
 sub rdx,rax
 cmp r8,rdx
 ja 5f
 sub rdx,r8
 add rax,rdx
4: mov dl,[rcx]
 test dl,dl
 jz 6f
 cmp dl,[rax]
 jne 5f
 inc rax
 inc rcx
 jmp 4b
5: xor eax,eax
 ret
6: mov eax,1
 ret
.section .text.tcp_exchange,"ax"
 # rax IPv4, rcx port, rdx request. One connection; response is bounded text.
 cmp rax,0
 jl net_error
 mov r8d,0xffffffff
 cmp rax,r8
 ja net_error
 cmp rcx,1
 jl net_error
 cmp rcx,65535
 ja net_error
 push rbp
 mov rbp,rsp
 sub rsp,48
 mov [rbp-8],rdx
 call net_connect
 mov [rbp-16],rax
 mov rax,[rbp-8]
 call strlen
 cmp rax,65535
 ja net_error
 mov [rbp-24],rax
 mov qword ptr [rbp-32],0
1: mov rdx,[rbp-24]
 sub rdx,[rbp-32]
 jz 2f
 mov rax,[rbp-16]
 mov rsi,[rbp-8]
 add rsi,[rbp-32]
 call net_send
 test rax,rax
 jle net_error
 add [rbp-32],rax
 jmp 1b
2: mov rax,[rbp-16]
 call net_shutdown
 mov eax,65536
 call alloc
 mov [rbp-40],rax
 mov qword ptr [rbp-48],0
3: mov rax,[rbp-16]
 mov rsi,[rbp-40]
 add rsi,[rbp-48]
 mov edx,65536
 sub rdx,[rbp-48]
 call net_recv
 test rax,rax
 js net_error
 jz 4f
 add [rbp-48],rax
 cmp qword ptr [rbp-48],65536
 jae net_error
 jmp 3b
4: mov rax,[rbp-16]
 call net_close
 mov rax,[rbp-40]
 mov rcx,[rbp-48]
 mov byte ptr [rax+rcx],0
 xor edx,edx
5: cmp rdx,rcx
 je 6f
 cmp byte ptr [rax+rdx],0
 je net_error
 inc rdx
 jmp 5b
6: mov edx,65535
 sub rdx,rcx
 lea r8,[rip+arena_state]
 sub [r8],rdx
 add [r8+8],rdx
 leave
 ret
.section .text.net_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,15
 jmp exit
1: .ascii "TRIPOL runtime error: TCP connection, transfer or response limit failed\n"
2:
.section .text.net_connect_linux,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,40
 mov qword ptr [rbp-16],0
 mov qword ptr [rbp-8],0
 mov word ptr [rbp-16],2
 bswap eax
 mov [rbp-12],eax
 rol cx,8
 mov [rbp-14],cx
 mov eax,41
 mov edi,2
 mov esi,1
 xor edx,edx
 syscall
 test rax,rax
 js net_error
 mov [rbp-24],rax
 # SO_RCVTIMEO and SO_SNDTIMEO (also bounds blocking connect on Linux).
 mov qword ptr [rbp-40],5
 mov qword ptr [rbp-32],0
 mov edi,eax
 mov eax,54
 mov esi,1
 mov edx,20
 lea r10,[rbp-40]
 mov r8d,16
 syscall
 test rax,rax
 js net_error
 mov edi,[rbp-24]
 mov eax,54
 mov esi,1
 mov edx,21
 lea r10,[rbp-40]
 mov r8d,16
 syscall
 test rax,rax
 js net_error
 mov edi,[rbp-24]
 mov eax,42
 lea rsi,[rbp-16]
 mov edx,16
 syscall
 test rax,rax
 js net_error
 mov rax,[rbp-24]
 leave
 ret
.section .text.net_send_linux,"ax"
 mov rdi,rax
 mov eax,44
 mov r10d,0x4000
 xor r8d,r8d
 xor r9d,r9d
 syscall
 ret
.section .text.net_recv_linux,"ax"
 mov rdi,rax
 mov eax,45
 xor r10d,r10d
 xor r8d,r8d
 xor r9d,r9d
 syscall
 ret
.section .text.net_shutdown_linux,"ax"
 mov rdi,rax
 mov eax,48
 mov esi,1
 syscall
 test rax,rax
 js net_error
 ret
.section .text.net_close_linux,"ax"
 mov rdi,rax
 mov eax,3
 syscall
 test rax,rax
 js net_error
 ret
.section .text.net_connect_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,512
 # WSAData at rbp-448, local sockaddr at rbp-16, handle at rbp-24.
 mov qword ptr [rbp-16],0
 mov qword ptr [rbp-8],0
 mov word ptr [rbp-16],2
 bswap eax
 mov [rbp-12],eax
 rol cx,8
 mov [rbp-14],cx
 mov ecx,0x0202
 lea rdx,[rbp-448]
 call qword ptr [rip+iat10]
 test eax,eax
 jnz net_error
 mov ecx,2
 mov edx,1
 mov r8d,6
 call qword ptr [rip+iat11]
 cmp rax,-1
 je net_error
 mov [rbp-24],rax
 # Winsock timeout value is milliseconds (DWORD).
 mov dword ptr [rbp-28],5000
 mov rcx,rax
 mov edx,0xffff
 mov r8d,0x1005
 lea r9,[rbp-28]
 mov qword ptr [rsp+32],4
 call qword ptr [rip+iat13]
 test eax,eax
 jnz net_error
 mov rcx,[rbp-24]
 mov edx,0xffff
 mov r8d,0x1006
 lea r9,[rbp-28]
 mov qword ptr [rsp+32],4
 call qword ptr [rip+iat13]
 test eax,eax
 jnz net_error
 mov rcx,[rbp-24]
 lea rdx,[rbp-16]
 mov r8d,16
 call qword ptr [rip+iat12]
 test eax,eax
 jnz net_error
 mov rax,[rbp-24]
 leave
 ret
.section .text.net_send_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rcx,rax
 mov r8d,edx
 mov rdx,rsi
 xor r9d,r9d
 call qword ptr [rip+iat14]
 cdqe
 leave
 ret
.section .text.net_recv_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rcx,rax
 mov r8d,edx
 mov rdx,rsi
 xor r9d,r9d
 call qword ptr [rip+iat16]
 cdqe
 leave
 ret
.section .text.net_shutdown_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rcx,rax
 mov edx,1
 call qword ptr [rip+iat15]
 test eax,eax
 jnz net_error
 leave
 ret
.section .text.net_close_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rcx,rax
 call qword ptr [rip+iat17]
 test eax,eax
 jnz net_error
 call qword ptr [rip+iat18]
 test eax,eax
 jnz net_error
 leave
 ret
.section .text.find_bytes,"ax"
 # rax haystack, rcx needle. Byte offset, -1 if not found.
 mov r8,rax
 xor eax,eax
 cmp byte ptr [rcx],0
 je 4f
1: cmp byte ptr [r8+rax],0
 je 5f
 lea rdx,[r8+rax]
 mov r9,rcx
2: mov r10b,[r9]
 test r10b,r10b
 jz 4f
 cmp [rdx],r10b
 jne 3f
 inc rdx
 inc r9
 jmp 2b
3: inc rax
 jmp 1b
4: ret
5: mov rax,-1
 ret

# Windows modal form: title, label, button, returns ANSI edit text or empty on cancel.
.section .text.ui_form_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,96
 mov [rsp+48],rax
 mov [rsp+56],rcx
 mov [rsp+64],rdx
 mov eax,4097
 call alloc
 mov byte ptr [rax],0
 mov [rsp+72],rax
 xor ecx,ecx
 call qword ptr [rip+iat40]
 mov rcx,rax
 lea rdx,[rip+ui_form_template]
 xor r8d,r8d
 lea r9,[rip+ui_form_proc]
 lea rax,[rsp+48]
 mov [rsp+32],rax
 call qword ptr [rip+iat20]
 cmp rax,-1
 je ui_error
 mov rax,[rsp+72]
 call strlen
 mov edx,4096
 sub rdx,rax
 lea r8,[rip+arena_state]
 sub [r8],rdx
 add [r8+8],rdx
 mov rax,[rsp+72]
 leave
 ret
.section .text.ui_form_linux,"ax"
 jmp ui_error
.section .text.ui_form_proc,"ax"
 # Microsoft x64 callback: hwnd, msg, wparam, lparam.
 push rbp
 mov rbp,rsp
 sub rsp,64
 mov [rsp+32],rcx
 mov [rsp+40],r9
 cmp edx,0x110
 je 1f
 cmp edx,0x111
 je 2f
 cmp edx,0x10
 je 4f
 xor eax,eax
 leave
 ret
1: mov edx,16
 mov r8,r9
 call qword ptr [rip+iat21]
 mov rcx,[rsp+32]
 mov rax,[rsp+40]
 mov rdx,[rax]
 call qword ptr [rip+iat23]
 mov rcx,[rsp+32]
 mov edx,100
 mov rax,[rsp+40]
 mov r8,[rax+8]
 call qword ptr [rip+iat24]
 mov rcx,[rsp+32]
 mov edx,1
 mov rax,[rsp+40]
 mov r8,[rax+16]
 call qword ptr [rip+iat24]
 mov eax,1
 leave
 ret
2: movzx eax,r8w
 cmp eax,2
 je 4f
 cmp eax,1
 jne 5f
 mov edx,16
 call qword ptr [rip+iat22]
 test rax,rax
 jz 4f
 mov rcx,[rsp+32]
 mov edx,101
 mov r8,[rax+24]
 mov r9d,4097
 call qword ptr [rip+iat25]
 mov rcx,[rsp+32]
 mov edx,1
 call qword ptr [rip+iat26]
 mov eax,1
 leave
 ret
4: mov rcx,[rsp+32]
 xor edx,edx
 call qword ptr [rip+iat26]
 mov eax,1
 leave
 ret
5: xor eax,eax
 leave
 ret

# Windows WinHTTP GET, strict OS certificate checks, bounded text, no redirects.
.section .text.https_get_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,2880
 mov [rsp+64],rax
 mov [rsp+72],rcx
 mov qword ptr [rsp+80],0
 mov qword ptr [rsp+88],0
 mov qword ptr [rsp+96],0
 # ASCII host/path to UTF-16, bounded by 253/1024 bytes.
 mov rsi,rax
 lea rdi,[rsp+256]
 xor ecx,ecx
1: movzx eax,byte ptr [rsi+rcx]
 mov [rdi+rcx*2],ax
 test eax,eax
 jz 2f
 inc ecx
 cmp ecx,253
 jbe 1b
 jmp 9f
2: mov rsi,[rsp+72]
 lea rdi,[rsp+768]
 xor ecx,ecx
3: movzx eax,byte ptr [rsi+rcx]
 mov [rdi+rcx*2],ax
 test eax,eax
 jz 4f
 inc ecx
 cmp ecx,1024
 jbe 3b
 jmp 9f
4: mov eax,65536
 call alloc
 mov [rsp+104],rax
 # Open with default proxy policy. No custom certificate exceptions.
 lea rcx,[rip+agent]
 xor edx,edx
 xor r8d,r8d
 xor r9d,r9d
 mov qword ptr [rsp+32],0
 call qword ptr [rip+iat30]
 test rax,rax
 jz 9f
 mov [rsp+80],rax
 mov rcx,rax
 mov edx,5000
 mov r8d,5000
 mov r9d,5000
 mov qword ptr [rsp+32],5000
 call qword ptr [rip+iat37]
 test eax,eax
 jz 9f
 mov rcx,[rsp+80]
 lea rdx,[rsp+256]
 mov r8d,443
 xor r9d,r9d
 call qword ptr [rip+iat31]
 test rax,rax
 jz 9f
 mov [rsp+88],rax
 mov rcx,rax
 lea rdx,[rip+method]
 lea r8,[rsp+768]
 xor r9d,r9d
 mov qword ptr [rsp+32],0
 mov qword ptr [rsp+40],0
 mov qword ptr [rsp+48],0x00800000
 call qword ptr [rip+iat32]
 test rax,rax
 jz 9f
 mov [rsp+96],rax
 # WINHTTP_OPTION_REDIRECT_POLICY=88, NEVER=0.
 mov rcx,rax
 mov edx,88
 mov dword ptr [rsp+144],0
 lea r8,[rsp+144]
 mov r9d,4
 call qword ptr [rip+iat38]
 test eax,eax
 jz 9f
 mov rcx,[rsp+96]
 xor edx,edx
 xor r8d,r8d
 xor r9d,r9d
 mov qword ptr [rsp+32],0
 mov qword ptr [rsp+40],0
 mov qword ptr [rsp+48],0
 call qword ptr [rip+iat33]
 test eax,eax
 jz 9f
 mov rcx,[rsp+96]
 xor edx,edx
 call qword ptr [rip+iat34]
 test eax,eax
 jz 9f
 # WINHTTP_QUERY_STATUS_CODE|WINHTTP_QUERY_FLAG_NUMBER.
 mov rcx,[rsp+96]
 mov edx,0x20000013
 xor r8d,r8d
 lea r9,[rsp+144]
 mov dword ptr [rsp+148],4
 lea rax,[rsp+148]
 mov [rsp+32],rax
 mov qword ptr [rsp+40],0
 call qword ptr [rip+iat39]
 test eax,eax
 jz 9f
 cmp dword ptr [rsp+144],200
 jb 9f
 cmp dword ptr [rsp+144],299
 ja 9f
 mov qword ptr [rsp+112],0
5: mov rcx,[rsp+96]
 mov rdx,[rsp+104]
 add rdx,[rsp+112]
 mov r8d,65536
 sub r8,[rsp+112]
 lea r9,[rsp+152]
 call qword ptr [rip+iat35]
 test eax,eax
 jz 9f
 mov eax,[rsp+152]
 test eax,eax
 jz 6f
 add [rsp+112],rax
 cmp qword ptr [rsp+112],65536
 jae 9f
 jmp 5b
6: mov rax,[rsp+104]
 mov rcx,[rsp+112]
 mov byte ptr [rax+rcx],0
 xor edx,edx
7: cmp rdx,rcx
 je 8f
 cmp byte ptr [rax+rdx],0
 je 9f
 inc rdx
 jmp 7b
8: mov qword ptr [rsp+120],1
 jmp 10f
9: mov qword ptr [rsp+120],0
10: mov rcx,[rsp+96]
 test rcx,rcx
 jz 11f
 call qword ptr [rip+iat36]
11: mov rcx,[rsp+88]
 test rcx,rcx
 jz 12f
 call qword ptr [rip+iat36]
12: mov rcx,[rsp+80]
 test rcx,rcx
 jz 13f
 call qword ptr [rip+iat36]
13: cmp qword ptr [rsp+120],0
 je https_error
 # Reclaim unused response buffer suffix.
 mov edx,65535
 sub rdx,[rsp+112]
 lea r8,[rip+arena_state]
 sub [r8],rdx
 add [r8+8],rdx
 mov rax,[rsp+104]
 leave
 ret
 .balign 2
agent: .short 84,82,73,80,79,76,0
method: .short 71,69,84,0
.section .text.https_get_linux,"ax"
 jmp https_error
.section .text.https_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,15
 jmp exit
1: .ascii "TRIPOL runtime error: HTTPS requires Windows, or request failed\n"
2:

# QUANTUM: one Windows window, one 32-bit staging buffer, no language arena references.
# Window is a built-in STATIC class, with the OS default message handler.
.section .text.game_open_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,128
 mov [rsp+104],rax
 mov [rsp+112],rcx
 mov [rsp+120],rdx
 cmp rcx,64
 jl game_error
 cmp rcx,2048
 jg game_error
 cmp rdx,64
 jl game_error
 cmp rdx,2048
 jg game_error
 call game_close
 xor ecx,ecx
 call qword ptr [rip+iat40]
 mov [rsp+80],rax
 xor ecx,ecx
 lea rdx,[rip+1f]
 mov r8,[rsp+104]
 mov r9d,0x10cf0000
 mov qword ptr [rsp+32],-2147483648
 mov qword ptr [rsp+40],-2147483648
 mov rax,[rsp+112]
 mov [rsp+48],rax
 mov rax,[rsp+120]
 mov [rsp+56],rax
 mov qword ptr [rsp+64],0
 mov qword ptr [rsp+72],0
 mov qword ptr [rsp+88],0
 call qword ptr [rip+iat50]
 test rax,rax
 jz game_error
 lea rdx,[rip+game_state]
 mov [rdx],rax
 xor ecx,ecx
 mov edx,262144
 mov r8d,0x3000
 mov r9d,4
 call qword ptr [rip+iat4]
 test rax,rax
 jz game_error
 lea rdx,[rip+game_state]
 mov [rdx+8],rax
 mov eax,1
 leave
 ret
1: .asciz "STATIC"
.section .text.game_close_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 lea rax,[rip+game_state]
 mov rcx,[rax]
 test rcx,rcx
 jz 1f
 call qword ptr [rip+iat51]
 test eax,eax
 jz 1f
 mov rcx,[rip+game_state]
 call qword ptr [rip+iat55]
1: lea rax,[rip+game_state]
 mov qword ptr [rax],0
 mov rcx,[rax+8]
 test rcx,rcx
 jz 2f
 xor edx,edx
 mov r8d,0x8000
 call qword ptr [rip+iat62]
 test eax,eax
 jz game_error
2: lea rdx,[rip+game_state]
 mov qword ptr [rdx+8],0
 mov eax,1
 leave
 ret
.section .text.game_running_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,96
1: lea rcx,[rsp+40]
 xor edx,edx
 xor r8d,r8d
 xor r9d,r9d
 mov qword ptr [rsp+32],1
 call qword ptr [rip+iat52]
 test eax,eax
 jz 3f
 cmp dword ptr [rsp+48],0x12
 je 4f
 cmp dword ptr [rsp+48],0x10
 jne 2f
 mov rax,[rsp+40]
 cmp rax,[rip+game_state]
 je 4f
2: lea rcx,[rsp+40]
 call qword ptr [rip+iat53]
 lea rcx,[rsp+40]
 call qword ptr [rip+iat54]
 jmp 1b
3: mov rcx,[rip+game_state]
 test rcx,rcx
 jz 4f
 call qword ptr [rip+iat51]
 test eax,eax
 jz 4f
 mov eax,1
 leave
 ret
4: call game_close
 xor eax,eax
 leave
 ret
.section .text.game_key_windows,"ax"
 cmp rax,0
 jl game_error
 cmp rax,255
 jg game_error
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,48
 mov [rsp+40],rax
 call qword ptr [rip+iat60]
 test rax,rax
 jz 1f
 cmp rax,[rip+game_state]
 jne 1f
 mov rcx,[rsp+40]
 call qword ptr [rip+iat59]
 shr eax,15
 and eax,1
 leave
 ret
1: xor eax,eax
 leave
 ret
.section .text.game_present_windows,"ax"
 # rax=array ptr, rcx=source width, rdx=source height; exact length, 0xRRGGBB.
 cmp rcx,1
 jl game_error
 cmp rcx,65535
 jg game_error
 cmp rdx,1
 jl game_error
 cmp rdx,65535
 jg game_error
 mov r8,rcx
 imul r8,rdx
 cmp r8,65535
 ja game_error
 cmp [rax],r8
 jne game_error
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,208
 mov [rsp+112],rcx
 mov [rsp+120],rdx
 mov rdi,[rip+game_buffer]
 test rdi,rdi
 jz game_error
 lea rsi,[rax+8]
 mov rcx,r8
 cld
1: lodsq
 cmp rax,16777215
 ja game_error
 stosd
 dec rcx
 jnz 1b
 mov rcx,[rip+game_state]
 test rcx,rcx
 jz game_error
 lea rdx,[rsp+136]
 call qword ptr [rip+iat58]
 test eax,eax
 jz game_error
 # skip minimized (zero area) without reporting failure
 cmp dword ptr [rsp+144],0
 jle 4f
 cmp dword ptr [rsp+148],0
 jle 4f
 mov rcx,[rip+game_state]
 call qword ptr [rip+iat56]
 test rax,rax
 jz game_error
 mov [rsp+128],rax
 # BITMAPINFOHEADER at 152: top-down 32-bit RGB DIB
 lea rdi,[rsp+152]
 xor eax,eax
 mov ecx,5
 cld
 rep stosq
 mov dword ptr [rsp+152],40
 mov eax,[rsp+112]
 mov [rsp+156],eax
 mov eax,[rsp+120]
 neg eax
 mov [rsp+160],eax
 mov dword ptr [rsp+164],0x00200001
 mov rcx,[rsp+128]
 xor edx,edx
 xor r8d,r8d
 mov r9d,[rsp+144]
 mov eax,[rsp+148]
 mov [rsp+32],rax
 mov qword ptr [rsp+40],0
 mov qword ptr [rsp+48],0
 mov rax,[rsp+112]
 mov [rsp+56],rax
 mov rax,[rsp+120]
 mov [rsp+64],rax
 mov rax,[rip+game_buffer]
 mov [rsp+72],rax
 lea rax,[rsp+152]
 mov [rsp+80],rax
 mov qword ptr [rsp+88],0
 mov qword ptr [rsp+96],0x00cc0020
 call qword ptr [rip+iat61]
 mov [rsp+192],eax
 mov rcx,[rip+game_state]
 mov rdx,[rsp+128]
 call qword ptr [rip+iat57]
 cmp dword ptr [rsp+192],0
 jle game_error
4: mov eax,1
 leave
 ret
.section .text.game_open_linux,"ax"
 jmp game_error
.section .text.game_close_linux,"ax"
 jmp game_error
.section .text.game_running_linux,"ax"
 jmp game_error
.section .text.game_key_linux,"ax"
 jmp game_error
.section .text.game_present_linux,"ax"
 jmp game_error
.section .text.game_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,17
 jmp exit
1: .ascii "TRIPOL runtime error: game window requires Windows or graphics operation failed\n"
2:
