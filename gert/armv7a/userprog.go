package main

import (
	"./embedded"
	"fmt"
)

func user_init() {

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
	//embedded.WB_JP4_4.SetOutput()
	//embedded.WB_JP4_4.SetLO()
	//embedded.WB_JP4_6.SetInput()
	//embedded.WB_JP4_6.EnableIntr(embedded.INTR_FALLING, inc)
	//embedded.Enable_interrupt(103, 0) //send GPIO3 interrupt to CPU0

	embedded.WB_PWM1.Begin(0x10)
	embedded.WB_PWM1.SetDuty(0.5)

	embedded.WB_PWM2.Begin(0xF000)
	embedded.WB_PWM2.SetDuty(0.5)

	embedded.WB_PWM3.Begin(0xFF00)
	embedded.WB_PWM3.SetDuty(0.5)

	//send the GPT interrupt to CPU1
	//embedded.Enable_interrupt(87, 1)

	//start the GPT
	//embedded.StartGPT()

	//start a little watch
	//go embedded.Gopherwatch()
}

func user_loop() {
	fmt.Printf("waiting for input: ")
	data := string(embedded.WB_DEFAULT_UART.Read(10)[:])
	fmt.Printf("got %s\n", data)
	//embedded.WB_JP4_6.SetLO()
	//embedded.WB_JP4_6.SetHI()
	//	embedded.WB_JP4_4.SetLO()
	//	embedded.WB_JP4_4.SetHI()
	//	if ping {
	//		fmt.Printf("count is %d\n", count)
	//		ping = false
	//	}
	//embedded.Sleep(2)
	//fmt.Printf("count is %d\n", count)
}

//var ping bool
var count uint32

//go:nosplit
func inc() {
	count += 1
	//	if count%5 == 0 {
	//		ping = true
	//	}
}
