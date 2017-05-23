#!/usr/bin/python
import matplotlib.pyplot as plt
import matplotlib

xaxis = [1,2,3,4]
GERT=[131,262,375,504]
LinuxC=[96, 186, 276, 355]
LinuxGo=[47, 93, 138, 174]

matplotlib.rcParams.update({'font.size': 33})
plt.ylim([0, 550])
plt.ylabel('events/sec')
plt.xlabel('# cpus')
plt.plot(xaxis, GERT, '-o', linewidth=2.0, label='GERT', markersize=10)
plt.plot(xaxis, LinuxC, '--o', linewidth=2.0, label='Linux C', markersize=10)
plt.plot(xaxis, LinuxGo, '-.o', linewidth=2.0, label='Linux Go', markersize=10)
plt.legend(loc=2,prop={'size':48})
plt.show()
