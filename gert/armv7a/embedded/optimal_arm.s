TEXT ·Set(SB), 4, $0
	MOVW ptr+0(FP), R2
	MOVW val+4(FP), R3
	MOVW R3, (R2)
	RET
