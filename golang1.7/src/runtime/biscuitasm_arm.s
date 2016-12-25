#include "go_asm.h"
#include "go_tls.h"
#include "funcdata.h"
#include "textflag.h"
#define TO	R8
#define TOE	R11
#define N	R12
#define TMP	R12 // N and TMP don't overlap

TEXT runtime·PutR0(SB), NOSPLIT, $0-4
	MOVW ret+4(FP), R0
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

TEXT runtime·ReplayTrapframe(SB), NOSPLIT, $0
	MOVW runtime·curthread(SB), R4
	MOVW 0(R4), LR
	MOVW 4(R4), R13
	MOVW 8(R4), R11
	MOVW 12(R4), R0
	MOVW 16(R4), R1
	MOVW 20(R4), R2
	MOVW 24(R4), R3
	MOVW 28(R4), g
	RET

TEXT runtime·Threadschedule(SB), NOSPLIT, $0
	MOVW R13, R5
	ADD  $4, R5, R5                // undo the push {lr}
	MOVW runtime·curthread(SB), R4
	MOVW LR, 0(R4)

	//	MOVW R13, 4(R4)
	MOVW R5, 4(R4)
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

TEXT runtime·loadvbar(SB), NOSPLIT, $0-4
	MOVW vbar_addr+0(FP), R0
	MOVW LR, R6

	// load vbar from R0
	WORD $0xee0c0f10

	B (R6)

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
