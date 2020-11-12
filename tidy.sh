#!/bin/bash
set -e

function isGoMod() {
	res=0
	for f in $1/*
	do
		if [ "${f##*.}"x = "mod"x ]; then
		  res=1
		  break
		fi
	done
	if [ res = 1 ]; then
		return 1
	fi
	#echo $res
}

function modTidy() {
	for file in $1/*
	do
		if test -d $file ; then
			isGoMod $file
			re2=`echo $?`
			if [ $re2 = 0 ]; then
				current=$(pwd)
				cd $file
				echo "go mod tidy $file"
				go mod tidy
				cd $current
			fi
			modTidy $file
		fi
	done
}

modTidy ./pkg