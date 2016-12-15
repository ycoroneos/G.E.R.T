package runtime

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

func trap_debug() {
	arg0 := RR0()
	arg1 := RR1()
	arg2 := RR2()
	arg3 := RR3()
	arg4 := RR4()
	arg5 := RR5()
	arg6 := RR6()
	trapno := RR7()
	print("unpatched trap: ", trapno, "\n")
	print("\targ0: ", hex(arg0), "\n")
	print("\targ1: ", hex(arg1), "\n")
	print("\targ2: ", hex(arg2), "\n")
	print("\targ3: ", hex(arg3), "\n")
	print("\targ4: ", hex(arg4), "\n")
	print("\targ5: ", hex(arg5), "\n")
	print("\targ6: ", hex(arg6), "\n")
	switch trapno {
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
		print("spoofing gettid\n")
		PutR0(0)
		return
	case 238:
		print("spoofing tkill\n")
		PutR0(0)
		return
	case 248:
		print("exit")
		for {
		}
	}
	for {
	}
	throwcritical("trap")
}
