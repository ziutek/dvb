#!/bin/bash
GOARCH=arm GOARM=7 go install && scp -P 10022 $GOPATH/bin/linux_arm/pcrjtter_example $1:/sdcard
