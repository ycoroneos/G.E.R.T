.equ SCTLR_ENABLE_DATA_CACHE, 0x4
.equ SCTLR_ENABLE_BRANCH_PREDICTION, 0x800
.equ SCTLR_ENABLE_INSTRUCTION_CACHE, 0x1000

// Make entry point global.
.globl _init_vectors
.align 5

_init_vectors:
	b   L_start
	ldr pc, L_undefined_handler
	ldr pc, L_svc_handler
	ldr pc, L_prefetch_abort_handler
	ldr pc, L_data_abort_handler
	b   .
	ldr pc, L_irq_handler
	ldr pc, L_fiq_handler

L_start:
	// load vbar
	adr r0, _init_vectors
	mcr p15, 0, r0, c12, c0, 0

	// Setup the stack.
	adr sp, _init_vectors

	// Clear out bss.
	ldr r4, L_bss_start
	ldr r9, L_bss_end
	mov r5, #0
	mov r6, #0
	mov r7, #0
	mov r8, #0
	b   2f

1:
	// store multiple at r4.
	stmia r4!, {r5-r8}

	// If we are still below bss_end, loop.
2:
	cmp r4, r9
	blo 1b

	// Enable L1 Cache -------------------------------------------------------
	// R0 = System Control Register
	mrc p15, 0, r0, c1, c0, 0

	// Enable caches and branch prediction
	orr r0, #SCTLR_ENABLE_BRANCH_PREDICTION
	orr r0, #SCTLR_ENABLE_DATA_CACHE
	orr r0, #SCTLR_ENABLE_INSTRUCTION_CACHE

	// System Control Register = R0
	mcr p15, 0, r0, c1, c0, 0

	// Enable VFP ------------------------------------------------------------

	// r1 = Access Control Register
	MRC p15, #0, r1, c1, c0, #2

	// enable full access for p10,11
	ORR r1, r1, #(0xf << 20)

	// access Control Register = r1
	MCR p15, #0, r1, c1, c0, #2
	MOV r1, #0

	// flush prefetch buffer because of FMXR below
	MCR p15, #0, r1, c7, c5, #4

	// and CP 10 & 11 were only just enabled
	// Enable VFP itself
	MOV r3, #0x40000000

	// FPEXC = r3
	FMXR FPEXC, r3

	// Call boot_main
	ldr r3, L_main_func
	blx r3

	// should never get here
	b .

L_bss_start:
	.long __bss_start

L_bss_end:
	.long __bss_end

L_main_func:
	.long main

L_reset_handler:
	.long reset_interrupt

L_undefined_handler:
	.long undefined_interrupt

L_svc_handler:
	.long svc_interrupt

L_prefetch_abort_handler:
	.long prefetch_abort_interrupt

L_data_abort_handler:
	.long data_abort_interrupt

L_irq_handler:
	.long irq_interrupt

L_fiq_handler:
	.long fiq_interrupt

	.globl  gobin_start
	.globl  gobin_end
	.balign 0x1000

gobin_start:
	.incbin "obj/kernel_hacked.bin"

gobin_end:
