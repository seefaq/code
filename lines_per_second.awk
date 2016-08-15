#!/usr/bin/awk -f
BEGIN { cnt=max=sec=line=0; oldt=""; }
{
    t=substr($0, 2, 8);  ## format ^(HH:MM:SS.ssss
    if (oldt=="") { oldt=t; sec=1; }
    if (oldt==t) {
        cnt++;
        if (max < cnt) max=cnt;
    } else {
        printf("%s  %d   > sec=%d  max=%d  avg=%d  (%d lines)\n", oldt, cnt, sec, max, line/sec, line);
        sec++;
        cnt=1;
    }
    oldt = t;
    line++;
}
END {
    printf("%s  %d   > sec=%d  max=%d  avg=%d  (%d lines)\n", oldt, cnt, sec, max, line/sec, line);
}
