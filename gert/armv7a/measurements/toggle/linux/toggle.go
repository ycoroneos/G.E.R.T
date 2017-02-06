package main

import "os"

func main() {
	exporter, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY,0777)
	if err != nil {
		panic(err)
	}
	exporter.WriteString("91")
	exporter.Close()

	direction, err := os.OpenFile("/sys/class/gpio/gpio91/direction", os.O_WRONLY,0777)
	if err != nil {
		panic(err)
	}
	direction.WriteString("out")
	direction.Close()

	value, err := os.OpenFile("/sys/class/gpio/gpio91/value", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000000; i++ {
		value.WriteString("0")
		value.WriteString("1")
	}
	value.Close()
}
