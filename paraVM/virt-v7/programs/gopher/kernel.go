// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"../../embedded"
	"fmt"
	"runtime"
	"syscall"
)

/*
* This is the entry point of GERT. dont try anything fancy
 */

//go:nosplit
func Entry() {
	runtime.Armhackmode = 1
	runtime.Runtime_main()
}

//the runtime calls main after it's done setting up
func main() {
	//test things like channels and whatnot
	fmt.Printf("self tests ...")
	self_tests()
	fmt.Printf("done!\n")

	//print out some warnings for myself so I dont forget possibly sketchy things I have done
	fmt.Printf("warnings ...")
	self_warnings()
	fmt.Printf("done!\n")

	//init the GIC and turn on interrupts
	fmt.Printf("pre-init ...")
	pre_init()
	syscall.Setenv("TZ", "UTC")
	runtime.Booted = 1
	fmt.Printf("done!\n")

	//user-provided init code
	fmt.Printf("user init ...")
	user_init()
	fmt.Printf("done!\n")

	//user main loop
	for {
		user_loop()
	}
	panic("user loop broke out")
}

//add things here if you think they are critical for functionality
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

//I never read the git logs. Now I dont have to
func self_warnings() {
	//fmt.Println("REMEMBER THAT SKETCHY THING YOU DID WITH MAPPING AN EXTRA PAGE IN MAP_REGION")
}

//If a user doesnt want IRQs then they should never enable one. The GIC will just be ON but do nothing
func pre_init() {
	//enable GIC
	embedded.GIC_init(false)

	//set IRQ callback function
	runtime.SetIRQcallback(irq)

	//Release spinning cpus
	runtime.Release(3)

	//unmap the first page
	runtime.Unmap_region(0x0, 0x0, 0x100000)
}
