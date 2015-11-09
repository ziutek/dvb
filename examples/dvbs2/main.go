package main

import (
	"flag"
	"fmt"
	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/linuxdvb/demux"
	"github.com/ziutek/dvb/linuxdvb/frontend"
	"github.com/ziutek/dvb/ts"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err != dvb.ErrOverflow && err != ts.ErrSync {
			os.Exit(1)
		}
	}
}

func die(info string) {
	fmt.Fprintln(os.Stderr, info)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		"Usage:\n  %s [flags] FreqMhz {h|v} SRkBaud PID1 [PID2...]\n",
		filepath.Base(os.Args[0]),
	)
	flag.PrintDefaults()
}

func main() {
	adapterPath := flag.String(
		"a", "/dev/dvb/adapter0",
		"path to the adapter directory",
	)
	frontendPath := flag.String(
		"f", "frontend0",
		"frontend path, relative to the adapter directory",
	)
	demuxPath := flag.String(
		"d", "demux0",
		"demux path, relative to the adapter directory",
	)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 4 {
		usage()
		os.Exit(1)
	}

	freq, err := strconv.ParseInt(args[0], 0, 64)
	checkErr(err)
	freq *= 1e6

	var polar byte
	switch args[1] {
	case "v", "V":
		polar = 'v'
	case "h", "H":
		polar = 'h'
	default:
		die("wrong polarisation: " + args[1])
	}

	sr, err := strconv.ParseUint(args[2], 0, 32)
	checkErr(err)
	sr *= 1e3

	args = args[3:]
	pids := make([]int16, len(args))
	for i, a := range args {
		pid, err := strconv.ParseInt(a, 0, 64)
		checkErr(err)
		if uint64(pid) > 8192 {
			die(a + " isn't in valid PID range [0, 8192]")
		}
		pids[i] = int16(pid)
	}

	fe, err := frontend.Open(filepath.Join(*adapterPath, *frontendPath))
	checkErr(err)

	checkErr(fe.SetDeliverySystem(dvb.SysDVBS2))
	checkErr(fe.SetModulation(dvb.PSK8))
	checkErr(fe.SetRolloff(dvb.RolloffAuto))
	checkErr(fe.SetPilot(dvb.PilotAuto))
	checkErr(fe.SetSymbolRate(uint32(sr)))
	checkErr(fe.SetInnerFEC(dvb.FECAuto))
	checkErr(fe.SetInversion(dvb.InversionAuto))
	ifreq, tone, volt := frontend.SecParam(freq, rune(polar))
	checkErr(fe.SetFrequency(ifreq))
	checkErr(fe.SetTone(tone))
	checkErr(fe.SetVoltage(volt))

	checkErr(fe.Tune())

	deadline := time.Now().Add(5 * time.Second)
	var ev frontend.Event
	for ev.Status()&frontend.HasLock == 0 {
		timedout, err := frontend.API3{fe}.WaitEvent(&ev, deadline)
		checkErr(err)
		if timedout {
			die("tuning timeout")
		}
		fmt.Fprintln(os.Stderr, ev.Status())
	}
	fmt.Fprintln(os.Stderr, "tuned!")

	dmx := demux.Device(filepath.Join(*adapterPath, *demuxPath))
	f, err := dmx.NewStreamFilter(
		&demux.StreamFilterParam{
			Pid:  pids[0],
			In:   demux.InFrontend,
			Out:  demux.OutTSDemuxTap,
			Type: demux.Other,
		},
	)
	checkErr(err)
	defer f.Close()

	for _, pid := range pids[1:] {
		checkErr(f.AddPid(pid))
	}
	checkErr(f.SetBufferSize(1024 * ts.PktLen))
	checkErr(f.Start())

	r := ts.NewPktStreamReader(f)
	pkt := new(ts.ArrayPkt)

	for {
		t := time.Now()
		n := 4000
		for i := 0; i < n; i++ {
			checkErr(r.ReadPkt(pkt))
		}
		dt := time.Now().Sub(t)
		pps := time.Duration(n) * time.Second / dt
		fmt.Printf("%d pkt/s (%d kb/s)\n", pps, pps*ts.PktLen*8/1000)
	}
}
