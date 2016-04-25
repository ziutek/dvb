package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ziutek/dvb/linuxdvb/demux"
	"github.com/ziutek/dvb/linuxdvb/frontend"

	"github.com/ziutek/dvb/examples/internal"
)

func die(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		die(err.Error())
	}
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		"Usage: %s [OPTION] PID [PID...]\nUse PID=8192 to obtain all PIDs on DVR.\nOptions:\n",
		filepath.Base(os.Args[0]),
	)
	flag.PrintDefaults()
	os.Exit(1)
}

var filter demux.StreamFilter

func main() {
	fpath := flag.String(
		"front", "/dev/dvb/adapter0/frontend0",
		"path to the frontend device",
	)
	dmxpath := flag.String(
		"demux", "/dev/dvb/adapter0/demux0",
		"path to the demux device",
	)
	sys := flag.String(
		"sys", "t",
		"delivery system type: t, s, s2, ca, cb, cc",
	)
	freq := flag.Float64(
		"freq", 0,
		"frequency [Mhz]",
	)
	sr := flag.Uint(
		"sr", 0,
		"symbol rate [kBd]",
	)
	pol := flag.String(
		"pol", "h",
		"polarization: h, v",
	)
	bw := flag.Uint(
		"bw", 0,
		"bandwidth [MHz] (default 0 means automatic - not supported by many tuners)",
	)
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	pids := make([]int16, flag.NArg())
	for i, a := range flag.Args() {
		pid, err := strconv.ParseInt(a, 0, 64)
		checkErr(err)
		if uint64(pid) > 8192 {
			die(a + " isn't in valid PID range [0, 8192]")
		}
		pids[i] = int16(pid)
	}

	fe, err := internal.Tune(*fpath, *sys, *pol, int64(*freq*1e6), int(*bw*1e6), *sr)
	checkErr(err)
	checkErr(internal.WaitForTune(fe, time.Now().Add(7*time.Second), true))

	filterParam := demux.StreamFilterParam{
		Pid:  pids[0],
		In:   demux.InFrontend,
		Out:  demux.OutTSTap,
		Type: demux.Other,
	}
	filter, err = demux.Device(*dmxpath).NewStreamFilter(&filterParam)
	for _, pid := range pids[1:] {
		checkErr(filter.AddPid(pid))
	}
	checkErr(err)
	checkErr(filter.Start())

	hasublk := true
	for {
		fe3 := frontend.API3{fe}
		sig, err := fe3.SignalStrength()
		checkErr(err)
		snr, err := fe3.SNR()
		checkErr(err)
		ber, err := fe3.BER()
		checkErr(err)
		var ublk uint32
		if hasublk {
			ublk, err = fe3.UncorrectedBlocks()
			hasublk = (err == nil)
		}
		if hasublk {
			log.Printf("sig: %d  snr: %d  ber: %d  ublk: %d", sig, snr, ber, ublk)
		} else {
			log.Printf("sig: %d  snr: %d  ber: %d", sig, snr, ber)
		}
		time.Sleep(2 * time.Second)
	}
}
