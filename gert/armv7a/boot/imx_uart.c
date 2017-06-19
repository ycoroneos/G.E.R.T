// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "imx_uart.h"

void uart_setup()
{
  //disable the UART
  *UART1_UCR1=0;
  //software reset
  *UART1_UCR2=0;
  //wait to reset
  while (!(*UART1_UCR2&0x1)) {};
  //enable the UART
  *UART1_UCR1|=0x1;
  //RXDMUXEL must be on
  *UART1_UCR3|=0x4;
  //8 bits, 1 stop bit, no parity,software flow control
  //              8bits    ignore RTS enable RX  enable tx   dont reset
  *UART1_UCR2|= (0x1<<5) | (0x1<<14) | 0x2      | (0x1<<2) | 0x1;
  *UART1_UFCR |= (0x04 << 7);
  //set the divider ratios
  *UART1_UBIR = 0x0F;
  *UART1_UBMR = (PLL3_FREQUENCY / ((*CCM_CSCDR1 & 0x3F) + 1)) / (2 * BAUD_RATE);
}

void uart_tx(uint8_t c)
{
  *UART1_UTXD = c; //write char to tx buffer
  //wait for send
  while (!(*UART1_UTS&(0x1<<6))) {};
  return;
}

void imx_puts(char *string)
{
  uint32_t i;
  for (i=0; string[i]!=0; ++i)
  {
    uart_tx(string[i]);
  }
}
