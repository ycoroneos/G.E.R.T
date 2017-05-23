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

var event_chan chan interface{}
var laser_chan chan LaserCommand

func user_init() {
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
	//	good, contents := bootdir.Fileread("P.TXT")
	//	if !good {
	//		panic("file read failure")
	//	}
	var points []CompactPoint
	contents, err := Asset("bindata.gob")
	if err != nil {
		panic("bindata not found")
	}
	r := bytes.NewBuffer(contents)
	d := gob.NewDecoder(r)
	err = d.Decode(&points)
	if err != nil {
		fmt.Printf("error de-GOBing:\n")
		panic(err)
	}
	fmt.Printf("%v points", len(points))
	laser_chan = make(chan LaserCommand, 10)
	//laser_chan <- LaserCommand{WAITTIME, 50}
	go lasermon(laser_chan, points)
}

func lasermon(commands chan LaserCommand, points []CompactPoint) {
	fmt.Printf("Hi from lasermon!\n")
	wait := 0 * time.Microsecond
	curpoint := 0
	skipcount := 0
	lastx := uint16(0)
	lasty := uint16(0)
	dac := embedded.MakeMCP4922(embedded.WB_SPI1)
	for {
		select {
		case command := <-commands:
			switch command.Type {
			case WAITTIME:
				wait = time.Duration(command.Value) * time.Microsecond
			}
		default:
			if skipcount == 0 {
				if curpoint >= len(points) {
					curpoint = 0
				}
				xd := uint16(0)
				yd := uint16(0)
				if lastx > points[curpoint].X {
					xd = lastx - points[curpoint].X
				} else {
					xd = points[curpoint].X - lastx
				}
				if lasty > points[curpoint].Y {
					yd = lasty - points[curpoint].Y
				} else {
					yd = points[curpoint].Y - lasty
				}
				skipcount = int((xd + yd)) / 8
				dac.Write(points[curpoint].X, 0)
				dac.Write(points[curpoint].Y, 1)
				lastx = points[curpoint].X
				lasty = points[curpoint].Y
				if wait > 0 {
					time.Sleep(wait)
				}
				curpoint += 1
			} else {
				skipcount -= 1
			}
		}
	}
}

func user_loop() {
}
