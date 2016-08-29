for f in `ls $*`; do
  [ -z $a ] && a=$f && continue
  if diff $a $f > /dev/null ; then
    echo "$a $f nodiff "
  else
    echo "$a $f diff "
    a=$f
  fi
done
