// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedded

import (
	"fmt"
)

type SPI_regs struct {
	rxdata  uint32
	txdata  uint32
	control uint32
	config  uint32
	intr    uint32
	dma     uint32
	status  uint32
	period  uint32
	test    uint32
	_       [7]uint32
	msgdata uint32
}

//either mosi, miso, sclk, or cs
type SPI_pin struct {
	name   string
	alt    uint8
	muxctl *uint32
	padctl *uint32
}

type SPI_periph struct {
	mosi       SPI_pin
	miso       SPI_pin
	sclk       SPI_pin
	cs         []SPI_pin
	regs       *SPI_regs
	mode       uint32
	frequency  uint32
	datalength uint32
}

//SPI has 3 types of modes which affect the polarity of the clock and the resting state of the data signals
//data length is how many bits each SPI frame contains. 7,8,16 are common amounts

//the imx6 can support up to 2^12 bits in a single frame!
func (spi *SPI_periph) Begin(mode, freq, datalength, channel uint32) {
	//put the gpio pins on push/pull mode with their appropriate alternate functions
	*spi.mosi.muxctl = makeGPIOmuxconfig(spi.mosi.alt)
	*spi.mosi.padctl = makeGPIOpadconfig(1, PULLDOWN_100K, 1, 1, 0, SPEED_FAST, DRIVE_260R, SLEW_FAST)

	*spi.miso.muxctl = makeGPIOmuxconfig(spi.miso.alt)
	*spi.miso.padctl = makeGPIOpadconfig(1, PULLDOWN_100K, 1, 1, 0, SPEED_FAST, DRIVE_260R, SLEW_FAST)

	*spi.sclk.muxctl = makeGPIOmuxconfig(spi.sclk.alt)
	*spi.sclk.padctl = makeGPIOpadconfig(1, PULLDOWN_100K, 1, 1, 0, SPEED_FAST, DRIVE_260R, SLEW_FAST)

	for i := 0; i < len(spi.cs); i++ {
		*spi.cs[i].muxctl = makeGPIOmuxconfig(spi.cs[i].alt)
		*spi.cs[i].padctl = makeGPIOpadconfig(1, PULLDOWN_100K, 1, 1, 0, SPEED_FAST, DRIVE_260R, SLEW_FAST)
	}

	//ungate the module clock
	//Put 0x3F into CCM_CCGR1
	//CCM_CCGR1 := ((*uint32)(unsafe.Pointer(uintptr(0x20C406C))))
	//*CCM_CCGR1 |= 0x3F

	//configure the SPI registers for the freq, mode, and datalength
	spi.regs.control = 0
	spi.regs.control |= datalength << 20
	if channel > 1 {
		fmt.Printf("bad spi channel, abort\n")
		return
	}
	spi.regs.control |= channel << 18

	//use a post divider
	spi.regs.control |= freq << 8

	//everyone is in master mode
	spi.regs.control |= 0xF << 4

	spi.regs.control |= 1 << 3

	spi.regs.config = 0
	if mode == 0 {
	} else if mode == 3 {
		spi.regs.config |= 0xFF
	} else {
		fmt.Printf("unsupported mode %d. abort\n", mode)
	}
	spi.regs.config |= 0xF << 12

	//toggle cs on every burst
	spi.regs.config |= 0xF << 8

	spi.regs.intr = 0

	spi.regs.control |= 1

	spi.regs.status = 0xff

	spi.mode = mode
	spi.frequency = freq
	spi.datalength = datalength
}

//assumes datalength < 32bits
func (spi *SPI_periph) Send(data uint32) {
	//mask := uint32(1<<spi.datalength) - uint32(1)
	//data = data & mask

	//wait for tx fifo to have space
	//for spi.regs.status&0x2 != 0 {
	//}

	spi.regs.txdata = data
}

//assumes datalength < 32 bits
func (spi *SPI_periph) Exchange(data uint32) uint32 {
	mask := uint32(1<<spi.datalength) - uint32(1)
	data = data & mask
	spi.regs.txdata = data

	//wait for rx fifo to get data
	for spi.regs.status&(0x1<<3) == 0 {
	}
	return spi.regs.rxdata
}
