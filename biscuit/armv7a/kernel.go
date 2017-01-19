package main

import (
	"./embedded"
	"fmt"
	"runtime"
)

/*
* This is the entry point of the go kernel. dont try anything fancy
 */

//go:nosplit
func Entry() {
	runtime.Armhackmode = 1
	runtime.Runtime_main()
}

//the runtime calls main after it's done setting up
func main() {
	//test things like channels and whatnot
	self_tests()

	//print out some warnings for myself so I dont forget possibly sketchy things I have done
	self_warnings()

	//init the GIC and turn on interrupts
	pre_init()

	//user-provided init code
	user_init()

	//user main loop
	for {
		user_loop()
	}
	panic("user loop broke out")
	//	embedded.GIC_init(false)
	//	runtime.SetIRQcallback(irq)
	//	runtime.Release()
	//	embedded.Enable_interrupt(87, 2)
	//	channel := make(chan string, 1)
	//	channel <- "channel test pass"
	//	val := <-channel
	//	fmt.Println(val)
	//	fmt.Println("REMEMBER THAT SKETCHY THING YOU DID WITH MAPPING AN EXTRA PAGE IN MAP_REGION")
	//	//	for i := 0; i < 20; i++ {
	//	//		go printer(channel)
	//	//	}
	//	//	for i := 0; i < 20; i++ {
	//	//		fmt.Println(<-channel)
	//	//	}
	//	//fmt.Println("waiting")
	//	//fmt.Println(makesgi())
	//	//fmt.Println("got it")
	//	fmt.Println(embedded.StartGPT())
	//	go embedded.Gopherwatch()
	//	//	for i := 0; i < 10000000; i++ {
	//	//		fmt.Println(i)
	//	//	}
	//	//	fmt.Println("try to init sd card 1")
	//	//	if !init_som_sdcard() {
	//	//		fmt.Println("init sd card failure")
	//	//	} else {
	//	//		fmt.Println("done init sd card")
	//	//		good, data := read_som_sdcard(8, 0x3000)
	//	//		if good {
	//	//			for i := 0; i < len(data); i++ {
	//	//				fmt.Printf("\tbyte read %x\n", data[i])
	//	//			}
	//	//		}
	//	//	}
	//	//	fir_main()
	//	//	fmt.Println("done with fir test")
	//	good, root := embedded.Fat32_som_start(embedded.Init_som_sdcard, embedded.Read_som_sdcard)
	//	if !good {
	//		fmt.Println("fat32 init failure")
	//	}
	//	fmt.Println(root.Getfilenames())
	//	fmt.Println(root.Getsubdirnames())
	//	good, bootdir := root.Direnter("BOOT")
	//	if !good {
	//		panic("dir entry failed")
	//	} else {
	//		fmt.Println(bootdir.Getfilenames())
	//		good, contents := bootdir.Fileread("UENV.TXT")
	//		if !good {
	//			panic("file read failure")
	//		}
	//		fmt.Println(string(contents))
	//	}
	//
	//	//do nothing
	//	for {
	//		//fmt.Println(<-irqchan)
	//	}
}

func self_tests() {
	fmt.Println("Hi from fmt")
	channel := make(chan string, 1)
	channel <- "channel test pass"
	val := <-channel
	fmt.Println(val)
	go func(resp chan string) {
		fmt.Println("print from inside goroutine")
		resp <- "send channel from inside a goroutine"
	}(channel)
	val = <-channel
	fmt.Println(val)
}

func self_warnings() {
	fmt.Println("REMEMBER THAT SKETCHY THING YOU DID WITH MAPPING AN EXTRA PAGE IN MAP_REGION")
}

func pre_init() {
	//enable GIC
	embedded.GIC_init(false)

	//set IRQ callback function
	runtime.SetIRQcallback(irq)

	//Release spinning cpus
	runtime.Release()
}
