package main

import (
	"errors"
	"os"
	"time"

	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/linuxdvb/demux"
	"github.com/ziutek/dvb/linuxdvb/frontend"
	"github.com/ziutek/dvb/ts"
)

var (
	filter demux.StreamFilter
	fe     frontend.Device
)

func tune(fpath, dmxpath, dvrpath, sys, pol string, freqHz int64, bwHz int, sr uint, pids []int16) ts.PktReader {
	var polar rune
	switch pol {
	case "h", "v":
		polar = rune((pol)[0])
	default:
		die("unknown polarization: " + pol)
	}
	var err error
	fe, err = frontend.Open(fpath)
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

	filterParam := demux.StreamFilterParam{
		Pid:  pids[0],
		In:   demux.InFrontend,
		Out:  demux.OutTSDemuxTap,
		Type: demux.Other,
	}
	if dvrpath != "" {
		filterParam.Out = demux.OutTSTap
	}
	filter, err = demux.Device(dmxpath).NewStreamFilter(&filterParam)
	checkErr(err)
	for _, pid := range pids[1:] {
		checkErr(filter.AddPid(pid))
	}
	if dvrpath == "" {
		checkErr(filter.SetBufferLen(1024 * 188))
		checkErr(filter.Start())
		return ts.NewPktStreamReader(filter)
	}
	dvr, err := os.Open(dvrpath)
	checkErr(err)
	checkErr(filter.Start())
	return ts.NewPktStreamReader(dvr)
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
