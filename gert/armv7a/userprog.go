package main

import (
	"./embedded"
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

const (
	WAITTIME = 0
)

type LaserCommand struct {
	Type  int
	Value int
}

type CompactPoint struct {
	X     uint16
	Y     uint16
	Color uint8
}

var points []CompactPoint

var event_chan chan interface{}
var laser_chan chan LaserCommand

func user_init() {
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

	//read the points for the laser off an sdcard
	//good, root := embedded.Fat32_som_start(embedded.Init_som_sdcard, embedded.Read_som_sdcard)
	//if !good {
	//	fmt.Println("fat32 init failure")
	//}
	//fmt.Println(root.Getfilenames())
	//fmt.Println(root.Getsubdirnames())
	good, contents := bootdir.Fileread("P.TXT")
	if !good {
		panic("file read failure")
	}
	r := bytes.NewBuffer(contents)
	d := gob.NewDecoder(r)
	err := d.Decode(&points)
	if err != nil {
		fmt.Printf("error de-GOBing:\n")
		panic(err)
	}
	fmt.Printf("%v", points)
	//	good, bootdir := root.Direnter("LASER")
	//	if !good {
	//		panic("dir entry failed")
	//	} else {
	//		fmt.Println(bootdir.Getfilenames())
	//		//good, contents := bootdir.Fileread("POINTS.GOB")
	//		good, contents := bootdir.Fileread("testfile.txt")
	//		if !good {
	//			panic("file read failure")
	//		}
	//		r := bytes.NewBuffer(contents)
	//		d := gob.NewDecoder(r)
	//		err := d.Decode(&points)
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	laser_chan = make(chan LaserCommand, 10)
	go lasermon(laser_chan)
	//	adc = embedded.MakeMCP3008(embedded.WB_SPI1)
	//	drive = embedded.MakeMDD10A(embedded.WB_PWM1, embedded.WB_PWM2, embedded.WB_JP4_4, embedded.WB_JP4_6)
	//	event_chan = make(chan interface{}, 10)
	//	_ = embedded.Poll(func() interface{} {
	//		return string(embedded.WB_DEFAULT_UART.Read(1)[:])
	//	}, 0, event_chan)
	//
	//	_ = embedded.Poll(func() interface{} {
	//		return adc.Read(0)
	//	}, 2*time.Second, event_chan)
	//
	//	//fmt.Printf("pi is %v \n", pi(50))
	//
	//	go func() {
	//		for {
	//			old := count
	//			time.Sleep(1 * time.Second)
	//			new := count
	//			event_chan <- new - old
	//		}
	//	}()
	//
	//	embedded.WB_JP4_10.SetOutput()
	//	//out := ((*uint32)(unsafe.Pointer(uintptr(0x209C000))))
	//	for {
	//		//embedded.WB_JP4_10.SetHInow()
	//		//embedded.WB_JP4_10.SetLOnow()
	//		//*out = 0xFFFF
	//		//*out = 0xFFFF0000
	//		embedded.WB_JP4_10.SetHI()
	//		embedded.WB_JP4_10.SetLO()
	//		//embedded.Set(unsafe.Pointer(uintptr(0x209C000)), uint32(0x0))
	//		//embedded.Set(unsafe.Pointer(uintptr(0x209C000)), uint32(0xFFFFFFFF))
	//	}
	//
	//	embedded.WB_JP4_10.SetInput()
	//	embedded.WB_JP4_10.EnableIntr(embedded.INTR_FALLING, inc)
	//	embedded.Enable_interrupt(99, 0) //send GPIO1 interrupt to CPU0
}

func lasermon(commands chan LaserCommand) {
	fmt.Printf("Hi from lasermon!\n")
	wait := 0 * time.Microsecond
	curpoint := 0
	dac := embedded.MakeMCP4922(embedded.WB_SPI1)
	for {
		select {
		case command := <-commands:
			switch command.Type {
			case WAITTIME:
				wait = time.Duration(command.Value) * time.Microsecond
			}
		default:
			if curpoint > len(points) {
				curpoint = 0
			}
			dac.Write(points[curpoint].X, 0)
			dac.Write(points[curpoint].Y, 0)
			if wait > 0 {
				time.Sleep(wait)
			}
		}
	}
}

func user_loop() {
}
