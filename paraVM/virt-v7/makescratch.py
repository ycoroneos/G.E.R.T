#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os
import sys

if __name__=="__main__":
    template = "DATA runtime·scratch+%d(SB)/4, $0\n"
    end = "GLOBL runtime·scratch(SB), $%d\n"
    size=int(sys.argv[1], 16)
    print(str(size))
    if (size%0x4 !=0):
        print('size is not multiple of 4 bytes')
        sys.exit(1)
    with open('scratchasm_arm.s', 'wt') as f:
        for i in range(0, size/0x4):
            f.write(template % (i*0x4))
        f.write(end % size)
    print('done!')

