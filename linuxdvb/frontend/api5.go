package frontend

import (
	"github.com/ziutek/dvb"
	"os"
	"syscall"
	"unsafe"
)

type Device struct {
	file *os.File
}

func Open(filepath string) (d Device, err error) {
	d.file, err = os.OpenFile(filepath, os.O_RDWR, 0)
	return
}

func OpenRO(filepath string) (d Device, err error) {
	d.file, err = os.Open(filepath)
	return
}

func (d Device) Close() error {
	return d.file.Close()
}

func (d Device) Fd() uintptr {
	return d.file.Fd()
}

type cmd uint32

// APIv5 properties
const (
	dtvUndefined cmd = iota
	dtvTune
	dtvClear
	dtvFrequency
	dtvModulation
	dtvBandwidthHz
	dtvInversion
	dtvDiseqcMaster
	dtvSymbolRate
	dtvInnerFEC
	dtvVoltage
	dtvTone
	dtvPilot
	dtvRolloff
	dtvDiseqcSlaveReply
	dtvFeCapabilityCount
	dtvFeCapability
	dtvDeliverySystem
	dtvISDBTPartialReception
	dtvISDBTSoundBroadcasting
	dtvISDBTSBSubchannelId
	dtvISDBTSBSegmentIdx
	dtvISDBTSBSegmentCount
	dtvISDBTLayeraFEC
	dtvISDBTLayeraModulation
	dtvISDBTLayeraSegmentCount
	dtvISDBTLayeraTimeInterleaving
	dtvISDBTLayerbFEC
	dtvISDBTLayerbModulation
	dtvISDBTLayerbSegmentCount
	dtvISDBTLayerbTimeInterleaving
	dtvISDBTLayercFEC
	dtvISDBTLayercModulation
	dtvISDBTLayercSegmentCount
	dtvISDBTLayercTimeInterleaving
	dtvAPIVersion
	dtvCodeRateHP
	dtvCodeRateLP
	dtvGuardInterval
	dtvTransmissionMode
	dtvHierarchy
	dtvISDBTLayerEnabled
	dtvISDBSTSId
	dtvDVBT2PLPId
)

type property struct {
	cmd      cmd
	reserved [3]uint32
	
	// union
	data     uint32
	bufData      [28]byte
	bufLen       uint32
	bufReserved1 [3]uint32
	bufReserved2 uintptr

	result int
}

type properties struct {
	num   uint32
	props *property
}

func (d Device) set(c cmd, data uint32) syscall.Errno {
	p := property{cmd: c, data: data}
	ps := properties{1, &p}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		d.Fd(),
		_FE_SET_PROPERTY,
		uintptr(unsafe.Pointer(&ps)),
	)
	return e
}

func (d Device) get(c cmd) (uint32, syscall.Errno) {
	p := property{cmd: c}
	ps := properties{1, &p}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		d.Fd(),
		_FE_GET_PROPERTY,
		uintptr(unsafe.Pointer(&ps)),
	)
	return p.data, e
}

func (d Device) Version() (major, minor int, err error) {
	var u uint32
	u, e := d.get(dtvAPIVersion)
	if e != 0 {
		err = Error{"get", "api version", e}
		return
	}
	major = int(u>>8) & 0xff
	minor = int(u) & 0xff
	return
}

func (d Device) DeliverySystem() (dvb.DeliverySystem, error) {
	ds, e := d.get(dtvDeliverySystem)
	if e != 0 {
		return 0, Error{"get", "dlivery system", e}
	}
	return dvb.DeliverySystem(ds), nil
}

func (d Device) SetDeliverySystem(ds dvb.DeliverySystem) error {
	e := d.set(dtvDeliverySystem, uint32(ds))
	if e != 0 {
		return Error{"set", "dlivery system", e}
	}
	return nil

}

func (d Device) Tune() error {
	e := d.set(dtvTune, 0)
	if e != 0 {
		return Error{"tune", "", e}
	}
	return nil
}

func (d Device) Clear() error {
	e := d.set(dtvClear, 0)
	if e != 0 {
		return Error{"clear", "", e}
	}
	return nil

}

func (d Device) Frequency() (uint32, error) {
	f, e := d.get(dtvFrequency)
	if e != 0 {
		return 0, Error{"get", "frequency", e}
	}
	return f, nil
}

func (d Device) SetFrequency(f uint32) error {
	e := d.set(dtvFrequency, f)
	if e != 0 {
		return Error{"set", "frequency", e}
	}
	return nil

}

func (d Device) Modulation() (dvb.Modulation, error) {
	m, e := d.get(dtvModulation)
	if e != 0 {
		return 0, Error{"get", "modulation", e}
	}
	return dvb.Modulation(m), nil
}

func (d Device) SetModulation(m dvb.Modulation) error {
	e := d.set(dtvModulation, uint32(m))
	if e != 0 {
		return Error{"set", "modulation", e}
	}
	return nil

}

// Bandwidth returns bandwidth in Hz
func (d Device) Bandwidth() (uint32, error) {
	b, e := d.get(dtvBandwidthHz)
	if e != 0 {
		return 0, Error{"get", "bandwidth", e}
	}
	return b, nil

}

func (d Device) SetBandwidth(hz uint32) error {
	e := d.set(dtvBandwidthHz, hz)
	if e != 0 {
		return Error{"set", "bandwidth", e}
	}
	return nil
}

func (d Device) Inversion() (dvb.Inversion, error) {
	i, e := d.get(dtvInversion)
	if e != 0 {
		return 0, Error{"get", "inversion", e}
	}
	return dvb.Inversion(i), nil
}

func (d Device) SetInversion(i dvb.Inversion) error {
	e := d.set(dtvInversion, uint32(i))
	if e != 0 {
		return Error{"set", "inversion", e}
	}
	return nil
}

func (d Device) SymbolRate() (uint32, error) {
	sr, e := d.get(dtvSymbolRate)
	if e != 0 {
		return 0, Error{"get", "symbol rate", e}
	}
	return sr, nil
}

func (d Device) SetSymbolRate(bd uint32) error {
	e := d.set(dtvSymbolRate, bd)
	if e != 0 {
		return Error{"set", "symbol rate", e}
	}
	return nil
}

func (d Device) InnerFEC() (dvb.CodeRate, error) {
	r, e := d.get(dtvInnerFEC)
	if e != 0 {
		return 0, Error{"get", "inner fec", e}
	}
	return dvb.CodeRate(r), nil
}

func (d Device) SetInnerFEC(r dvb.CodeRate) error {
	e := d.set(dtvInnerFEC, uint32(r))
	if e != 0 {
		return Error{"set", "inner fec", e}
	}
	return nil
}

func (d Device) Voltage() (Voltage, error) {
	v, e := d.get(dtvVoltage)
	if e != 0 {
		return 0, Error{"get", "voltage", e}
	}
	return Voltage(v), nil
}

func (d Device) SetVoltage(v Voltage) error {
	e := d.set(dtvVoltage, uint32(v))
	if e != 0 {
		return Error{"set", "voltage", e}
	}
	return nil
}

func (d Device) Tone() (Tone, error) {
	t, e := d.get(dtvTone)
	if e != 0 {
		return 0, Error{"get", "tone", e}
	}
	return Tone(t), nil
}

func (d Device) SetTone(t Tone) error {
	e := d.set(dtvTone, uint32(t))
	if e != 0 {
		return Error{"set", "tone", e}
	}
	return nil

}

func (d Device) Pilot() (dvb.Pilot, error) {
	p, e := d.get(dtvPilot)
	if e != 0 {
		return 0, Error{"get", "rolloff", e}
	}
	return dvb.Pilot(p), nil
}

func (d Device) SetPilot(p dvb.Pilot) error {
	e := d.set(dtvPilot, uint32(p))
	if e != 0 {
		return Error{"set", "pilot", e}
	}
	return nil
}

func (d Device) Rolloff() (dvb.Rolloff, error) {
	r, e := d.get(dtvRolloff)
	if e != 0 {
		return 0, Error{"get", "rolloff", e}
	}
	return dvb.Rolloff(r), nil
}

func (d Device) SetRolloff(r dvb.Rolloff) error {
	e := d.set(dtvRolloff, uint32(r))
	if e != 0 {
		return Error{"set", "rolloff", e}
	}
	return nil
}

func (d Device) CodeRateHP() (dvb.CodeRate, error) {
	r, e := d.get(dtvCodeRateHP)
	if e != 0 {
		return 0, Error{"get", "code rate HP", e}
	}
	return dvb.CodeRate(r), nil
}

func (d Device) SetCodeRateHP(r dvb.CodeRate) error {
	e := d.set(dtvCodeRateHP, uint32(r))
	if e != 0 {
		return Error{"set", "code rate HP", e}
	}
	return nil
}

func (d Device) CodeRateLP() (dvb.CodeRate, error) {
	r, e := d.get(dtvCodeRateLP)
	if e != 0 {
		return 0, Error{"get", "code rate LP", e}
	}
	return dvb.CodeRate(r), nil
}

func (d Device) SetCodeRateLP(r dvb.CodeRate) error {
	e := d.set(dtvCodeRateLP, uint32(r))
	if e != 0 {
		return Error{"set", "code rate LP", e}
	}
	return nil
}

func (d Device) Guard() (dvb.Guard, error) {
	g, e := d.get(dtvGuardInterval)
	if e != 0 {
		return 0, Error{"get", "guard interval", e}
	}
	return dvb.Guard(g), nil
}

func (d Device) SetGuard(g dvb.Guard) error {
	e := d.set(dtvGuardInterval, uint32(g))
	if e != 0 {
		return Error{"set", "guard interval", e}
	}
	return nil
}

func (d Device) TxMode() (dvb.TxMode, error) {
	m, e := d.get(dtvTransmissionMode)
	if e != 0 {
		return 0, Error{"get", "transmission mode", e}
	}
	return dvb.TxMode(m), nil
}

func (d Device) SetTxMode(m dvb.TxMode) error {
	e := d.set(dtvTransmissionMode, uint32(m))
	if e != 0 {
		return Error{"set", "transmission mode", e}
	}
	return nil
}

func (d Device) Hierarchy() (dvb.Hierarchy, error) {
	h, e := d.get(dtvHierarchy)
	if e != 0 {
		return 0, Error{"get", "hierarchy", e}
	}
	return dvb.Hierarchy(h), nil
}

func (d Device) SetHierarchy(h dvb.Hierarchy) error {
	e := d.set(dtvHierarchy, uint32(h))
	if e != 0 {
		return Error{"set", "hierarchy", e}
	}
	return nil

}
