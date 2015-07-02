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
	return Table(nit).Descriptors(0)
}

func (nit NIT) MuxInfo() MuxInfoList {
	return MuxInfoList{Table(nit).Cursor()}
}

type MuxInfo []byte

func (i MuxInfo) MuxId() uint16 {
	return decodeU16(i[0:2])
}

func (i MuxInfo) OrgNetId() uint16 {
	return decodeU16(i[2:4])
}

func (i MuxInfo) Descriptors() DescriptorList {
	return DescriptorList(i[6:])
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