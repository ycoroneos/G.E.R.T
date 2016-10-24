package runtime

func trap_debug(arg1, arg2, arg3, arg4, arg5, arg7, trapno uint32) {
	print("unhandled trap: ", hex(trapno))
	throw("panic")
}
