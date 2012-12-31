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
	dtvFreqency
	dtvModulation
	dtvBandwidthHz
	dtvInversion
	dtvDiseqcMaster
	dtvSymbolRate
	dtvInnerFec
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

// Device delivery sytem
const (
	sysUndefined = iota
	sysDVBCAnnexAC
	sysDVBCAnnexB
	sysDVBT
	sysDSS
	sysDVBS
	sysDVBS2
	sysDVBH
	sysISDBT
	sysISDBS
	sysISDBC
	sysATSC
	sysATSCMH
	sysDMBTH
	sysCMMB
	sysDAB
	sysDVBT2
	sysTURBO
)

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


