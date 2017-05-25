#!/usr/bin/python

from svgpathtools import Path, Line, QuadraticBezier, CubicBezier, Arc
from svgpathtools import parse_path
from svgpathtools import svg2paths
import os
import sys
import matplotlib.pyplot as plt
import numpy as np
import json

def sweepbez(c, steps):
    return map(lambda t: c[0]*(1.0-t)**3 +3.0*c[1]*t*(1.0-t)**2 + 3.0*c[2]*(1.0-t)*t**2 + c[3]*t**3, np.linspace(0,1,steps))

def writepoints(xpoints, ypoints, filename):
    d=dict()
    d['X']=map(int, xpoints)
    d['Y']=map(int, ypoints)
    with open(filename, 'wt') as output:
        output.write(json.dumps(d))

if __name__=="__main__":
    if (len(sys.argv)!=2):
        sys.exit(0)
    xpoints=list()
    ypoints=list()
    paths, attributes = svg2paths(sys.argv[1])
    for path in paths:
        for trace in path:
            if (len(trace)==2):
                #line
                xpoints.append(trace.start.real)
                xpoints.append(trace.end.real)
                ypoints.append(trace.start.imag)
                ypoints.append(trace.end.imag)
            elif (len(trace)==4):
                #cubic bezier
                newx = sweepbez(map(lambda x:x.real, trace), 2)
                newy = sweepbez(map(lambda x:x.imag, trace), 2)
                xpoints = xpoints + newx
                ypoints = ypoints + newy
            else:
                print trace
    writepoints(xpoints, ypoints, 'points.txt')
    plt.plot(xpoints, ypoints, marker='o', color='r')
    plt.show()
