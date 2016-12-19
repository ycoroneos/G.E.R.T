package runtime

import "unsafe"

const PERIPH_BASE uint32 = 0xA00000
const PERIPH_BASE_SIZE uint32 = 0x2000
const GLOBAL_TIMER_BASE uint32 = PERIPH_BASE + 0x200

var timer_lo *uint32 = ((*uint32)(unsafe.Pointer(uintptr(GLOBAL_TIMER_BASE + 0x0))))
var timer_hi *uint32 = ((*uint32)(unsafe.Pointer(uintptr(GLOBAL_TIMER_BASE + 0x4))))
var timer_ctrl *uint32 = ((*uint32)(unsafe.Pointer(uintptr(GLOBAL_TIMER_BASE + 0x8))))
var timer_isr *uint32 = ((*uint32)(unsafe.Pointer(uintptr(GLOBAL_TIMER_BASE + 0xC))))
var timer_cmplo *uint32 = ((*uint32)(unsafe.Pointer(uintptr(GLOBAL_TIMER_BASE + 0x10))))
var timer_cmphi *uint32 = ((*uint32)(unsafe.Pointer(uintptr(GLOBAL_TIMER_BASE + 0x14))))
var timer_autoinc *uint32 = ((*uint32)(unsafe.Pointer(uintptr(GLOBAL_TIMER_BASE + 0x18))))

var count uint64
var global_time timespec

//starts the timer to count up with no rollover since interrupts arent enabled yet
//go:nosplit
func timer_start() {
	//disable it
	*timer_ctrl = 0x0

	//set the prescaler
	*timer_ctrl |= 0 << 8

	//init global count
	count = 0
	*timer_lo = 0
	*timer_hi = 0

	//enable timer
	*timer_ctrl |= 0x1
}

//go:nosplit
func timer_read() timespec {
	hi := uint32(0)
	lo := uint32(0)
	done := false
	for !done {
		hi = *timer_hi
		lo = *timer_lo
		done = *timer_hi == hi
	}
	global_time.tv_sec = int32(hi)
	global_time.tv_nsec = int32(lo)
	return global_time
}

//go:nosplit
func timer_test() {
	for {
		ts := timer_read()
		print("sec ", ts.tv_sec, " nsec ", ts.tv_nsec, "\n")
	}
}
