#Makefile for ARM barebones development
SDCARD?=INVALID

GCCPREFIX ?= arm-none-eabi-
BUILDIR := obj

GOBIN := ../../golang1.7/bin/go
GOARCH := GOOS=linux GOARCH=arm GOARM=7
GOLINKFLAGS := "-T 0x11000000 -E main.Entry"
GOCMD := $(GOARCH) $(GOBIN)
GO_SRC = $(wildcard *.go)
#GO_SRC = kernel.go
GO_DIR := go/
GO_OBJ := $(BUILDIR)/kernel

C_SRC += $(wildcard boot/*.c)
AS_SRC += $(wildcard boot/*.s)
C_OBJ := $(patsubst boot/%.c,$(BUILDIR)/%.c.o,$(C_SRC))
AS_OBJ := $(patsubst boot/%.s,$(BUILDIR)/%.s.o,$(AS_SRC))
BOOT_OBJ := $(AS_OBJ) $(C_OBJ)
BOOT_TARGET := $(BUILDIR)/bootloader

QEMU := qemu-system-arm
#QEMU := ~/qemu/build/arm-softmmu/qemu-system-arm

CFLAGS := -O0 -std=gnu99 -Wall -Wextra -fno-builtin -fno-pic -mfpu=neon-vfpv4 -mfloat-abi=hard \
          -march=armv7-a -mtune=cortex-a9 -ffreestanding -nostartfiles -Iboot/

LFLAGS := $(CFLAGS) -Tboot/link.ld


# Targets
.PHONY : all runtime $(GO_OBJ).elf

all: $(BOOT_TARGET).elf
	$(GCCPREFIX)objdump -D $(BOOT_TARGET).elf > $(BOOT_TARGET).dump
	$(GCCPREFIX)readelf -a $(BOOT_TARGET).elf > $(BOOT_TARGET).read

$(GO_OBJ).elf: $(GO_SRC)
	@echo "+GO "$(GO_SRC)
	@$(GOCMD) build -ldflags $(GOLINKFLAGS) -o $@
	$(GCCPREFIX)objdump -D $(GO_OBJ).elf > $(GO_OBJ).dump
	$(GCCPREFIX)readelf -a $(GO_OBJ).elf > $(GO_OBJ).read
	$(GCCPREFIX)objcopy -O binary $(GO_OBJ).elf $(GO_OBJ).bin
	./hackbin.py $(GO_OBJ).bin $(GO_OBJ)_hacked.bin
	hexdump $(GO_OBJ).bin > $(GO_OBJ).hex
	hexdump $(GO_OBJ)_hacked.bin > $(GO_OBJ)_hacked.hex

#linker
$(BOOT_TARGET).elf: $(GO_OBJ).elf $(BOOT_OBJ)
	@echo "linking " $(BOOT_OBJ)
	@$(GCCPREFIX)gcc $(LFLAGS) -o $@ $(BOOT_OBJ)

#build a c file
$(BUILDIR)/%.c.o : boot/%.c
	@mkdir -p $(BUILDIR)
	@echo "+CC "$< " -> " $@
	@$(GCCPREFIX)gcc $(CFLAGS) -c $< -o $@

#build an asm file
$(BUILDIR)/%.s.o : boot/%.s
	@mkdir -p $(BUILDIR)
	@echo "+AS "$< " -> " $@
	@$(GCCPREFIX)gcc $(CFLAGS) -c $< -o $@

clean:
	-rm -rf $(BUILDIR)

gdb-trace:
	$(GCCPREFIX)gdb -x gdb-trace.script > output
jtag0:
	JLinkGDBServer -ScriptFile core0.jlink -device MCIMX6Q5 -port 3330
jtag1:
	JLinkGDBServer -ScriptFile core1.jlink -device MCIMX6Q5 -port 3331
jtag2:
	JLinkGDBServer -ScriptFile core2.jlink -device MCIMX6Q5 -port 3332
jtag3:
	JLinkGDBServer -ScriptFile core3.jlink -device MCIMX6Q5 -port 3333
gdbj:
	$(GCCPREFIX)gdb -x gdb-jtag.script
qemud: all
	#$(QEMU) -machine sabrelite -m 3G -kernel $(BOOT_TARGET).elf -nographic -s -S
	$(QEMU) -machine sabrelite -cpu cortex-a9 -smp cpus=4 -m 3G -kernel $(BOOT_TARGET).elf -nographic -s -S
qemu: all
	$(QEMU) -machine sabrelite -cpu cortex-a9 -smp cpus=4 -m 3G -kernel $(BOOT_TARGET).elf -nographic
runtime:
	cd ../../golang1.7/src && ./rebuild_arm.sh
runtimeq:
	cd ../src && ./rebuild_arm_quick.sh
bytecode: test.s
	$(GCCPREFIX)gcc $(CFLAGS) -c test.s -o test.o
	$(GCCPREFIX)objdump -d -S test.o
sdcard:
	sudo umount $(SDCARD)1
	sudo mkdosfs $(SDCARD)1 -F 32
	sudo dd if=./u-boot/SPL of=$(SDCARD) seek=1 bs=1k
	sudo dd if=./u-boot/u-boot.img of=$(SDCARD) seek=69 bs=1k
	sudo mount $(SDCARD)1 /mnt/usb
	sudo cp ./u-boot/uEnv.txt /mnt/usb/
	sudo cp $(BOOT_TARGET).elf /mnt/usb/kernel2.elf
	sudo mkdir -p /mnt/usb/boot
	sudo cp ./u-boot/uEnv.txt /mnt/usb/boot/
	sudo umount /mnt/usb
	sync
uboot:
	cd $(shell pwd)/u-boot && make ARCH=arm CROSS_COMPILE=$(GCCPREFIX) distclean && make ARCH=arm CROSS_COMPILE=$(GCCPREFIX) wandboard_defconfig && make ARCH=arm CROSS_COMPILE=$(GCCPREFIX)
print-%: ; @echo $*=$($*)
