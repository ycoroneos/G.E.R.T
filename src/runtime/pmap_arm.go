//+build !linux

package runtime

//This file will have all the things to do with the arm MMU and page tables
//assume we will be addressing 4gb of memory
//using the short descriptor page format

type PageInfo struct {
	next *PageInfo
	ref  uint32
}

//linear array of struct PageInfos
var pages *PageInfo

//pointer to the next PageInfo to give away
var nextfree *PageInfo

//linear array of page directory entries that form the kernel pgdir
var kernpgdir *uint32

//go:nosplit
func page_init() {
	for i := 0; i < uint(1<<31); i++ {
	}
}
