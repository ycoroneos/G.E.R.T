package runtime

var time int64 = 0

//go:nosplit
func armnanotime(clocktype int, timespec *int64) {
	print("armnanotime clock type: ", clocktype, "\n")
	time = time + 1
	*timespec = time
	print("done nanotime\n")
	return
}
