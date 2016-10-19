.equ SCTLR_ENABLE_DATA_CACHE, 0x4
.equ SCTLR_ENABLE_BRANCH_PREDICTION, 0x800
.equ SCTLR_ENABLE_INSTRUCTION_CACHE, 0x1000

// Make entry point global.
.globl _init_vectors
.align 5

// _init_vectors:
//	b   L_start
//	ldr pc, L_undefined_handler
//	ldr pc, L_svc_handler
//	ldr pc, L_prefetch_abort_handler
//	ldr pc, L_data_abort_handler
//	b   .
//	ldr pc, L_irq_handler
//	ldr pc, L_fiq_handler

_init_vectors:
	b L_start
	b L_undefined_handler
	b L_svc_handler
	b L_prefetch_abort_handler
	b L_data_abort_handler
	b .
	b L_irq_handler
	b L_fiq_handler

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

	// From the cortex A9 programming manual

	// Disable MMU and enable unaligned access
	MRC p15, 0, r1, c1, c0, 0 // Read Control Register configuration data
	BIC r1, r1, #0x1
	BIC r1, r1, #0x2
	MCR p15, 0, r1, c1, c0, 0 // Write Control Register configuration data

	// Disable L1 Caches
	MRC p15, 0, r1, c1, c0, 0 // Read Control Register configuration data
	BIC r1, r1, #(0x1 << 12)  // Disable I Cache
	BIC r1, r1, #(0x1 << 2)   // Disable D Cache
	MCR p15, 0, r1, c1, c0, 0 // Write Control Register configuration data

	// Invalidate L1 Caches
	// Invalidate Instruction cache
	MOV r1, #0
	MCR p15, 0, r1, c7, c5, 0

	// Invalidate Data cache
	// to make the code general purpose, we calculate the
	// cache size first and loop through each set + way
	MRC  p15, 1, r0, c0, c0, 0 // Read Cache Size ID
	MOVW r3, #0x1ff
	AND  r0, r3, r0, LSR #13   // r0 = no. of sets - 1
	MOV  r1, #0                // r1 = way counter way_loop

way_loop:
	MOV r3, #0 // r3 = set counter set_loop

set_loop:
	MOV r2, r1, LSL #30
	ORR r2, r3, LSL #5        // r2 = set/way cache operation format
	MCR p15, 0, r2, c7, c6, 2 // Invalidate line described by r2
	ADD r3, r3, #1            // Increment set counter
	CMP r0, r3                // Last set reached yet?
	BNE set_loop              // if not, iterate set_loop
	ADD r1, r1, #1            // else, next
	CMP r1, #4                // Last way reached yet?
	BNE way_loop              // if not, iterate way_loop

	// Invalidate TLB
	MCR p15, 0, r1, c8, c7, 0

	// Branch Prediction Enable
	MOV r1, #0
	MRC p15, 0, r1, c1, c0, 0 // Read Control Register configuration data
	ORR r1, r1, #(0x1 << 11)  // Global BP Enable bit
	MCR p15, 0, r1, c1, c0, 0 // Write Control Register configuration data

	// Enable D-side Prefetch
	MRC p15, 0, r1, c1, c0, 1 // Read Auxiliary Control Register
	ORR r1, r0, #(0x1 <<2)    // Enable D-side prefetch
	MCR p15, 0, r1, c1, c0, 1 // Write Auxiliary Control Register
	DSB
	ISB

	// Enable L1 Cache -------------------------------------------------------
	// R0 = System Control Register
	//	mrc p15, 0, r0, c1, c0, 0
	//
	//	// Enable caches and branch prediction
	//	orr r0, #SCTLR_ENABLE_BRANCH_PREDICTION
	//	orr r0, #SCTLR_ENABLE_DATA_CACHE
	//	orr r0, #SCTLR_ENABLE_INSTRUCTION_CACHE
	//
	//	// System Control Register = R0
	//	mcr p15, 0, r0, c1, c0, 0

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

	// L_reset_handler:
	//	.long reset_interrupt

L_undefined_handler:
	push {r0-r12}
	push {sp}
	push {lr}
	mov  r0, #0x1
	push {r0}
	mov  r0, sp
	b    trap

	//	push #0x1
	//	b _alltraps

	//	.long undefined_interrupt

L_svc_handler:
	push {r0-r12}
	push {sp}
	push {lr}
	mov  r0, #0x2
	push {r0}
	mov  r0, sp
	b    trap

	//	push #0x2
	//	b    _alltraps

	//	.long svc_interrupt

L_prefetch_abort_handler:
	push {r0-r12}
	push {sp}
	push {lr}
	mov  r0, #0x3
	push {r0}
	mov  r0, sp
	b    trap

	//	push #0x3
	//	b    _alltraps

	//	.long prefetch_abort_interrupt

L_data_abort_handler:
	push {r0-r12}
	push {sp}
	push {lr}
	mov  r0, #0x4
	push {r0}
	mov  r0, sp
	b    trap

	//	push #0x4
	//	b    _alltraps

	//	.long data_abort_interrupt

L_irq_handler:
	push {r0-r12}
	push {sp}
	push {lr}
	mov  r0, #0x5
	push {r0}
	mov  r0, sp
	b    trap

	//	push #0x5
	//	b    _alltraps

	//	.long irq_interrupt

L_fiq_handler:
	push {r0-r12}
	push {sp}
	push {lr}
	mov  r0, #0x6
	push {r0}
	mov  r0, sp
	b    trap

	//	push #0x6
	//	b    _alltraps

	//	.long fiq_interrupt

	// _alltraps:
	//	push {R0-R10}
	//	push lr
	//	push sp
	//	mov  r0, sp
	//	b    trap

	.globl  gobin_start
	.globl  gobin_end
	.balign 0x1000

gobin_start:
	.incbin "obj/kernel_hacked.bin"

gobin_end:
