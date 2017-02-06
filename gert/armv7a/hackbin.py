#!/usr/bin/python
import os
import sys
import struct

lookfor=0xe52de004 #push lr
replace=0xe320f000 #nop

if (__name__=="__main__"):
    binname=sys.argv[1]
    outfile=sys.argv[2]
    hacked=False
    with open(binname, "rb") as bfile:
        bindata=bfile.read()
    with open(outfile, "wb") as nfile:
        for i in range(0,len(bindata),4):
            word=struct.unpack('I',bindata[i:i+4])[0]
            if (word==lookfor and hacked==False):
                #write new instruction
                nfile.write(struct.pack('I',replace))
                hacked=True
            else:
                if (hacked==True):
                    #write old instruction
                    nfile.write(struct.pack('I',word))
