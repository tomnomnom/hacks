global _start

section .text
_start:

	mov	rax, 0
	mov	rdi, 0
	mov rsi, buf
	mov rdx, length
	syscall
	push rax

    mov rax, 1
    mov rdi, 1
    mov rsi, buf
	pop	rdx
    ;mov rdx, length
    syscall

    mov rax, 60
    mov rdi, 0
    syscall

section .data
    ;msg: db 'Hello, ASM!',0x0A
    length: dd 16

section .bss
	buf: resb 16
