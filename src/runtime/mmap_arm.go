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

	//if fixed mapping check that those addresses are not in use
	//	if (uint32(flags) & MMAP_FIXED) > 0 {
	//		for start := va; start < (va + uintptr(size)); start += uintptr(PGSIZE) {
	//			pte := walk_pgdir(kernpgdir, uint32(start))
	//			if *pte&0x10 > 0 {
	//				print("mmap_fixed failure for va: ", hex(start), " because it's already mapped\n")
	//			}
	//		}
	//	} else {
	//		throw("mmap_arm: cant handle non-fixed mappings yet\n")
	//	}
	for start := va; start < (va + uintptr(size)); start += uintptr(PGSIZE) {
		pte := walk_pgdir(kernpgdir, uint32(start))
		if *pte&0x10 > 0 {
			print("mmap_fixed failure for va: ", hex(start), " because it's already mapped\n")
		}
		page := page_alloc()
		//print("got page\n")
		if page == nil {
			throw("mmap_arm: out of memory\n")
		}
		pa := pageinfo2pa(page) & 0xFFF00000
		*pte = uint32(pa) | 0x2
		//print("mapped va: ", hex(start), " to pa: ", hex(pa), "\n")
	}
	//showl1table()
	print("reloading page table\n")
	invallpages()
	//memclrbytes(unsafe.Pointer(va), uintptr(size))
	memclr(unsafe.Pointer(va), uintptr(size))
	Spunlock(maplock)
	print("updated page tables\n")
	return addr
}
