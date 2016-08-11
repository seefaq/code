#!/usr/bin/perl
use POSIX;
use Time::Local;

%dayOff = (
  '20160101' =>1,
  '20160222' =>1,
  '20160406' =>1,
  '20160413' =>1,
  '20160414' =>1,
  '20160415' =>1,
  '20160502' =>1,
  '20160505' =>1,
  '20160520' =>1,
  '20160701' =>1,
  '20160718' =>1,
  '20160719' =>1,
  '20160812' =>1,
  '20161024' =>1,
  '20161205' =>1,
  '20161212' =>1,
);

print "syntax: workday.pl [YYYYMMDD]\n";
if ($ARGV[0] ne '') {
  $dt = $ARGV[0];
} else {
  $dt='20160824';
}

&show('20160824', -10, \%dayOff);
&show('20160813', 0, \%dayOff);
&show('20160813', 1, \%dayOff);
print "\n";
&show($dt, -10, \%dayOff);
&show($dt, -1, \%dayOff);
&show($dt, 0, \%dayOff);
&show($dt, 1, \%dayOff);
&show($dt, 10, \%dayOff);

##################
##  SUBROUTINE
##################

sub show {
  local($dt, $day, $holiday) = @_;
  ## dt format: YYYYMMDD
  ## holiday hash [by reference] format: YYYYMMDD => 1
  ## day: number of workdays before/after dt

  if ($holiday->{$dt} || &isWeekEnd($dt)) {
    print "today:$dt(_HOLIDAY_) ";
  } else {
    print "today:$dt(-WORKDAY-) ";
  }
  printf("yesterday:%s tomorrow:%s, next %d workday:%s\n",
      &nextDate($dt,-1), &nextDate($dt,1), $day, &workday($dt, $day, $holiday));
}

sub workday {
  local($dt, $day, $holiday) = @_;
  ## dt format: YYYYMMDD
  ## holiday hash [by reference] format: YYYYMMDD => 1
  ## day: < 0; find work date before dt
  ## day: >= 0; find work date after dt

  if ($day < 0) {
    $day *= (-1);
    $dt=&nextDate($dt, 1) while ($holiday->{$dt} || &isWeekEnd($dt));
    while ($day) {
      if ($holiday->{$dt} || &isWeekEnd($dt)) {
        $dt=&nextDate($dt, -1);
      } else {
        $day--;
        if ($day) {
          $dt=&nextDate($dt, -1);
        } else {
          $dt=&nextDate($dt, -1);
          last;
        }
      }
    }
    $dt=&nextDate($dt, -1) while ($holiday->{$dt} || &isWeekEnd($dt));
  } elsif ($day==0) {
    $dt=&nextDate($dt, 1) while ($holiday->{$dt} || &isWeekEnd($dt));
  } else {
    $dt=&nextDate($dt, -1) while ($holiday->{$dt} || &isWeekEnd($dt));
    while ($day) {
      if ($holiday->{$dt} || &isWeekEnd($dt)) {
        $dt=&nextDate($dt, 1);
      } else {
        $day--;
        if ($day) {
          $dt=&nextDate($dt, 1);
        } else {
          $dt=&nextDate($dt, 1);
          last;
        }
      }
    }
    $dt=&nextDate($dt, 1) while ($holiday->{$dt} || &isWeekEnd($dt));
  }
  $dt;
}

sub nextDate {
  local($dt, $offset) = @_;  ## format: YYYYMMDD
  local($tm, $ret);

  $tm = timelocal(0,0,0,substr($dt,6,2), substr($dt,4,2)-1, substr($dt,0,4));
  $ret = strftime "%Y%m%d", localtime($tm + ($offset*86400));
  $ret;
}

sub isWeekEnd {
  local($dt) = @_;   ## format: YYYYMMDD
  local($tm, $ret);

  $tm = timelocal(0,0,0,substr($dt,6,2), substr($dt,4,2)-1, substr($dt,0,4));
  $wday = (localtime($tm))[6];

  return 1 if ($wday==0 || $wday==6);
  return 0;
}
