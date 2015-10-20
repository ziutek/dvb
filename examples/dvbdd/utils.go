package main

import (
	"fmt"
	"os"

	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/ts"
)

func die(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		if err == dvb.ErrOverflow || err == ts.ErrSync {
			return
		}
		die(err.Error())
	}
}
