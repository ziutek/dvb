package psi

import (
	"github.com/ziutek/dvb"
)

type NIT Table

func (nit NIT) Version() int8 {
	return Table(nit).Version()
}

func (nit NIT) Current() bool {
	return Table(nit).Current()
}

func (nit NIT) NetId() uint16 {
	return Table(nit).TableIdExt()
}

func (nit NIT) s() Section {
	if len(nit) == 0 {
		panic("NIT doesn't contain valid data")
	}
	return nit[0]
}

var ErrNITSectionLen = dvb.TemporaryError("incorrect NIT section length")

// Update reads next NIT from r
func (nit *NIT) Update(r SectionReader, actualMux bool, current bool) error {
	tableId := byte(0x41)
	if actualMux {
		tableId = 0x40
	}
	return (*Table)(nit).Update(r, tableId, true, current, ISOSectionMaxLen)
}

// Descriptors returns network descriptors list
func (nit NIT) Descriptors() TableDescriptors {
	// BUG: Network descriptros can be in more than one (first) section. They
	// should be located in the first sections of table and they should end
	// before first not empty transport stream looop.
	return Table(nit).Descriptors(0)
}

func (nit NIT) MuxInfo() MuxInfoList {
	return MuxInfoList{Table(nit).Cursor()}
}

var (
	nitCfgActual = &TableConfig{
		TableId:       0x40,
		GenericSyntax: true,
		PrivateSyntax: true,
		SectionMaxLen: ISOSectionMaxLen,
		NumLenFields:  2,
	}
	nitCfgOther = &TableConfig{
		TableId:       0x41,
		GenericSyntax: true,
		PrivateSyntax: true,
		SectionMaxLen: ISOSectionMaxLen,
		NumLenFields:  2,
	}
)

// AppendNetDescriptors appends network descriptors ds to nit. It can be called
// only before first AppendMuxInfo call.
func (nit *NIT) AppendNetDescriptor(ds ...Descriptor) {
	for _, d := range ds {
		data := (*Table)(nit).Alloc(len(d), nitCfgActual, 0, nil)
		copy(data, d)
	}
}

// AppendMuxInfos appends information about transport stream to nit.
func (nit *NIT) AppendMuxInfo(mis ...MuxInfo) {
	for _, mi := range mis {
		data := (*Table)(nit).Alloc(len(mi), nitCfgActual, 1, nil)
		copy(data, mi)
	}
}

func (nit NIT) Close(netid uint16, actualMux, current bool, version int8) {
	nitcfg := nitCfgOther
	if actualMux {
		nitcfg = nitCfgActual
	}
	Table(nit).Close(nitcfg, netid, current, version)
}

type MuxInfo []byte

func (mi MuxInfo) MuxId() uint16 {
	return decodeU16(mi[0:2])
}

func (mi MuxInfo) SetMuxId(tsid uint16) {
	encodeU16(mi[0:2], tsid)
}

func (mi MuxInfo) OrgNetId() uint16 {
	return decodeU16(mi[2:4])
}

func (mi MuxInfo) SetOrgNetId(onid uint16) {
	encodeU16(mi[2:4], onid)
}

func (mi MuxInfo) Descriptors() DescriptorList {
	return DescriptorList(mi[6:])
}

type MuxInfoList struct {
	TableCursor
}

// Pop returns first MuxInfo element from il. If there is no more data to read
// Pop returns empty ServiceInfoList. If an error occurs it returns nil MuxInfo
// and non-empty MuxInfoList.
func (il MuxInfoList) Pop() (MuxInfo, MuxInfoList) {
	if len(il.Data) == 0 {
		if len(il.Tab) == 0 {
			return nil, il
		}
		il.TableCursor = il.NextSection()
		// Skip network descriptors
		if len(il.Data) < 2 {
			return nil, il
		}
		n := int(decodeU16(il.Data[0:2])&0x0fff) + 2
		if len(il.Data) < n {
			return nil, il
		}
		il.Data = il.Data[n:]
		// Decode transport_stream_loop_length
		if len(il.Data) < 2 {
			return nil, il
		}
		n = int(decodeU16(il.Data[0:2])&0x0fff) + 2
		if len(il.Data) < n {
			return nil, il
		}
		il.Data = il.Data[2:n]
	}
	if len(il.Data) < 6 {
		return nil, il
	}
	n := int(decodeU16(il.Data[4:6])&0x0fff) + 6
	if len(il.Data) < n {
		return nil, il
	}
	data := il.Data[:n]
	il.Data = il.Data[n:]
	return data, il
}

func MakeMuxInfo() MuxInfo {
	mi := make(MuxInfo, 6)
	mi[4] = 0xf0
	return mi
}

func (mi MuxInfo) descrLoopLen() int {
	return loopLen(mi[4:6])
}

func (mi MuxInfo) setDescrLoopLen(n int) {
	setLoopLen(mi[4:6], n)
}

func (mi MuxInfo) ClearDescriptors() {
	mi.setDescrLoopLen(0)
}

func (mi *MuxInfo) AppendDescriptor(ds ...Descriptor) {
	n := mi.descrLoopLen()
	for _, d := range ds {
		*mi = append((*mi)[:6+n], d...)
		n += len(d)
	}
	mi.setDescrLoopLen(n)
}
