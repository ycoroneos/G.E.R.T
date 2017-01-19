package embedded

import (
	"unsafe"
)

/*
* In the iMX6 all the gpios are layered on top of a muxer, the IOMUX peripheral.
*	As a result, the settings you enter into the GPIO register may or may not have any effect.
* Also make sure to set the correct IOMUX settings.
* Section 28.5 from iMX6 Quad Applications Manual
 */

const (
	GPIO1_BASE = 0x209C000
	GPIO2_BASE = 0x20A0000
	GPIO3_BASE = 0x20A4000
	GPIO4_BASE = 0x20A8000
	GPIO5_BASE = 0x20AC000
	GPIO6_BASE = 0x20B0000
	GPIO7_BASE = 0x20B4000

	//GPIO3_12
	//WB_JP4_4 = 76
)

type gpio struct {
	dr       uint32
	gdir     uint32
	psr      uint32
	icr1     uint32
	icr2     uint32
	imr      uint32
	isr      uint32
	edge_sel uint32
}

var gpios = [...]*gpio{
	((*gpio)(unsafe.Pointer(uintptr(GPIO1_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO2_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO3_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO4_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO5_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO7_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO6_BASE)))),
}

/*
* Turns out individual pins are complicated things in the iMX6.
* Lets use go to make them easier to use
 */
type GPIO_pin struct {
	name     string
	base     uint32
	offset   uint32
	gpioregs *gpio
	muxctl   *uint32
	padctl   *uint32
}

///user api
const (
	INTR_LOW     = 0
	INTR_HIGH    = 1
	INTR_RISING  = 2
	INTR_FALLING = 3
)

//section 28.4.3.1
func (pin GPIO_pin) SetInput() {
	//set gpio mode in iomux
	mux := makeGPIOmuxconfig(MUX_ALT5)
	*pin.muxctl = mux
	pad := makeGPIOpadconfig(0, PULLDOWN_100K, 1, 1, 0, DRIVE_HIZ, SLEW_FAST)
	*pin.padctl = pad

	//set gdir to 0
	pin.gpioregs.gdir &= ^(0x1 << pin.offset)
}

//section 28.4.3.2
func (pin GPIO_pin) SetOutput() {
	//set gpio mode in iomux
	mux := makeGPIOmuxconfig(MUX_ALT5)
	*pin.muxctl = mux
	pad := makeGPIOpadconfig(0, PULLDOWN_100K, 1, 1, 0, DRIVE_HIZ, SLEW_FAST)
	*pin.padctl = pad

	//set gdir to 1
	pin.gpioregs.gdir |= (0x1 << pin.offset)
}

//this is either hi or lo
//section 28.4.3.2
func (pin GPIO_pin) Write(val uint8) {
	if (val & 0x1) == 1 {
		pin.gpioregs.dr |= (0x1 << pin.offset)
	} else {
		pin.gpioregs.dr &= ^(0x1 << pin.offset)
	}
}

//this is either hi or lo
//section 28.4.3.1
func (pin GPIO_pin) Read() uint8 {
	return uint8((pin.gpioregs.dr >> pin.offset) & 0x1)
}

//converts the imx pin numbering into a linear number that I use and the registers use
func GetPinNum(base, number uint32) uint32 {
	return ((base - 1) * 32) + number
}

func (pin GPIO_pin) GetPinNum() uint32 {
	return GetPinNum(pin.base, pin.offset)
}

func (pin GPIO_pin) EnableIntr(mode uint8) {
	mode &= 0x3
	if pin.offset >= 16 {
		icr := &pin.gpioregs.icr2
		offset := pin.offset - 16
		*icr |= uint32(mode) << (2 * offset)
	} else {
		icr := &pin.gpioregs.icr1
		offset := pin.offset
		*icr |= uint32(mode) << (2 * offset)
	}
	pin.gpioregs.imr |= 0x1 << pin.offset
}

func (pin GPIO_pin) DisableIntr() {
	//just mask the interrupt
	pin.gpioregs.imr &= ^(0x1 << pin.offset)
}

////
