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
