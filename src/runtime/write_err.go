// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !android

package runtime

import "unsafe"

const UART1_UTXD uint32 = 0x02020040
const UART1_UTS uint32 = 0x020200B4

//go:nosplit
func uart_putc(c byte) {
	*(*byte)(unsafe.Pointer(uintptr(UART1_UTXD))) = c
	for (*(*uint32)(unsafe.Pointer(uintptr(UART1_UTS))) & (0x1 << 6)) <= 0 {
	}
}

func writeErr(b []byte) {
	var i = 0
	for i = 0; i < len(b); i++ {
		uart_putc(b[i])
	}
	//write(2, unsafe.Pointer(&b[0]), int32(len(b)))
}

//go:nowritebarrierrec
//go:nosplit
func writeUnsafe(b []byte) {
	var i = 0
	for i = 0; i < len(b); i++ {
		uart_putc(b[i])
	}
}
