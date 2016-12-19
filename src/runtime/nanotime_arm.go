package runtime

//go:nosplit
func armtime(clocktype int, dest *timespec) {
	//print("arm time\n")
	ts := timer_read()
	//print("time is ", ts.tv_nsec, "\n")
	PutR0(uint32(ts.tv_sec))
	PutR2(uint32(ts.tv_nsec))
	PutR0(uint32(0))
	PutR2(uint32(0))
	return
}

//go:nosplit
func armnanotime(clocktype int, dest *timespec) int64 {
	print("arm nanotime\n")
	ts := timer_read()
	out := int64(0)
	out += int64(1000000000*int64(ts.tv_sec)<<32) + int64(ts.tv_nsec)
	print("time is ", out, "\n")
	return out
}
