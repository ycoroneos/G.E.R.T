package main

import "os"
import "fmt"
import "time"
import "runtime"
import "strconv"

func main() {

	n, _ := strconv.Atoi(os.Args[1])
	runtime.GOMAXPROCS(n)
	result := make(chan int)
	if n > 0 {
		go pollpin("91", result)
	}
	if n > 1 {
		go pollpin("90", result)
	}
	if n > 2 {
		go pollpin("24", result)
	}
	if n > 3 {
		go pollpin("200", result)
	}
	count := 0
	for i := 0; i < n; i++ {
		count += <-result
	}
	fmt.Printf("count is %d\n", count)
}

func pollpin(pin string, resp chan int) {
	exporter, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	exporter.WriteString(pin)
	exporter.Close()

	direction, err := os.OpenFile("/sys/class/gpio/gpio"+pin+"/direction", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	direction.WriteString("in")
	direction.Close()

	value, err := os.OpenFile("/sys/class/gpio/gpio"+pin+"/value", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	count := uint32(0)
	oldc := []byte{'0'}
	c := []byte{'0'}
	start := time.Now()
	curtime := start
	goodcount := uint32(0)
	for time.Now().Sub(start).Seconds() <= 10.0 {
		datime := time.Now()
		if datime.Sub(curtime).Seconds() > 1.0 {
			curtime = datime
			fmt.Printf("%v count -> %v\n", pin, count)
			goodcount = count
			count = 0
		}
		value.Read(c)
		value.Seek(0, 0)
		if c[0] == '1' && oldc[0] == '0' {
			count++
		}
		oldc[0] = c[0]
	}
	value.Close()
	resp <- int(goodcount)
	return
}
