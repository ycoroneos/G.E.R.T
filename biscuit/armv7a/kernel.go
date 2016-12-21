package main

import "runtime"

//import "fmt"

//import "unsafe"

/*
* This is the entry point of the go kernel. dont try anything fancy
 */
func Getstack()
func SWI()

//go:nosplit
func Entry() {
	//*(*uint32)(unsafe.Pointer(uintptr(UART1_UTXD))) = 97
	//uart_putc(97)
	//runtime.SWIcall()
	//runtime.Traphandle = unsafe.Pointer(&f)
	runtime.Runtime_main()
	//main()
	//SWI()
	//for {
	//	}
}

//go:nosplit
func main() {
	for {
		uart_print([]byte("fmt print\n"))
		//fmt.Printf("hello, world\n")
		uart_print([]byte("fmt print done\n"))
	}
	SWI()
}
