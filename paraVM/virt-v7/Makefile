#Makefile for ARM barebones development

UPROG ?= .

GCCPREFIX ?= arm-none-eabi-
BUILDIR := obj
XVISOR := $(realpath ../../xvisor/)

GOLINK := 0x40000000
GOBIN := $(realpath ../../golang_embedded/bin/go)
GOARCH := GOOS=linux GOARCH=arm GOARM=7
GOLINKFLAGS := "-T $(GOLINK) -E main.Entry"
GOCMD := $(GOARCH) $(GOBIN)
GO_SRC = $(wildcard $(UPROG)/*.go)
GO_OBJ := $(BUILDIR)/kernel

LINKADDR := 0x45000000
C_SRC += $(wildcard boot/*.c)
AS_SRC += $(wildcard boot/*.S)
C_OBJ := $(patsubst boot/%.c,$(BUILDIR)/%.c.o,$(C_SRC))
AS_OBJ := $(patsubst boot/%.S,$(BUILDIR)/%.S.o,$(AS_SRC))
BOOT_OBJ := $(AS_OBJ) $(C_OBJ)
BOOT_TARGET := $(BUILDIR)/bootloader


CFLAGS := -O0 -std=gnu99 -Wall -Wextra -fno-builtin -mfpu=neon-vfpv4 -mfloat-abi=soft \
          -march=armv7ve -mtune=cortex-a15 -ffreestanding -nostartfiles -fno-pic -Iboot/ \
					-DLINK=$(LINKADDR)

LFLAGS := $(CFLAGS) -Ttext $(LINKADDR) -e_init_vectors


# Targets
.PHONY : all runtime uboot run

all: $(BOOT_TARGET).elf
	$(GCCPREFIX)objcopy -O binary $(BOOT_TARGET).elf $(BOOT_TARGET).bin
	$(GCCPREFIX)objdump -D $(BOOT_TARGET).elf > $(BOOT_TARGET).dump
	$(GCCPREFIX)readelf -a $(BOOT_TARGET).elf > $(BOOT_TARGET).read
	mkdir -p $(BUILDIR)/build/disk/tmp
	mkdir -p $(BUILDIR)/build/disk/system
	cp -f $(XVISOR)/docs/banner/roman.txt $(BUILDIR)/build/disk/system/banner.txt
	cp -f $(XVISOR)/docs/logo/xvisor_logo_name.ppm $(BUILDIR)/build/disk/system/logo.ppm
	mkdir -p $(BUILDIR)/build/disk/images/arm32/virt-v7
	$(XVISOR)/build/tools/dtc/bin/dtc -I dts -O dtb -o $(BUILDIR)/build/disk/images/arm32/virt-v7x2.dtb $(XVISOR)/tests/arm32/virt-v7/virt-v7x2.dts
	cp -f $(BOOT_TARGET).bin $(BUILDIR)/build/disk/images/arm32/virt-v7/firmware.bin
	cp -f boot/nor_flash.list $(BUILDIR)/build/disk/images/arm32/virt-v7/nor_flash.list
	genext2fs -B 1024 -b 16384 -d $(BUILDIR)/build/disk $(BUILDIR)/build/disk.img
	arm-linux-gnueabi-gcc -nostdlib -march=armv7ve -mcpu=cortex-a15 -e start_boot -Wl,--build-id=none -Wl,-Ttext=0x80000000 \
		-DSPIN_LOCATION=0x1c010030 -DSPIN_LOOP_ADDR=0x14000000 -DUART_PL011 -DUART_PL011_BASE=0x1c090000 \
		-DGENTIMER_FREQ=100000000 -DGICv2 -DGIC_DIST_BASE=0x2c001000 -DGIC_CPU_BASE=0x2c002000 \
		-DIMAGE=$(XVISOR)/build/vmm.bin -DDTB=$(XVISOR)/build/arch/arm/board/generic/dts/vexpress/a15/one_guest_virt-v7.dtb \
		-DINITRD=$(BUILDIR)/build/disk.img $(XVISOR)/docs/arm/fast_model_boot.S -o $(BUILDIR)/build/fast_model_boot.axf

#linker
$(BOOT_TARGET).elf: $(GO_OBJ).elf $(BOOT_OBJ)
	@echo "linking " $(BOOT_OBJ)
	@$(GCCPREFIX)gcc $(LFLAGS) -o $@ $(BOOT_OBJ)

#build GERT program
$(GO_OBJ).elf: $(GO_SRC)
	@echo "+GO "$(GO_SRC)
	@$(GOCMD) build -ldflags $(GOLINKFLAGS) -o $@ $(GO_SRC)
	$(GCCPREFIX)objdump -D $(GO_OBJ).elf > $(GO_OBJ).dump
	$(GCCPREFIX)readelf -a $(GO_OBJ).elf > $(GO_OBJ).read
	$(GCCPREFIX)objcopy -O binary $(GO_OBJ).elf $(GO_OBJ).bin
	./hackbin.py $(GO_OBJ).bin $(GO_OBJ)_hacked.bin
	hexdump $(GO_OBJ).bin > $(GO_OBJ).hex
	hexdump $(GO_OBJ)_hacked.bin > $(GO_OBJ)_hacked.hex

#build a c file
$(BUILDIR)/%.c.o : boot/%.c
	@mkdir -p $(BUILDIR)
	@echo "+CC "$< " -> " $@
	@$(GCCPREFIX)gcc $(CFLAGS) -c $< -o $@

#build an asm file
$(BUILDIR)/%.S.o : boot/%.S
	@mkdir -p $(BUILDIR)
	@echo "+AS "$< " -> " $@
	$(GCCPREFIX)gcc $(CFLAGS) -c $< -o $@

clean:
	-rm -rf $(BUILDIR)

bytecode: test.s
	$(GCCPREFIX)gcc $(CFLAGS) -c test.s -o test.o
	$(GCCPREFIX)objdump -d -S test.o
print-%: ; @echo $*=$($*)

runtime:
	cd ../../golang_embedded/src && ./rebuild_arm.sh

run:
	ARMLMD_LICENSE_FILE=~/Downloads/license.dat ~/FVP_VE_Cortex-A15x124/models/Linux64_GCC-4.9/FVP_VE_Cortex-A15x4 $(BUILDIR)/build/fast_model_boot.axf

rund:
	ARMLMD_LICENSE_FILE=~/Downloads/license.dat ~/FVP_VE_Cortex-A15x124/models/Linux64_GCC-4.9/FVP_VE_Cortex-A15x4 $(BUILDIR)/build/fast_model_boot.axf -S
