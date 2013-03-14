package main

import (
	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/linuxdvb/frontend"
	"log"
	"time"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fe, err := frontend.Open("/dev/dvb/adapter0/frontend0")
	checkErr(err)
	checkErr(fe.SetDeliverySystem(dvb.SysDVBT))
	checkErr(fe.SetFrequency(778e6))
	checkErr(fe.SetBandwidth(8e6))
	checkErr(fe.SetCodeRateHP(dvb.FECAuto))
	checkErr(fe.SetCodeRateLP(dvb.FECAuto))
	checkErr(fe.SetTxMode(dvb.TxModeAuto))
	checkErr(fe.SetGuard(dvb.GuardAuto))
	checkErr(fe.SetHierarchy(dvb.HierarchyNone))
	checkErr(fe.Tune())

	deadline := time.Now().Add(5 * time.Second)
	var ev frontend.Event
	for ev.Status()&frontend.HasLock == 0 {
		timedout, err := frontend.API3{fe}.WaitEvent(&ev, deadline)
		checkErr(err)
		if timedout {
			log.Fatal("tuning timeout")
		}
		log.Println(ev.Status())
	}
	log.Println("tuned!")
}
