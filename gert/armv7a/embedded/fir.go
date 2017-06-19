// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedded

import "fmt"
import "container/ring"

func fir_main() {
	fmt.Println("making input")
	inputsize := uint32(1)
	input := make(chan uint32, inputsize)
	output := make(chan uint32, inputsize)
	coeffs := []uint32{1, 2, 3}
	go fir(input, output, coeffs)
	for i := uint32(1); i < 20; i++ {
		input <- i
		fmt.Println("MAC result ", <-output)
	}
}

func fir(line_in, line_out chan uint32, coeffs []uint32) {
	r := ring.New(len(coeffs))
	//prime it with 0
	for i := 0; i < len(coeffs); i++ {
		r.Value = uint32(0)
		r = r.Next()
	}
	r = r.Next()
	for {
		//get input and add to ring
		val := <-line_in
		r.Value = val
		r = r.Next()
		fmt.Printf("got %d, here is the buffer now : ", val)
		out := uint32(0)
		for i := 0; i < len(coeffs); i++ {
			fmt.Printf(" %d ", r.Value)
			sample := r.Value.(uint32)
			out = out + sample*coeffs[i]
			r = r.Next()
		}
		fmt.Println("")
		line_out <- out
	}
}
