package main

const (
	SYS_BASE = 0x0

	SYS_exit              = (SYS_BASE + 1)
	SYS_read              = (SYS_BASE + 3)
	SYS_write             = (SYS_BASE + 4)
	SYS_open              = (SYS_BASE + 5)
	SYS_close             = (SYS_BASE + 6)
	SYS_getpid            = (SYS_BASE + 20)
	SYS_kill              = (SYS_BASE + 37)
	SYS_gettimeofday      = (SYS_BASE + 78)
	SYS_clone             = (SYS_BASE + 120)
	SYS_rt_sigreturn      = (SYS_BASE + 173)
	SYS_rt_sigaction      = (SYS_BASE + 174)
	SYS_rt_sigprocmask    = (SYS_BASE + 175)
	SYS_sigaltstack       = (SYS_BASE + 186)
	SYS_mmap2             = (SYS_BASE + 192)
	SYS_futex             = (SYS_BASE + 240)
	SYS_exit_group        = (SYS_BASE + 248)
	SYS_munmap            = (SYS_BASE + 91)
	SYS_madvise           = (SYS_BASE + 220)
	SYS_setitimer         = (SYS_BASE + 104)
	SYS_mincore           = (SYS_BASE + 219)
	SYS_gettid            = (SYS_BASE + 224)
	SYS_tkill             = (SYS_BASE + 238)
	SYS_sched_yield       = (SYS_BASE + 158)
	SYS_select            = (SYS_BASE + 142) // newselect
	SYS_ugetrlimit        = (SYS_BASE + 191)
	SYS_sched_getaffinity = (SYS_BASE + 242)
	SYS_clock_gettime     = (SYS_BASE + 263)
	SYS_epoll_create      = (SYS_BASE + 250)
	SYS_epoll_ctl         = (SYS_BASE + 251)
	SYS_epoll_wait        = (SYS_BASE + 252)
	SYS_epoll_create1     = (SYS_BASE + 357)
	SYS_fcntl             = (SYS_BASE + 55)
	SYS_access            = (SYS_BASE + 33)
	SYS_connect           = (SYS_BASE + 283)
	SYS_socket            = (SYS_BASE + 281)
)

//go:nosplit
func delay(ticks uint32) {
	for ; ticks > 0; ticks-- {
	}
}

type Trapfunc func(int32, int32, int32, int32, int32) int32

var f Trapfunc = trap

//go:nosplit
func trap(callnum int32, arg1 int32, arg2 int32, arg3 int32, arg4 int32) int32 {
	uart_putc('!')
	switch callnum {
	case SYS_read:
		return 0
	case SYS_write:
		straddr := uint32(arg1)
		len := arg2
		return uart_print(straddr, len)
	default:
		print("unimplemented syscall %d\n", callnum)
	}
	return 0
}
