package embedded

import (
	"unsafe"
)

/*
*
* This should serve as an example of how to declare some driver structs. I only declare ones I use.
* If GERT were to run on any other imx6-based system, the only thng that would need to change is this file...
*
 */

var WB_JP4_4 = GPIO_pin{"JP4_4", 3, 11, gpios[3-1], IOMUX_MUX_CTL_EIM_DA11, IOMUX_PAD_CTL_EIM_DA11}
var WB_JP4_6 = GPIO_pin{"JP4_6", 3, 27, gpios[3-1], IOMUX_MUX_CTL_EIM_D27, IOMUX_PAD_CTL_EIM_D27}
var WB_JP4_8 = GPIO_pin{"JP4_8", 6, 31, gpios[6-1], IOMUX_MUX_CTL_EIM_BCLK, IOMUX_PAD_CTL_EIM_BCLK}
var WB_JP4_10 = GPIO_pin{"JP4_10", 1, 24, gpios[1-1], IOMUX_MUX_CTL_ENET_RX_ER, IOMUX_PAD_CTL_ENET_RX_ER}

//SPI clock is 59.2MHz
var WB_SPI1 = SPI_periph{SPI_pin{"mosi", 1, IOMUX_MUX_CTL_EIM_D18, IOMUX_PAD_CTL_EIM_D18},
	SPI_pin{"miso", 1, IOMUX_MUX_CTL_EIM_D17, IOMUX_PAD_CTL_EIM_D17},
	SPI_pin{"sclk", 1, IOMUX_MUX_CTL_EIM_D16, IOMUX_PAD_CTL_EIM_D16},
	[]SPI_pin{
		SPI_pin{"channel0", 1, IOMUX_MUX_CTL_EIM_EB2, IOMUX_PAD_CTL_EIM_EB2},
		SPI_pin{"channel1", 0, IOMUX_MUX_CTL_KEY_COL2, IOMUX_PAD_CTL_KEY_COL2},
	},
	((*SPI_regs)(unsafe.Pointer(uintptr(0x2008000)))),
	0,
	0,
	0}

var WB_PWM1 = PWM_periph{PWM_pin{"JP1_17", 2, IOMUX_MUX_CTL_DISP0_DAT8, IOMUX_PAD_CTL_DISP0_DAT8}, ((*PWM_regs)(unsafe.Pointer(uintptr(0x2080000)))), 0, 0.0}
var WB_PWM2 = PWM_periph{PWM_pin{"JP1_19", 2, IOMUX_MUX_CTL_DISP0_DAT9, IOMUX_PAD_CTL_DISP0_DAT9}, ((*PWM_regs)(unsafe.Pointer(uintptr(0x2084000)))), 0, 0.0}

var WB_PWM3 = PWM_periph{PWM_pin{"JP1_3", 2, IOMUX_MUX_CTL_SD4_DATA1, IOMUX_PAD_CTL_SD4_DATA1}, ((*PWM_regs)(unsafe.Pointer(uintptr(0x2088000)))), 0, 0.0}

//var WB_PWM4 = PWM_periph{PWM_pin{"JP1_5", 2, IOMUX_MUX_CTL_SD4_DATA2, IOMUX_PAD_CTL_SD4_DATA2}, ((*PWM_regs)(unsafe.Pointer(uintptr(0x208C000))))}

var WB_DEFAULT_UART = UART{((*UART_regs)(unsafe.Pointer(uintptr(0x2020000))))}
