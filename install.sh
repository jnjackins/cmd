#!/bin/sh
export PATH=$HOME/bin:$GOROOT/bin:/bin:/usr/bin:/usr/local/bin
export CGO_ENABLED=0
if uname |grep -q Darwin; then export CGO_ENABLED=1; fi
for d in */
do
	cmd=$(echo $d |sed 's,/,,')
	before=$(ls -s $HOME/bin/$cmd |awk '{print $1}')
	(cd $d && go install)
	after=$(ls -s $HOME/bin/$cmd |awk '{print $1}')
	printf "$cmd\t%5d -> %5d\n" $before $after
done

