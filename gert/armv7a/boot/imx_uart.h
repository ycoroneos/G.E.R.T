// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef IMX_UART_H
#define IMX_UART_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

//in which I write my own periperal library...
#define BAUD_RATE 115200
#define  UCR4_REF16 (1<<6)  /* Ref freq 16 MHz */
#define  UFCR_RFDIV      (7<<7)  /* Reference freq divider mask */
#define CONFIG_SYSPLL_CLK_FREQ 792000000
#define PLL3_FREQUENCY 80000000

#define IOBASE 0x0

#define CCM_CSCDR1   (uint32_t*)(0x020C4024+IOBASE)

#define UART1_UTXD   (uint32_t*)(0x02020040+IOBASE)
#define UART1_UCR1   (uint32_t*)(0x02020080+IOBASE)
#define UART1_UCR2   (uint32_t*)(0x02020084+IOBASE)
#define UART1_UCR3   (uint32_t*)(0x02020088+IOBASE)
#define UART1_UCR4   (uint32_t*)(0x0202008C+IOBASE)
#define UART1_UFCR   (uint32_t*)(0x02020090+IOBASE)
#define UART1_USR1   (uint32_t*)(0x02020094+IOBASE)
#define UART1_USR2   (uint32_t*)(0x02020098+IOBASE)
#define UART1_UESC   (uint32_t*)(0x0202009C+IOBASE)
#define UART1_UTIM   (uint32_t*)(0x020200A0+IOBASE)
#define UART1_UBIR   (uint32_t*)(0x020200A4+IOBASE)
#define UART1_UBMR   (uint32_t*)(0x020200A8+IOBASE)
#define UART1_UBRC   (uint32_t*)(0x020200AC+IOBASE)
#define UART1_ONEMS  (uint32_t*)(0x020200B0+IOBASE)
#define UART1_UTS    (uint32_t*)(0x020200B4+IOBASE)
#define UART1_UMCR   (uint32_t*)(0x020200B8+IOBASE)

//nope later...
#define UART2_UTXD   (uint32_t*)0x021E8040
#define UART2_UCR1   (uint32_t*)0x021E8080
#define UART2_UCR2   (uint32_t*)0x021E8084
#define UART2_UCR3   (uint32_t*)0x021E8088
#define UART2_UCR4   (uint32_t*)0x021E808C
#define UART2_UFCR   (uint32_t*)0x021E8090
#define UART2_USR1   (uint32_t*)0x021E8094
#define UART2_USR2   (uint32_t*)0x021E8098
#define UART2_UESC   (uint32_t*)0x021E809C
#define UART2_UTIM   (uint32_t*)0x021E80A0
#define UART2_UBIR   (uint32_t*)0x021E80A4
#define UART2_UBMR   (uint32_t*)0x021E80A8
#define UART2_UBRC   (uint32_t*)0x021E80AC
#define UART2_ONEMS  (uint32_t*)0x021E80B0
#define UART2_UTS    (uint32_t*)0x021E80B4
#define UART2_UMCR   (uint32_t*)0x021E80B8

void uart_setup();
void uart_tx(uint8_t);
void imx_puts(char *);
#endif
