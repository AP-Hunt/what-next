#!/bin/bash
. "$PWD/scripts/demo-magic.sh"
export TYPE_SPEED=13
clear

pei "what-next calendar add demo demo.ical"
echo ""

pei "what-next todo add \"begin adding things to my todo list\" --due @today --duration 15m"
sleep 3

clear

pei "what-next"
sleep 10
echo ""