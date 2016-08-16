#!/usr/bin/awk -f
BEGIN { max=0 }
{
  l = length($0);
  if (max < l) max=l;
}
END { print max }
