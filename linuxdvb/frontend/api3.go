package frontend

import (
	"bytes"
	"fmt"
	"github.com/ziutek/dvb"
	"strings"
	"syscall"
	"time"
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

func (f API3) Info() (*Info, error) {
	i := new(Info)
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_GET_INFO,
		uintptr(unsafe.Pointer(i)),
	)
	if e != 0 {
		return nil, Error{"get", "info", e}
	}
	return i, nil
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

var bandwidthName = []string{
	"8Mhz",
	"7MHz",
	"6MHz",
	"auto",
	"5Mhz",
	"10Mhz",
	"1712kHz",
}

func (b Bandwidth) String() string {
	if b > Bandwidth1712kHz {
		return "unknown"
	}
	return bandwidthName[b]
}

func (f API3) tune(p unsafe.Pointer) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_SET_FRONTEND,
		uintptr(p),
	)
	if e != 0 {
		return e
	}
	return nil
}

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

func (p *ParamDVBT) Tune(f API3) error {
	return f.tune(unsafe.Pointer(p))
}

type ParamDVBS struct {
	Freq       uint32        // frequency in Hz
	Inversion  dvb.Inversion // spectral inversion
	SymbolRate uint32
	InnerFEC   dvb.CodeRate
}

func (p *ParamDVBS) Tune(f API3) error {
	return f.tune(unsafe.Pointer(p))
}

type ParamDVBC struct {
	Freq       uint32        // frequency in Hz
	Inversion  dvb.Inversion // spectral inversion
	SymbolRate uint32
	CodeRate   dvb.CodeRate
	Modulation dvb.Modulation
}

func (p *ParamDVBC) Tune(f API3) error {
	return f.tune(unsafe.Pointer(p))
}

// DefaultParamDVBT returns pointer to ParamDVBT initialized according to
// regulations in specific country.
// TODO: DefaultParamDVBT always setup ParamDVBT for Poland. Add support for
// other countries.
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

func (f API3) Status() (status Status, err error) {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_READ_STATUS,
		uintptr(unsafe.Pointer(&status)),
	)
	if e != 0 {
		err = Error{"get", "status", e}
	}
	return
}

func (f API3) BER() (ber uint32, err error) {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_READ_BER,
		uintptr(unsafe.Pointer(&ber)),
	)
	if e != 0 {
		err = Error{"get", "ber", e}
	}
	return
}

func (f API3) SNR() (snr int16, err error) {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_READ_SNR,
		uintptr(unsafe.Pointer(&snr)),
	)
	if e != 0 {
		err = Error{"get", "snr", e}
	}
	return
}

func (f API3) SignalStrength() (ss int16, err error) {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_READ_SIGNAL_STRENGTH,
		uintptr(unsafe.Pointer(&ss)),
	)
	if e != 0 {
		err = Error{"get", "rssi", e}
	}
	return
}

func (f API3) UncorrectedBlocks() (ublocks uint32, err error) {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_READ_UNCORRECTED_BLOCKS,
		uintptr(unsafe.Pointer(&ublocks)),
	)
	if e != 0 {
		err = Error{"get", "uncorrected_blocks", e}
	}
	return
}

type Event struct {
	status Status
	param  [36]byte // enough to save longest parameters (for DVB-T)
}

func (e *Event) Status() Status {
	return e.status
}

func (e *Event) ParamDVBT() *ParamDVBT {
	return (*ParamDVBT)(unsafe.Pointer(&e.param))
}

func (e *Event) ParamDVBS() *ParamDVBS {
	return (*ParamDVBS)(unsafe.Pointer(&e.param))
}

func (e *Event) ParamDVBC() *ParamDVBC {
	return (*ParamDVBC)(unsafe.Pointer(&e.param))
}

// WaitEvent can return dvb.ErrOverflow. If deadline is non zero time WaitEvent
// returns true if it doesn't receive any event up to deatline.
func (f API3) WaitEvent(ev *Event, deadline time.Time) (bool, error) {
	fd := f.Fd()
	if !deadline.IsZero() {
		timeout := deadline.Sub(time.Now())
		if timeout <= 0 {
			return true, nil
		}
		var r syscall.FdSet
		var n int
		var err error
		r.Bits[fd/64] = 1 << (fd % 64)
		tv := syscall.NsecToTimeval(int64(timeout))
		for {
			n, err = syscall.Select(int(fd+1), &r, nil, nil, &tv)
			if err != syscall.EINTR {
				break
			}
		}
		if err != nil {
			return false, Error{"get", "event (select)", err}
		}
		if n == 0 {
			return true, nil
		}
	}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		_FE_GET_EVENT,
		uintptr(unsafe.Pointer(ev)),
	)
	if e != 0 {
		if e == syscall.EOVERFLOW {
			return false, dvb.ErrOverflow
		}
		return false, Error{"get", "event", e}
	}
	return false, nil
}
