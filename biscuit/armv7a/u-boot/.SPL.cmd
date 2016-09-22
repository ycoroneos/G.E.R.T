cmd_SPL := ./tools/mkimage -n arch/arm/imx-common/spl_sd.cfg.cfgtmp -T imximage -e 0x00908000 -d spl/u-boot-spl.bin SPL  >/dev/null
