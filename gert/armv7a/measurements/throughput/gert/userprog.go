package main

import (
	"./embedded"
	"fmt"
	"time"
)

var count1, count2, count3, count4 uint32

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
	//embedded.WB_JP4_6.SetOutput()
	count1 = 0
	count2 = 0
	count3 = 0
	count4 = 0
	//	ping = false
	embedded.WB_JP4_6.SetInput()
	embedded.WB_JP4_6.EnableIntr(embedded.INTR_RISING, inc)
	embedded.WB_JP4_8.SetInput()
	embedded.WB_JP4_8.EnableIntr(embedded.INTR_RISING, inc)
	embedded.WB_JP4_10.SetInput()
	embedded.WB_JP4_10.EnableIntr(embedded.INTR_RISING, inc)
	embedded.WB_JP4_12.SetInput()
	embedded.WB_JP4_12.EnableIntr(embedded.INTR_RISING, inc)
	embedded.Enable_interrupt_mask(103, embedded.GICcpumask([]uint8[0, 1, 2, 3]), 0) //send GPIO3 interrupt to CPU0
	//	embedded.Enable_interrupt(109, 1, 1) //send GPIO6 interrupt to CPU1
	//	embedded.Enable_interrupt(99, 2, 2)  //send GPIO1 interrupt to CPU2
	//	embedded.Enable_interrupt(110, 3, 3) //send GPIO7 interrupt to CPU3

	//send the GPT interrupt to CPU1
	//embedded.Enable_interrupt(87, 1)

	//start the GPT
	//embedded.StartGPT()

	//start a little watch
	//go embedded.Gopherwatch()
	fmt.Println("about to loop")
}

var oldcount uint32

func user_loop() {
	//	embedded.WB_JP4_6.SetLO()
	//	embedded.WB_JP4_6.SetHI()
	//	if ping {
	//		fmt.Printf("count is %d\n", count)
	//		ping = false
	//	}
	total := count1 + count2 + count3 + count4
	diff := total - oldcount
	fmt.Printf("count is %d\n", diff)
	oldcount = total
	time.Sleep(1 * time.Second)
}

//go:nosplit
func inc() {
	//	count += 1
	//	//	if count%5 == 0 {
	//	//		ping = true
	//	//	}
}
