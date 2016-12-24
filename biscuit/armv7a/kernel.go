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

//go:nosplit
func main() {
	fmt.Println("hi from fmt")
	channel := make(chan string, 1)
	channel <- "channel test pass"
	val := <-channel
	fmt.Println(val)
	go func(resp chan string) {
		fmt.Println("goprint from inside a go routine!")
		resp <- "done"
	}(channel)
	<-channel
	fmt.Println("start GC")
	runtime.GC()
	fmt.Println("done GC, sleeping forever")
	for {
	}
}
