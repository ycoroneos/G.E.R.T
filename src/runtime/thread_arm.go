package runtime

import "unsafe"

//go:nosplit
func Threadschedule()

//go:nosplit
func RecordTrapframe()

//go:nosplit
func ReplayTrapframe()

var cpunum = 0

const maxcpus = 4

type trapframe struct {
	lr  uintptr
	sp  uintptr
	fp  uintptr
	r0  uint32
	r1  uint32
	r2  uint32
	r3  uint32
	r10 uint32
}

type thread_t struct {
	tf      trapframe
	state   uint32
	futaddr uintptr
}

// maximum # of runtime "OS" threads
const maxthreads = 64

var threads [maxthreads]thread_t
var curthread *thread_t

// thread states
const (
	ST_INVALID   = 0
	ST_RUNNABLE  = 1
	ST_RUNNING   = 2
	ST_WAITING   = 3
	ST_SLEEPING  = 4
	ST_WILLSLEEP = 5
)

//go:nosplit
func thread_init() {
	threads[0].state = ST_RUNNING
	curthread = &threads[0]
	RecordTrapframe()
	print("made thread 0\n")
}

//go:nosplit
func makethread(flags uint32, stack uintptr, entry uintptr) int {
	CLONE_VM := 0x100
	CLONE_FS := 0x200
	CLONE_FILES := 0x400
	CLONE_SIGHAND := 0x800
	CLONE_THREAD := 0x10000
	chk := uint32(CLONE_VM | CLONE_FS | CLONE_FILES | CLONE_SIGHAND |
		CLONE_THREAD)
	if flags != chk {
		print("unexpected clone args ", hex(uintptr(flags)), " expected ", hex(chk))
		throw("clone error")
	}
	i := 0
	for ; i < maxthreads; i++ {
		if threads[i].state == ST_INVALID {
			break
		}
	}
	if i == maxthreads {
		throw("out of threads to use\n")
	}
	threads[i].state = ST_RUNNABLE
	threads[i].tf.lr = entry
	threads[i].tf.sp = stack
	threads[i].tf.r0 = uint32(i)
	print("\t\t\t\t new thread ", i, "\n")
	return 0
}

var lastrun = 0

//go:nosplit
func thread_schedule() {
	//RecordTrapframe()
	print("thread scheduler\n")
	for ; lastrun < maxthreads; lastrun = (lastrun + 1) % maxthreads {
		if threads[lastrun].state == ST_RUNNABLE {
			threads[lastrun].state = ST_RUNNING
			curthread = &threads[lastrun]
			print("\t\t\t\tschedule thread ", lastrun, "\n")
			print("\t\t\t\tLR ", hex(curthread.tf.lr), " sp ", hex(curthread.tf.sp), "\n")
			ReplayTrapframe()
			throw("should never be here\n")
		}
	}
}

//go:nosplit
func thread_current() uint32 {
	return uint32(lastrun)
}

//go:nosplit
func hack_futex_arm(uaddr *int32, op, val int32, to *timespec, uaddr2 *int32, val2 int32) int32 {
	FUTEX_WAIT := int32(0)
	FUTEX_WAKE := int32(1)
	uaddrn := uintptr(unsafe.Pointer(uaddr))
	ret := 0
	switch op {
	case FUTEX_WAIT:
		dosleep := *uaddr == val
		if dosleep {
			//enter thread scheduler
			curthread.state = ST_SLEEPING
			curthread.futaddr = uaddrn
			Threadschedule()
			//returns with retval in r0
			ret = int(RR0())
		} else {
			//lost wakeup?
			eagain := -11
			ret = eagain
		}
	case FUTEX_WAKE:
		woke := 0
		for i := 0; i < maxthreads && val > 0; i++ {
			t := &threads[i]
			st := t.state
			if t.futaddr == uaddrn && st == ST_SLEEPING {
				t.state = ST_RUNNABLE
				//t.sleepfor = 0
				t.futaddr = 0
				t.tf.r0 = 0
				//t.sleepret = 0
				val--
				woke++
			}
		}
		print("futex woke ", woke, " threads\n")
		ret = woke
	default:
		print("futex op ", hex(uintptr(op)))
		throw("unexpected futex op")
	}
	return int32(ret)
}
