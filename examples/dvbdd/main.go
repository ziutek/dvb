package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
		"Usage: %s [OPTION]\nOptions:\n",
		filepath.Base(os.Args[0]),
	)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	fpath := flag.String("front", "/dev/dvb/adapter0/frontend0", "path to frontend device")
	dpath := flag.String("demux", "/dev/dvb/adapter0/demux0", "path to demux device")
	sys := flag.String("sys", "t", "name of delivery system: t, s, s2, ca, cb, cc")
	freq := flag.Uint("freq", 0, "frequency [Mhz]")
	sr := flag.Uint("sr", 0, "symbol rate [kBd]")
	pol := flag.String("pol", "h", "polarization: h, v")
	flag.Usage = usage
	flag.Parse()

	var polar rune
	switch *pol {
	case "h", "v":
		polar = rune((*pol)[0])
	default:
		die("unknown polarization: " + *pol)
	}

	fe, err := frontend.Open(*fpath)
	checkErr(err)
	defer fe.Close()

	switch *sys {
	case "t":
		checkErr(fe.SetDeliverySystem(dvb.SysDVBT))
		checkErr(fe.SetModulation(dvb.QAMAuto))
		checkErr(fe.SetFrequency(uint32(*freq) * 1e6))
		checkErr(fe.SetInversion(dvb.InversionAuto))
		//checkErr(fe.SetBandwidth(8e6))
		checkErr(fe.SetCodeRateHP(dvb.FECAuto))
		checkErr(fe.SetCodeRateLP(dvb.FECAuto))
		checkErr(fe.SetTxMode(dvb.TxModeAuto))
		checkErr(fe.SetGuard(dvb.GuardAuto))
		checkErr(fe.SetHierarchy(dvb.HierarchyNone))
	case "s", "s2":
		if *sys == "s" {
			checkErr(fe.SetDeliverySystem(dvb.SysDVBS))
			checkErr(fe.SetModulation(dvb.QPSK))
		} else {
			checkErr(fe.SetDeliverySystem(dvb.SysDVBS2))
			checkErr(fe.SetModulation(dvb.PSK8))
			checkErr(fe.SetRolloff(dvb.RolloffAuto))
			checkErr(fe.SetPilot(dvb.PilotAuto))
		}
		checkErr(fe.SetSymbolRate(uint32(*sr)))
		checkErr(fe.SetInnerFEC(dvb.FECAuto))
		checkErr(fe.SetInversion(dvb.InversionAuto))
		ifreq, tone, volt := frontend.SecParam(uint64(*freq)*1e6, polar)
		checkErr(fe.SetFrequency(ifreq))
		checkErr(fe.SetTone(tone))
		checkErr(fe.SetVoltage(volt))
	case "ca", "cb", "cc":
		switch *sys {
		case "ca":
			checkErr(fe.SetDeliverySystem(dvb.SysDVBCAnnexA))
		case "cb":
			checkErr(fe.SetDeliverySystem(dvb.SysDVBCAnnexB))
		case "cc":
			checkErr(fe.SetDeliverySystem(dvb.SysDVBCAnnexC))
		}
		checkErr(fe.SetModulation(dvb.QAMAuto))
		checkErr(fe.SetFrequency(uint32(*freq) * 1e6))
		checkErr(fe.SetInversion(dvb.InversionAuto))
		checkErr(fe.SetSymbolRate(uint32(*sr)))
		checkErr(fe.SetInnerFEC(dvb.FECAuto))
	default:
		die("unknown delivery system: " + *sys)
	}

	checkErr(fe.Tune())
	checkErr(waitForTune(fe))

	var filterParam = demux.StreamFilterParam{
		Pid:  8192,
		In:   demux.InFrontend,
		Out:  demux.OutTSDemuxTap,
		Type: demux.Other,
	}
	f, err := demux.Device(*dpath).StreamFilter(&filterParam)
	checkErr(err)
	defer f.Close()
	checkErr(f.SetBufferLen(1024 * 188))
	checkErr(f.Start())

	r := ts.NewPktStreamReader(f)
	pkt := new(ts.ArrayPkt)

	for {
		checkErr(r.ReadPkt(pkt))
		_, err = os.Stdout.Write(pkt.Bytes())
		checkErr(err)
	}
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
