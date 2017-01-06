package main

import "unsafe"
import "fmt"

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
