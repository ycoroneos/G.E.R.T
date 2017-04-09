package main

//import "./embedded"

/*
* This is the interrupt handler for all IRQs in GERT.
* It is extremely important that nothing in this code causes the scheduler to run
* or trigger a garbage collection. This is because this code may run while locks are held
* or even when the garbage collector is running too. This runs with CPSR=IRQ mode.
*
* To use this code, just modify the switch statement to look out for your IRQ
 */

//go:nosplit
//go:nowritebarrierec
func irq(irqnum uint32) {
	switch irqnum {
	default:
		//fmt.Printf("IRQ %d on cpu %d\n", irqnum, runtime.Cpunum())
	}
}
