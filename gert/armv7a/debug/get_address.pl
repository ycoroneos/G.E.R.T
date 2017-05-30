#!/usr/bin/perl -w

# A perl script I found at https://rwmj.wordpress.com/2016/03/17/tracing-qemu-guest-execution/
#

#
# Find everything that looks like a kernel address in the input
# and turn it into a symbol using gdb.
#
# Usage:
#   ksyms.pl vmlinux < input > output
# where 'vmlinux' is the kernel image which must contain debug
# symbols (on Fedora, find this in kernel-debuginfo).

use warnings;
use strict;

my $vmlinux = shift;
my %cache = ();

while (<>) {
    s{(^|\s)([0-9a-f]{6,16})(\s|$)}{ $1 . lookup ($2) . $3 }gei;
    print
}

sub lookup
{
    local $_;
    my $addr = $_[0];

    return $cache{$addr} if exists $cache{$addr};

    # Run gdb to lookup this symbol.
    my $cmd = "arm-none-eabi-gdb -batch -s '$vmlinux' -ex 'info symbol 0x$addr'";
    open PIPE, "$cmd 2>&- |" or die "$cmd: $!";
    my $r = <PIPE>;
    close PIPE;
    chomp $r;
    if ($r =~ m/^No symbol/) {
        # No match, just return the original string, but add the original
        # string to the cache so we don't do the lookup again.
        $cache{$addr} = $addr;
        return $addr;
    }

    # Memoize the match and return it.
    $cache{$addr} = $r;
    return $r;
}
