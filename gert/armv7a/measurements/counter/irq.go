package main

import "./embedded"

/*
* This is the interrupt handler for all IRQs in GERT.
* It is extremely important that nothing in this code causes the scheduler to run
* or trigger a garbage collection. This is because this code may run while locks are held
* or even when the garbage collector is running too. This runs with CPSR=IRQ mode.
*
* To amend this code, just modify the switch statement to look out for your IRQ
 */

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
		embedded.Addtime(1)
		embedded.ClearGPTIntr()
	case 103:
		//embedded.Addtime(1)
		//embedded.GPIO_ISR(irqnum)
		inc()
		embedded.ClearIntr(3)
	default:
		//fmt.Printf("IRQ %d on cpu %d\n", irqnum, runtime.Cpunum())
	}
}
