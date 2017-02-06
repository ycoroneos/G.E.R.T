#!/usr/bin/python2.7
#I shouldn't need to comment any of this
import serial, sys, os, signal, getch


if (len(sys.argv)<1):
    print 'you done <redacted>, try giving me a serial port <redacted>\n'
    sys.exit(1)
ser=serial.Serial(sys.argv[1], 115200, timeout=1)
pid=None

def signal_handler(signal, frame):
    print 'caught sigint... killing everything\n'
    ser.close()
    if (pid!=0):
        os.kill(pid, 9)
        print 'children have been <redacted>...\n'
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)

input='' #lol not using this anymore, ctrl-C is the official way of exiting
pid=os.fork()
if (pid==0):
    #child process
    while (input!='bye'):
        c=ser.read()
        if (c!=''):
            sys.stdout.write(c)
            sys.stdout.flush()
    sys.exit(0)
else:
    #parent process
    while 1:
        ser.write(str(getch.getch()))
