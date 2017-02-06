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
	fmt.Println("REMEMBER THAT SKETCHY THING YOU DID WITH MAPPING AN EXTRA PAGE IN MAP_REGION")
}

//If a user doesnt want IRQs then they should never enable one. The GIC will just be ON but do nothing
func pre_init() {
	//enable GIC
	embedded.GIC_init(false)

	//set IRQ callback function
	runtime.SetIRQcallback(irq)

	//Release spinning cpus
	runtime.Release()
}
