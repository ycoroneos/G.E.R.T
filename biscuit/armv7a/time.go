package main

import "runtime"
import "fmt"

type Counter struct {
	clock uint32
}

var time Counter

func addtime(amt uint32) {
	time.clock += amt
}
func gettime() uint32 {
	curtime := time.clock
	return curtime
}

func sleep(sleeptime uint32) uint32 {
	curtime := gettime()
	for gettime()-curtime < sleeptime {
		runtime.Gosched()
	}
	return gettime()
}

func busysleep(sleeptime uint32) uint32 {
	curtime := gettime()
	for gettime()-curtime < sleeptime {
	}
	return gettime()
}

func gopherwatch() {
	for {
		time := int(sleep(2))
		fmt.Printf("time is %x\r\n", time)
		fmt.Printf("last irq from %d\n", <-irqchan)
	}
}
