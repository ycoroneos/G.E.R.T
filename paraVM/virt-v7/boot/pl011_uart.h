
// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef PL011_UART_H
#define PL011_UART_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

#define BAUD_RATE 115200
#define CLOCK 24000000

//register settings
#define WLEN_8  0x11<<5
#define FEN     0x1<<4
#define RXE     0x1<<9
#define TXE     0x1<<8
#define EN      0x1
#define TXFF    0x1<<5
#define RXFE 0x10


//register definitions
#define UARTBASE    0x9000000
//#define UARTBASE    0x1c090000
#define UARTDR      (uint32_t*)(0x0+UARTBASE)
#define UARTRSR     (uint32_t*)(0x4+UARTBASE)
#define UARTFR      (uint32_t*)(0x18+UARTBASE)
#define UARTILPR    (uint32_t*)(0x20+UARTBASE)
#define UARTIBRD    (uint32_t*)(0x24+UARTBASE)
#define UARTFBRD    (uint32_t*)(0x28+UARTBASE)
#define UARTLCR_H   (uint32_t*)(0x2C+UARTBASE)
#define UARTCR      (uint32_t*)(0x30+UARTBASE)
#define UARTIFLS    (uint32_t*)(0x34+UARTBASE)
#define UARTIMSC    (uint32_t*)(0x38+UARTBASE)
#define UARTRIS     (uint32_t*)(0x3C+UARTBASE)
#define UARTMIS     (uint32_t*)(0x40+UARTBASE)
#define UARTICR     (uint32_t*)(0x44+UARTBASE)
#define UARTDMACR   (uint32_t*)(0x48+UARTBASE)
#define UARTPID0    (uint32_t*)(0xFE0+UARTBASE)
#define UARTPID1    (uint32_t*)(0xFE4+UARTBASE)
#define UARTPID2    (uint32_t*)(0xFE8+UARTBASE)
#define UARTPID3    (uint32_t*)(0xFEC+UARTBASE)


void uart_setup();
void uart_tx(uint8_t);
void puts(char *);
char uart_getc();

#endif
