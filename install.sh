#!/bin/sh
for d in */
do
	echo $d |sed 's,/,,'
	(cd $d && CGO_ENABLED=0 go install)
done

