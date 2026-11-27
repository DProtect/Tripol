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
.section .text.array_fill_range,"ax"
 # rax=array, rcx=start, rdx=count, r8=value; validate before writing.
 test rcx,rcx
 js bounds
 test rdx,rdx
 js bounds
 mov r9,[rax]
 cmp rcx,r9
 ja bounds
 sub r9,rcx
 cmp rdx,r9
 ja bounds
 lea rdi,[rax+rcx*8+8]
 mov rax,r8
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
 cmp edx,2
 je 1f
 mov esi,193 # O_WRONLY|O_CREAT|O_EXCL
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
 cmp edx,3
 je 4f
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
4: mov edx,0x40000000 # GENERIC_WRITE
 mov qword ptr [rsp+32],1 # CREATE_NEW
 jmp 3b
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
 cmp rax,4194304
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

# Nanosecond monotonic clock, not Unix time. On Windows QPC is converted
# using a 128-bit product to avoid overflowing before division.
.section .text.clock_ns_linux,"ax"
 sub rsp,16
 mov eax,228
 mov edi,1
 mov rsi,rsp
 syscall
 test rax,rax
 js text_error
 imul rax,[rsp],1000000000
 jo text_error
 add rax,[rsp+8]
 jo text_error
 add rsp,16
 ret
.section .text.clock_ns_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,48
 lea rcx,[rsp+32]
 call [rip+iat63]
 test eax,eax
 jz text_error
 lea rcx,[rsp+40]
 call [rip+iat64]
 test eax,eax
 jz text_error
 mov rax,[rsp+32]
 mov rcx,[rsp+40]
 test rax,rax
 js text_error
 test rcx,rcx
 jle text_error
 mov edx,1000000000
 mul rdx
 cmp rdx,rcx
 jae text_error
 div rcx
 test rax,rax
 js text_error
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
 mov edx,16777216
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
 cmp r8,4194304
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

# Encode a JSON string with one arena allocation and the same UTF-8 bytes.
.section .text.json_quote,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,16
 mov [rbp-8],rax
 mov rsi,rax
 mov eax,3
.jq_measure:
 mov dl,[rsi]
 test dl,dl
 jz .jq_allocate
 inc rsi
 inc eax
 cmp dl,32
 jb .jq_control
 cmp dl,34
 je .jq_extra
 cmp dl,92
 jne .jq_bound
.jq_extra:
 inc eax
 jmp .jq_bound
.jq_control:
 add eax,5
.jq_bound:
 cmp eax,65536
 ja text_error
 jmp .jq_measure
.jq_allocate:
 call alloc
 mov [rbp-16],rax
 mov rdi,rax
 mov rsi,[rbp-8]
 lea r8,[rip+.jq_digits]
 mov byte ptr [rdi],34
 inc rdi
.jq_pass:
 movzx edx,byte ptr [rsi]
 test dl,dl
 jz .jq_finish
 inc rsi
 cmp dl,32
 jb .jq_escape_control
 cmp dl,34
 je .jq_escape
 cmp dl,92
 jne .jq_copy
.jq_escape:
 mov byte ptr [rdi],92
 inc rdi
.jq_copy:
 mov [rdi],dl
 inc rdi
 jmp .jq_pass
.jq_escape_control:
 mov dword ptr [rdi],0x3030755c
 mov eax,edx
 shr eax,4
 and edx,15
 mov al,[r8+rax]
 mov dl,[r8+rdx]
 mov [rdi+4],al
 mov [rdi+5],dl
 add rdi,6
 jmp .jq_pass
.jq_finish:
 mov word ptr [rdi],34
 mov rax,[rbp-16]
 leave
 ret
.jq_digits: .ascii "0123456789abcdef"

# Encode one non-NUL Unicode scalar. Use a four-byte local staging buffer.
.section .text.utf8_char,"ax"
 cmp rax,0
 jle text_error
 cmp rax,0x10ffff
 ja text_error
 cmp eax,0xd800
 jb .uc_checked
 cmp eax,0xdfff
 jbe text_error
.uc_checked:
 push rbp
 mov rbp,rsp
 sub rsp,16
 mov r8d,1
 cmp eax,0x80
 jb .uc_first
 inc r8d
 cmp eax,0x800
 jb .uc_two
 inc r8d
 cmp eax,0x10000
 jb .uc_three
 inc r8d
 mov r9d,240
 jmp .uc_tail
.uc_two:
 mov r9d,192
 jmp .uc_tail
.uc_three:
 mov r9d,224
.uc_tail:
 mov ecx,r8d
 dec ecx
.uc_cont:
 mov edx,eax
 and edx,63
 or edx,128
 mov [rbp+rcx-8],dl
 shr eax,6
 dec ecx
 jnz .uc_cont
 or eax,r9d
.uc_first:
 mov [rbp-8],al
 mov byte ptr [rbp+r8-8],0
 inc r8d
 mov [rbp-16],r8
 mov eax,r8d
 call alloc
 mov rdi,rax
 mov rcx,[rbp-16]
 lea rsi,[rbp-8]
 cld
 rep movsb
 leave
 ret

.section .text.discord_edit_windows,"ax"
 mov ecx,1
 jmp discord_request

.section .text.discord_callback_windows,"ax"
 xor ecx,ecx
 jmp discord_request

.section .text.discord_request_windows,"ax"
 xor r8d,r8d
 jmp discord_rest
.section .text.discord_rest_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,2048
 test byte ptr [rsp],0
 sub rsp,2048
 test byte ptr [rsp],0
 mov [rsp+72],rax # path
 mov [rsp+64],rcx # token address, 0 for anonymous POST, 1 for anonymous PATCH/DELETE
 mov [rsp+160],rdx # body
 mov [rsp+176],r8 # method: 0 legacy, 1 GET, 2 POST, 3 PUT, 4 PATCH, 5 DELETE
 mov qword ptr [rsp+168],0 # body length
 mov qword ptr [rsp+80],0
 mov qword ptr [rsp+88],0
 mov qword ptr [rsp+96],0
 mov qword ptr [rsp+104],0
 cmp r8,5
 ja 9f
 # Fixed trusted host, path ASCII and bounded, token restricted to safe chars.
 mov rsi,[rsp+72]
 lea rdi,[rsp+256]
 xor ecx,ecx
 cmp byte ptr [rsi],47
 jne 9f
3: movzx eax,byte ptr [rsi+rcx]
 mov [rdi+rcx*2],ax
 test eax,eax
 jz 4f
 cmp eax,33
 jb 9f
 cmp eax,126
 ja 9f
 inc ecx
 cmp ecx,1024
 jbe 3b
 jmp 9f
4: # Build wide Authorization header; no CR/LF in token.
 lea rdi,[rsp+2400]
 cmp qword ptr [rsp+64],1
 jbe .dcontent_type
 lea rsi,[rip+discord_header]
 mov ecx,19
 rep movsw
 mov rsi,[rsp+64]
 xor ecx,ecx
 cmp byte ptr [rsi],0
 je 9f
.dtoken:
 movzx eax,byte ptr [rsi+rcx]
 test eax,eax
 jz .dheader_end
 cmp eax,45
 je .dtok_ok
 cmp eax,46
 je .dtok_ok
 cmp eax,95
 je .dtok_ok
 cmp eax,48
 jb 9f
 cmp eax,57
 jbe .dtok_ok
 cmp eax,65
 jb 9f
 cmp eax,90
 jbe .dtok_ok
 cmp eax,97
 jb 9f
 cmp eax,122
 ja 9f
.dtok_ok:
 mov [rdi+rcx*2],ax
 inc ecx
 cmp ecx,512
 jbe .dtoken
 jmp 9f
.dheader_end:
 lea rdi,[rdi+rcx*2]
 lea rsi,[rip+discord_header_end]
 mov ecx,2
 rep movsw
.dcontent_type:
 lea rsi,[rip+discord_content_type]
 mov ecx,33
 rep movsw
 mov rax,[rsp+160]
 call strlen
 cmp rax,65535
 ja 9f
 mov [rsp+168],rax
 mov eax,65536
 call alloc
 mov [rsp+104],rax
 # Open with default proxy policy. No custom certificate exceptions.
 lea rcx,[rip+discord_agent]
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
 cmp qword ptr [rsp+64],1
 ja .dtimeouts
 mov edx,2000
.dtimeouts:
 mov r8d,edx
 mov r9d,edx
 mov [rsp+32],rdx
 call qword ptr [rip+iat37]
 test eax,eax
 jz 9f
 mov rcx,[rsp+80]
 lea rdx,[rip+discord_host]
 mov r8d,443
 xor r9d,r9d
 call qword ptr [rip+iat31]
 test rax,rax
 jz 9f
 mov [rsp+88],rax
 mov rcx,rax
 mov eax,[rsp+176]
 lea rdx,[rip+discord_get]
 cmp eax,1
 je .dverb
 cmp eax,2
 je .dpost
 cmp eax,3
 je .dput
 cmp eax,4
 je .dpatch
 cmp eax,5
 je .ddelete
 cmp qword ptr [rsp+64],1
 jne .dregular
 cmp qword ptr [rsp+168],0
 je .ddelete
 jmp .dpatch
.dregular:
 cmp qword ptr [rsp+168],0
 je .dverb
.dpost:
 lea rdx,[rip+discord_post]
 jmp .dverb
.dput:
 lea rdx,[rip+discord_put]
 jmp .dverb
.dpatch:
 lea rdx,[rip+discord_patch]
 jmp .dverb
.ddelete:
 lea rdx,[rip+discord_delete]
.dverb:
 lea r8,[rsp+256]
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
 lea rdx,[rsp+2400]
 mov r8d,-1
 mov r9,[rsp+160]
 mov rax,[rsp+168]
 test rax,rax
 jnz .dbody
 xor r9d,r9d
.dbody:
 mov [rsp+32],rax
 mov [rsp+40],rax
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
 jne .dresult
 cmp qword ptr [rsp+176],0
 jne .drecover
 cmp qword ptr [rsp+64],1
 ja discord_error
.drecover:
 cmp qword ptr [rsp+104],0
 je .dfail
 lea r8,[rip+arena_state]
 sub qword ptr [r8],65536
 add qword ptr [r8+8],65536
.dfail:
 lea rax,[rip+.dempty]
 leave
 ret
.dresult:
 # Reclaim unused response buffer suffix.
 mov edx,65535
 sub rdx,[rsp+112]
 lea r8,[rip+arena_state]
 sub [r8],rdx
 add [r8+8],rdx
 mov rax,[rsp+104]
 cmp qword ptr [rsp+176],0
 jne .dno_content
 cmp qword ptr [rsp+64],1
 ja .dreturn
.dno_content:
 cmp qword ptr [rsp+112],0
 jne .dreturn
 lea rax,[rip+.daccepted]
.dreturn:
 leave
 ret
.dempty: .byte 0
.daccepted: .asciz "{}"
 .balign 2

discord_header: .short 65,117,116,104,111,114,105,122,97,116,105,111,110,58,32,66,111,116,32
discord_header_end: .short 13,10
discord_content_type: .short 67,111,110,116,101,110,116,45,84,121,112,101,58,32,97,112,112,108,105,99,97,116,105,111,110,47,106,115,111,110,13,10,0
discord_host: .short 100,105,115,99,111,114,100,46,99,111,109,0
discord_agent: .short 84,82,73,80,79,76,0
discord_get: .short 71,69,84,0
discord_post: .short 80,79,83,84,0
discord_put: .short 80,85,84,0
discord_patch: .short 80,65,84,67,72,0
discord_delete: .short 68,69,76,69,84,69,0

.section .text.discord_edit_linux,"ax"
 mov r8,rax
 xor r9d,r9d
 mov r10,rdx
 mov eax,6
 jmp discord_pipe
.section .text.discord_callback_linux,"ax"
 mov r8,rax
 xor r9d,r9d
 mov r10,rdx
 mov eax,5
 jmp discord_pipe
.section .text.discord_request_linux,"ax"
 mov r8,rax
 mov r9,rcx
 mov r10,rdx
 mov eax,4
 jmp discord_pipe
.section .text.discord_rest_linux,"ax"
 cmp r8,1
 jb .drest_empty
 cmp r8,5
 ja .drest_empty
 lea r11d,[r8+6] # Bridge operations 7 through 11.
 mov r8,rax
 mov r9,rcx
 mov r10,rdx
 mov eax,r11d
 jmp discord_pipe
.drest_empty:
 lea rax,[rip+.drest_zero]
 ret
.drest_zero: .byte 0
.section .text.discord_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,15
 jmp exit
1: .ascii "TRIPOL Discord request failed (Windows HTTPS, token, permissions, rate limit or response limit). Stopped; do not retry blindly.\n"
2:

.section .text.env_get_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,2048
 test byte ptr [rsp],0
 sub rsp,2112
 test byte ptr [rsp],0
 mov rcx,rax
 lea rdx,[rsp+64]
 mov r8d,4096
 call qword ptr [rip+iat65]
 cmp eax,4096
 jae 1f
 mov [rsp+40],rax
 lea rax,[rax+1]
 call alloc
 mov rcx,[rsp+40]
 lea rsi,[rsp+64]
 mov rdi,rax
 rep movsb
 mov byte ptr [rdi],0
 leave
 ret
1: xor eax,eax
 inc eax
 call alloc
 mov byte ptr [rax],0
 leave
 ret
.section .text.env_get_linux,"ax"
 cmp byte ptr [rax],0
 je .environ_empty
 mov r8,rax
 mov r9,[rip+linux_env]
.environ_next:
 mov rax,[r9]
 add r9,8
 test rax,rax
 jz .environ_empty
 mov rdx,r8
.environ_compare:
 mov cl,[rdx]
 test cl,cl
 jz .environ_match
 cmp cl,[rax]
 jne .environ_next
 inc rax
 inc rdx
 jmp .environ_compare
.environ_match:
 cmp byte ptr [rax],61
 jne .environ_next
 inc rax
 mov rdi,rax
 mov ecx,4096
 xor edx,edx
.environ_bound:
 cmp byte ptr [rdi+rdx],0
 je .environ_found
 inc rdx
 dec ecx
 jnz .environ_bound
.environ_empty:
 lea rax,[rip+.environ_none]
.environ_found:
 ret
.environ_none: .byte 0

# SSE2 equality of dynamic int64 arrays. Two elements per vector; scalar tail.
.section .text.array_equal,"ax"
 mov r8,[rax]
 cmp r8,[rcx]
 jne 4f
 add rax,8
 add rcx,8
1: cmp r8,2
 jb 2f
 movdqu xmm0,[rax]
 movdqu xmm1,[rcx]
 pcmpeqd xmm0,xmm1
 pmovmskb edx,xmm0
 cmp edx,65535
 jne 4f
 add rax,16
 add rcx,16
 sub r8,2
 jmp 1b
2: test r8,r8
 jz 3f
 mov rdx,[rax]
 cmp rdx,[rcx]
 jne 4f
3: mov eax,1
 ret
4: xor eax,eax
 ret

# Bounded state slots. Never retain references to temporary arena objects.
.section .text.bot_slot_set,"ax"
 cmp rax,38
 ja text_error
 imul rax,rax,528
 lea rdi,[rip+bot_state]
 add rdi,rax
 mov rsi,rcx
 xor edx,edx
1: cmp byte ptr [rsi+rdx],0
 je 2f
 inc rdx
 cmp rdx,512
 jbe 1b
 jmp text_error
2: lea rcx,[rdx+1]
 cld
 rep movsb
 # Clear old suffix, including secrets, on replacement.
 mov ecx,527
 sub rcx,rdx
 xor eax,eax
 rep stosb
 mov eax,1
 ret

.section .text.bot_slot_get,"ax"
 cmp rax,38
 ja text_error
 imul rax,rax,528
 lea rdx,[rip+bot_state]
 add rax,rdx
 push rbp
 mov rbp,rsp
 sub rsp,16
 mov [rbp-8],rax
 call strlen
 inc rax
 mov [rbp-16],rax
 call alloc
 mov rdi,rax
 mov rsi,[rbp-8]
 mov rcx,[rbp-16]
 cld
 rep movsb
 leave
 ret

# Copy independent int array in one allocation, no redundant zero fill.
.section .text.array_copy,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,16
 mov [rbp-8],rax
 mov rcx,[rax]
 cmp rcx,4194304
 ja text_error
 inc rcx
 mov [rbp-16],rcx
 lea rax,[rcx*8]
 call alloc
 mov rdi,rax
 mov rsi,[rbp-8]
 mov rcx,[rbp-16]
 cld
 rep movsq
 leave
 ret

# Reuse destination storage; reject different lengths before touching data.
.section .text.array_copy_into,"ax"
 mov rdx,[rax]
 cmp rdx,[rcx]
 jne .aci_bad
 cmp rax,rcx
 je .aci_ok
 lea rdi,[rax+8]
 lea rsi,[rcx+8]
 mov rcx,rdx
 cld
 rep movsq
.aci_ok:
 mov eax,1
 ret
.aci_bad:
 xor eax,eax
 ret

# Find first matching int64; bounded and read-only.
# AVX2 first-match search with CPUID/XGETBV gate; scalar fallback and tail.
.section .text.array_find,"ax"
 mov r8,[rax]
 test r8,r8
 jz .af_none
 lea r10,[rax+8]
 mov r9,rcx
 mov rdi,r10
 cmp r8,32
 jb .af_scalar
 call cpu_has_avx2
 test eax,eax
 jz .af_scalar
 vmovq xmm0,r9
 vpbroadcastq ymm0,xmm0
 mov rcx,r8
.af_pair:
 cmp rcx,8
 jb .af_vector
 vpcmpeqq ymm1,ymm0,[rdi]
 vpcmpeqq ymm2,ymm0,[rdi+32]
 vpor ymm3,ymm1,ymm2
 vptest ymm3,ymm3
 jnz .af_pair_hit
 add rdi,64
 sub rcx,8
 jmp .af_pair
.af_pair_hit:
 vpmovmskb edx,ymm1
 test edx,edx
 jnz .af_hit
 add rdi,32
 vpmovmskb edx,ymm2
 jmp .af_hit
.af_vector:
 cmp rcx,4
 jb .af_tail
 vpcmpeqq ymm1,ymm0,[rdi]
 vpmovmskb edx,ymm1
 test edx,edx
 jnz .af_hit
 add rdi,32
 sub rcx,4
 jmp .af_vector
.af_hit:
 bsf edx,edx
 shr edx,3
 mov rax,rdi
 sub rax,r10
 shr rax,3
 add rax,rdx
 vzeroupper
 ret
.af_tail:
 vzeroupper
 test rcx,rcx
 jz .af_none
 mov rax,r9
 jmp .af_scan
.af_scalar:
 mov rax,r9
 mov rcx,r8
.af_scan:
 cld
 repne scasq
 jne .af_none
 lea rax,[rdi-8]
 sub rax,r10
 shr rax,3
 ret
.af_none:
 mov rax,-1
 ret

# Returns EAX=1 only when the processor and operating system allow AVX2.
# Clobbers RAX/RCX/RDX; preserves RBX and R8 through R11.
.section .text.cpu_has_avx2,"ax"
 push rbx
 xor eax,eax
 cpuid
 cmp eax,7
 jb .no_avx2
 mov eax,1
 cpuid
 and ecx,0x18000000 # AVX + OSXSAVE
 cmp ecx,0x18000000
 jne .no_avx2
 xor ecx,ecx
 xgetbv
 and eax,6
 cmp eax,6
 jne .no_avx2
 mov eax,7
 xor ecx,ecx
 cpuid
 test ebx,32
 setnz al
 movzx eax,al
 pop rbx
 ret
.no_avx2:
 xor eax,eax
 pop rbx
 ret

# Fill exactly RCX int64 elements at RDX with RAX. Preserve the caller's
# bounds flag in R8. Only enter the AVX2 path for sufficiently long fills.
.section .text.vector_fill,"ax"
 cmp rcx,32
 jb .vf_sse
 mov r9,rax
 mov r10,rdx
 mov r11,rcx
 call cpu_has_avx2
 test eax,eax
 mov rax,r9
 mov rdx,r10
 mov rcx,r11
 jz .vf_sse
 vmovq xmm0,rax
 vpbroadcastq ymm0,xmm0
.vf_avx_loop:
 vmovdqu [rdx],ymm0
 add rdx,32
 sub rcx,4
 cmp rcx,4
 jae .vf_avx_loop
 vzeroupper
.vf_sse:
 movq xmm0,rax
 punpcklqdq xmm0,xmm0
.vf_sse_loop:
 cmp rcx,2
 jb .vf_tail
 movdqu [rdx],xmm0
 add rdx,16
 sub rcx,2
 jmp .vf_sse_loop
.vf_tail:
 test rcx,rcx
 jz .vf_done
 mov [rdx],rax
.vf_done:
 ret

# NEXUS Gateway. Fixed TLS host, no redirects. One synchronous receive worker
# at a time; no arena use on worker. Main waits locally, not network polling.
# gw_state: session, connection, websocket, thread, signed length, padding, data.
.section .text.gw_open_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,96
 lea rax,[rip+gw_state]
 mov [rsp+72],rax
 cmp qword ptr [rax],0
 jne .gwopen_bad
 mov qword ptr [rsp+64],0
 lea rcx,[rip+.gwagent]
 xor edx,edx
 xor r8d,r8d
 xor r9d,r9d
 mov qword ptr [rsp+32],0
 call qword ptr [rip+iat30]
 test rax,rax
 jz .gwopen_bad
 mov r10,[rsp+72]
 mov [r10],rax
 mov rcx,rax
 mov edx,10000
 mov r8d,10000
 mov r9d,10000
 mov qword ptr [rsp+32],120000
 call qword ptr [rip+iat37]
 test eax,eax
 jz .gwopen_fail
 mov r10,[rsp+72]
 mov rcx,[r10]
 lea rdx,[rip+.gwhost]
 mov r8d,443
 xor r9d,r9d
 call qword ptr [rip+iat31]
 test rax,rax
 jz .gwopen_fail
 mov r10,[rsp+72]
 mov [r10+8],rax
 mov rcx,rax
 lea rdx,[rip+.gwget]
 lea r8,[rip+.gwpath]
 xor r9d,r9d
 mov qword ptr [rsp+32],0
 mov qword ptr [rsp+40],0
 mov qword ptr [rsp+48],0x00800000
 call qword ptr [rip+iat32]
 test rax,rax
 jz .gwopen_fail
 mov [rsp+64],rax
 mov rcx,rax
 mov edx,88
 mov dword ptr [rsp+80],0
 lea r8,[rsp+80]
 mov r9d,4
 call qword ptr [rip+iat38]
 test eax,eax
 jz .gwopen_fail
 mov rcx,[rsp+64]
 mov edx,114
 xor r8d,r8d
 xor r9d,r9d
 call qword ptr [rip+iat38]
 test eax,eax
 jz .gwopen_fail
 mov rcx,[rsp+64]
 xor edx,edx
 xor r8d,r8d
 xor r9d,r9d
 mov qword ptr [rsp+32],0
 mov qword ptr [rsp+40],0
 mov qword ptr [rsp+48],0
 call qword ptr [rip+iat33]
 test eax,eax
 jz .gwopen_fail
 mov rcx,[rsp+64]
 xor edx,edx
 call qword ptr [rip+iat34]
 test eax,eax
 jz .gwopen_fail
 mov rcx,[rsp+64]
 xor edx,edx
 call qword ptr [rip+iat66]
 test rax,rax
 jz .gwopen_fail
 mov r10,[rsp+72]
 mov [r10+16],rax
 mov rcx,[rsp+64]
 call qword ptr [rip+iat36]
 mov eax,1
 leave
 ret
.gwopen_fail:
 mov rcx,[rsp+64]
 test rcx,rcx
 jz .gwopen_cleanup
 call qword ptr [rip+iat36]
.gwopen_cleanup:
 call gw_close
.gwopen_bad:
 xor eax,eax
 leave
 ret
 .balign 2
.gwagent: .short 84,82,73,80,79,76,0
.gwhost: .short 103,97,116,101,119,97,121,46,100,105,115,99,111,114,100,46,103,103,0
.gwget: .short 71,69,84,0
.gwpath: .short 47,63,118,61,49,48,38,101,110,99,111,100,105,110,103,61,106,115,111,110,0

.section .text.gw_send_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,48
 mov [rsp+32],rax
 call strlen
 cmp rax,4096
 ja gw_abort
 mov r9d,eax
 mov r8,[rsp+32]
 lea rcx,[rip+gw_state]
 mov rcx,[rcx+16]
 test rcx,rcx
 jz gw_abort
 mov edx,2
 call qword ptr [rip+iat67]
 test eax,eax
 jnz gw_abort
 mov eax,1
 leave
 ret

# The OS thread uses Win64 ABI, preserves nonvolatile registers, and writes
# only its own bounded buffer. Completion is synchronized by thread wait.
.section .text.gw_worker,"ax"
 push rbx
 sub rsp,64
 mov rbx,rcx
 mov qword ptr [rbx+32],0
.gww_loop:
 mov rcx,[rbx+16]
 lea rdx,[rbx+48]
 add rdx,[rbx+32]
 mov r8d,65535
 sub r8d,[rbx+32]
 jz .gww_bad
 lea r9,[rsp+40]
 lea rax,[rsp+44]
 mov [rsp+32],rax
 call qword ptr [rip+iat68]
 test eax,eax
 jnz .gww_bad
 mov eax,[rsp+40]
 add [rbx+32],rax
 cmp qword ptr [rbx+32],65535
 ja .gww_bad
 cmp dword ptr [rsp+44],3
 je .gww_loop
 cmp dword ptr [rsp+44],2
 jne .gww_bad
 mov rax,[rbx+32]
 test rax,rax
 jz .gww_bad
 mov byte ptr [rbx+rax+48],0
 xor ecx,ecx
.gww_nul:
 cmp rcx,rax
 je .gww_done
 cmp byte ptr [rbx+rcx+48],0
 je .gww_bad
 inc rcx
 jmp .gww_nul
.gww_bad:
 mov qword ptr [rbx+32],-1
.gww_done:
 xor eax,eax
 add rsp,64
 pop rbx
 ret

.section .text.gw_receive_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,80
 cmp rax,1000
 ja gw_abort
 mov [rsp+56],rax
 lea r10,[rip+gw_state]
 mov [rsp+64],r10
 cmp qword ptr [r10+16],0
 je gw_abort
 cmp qword ptr [r10+24],0
 jne .gwr_wait
 xor ecx,ecx
 xor edx,edx
 lea r8,[rip+gw_worker]
 mov r9,r10
 mov qword ptr [rsp+32],0
 mov qword ptr [rsp+40],0
 call qword ptr [rip+iat70]
 test rax,rax
 jz gw_abort
 mov r10,[rsp+64]
 mov [r10+24],rax
.gwr_wait:
 mov rcx,[r10+24]
 mov edx,[rsp+56]
 call qword ptr [rip+iat71]
 cmp eax,258
 je .gwr_empty
 test eax,eax
 jnz gw_abort
 mov r10,[rsp+64]
 mov rcx,[r10+24]
 call qword ptr [rip+iat7]
 mov r10,[rsp+64]
 mov qword ptr [r10+24],0
 mov rax,[r10+32]
 test rax,rax
 js .gwr_failed
 inc rax
 mov [rsp+56],rax
 call alloc
 mov rdi,rax
 mov r10,[rsp+64]
 lea rsi,[r10+48]
 mov rcx,[rsp+56]
 cld
 rep movsb
 leave
 ret
.gwr_empty:
 lea rax,[rip+.gwempty]
 leave
 ret
.gwr_failed:
 lea rax,[rip+.gwfailed]
 leave
 ret
.gwempty: .byte 0
.gwfailed: .asciz "!"

.section .text.gw_close_windows,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,48
 lea rax,[rip+gw_state]
 mov [rsp+32],rax
 cmp qword ptr [rax+24],0
 jne gw_abort
 mov rcx,[rax+16]
 test rcx,rcx
 jz .gwc_conn
 call qword ptr [rip+iat36]
.gwc_conn:
 mov rax,[rsp+32]
 mov rcx,[rax+8]
 test rcx,rcx
 jz .gwc_session
 call qword ptr [rip+iat36]
.gwc_session:
 mov rax,[rsp+32]
 mov rcx,[rax]
 test rcx,rcx
 jz .gwc_done
 call qword ptr [rip+iat36]
.gwc_done:
 mov rax,[rsp+32]
 mov qword ptr [rax],0
 mov qword ptr [rax+8],0
 mov qword ptr [rax+16],0
 mov eax,1
 leave
 ret

# Linux launcher requests: op, three byte lengths, then unpadded text.
# Replies: one length followed by text. 0xffffffff denotes a diagnosed error.
.section .text.discord_pipe,"ax"
 push rbp
 mov rbp,rsp
 sub rsp,64
 mov [rbp-16],eax
 mov [rbp-32],r8
 mov [rbp-40],r9
 mov [rbp-48],r10
 lea rax,[rip+.discord_pipe_env]
 call env_get
 cmp byte ptr [rax],49
 jne .discord_pipe_missing
 cmp byte ptr [rax+1],0
 jne .discord_pipe_missing
 mov qword ptr [rbp-56],0
.discord_pipe_lengths:
 mov rdi,[rbp-56]
 lea rsi,[rbp-32]
 mov rax,[rsi+rdi]
 test rax,rax
 jz .discord_pipe_length_ready
 call strlen
 cmp rax,65535
 ja gw_abort
.discord_pipe_length_ready:
 mov rdi,[rbp-56]
 neg rdi
 shr rdi,1
 mov [rbp+rdi-12],eax
 sub qword ptr [rbp-56],8
 cmp qword ptr [rbp-56],-24
 jne .discord_pipe_lengths
 mov edi,3
 lea rsi,[rbp-16]
 mov edx,16
 mov r8d,1
 call discord_pipe_io
 mov qword ptr [rbp-56],0
.discord_pipe_arguments:
 mov rcx,[rbp-56]
 mov r10,rcx
 neg rcx
 shr rcx,1
 mov edx,[rbp+rcx-12]
 test edx,edx
 jz .discord_pipe_next
 mov rsi,[rbp+r10-32]
 mov edi,3
 mov r8d,1
 call discord_pipe_io
.discord_pipe_next:
 sub qword ptr [rbp-56],8
 cmp qword ptr [rbp-56],-24
 jne .discord_pipe_arguments
 mov edi,4
 lea rsi,[rbp-16]
 mov edx,4
 xor r8d,r8d
 call discord_pipe_io
 mov eax,[rbp-16]
 cmp eax,65535
 ja gw_abort
 test eax,eax
 jz .discord_pipe_empty
 inc eax
 call alloc
 mov [rbp-24],rax
 mov rsi,rax
 mov edi,4
 mov edx,[rbp-16]
 xor r8d,r8d
 call discord_pipe_io
 mov rax,[rbp-24]
 mov edx,[rbp-16]
 mov byte ptr [rax+rdx],0
 leave
 ret
.discord_pipe_empty:
 lea rax,[rip+.discord_pipe_none]
 leave
 ret
.discord_pipe_missing:
 lea rsi,[rip+.discord_pipe_message]
 mov edx,.discord_pipe_message_end-.discord_pipe_message
 call writeerr
 mov eax,15
 jmp exit
.discord_pipe_env: .asciz "TRIPOL_DISCORD_PIPE"
.discord_pipe_none: .byte 0
.discord_pipe_message: .ascii "Discord transport on Linux requires the TRIPOL launcher. Run the source with tripol run.\n"
.discord_pipe_message_end:

.section .text.discord_pipe_io,"ax"
 mov eax,r8d
 syscall
 cmp rax,-4
 je discord_pipe_io
 test rax,rax
 jle gw_abort
 add rsi,rax
 sub rdx,rax
 jnz discord_pipe_io
 ret

.section .text.gw_host_start,"ax"
 mov r8,rax
 xor r9d,r9d
 xor r10d,r10d
 mov eax,1
 jmp discord_pipe
.section .text.gw_host_poll,"ax"
 xor r8d,r8d
 xor r9d,r9d
 xor r10d,r10d
 mov eax,2
 jmp discord_pipe
.section .text.gw_host_stop,"ax"
 xor r8d,r8d
 xor r9d,r9d
 xor r10d,r10d
 mov eax,3
 jmp discord_pipe
.section .text.gw_open_linux,"ax"
 xor eax,eax
 ret
.section .text.gw_close_linux,"ax"
 mov eax,1
 ret
.section .text.gw_send_linux,"ax"
 jmp gw_abort
.section .text.gw_receive_linux,"ax"
 jmp gw_abort
.section .text.gw_abort,"ax"
 lea rsi,[rip+.gwerr]
 mov edx,.gwerrend-.gwerr
 call writeerr
 mov eax,15
 jmp exit
.gwerr: .ascii "TRIPOL Discord Gateway stopped (transport, heartbeat or protocol error). Check connection and bot configuration.\n"
.gwerrend:


.section .text.bytes_write,"ax"
# rax=path, rcx=ints, edx=open mode; validate before opening the file
 push rbp
 mov rbp,rsp
 sub rsp,320
 mov [rbp-8],rax
 mov [rbp-16],rcx
 mov [rbp-64],edx
 mov rdx,[rcx]
 cmp rdx,65536
 ja file_error
 mov [rbp-24],rdx
 xor r8d,r8d
1: cmp r8,rdx
 jae 2f
 cmp qword ptr [rcx+r8*8+8],255
 ja file_error
 inc r8
 jmp 1b
2: mov rax,[rbp-8]
 mov edx,[rbp-64]
 call file_open
 mov [rbp-32],rax
 mov qword ptr [rbp-40],0
3: mov rdx,[rbp-24]
 sub rdx,[rbp-40]
 jz 7f
 cmp rdx,256
 jbe 4f
 mov edx,256
4: mov [rbp-48],rdx
 mov rcx,[rbp-16]
 mov r8,[rbp-40]
 lea rcx,[rcx+r8*8+8]
 lea rsi,[rbp-320]
 xor r8d,r8d
5: mov rax,[rcx+r8*8]
 mov [rsi+r8],al
 inc r8
 cmp r8,rdx
 jb 5b
 mov qword ptr [rbp-56],0
6: mov rax,[rbp-32]
 lea rsi,[rbp-320]
 add rsi,[rbp-56]
 mov rdx,[rbp-48]
 sub rdx,[rbp-56]
 call file_chunk_write
 test rax,rax
 jle file_error
 add [rbp-56],rax
 mov rcx,[rbp-56]
 cmp rcx,[rbp-48]
 jb 6b
 add [rbp-40],rcx
 jmp 3b
7: mov rax,[rbp-32]
 call file_close
 mov rax,[rbp-24]
 leave
 ret
.section .text.bytes_read,"ax"
# rax=path, rcx=maximum bytes
 push rbp
 mov rbp,rsp
 sub rsp,320
 mov [rbp-8],rax
 cmp rcx,65536
 ja file_error
 mov [rbp-16],rcx
 mov rax,rcx
 call ints_alloc
 mov [rbp-24],rax
 mov rax,[rbp-8]
 xor edx,edx
 call file_open
 mov [rbp-32],rax
 mov qword ptr [rbp-40],0
1: mov rdx,[rbp-16]
 sub rdx,[rbp-40]
 jz 5f
 cmp rdx,256
 jbe 2f
 mov edx,256
2: mov rax,[rbp-32]
 lea rsi,[rbp-320]
 call file_chunk_read
 test rax,rax
 js file_error
 jz 6f
 mov rcx,[rbp-24]
 mov r8,[rbp-40]
 lea rcx,[rcx+r8*8+8]
 lea rsi,[rbp-320]
 xor edx,edx
3: movzx r8d,byte ptr [rsi+rdx]
 mov [rcx+rdx*8],r8
 inc rdx
 cmp rdx,rax
 jb 3b
 add [rbp-40],rax
 jmp 1b
5: mov rax,[rbp-32]
 lea rsi,[rbp-320]
 mov edx,1
 call file_chunk_read
 test rax,rax
 jnz file_error
6: mov rax,[rbp-32]
 call file_close
 mov rax,[rbp-24]
 mov rcx,[rbp-40]
 mov [rax],rcx
 leave
 ret


.section .text.list_new,"ax"
# rax=capacity, stable list header = [capacity, ints backing pointer]
 push rbp
 mov rbp,rsp
 sub rsp,24
 cmp rax,1048576
 ja bounds
 mov [rbp-8],rax
 mov eax,16
 call alloc
 mov [rbp-16],rax
 mov rax,[rbp-8]
 call ints_alloc
 mov qword ptr [rax],0
 mov rcx,[rbp-16]
 mov [rcx+8],rax
 mov rdx,[rbp-8]
 mov [rcx],rdx
 mov rax,rcx
 leave
 ret
.section .text.list_add,"ax"
# rax=header, rcx=value; grow backing without invalidating aliases
 push rbp
 mov rbp,rsp
 sub rsp,48
 mov [rbp-8],rax
 mov [rbp-16],rcx
 mov rdx,[rax+8]
 mov [rbp-24],rdx
 mov rcx,[rdx]
 cmp rcx,[rax]
 jb 2f
 mov [rbp-32],rcx
 mov rax,[rax]
 cmp rax,1048576
 jae bounds
 add rax,rax
 jnz 1f
 inc rax
1: cmp rax,1048576
 jbe 3f
 mov eax,1048576
3: mov [rbp-40],rax
 call ints_alloc
 mov rcx,[rbp-32]
 mov [rax],rcx
 lea rdi,[rax+8]
 mov rsi,[rbp-24]
 add rsi,8
 cld
 rep movsq
 mov rcx,[rbp-8]
 mov [rcx+8],rax
 mov rdx,[rbp-40]
 mov [rcx],rdx
 mov rdx,rax
2: mov rcx,[rdx]
 mov rax,[rbp-16]
 mov [rdx+rcx*8+8],rax
 inc rcx
 mov [rdx],rcx
 mov rax,rcx
 leave
 ret
.section .text.list_remove,"ax"
# rax=header, rcx=index; return removed integer
 mov rdx,[rax+8]
 cmp rcx,[rdx]
 jae bounds
 mov rax,[rdx+rcx*8+8]
 mov r8,[rdx]
 dec r8
 mov [rdx],r8
1: cmp rcx,r8
 jae 2f
 mov r9,[rdx+rcx*8+16]
 mov [rdx+rcx*8+8],r9
 inc rcx
 jmp 1b
2: ret
