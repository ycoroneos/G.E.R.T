cmd_u-boot.img := ./tools/mkimage -A arm -T firmware -C none -O u-boot -a 0x17800000 -e 0 -n "U-Boot 2016.01 for wandboard board" -d u-boot.bin u-boot.img  >/dev/null
