#!/bin/sh
for d in */
do
	echo $d |sed 's,/,,'
	(cd $d && go install)
done

