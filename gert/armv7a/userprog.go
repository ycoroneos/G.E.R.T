package main

import (
	"./embedded"
	"fmt"
	"time"
)

var c1, c2, c3, c4 uint32

func user_init() {

	embedded.WB_JP4_6.SetInput()
	embedded.WB_JP4_6.EnableIntr(embedded.INTR_RISING, nil)
	embedded.Enable_interrupt(103, 0) //send GPIO3:27 interrupt to CPU0

	embedded.WB_JP4_8.SetInput()
	embedded.WB_JP4_8.EnableIntr(embedded.INTR_RISING, nil)
	embedded.Enable_interrupt(109, 1) //send GPIO6:31 interrupt to CPU1

	embedded.WB_JP4_10.SetInput()
	embedded.WB_JP4_10.EnableIntr(embedded.INTR_RISING, nil)
	embedded.Enable_interrupt(99, 2) //send GPIO1:24 interrupt to CPU2

	embedded.WB_JP4_12.SetInput()
	embedded.WB_JP4_12.EnableIntr(embedded.INTR_RISING, nil)
	embedded.Enable_interrupt(110, 3) //send GPIO1:24 interrupt to CPU3
}

var state bool

func user_loop() {
	time.Sleep(3 * time.Second)
	fmt.Printf("count : %v\n", c1+c2+c3+c4)
}
