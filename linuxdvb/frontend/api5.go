package frontend

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/ziutek/dvb"
)

type Device struct {
	file *os.File
}

func Open(filepath string) (d Device, err error) {
	for {
		d.file, err = os.OpenFile(filepath, os.O_RDWR, 0)
		if err == nil || err.(*os.PathError).Err != syscall.EINTR {
			break
		}
	}
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
	dtvUndefined                   cmd = 0
	dtvTune                        cmd = 1
	dtvClear                       cmd = 2
	dtvFrequency                   cmd = 3
	dtvModulation                  cmd = 4
	dtvBandwidthHz                 cmd = 5
	dtvInversion                   cmd = 6
	dtvDiseqcMaster                cmd = 7
	dtvSymbolRate                  cmd = 8
	dtvInnerFEC                    cmd = 9
	dtvVoltage                     cmd = 10
	dtvTone                        cmd = 11
	dtvPilot                       cmd = 12
	dtvRolloff                     cmd = 13
	dtvDiseqcSlaveReply            cmd = 14
	dtvFeCapabilityCount           cmd = 15
	dtvFeCapability                cmd = 16
	dtvDeliverySystem              cmd = 17
	dtvISDBTPartialReception       cmd = 18
	dtvISDBTSoundBroadcasting      cmd = 19
	dtvISDBTSBSubchannelId         cmd = 20
	dtvISDBTSBSegmentIdx           cmd = 21
	dtvISDBTSBSegmentCount         cmd = 22
	dtvISDBTLayeraFEC              cmd = 23
	dtvISDBTLayeraModulation       cmd = 24
	dtvISDBTLayeraSegmentCount     cmd = 25
	dtvISDBTLayeraTimeInterleaving cmd = 26
	dtvISDBTLayerbFEC              cmd = 27
	dtvISDBTLayerbModulation       cmd = 28
	dtvISDBTLayerbSegmentCount     cmd = 29
	dtvISDBTLayerbTimeInterleaving cmd = 30
	dtvISDBTLayercFEC              cmd = 31
	dtvISDBTLayercModulation       cmd = 32
	dtvISDBTLayercSegmentCount     cmd = 33
	dtvISDBTLayercTimeInterleaving cmd = 34
	dtvAPIVersion                  cmd = 35
	dtvCodeRateHP                  cmd = 36
	dtvCodeRateLP                  cmd = 37
	dtvGuardInterval               cmd = 38
	dtvTransmissionMode            cmd = 39
	dtvHierarchy                   cmd = 40
	dtvISDBTLayerEnabled           cmd = 41
	dtvISDBSTSId                   cmd = 42
	dtvDVBT2PLPId                  cmd = 43
	dtvEnumDelSys                  cmd = 44

	dtvStatSignalStrength    cmd = 62
	dtvStatCNR               cmd = 63
	dtvStatPreErrorBitCount  cmd = 64
	dtvStatPreTotalBitCount  cmd = 65
	dtvStatPostErrorBitCount cmd = 66
	dtvStatPostTotalBitCount cmd = 67
	dtvStatErrorBlockCount   cmd = 68
	dtvStatTotalBlockCount   cmd = 69
)

type property struct {
	cmd      cmd
	reserved [3]uint32

	// union
	data         uint32
	bufData      [28]byte
	bufLen       uint32
	bufReserved1 [3]uint32
	bufReserved2 [unsafe.Sizeof(uintptr(0))]byte

	result int32
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

type Scale byte

const (
	ScaleNotAvailable Scale = iota
	ScaleDecibel
	ScaleRelative
	ScaleCounter
)

var scaleStr = []string{"(unk)", "dB", "%", ""}

func (s Scale) String() string {
	if s > ScaleCounter {
		return scaleStr[0]
	}
	return scaleStr[s]
}

type Param struct {
	scale Scale
	value int64
}

func (p Param) Scale() Scale {
	return p.scale
}

func (p Param) Decibel() float64 {
	return float64(p.value) * 0.001
}

func (p Param) Relative() float64 {
	return float64(p.value) * 100 / 0xffff
}

func (p Param) Counter() uint64 {
	return uint64(p.value)
}

func (p Param) Format(f fmt.State, _ rune) {
	switch p.scale {
	case ScaleDecibel:
		fmt.Fprintf(f, "%.3f dB", p.Decibel())
	case ScaleRelative:
		fmt.Fprintf(f, "%.3f %%", p.Relative())
	case ScaleCounter:
		fmt.Fprintf(f, "%d", p.Counter())
	default:
		fmt.Fprintf(f, "0x%x", p.Counter())
	}
}

type Stat struct {
	Signal     []Param // Signal strength level at the analog part of the tuner or of the demod.
	CNR        []Param // Signal to Noise ratio for the main carrier.
	PreErrBit  []Param // Number of bit errors before FEC on the inner coding block (Viterbi, LDPC, ...).
	PreTotBit  []Param // Number of bits received before the inner code block.
	PostErrBit []Param // Number of bit errors after FEC done by inner code block (Viterbi, LDPC, ...).
	PostTotBit []Param // Number of bits received after the inner coding.
	ErrBlk     []Param // Number of block errors after the outer FEC (Reed-Solomon, ...).
	TotBlk     []Param // Total number of blocks received.
}

type statProp struct {
	cmd      cmd
	reserved [3]uint32

	// union
	slen      byte
	stat      [4][9]byte
	reserved1 [11]byte
	reserved2 [unsafe.Sizeof(uintptr(0))]byte

	result int32
}

type statProps struct {
	num   uint32
	props *statProp
}

func getStateParam(p *statProp) []Param {
	ret := make([]Param, p.slen)
	for i := range ret {
		stat := &p.stat[i]
		ret[i].scale = Scale(stat[0])
		val := (*[8]byte)(unsafe.Pointer(&ret[i].value))
		copy(val[:], stat[1:])
	}
	return ret
}

func (d Device) Stat() (*Stat, error) {
	p := [...]statProp{
		{cmd: dtvStatSignalStrength},
		{cmd: dtvStatCNR},
		{cmd: dtvStatPreErrorBitCount},
		{cmd: dtvStatPreTotalBitCount},
		{cmd: dtvStatPostErrorBitCount},
		{cmd: dtvStatPostTotalBitCount},
		{cmd: dtvStatErrorBlockCount},
		{cmd: dtvStatTotalBlockCount},
	}
	ps := statProps{uint32(len(p)), &p[0]}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		d.Fd(),
		_FE_GET_PROPERTY,
		uintptr(unsafe.Pointer(&ps)),
	)
	if e != 0 {
		return nil, e
	}
	stat := &Stat{
		Signal:     getStateParam(&p[0]),
		CNR:        getStateParam(&p[1]),
		PreErrBit:  getStateParam(&p[2]),
		PreTotBit:  getStateParam(&p[3]),
		PostErrBit: getStateParam(&p[4]),
		PostTotBit: getStateParam(&p[5]),
		ErrBlk:     getStateParam(&p[6]),
		TotBlk:     getStateParam(&p[7]),
	}
	return stat, nil
}
