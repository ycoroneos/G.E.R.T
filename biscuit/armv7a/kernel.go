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

var irqchan chan int = make(chan int, 1)

//go:nosplit
//go:nowritebarrierec
func irq() {
	irqnum := gic_cpu.interrupt_acknowledge_register
	select {
	case irqchan <- runtime.Cpunum():
	default:
	}
	switch irqnum {
	case 87:
		addtime(1)
		gpt.SR = 0x1
	default:
		//fmt.Printf("IRQ %d on cpu %d\n", irqnum, runtime.Cpunum())
	}
	gic_cpu.end_of_interrupt_register = irqnum
}

func printer(resp chan string) {
	fmt.Println("hiii from cpu ", runtime.Cpunum(), "\n")
	resp <- "done"
}

func makesgi() int {
	sgi(0, 0xFF)
	return 1
}

//the runtime calls main after it's done setting up
func main() {
	GIC_init(false)
	runtime.SetIRQcallback(irq)
	runtime.Release()
	channel := make(chan string, 1)
	channel <- "channel test pass"
	val := <-channel
	fmt.Println(val)
	//	for i := 0; i < 20; i++ {
	//		go printer(channel)
	//	}
	//	for i := 0; i < 20; i++ {
	//		fmt.Println(<-channel)
	//	}
	fmt.Println("waiting")
	fmt.Println(makesgi())
	//fmt.Println("got it")
	fmt.Println(startGPT())
	go gopherwatch()
	//	for i := 0; i < 10000000; i++ {
	//		fmt.Println(i)
	//	}
	for {
		//fmt.Println(<-irqchan)
	}
}
