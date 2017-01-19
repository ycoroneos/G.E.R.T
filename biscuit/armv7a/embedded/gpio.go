package embedded

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

var gpios = [...]*gpio {
	((*gpio)(unsafe.Pointer(uintptr(GPIO1_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO2_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO3_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO4_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO5_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO7_BASE)))),
	((*gpio)(unsafe.Pointer(uintptr(GPIO6_BASE))))
}


///user api
//section 28.4.3.1
func PinSetInput(pnum uint32) {
}

//section 28.4.3.2
func PinSetOutput(pnum uint32) {
}

//this is either hi or lo
//section 28.4.3.2
func PinWrite(pnum uint32, val uint8) {
}

//this is either hi or lo
//section 28.4.3.1
func PinRead(pnum uint32) uint8 {
}
////
