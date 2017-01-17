// Do not edit. Bootstrap copy of /home/yanni/biscuit/golang1.7/src/cmd/compile/internal/ssa/regalloc_test.go

//line /home/yanni/biscuit/golang1.7/src/cmd/compile/internal/ssa/regalloc_test.go:1
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

import "testing"

func TestLiveControlOps(t *testing.T) {
	c := testConfig(t)
	f := Fun(c, "entry",
		Bloc("entry",
			Valu("mem", OpInitMem, TypeMem, 0, nil),
			Valu("x", OpAMD64MOVLconst, TypeInt8, 1, nil),
			Valu("y", OpAMD64MOVLconst, TypeInt8, 2, nil),
			Valu("a", OpAMD64TESTB, TypeFlags, 0, nil, "x", "y"),
			Valu("b", OpAMD64TESTB, TypeFlags, 0, nil, "y", "x"),
			Eq("a", "if", "exit"),
		),
		Bloc("if",
			Eq("b", "plain", "exit"),
		),
		Bloc("plain",
			Goto("exit"),
		),
		Bloc("exit",
			Exit("mem"),
		),
	)
	flagalloc(f.f)
	regalloc(f.f)
	checkFunc(f.f)
}
