package internal

import (
	"errors"
	"log"
	"time"

	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/linuxdvb/frontend"
)

func Tune(fpath, sys, pol string, freqHz int64, bwHz int, sr uint) (fe frontend.Device, err error) {
	var polar rune
	switch pol {
	case "h", "v":
		polar = rune((pol)[0])
	default:
		err = errors.New("unknown polarization: " + pol)
		return
	}
	fe, err = frontend.Open(fpath)
	if err != nil {
		return
	}

	switch sys {
	case "t":
		if err = fe.SetDeliverySystem(dvb.SysDVBT); err != nil {
			return
		}
		if err = fe.SetModulation(dvb.QAMAuto); err != nil {
			return
		}
		if err = fe.SetFrequency(uint32(freqHz)); err != nil {
			return
		}
		if err = fe.SetInversion(dvb.InversionAuto); err != nil {
			return
		}
		if bwHz != 0 {
			if err = fe.SetBandwidth(uint32(bwHz)); err != nil {
				return
			}
		}
		if err = fe.SetCodeRateHP(dvb.FECAuto); err != nil {
			return
		}
		if err = fe.SetCodeRateLP(dvb.FECAuto); err != nil {
			return
		}
		if err = fe.SetTxMode(dvb.TxModeAuto); err != nil {
			return
		}
		if err = fe.SetGuard(dvb.GuardAuto); err != nil {
			return
		}
		if err = fe.SetHierarchy(dvb.HierarchyNone); err != nil {
			return
		}
	case "s", "s2":
		if sys == "s" {
			if err = fe.SetDeliverySystem(dvb.SysDVBS); err != nil {
				return
			}
			if err = fe.SetModulation(dvb.QPSK); err != nil {
				return
			}
		} else {
			if err = fe.SetDeliverySystem(dvb.SysDVBS2); err != nil {
				return
			}
			if err = fe.SetModulation(dvb.PSK8); err != nil {
				return
			}
			if err = fe.SetRolloff(dvb.RolloffAuto); err != nil {
				return
			}
			if err = fe.SetPilot(dvb.PilotAuto); err != nil {
				return
			}
		}
		if err = fe.SetSymbolRate(uint32(sr)); err != nil {
			return
		}
		if err = fe.SetInnerFEC(dvb.FECAuto); err != nil {
			return
		}
		if err = fe.SetInversion(dvb.InversionAuto); err != nil {
			return
		}
		ifreq, tone, volt := frontend.SecParam(freqHz, polar)
		if err = fe.SetFrequency(ifreq); err != nil {
			return
		}
		if err = fe.SetTone(tone); err != nil {
			return
		}
		if err = fe.SetVoltage(volt); err != nil {
			return
		}
	case "ca", "cb", "cc":
		switch sys {
		case "ca":
			err = fe.SetDeliverySystem(dvb.SysDVBCAnnexA)
		case "cb":
			err = fe.SetDeliverySystem(dvb.SysDVBCAnnexB)
		case "cc":
			err = fe.SetDeliverySystem(dvb.SysDVBCAnnexC)
		}
		if err != nil {
			return
		}
		if err = fe.SetModulation(dvb.QAMAuto); err != nil {
			return
		}
		if err = fe.SetFrequency(uint32(freqHz)); err != nil {
			return
		}
		if err = fe.SetInversion(dvb.InversionAuto); err != nil {
			return
		}
		if err = fe.SetSymbolRate(uint32(sr)); err != nil {
			return
		}
		if err = fe.SetInnerFEC(dvb.FECAuto); err != nil {
			return
		}
	default:
		err = errors.New("unknown delivery system: " + sys)
		return
	}

	err = fe.Tune()
	return
}

func WaitForTune(fe frontend.Device, deadline time.Time, debug bool) error {
	fe3 := frontend.API3{fe}
	var ev frontend.Event
	for ev.Status()&frontend.HasLock == 0 {
		timedout, err := fe3.WaitEvent(&ev, deadline)
		if err != nil {
			return err
		}
		if timedout {
			return errors.New("tuning timeout")
		}
		if debug {
			log.Println(ev.Status())
		}
	}
	return nil
}
