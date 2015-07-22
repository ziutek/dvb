package psi

import (
	"github.com/ziutek/dvb"
	"strconv"
)

type ServiceType byte

const (
	DigitalTelevisionService = ServiceType(0x01)
	DigitalRadioSoundService = ServiceType(0x02)
	TeletextService          = ServiceType(0x03)
	NVODReferenceService     = ServiceType(0x04)
	NVODTimeShiftedService   = ServiceType(0x05)
	MosaicService            = ServiceType(0x06)
	FMRadioService           = ServiceType(0x07)
	DVBSRMService            = ServiceType(0x08)
	// 0x09
	AdvancedCodecDigitalRadioSoundService = ServiceType(0x0a)
	AdvancedCodecMosaicService            = ServiceType(0x0b)
	DataBroadcastService                  = ServiceType(0x0c)
	// 0x0d
	RCSMapService                   = ServiceType(0x0e)
	RCSFLSService                   = ServiceType(0x0f)
	DVBMHPService                   = ServiceType(0x10)
	MPEG2HDDigitalTelevisionService = ServiceType(0x11)
	// 0x12
	// 0x13
	// 0x14
	// 0x15
	AdvancedCodecSDDigitalTelevisionService = ServiceType(0x16)
	AdvancedCodecSDNVODTimeShiftedService   = ServiceType(0x17)
	AdvancedCodecSDNVODReferenceService     = ServiceType(0x18)
	AdvancedCodecHDDigitalTelevisionService = ServiceType(0x19)
	AdvancedCodecHDNVODTimeShiftedService   = ServiceType(0x1a)
	AdvancedCodecHDNVODReferenceService     = ServiceType(0x1b)
)

var stn = []string{
	"digital television",
	"digital radio sound",
	"teletext",
	"NVOD reference",
	"NVOD time-shifted",
	"mosaic",
	"FM radio",
	"DVB SRM",
	"reserved",
	"advanced codec digital radio sound",
	"advanced codec mosaic",
	"data broadcast",
	"reserved",
	"RCS Map",
	"RCS FLS",
	"DVB MHP",
	"MPEG-2 HD digital television",
	"reserved",
	"reserved",
	"reserved",
	"reserved",
	"advanced codec SD digital television",
	"advanced codec SD NVOD time-shifted",
	"advanced codec SD NVOD reference",
	"advanced codec HD digital television",
	"advanced codec HD NVOD time-shifted",
	"advanced codec HD NVOD reference",
}

func (t ServiceType) String() string {
	if t == 0 || t == 0xff || t > AdvancedCodecHDNVODReferenceService && t <= 0x7F {
		return "reserved"
	}
	if t > 0x7F {
		return "user defined"
	}
	return stn[t-1]
}

type ServiceDescriptor struct {
	Type         ServiceType
	ProviderName []byte
	ServiceName  []byte
}

func ParseServiceDescriptor(d Descriptor) (sd ServiceDescriptor, ok bool) {
	if d.Tag() != ServiceTag {
		return
	}
	data := d.Data()
	if len(data) < 2 {
		return
	}
	sd.Type = ServiceType(data[0])
	providerNameLen := data[1]
	data = data[2:]

	if len(data) < int(providerNameLen+1) {
		return
	}
	sd.ProviderName = data[:providerNameLen]
	serviceNameLen := data[providerNameLen]
	data = data[providerNameLen+1:]

	if len(data) < int(serviceNameLen) {
		return
	}
	sd.ServiceName = data[:serviceNameLen]

	ok = true
	return
}

/*func (sd ServiceDescriptor) Append(dl *DescriptorList)  {
	d := dl.Alloc(1 + 1 + len(sd.ProviderName) + 1 + len(sd.ServiceName)
	d.SetTag(ServiceTag)
	data := d.Data()
	data[0] = len(sd.ProviderName)
	copy(data[1:],
}*/

func MakeServiceDescriptor(typ ServiceType, provider, service string) Descriptor {
	p := EncodeText(provider)
	s := EncodeText(service)
	d := MakeDescriptor(ServiceTag, 1+1+len(p)+1+len(s))
	data := d.Data()[:0]
	data = append(data, byte(typ))
	data = append(data, byte(len(p)))
	data = append(data, p...)
	data = append(data, byte(len(s)))
	data = append(data, s...)
	return d
}

type NetworkNameDescriptor []byte

func ParseNetworkNameDescriptor(d Descriptor) (nnd NetworkNameDescriptor, ok bool) {
	if d.Tag() != NetworkNameTag {
		return
	}
	nnd = NetworkNameDescriptor(d.Data())
	ok = true
	return
}

func MakeNetworkNameDescriptor(netname string) Descriptor {
	nn := EncodeText(netname)
	d := MakeDescriptor(NetworkNameTag, len(nn))
	copy(d.Data(), nn)
	return d
}

type ServiceListDescriptor []byte

func ParseServiceListDescriptor(d Descriptor) (sld ServiceListDescriptor, ok bool) {
	if d.Tag() != ServiceListTag {
		return
	}
	sld = ServiceListDescriptor(d.Data())
	ok = true
	return
}

// Pop returns first (sid, typ) pair from d. Remaining pairs are returned in rd.
// If there is no more pairs to read len(rd) == 0. If an error occurs rd = nil
func (d ServiceListDescriptor) Pop() (sid uint16, typ ServiceType, rd ServiceListDescriptor) {
	if len(d) < 3 {
		return
	}
	sid = decodeU16(d[0:2])
	typ = ServiceType(d[2])
	rd = d[3:]
	return
}

type TerrestrialDeliverySystemDescriptor struct {
	Freq          int64 // center frequency [Hz]
	Bandwidth     int32 // bandwidth [Hz]
	HighPrio      bool
	TimeSlicing   bool
	MPEFEC        bool
	Constellation dvb.Modulation
	Hierarchy     dvb.Hierarchy
	CodeRateHP    dvb.CodeRate
	CodeRateLP    dvb.CodeRate
	Guard         dvb.Guard
	TxMode        dvb.TxMode
	OtherFreq     bool
}

func decodeCodeRate(b byte) dvb.CodeRate {
	switch b {
	case 0:
		return dvb.FEC12
	case 1:
		return dvb.FEC23
	case 2:
		return dvb.FEC34
	case 3:
		return dvb.FEC56
	case 4:
		return dvb.FEC78
	case 5, 6:
		return dvb.FECAuto

	}
	return dvb.FECNone
}

func ParseTerrestrialDeliverySystemDescriptor(d Descriptor) (tds TerrestrialDeliverySystemDescriptor, ok bool) {
	if d.Tag() != TerrestrialDeliverySystemTag {
		return
	}
	data := d.Data()
	if len(data) < 11 {
		return
	}
	tds.Freq = int64(decodeU32(data[0:4])) * 10
	switch data[4] >> 5 {
	case 0:
		tds.Bandwidth = 8e6
	case 1:
		tds.Bandwidth = 7e6
	case 2:
		tds.Bandwidth = 6e6
	case 3:
		tds.Bandwidth = 5e6
	default:
		return
	}
	tds.HighPrio = data[4]&0x10 != 0
	tds.TimeSlicing = data[4]&0x08 == 0
	tds.MPEFEC = data[4]&0x04 == 0
	switch data[5] >> 6 {
	case 0:
		tds.Constellation = dvb.QPSK
	case 1:
		tds.Constellation = dvb.QAM16
	case 2:
		tds.Constellation = dvb.QAM64
	default:
		return
	}
	tds.Hierarchy = dvb.Hierarchy((data[5] >> 3) & 0x07)
	if tds.Hierarchy > dvb.Hierarchy4 {
		return
	}
	tds.CodeRateHP = decodeCodeRate(data[5] & 0x07)
	if tds.CodeRateHP == dvb.FECAuto || tds.CodeRateHP == dvb.FECNone {
		return
	}
	tds.CodeRateLP = decodeCodeRate(data[6] >> 5)
	if tds.CodeRateLP == dvb.FECAuto {
		return
	}
	tds.Guard = dvb.Guard((data[6] >> 3) & 0x03)
	tds.TxMode = dvb.TxMode((data[6] >> 1) & 0x03)
	if tds.TxMode > dvb.TxMode8k {
		return
	}
	tds.OtherFreq = data[6]&0x01 != 0
	ok = true
	return
}

func encodeCodeRate(fec dvb.CodeRate) byte {
	switch fec {
	case dvb.FEC12:
		return 0
	case dvb.FEC23:
		return 1
	case dvb.FEC34:
		return 2
	case dvb.FEC56:
		return 3
	case dvb.FEC78:
		return 4
	}
	return 7
}

func (tds TerrestrialDeliverySystemDescriptor) Make() Descriptor {
	d := MakeDescriptor(TerrestrialDeliverySystemTag, 11)
	data := d.Data()
	encodeU32(data[0:4], uint32((tds.Freq+5)/10))
	var b byte
	switch tds.Bandwidth {
	case 8e6:
		b = 0 << 5
	case 7e6:
		b = 1 << 5
	case 6e6:
		b = 2 << 5
	case 5e6:
		b = 3 << 5
	default: // Unknown
		b = 7 << 5
	}
	if tds.HighPrio || tds.Hierarchy == dvb.HierarchyNone {
		b |= 0x10
	}
	if !tds.TimeSlicing {
		b |= 0x08
	}
	if !tds.MPEFEC {
		b |= 0x04
	}
	b |= 0x03 // Reserved
	data[4] = b
	switch tds.Constellation {
	case dvb.QPSK:
		b = 0 << 6
	case dvb.QAM16:
		b = 1 << 6
	case dvb.QAM64:
		b = 2 << 6
	default: // Unknown
		b = 3 << 6
	}
	b |= byte(tds.Hierarchy&0x07) << 3
	b |= encodeCodeRate(tds.CodeRateHP)
	data[5] = b
	b = encodeCodeRate(tds.CodeRateLP) << 5
	b |= byte(tds.Guard&0x03) << 3
	b |= byte(tds.TxMode&0x03) << 1
	if tds.OtherFreq {
		b |= 0x01
	}
	data[6] = b
	// Reserved
	data[7] = 0xff
	data[8] = 0xff
	data[9] = 0xff
	data[10] = 0xff
	return d
}

type CAS uint16

var casn = map[CAS]string{
	0x0100: "Mediaguard",

	0x0b00: "Conax",
	0x0b01: "Conax",
	0x0b02: "Conax",
	0x0b03: "Conax",
	0x0b04: "Conax",
	0x0b05: "Conax",
	0x0b06: "Conax",
	0x0b07: "Conax",
	0x0baa: "Conax",

	0x0d00: "Cryptoworks",
	0x0d02: "Cryptoworks",
	0x0d03: "Cryptoworks",
	0x0d05: "Cryptoworks",
	0x0d07: "Cryptoworks",
	0x0d20: "Cryptoworks",

	0x0500: "Viaccess",

	0x0602: "Irdeto",
	0x0604: "Irdeto",
	0x0606: "Irdeto",
	0x0608: "Irdeto",
	0x0622: "Irdeto",
	0x0626: "Irdeto",
	0x0664: "Irdeto",
	0x0614: "Irdeto",
	0x0692: "Irdeto",

	0x0911: "Videoguard",
	0x0919: "Videoguard",
	0x0960: "Videoguard",
	0x0961: "Videoguard",
	0x093b: "Videoguard",
	0x0963: "Videoguard",
	0x09AC: "Videoguard",
	0x0927: "Videoguard",

	0x0700: "DigiCipher2",

	0x0E00: "PowerVu",

	0x1702: "Nagravision",
	0x1722: "Nagravision",
	0x1762: "Nagravision",
	0x1800: "Nagravision",
	0x1801: "Nagravision",
	0x1810: "Nagravision",
	0x1830: "Nagravision",

	0x4AEA: "Cryptoguard",
}

func (cas CAS) String() string {
	name := casn[cas]
	s := strconv.FormatUint(uint64(cas), 16)
	if name == "" {
		return s
	}
	return name + "(" + s + ")"
}

type CADescriptor struct {
	Sys CAS
	Pid int16
}

func ParseCADescriptor(d Descriptor) (cad CADescriptor, ok bool) {
	if d.Tag() != CATag {
		return
	}
	data := d.Data()
	if len(data) < 4 {
		return
	}
	cad.Sys = CAS(decodeU16(data[0:2]))
	cad.Pid = int16(decodeU16(data[2:4]) & 0x1fff)
	ok = true
	return
}

type ISO639LangDescriptor []byte

func ParseISO639LangDescriptor(d Descriptor) (ld ISO639LangDescriptor, ok bool) {
	if d.Tag() != ISO639LangTag {
		return
	}
	return ISO639LangDescriptor(d.Data()), true
}

type AudioType byte

const (
	UndefinedAudio AudioType = iota
	CleanEffectsAudio
	HearingImpairedAudio
	VisualImpairedCommentaryAudio
)

var atn = []string{
	"undefined",
	"clean effects",
	"hearing impaired",
	"visual impaired commentary",
}

func (t AudioType) String() string {
	if t > VisualImpairedCommentaryAudio {
		return "reserved"
	}
	return atn[t]
}

type ISO639LangCode uint32

// Pop returns first (lc, at) pair from d. Remaining pairs are returned in rd.
// If there is no more pairs to read len(rd) == 0. If an error occurs rd = nil
func (d ISO639LangDescriptor) Pop() (lc ISO639LangCode, at AudioType, rd ISO639LangDescriptor) {
	if len(d) < 4 {
		return
	}
	lc = ISO639LangCode(decodeU24(d[0:3]))
	at = AudioType(d[3])
	rd = d[4:]
	return
}

type StreamIdentifierDescriptor byte

func ParseStreamIdentifierDescriptor(d Descriptor) (sid StreamIdentifierDescriptor, ok bool) {
	if d.Tag() != StreamIdentifierTag {
		return
	}
	data := d.Data()
	if len(data) != 1 {
		return
	}
	return StreamIdentifierDescriptor(data[0]), true
}
