// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
)

const infile = "points.txt"
const outfile = "bindata.gob"

type PPoint struct {
	X []int
	Y []int
}

type CompactPoint struct {
	X     uint16
	Y     uint16
	Color uint8
}

func main() {
	file, err := os.Open(infile)
	check(err)
	defer file.Close()
	dec := json.NewDecoder(file)
	var p PPoint
	dec.Decode(&p)
	newpoints := make([]CompactPoint, len(p.X))
	for i := 0; i < len(p.X); i++ {
		newpoints[i].X = uint16(p.X[i])
		newpoints[i].Y = uint16(p.Y[i])
	}

	fmt.Printf("GOBing points\n")
	file, err = os.Create(outfile)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(newpoints)
	}
	fmt.Printf("done\n")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
