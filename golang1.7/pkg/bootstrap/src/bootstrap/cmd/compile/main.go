// Do not edit. Bootstrap copy of /Users/fruit/Documents/biscuit/golang1.7/src/cmd/compile/main.go

//line /Users/fruit/Documents/biscuit/golang1.7/src/cmd/compile/main.go:1
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bootstrap/cmd/compile/internal/amd64"
	"bootstrap/cmd/compile/internal/arm"
	"bootstrap/cmd/compile/internal/arm64"
	"bootstrap/cmd/compile/internal/gc"
	"bootstrap/cmd/compile/internal/mips"
	"bootstrap/cmd/compile/internal/mips64"
	"bootstrap/cmd/compile/internal/ppc64"
	"bootstrap/cmd/compile/internal/s390x"
	"bootstrap/cmd/compile/internal/x86"
	"bootstrap/cmd/internal/obj"
	"fmt"
	"log"
	"os"
)

func main() {
	// disable timestamps for reproducible output
	log.SetFlags(0)
	log.SetPrefix("compile: ")

	switch obj.GOARCH {
	default:
		fmt.Fprintf(os.Stderr, "compile: unknown architecture %q\n", obj.GOARCH)
		os.Exit(2)
	case "386":
		x86.Init()
	case "amd64", "amd64p32":
		amd64.Init()
	case "arm":
		arm.Init()
	case "arm64":
		arm64.Init()
	case "mips", "mipsle":
		mips.Init()
	case "mips64", "mips64le":
		mips64.Init()
	case "ppc64", "ppc64le":
		ppc64.Init()
	case "s390x":
		s390x.Init()
	}

	gc.Main()
	gc.Exit(0)
}
