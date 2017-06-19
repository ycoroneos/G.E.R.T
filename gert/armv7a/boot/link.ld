/* Copyright 2017 Yanni Coroneos. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
*/

ENTRY(_init_vectors)

SECTIONS
{
    . = 0x50000000;
    . = ALIGN(4096); /* align to page size */
    .text :
    {
        *(.text)
    }
    __text_end = .;

    . = ALIGN(4096); /* align to page size */
    __rodata_start = .;
    .rodata :
    {
        *(.rodata)
    }
    __rodata_end = .;

    . = ALIGN(4096); /* align to page size */
    __data_start = .;
    .data :
    {
        *(.data)
    }
    __data_end = .;

    . = ALIGN(4096); /* align to page size */
    __bss_start = .;
    .bss :
    {
        bss = .;
        *(.bss)
    }
    __bss_end = .;
}
