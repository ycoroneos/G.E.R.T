cmd_u-boot.srec := arm-none-eabi-objcopy --gap-fill=0xff  -j .text -j .secure_text -j .rodata -j .hash -j .data -j .got -j .got.plt -j .u_boot_list -j .rel.dyn -O srec u-boot u-boot.srec
