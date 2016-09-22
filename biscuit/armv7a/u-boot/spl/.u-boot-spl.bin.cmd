cmd_spl/u-boot-spl.bin := arm-none-eabi-objcopy  -j .text -j .secure_text -j .rodata -j .hash -j .data -j .got -j .got.plt -j .u_boot_list -j .rel.dyn  -O binary spl/u-boot-spl spl/u-boot-spl.bin
