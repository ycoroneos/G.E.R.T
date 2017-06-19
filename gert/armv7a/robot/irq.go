// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	switch irqnum {
	case 87:
		embedded.Addtime(1)
		embedded.ClearGPTIntr()
	case 99:
		inc()
		embedded.ClearIntr(1)
	case 103:
		embedded.ClearIntr(3)
	default:
		//fmt.Printf("IRQ %d on cpu %d\n", irqnum, runtime.Cpunum())
	}
}
