MEMORY { .sram : ORIGIN = 0x00908000, LENGTH = 0x10000 }
MEMORY { .sdram : ORIGIN = 0x18200000, LENGTH = 0x100000 }
OUTPUT_FORMAT("elf32-littlearm", "elf32-littlearm", "elf32-littlearm")
OUTPUT_ARCH(arm)
ENTRY(_start)
SECTIONS
{
 .text :
 {
  __start = .;
  *(.vectors)
  arch/arm/cpu/armv7/start.o (.text*)
  *(.text*)
 } >.sram
 . = ALIGN(4);
 .rodata : { *(SORT_BY_ALIGNMENT(.rodata*)) } >.sram
 . = ALIGN(4);
 .data : { *(SORT_BY_ALIGNMENT(.data*)) } >.sram
 . = ALIGN(4);
 .u_boot_list : {
  KEEP(*(SORT(.u_boot_list*_i2c_*)));
 } >.sram
 . = ALIGN(4);
 __image_copy_end = .;
 .end :
 {
  *(.__end)
 }
 .bss :
 {
  . = ALIGN(4);
  __bss_start = .;
  *(.bss*)
  . = ALIGN(4);
  __bss_end = .;
 } >.sdram
}
