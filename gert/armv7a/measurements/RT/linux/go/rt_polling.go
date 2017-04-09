package main

import "os"

func main() {

	//make 91 an input
	exporter, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	exporter.WriteString("91")
	exporter.Close()

	direction, err := os.OpenFile("/sys/class/gpio/gpio91/direction", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	direction.WriteString("in")
	direction.Close()

	value, err := os.OpenFile("/sys/class/gpio/gpio91/value", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}

	//make 191 the output
	exporter, err = os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	exporter.WriteString("191")
	exporter.Close()

	direction, err = os.OpenFile("/sys/class/gpio/gpio191/direction", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	direction.WriteString("out")
	direction.Close()

	writeval, err := os.OpenFile("/sys/class/gpio/gpio191/value", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	writeval.WriteString("0")

	c := []byte{'0'}
	for c[0] != '0' {
		value.Read(c)
		value.Seek(0, 0)
	}
	writeval.WriteString("1")
	value.Close()
	writeval.Close()
}
