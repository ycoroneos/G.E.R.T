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
	fmt.Println("hiii from ", runtime.Cpunum(), "\n")
	resp <- "done"
}

func gcdone() {
	fmt.Println("stub")
	return
}

func main() {
	runtime.Release()
	fmt.Println("hi from fmt")
	channel := make(chan string, 1)
	channel <- "channel test pass"
	val := <-channel
	fmt.Println(val)
	for i := 0; i < 10; i++ {
		go printer(channel)
	}
	for i := 0; i < 10; i++ {
		<-channel
	}
	fmt.Println("start GC")
	runtime.GC()
	fmt.Println("done GC, sleeping forever")
	runtime.Crash = true
	gcdone()
	for {
	}
}
