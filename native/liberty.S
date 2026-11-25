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
 mov rcx,rax
 xor eax,eax
1: cmp byte ptr [rcx+rax],0
 je 2f
 inc rax
 jmp 1b
2: ret
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
 # 8 MiB arena allocated lazily. Immutable text remains valid until exit.
 push rbp
 mov rbp,rsp
 sub rsp,8
 mov [rbp-8],rax
 lea rdx,[rip+arena_state]
 cmp qword ptr [rdx],0
 jne 1f
 mov eax,8388608
 call osalloc
 lea rdx,[rip+arena_state]
 mov [rdx],rax
 mov qword ptr [rdx+8],8388608
1: mov rcx,[rbp-8]
 cmp rcx,[rdx+8]
 ja text_error
 mov rax,[rdx]
 add [rdx],rcx
 sub [rdx+8],rcx
 leave
 ret
.section .text.osalloc_linux,"ax"
 mov rsi,rax
 xor edi,edi
 mov edx,3
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
 mov r8d,0x3000
 mov r9d,4
 call qword ptr [rip+iat4]
 test rax,rax
 jz text_error
 leave
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
 mov eax,1
 jmp exit
1: .ascii "TRIPOL runtime error: invalid input, text too long or allocation failed\n"
2:
.section .text.float_error,"ax"
 lea rsi,[rip+1f]
 mov edx,2f-1f
 call writeerr
 mov eax,1
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
 mov eax,1
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
