// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9
// +build !solaris
// +build !windows
// +build !nacl
// +build !linux !amd64

package runtime

import "unsafe"

const MMAP_FIXED = uint32(0x10)

//var maplock = &Spinlock_t{}
var maplock *Spinlock_t

// mmap calls the mmap system call.  It is implemented in assembly.
//go:nosplit
func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) unsafe.Pointer {
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
		memclr(unsafe.Pointer(va), uintptr(size))
	}
	if armhackmode == 0 {
		print("hackmode nuked\n")
	}
	Spunlock(maplock)
	print("updated page tables -> ", hex(va), "\n")
	return unsafe.Pointer(va)
}
