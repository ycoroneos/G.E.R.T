// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

TEXT ·Set(SB), 4, $0
	MOVW ptr+0(FP), R2
	MOVW val+4(FP), R3
	MOVW R3, (R2)
	RET
