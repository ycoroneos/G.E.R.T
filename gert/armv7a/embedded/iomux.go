package embedded

import "unsafe"
import "fmt"

/*
* The IOMUX peripheral essentially has hundreds of different 32bit registers and each one controls properties of a specific pin in the SOC.
* I cant support everything in this lifetime so I will pick some pins that are exposed on the header of the Wandboard Quad
 */

//for sdcard
//signal 		| pad 			| mode | direction
//SD1_CD_B 	  GPIO_1  		ALT6  	IN
//SD1_CLK     SD1_CLK  		ALT0 	  OUT
//SD1_CMD 		SD1_CMD 		ALT0    IN/OUT
//SD1_DATA0 	SD1_DAT0 		ALT0 		IN/OUT
//SD1_DATA1 	SD1_DAT1 		ALT0 		IN/OUT
//SD1_DATA2 	SD1_DAT2 		ALT0 		IN/OUT
//SD1_DATA3 	SD1_DAT3 		ALT0 		IN/OUT
//SD1_WP 			GPIO4_15 		ALT0 		IN/OUT ?????

var IOMUX_GPIO1 = ((*uint32)(unsafe.Pointer((uintptr(0x20E0224)))))
var IOMUX_SD1_CLK = ((*uint32)(unsafe.Pointer((uintptr(0x20E0350)))))
var IOMUX_SD1_CMD = ((*uint32)(unsafe.Pointer((uintptr(0x20E0348)))))
var IOMUX_SD1_DATA0 = ((*uint32)(unsafe.Pointer((uintptr(0x20E0340)))))
var IOMUX_SD1_DATA1 = ((*uint32)(unsafe.Pointer((uintptr(0x20E033C)))))
var IOMUX_SD1_DATA2 = ((*uint32)(unsafe.Pointer((uintptr(0x20E034C)))))
var IOMUX_SD1_DATA3 = ((*uint32)(unsafe.Pointer((uintptr(0x20E0344)))))

//gpios
var IOMUX_MUX_CTL_GPIO0 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0220))))
var IOMUX_MUX_CTL_GPIO1 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0224))))
var IOMUX_MUX_CTL_GPIO2 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0234))))
var IOMUX_MUX_CTL_GPIO3 = ((*uint32)(unsafe.Pointer(uintptr(0x20E022C))))
var IOMUX_MUX_CTL_GPIO4 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0238))))
var IOMUX_MUX_CTL_GPIO5 = ((*uint32)(unsafe.Pointer(uintptr(0x20E023C))))
var IOMUX_MUX_CTL_GPIO6 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0230))))
var IOMUX_MUX_CTL_GPIO7 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0240))))
var IOMUX_MUX_CTL_GPIO8 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0244))))
var IOMUX_MUX_CTL_GPIO9 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0228))))
var IOMUX_MUX_CTL_GPIO16 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0248))))
var IOMUX_MUX_CTL_GPIO17 = ((*uint32)(unsafe.Pointer(uintptr(0x20E024C))))
var IOMUX_MUX_CTL_GPIO18 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0250))))
var IOMUX_MUX_CTL_GPIO19 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0254))))
var IOMUX_MUX_CTL_EIM_DA11 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0140))))
var IOMUX_MUX_CTL_EIM_D27 = ((*uint32)(unsafe.Pointer(uintptr(0x20E00C0))))
var IOMUX_MUX_CTL_EIM_BCLK = ((*uint32)(unsafe.Pointer(uintptr(0x20E0158))))
var IOMUX_MUX_CTL_ENET_RX_ER = ((*uint32)(unsafe.Pointer(uintptr(0x20E01D8))))

var IOMUX_PAD_CTL_GPIO0 = ((*uint32)(unsafe.Pointer(uintptr(0x20E05F0))))
var IOMUX_PAD_CTL_GPIO1 = ((*uint32)(unsafe.Pointer(uintptr(0x20E05F4))))
var IOMUX_PAD_CTL_GPIO2 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0604))))
var IOMUX_PAD_CTL_GPIO3 = ((*uint32)(unsafe.Pointer(uintptr(0x20E05FC))))
var IOMUX_PAD_CTL_GPIO4 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0608))))
var IOMUX_PAD_CTL_GPIO5 = ((*uint32)(unsafe.Pointer(uintptr(0x20E060C))))
var IOMUX_PAD_CTL_GPIO6 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0600))))
var IOMUX_PAD_CTL_GPIO7 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0610))))
var IOMUX_PAD_CTL_GPIO8 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0614))))
var IOMUX_PAD_CTL_GPIO9 = ((*uint32)(unsafe.Pointer(uintptr(0x20E05F8))))
var IOMUX_PAD_CTL_GPIO16 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0618))))
var IOMUX_PAD_CTL_GPIO17 = ((*uint32)(unsafe.Pointer(uintptr(0x20E061C))))
var IOMUX_PAD_CTL_GPIO18 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0620))))
var IOMUX_PAD_CTL_GPIO19 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0624))))
var IOMUX_PAD_CTL_EIM_DA11 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0454))))
var IOMUX_PAD_CTL_EIM_D27 = ((*uint32)(unsafe.Pointer(uintptr(0x20E03D4))))
var IOMUX_PAD_CTL_EIM_BCLK = ((*uint32)(unsafe.Pointer(uintptr(0x20E046C))))
var IOMUX_PAD_CTL_ENET_RX_ER = ((*uint32)(unsafe.Pointer(uintptr(0x20E04EC))))

//SPI1
var IOMUX_MUX_CTL_EIM_D17 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0094))))
var IOMUX_MUX_CTL_EIM_D18 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0098))))
var IOMUX_MUX_CTL_EIM_D16 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0090))))
var IOMUX_MUX_CTL_KEY_COL2 = ((*uint32)(unsafe.Pointer(uintptr(0x20E0208))))

var IOMUX_PAD_CTL_EIM_D17 = ((*uint32)(unsafe.Pointer(uintptr(0x20E03A8))))
var IOMUX_PAD_CTL_EIM_D18 = ((*uint32)(unsafe.Pointer(uintptr(0x20E03AC))))
var IOMUX_PAD_CTL_EIM_D16 = ((*uint32)(unsafe.Pointer(uintptr(0x20E03A4))))
var IOMUX_PAD_CTL_KEY_COL2 = ((*uint32)(unsafe.Pointer(uintptr(0x20E05D8))))

//go:nosplit
func usdhc_iomux_config(instance uint32) {
	switch instance {
	case 0:
		fmt.Println("SDHC instance 0 does not exist, perhaps you mean 1?")
	case 1:
		*IOMUX_GPIO1 = 0x1<<4 | 0x6 //SION and ALT6
		*IOMUX_SD1_CLK = 0x1 << 4
		*IOMUX_SD1_CMD = 0x1 << 4
		*IOMUX_SD1_DATA0 = 0x1 << 4
		*IOMUX_SD1_DATA1 = 0x1 << 4
		*IOMUX_SD1_DATA2 = 0x1 << 4
		*IOMUX_SD1_DATA3 = 0x1 << 4
	case 2:
		fmt.Println("SDHC instance 2 has no sdcard slot on the wandboard")
	case 3:
		fmt.Println("SDHC instance 3 not supported yet")
	}
}

//this function applies to pins 4,6,8,10,12,14,16,18 on JP4 of the wandboard
const (
	PULLDOWN_100K = 0
	PULLUP_47K    = 1
	PULLUP_100K   = 2
	PULLUP_22K    = 3
	DRIVE_HIZ     = 0
	DRIVE_260R    = 1
	DRIVE_130R    = 2
	DRIVE_90R     = 3
	DRIVE_60R     = 4
	DRIVE_50R     = 5
	DRIVE_40R     = 6
	DRIVE_33R     = 7
	SLEW_SLOW     = 0
	SLEW_FAST     = 1
	MUX_ALT0      = 0
	MUX_ALT1      = 1
	MUX_ALT2      = 2
	MUX_ALT3      = 3
	MUX_ALT4      = 4
	MUX_ALT5      = 5
	MUX_ALT6      = 6
	MUX_ALT7      = 7
	SPEED_LOW     = 0
	SPEED_MEDIUM  = 1
	SPEED_FAST    = 3
)

var pinmap = map[uint32]*uint32{}

func makeGPIOmuxconfig(muxmode uint8) uint32 {
	muxmode &= 7
	return uint32(muxmode)
}

func makeGPIOpadconfig(hysteresis, pull, pull_keep_mode, pull_keep_enabled, open_drain, speed, drive_strength, slewrate uint32) uint32 {
	hysteresis &= 0x1
	pull &= 0x3
	pull_keep_mode &= 0x1
	pull_keep_enabled &= 0x1
	open_drain &= 0x1
	drive_strength &= 0x7
	slewrate &= 0x1
	speed &= 0x3
	return (uint32(hysteresis) << 16) | (uint32(pull) << 14) | (uint32(pull_keep_mode) << 13) | (uint32(pull_keep_enabled) << 12) | (uint32(open_drain) << 11) | (uint32(speed) << 6) | (uint32(drive_strength) << 3) | (uint32(slewrate))
}
