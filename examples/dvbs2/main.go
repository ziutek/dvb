package main

import (
	"errors"
	"fmt"
	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/linuxdvb/demux"
	"github.com/ziutek/dvb/linuxdvb/frontend"
	"github.com/ziutek/dvb/ts"
	"os"
	"time"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		if err != dvb.ErrOverflow {
			os.Exit(1)
		}
	}
}

func main() {
	adapter := "/dev/dvb/adapter0"

	fe, err := frontend.Open(adapter + "/frontend0")
	checkErr(err)

	checkErr(fe.SetDeliverySystem(dvb.SysDVBS2))
	checkErr(fe.SetModulation(dvb.PSK8))
	checkErr(fe.SetRolloff(dvb.RolloffAuto))
	checkErr(fe.SetPilot(dvb.PilotAuto))
	checkErr(fe.SetSymbolRate(27500e3))
	checkErr(fe.SetInnerFEC(dvb.FECAuto))
	checkErr(fe.SetInversion(dvb.InversionAuto))
	freq, tone, volt := frontend.SecParam(10911e6, 'v')
	checkErr(fe.SetFrequency(freq))
	checkErr(fe.SetTone(tone))
	checkErr(fe.SetVoltage(volt))

	checkErr(fe.Tune())

	deadline := time.Now().Add(5 * time.Second)
	var ev frontend.Event
	for ev.Status()&frontend.HasLock == 0 {
		timedout, err := frontend.API3{fe}.WaitEvent(&ev, deadline)
		checkErr(err)
		if timedout {
			checkErr(errors.New("tuning timeout"))
		}
		fmt.Println(ev.Status())
	}
	fmt.Println("tuned!")

	filterParam := demux.StreamFilterParam{
		Pid:  8192,
		In:   demux.InFrontend,
		Out:  demux.OutTSDemuxTap,
		Type: demux.Other,
	}
	f, err := demux.Device(adapter + "/demux0").StreamFilter(&filterParam)
	checkErr(err)
	defer f.Close()
	checkErr(f.Start())

	r := ts.NewPktStreamReader(f)
	pkt := new(ts.ArrayPkt)

	n := 10000
	for {
		t := time.Now()
		for i := 0; i < n; i++ {
			checkErr(r.ReadPkt(pkt))
		}
		dt := time.Now().Sub(t)
		pps := time.Duration(n) * time.Second / dt
		fmt.Printf(
			"%d pkt in %s: %d pkt/s (%d B/s)\n",
			n, dt, pps, pps*ts.PktLen,
		)
	}
}
