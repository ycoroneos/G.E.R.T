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

var irqchan chan int = make(chan int, 20)

//go:nosplit
//go:nowritebarrierec
func irq(irqnum uint32) {
	//	//irqnum := gic_cpu.interrupt_acknowledge_register
	//	if len(irqchan) == cap(irqchan) {
	//		// Channel was full, but might not be by now
	//	} else {
	//		// Channel wasn't full, but might be by now
	//		irqchan <- runtime.Cpunum()
	//	}
	switch irqnum {
	case 87:
		addtime(1)
		gpt.SR = 0x1
	default:
		//fmt.Printf("IRQ %d on cpu %d\n", irqnum, runtime.Cpunum())
	}
	//gic_cpu.end_of_interrupt_register = irqnum
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
	enable_interrupt(87, 2)
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
	//fmt.Println("waiting")
	//fmt.Println(makesgi())
	//fmt.Println("got it")
	//fmt.Println(startGPT())
	//go gopherwatch()
	//	for i := 0; i < 10000000; i++ {
	//		fmt.Println(i)
	//	}
	//fir_main()
	//fmt.Println("done with fir test")
	fmt.Println("try to init sd card 1")
	if card_init(3, 4) < 0 {
		fmt.Println("init sd card failure")
	} else {
		fmt.Println("done init sd card")
		data := make([]uint32, 25, 25)
		if card_data_read(uint32(3), &data, 4, 512) > 0 {
			for i := 0; i < 25; i++ {
				fmt.Printf("byte %x\n", data[i])
			}
		}
	}
	for {
		//fmt.Println(<-irqchan)
	}
}
