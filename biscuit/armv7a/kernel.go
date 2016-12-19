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
	//*(*uint32)(unsafe.Pointer(uintptr(UART1_UTXD))) = 97
	//uart_putc(97)
	//runtime.SWIcall()
	//runtime.Traphandle = unsafe.Pointer(&f)
	runtime.Runtime_main()
	//SWI()
	for {
	}
}

//go:nosplit
func main() {
	fmt.Printf("hello, world\n")
	for {
		uart_print([]byte("heloooo\n"))
	}
	SWI()
}
