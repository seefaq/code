#!/usr/bin/perl

###### from "$130.00","$2,200.00",,"See Co.,",1234,,"$1,230.63"
###### to   $130.00|$2,200.00||See Co.,|1234||$1,230.63

$in=0;
while ($line = <>) {
  chop($line);
  @chars = split(//, $line);
  foreach $ch (@chars) {
    if ($ch eq '"') {
      $in = (!$in)? 1:0;
      next;
    } elsif ($ch eq ',' && !$in) {
      print '|';
      next;
    }
    print $ch;
  }
  print "\n";
}
