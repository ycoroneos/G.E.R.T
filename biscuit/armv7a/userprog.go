package main

import (
	"./embedded"
	"fmt"
)

func user_init() {
	//send the GPT interrupt to CPU2
	embedded.Enable_interrupt(87, 2)

	//start the GPT
	embedded.StartGPT()

	//start a little watch
	go embedded.Gopherwatch()

	//play with the SD card a bit
	good, root := embedded.Fat32_som_start(embedded.Init_som_sdcard, embedded.Read_som_sdcard)
	if !good {
		fmt.Println("fat32 init failure")
	}
	fmt.Println(root.Getfilenames())
	fmt.Println(root.Getsubdirnames())
	good, bootdir := root.Direnter("BOOT")
	if !good {
		panic("dir entry failed")
	} else {
		fmt.Println(bootdir.Getfilenames())
		good, contents := bootdir.Fileread("UENV.TXT")
		if !good {
			panic("file read failure")
		}
		fmt.Println(string(contents))
	}
	embedded.WB_JP4_4.SetOutput()
}

var toggle = uint8(0)

func user_loop() {
	embedded.WB_JP4_4.Write(toggle)
	toggle = ^toggle
}
