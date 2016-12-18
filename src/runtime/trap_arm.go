package runtime

import "unsafe"

//go:nosplit
func PutR0(val uint32)

//go:nosplit
func RR0() uint32

//go:nosplit
func RR1() uint32

//go:nosplit
func RR2() uint32

//go:nosplit
func RR3() uint32

//go:nosplit
func RR4() uint32

//go:nosplit
func RR5() uint32

//go:nosplit
func RR6() uint32

//go:nosplit
func RR7() uint32

var firstexit = true

//go:nosplit
func trap_debug() {
	arg0 := RR0()
	arg1 := RR1()
	arg2 := RR2()
	arg3 := RR3()
	arg4 := RR4()
	arg5 := RR5()
	arg6 := RR6()
	trapno := RR7()
	switch trapno {
	case 120:
		print("spoofing clone\n")
		thread_id := makethread(uint32(arg0), uintptr(arg1), uintptr(arg2))
		PutR0(uint32(thread_id))
		return
	case 142:
		print("spoofing select\n")
		PutR0(0)
		return
	case 174:
		print("spoofing rtsigproc\n")
		PutR0(0)
		return
	case 175:
		print("spoofing rtsigaction\n")
		PutR0(0)
		return
	case 186:
		print("spoofing sigaltstack\n")
		PutR0(0)
		return
	case 224:
		print("gettid\n")
		PutR0(thread_current())
		return
	case 238:
		print("spoofing tkill\n")
		PutR0(0)
		return
	case 240:
		print("spoofing futex\n")
		uaddr := ((*int32)(unsafe.Pointer(uintptr(arg0))))
		ts := ((*timespec)(unsafe.Pointer(uintptr(arg3))))
		uaddr2 := ((*int32)(unsafe.Pointer(uintptr(arg4))))
		ret := hack_futex_arm(uaddr, int32(arg1), int32(arg2), ts, uaddr2, int32(arg5))
		PutR0(uint32(ret))
		return
	case 248:
		if firstexit == true {
			firstexit = false
			throw("exit")
		}
		for {
		}
	}
	print("unpatched trap: ", trapno, "\n")
	print("\tr0: ", hex(arg0), "\n")
	print("\tr1: ", hex(arg1), "\n")
	print("\tr2: ", hex(arg2), "\n")
	print("\tr3: ", hex(arg3), "\n")
	print("\tr4: ", hex(arg4), "\n")
	print("\tr5: ", hex(arg5), "\n")
	print("\tr6: ", hex(arg6), "\n")
	throw("trap")
}
