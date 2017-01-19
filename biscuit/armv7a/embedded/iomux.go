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
)

func MakeGPIOconfig(hysteresis, pull, pull_keep_mode, pull_keep_enabled, open_drain, drive_strength, slewrate uint8) uint32 {
	hysteresis &= 0x1
	pull &= 0x3
	pull_keep_mode &= 0x1
	pull_keep_enable &= 0x1
	open_drain &= 0x1
	drive_strength &= 0x7
	slewrate &= 0x1
	return (uint32(hysteresis) << 16) | (uint32(pull) << 14) | (uint32(pull_keep_mode) << 13) | (uint32(pull_keep_enabled) << 12) | (uint32(open_drain) << 11) | (uint32(drive_strength) << 3) | (uint32(slewrate))
}
