#include "go_asm.h"
#include "go_tls.h"
#include "funcdata.h"
#include "textflag.h"
#define TO	R8
#define TOE	R11
#define N	R12
#define TMP	R12 // N and TMP don't overlap
#define SCTLR_ENABLE_DATA_CACHE $0x4
#define SCTLR_ENABLE_BRANCH_PREDICTION $0x800
#define SCTLR_ENABLE_INSTRUCTION_CACHE $0x1000
#define ACTLR_SMP $0x40
#define ACTLR_L1_PREFETCH $0x4
#define ACTLR_L2_PREFETCH $0x2
#define ACTLR_FW $0x1

TEXT runtime·PutR0(SB), NOSPLIT, $0-4
	MOVW ret+4(FP), R0
	RET

TEXT runtime·PutSP(SB), NOSPLIT, $0-4
	MOVW ret+4(FP), R13
	RET

TEXT runtime·PutR2(SB), NOSPLIT, $0-4
	MOVW ret+4(FP), R2
	RET

TEXT runtime·RR0(SB), NOSPLIT, $0
	MOVW R0, ret+0(FP)
	RET

TEXT runtime·RR1(SB), NOSPLIT, $0
	MOVW R1, ret+0(FP)
	RET

TEXT runtime·RR2(SB), NOSPLIT, $0
	MOVW R2, ret+0(FP)
	RET

TEXT runtime·RR3(SB), NOSPLIT, $0
	MOVW R3, ret+0(FP)
	RET

TEXT runtime·RR4(SB), NOSPLIT, $0
	MOVW R4, ret+0(FP)
	RET

TEXT runtime·RR5(SB), NOSPLIT, $0
	MOVW R5, ret+0(FP)
	RET

TEXT runtime·RR6(SB), NOSPLIT, $0
	MOVW R6, ret+0(FP)
	RET

TEXT runtime·RSP(SB), NOSPLIT, $0
	MOVW R13, ret+0(FP)
	RET

TEXT runtime·RR7(SB), NOSPLIT, $0
	MOVW R7, ret+0(FP)
	RET

TEXT runtime·DMB(SB), NOSPLIT, $0
	WORD $0xf57ff05e // DMB ST
	RET

TEXT runtime·RecordTrapframe(SB), NOSPLIT, $0
	MOVW runtime·curthread(SB), R4
	MOVW LR, 0(R4)
	MOVW R13, 4(R4)
	MOVW R11, 8(R4)
	MOVW R0, 12(R4)
	MOVW R1, 16(R4)
	MOVW R2, 20(R4)
	MOVW R3, 24(R4)
	MOVW g, 28(R4)
	RET

// TEXT runtime·ReplayTrapframe(SB), NOSPLIT, $0
//	MOVW runtime·curthread(SB), R4
//	MOVW 0(R4), LR
//	MOVW 4(R4), R13
//	MOVW 8(R4), R11
//	MOVW 12(R4), R0
//	MOVW 16(R4), R1
//	MOVW 20(R4), R2
//	MOVW 24(R4), R3
//	MOVW 28(R4), g
//	RET

TEXT runtime·ReplayTrapframe(SB), NOSPLIT, $0-4
	MOVW thread_ptr+0(FP), R4
	MOVW 0(R4), LR
	MOVW 4(R4), R13
	MOVW 8(R4), R11
	MOVW 12(R4), R0
	MOVW 16(R4), R1
	MOVW 20(R4), R2
	MOVW 24(R4), R3
	MOVW 28(R4), g

	RET

TEXT runtime·disable_interrupts(SB), NOSPLIT, $0
	// disable interrupts while in thread scheduler
	WORD $0xe321f0d3 // msr	CPSR_c, #211	; 0xd3
	RET

TEXT runtime·enable_interrupts(SB), NOSPLIT, $0
	// enable interrupts when running a thread
	WORD $0xe321f053 // msr	CPSR_c, #83	; 0x53
	RET

TEXT runtime·Threadschedule(SB), NOSPLIT, $0

	MOVW R13, R5
	ADD  $4, R5, R5 // undo the push {lr}

	// get the cpunum
	WORD $0xee106fb0               // mrc	15, 0, r6, cr0, cr0, {5}
	AND  $3, R6                    // get rid of everything except cpuid
	MOVW $4, R4
	MUL  R4, R6
	MOVW runtime·curthread(SB), R4
	MOVW R11, R4                   // the golang assember likes to dereference pointers without telling me
	ADD  R6, R4                    // offset into curthread list
	MOVW (R4), R4                  // dereference
	MOVW LR, 0(R4)

	//	MOVW R13, 4(R4)
	MOVW R5, 4(R4)                   // save the modified stack pointer
	MOVW R11, 8(R4)
	MOVW R0, 12(R4)
	MOVW R1, 16(R4)
	MOVW R2, 20(R4)
	MOVW R3, 24(R4)
	MOVW g, 28(R4)
	CALL runtime·thread_schedule(SB)

	// this should never happen
	RET

TEXT runtime·SWIcall(SB), NOSPLIT, $0
	MOVW LR, R6
	SWI  $0x0
	B    (R6)

TEXT runtime·Readvbar(SB), NOSPLIT, $0
	WORD $0xee1c0f10   // mrc	15, 0, r0, cr12, cr0, {0}
	MOVW R0, ret+0(FP)
	RET

TEXT runtime·loadvbar(SB), NOSPLIT, $0-4
	MOVW vbar_addr+0(FP), R0

	// load vbar from R0
	WORD $0xee0c0f10
	RET

TEXT runtime·loadttbr0(SB), NOSPLIT, $0-4
	MOVW l1base+0(FP), R0
	MOVW runtime·l1_table(SB), R0
	MOVW LR, R6

	// instruction barrier
	WORD $0xf57ff06f

	// data barrier
	WORD $0xf57ff04f

	// invalidate instruction tlb
	WORD $0xee080f15

	// invalidate data tlb
	WORD $0xee080f16

	// invalidate whole tlb
	WORD $0xee080f17

	// put r0 into ttbr0
	WORD $0xee020f10

	// clear TTBCR
	MOVW $0x0, R0
	WORD $0xee020f50

	// all domain access
	MOVW $0x3, R0
	WORD $0xee030f10

	// read mmu config into R0
	WORD $0xee110f10

	// enable MMU
	ORR $0x1, R0

	// instruction barrier
	WORD $0xf57ff06f

	// data barrier
	WORD $0xf57ff04f

	// put R0 into mmu config
	WORD $0xee010f10
	B    (R6)

TEXT runtime·invallpages(SB), NOSPLIT, $0
	MOVW LR, R6

	// instruction barrier
	WORD $0xf57ff06f

	// data barrier
	WORD $0xf57ff04f

	// invalidate instruction tlb
	WORD $0xee080f15

	// invalidate data tlb
	WORD $0xee080f16

	// invalidate whole tlb
	WORD $0xee080f17
	B    (R6)

TEXT runtime·Runtime_main(SB), NOSPLIT, $0
	MOVW $runtime·rt0_go(SB), R4
	B    (R4)

TEXT runtime·memclrbytes(SB), NOSPLIT, $0-8
	MOVW ptr+0(FP), TO
	MOVW n+4(FP), N
	MOVW $0, R0
	MOVW LR, R6

	ADD N, TO, TOE // to end pointer
	CMP TOE, TO
	BLT _zero

	// return if nothing to do
	B (R6)

_zero:
	MOVB R0, (TO)
	ADD  $0x1, TO, TO
	CMP  TOE, TO
	BLT  _zero

	// return
	B (R6)

TEXT runtime·cleardcache(SB), NOSPLIT, $0
	WORD $0xee300f10 // mrc	15, 1, r0, cr0, cr0, {0}
	WORD $0xe30031ff // movw	r3, #511	; 0x1ff
	WORD $0xe00306a0 // and	r0, r3, r0, lsr #13
	WORD $0xe3a01000 // mov	r1, #0
	WORD $0xe3a03000 // mov	r3, #0
	WORD $0xe1a02f01 // lsl	r2, r1, #30
	WORD $0xe1822283 // orr	r2, r2, r3, lsl #5
	WORD $0xee072f56 // mcr	15, 0, r2, cr7, cr6, {2}
	WORD $0xe2833001 // add	r3, r3, #1
	WORD $0xe1500003 // cmp	r0, r3
	WORD $0x1afffff9 // bne	14 <set_loop>
	WORD $0xe2811001 // add	r1, r1, #1
	WORD $0xe3510004 // cmp	r1, #4
	WORD $0x1afffff5 // bne	10 <way_loop>
	RET

TEXT runtime·cpunum(SB), NOSPLIT, $0
	// first read cpu id into r0
	WORD $0xee100fb0   // mrc	15, 0, r0, cr0, cr0, {5}
	AND  $3, R0        // get rid of everything except cpuid
	MOVW R0, ret+0(FP)
	RET

TEXT runtime·scu_enable(SB), NOSPLIT, $0
	// turn on the scu
	MOVW runtime·scubase(SB), R1
	MOVW (R1), R0                // read scu config register

	ORR  $1, R0   // enable bit
	MOVW R0, (R1) // write it in

	// data cache was cleared and disabled in the bootloader
	// enable it here
	WORD $0xee110f10 // mrc	15, 0, r0, cr1, cr0, {0}

	ORR  SCTLR_ENABLE_DATA_CACHE, R0
	ORR  SCTLR_ENABLE_BRANCH_PREDICTION, R0
	ORR  SCTLR_ENABLE_INSTRUCTION_CACHE, R0
	WORD $0xee010f10                        // mcr	15, 0, r0, cr1, cr0, {0}

	// read actlr into r0
	WORD $0xee110f30 // mrc	15, 0, r0, cr1, cr0, {1}

	ORR  ACTLR_SMP, R0
	ORR  ACTLR_L1_PREFETCH, R0
	ORR  ACTLR_L2_PREFETCH, R0
	ORR  ACTLR_FW, R0
	WORD $0xee010f30           // mcr	15, 0, r0, cr1, cr0, {1}

	RET

// sets up isr stack on cpu0
TEXT runtime·isr_setup(SB), NOSPLIT, $0
	MOVW runtime·isr_stack(SB), R0

	// set up isr stack
	WORD $0xe321f0d2 // msr	CPSR_c, #210	; 0xd2
	MOVW R0, R13
	WORD $0xe321f0d3 // msr CPSR_c, #211; 0xd3

	// set up abort stack
	WORD $0xe321f0d7 // msr	CPSR_c, #215	; 0xd7
	MOVW R0, R13
	WORD $0xe321f0d3 // msr	CPSR_c, #211	; 0xd3
	RET

TEXT runtime·boot_any(SB), NOSPLIT, $0
	MOVW $0x02020040, R0
	MOVW $64, R1
	MOVW R1, (R0)

	// first read cpu id into r0
	WORD $0xee100fb0                 // mrc	15, 0, r0, cr0, cr0, {5}
	AND  $3, R0                      // get rid of everything except cpuid
	MOVW R0, R2
	MOVW $8, R1
	MUL  R1, R2
	SUB  $8, R2
	MOVW runtime·cpu1bootarg(SB), R1
	MOVW R11, R1

	ADD R1, R2

	MOVW (R2), R2 // now r2 contains *sp
	MOVW (R2), R2 // now r2 contains sp
	MOVW R2, R13
	MOVW R2, R0

	// set up isr stack
	WORD $0xe321f0d2 // msr	CPSR_c, #210	; 0xd2
	MOVW R0, R13
	WORD $0xe321f0d3 // msr CPSR_c, #211; 0xd3

	// set up abort stack
	WORD $0xe321f0d7 // msr	CPSR_c, #215	; 0xd7
	MOVW R0, R13
	WORD $0xe321f0d3 // msr	CPSR_c, #211	; 0xd3

	// set up undefined stack
	WORD $0xe321f0db // msr	CPSR_c, #219	; 0xdb
	MOVW R0, R13
	WORD $0xe321f0d3 // msr	CPSR_c, #211	; 0xd3

	CALL runtime·cleardcache(SB)
	MOVW $0x02020040, R0
	MOVW $65, R1
	MOVW R1, (R0)

	// enable data cache
	WORD $0xee110f10 // mrc	15, 0, r0, cr1, cr0, {0}

	ORR  SCTLR_ENABLE_DATA_CACHE, R0
	ORR  SCTLR_ENABLE_BRANCH_PREDICTION, R0
	ORR  SCTLR_ENABLE_INSTRUCTION_CACHE, R0
	WORD $0xee010f10                        // mcr	15, 0, r0, cr1, cr0, {0}

	// set SMP bit
	// read actlr into r0
	WORD $0xee110f30 // mrc	15, 0, r0, cr1, cr0, {1}

	ORR  ACTLR_SMP, R0
	ORR  ACTLR_L1_PREFETCH, R0
	ORR  ACTLR_L2_PREFETCH, R0
	ORR  ACTLR_FW, R0
	WORD $0xee010f30           // mcr	15, 0, r0, cr1, cr0, {1}

	// enable floating point
	WORD $0xee111f50 // mrc	15, 0, r1, cr1, cr0, {2}
	WORD $0xe381160f // orr	r1, r1, #15728640	; 0xf00000
	WORD $0xee011f50 // mcr	15, 0, r1, cr1, cr0, {2}

	MOVW $0, R1
	WORD $0xee071f95     // mcr	15, 0, r1, cr7, cr5, {4}
	MOVW $0x40000000, R3
	WORD $0xeee83a10     // vmsr	fpexc, r3

	MOVW $0x02020040, R0
	MOVW $66, R1
	MOVW R1, (R0)

	// load vectors
	//	CALL runtime·loadvbar(SB)

	// load the page tables
	CALL runtime·loadttbr0(SB)

	MOVW $0x02020040, R0
	MOVW $67, R1
	MOVW R1, (R0)
	MOVW $0x02020040, R0
	MOVW $68, R1
	MOVW R1, (R0)

	// enter holding pen
	CALL runtime·mp_pen(SB)

wait:
	B wait
	RET

TEXT runtime·getentry(SB), NOSPLIT, $0
	MOVW runtime·boot_any(SB), R2
	MOVW R11, R2
	MOVW R2, ret+0(FP)
	RET

TEXT runtime·catch(SB), NOSPLIT, $0
	WORD $0xe24ee004 // sub	lr, lr, #4
	WORD $0xe92d5fff // push	{r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, lr}

	// MOVW $0x02020040, R0
	// MOVW $64, R1
	MOVW R14, R0
	MOVW $runtime·cpucatch(SB), R11
	BL   (R11)

	// MOVW $0x02020040, R0
	// MOVW $65, R1
	WORD $0xe8fd9fff // ldm	sp!, {r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, pc}^
	RET              // wont get here

TEXT runtime·getcatch(SB), NOSPLIT, $0
	MOVW runtime·catch(SB), R2
	MOVW R11, R2
	MOVW R2, ret+0(FP)
	RET

TEXT runtime·abort_int(SB), NOSPLIT, $0
	WORD $0xe24ee008                // sub	lr, lr, #8
	WORD $0xe92d5fff                // push	{r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, lr}
	MOVW R14, R0
	MOVW $runtime·cpuabort(SB), R11
	BL   (R11)
	WORD $0xe8fd9fff                // ldm	sp!, {r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, pc}^
	RET                             // wont get here

TEXT runtime·getabort_int(SB), NOSPLIT, $0
	MOVW runtime·abort_int(SB), R2
	MOVW R11, R2
	MOVW R2, ret+0(FP)
	RET

TEXT runtime·pref_abort(SB), NOSPLIT, $0
	WORD $0xe24ee004                    // sub	lr, lr, #4
	WORD $0xe92d5fff                    // push	{r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, lr}
	MOVW R14, R0
	MOVW $runtime·cpuprefabort(SB), R11
	BL   (R11)
	WORD $0xe8fd9fff                    // ldm	sp!, {r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, pc}^
	RET                                 // wont get here

TEXT runtime·getpref_abort(SB), NOSPLIT, $0
	MOVW runtime·pref_abort(SB), R2
	MOVW R11, R2
	MOVW R2, ret+0(FP)
	RET

TEXT runtime·undefined(SB), NOSPLIT, $0
	WORD $0xe24ee004                    // sub	lr, lr, #4
	WORD $0xe92d5fff                    // push	{r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, lr}
	MOVW R14, R0
	MOVW $runtime·cpuundefined(SB), R11
	BL   (R11)
	WORD $0xe8fd9fff                    // ldm	sp!, {r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, sl, fp, ip, pc}^
	RET                                 // wont get here

TEXT runtime·getundefined(SB), NOSPLIT, $0
	MOVW runtime·undefined(SB), R2
	MOVW R11, R2
	MOVW R2, ret+0(FP)
	RET

TEXT runtime·EnableIRQ(SB), NOSPLIT, $0
	// just flips the I bit
	WORD $0xe321f053 // msr CPSR_c, #83; 0x53
	RET

TEXT runtime·DisableIRQ(SB), NOSPLIT, $0
	// just sets the I bit
	WORD $0xe321f0d3 // msr CPSR_c, #211; 0xd3
	RET

TEXT runtime·Getmpcorebase(SB), NOSPLIT, $0
	WORD $0xee9f0f10   // mrc 15, 4, r0, cr15, cr0, {0}
	MOVW R0, ret+0(FP)
	RET

