// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !android

package runtime

import "unsafe"

func writeErr(b []byte) {
	write(2, unsafe.Pointer(&b[0]), int32(len(b)))
	//	if Armhackmode > 0 {
	//
	//		write_uart(b)
	//
	//	} else {
	//		write(2, unsafe.Pointer(&b[0]), int32(len(b)))
	//	}
}

////printing
const UART1_UTXD uint32 = 0x02020040
const UART1_UTS uint32 = 0x020200B4

//go:nosplit
func uart_putc(c byte) {
	*(*byte)(unsafe.Pointer(uintptr(UART1_UTXD))) = c
	for (*(*uint32)(unsafe.Pointer(uintptr(UART1_UTS))) & (0x1 << 6)) <= 0 {
	}
}

//go:nosplit
//go:nobarrierec
func write_uart(b []byte) {
	var i = 0
	for i = 0; i < len(b); i++ {
		uart_putc(b[i])
	}
}

//go:nosplit
//go:nobarrierec
func write_uart_unsafe(buf uintptr, count uint32) uint32 {
	for i := uint32(0); i < count; i++ {
		c := ((*byte)(unsafe.Pointer(buf + uintptr(i))))
		uart_putc(*c)
	}
	return count
}
