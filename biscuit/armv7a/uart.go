package main

import "unsafe"

const UART1_UTXD uint32 = 0x02020040
const UART1_UTS uint32 = 0x020200B4

//go:nosplit
func uart_putc(c byte) {
	*(*byte)(unsafe.Pointer(uintptr(UART1_UTXD))) = c
	for (*(*uint32)(unsafe.Pointer(uintptr(UART1_UTS))) & (0x1 << 6)) <= 0 {
	}
	//
	//	uart := (*uint32)(unsafe.pointer(uintptr(uart1_utxd)))
	//	utrstat1 := (*uint32)(unsafe.pointer(uintptr(uart1_uts)))
	//	*uart = uint32(c)
	//	for (*utrstat1 & (0x1 << 6)) <= 0 {
	//	}
}

//go:nosplit
func uart_print(b []byte) {
	var i = 0
	for i = 0; i < len(b); i++ {
		uart_putc(b[i])
	}
}

//const (
//	UART_BASE = PERIPHERAL_BASE + 0x201000
//	UART_SIZE = 0x8C
//)
//const (
//	DR     = 0x0
//	RSRECR = 0x4
//	FR     = 0x8
//	ILPR   = 0xC
//	IBRD   = 0x10
//	FBRD   = 0x14
//	LCRH   = 0x18
//	CR     = 0x1C
//	IFLS   = 0x20
//	IMSC   = 0x24
//	RIS    = 0x28
//	MIS    = 0x2C
//	ICR    = 0x30
//	DMACR  = 0x34
//	ITCR   = 0x38
//	ITIP   = 0x3C
//	ITOP   = 0x40
//	TDR    = 0x44
//)
//
////go:nosplit
//func uart_init() {
//	//first disable uart
//	cr := (*uint32)(unsafe.Pointer(uintptr(CR + UART_BASE)))
//	*cr = 0
//
//	//set no pullup/pulldown for gpios
//	gppud := (*uint32)(unsafe.Pointer(uintptr(GPPUD + GPIO_BASE)))
//	*gppud = 0
//	delay(150)
//
//	//enable it on gpios 14,15
//	gppudclk0 := (*uint32)(unsafe.Pointer(uintptr(GPPUDCLK0 + GPIO_BASE)))
//	*gppudclk0 = 1<<14 | 1<<15
//
//	//commit effects
//	*gppud = 0
//
//	//clear uart interrupts
//	icr := (*uint32)(unsafe.Pointer(uintptr(ICR + UART_BASE)))
//	*icr = 0x7FF
//
//	//Set integer & fractional part of baud rate.
//	// Divider = UART_CLOCK/(16 * Baud)
//	// Fraction part register = (Fractional part * 64) +   0.5
//	// UART_CLOCK = 3000000; Baud = 115200.
//
//	// Divider = 3000000 /    (16 * 115200) = 1.627 = ~1.
//	// Fractional part register = (.627 * 64)   + 0.5 = 40.6 = ~40.
//	ibrd := (*uint32)(unsafe.Pointer(uintptr(IBRD + UART_BASE)))
//	fbrd := (*uint32)(unsafe.Pointer(uintptr(FBRD + UART_BASE)))
//	*ibrd = 1
//	*fbrd = 40
//
//	// Enable FIFO & 8 bit data transmissio (1 stop bit, no parity).
//	lcrh := (*uint32)(unsafe.Pointer(uintptr(LCRH + UART_BASE)))
//	*lcrh = 1<<4 | 1<<5 | 1<<6
//
//	// Mask all interrupts.
//	imsc := (*uint32)(unsafe.Pointer(uintptr(IMSC + UART_BASE)))
//	*imsc = (1 << 1) | (1 << 4) | (1 << 5) | (1 << 6) | (1 << 7) | (1 << 8) | (1 << 9) | (1 << 10)
//
//	//enable uart
//	*cr = (1 << 0) | (1 << 8) | (1 << 9)
//}
//
////go:nosplit
//func uart_putc(data byte) {
//	fr := (*uint32)(unsafe.Pointer(uintptr(FR + UART_BASE)))
//	dr := (*uint32)(unsafe.Pointer(uintptr(DR + UART_BASE)))
//	for (*fr & (1 << 4)) > 0 {
//	}
//	*dr = uint32(data)
//}
//
////go:nosplit
//func uart_getc() byte {
//	fr := (*uint32)(unsafe.Pointer(uintptr(FR + UART_BASE)))
//	dr := (*uint32)(unsafe.Pointer(uintptr(DR + UART_BASE)))
//	for (*fr & (1 << 5)) > 0 {
//	}
//	return byte(*dr)
//}
//
////go:nosplit
//func uart_print(data uint32, len int32) int32 {
//	for i := uint32(0); i < uint32(len); i++ {
//		uart_putc(*(*byte)(unsafe.Pointer(uintptr(data + i))))
//	}
//	return len
//}
