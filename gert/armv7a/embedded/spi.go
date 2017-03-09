package embedded

import (
	"unsafe"
)

/*
* Translated again from the imx6 platform sdk
*
 */

type SPI_dev uint32

//type SPI_regs struct {
//	channel    uint8
//	mode       uint8
//	ss_pol     uint8
//	sclk_pol   uint8
//	sclk_phase uint8
//	pre_div    uint8
//	post_div   uint8
//}

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
	mosi      SPI_pin
	miso      SPI_pin
	sclk      SPI_pin
	cs        []SPI_pin
	regs      *SPI_regs
	mode      uint8
	frequency uint8
}

//SPI has 3 types of modes which affect the polarity of the clock and the resting state of the data signals
//data length is how many bits each SPI frame contains. 7,8,16 are common amounts

//the imx6 can support up to 2^12 bits in a single frame!
func (spi *SPI_periph) Begin(mode, freq, datalength uint16) {
	//put the gpio pins on push/pull mode with their appropriate alternate functions
	*spi.mosi.muxctl = makeGPIOmuxconfig(spi.mosi.alt)
	*spi.mosi.padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)

	*spi.miso.muxctl = makeGPIOmuxconfig(spi.miso.alt)
	*spi.miso.padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)

	*spi.sclk.muxctl = makeGPIOmuxconfig(spi.sclk.alt)
	*spi.sclk.padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)

	for i := 0; i < len(spi.cs); i++ {
		*spi.cs[i].muxctl = makeGPIOmuxconfig(spi.cs[i].alt)
		*spi.cs[i].padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)
	}

	//ungate the module clock
	//Put 0x3F into CCM_CCGR1
	CCM_CCGR1 := ((*uint32)(unsafe.Pointer(uintptr(0x20C406C))))
	*CCM_CCGR1 |= 0x3F

	//configure the SPI registers for the freq, mode, and datalength
}
