package main

import (
	"fmt"
	"time"
)

func user_init() {
	fmt.Printf("Hello World!\n")
}

func user_loop() {
	fmt.Printf("time is %v\n", time.Now())
	time.Sleep(2 * time.Second)
}
