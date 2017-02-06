package embedded

import (
	"unsafe"
)

/*
* Translated again from the imx6 platform sdk
*
 */

type SPI_dev uint32

const (
	DEV_SPI1 = 1
	DEV_SPI2 = 2
	DEV_SPI3 = 3
	DEV_SPI4 = 4
	DEV_SPI5 = 5
)

type spi_desc struct {
	channel    uint8
	mode       uint8
	ss_pol     uint8
	sclk_pol   uint8
	sclk_phase uint8
	pre_div    uint8
	post_div   uint8
}

//either mosi, miso, or sclk
type SPI_pin struct {
	name   string
	alt    uint8
	muxctl *uint32
	padctl *uint32
}

type SPI_periph struct {
	mosi           SPI_pin
	miso           SPI_pin
	sclk           SPI_pin
	cs             SPI_pin
	mode           uint8
	frequency      uint8
	channel_select uint8
}

//SPI has 3 types of modes which affect the polarity of the clock and the resting state of the data signals
//data length is how many bits each SPI frame contains. 7,8,16 are common amounts
func (spi SPI_periph) Begin(mode, freq, datalength uint8) {
	//put the gpio pins on push/pull mode with their appropriate alternate functions
	*spi.mosi.muxctl = makeGPIOmuxconfig(spi.mosi.alt)
	*spi.mosi.padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)

	*spi.miso.muxctl = makeGPIOmuxconfig(spi.miso.alt)
	*spi.miso.padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)

	*spi.sclk.muxctl = makeGPIOmuxconfig(spi.sclk.alt)
	*spi.sclk.padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)

	*spi.cs.muxctl = makeGPIOmuxconfig(spi.cs.alt)
	*spi.cs.padctl = makeGPIOpadconfig(0, PULLDOWN_100K, 0, 0, 0, SPEED_FAST, DRIVE_33R, SLEW_FAST)

	//ungate the module clock
	//Put 0x3F into CCM_CCGR1
	CCM_CCGR1 := ((*uint32)(unsafe.Pointer(uintptr(0x20C406C))))
	*CCM_CCGR1 |= 0x3F

	//configure the SPI registers for the freq, mode, and datalength
}

//ecspi_open(dev SPI_dev, param_ecspi_t spi_desc) int {
//    // Configure IO signals
//    ecspi_iomux_config(dev)
//
//    // Ungate the module clock.
//		//Put 0x3F into CCM_CCGR1
//    clock_gating_config(REGS_ECSPI_BASE(dev), CLOCK_ON);
//
//    // Configure eCSPI registers
//    ecspi_configure(dev, param);
//
//    return SUCCESS;
//}
