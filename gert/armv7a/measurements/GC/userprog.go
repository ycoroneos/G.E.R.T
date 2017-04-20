package main

import "fmt"

func user_init() {
}

func user_loop() {
	for i := 0; i < 100000; i++ {
		fmt.Printf("iter %v\n", i)
		junk := make([]int, 100000)
		junk[0] = 10
		junk[9999] = 10
	}

}

//import (
//	"./embedded"
//	"fmt"
//	"runtime"
//	//	"math"
//	//	"time"
//	//"unsafe"
//)
//
//var event_chan chan interface{}
//var drive *embedded.MDD10A_controller
//var adc *embedded.MCP3008_controller
//
//func user_init() {
//
//	//play with the SD card a bit
//	//	good, root := embedded.Fat32_som_start(embedded.Init_som_sdcard, embedded.Read_som_sdcard)
//	//	if !good {
//	//		fmt.Println("fat32 init failure")
//	//	}
//	//	fmt.Println(root.Getfilenames())
//	//	fmt.Println(root.Getsubdirnames())
//	//	good, bootdir := root.Direnter("BOOT")
//	//	if !good {
//	//		panic("dir entry failed")
//	//	} else {
//	//		fmt.Println(bootdir.Getfilenames())
//	//		good, contents := bootdir.Fileread("UENV.TXT")
//	//		if !good {
//	//			panic("file read failure")
//	//		}
//	//		fmt.Println(string(contents))
//	//	}
//	adc = embedded.MakeMCP3008(embedded.WB_SPI1)
//	drive = embedded.MakeMDD10A(embedded.WB_PWM1, embedded.WB_PWM2, embedded.WB_JP4_4, embedded.WB_JP4_6)
//	event_chan = make(chan interface{}, 10)
//	_ = embedded.Poll(func() interface{} {
//		return string(embedded.WB_DEFAULT_UART.Read(1)[:])
//	}, 0, event_chan)
//
//	//	_ = embedded.Poll(func() interface{} {
//	//		return adc.Read(0)
//	//	}, 2*time.Second, event_chan)
//	//
//	//	//fmt.Printf("pi is %v \n", pi(50))
//	//
//	//	go func() {
//	//		for {
//	//			old := count
//	//			time.Sleep(1 * time.Second)
//	//			new := count
//	//			event_chan <- new - old
//	//		}
//	//	}()
//
//	//embedded.WB_JP4_10.SetOutput()
//	//out := ((*uint32)(unsafe.Pointer(uintptr(0x209C000))))
//	//	for {
//	//		//embedded.WB_JP4_10.SetHInow()
//	//		//embedded.WB_JP4_10.SetLOnow()
//	//		//*out = 0xFFFF
//	//		//*out = 0xFFFF0000
//	//		embedded.WB_JP4_10.SetHI()
//	//		embedded.WB_JP4_10.SetLO()
//	//		//embedded.Set(unsafe.Pointer(uintptr(0x209C000)), uint32(0x0))
//	//		//embedded.Set(unsafe.Pointer(uintptr(0x209C000)), uint32(0xFFFFFFFF))
//	//	}
//
//	//embedded.WB_JP4_10.SetInput()
//	//embedded.WB_JP4_10.EnableIntr(embedded.INTR_FALLING, inc)
//	//embedded.Enable_interrupt(99, 0) //send GPIO1 interrupt to CPU0
//	runtime.GC()
//}
//
//func user_loop() {
//	//make an event loop
//	select {
//	case event := <-event_chan:
//		fmt.Printf("%v\n", event)
//		//	switch event {
//		//	case "p":
//		//		val := adc.Read(0)
//		//		fmt.Printf("adc reads %v\n", val)
//		//	case "w":
//		//		drive.Forward(0.2)
//		//	case "s":
//		//		drive.Backward(0.2)
//		//	case "a":
//		//		drive.TurnRight(0.2)
//		//	case "d":
//		//		drive.TurnLeft(0.2)
//		//	case " ":
//		//		drive.Stop()
//		//	}
//	}
//
//}
//
//var count uint32
//
//func inc() {
//	count += 1
//}
//
//// pi launches n goroutines to compute an
//// approximation of pi.
////func pi(n int) float64 {
////	ch := make(chan float64)
////	for k := 0; k <= n; k++ {
////		go term(ch, float64(k))
////	}
////	f := 0.0
////	for k := 0; k <= n; k++ {
////		f += <-ch
////	}
////	return f
////}
////
////func term(ch chan float64, k float64) {
////	ch <- 4 * math.Pow(-1, k) / (2*k + 1)
////}
