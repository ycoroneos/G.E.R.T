package main

import "unsafe"
import "fmt"

/*
209_8000 GPT Control Register (GPT_CR) 32 R/W 0000_0000h 30.6.1/1499
209_8004 GPT Prescaler Register (GPT_PR) 32 R/W 0000_0000h 30.6.2/1503
209_8008 GPT Status Register (GPT_SR) 32 R/W 0000_0000h 30.6.3/1504
209_800C GPT Interrupt Register (GPT_IR) 32 R/W 0000_0000h 30.6.4/1505
209_8010 GPT Output Compare Register 1 (GPT_OCR1) 32 R/W FFFF_FFFFh 30.6.5/1506
209_8014 GPT Output Compare Register 2 (GPT_OCR2) 32 R/W FFFF_FFFFh 30.6.6/1507
209_8018 GPT Output Compare Register 3 (GPT_OCR3) 32 R/W FFFF_FFFFh 30.6.7/1507
209_801C GPT Input Capture Register 1 (GPT_ICR1) 32 R 0000_0000h 30.6.8/1508
209_8020 GPT Input Capture Register 2 (GPT_ICR2) 32 R 0000_0000h 30.6.9/1508
209_8024 GPT Counter Register (GPT_CNT) 32 R 0000_0000h 30.6.10/1509
*/

type GPT struct {
	CR   uint32
	PR   uint32
	SR   uint32
	IR   uint32
	OCR1 uint32
	OCR2 uint32
	OCR3 uint32
	ICR1 uint32
	ICR2 uint32
	CNT  uint32
}

var gpt *GPT = (*GPT)(unsafe.Pointer(uintptr(0x2098000)))

func startGPT() bool {
	fmt.Printf("GPT lives at %x\r\n", gpt)
	gpt.CR = 0
	gpt.CR |= (0x7 << 6) | (0x1 << 5)
	gpt.PR = 0xFFF
	gpt.IR = 0x1
	gpt.OCR1 = 0x1FFF
	gpt.CR |= 0x1
	return true
}
