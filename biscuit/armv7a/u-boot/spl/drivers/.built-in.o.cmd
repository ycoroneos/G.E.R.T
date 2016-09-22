cmd_spl/drivers/built-in.o :=  arm-none-eabi-ld.bfd     -r -o spl/drivers/built-in.o spl/drivers/i2c/built-in.o spl/drivers/gpio/built-in.o spl/drivers/mmc/built-in.o spl/drivers/serial/built-in.o 
