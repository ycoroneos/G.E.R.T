package runtime

type thread_t struct {
	pc    uintptr
	sp    uintptr
	lr    uintptr
	state uint32
	r0    uint32
	r1    uint32
	r2    uint32
	r3    uint32
	r4    uint32
	r5    uint32
	r6    uint32
	r7    uint32
	r8    uint32
	r9    uint32
	r10   uint32
	r11   uint32
	r12   uint32
}

// maximum # of runtime "OS" threads
const maxthreads = 64

var threads [maxthreads]thread_t

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
func makethread(entry, stack uintptr, flags uint32) int {
	CLONE_VM := 0x100
	CLONE_FS := 0x200
	CLONE_FILES := 0x400
	CLONE_SIGHAND := 0x800
	CLONE_THREAD := 0x10000
	chk := uint32(CLONE_VM | CLONE_FS | CLONE_FILES | CLONE_SIGHAND |
		CLONE_THREAD)
	if flags != chk {
		print("unexpected clone args", uintptr(flags))
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
	threads[i].pc = entry
	threads[i].sp = stack
	threads[i].state = ST_RUNNABLE
	return i
}

var lastrun = 0

//never return
//go:nosplit
func schedule_thread() {
	for ; lastrun < maxthreads; lastrun = (lastrun + 1) % maxthreads {
		if threads[lastrun].state == ST_RUNNABLE {
			switch_thread(&threads[lastrun])
		}
	}
}
