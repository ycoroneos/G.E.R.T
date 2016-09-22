package main

import "unsafe"

const (
	GPIO_BASE = PERIPHERAL_BASE + 0x200000
	GPIO_SIZE = 0xB0
)
const (
	GPFSEL0   = 0x0
	GPFSEL1   = 0x4
	GPFSEL2   = 0x8
	GPFSEL3   = 0xC
	GPFSEL4   = 0x10
	GPFSEL5   = 0x14
	GPSET0    = 0x1C
	GPSET1    = 0x20
	GPCLR0    = 0x28
	GPCLR1    = 0x2C
	GPLEV0    = 0x34
	GPLEV1    = 0x38
	GPSED0    = 0x40
	GPSED1    = 0x44
	GPREN0    = 0x4C
	GPREN1    = 0x50
	GPFEN0    = 0x58
	GPFEN1    = 0x5C
	GPHEN0    = 0x64
	GPHEN1    = 0x68
	GPLEN0    = 0x70
	GPLEN1    = 0x74
	GPAREN0   = 0x7C
	GPAREN1   = 0x80
	GPAFEN0   = 0x88
	GPAFEN1   = 0x8C
	GPPUD     = 0x94
	GPPUDCLK0 = 0x98
	GPPUDCLK1 = 0x9C
)

//go:nosplit
func Led_init() {
	*(*uint32)(unsafe.Pointer(uintptr(GPIO_BASE + GPFSEL4))) |= FS_OUTPUT << 21
}

//go:nosplit
func Led_off() {
	*(*uint32)(unsafe.Pointer(uintptr(GPIO_BASE + GPCLR1))) = 1 << 15
}

//go:nosplit
func Led_on() {
	*(*uint32)(unsafe.Pointer(uintptr(GPIO_BASE + GPSET1))) = 1 << 15
}
