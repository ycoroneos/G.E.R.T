package main

import "os"
import "fmt"

func main() {

	result := make(chan int)
	go pollpin("91", result)
	go pollpin("191", result)
	go pollpin("24", result)
	go pollpin("200", result)
	count := 0
	for i := 0; i < 4; i++ {
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
	count := 0
	oldc := []byte{'0'}
	c := []byte{'0'}
	for i := 0; i < 1000000; i++ {
		value.Read(c)
		value.Seek(0, 0)
		if c[0] == oldc[0] {
			continue
		} else {
			oldc[0] = c[0]
			count++
		}
	}
	value.Close()
	resp <- count
	return
}
