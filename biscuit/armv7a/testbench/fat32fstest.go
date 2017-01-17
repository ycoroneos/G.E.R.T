package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	dat, err := ioutil.ReadFile("sdcard.img")
	if err != nil {
		panic(err)
	}
	fmt.Println("Opened sdcard image")
	reader := func(len, offset uint32) (bool, []byte) {
		return true, dat[offset : offset+len]
	}
	good, _ := fat32_som_start(func() bool {
		fmt.Println("init sdcard yay")
		return true
	}, reader)
	if !good {
		fmt.Println("fat32 init failure")
	}
}
