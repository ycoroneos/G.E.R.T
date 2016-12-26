package runtime

import "unsafe"
import "runtime/internal/atomic"

//for booting
func Runtime_main()

//go:nosplit
func PutR0(val uint32)

//go:nosplit
func PutR2(val uint32)

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

////catching traps
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
	//print("incoming trap: ", trapno, "\n")
	switch trapno {
	case 3:
		//print("spoofing read on: ", arg0, "\n")
		ret := uint32(0xffffffff)
		PutR0(ret)
		return
	case 4:
		//print("spoofing write on: ", arg0, "\n")
		ret := uint32(0xffffffff)
		if arg0 == 1 || arg0 == 2 {
			ret = write_uart_unsafe(uintptr(arg1), arg2)
		} else {
		}
		PutR0(ret)
		return
	case 5:
		print("spoofing open on: ", arg0, "\n")
		PutR0(0xffffffff)
		return
	case 6:
		print("spoofin close on: ", arg0, "\n")
		PutR0(0)
		return
	case 120:
		print("spoofing clone\n")
		thread_id := makethread(uint32(arg0), uintptr(arg1), uintptr(arg2))
		PutR0(uint32(thread_id))
		return
	case 142:
		print("spoofing select\n")
		//throw("select")
		if !panicpanic {
			curthread.state = ST_RUNNABLE
			Threadschedule()
		}
		PutR0(0)
		return
	case 158:
		print("spoofing sys sched yield\n")
		curthread.state = ST_RUNNABLE
		Threadschedule()
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
	case 192:
		throw("mmap trap\n")
		addr := unsafe.Pointer(uintptr(arg0))
		n := uintptr(arg1)
		prot := int32(arg2)
		flags := int32(arg3)
		fd := int32(arg4)
		off := uint32(arg5)
		ret := uint32(uintptr(hack_mmap(addr, n, prot, flags, fd, off)))
		print("mmap returns ", hex(ret), "\n")
		PutR0(ret)
		return
	case 220:
		print("spoofing madvise\n")
		PutR0(0)
		return
	case 224:
		print("gettid\n")
		//throw("gettid")
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
	case 263:
		clock_type := arg0
		ts := ((*timespec)(unsafe.Pointer(uintptr(arg1))))
		clk_gettime(clock_type, ts)
		PutR0(0)
		return
	case 0xbadbabe:
		throw("abort")
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////threading

//go:nosplit
func Threadschedule()

//go:nosplit
func RecordTrapframe()

//go:nosplit
func ReplayTrapframe()

const maxcpus = 4

//cpu states
const (
	CPU_WFI     = 0
	CPU_STARTED = 1
)

var cpustatus [maxcpus]uint32

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
	print("\t\t\t\t LR ", hex(threads[i].tf.lr), " sp: ", hex(threads[i].tf.sp), "\n")
	return 0
}

var lastrun = 0

//go:nosplit
func thread_schedule() {
	print("thread scheduler\n")
	//start looking after the current thread id
	lastrun = (lastrun + 1) % maxthreads
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
	throw("no runnable threads. what happened?\n")
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
		//print("futex woke ", woke, " threads\n")
		ret = woke
	default:
		print("futex op ", hex(uintptr(op)))
		throw("unexpected futex op")
	}
	return int32(ret)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
/////vmem system
type physaddr uint32

//go:nosplit
func loadttbr0(l1base unsafe.Pointer)

//go:nosplit
func loadvbar(vbar_addr unsafe.Pointer)

//go:nosplit
func invallpages()

//go:nosplit
func DMB()

//This file will have all the things to do with the arm MMU and page tables
//assume we will be addressing 4gb of memory
//using the short descriptor page format

const RAM_START = physaddr(0x10000000)
const RAM_SIZE = uint32(0x80000000)
const ONE_MEG = uint32(0x00100000)

const PERIPH_START = physaddr(0x0)
const PERIPH_SIZE = uint32(0x0FFFFFFF)

//1MB pages
const PGSIZE = uint32(0x100000)
const PGSHIFT = uint32(20)
const L1_ALIGNMENT = uint32(0x4000)
const VBAR_ALIGNMENT = uint32(0x20)

var kernelstart physaddr
var kernelsize physaddr
var bootstack physaddr

type Interval struct {
	start uint32
	size  uint32
}

const ELF_MAGIC = 0x464C457F
const ELF_PROG_LOAD = 1

type Elf struct {
	magic       uint32
	e_elf       [12]uint8
	e_type      uint16
	e_machine   uint16
	e_version   uint32
	e_entry     uint32
	e_phoff     uint32
	e_shoff     uint32
	e_flags     uint32
	e_ehsize    uint16
	e_phentsize uint16
	e_phnum     uint16
	e_shentsize uint16
	e_shnum     uint16
	e_shstrndx  uint16
}

type Proghdr struct {
	p_type   uint32
	p_offset uint32
	p_va     uint32
	p_pa     uint32
	p_filesz uint32
	p_memsz  uint32
	p_flags  uint32
	p_align  uint32
}

const KERNEL_END = physaddr(0x40200000)

var boot_end physaddr

const PageInfoSz = uint32(8)

type PageInfo struct {
	next_pageinfo uintptr
	ref           uint32
}

//linear array of struct PageInfos
var npages uint32
var pages uintptr
var pgnfosize uint32 = uint32(8)

//pointer to the next PageInfo to give away
var nextfree uintptr

//L1 table
var l1_table physaddr

//each cpu gets an interrupt stack
//and a boot stack
var isr_stack [4]physaddr
var isr_stack_pt *physaddr = &isr_stack[0]

//vector table
//8 things
//reset, undefined, svc, prefetch abort, data abort, unused, irq, fiq
type vector_table struct {
	reset          uint32
	undefined      uint32
	svc            uint32
	prefetch_abort uint32
	data_abort     uint32
	_              uint32
	irq            uint32
	fiq            uint32
}

var vectab physaddr

//linear array of page directory entries that form the kernel pgdir
var kernpgdir uintptr

//go:nosplit
func roundup(val, upto uint32) uint32 {
	result := (val + (upto - 1)) & ^(upto - 1)
	//	print("rounded ", hex(val), " to ", hex(result), "\n")
	return result
}

//go:nosplit
func verifyzero(addr uintptr, n uint32) {
	print("inside verifyzero\n")
	for i := 0; i < int(n); i++ {
		test := *((*byte)(unsafe.Pointer(addr + uintptr(n))))
		if test > 0 {
			print("verify zero failure\n")
		}
	}
}

//go:nosplit
func memclrbytes(ptr unsafe.Pointer, n uintptr)

//go:nosplit
func boot_alloc(size uint32) physaddr {
	//allocate ROUNDUP(size, PGSIZE) bytes from the boot region
	result := boot_end
	newsize := uint32(roundup(uint32(size), 0x4))
	boot_end = boot_end + physaddr(newsize)
	//print("boot alloc clearing ", hex(uint32(result)), " up to ", hex(uint32(boot_end)), "\n")
	memclrNoHeapPointers(unsafe.Pointer(uintptr(result)), uintptr(newsize))
	//memclrbytes(unsafe.Pointer(uintptr(result)), uintptr(newsize))
	DMB()
	verifyzero(uintptr(result), newsize)
	return result
}

//go:nosplit
func mem_init() {
	print("mem init: ", hex(RAM_SIZE), " bytes of ram\n")
	print("mem init: kernel elf start: ", hex(kernelstart), " kernel elf end: ", hex(kernelstart+kernelsize), "\n")
	print("stack at: ", hex(uint32(bootstack)), "\n")
	//calculate how many pages we can have
	npages = RAM_SIZE / PGSIZE
	print("\t npages: ", npages, "\n")

	//allocate the l1 table
	//4 bytes each and 4096 entries
	boot_end = physaddr(roundup(uint32(bootstack), L1_ALIGNMENT))
	l1_table = boot_alloc(4 * 4096)
	print("\tl1 page table at: ", hex(l1_table), "\n")

	//allocate the vector table
	boot_end = physaddr(roundup(uint32(boot_end), VBAR_ALIGNMENT))
	vectab = boot_alloc(uint32(unsafe.Sizeof(vector_table{})))
	print("\tvector table at: ", hex(vectab), " \n")

	//allocate the spinlock for mmap
	maplock = (*Spinlock_t)(unsafe.Pointer(uintptr(boot_alloc(uint32(unsafe.Sizeof(Spinlock_t{}))))))
	print("\tmap spinlock at: ", hex(uintptr(unsafe.Pointer(maplock))), " \n")

	//allocate the spinlock for mmap
	bootlock = (*Spinlock_t)(unsafe.Pointer(uintptr(boot_alloc(uint32(unsafe.Sizeof(Spinlock_t{}))))))
	print("\tboot spinlock at: ", hex(uintptr(unsafe.Pointer(bootlock))), " \n")

	//allocate pages array outside the runtime's knowledge
	//boot_end = boot_end + physaddr(8*4)
	//boot_end = physaddr(roundup(uint32(boot_end), PGSIZE))
	pages = uintptr(boot_alloc(npages * 8))
	//print("pages at: ", hex(uintptr(unsafe.Pointer(pages))), " sizeof(struct PageInfo) is ", hex(unsafe.Sizeof(*pages)), "\n")
	print("pages at: ", hex(pages), "\n")
	physPageSize = uintptr(PGSIZE)

}

var bootlock *Spinlock_t

//go:nosplit
func cpunum() int

//go:nosplit
func boot_any()

//go:nosplit
func getentry() uint32

//1:0x20d8028 2:0x20d8030 3:0x20d8038 scr:0x20d8000

var scr *uint32 = ((*uint32)(unsafe.Pointer(uintptr(0x20d8000))))
var cpu1bootaddr *uint32 = ((*uint32)(unsafe.Pointer(uintptr(0x20d8028))))
var cpu2bootaddr *uint32 = ((*uint32)(unsafe.Pointer(uintptr(0x20d8030))))
var cpu3bootaddr *uint32 = ((*uint32)(unsafe.Pointer(uintptr(0x20d8038))))

//go:nosplit
func mp_init() {
	//set up stacks, they must be 8 byte aligned

	//first set up isr_stack
	//other cores will boot off the isr stack
	start := uint32(boot_alloc(4 * 1028))
	end := uint32(boot_alloc(0))

	print("start stack: ", hex(start), " end stack: ", hex(end), "\n")
	for i := uint32(0); i < 4; i++ {
		isr_stack[i] = physaddr((end - 1024*i) & 0xFFFFFFF8)
		print("cpu[", i, "] isr stack at ", hex(isr_stack[i]), "\n")
	}
	print("cur cpu: ", cpunum(), "\n")

	entry := getentry()

	Splock(bootlock)
	print("scr reads: ", hex(*scr), "\n")
	print("trying to boot cpu 1, entry is ", hex(entry), "\n")
	//do the boots
	*cpu1bootaddr = entry
	val := *scr
	val |= 0x1<<22 | 0x1<<14 | 0x1<<18
	*scr = val
	//print("cpu1bootaddr reads: ", hex(*cpu1bootaddr), "\n")
	//print("scr reads: ", hex(*scr), "\n")
	for cpustatus[1] == CPU_WFI {
	}
	//Spunlock(bootlock)

}

//go:nosplit
func mp_pen() {
	print("cpu ", cpunum(), " is in the pen\n")
	cpustatus[cpunum()] = CPU_STARTED
	Splock(bootlock)
	Spunlock(bootlock)
}

//go:nosplit
func pgnum2pa(pgnum uint32) physaddr {
	return physaddr(PGSIZE * pgnum)
}

//go:nosplit
func pa2page(pa physaddr) *PageInfo {
	pgnum := uint32(uint32(pa) / PGSIZE)
	return (*PageInfo)(unsafe.Pointer((uintptr(unsafe.Pointer(uintptr(pages))) + uintptr(pgnum*pgnfosize))))
	//return uintptr(pages) + uintptr(pgnum*pgnfosize)
}

//go:nosplit
func pa2pgnum(pa physaddr) uint32 {
	return uint32(pa) / PGSIZE
}

//go:nosplit
func pageinfo2pa(pgnfo *PageInfo) physaddr {
	pgnum := uint32((uintptr(unsafe.Pointer(pgnfo)) - pages) / unsafe.Sizeof(PageInfo{}))
	return pgnum2pa(pgnum)
}

//go:nosplit
func check_page_free(pgnfo *PageInfo) bool {
	curpage := (*PageInfo)(unsafe.Pointer(nextfree))
	for {
		if pgnfo == curpage {
			return true
		}
		if curpage.next_pageinfo == 0 {
			break
		}
		curpage = (*PageInfo)(unsafe.Pointer(curpage.next_pageinfo))
	}
	return false
}

//go:nosplit
func walk_pgdir(pgdir uintptr, va uint32) *uint32 {
	table_index := va >> PGSHIFT
	pte := (*uint32)(unsafe.Pointer(pgdir + uintptr(4*table_index)))
	return pte
}

//go:nosplit
func page_init() {
	//construct a linked-list of free pages
	nfree := uint32(0)
	nextfree = 0
	for i := pa2pgnum(RAM_START); i < pa2pgnum(physaddr(uint32(RAM_START)+RAM_SIZE)); i++ {
		pa := pgnum2pa(i)
		pagenfo := pa2page(pa)
		if pa >= physaddr(RAM_START) && pa < kernelstart {
			pagenfo.next_pageinfo = nextfree
			pagenfo.ref = 0
			nextfree = uintptr(unsafe.Pointer(pagenfo))
			nfree += 1
		} else if pa >= physaddr(KERNEL_END) && pa < physaddr(uint32(RAM_START)+uint32(RAM_SIZE)-uint32(ONE_MEG)) {
			pagenfo.next_pageinfo = nextfree
			pagenfo.ref = 0
			nextfree = uintptr(unsafe.Pointer(pagenfo))
			nfree += 1
		} else {
			pagenfo.ref = 1
			pagenfo.next_pageinfo = 0
		}
	}
	print("page init done\n")
	print("free pages: ", nfree, "\n")
	npagenfo := (*PageInfo)(unsafe.Pointer(nextfree))
	print("next free page is for pa: ", hex(pageinfo2pa(npagenfo)), "\n")
}

//go:nosplit
func page_alloc() *PageInfo {
	freepage := (*PageInfo)(unsafe.Pointer(nextfree))
	nextfree = freepage.next_pageinfo
	return freepage
}

//go:nosplit
func checkcontiguousfree(pgdir uintptr, va, size uint32) bool {
	for start := va; start < va+size; start += PGSIZE {
		//print("checkcontiguous va: ", hex(start), " size: ", hex(size), "\n")
		pgnum := start >> PGSHIFT
		pde := (*uint32)(unsafe.Pointer(pgdir + uintptr(pgnum*4)))
		if *pde&0x2 > 0 {
			//print("\tfalse: ", hex(*pde), "\n")
			return false
		}
	}
	//print("found contiguous\n")
	return true
}

//go:nosplit
func map_region(pa uint32, va uint32, size uint32, perms uint32) {
	//section entry bits
	pa = pa & 0xFFF00000
	va = va & 0xFFF00000
	perms = perms | 0x2
	//realsize := roundup(size, PGSIZE)
	realsize := roundup(size, PGSIZE)
	//print("realsize is ", hex(realsize), "\n")
	i := uint32(0)
	for ; i < realsize; i += PGSIZE {
		//pgnum := pa2pgnum(physaddr(i + pa))
		nextpa := pa + i
		l1offset := nextpa >> 18
		//entry := (*uint32)(unsafe.Pointer((uintptr(unsafe.Pointer(l1_table))) + uintptr(pgnum*4)))
		//print("l1_table: ", hex(uintptr(l1_table)), " offset: ", hex(uint32(l1offset)), "\n")
		//print("entry addr: ", hex(uintptr(l1_table+physaddr(l1offset))), "\n")
		entry := (*uint32)(unsafe.Pointer(uintptr(l1_table + physaddr(l1offset))))
		base_addr := (va + i)
		*entry = base_addr | perms
	}
	print("mapped region va from ", hex(va), " to ", hex(va+i), "\n")
}

//go:nosplit
func map_kernel() {
	//read the kernel elf to find the regions of the kernel
	elf := ((*Elf)(unsafe.Pointer(uintptr(kernelstart))))
	if elf.magic != ELF_MAGIC {
		print("bad kernel elf header\n")
		throw("bad elf header in the kernel\n")
	}

	print("mapping kernel:\n")
	for i := uintptr(0); i < uintptr(elf.e_phnum); i++ {
		ph := ((*Proghdr)(unsafe.Pointer(uintptr(elf.e_phoff) + uintptr(i*unsafe.Sizeof(Proghdr{})) + uintptr(kernelstart))))
		if ph.p_type == ELF_PROG_LOAD {
			filesz := ph.p_filesz
			pa := ph.p_pa
			va := ph.p_va
			print("\tkernel pa: ", hex(pa), " va: ", hex(va), " size: ", hex(filesz), "\n")
			map_region(pa, va, filesz, 0x0)
		}
	}

	//install the kernel page table

	//map the uart
	print("mapping uart\n")
	map_region(0x02000000, 0x02000000, PGSIZE, 0x0)

	//map the timer
	print("mapping peripherals + timer\n")
	map_region(uint32(PERIPH_START), uint32(PERIPH_START), PGSIZE, 0x0)

	//map the stack and boot_alloc scratch space
	print("mapping stack + page tables\n")
	nextfree := boot_alloc(0)
	if uint32(nextfree) < (uint32(RAM_START) + RAM_SIZE - ONE_MEG) {
		throw("out of scratch space\n")
	}
	map_region(uint32(uint32(RAM_START)+RAM_SIZE-ONE_MEG), uint32(RAM_START)+RAM_SIZE-ONE_MEG, PGSIZE, 0x0)

	//map the boot rom
	map_region(uint32(0x0), uint32(0x0), PGSIZE, 0x0)

	//identity map [kernelstart, boot_alloc(0))
	//	print("kernel start is ", hex(uint32(kernelstart)), "\n")
	//
	//	map_region(uint32(kernelstart), uint32(kernelstart), uint32(KERNEL_END-kernelstart), 0x0)
	//
	//	map_region(uint32(uint32(RAM_START)+RAM_SIZE-ONE_MEG), uint32(RAM_START)+RAM_SIZE-ONE_MEG, PGSIZE, 0x0)
	//	print("boot_alloc(0) is ", hex(uint32(boot_alloc(0))), "\n")
	//	showl1table()
	//loadvbar(unsafe.Pointer(uintptr(vectab)))
	loadttbr0(unsafe.Pointer(uintptr(l1_table)))
	kernpgdir = (uintptr)(unsafe.Pointer(uintptr(l1_table)))
	print("mapped kernel identity\n")
}

//go:nosplit
func showl1table() {
	print("l1 table: ", hex(uint32(l1_table)), "\n")
	print("__________________________\n")
	for i := uint32(0); i < 4096; i += 1 {
		entry := *(*uint32)(unsafe.Pointer((uintptr(l1_table)) + uintptr(i*4)))
		if entry == 0 {
			continue
		}
		base := entry & 0xFFF00000
		perms := entry & 0x3
		print("\t| entry: ", i, ", base: ", hex(base), " perms: ", hex(perms), "\n")
	}
	print("__________________________\n")
}

//go:nosplit
func l1_walk() {
}

type Spinlock_t struct {
	v uint32
}

//go:nosplit
func Splock(l *Spinlock_t) {
	for {
		if atomic.Xchg(&l.v, 1) == 1 {
			break
		}
		//for l.v != 0 {
		//htpause()
		//}
	}
}

//go:nosplit
func Spunlock(l *Spinlock_t) {
	atomic.Store(&l.v, 0)
	//	l.v = 0
}

const MMAP_FIXED = uint32(0x10)

//var maplock = &Spinlock_t{}
var maplock *Spinlock_t

// mmap calls the mmap system call.  It is implemented in assembly.
//go:nosplit
func hack_mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) unsafe.Pointer {
	Splock(maplock)
	va := uintptr(addr)
	size := uint32(roundup(uint32(n), PGSIZE))
	print("mmap_arm: ", hex(va), " ", hex(n), " ", hex(prot), " ", hex(flags), "\n")

	if va == 0 {
		//throw("cowardly refusing to map 0\n")
		//need to find a contiguous amount of virtual mem with size
		//find the first chunk of contiguous free space in virtual memory
		for pgnum := uint32(KERNEL_END >> PGSHIFT); pgnum < 4096; pgnum += 1 {
			if checkcontiguousfree(kernpgdir, uint32(pgnum<<PGSHIFT), size) == true {
				va = uintptr(pgnum << PGSHIFT)
				break
			}
		}
		if va == 0 {
			throw("cant find large enough chunk of contiguous virtual memory\n")
		}
	}

	clear := true
	for start := va; start < (va + uintptr(size)); start += uintptr(PGSIZE) {
		pte := walk_pgdir(kernpgdir, uint32(start))
		if *pte&0x2 > 0 {
			print("mmap_fixed failure for va: ", hex(start), " because it's already mapped\n")
			print("pte addr ", hex(uintptr(unsafe.Pointer(pte))), " contents ", hex(*pte), "\n")
			clear = false
			continue
		}
		page := page_alloc()
		if page == nil {
			throw("mmap_arm: out of memory\n")
		}
		pa := pageinfo2pa(page) & 0xFFF00000
		*pte = uint32(pa) | 0x2
	}
	//showl1table()
	print("reloading page table\n")
	invallpages()
	//memclrbytes(unsafe.Pointer(va), uintptr(size))
	if clear == true {
		print("clearing... ")
		memclrNoHeapPointers(unsafe.Pointer(va), uintptr(size))
		//memclrbytes(unsafe.Pointer(va), uintptr(size))
	}
	Spunlock(maplock)
	print("updated page tables -> ", hex(va), "\n")
	return unsafe.Pointer(va)
}

//go:nosplit
func clk_gettime(clock_type uint32, ts *timespec) {
	print("spoof clock_gettime\n")
	ts.tv_sec = 0
	ts.tv_nsec = 0
}
