#!/usr/bin/env bash

:<<!
[case]

title=check string matches pattern
cid=0
pid=0

[group]
1. exactly match            >> hello
2. regular expression match >> 1[0-9]{10}
3. format string match      >> %s%d

[esac]
!

str="abc123"

echo ">> hello"
echo ">> 13905120512"
echo ">> abc123"
