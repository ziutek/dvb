package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/linuxdvb/demux"
	"github.com/ziutek/dvb/linuxdvb/frontend"
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

func usage() {
	fmt.Fprintf(
		os.Stderr,
		"Usage: %s [OPTION] PID [PID...]\nOptions:\n",
		filepath.Base(os.Args[0]),
	)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	src := flag.String("src", "rf", "source: rf, udp")
	laddr := flag.String("laddr", "0.0.0.0:1234", "listen on laddr")
	fpath := flag.String("front", "/dev/dvb/adapter0/frontend0", "path to frontend device")
	dpath := flag.String("demux", "/dev/dvb/adapter0/demux0", "path to demux device")
	sys := flag.String("sys", "t", "name of delivery system: t, s, s2, ca, cb, cc")
	freq := flag.Float64("freq", 0, "frequency [Mhz]")
	sr := flag.Uint("sr", 0, "symbol rate [kBd]")
	pol := flag.String("pol", "h", "polarization: h, v")
	count := flag.Uint64("count", 0, "number of MPEG-TS packets to process (0 means infinity)")
	bw := flag.Int("bw", 0, "bandwidth [MHz] (0 == auto)")
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

	var r ts.PktReader

	switch *src {
	case "rf":
		r = tune(*fpath, *dpath, *sys, *pol, int64(*freq*1e6), *bw*1e6, *sr, pids)
	case "udp":
		r = listenUDP(*laddr, pids)
	default:
		die("Unknown source: " + *src)
	}

	pkt := new(ts.ArrayPkt)

	if *count == 0 {
		for {
			checkErr(r.ReadPkt(pkt))
			_, err := os.Stdout.Write(pkt.Bytes())
			checkErr(err)
		}
		return
	}
	for *count != 0 {
		checkErr(r.ReadPkt(pkt))
		_, err := os.Stdout.Write(pkt.Bytes())
		checkErr(err)
		*count--
	}
}

type pidFilter struct {
	r    ts.PktReader
	pids []int16
}

func (f *pidFilter) ReadPkt(pkt ts.Pkt) error {
	for {
		if err := f.r.ReadPkt(pkt); err != nil {
			return err
		}
		pid := pkt.Pid()
		// TODO: sort f.pids to use more effecitve search method.
		for _, p := range f.pids {
			if p == 8192 || p == pid {
				return nil
			}
		}
	}
}

func listenUDP(laddr string, pids []int16) ts.PktReader {
	la, err := net.ResolveUDPAddr("udp", laddr)
	checkErr(err)
	c, err := net.ListenUDP("udp", la)
	checkErr(err)
	checkErr(c.SetReadBuffer(2 * 1024 * 1024))
	return &pidFilter{
		r:    ts.NewPktPktReader(c, make([]byte, 7*ts.PktLen)),
		pids: pids,
	}
}

func tune(fpath, dpath, sys, pol string, freqHz int64, bwHz int, sr uint, pids []int16) ts.PktReader {
	var polar rune
	switch pol {
	case "h", "v":
		polar = rune((pol)[0])
	default:
		die("unknown polarization: " + pol)
	}

	fe, err := frontend.Open(fpath)
	checkErr(err)

	switch sys {
	case "t":
		checkErr(fe.SetDeliverySystem(dvb.SysDVBT))
		checkErr(fe.SetModulation(dvb.QAMAuto))
		checkErr(fe.SetFrequency(uint32(freqHz)))
		checkErr(fe.SetInversion(dvb.InversionAuto))
		if bwHz != 0 {
			checkErr(fe.SetBandwidth(uint32(bwHz)))
		}
		checkErr(fe.SetCodeRateHP(dvb.FECAuto))
		checkErr(fe.SetCodeRateLP(dvb.FECAuto))
		checkErr(fe.SetTxMode(dvb.TxModeAuto))
		checkErr(fe.SetGuard(dvb.GuardAuto))
		checkErr(fe.SetHierarchy(dvb.HierarchyNone))
	case "s", "s2":
		if sys == "s" {
			checkErr(fe.SetDeliverySystem(dvb.SysDVBS))
			checkErr(fe.SetModulation(dvb.QPSK))
		} else {
			checkErr(fe.SetDeliverySystem(dvb.SysDVBS2))
			checkErr(fe.SetModulation(dvb.PSK8))
			checkErr(fe.SetRolloff(dvb.RolloffAuto))
			checkErr(fe.SetPilot(dvb.PilotAuto))
		}
		checkErr(fe.SetSymbolRate(uint32(sr)))
		checkErr(fe.SetInnerFEC(dvb.FECAuto))
		checkErr(fe.SetInversion(dvb.InversionAuto))
		ifreq, tone, volt := frontend.SecParam(freqHz, polar)
		checkErr(fe.SetFrequency(ifreq))
		checkErr(fe.SetTone(tone))
		checkErr(fe.SetVoltage(volt))
	case "ca", "cb", "cc":
		switch sys {
		case "ca":
			checkErr(fe.SetDeliverySystem(dvb.SysDVBCAnnexA))
		case "cb":
			checkErr(fe.SetDeliverySystem(dvb.SysDVBCAnnexB))
		case "cc":
			checkErr(fe.SetDeliverySystem(dvb.SysDVBCAnnexC))
		}
		checkErr(fe.SetModulation(dvb.QAMAuto))
		checkErr(fe.SetFrequency(uint32(freqHz)))
		checkErr(fe.SetInversion(dvb.InversionAuto))
		checkErr(fe.SetSymbolRate(uint32(sr)))
		checkErr(fe.SetInnerFEC(dvb.FECAuto))
	default:
		die("unknown delivery system: " + sys)
	}

	checkErr(fe.Tune())
	checkErr(waitForTune(fe))

	var filterParam = demux.StreamFilterParam{
		Pid:  pids[0],
		In:   demux.InFrontend,
		Out:  demux.OutTSDemuxTap,
		Type: demux.Other,
	}
	f, err := demux.Device(dpath).StreamFilter(&filterParam)
	checkErr(err)
	for _, pid := range pids[1:] {
		checkErr(f.AddPid(pid))
	}
	checkErr(f.SetBufferLen(1024 * 188))
	checkErr(f.Start())

	return ts.NewPktStreamReader(f)
}

func waitForTune(fe frontend.Device) error {
	fe3 := frontend.API3{fe}
	deadline := time.Now().Add(5 * time.Second)
	var ev frontend.Event
	for ev.Status()&frontend.HasLock == 0 {
		timedout, err := fe3.WaitEvent(&ev, deadline)
		if err != nil {
			return err
		}
		if timedout {
			return errors.New("tuning timeout")
		}
	}
	return nil
}
