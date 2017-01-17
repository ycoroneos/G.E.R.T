// Do not edit. Bootstrap copy of /Users/fruit/Documents/biscuit/golang1.7/src/cmd/compile/internal/ssa/stackframe.go

//line /Users/fruit/Documents/biscuit/golang1.7/src/cmd/compile/internal/ssa/stackframe.go:1
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

// stackframe calls back into the frontend to assign frame offsets.
func stackframe(f *Func) {
	f.Config.fe.AllocFrame(f)
}
