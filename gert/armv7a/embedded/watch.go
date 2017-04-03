package embedded

import "runtime"
import "fmt"

type Counter struct {
	clock uint32
}

var timer Counter

func Addtime(amt uint32) {
	timer.clock += amt
}
func Gettime() uint32 {
	curtime := timer.clock
	return curtime
}

func Sleep(sleeptime uint32) uint32 {
	curtime := Gettime()
	for Gettime()-curtime < sleeptime {
		runtime.Gosched()
	}
	return Gettime()
}

func Busysleep(sleeptime uint32) uint32 {
	curtime := Gettime()
	for Gettime()-curtime < sleeptime {
	}
	return Gettime()
}

func Gopherwatch() {
	for {
		time := int(Sleep(2))
		fmt.Printf("time is %x\r\n", time)
		//		fmt.Printf("last irq from %d\n", <-irqchan)
	}
}
