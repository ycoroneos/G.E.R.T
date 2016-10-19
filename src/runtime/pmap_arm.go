package runtime

import "unsafe"

type physaddr uint32

//go:nosplit
func loadttbr0(l1base unsafe.Pointer)

//go:nosplit
func loadvbar(vbar_addr unsafe.Pointer)

//This file will have all the things to do with the arm MMU and page tables
//assume we will be addressing 4gb of memory
//using the short descriptor page format

const RAM_START = physaddr(0x10000000)
const RAM_SIZE = uint32(0x80000000)

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
func memclrbytes(ptr unsafe.Pointer, n uintptr)

//go:nosplit
func boot_alloc(size uint32) physaddr {
	//allocate ROUNDUP(size, PGSIZE) bytes from the boot region
	result := boot_end
	newsize := uint32(roundup(uint32(size), 0x4))
	boot_end = boot_end + physaddr(newsize)
	//print("boot alloc clearing ", hex(uint32(result)), " up to ", hex(uint32(boot_end)), "\n")
	memclrbytes(unsafe.Pointer(uintptr(result)), uintptr(newsize))
	return result
}

//go:nosplit
func mem_init() {
	print("mem init: ", hex(RAM_SIZE), " bytes of ram\n")
	print("mem init: kernel start: ", hex(kernelstart), " kernel end: ", hex(kernelstart+kernelsize), "\n")
	print("stack at: ", hex(uint32(bootstack)), "\n")
	//calculate how many pages we can have
	npages = RAM_SIZE / PGSIZE
	print("\t npages: ", npages, "\n")

	//allocate the l1 table
	//4 bytes each and 4096 entries
	//l1_table = boot_end
	//memclr(unsafe.Pointer(uintptr(l1_table)), uintptr(4*4096))
	//boot_end = boot_end + physaddr(4*4096)
	boot_end = physaddr(roundup(uint32(kernelstart+kernelsize), L1_ALIGNMENT))
	l1_table = boot_alloc(4 * 4096)
	print("\tl1 page table at: ", hex(l1_table), "\n")

	//allocate the vector table
	boot_end = physaddr(roundup(uint32(boot_end), VBAR_ALIGNMENT))
	vectab = boot_alloc(uint32(unsafe.Sizeof(vector_table{})))
	print("\tvector table at: ", hex(vectab), " \n")

	//allocate the spinlock for mmap
	maplock = (*Spinlock_t)(unsafe.Pointer(uintptr(boot_alloc(uint32(unsafe.Sizeof(Spinlock_t{}))))))
	print("\tmap spinlock at: ", hex(uintptr(unsafe.Pointer(maplock))), " \n")

	//allocate pages array outside the runtime's knowledge
	//boot_end = boot_end + physaddr(8*4)
	//boot_end = physaddr(roundup(uint32(boot_end), PGSIZE))
	pages = uintptr(boot_alloc(npages * 8))
	//print("pages at: ", hex(uintptr(unsafe.Pointer(pages))), " sizeof(struct PageInfo) is ", hex(unsafe.Sizeof(*pages)), "\n")
	print("pages at: ", hex(pages), "\n")
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
		} else if pa >= boot_alloc(0) && pa < (RAM_START+physaddr(RAM_SIZE)) {
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
func map_region(pa uint32, va uint32, size uint32, perms uint32) {
	//section entry bits
	pa = pa & 0xFFF00000
	va = va & 0xFFF00000
	perms = perms | 0x2
	//realsize := roundup(size, PGSIZE)
	realsize := roundup(size, PGSIZE)
	print("realsize is ", hex(realsize), "\n")
	for i := uint32(0); i < realsize; i += PGSIZE {
		//pgnum := pa2pgnum(physaddr(i + pa))
		nextpa := pa + i
		l1offset := nextpa >> 18
		//entry := (*uint32)(unsafe.Pointer((uintptr(unsafe.Pointer(l1_table))) + uintptr(pgnum*4)))
		entry := (*uint32)(unsafe.Pointer(uintptr(l1_table + physaddr(l1offset))))
		base_addr := (va + i)
		*entry = base_addr | perms
	}
}

//go:nosplit
func map_kernel() {
	//install the kernel page table

	//map the uart
	map_region(0x02000000, 0x02000000, PGSIZE, 0x0)

	//identity map [kernelstart, boot_alloc(0))
	print("kernel start is ", hex(uint32(kernelstart)), "\n")
	map_region(uint32(kernelstart), uint32(kernelstart), uint32(boot_alloc(0)-kernelstart), 0x0)
	print("boot_alloc(0) is ", hex(uint32(boot_alloc(0))), "\n")
	//	showl1table()
	loadvbar(unsafe.Pointer(uintptr(vectab)))
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
