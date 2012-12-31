package frontend

import (
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
	dtvRollOff
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

func (f Device) set(c cmd, data uint32) error {
	p := property{cmd: c, data: data}
	ps := properties{1, &p}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_SET_PROPERTY,
		uintptr(unsafe.Pointer(&ps)),
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f Device) get(c cmd) (uint32, error) {
	p := property{cmd: c}
	ps := properties{1, &p}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		_FE_GET_PROPERTY,
		uintptr(unsafe.Pointer(&ps)),
	)
	if e != 0 {
		return 0, e
	}
	return p.data, nil
}

func (f Device) Version() (major, minor int, err error) {
	var u uint32
	u, err = f.get(dtvAPIVersion)
	if err != nil {
		return
	}
	major = int(u>>8) & 0xff
	minor = int(u) & 0xff
	return
}

// Device delivery sytem
type DeliverySystem uint32

const (
	SysUndefined = iota
	SysDVBCAnnexAC
	SysDVBCAnnexB
	SysDVBT
	SysDSS
	SysDVBS
	SysDVBS2
	SysDVBH
	SysISDBT
	SysISDBS
	SysISDBC
	SysATSC
	SysATSCMH
	SysDMBTH
	SysCMMB
	SysDAB
	SysDVBT2
	SysTURBO
)

func (ds DeliverySystem) String() string {
	dsn := []string{
		"Undefined",
		"DVB-C Annex AC",
		"DVB-C Annex B",
		"DVB-T",
		"DSS",
		"DVB-S",
		"DVB-S2",
		"DVB-H",
		"ISDB-T",
		"ISDB-S",
		"ISDB-C",
		"ATSC",
		"ATSC-MH",
		"DMBT-H",
		"CMMB",
		"DAB",
		"DVB-T2",
		"TURBO",
	}
	if ds > DeliverySystem(len(dsn)) {
		return "Unknown"
	}
	return dsn[ds]
}

func (f Device) DeliverySystem() (DeliverySystem, error) {
	ds, err := f.get(dtvDeliverySystem)
	return DeliverySystem(ds), err
}

func (f Device) SetDeliverySystem(d DeliverySystem) error {
	return f.set(dtvDeliverySystem, uint32(d))
}

func (f Device) Tune() error {
	return f.set(dtvTune, 0)
}

func (f Device) Clear() error {
	return f.set(dtvClear, 0)
}

func (f Device) Frequency() (uint32, error) {
	return f.get(dtvFrequency)
}

func (f Device) SetFrequency(freq uint32) error {
	return f.set(dtvFrequency, freq)
}

func (f Device) Modulation() (Modulation, error) {
	m, err := f.get(dtvModulation)
	return Modulation(m), err
}

func (f Device) SetModulation(m Modulation) error {
	return f.set(dtvModulation, uint32(m))
}

func (f Device) BandwidthHz() (uint32, error) {
	return f.get(dtvBandwidthHz)
}

func (f Device) SetBandwidthHz(bw uint32) error {
	return f.set(dtvBandwidthHz, bw)
}

func (f Device) Inversion() (Inversion, error) {
	i, err := f.get(dtvInversion)
	return Inversion(i), err
}

func (f Device) SetInversion(i Inversion) error {
	return f.set(dtvInversion, uint32(i))
}

func (f Device) SymbolRate() (uint32, error) {
	return f.get(dtvSymbolRate)
}

func (f Device) SetSymbolRate(bd uint32) error {
	return f.set(dtvSymbolRate, bd)
}

func (f Device) InnerFEC() (CodeRate, error) {
	r, err := f.get(dtvInnerFEC)
	return CodeRate(r), err
}

func (f Device) SetInnerFEC(r CodeRate) error {
	return f.set(dtvInnerFEC, uint32(r))
}
