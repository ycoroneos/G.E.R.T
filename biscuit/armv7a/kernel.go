package main

import "runtime"

import "fmt"

//import "unsafe"

/*
* This is the entry point of the go kernel. dont try anything fancy
 */
func Getstack()
func SWI()

//go:nosplit
func Entry() {
	runtime.Armhackmode = 1
	runtime.Runtime_main()
}

func printer(resp chan string) {
	fmt.Println("hiii from cpu ", runtime.Cpunum(), "\n")
	resp <- "done"
}

//the runtime calls main after it's done setting up
func main() {
	GIC_init(false)
	//runtime.Release()
	channel := make(chan string, 1)
	channel <- "channel test pass"
	val := <-channel
	fmt.Println(val)
	for i := 0; i < 20; i++ {
		go printer(channel)
	}
	for i := 0; i < 20; i++ {
		fmt.Println(<-channel)
	}
	fmt.Println(startGPT())
	fmt.Println("done now, wait for interrupt")
	//sgi(0, 0xFF)
	for {
	}
}
