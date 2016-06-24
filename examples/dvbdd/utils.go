package main

import (
	"fmt"
	"os"

	"github.com/ziutek/dvb"
)

func die(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		if _, ok := err.(dvb.TemporaryError); ok {
			return
		}
		die(err.Error())
	}
}
