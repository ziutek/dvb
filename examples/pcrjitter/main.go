package main

import (
	"errors"
	"fmt"
	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/linuxdvb/demux"
	"github.com/ziutek/dvb/linuxdvb/frontend"
	"github.com/ziutek/dvb/ts"
	"github.com/ziutek/sched"
	"github.com/ziutek/thread"
	"os"
	"runtime"
	"time"
)

const (
	adpath  = "/dev/dvb/adapter0"
	fepath  = adpath + "/frontend0"
	dmxpath = adpath + "/demux0"
	dvrpath = adpath + "/dvr0"
	freq    = 778 // MHz
	pcrpid  = 202
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err != dvb.ErrOverflow && err != ts.ErrSync {
			os.Exit(1)
		}
	}
}

type PCR struct {
	lastPCR  ts.PCR
	lastRead time.Time
	firstPCR time.Time

	cnt       uint32
	jitterSum time.Duration
	jitterMax time.Duration

	discard bool
}

func (p *PCR) reset() {
	p.cnt = 0
	p.jitterSum = 0
	p.jitterMax = 0
	p.discard = false
}

func (p *PCR) PrintReport() {
	defer p.reset()
	p.discard = true

	if p.cnt == 0 {
		fmt.Println("no PCR data")
		return
	}
	if p.cnt < 5 {
		fmt.Println("too few PCR values received: ", p.cnt)
		return
	}
	cnt := time.Duration(p.cnt)
	fmt.Printf(
		"period: %s, jitter: avg=%s, max=%s\n",
		p.lastRead.Sub(p.firstPCR)/cnt, p.jitterSum/(cnt-1), p.jitterMax,
	)
}

func (p *PCR) Loop(dvr ts.PktReader) {
	if os.Geteuid() == 0 {
		runtime.LockOSThread()
		t := thread.Current()
		fmt.Println("Setting realtime sheduling for thread:", t)
		p := sched.Param{
			Priority: sched.FIFO.MaxPriority(),
		}
		checkErr(t.SetSchedPolicy(sched.FIFO, &p))
	} else {
		fmt.Println(
			"Running without root privileges: realtime scheduling disabled",
		)
	}
	fmt.Println()

	pkt := new(ts.ArrayPkt)

	for {
		err := dvr.ReadPkt(pkt)
		now := time.Now()
		checkErr(err)

		if p.discard {
			continue
		}
		if !pkt.Flags().ContainsAF() {
			continue
		}
		af := pkt.AF()
		if af.Flags()&ts.ContainsPCR == 0 {
			continue
		}
		pcr, err := af.PCR()
		if err != nil {
			continue
		}

		if p.cnt == 0 {
			p.firstPCR = now
		} else {
			pcrDiff := (pcr - p.lastPCR).Nanosec()
			readDiff := now.Sub(p.lastRead)
			var jitter time.Duration
			if pcrDiff > readDiff {
				jitter = pcrDiff - readDiff
			} else {
				jitter = readDiff - pcrDiff
			}
			p.jitterSum += jitter
			if p.jitterMax < jitter {
				p.jitterMax = jitter
			}
		}

		p.lastPCR = pcr
		p.lastRead = now
		p.cnt++
	}
}

func main() {
	fe, err := frontend.Open(fepath)
	checkErr(err)
	fe3 := frontend.API3{fe}

	feInfo, err := fe3.Info()
	checkErr(err)
	fmt.Println("Frontend information\n")
	fmt.Println(feInfo)

	if feInfo.Type != frontend.DVBT {
		fmt.Fprintln(
			os.Stderr,
			"This application supports only DVB-T frontend.",
		)
		os.Exit(1)
	}

	fmt.Printf("Tuning to %d MHz...\n", uint(freq))
	feParam := frontend.DefaultParamDVBT(feInfo.Caps, "pl")
	feParam.Freq = freq * 1e6
	checkErr(feParam.Tune(fe3))

	deadline := time.Now().Add(5 * time.Second)
	var ev frontend.Event
	for ev.Status()&frontend.HasLock == 0 {
		timedout, err := fe3.WaitEvent(&ev, deadline)
		checkErr(err)
		if timedout {
			checkErr(errors.New("tuning timeout"))
		}
		fmt.Println("FE status:", ev.Status())
	}
	fmt.Println()

	dmx := demux.Device(dmxpath)
	filterParam := demux.StreamFilterParam{
		Pid:  pcrpid,
		In:   demux.InFrontend,
		Out:  demux.OutTSTap,
		Type: demux.Other,
	}
	filter, err := dmx.NewStreamFilter(&filterParam)
	checkErr(err)
	checkErr(filter.Start())

	file, err := os.Open(dvrpath)
	checkErr(err)

	var pcr PCR

	go pcr.Loop(ts.NewPktStreamReader(file))

	for {
		time.Sleep(5 * time.Second)
		pcr.PrintReport()
	}
}
