#include "pl011_uart.h"

void uart_setup()
{
  //disable UART
  *UARTCR = 0x0;

  /*
	 * Set baud rate
	 *
	 * IBRD = UART_CLK / (16 * BAUD_RATE)                => 3
	 * FBRD = RND((64 * MOD(UART_CLK,(16 * BAUD_RATE)))  => 1
	 * 	  / (16 * BAUD_RATE))
	 */
//	temp = 16 * BAUDRATE;
//	divider = CLOCK / temp;
//	remainder = arm_umod32(CLOCK, temp);
//	temp = arm_udiv32((8 * remainder), BAUDRATE);
//	fraction = (temp >> 1) + (temp & 1);

  *UARTIBRD=3;
  *UARTFBRD=1;

  //8data bits, 1 stop bit, FIFO enable
  *UARTLCR_H=WLEN_8 | FEN;

  //enable UART
  *UARTCR = EN | RXE | TXE;
}

void uart_tx(uint8_t c)
{
  //wait for space in the fifo
  while (*UARTFR & TXFF);
  *UARTDR=(uint32_t)c;
  return;
}

void puts(char *string)
{
  uint32_t i;
  for (i=0; string[i]!=0; ++i)
  {
    uart_tx(string[i]);
  }
}
