package frontend

import (
	"bytes"
	"fmt"
	"github.com/ziutek/dvb"
	"strings"
	"syscall"
	"unsafe"
)

type Caps uint32

const IsStupid Caps = 0

const (
	CanInversionAuto Caps = 1 << iota
	CanFEC12
	CanFEC23
	CanFEC34
	CanFEC45
	CanFEC56
	CanFEC67
	CanFEC78
	CanFEC89
	CanFECAuto
	CanQPSK
	CanQAM16
	CanQAM32
	CanQAM64
	CanQAM128
	CanQAM256
	CanQAMAuto
	CanTxModeAuto
	CanBandwidthAuto
	CanGuardAuto
	CanHierarchyAuto
	Can8VSB
	Can16VSB
	HasExtendedCaps
)

const (
	CanTurboFEC Caps = 0x8000000 << iota
	Can2GModulation
	NeedsBending
	CanRecover
	CanMuteTS
)

// Returns string representation of Caps.
func (c Caps) String() string {
	var can string
	if c&CanInversionAuto != 0 {
		can += "Auto inversion\n"
	}
	var s []string
	if c&CanFEC12 != 0 {
		s = append(s, "1/2")
	}
	if c&CanFEC23 != 0 {
		s = append(s, "2/3")
	}
	if c&CanFEC34 != 0 {
		s = append(s, "3/4")
	}
	if c&CanFEC45 != 0 {
		s = append(s, "4/5")
	}
	if c&CanFEC45 != 0 {
		s = append(s, "4/5")
	}
	if c&CanFEC56 != 0 {
		s = append(s, "5/6")
	}
	if c&CanFEC67 != 0 {
		s = append(s, "6/7")
	}
	if c&CanFEC78 != 0 {
		s = append(s, "7/8")
	}
	if c&CanFEC89 != 0 {
		s = append(s, "8/9")
	}
	if c&CanFECAuto != 0 {
		s = append(s, "auto")
	}
	if len(s) > 0 {
		can += "FEC: " + strings.Join(s, ",") + "\n"
	}
	var mod []string
	if c&CanQPSK != 0 {
		mod = append(mod, "QPSK")
	}
	s = s[:0]
	if c&CanQAM16 != 0 {
		s = append(s, "16")
	}
	if c&CanQAM32 != 0 {
		s = append(s, "32")
	}
	if c&CanQAM64 != 0 {
		s = append(s, "64")
	}
	if c&CanQAM128 != 0 {
		s = append(s, "128")
	}
	if c&CanQAM256 != 0 {
		s = append(s, "256")
	}
	if c&CanQAMAuto != 0 {
		s = append(s, "auto")
	}
	if len(s) > 0 {
		mod = append(mod, "QAM:"+strings.Join(s, ","))
	}
	if len(mod) > 0 {
		can += strings.Join(mod, ", ") + "\n"
	}
	if c&CanTxModeAuto != 0 {
		can += "Auto transmission mode\n"
	}
	if c&CanBandwidthAuto != 0 {
		can += "Auto bandwidth\n"
	}
	if c&CanGuardAuto != 0 {
		can += "Auto guard interval\n"
	}
	if c&CanHierarchyAuto != 0 {
		can += "Auto hierarchy\n"
	}
	if c&Can8VSB != 0 {
		can += "8 VSB\n"
	}
	if c&Can16VSB != 0 {
		can += "16 VSB\n"
	}
	if c&HasExtendedCaps != 0 {
		can += "Extended caps\n"
	}
	if c&CanTurboFEC != 0 {
		can += "Turbo FEC\n"
	}
	if c&Can2GModulation != 0 {
		can += "2G modulation"
	}
	if c&NeedsBending != 0 {
		can += "Needs bending"
	}
	if c&CanRecover != 0 {
		can += "Recover"
	}
	if c&CanMuteTS != 0 {
		can += "TS mute\n"
	}
	return can
}

type Type uint32

const (
	DVBS Type = iota
	DVBC
	DVBT
	ATSC
)

func (t Type) String() string {
	switch t {
	case DVBS:
		return "DVB-S"
	case DVBC:
		return "DVB-C"
	case DVBT:
		return "DVB-T"
	case ATSC:
		return "ATSC"
	}
	return "unknown"
}

// API3 provides Linux DVB API v3 interface to frontend device
type API3 struct {
	Device
}

type Info struct {
	Name          [128]byte
	Type          Type // DEPRECATED. Use Device.DeliverySystem() instead
	FreqMin       uint32
	FreqMax       uint32
	FreqStepSize  uint32
	FreqTolerance uint32
	SRMin         uint32
	SRMax         uint32
	SRTolerance   uint32 // ppm
	NotiferDelay  uint32 // DEPRECATED
	Caps          Caps
}

func (i *Info) String() string {
	name := i.Name[:]
	if k := bytes.IndexByte(name, 0); k != -1 {
		name = name[:k]
	}
	return fmt.Sprintf(
		"Name: %s\n"+
			"Type: %s\n"+
			"Freq: %d - %d kHz (step: %d Hz, tolerance: %d Hz)\n"+
			"Symbol rate: %d - %d kBd (tolerance: %d Bd)\n"+
			"Notifer delay: %d ms\n"+
			"Capabilities:\n%s",
		name, i.Type, i.FreqMin/1000, i.FreqMax/1000,
		i.FreqStepSize, i.FreqTolerance,
		i.SRMin/1000, i.SRMax/1000, i.SRTolerance,
		i.NotiferDelay, i.Caps,
	)
}

// GetInfo works like Info but doesn't allocates memory for Info struct
func (f API3) GetInfo(i *Info) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_GET_INFO,
		uintptr(unsafe.Pointer(i)),
	)
	if e != 0 {
		return e
	}
	return nil
}

// Info returns frontend informations
func (f API3) Info() (*Info, error) {
	i := new(Info)
	err := f.GetInfo(i)
	return i, err
}

type Bandwidth uint32

const (
	Bandwidth8MHz Bandwidth = iota
	Bandwidth7MHz
	Bandwidth6MHz
	BandwidthAuto
	Bandwidth5Mhz
	Bandwidth10Mhz
	Bandwidth1712kHz
)

type ParamDVBT struct {
	Freq       uint32        // frequency in Hz
	Inversion  dvb.Inversion // spectral inversion
	Bandwidth  Bandwidth
	CodeRateHP dvb.CodeRate
	CodeRateLP dvb.CodeRate
	Modulation dvb.Modulation
	TxMode     dvb.TxMode
	Guard      dvb.Guard
	Hierarchy  dvb.Hierarchy
}

func DefaultParamDVBT(c Caps, country string) *ParamDVBT {
	var p ParamDVBT
	if c&CanInversionAuto != 0 {
		p.Inversion = dvb.InversionAuto
	} else {
		p.Inversion = dvb.InversionOff
	}
	if c&CanBandwidthAuto != 0 {
		p.Bandwidth = BandwidthAuto
	} else {
		p.Bandwidth = Bandwidth8MHz
	}
	if c&CanFECAuto != 0 {
		p.CodeRateHP = dvb.FECAuto
		p.CodeRateLP = dvb.FECAuto
	} else {
		p.CodeRateHP = dvb.FEC34
		p.CodeRateLP = dvb.FEC34
	}
	if c&CanQAMAuto != 0 {
		p.Modulation = dvb.QAMAuto
	} else {
		p.Modulation = dvb.QAM64
	}
	if c&CanTxModeAuto != 0 {
		p.TxMode = dvb.TxModeAuto
	} else {
		p.TxMode = dvb.TxMode8k
	}
	if c&CanGuardAuto != 0 {
		p.Guard = dvb.GuardAuto
	} else {
		p.Guard = dvb.Guard8
	}
	if c&CanHierarchyAuto != 0 {
		p.Hierarchy = dvb.HierarchyAuto
	} else {
		p.Hierarchy = dvb.HierarchyNone
	}
	return &p
}

func (f API3) TuneDVBT(p *ParamDVBT) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_SET_FRONTEND,
		uintptr(unsafe.Pointer(p)),
	)
	if e != 0 {
		return e
	}
	return nil
}

type Status uint32

const (
	HasSignal  Status = 1 << iota // found something above the noise level
	HasCarrier                    // found a DVB signal
	HasViterbi                    // FEC is stable
	HasSync                       // found sync bytes
	HasLock                       // everything's working...
	Timedout                      // no lock within the last ~2 seconds
	Reinit                        // frontend was reinitialized, application is recommned to reset DiSEqC, tone and parameters
)

func (s Status) String() string {
	stat := []byte("-sign -carr -vite -sync -lock")
	if s&HasSignal != 0 {
		stat[0] = '+'
	}
	if s&HasCarrier != 0 {
		stat[6] = '+'
	}
	if s&HasViterbi != 0 {
		stat[12] = '+'
	}
	if s&HasSync != 0 {
		stat[18] = '+'
	}
	if s&HasLock != 0 {
		stat[24] = '+'
	}
	ret := string(stat)
	if s&Timedout != 0 {
		ret += " timeout"
	}
	if s&Reinit != 0 {
		ret += " reinit"
	}
	return ret
}

type Event interface {
	Status() Status
}

type EventDVBT struct {
	Status Status
	Param  ParamDVBT
}

// GetEvent can return dvb.OverflowError
func (f API3) GetEventDVBT(ev *EventDVBT) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_GET_EVENT,
		uintptr(unsafe.Pointer(ev)),
	)
	if e != 0 {
		return e
	}
	return nil
}
