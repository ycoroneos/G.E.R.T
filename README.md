# G.E.R.T : Golang Embedded Run-Time

GERT is a modified version of Go that runs bare-metal on armv7a SOCs. The minimal
set of OS primitives that Go relies on have been re-implemented entirely in Go
and Plan 9 assembly inside the modified runtime. The goal of this project is to bring
the benefits of a high-level, type-safe, and garbage-collected language to bare-metal
embedded environments. GERT has been developed for the Wandboard Quad (iMX6 Quad SOC), but
GERT can be easily ported to any armv7a SOC with adequate documentation.

## Index


## Quickstart

### Materials Needed

GERT can either run in the QEMU emulator or on real hardware. To emulate GERT follow the install instructions
below, because they include a QEMU installation. To run GERT on real hardware,
you will need to get a **Wandboard Quad** and an SD card, or any other dev kit which uses the Freescale
iMX6 Quad SOC. GERT only works with the iMX6 right now because its memory map is hard-coded into
the kernel.


### Directory Layout
  |Directory | Function |
  |----------|----------|
  |`golang_embedded`| Modified Go runtime which runs on bare-metal. It has its own repo which explains modifications in detail|
  |`qemu`| QEMU git master branch. submoduled|
  |`thesis`| My master's thesis|
  |`gert/armv7a`| Contains the user-facing code for running GERT on the Wandboard Quad dev board|
  |`gert/armv7a/uboot`| U-boot bootloader which configures basic device clocks and loads GERT off the sd card|
  |`gert/armv7a/boot`|  Second-stage bootloader written in C that prepares the initial Go stack|
  |`gert/armv7a/embedded`| Go package which contains many drivers for the iMX6 Quad cpu|
  |`gert/armv7a/doc`| Technical reference manuals on the armv7a, cortex-a9 mpcore architectures, and iMX6 SOC|
  |`gert/armv7a/programs`| A storage directory for some GERT programs|
  |`gert/armv7a/measurements`| Benchmarks and measurements I took for my thesis|
  |`gert/armv7a/debug`| JTAG scripts for the JLink EDU|


### Installation

#### Ubuntu
  <!-- language: lang-none -->

     sudo apt install gcc-arm-none-eabi arm-none-eabi-gdb golang git
     sudo apt-get build-dep golang qemu
     git clone git@github.com:ycoroneos/G.E.R.T.git
     cd G.E.R.T
     git submodule init
     git submodule update --recursive
     cd qemu && git submodule init && git submodule update --recursive && ./configure --target-list=arm-softmmu && make -j4 && cd ..
     cd gert/armv7a && make runtime && make uboot && UPROG=programs/hello make && make qemu

If all went well, you should be running the 'hello' program in QEMU

### Programming With GERT

GERT programs live inside the `gert/armv7a/programs` directory. Each program folder
must contain three things: `kernel.go`, `irq.go`, and `userprog.go`. Take a look at
programs/hello to get a feel for the layout.

#### kernel.go

This is not the actual GERT kernel, but it contains the GERT entry point,
initialization for the interrupt controller, and code to enable all 4 cpus
on the iMX6. It's boilerplate code you would have written anyway.

#### irq.go

Contains the interrupt service routine that GERT executes when a cpu
gets an interrupt (switches to ISR mode). Every cpu can concurrently execute
in the interrupt handler; interrupts are not serialized. There are some basic rules
you must obey inside the interrupt handler though : **no blocking operations
and no allocations on the heap**. This is because the garbage collector might be running
while an interrupt is being serviced. The *irqnum* input is the ID of the SPI
that was received. You must enable the specific interrupts you want to receive
in the GIC before the interrupt handler will ever execute.

#### userprog.go

This contains your GERT program. There are at least two functions you must implement:
*user_init*, which is called once, and *user_loop*, which is called repeatedly. Besides that,
Go crazy. Most of the standard libraries work as well as channels and goroutines. The *embedded* package
contains drivers for many iMX6 peripherals for you to play around with.


### Working With GERT

`gert/armv7a` is the working directory so you should execute all commands from in there. Everytime you change the runtime run `make runtime`.
To build a GERT program, run `UPROG=<your prog dir> make`. To run your GERT program in QEMU do `make qemu`.
To put your GERT program on an sd card and boot it run `SDCARD=/dev/<your sdcard> make sdcard`


## 'Hello World!'

Lets make hello world in GERT and run it on bare-metal. This is exactly what your quad-core SOC was made for!
1. Navigate to `gert/armv7a` because this is the working directory
2. Copy `programs/blank` to `programs/hello`

## Debugging

Many bugs in a GERT program result in a Go panic, which prints a very useful backtrace along with a useful
error message. This is usually enough to fix the problem.

### QEMU

Try to reproduce your error in QEMU. It will significantly reduce debugging time.
Run your faulty GERT program with `make qemud` and connect to it with `gdb-arm-none-eabi`. GERT programs
are compiled with debugging symbols and Go has excellent support for GDB.


### JTAG

QEMU does not emulate most of the iMX6 peripherals or cpu very well. You may need to use a JTAG adaptor
to step single instructions on the actual hardware. There are JTAG scripts for attaching to specific
Cortex A9 cores inside the `gert/armv7a/debug` directory. These scripts each setup a gdb server that
you can connect to with `gdb-arm-none-eabi` again. These scripts work for the Segger JLink
EDU. They can probably work with other Segger products but not with OpenOCD.

