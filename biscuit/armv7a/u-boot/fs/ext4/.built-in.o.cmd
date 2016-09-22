cmd_fs/ext4/built-in.o :=  arm-none-eabi-ld.bfd     -r -o fs/ext4/built-in.o fs/ext4/ext4fs.o fs/ext4/ext4_common.o fs/ext4/dev.o fs/ext4/ext4_write.o fs/ext4/ext4_journal.o fs/ext4/crc16.o 
