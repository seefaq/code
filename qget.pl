#!/usr/bin/perl
$| = 1; ### turn on autoflush

$qid = $ARGV[0];

printf("Get messages from Qid: $qid\n");
$msgSize = 10000;
$getPriority = -4000;
$ipcState = 5;
while ($qid > 0) {
  msgrcv($qid, $rcvd, $msgSize, $getPriority, $ipcState);
  ($type_rcvd, $msg) = unpack("l! a*", $rcvd);
  $time_st        = `date '+%T'`;
  chop($time_st);
  printf(STDOUT "($time_st) [%s]\n", $msg);
}
