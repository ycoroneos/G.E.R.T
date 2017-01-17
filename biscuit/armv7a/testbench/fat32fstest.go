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
	good, root := fat32_som_start(func() bool {
		fmt.Println("init sdcard yay")
		return true
	}, reader)
	if !good {
		fmt.Println("fat32 init failure")
	}
	fmt.Println(root.getfilenames())
	fmt.Println(root.getsubdirnames())
	good, bootdir := root.direnter("BOOT")
	if !good {
		panic("dir entry failed")
	} else {
		fmt.Println(bootdir.getfilenames())
		good, contents := bootdir.fileread("UENV.TXT")
		if !good {
			panic("file read failure")
		}
		fmt.Println(string(contents))
	}
}
