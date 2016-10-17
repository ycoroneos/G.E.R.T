#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include "console.h"

// go elf kernel
extern char gobin_start[], gobin_end[];

// go bin load address
uintptr_t go_load_addr=(uintptr_t)0x10000000;


static void backtrace(uint32_t *fp)
{
  unsigned i=0;
  while (*fp)
  {
    uint32_t *lr = fp-1;
    uint32_t *newfp = fp-3;
    //cprintf("frame %d: lr 0x%x
    ++i;
  }
}

static void boot_memcpy(char *dst, char *src, size_t size)
{
  for (size_t i=0; i<size; ++i)
  {
    dst[i]=src[i];
  }
}

static void boot_memset(char *dst, char val, size_t size)
{
  for (size_t i=0; i<size; ++i)
  {
    dst[i]=val;
  }
}

void reset_interrupt()
{
  panic("unintended reset");
  volatile short stub=1;
  while (stub) {};
}

void undefined_interrupt()
{
  volatile short stub=1;
  while (stub) {};
}

void svc_interrupt()
{
  uint32_t lr;
  asm volatile("mov %[a], lr" : [a] "=r" (lr) :);
  lr;
  cprintf("svc from addr 0x%x\r\n", lr-4);
  ((void (*)(void)) (lr))();
  volatile short stub=1;
  while (stub) {};
}

void prefetch_abort_interrupt()
{
  cprintf("prefetch abort");
  volatile short stub=1;
  while (stub) {};
}

void data_abort_interrupt()
{
  uint32_t lr;
  asm volatile("mov %[a], lr" : [a] "=r" (lr) :);
  lr-=8;
  cprintf("data abort from addr 0x%x\r\n", lr);
  volatile short stub=1;
  while (stub) {};
}

void irq_interrupt()
{
  cprintf("irq");
  volatile short stub=1;
  while (stub) {};
}

void fiq_interrupt()
{
  cprintf("fiq");
  volatile short stub=1;
  while (stub) {};
}

static void load_go()
{
  size_t binsize = gobin_end - gobin_start;
  cprintf("go bin size : 0x%x\r\n", binsize);
  //boot_memset((char *)go_load_addr, 0, binsize);
  boot_memcpy((char *)go_load_addr, (char *)gobin_start, binsize);
  cprintf("loaded at 0x%x\r\n", go_load_addr);
  cprintf("first 10 words are:\r\n");
  uint32_t *in = (uint32_t *)(gobin_start);
  uint32_t *out = (uint32_t *)(go_load_addr);
  for (unsigned i=0; i<10; ++i)
  {
    cprintf("\t0x%x vs 0x%x\r\n", out[i], in[i]);
  }
  uint32_t nwords = binsize /4;
  cprintf("verifying: ");
  for (unsigned i=0; i<nwords; ++i)
  {
    if (in[i] != out[i])
    {
      cprintf("mismatch at address 0x%x\r\n", &out[i]);
      panic("verify failed");
    }
  }
  cprintf(" complete!\r\n");
}

int main()
{
  //get a new stack
 // asm("mov sp,%[newsp]" : : [newsp]"r"(&stack[1023]));

  consoleinit();
  cprintf("---------------------------------------------\r\n");
  cprintf("Welcome to Biscuit-ARM Bootloader, 16 is %x hex!\r\n",16);

  float a=1.25;
  float b =a*3;
  cprintf("float multiplication %x\r\n", b);
  //load our kernel
  load_go();
  uint32_t kernel_start = go_load_addr;
  uint32_t kernel_size = gobin_end - gobin_start;
  asm volatile("mov r0, %0"
      :
      :"r"(kernel_start)
      :"r0"
      );
  asm volatile("mov r1, %0"
      :
      :"r"(kernel_size)
      :"r0"
      );
  //sub sp, 16
  asm volatile("sub sp, #16"
      :
      :
      :
      );
  //and sp -16
  asm volatile("and sp, #-16"
      :
      :
      :
      );
  ((void (*)(void)) (go_load_addr))();
  panic("should not be here");
}
