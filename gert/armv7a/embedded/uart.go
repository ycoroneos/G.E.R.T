// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedded

const (
	READ_MASK = 0xFF
	RRDY      = 1 << 9
)

type UART_regs struct {
	urxd  uint32
	_     [15]uint32 //thanks freescale
	utxd  uint32
	_     [15]uint32 //thanks freescale
	ucr1  uint32
	ucr2  uint32
	ucr3  uint32
	ucr4  uint32
	ufcr  uint32
	usr1  uint32
	usr2  uint32
	uesc  uint32
	utim  uint32
	ubir  uint32
	ubmr  uint32
	ubrc  uint32
	onems uint32
	uts   uint32
	umcr  uint32
}

type UART struct {
	regs *UART_regs
}

//blocking read
func (u *UART) getchar() byte {
	for (u.regs.usr1 & RRDY) == 0 {
	}
	return byte(u.regs.urxd & READ_MASK)
}

func (u *UART) Read(n int) []byte {
	output := make([]byte, n)
	for i := 0; i < n; i++ {
		output[i] = u.getchar()
	}
	return output
}
