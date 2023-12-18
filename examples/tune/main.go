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
		"delivery system type: t, t2, s, s2, ca, cb, cc",
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

	fe, err := internal.Tune(*fpath, *sys, *pol, int64(*freq*1e6), int(*bw*1e6), *sr*1e3)
	checkErr(err)
	checkErr(internal.WaitForTune(fe, time.Now().Add(15*time.Second), true))

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

	var rssi, snr, ber, ublk string
	for {
		fe3 := frontend.API3{fe}
		if rssi != "-" {
			val, err := fe3.SignalStrength()
			if err != nil {
				if e, ok := err.(frontend.Error); ok && e.What == "rssi" {
					rssi = "-"
				} else {
					die(err.Error())
				}
			} else {
				rssi = strconv.FormatInt(int64(val), 10)
			}
		}
		if snr != "-" {
			val, err := fe3.SNR()
			if err != nil {
				if e, ok := err.(frontend.Error); ok && e.What == "snr" {
					snr = "-"
				} else {
					die(err.Error())
				}
			} else {
				snr = strconv.FormatInt(int64(val), 10)
			}
		}
		if ber != "-" {
			val, err := fe3.BER()
			if err != nil {
				if e, ok := err.(frontend.Error); ok && e.What == "ber" {
					ber = "-"
				} else {
					die(err.Error())
				}
			} else {
				ber = strconv.FormatInt(int64(val), 10)
			}
		}
		if ublk != "-" {
			val, err := fe3.UncorrectedBlocks()
			if err != nil {
				if e, ok := err.(frontend.Error); ok && e.What == "uncorrected_blocks" {
					ublk = "-"
				} else {
					die(err.Error())
				}
			} else {
				ublk = strconv.FormatUint(uint64(val), 10)
			}
		}
		log.Printf("API3 RSSI: %s  SNR: %s  BER: %s  UBLK: %s", rssi, snr, ber, ublk)
		s, err := fe.Stat()
		if err == nil {
			log.Printf(
				"API5 Sig: %v  CNR: %v  PreIFEC: %v/%v bit/bit  PostIFEC: %v/%v bit/bit  PostOFEC %v/%v blk/blk\n",
				s.Signal, s.CNR, s.PreErrBit, s.PreTotBit, s.PostErrBit, s.PostTotBit, s.ErrBlk, s.TotBlk,
			)
		}
		time.Sleep(2 * time.Second)
	}
}
