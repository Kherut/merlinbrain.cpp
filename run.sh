#!/bin/sh

./pulse.sh &
echo $!
/usr/bin/go run main.go -- $!