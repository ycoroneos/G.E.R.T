// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

/*
 * Lets set up a stack
 * Turn on the MMU and caches
 * Enable VFP and NEON
 * Load a vector table into vbar
 */
#define CPSR_MODE_USER         0x10
#define CPSR_MODE_FIQ          0x11
#define CPSR_MODE_IRQ          0x12
#define CPSR_MODE_SVR          0x13
#define CPSR_MODE_ABORT        0x17
#define CPSR_MODE_UNDEFINED    0x1B
#define CPSR_MODE_SYSTEM       0x1F

#define CPSR_IRQ_INHIBIT       0x80
#define CPSR_FIQ_INHIBIT       0x40
#define CPSR_THUMB             0x20

#define SCTLR_ENABLE_DATA_CACHE        0x4
#define SCTLR_ENABLE_BRANCH_PREDICTION 0x800
#define SCTLR_ENABLE_INSTRUCTION_CACHE 0x1000
TEXT ·Getstack(SB), NOSPLIT, $-4
	SWI  $0x0
	MOVW $CPSR_MODE_SVR, R0
	ORR  $CPSR_IRQ_INHIBIT, R0
	ORR  $CPSR_FIQ_INHIBIT, R0
	WORD $0xe121f000           // msr cpsr_c, r0
	MOVW $0x8000, R13

	// Enable L1 Cache first
	WORD $0xee110f10                         // MRC p15, 0, R0, c1, c0, 0
	MOVW $SCTLR_ENABLE_BRANCH_PREDICTION, R1
	ORR  R1, R0
	MOVW $SCTLR_ENABLE_DATA_CACHE, R1
	ORR  R1, R0
	MOVW $SCTLR_ENABLE_INSTRUCTION_CACHE, R1
	ORR  R1, R0
	WORD $0xee010f10                         // mcr p15,0,r0,c1,c0,0

	// Enable VFP
	WORD $0xee110f50     // MRC p15,#0,r0,c1,c0,#2
	MOVW $0xf << 20, R1
	ORR  R1, R0
	WORD $0xee010f50     // MCR p15,#0,r0,c1,c0,#2
	MOVW $0x0, R0
	WORD $0xee070f95     // MCR p15, #0, r0, c7, c5, #4
	MOVW $0x40000000, R3
	WORD $0xeee83a10     // FMXR FPEXC, r3
	RET
