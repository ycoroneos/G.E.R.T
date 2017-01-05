package runtime

var Armhackmode uint32
var btrace = false

func biscuitproccount() int32 {
	return 4
}

func BTrace() int {
	btrace = true
	return 0
}
func brk() {
	func() {
		print("bkpt")
	}()
	for {
	}
}
