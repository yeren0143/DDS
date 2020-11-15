#!/bin/bash
set -e

function format() {
	for file in $1/*
	do
		if test -d $file ; then
			current=$(pwd)
			cd $file
			echo "go mod format $file"
			gofmt -w .
			cd $current
			format $file
		fi
	done
}

format ./pkg
