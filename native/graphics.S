.intel_syntax noprefix
# Windows x64; shadow space and local data.
.section .text.gl_close,"ax"
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,32
 mov rcx,[rip+gl_context]
 test rcx,rcx
 jz 1f
 xor ecx,ecx
 xor edx,edx
 call qword ptr [rip+iat81]
 mov rcx,[rip+gl_context]
 call qword ptr [rip+iat82]
1: mov rdx,[rip+gl_state]
 test rdx,rdx
 jz 2f
 mov rcx,[rip+game_state]
 call qword ptr [rip+iat57]
2: xor eax,eax
 mov [rip+gl_state],rax
 mov [rip+gl_context],rax
 leave
 ret
.section .text.gl_frame,"ax"
# rax=array [x,y,color] triples, rcx=clear color; coordinates -1000..1000.
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,144
 mov [rsp+96],rax
 mov [rsp+104],rcx
 cmp rcx,16777215
 ja bounds
 mov r8,[rax]
 cmp r8,90000
 ja bounds
 mov rax,r8
 xor edx,edx
 mov ecx,9
 div rcx
 test rdx,rdx
 jnz bounds
 mov rdx,[rsp+96]
 xor ecx,ecx
1: cmp rcx,r8
 jae 2f
 mov rax,[rdx+rcx*8+8]
 cmp rax,-1000
 jl bounds
 cmp rax,1000
 jg bounds
 mov rax,[rdx+rcx*8+16]
 cmp rax,-1000
 jl bounds
 cmp rax,1000
 jg bounds
 cmp qword ptr [rdx+rcx*8+24],16777215
 ja bounds
 add rcx,3
 jmp 1b
2: mov rcx,[rip+game_state]
 test rcx,rcx
 jz 9f
 lea rdx,[rsp+112]
 call qword ptr [rip+iat58]
 test eax,eax
 jz 9f
 cmp dword ptr [rsp+120],0
 jle 8f
 cmp dword ptr [rsp+124],0
 jle 8f
 mov rax,[rip+gl_context]
 test rax,rax
 jne 4f
 mov rcx,[rip+game_state]
 call qword ptr [rip+iat56]
 test rax,rax
 jz 9f
 mov [rip+gl_state],rax
 xor eax,eax
 mov [rsp+48],rax
 mov [rsp+56],rax
 mov [rsp+64],rax
 mov [rsp+72],rax
 mov [rsp+80],rax
 mov dword ptr [rsp+48],0x00010028
 mov dword ptr [rsp+52],0x25
 mov byte ptr [rsp+57],24
 mov rcx,[rip+gl_state]
 call qword ptr [rip+iat76]
 test eax,eax
 jnz 3f
 mov rcx,[rip+gl_state]
 lea rdx,[rsp+48]
 call qword ptr [rip+iat77]
 test eax,eax
 jz 9f
 mov edx,eax
 mov rcx,[rip+gl_state]
 lea r8,[rsp+48]
 call qword ptr [rip+iat78]
 test eax,eax
 jz 9f
3: mov rcx,[rip+gl_state]
 call qword ptr [rip+iat80]
 test rax,rax
 jz 9f
 mov [rip+gl_context],rax
4: mov rcx,[rip+gl_state]
 mov rdx,[rip+gl_context]
 call qword ptr [rip+iat81]
 test eax,eax
 jz 9f
 xor ecx,ecx
 xor edx,edx
 mov r8d,[rsp+120]
 mov r9d,[rsp+124]
 call qword ptr [rip+iat83]
 mov eax,0x437f0000
 movd xmm4,eax
 mov eax,[rsp+104]
 shr eax,16
 and eax,255
 cvtsi2ss xmm0,eax
 divss xmm0,xmm4
 mov eax,[rsp+104]
 shr eax,8
 and eax,255
 cvtsi2ss xmm1,eax
 divss xmm1,xmm4
 movzx eax,byte ptr [rsp+104]
 cvtsi2ss xmm2,eax
 divss xmm2,xmm4
 mov eax,0x3f800000
 movd xmm3,eax
 call qword ptr [rip+iat84]
 mov ecx,0x4000
 call qword ptr [rip+iat85]
 mov ecx,4
 call qword ptr [rip+iat86]
 mov qword ptr [rsp+128],0
5: mov rax,[rsp+96]
 mov rdx,[rsp+128]
 cmp rdx,[rax]
 jae 6f
 lea rax,[rax+rdx*8+8]
 mov ecx,[rax+16]
 mov edx,ecx
 mov r8d,ecx
 shr ecx,16
 shr edx,8
 and ecx,255
 and edx,255
 and r8d,255
 call qword ptr [rip+iat87]
 mov rax,[rsp+96]
 mov rdx,[rsp+128]
 lea rax,[rax+rdx*8+8]
 cvtsi2ss xmm0,qword ptr [rax]
 cvtsi2ss xmm1,qword ptr [rax+8]
 mov eax,0x447a0000
 movd xmm4,eax
 divss xmm0,xmm4
 divss xmm1,xmm4
 call qword ptr [rip+iat88]
 add qword ptr [rsp+128],3
 jmp 5b
6: call qword ptr [rip+iat89]
 mov rcx,[rip+gl_state]
 call qword ptr [rip+iat79]
 mov [rsp+136],eax
 xor ecx,ecx
 xor edx,edx
 call qword ptr [rip+iat81]
 cmp dword ptr [rsp+136],0
 je 9f
8: mov eax,1
 leave
 ret
9: call gl_close
 xor eax,eax
 leave
 ret
.section .text.pad_read,"ax"
# Controller 0..3, field 0..8; disconnected returns -1.
 cmp rax,3
 ja bounds
 cmp rcx,8
 ja bounds
 push rbp
 mov rbp,rsp
 and rsp,-16
 sub rsp,64
 mov [rsp+48],rcx
 mov ecx,eax
 lea rdx,[rsp+32]
 call qword ptr [rip+iat75]
 test eax,eax
 jnz 9f
 mov rcx,[rsp+48]
 test rcx,rcx
 jz 8f
 cmp rcx,1
 je 1f
 cmp rcx,2
 je 2f
 cmp rcx,3
 je 3f
 cmp rcx,4
 je 4f
 cmp rcx,5
 je 5f
 cmp rcx,6
 je 6f
 cmp rcx,7
 je 7f
 mov eax,[rsp+32]
 leave
 ret
1: movzx eax,word ptr [rsp+36]
 leave
 ret
2: movzx eax,byte ptr [rsp+38]
 leave
 ret
3: movzx eax,byte ptr [rsp+39]
 leave
 ret
4: movsx rax,word ptr [rsp+40]
 leave
 ret
5: movsx rax,word ptr [rsp+42]
 leave
 ret
6: movsx rax,word ptr [rsp+44]
 leave
 ret
7: movsx rax,word ptr [rsp+46]
 leave
 ret
8: mov eax,1
 leave
 ret
9: mov rax,-1
 leave
 ret
