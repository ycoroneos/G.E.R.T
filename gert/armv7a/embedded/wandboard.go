package embedded

import (
	"unsafe"
)

//SPI1 Pins
// MOSI EIM D18

//3-11 EIM_DA11
//3-27 EIM_D27
//6-31 EIM_BCLK
//1-24 ENET_RX_ER

var WB_JP4_4 = GPIO_pin{"JP4_4", 3, 11, gpios[3-1], IOMUX_MUX_CTL_EIM_DA11, IOMUX_PAD_CTL_EIM_DA11}
var WB_JP4_6 = GPIO_pin{"JP4_6", 3, 27, gpios[3-1], IOMUX_MUX_CTL_EIM_D27, IOMUX_PAD_CTL_EIM_D27}
var WB_JP4_8 = GPIO_pin{"JP4_8", 6, 31, gpios[6-1], IOMUX_MUX_CTL_EIM_BCLK, IOMUX_PAD_CTL_EIM_BCLK}
var WB_JP4_10 = GPIO_pin{"JP4_10", 1, 24, gpios[1-1], IOMUX_MUX_CTL_ENET_RX_ER, IOMUX_PAD_CTL_ENET_RX_ER}

var WB_SPI1 = SPI_periph{SPI_pin{"mosi", 1, IOMUX_MUX_CTL_EIM_D17, IOMUX_PAD_CTL_EIM_D17},
	SPI_pin{"miso", 1, IOMUX_MUX_CTL_EIM_D18, IOMUX_PAD_CTL_EIM_D18},
	SPI_pin{"sclk", 1, IOMUX_MUX_CTL_EIM_D16, IOMUX_PAD_CTL_EIM_D16},
	SPI_pin{"chip select", 0, IOMUX_MUX_CTL_KEY_COL2, IOMUX_PAD_CTL_KEY_COL2},
	0,
	0,
	1}

var WB_PWM3 = PWM_periph{PWM_pin{"JP1_3", 2, IOMUX_MUX_CTL_SD4_DATA1, IOMUX_PAD_CTL_SD4_DATA1}, ((*PWM_regs)(unsafe.Pointer(uintptr(0x2088000))))}
var WB_PWM4 = PWM_periph{PWM_pin{"JP1_5", 2, IOMUX_MUX_CTL_SD4_DATA2, IOMUX_PAD_CTL_SD4_DATA2}, ((*PWM_regs)(unsafe.Pointer(uintptr(0x208C000))))}

var WB_DEFAULT_UART = UART{((*UART_regs)(unsafe.Pointer(uintptr(0x2020000))))}
