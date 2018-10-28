#!/bin/sh

while true
do
	echo 1 > /sys/class/leds/green_led/brightness
	sleep 0.9
	echo 0 > /sys/class/leds/green_led/brightness
	sleep 0.5
done