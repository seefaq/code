#!/usr/bin/perl
$qkey = 0x0;
$flag = 1000 | 2;
$qid  = msgget($qkey, $flag);
print "got Qid [$qid], <Enter> to remove queue";
<>;
msgctl($qid, IPC_RMID, 0);
print STDOUT $!;
