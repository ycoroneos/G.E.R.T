#!/usr/bin/perl -w

#A perl script I found at https://rwmj.wordpress.com/2016/03/17/tracing-qemu-guest-execution/
#

use warnings;
use strict;

my $ts = 0;

while (<>) {
    my $ts_delta;
    my $pc;

    if (m{^exec_tb(_nocache)? ([-\d.]+).*pc=0x([a-fA-F0-9]+)}) {
        $ts_delta = $2;
        $pc = "$3";
    }
    elsif (m{^exec_tb_exit ([-\d.]+)}) {
        $ts_delta = $1;
    }
    elsif (m{^Dropped_Event ([-\d.]+)}) {
        $ts_delta = $1;
    }
    else {
        die "could not parse output: $_"
    }
    $ts += $ts_delta;

    if (defined $pc) {
        print "$ts $pc\n";
    }
}
