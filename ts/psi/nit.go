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
	return MuxInfoList{nit: nit}
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
	nit   NIT
	data []byte
}

func (il MuxInfoList) IsEmpty() bool {
	return len(il.nit) == 0 && len(il.data) == 0
}

// Pop returns first multiplex information element in i and remaining
// elements in ril.  If there is no more elements then len(ril) == 0. If an
// error occurs i == nil.
func (il MuxInfoList) Pop() (i MuxInfo, ril MuxInfoList) {
	if len(il.data) == 0 {
		if len(il.nit) == 0 {
			return
		}
		il.data = il.nit[0].Data()
		if len(il.data) < 2 {
			return
		}
		l := int(decodeU16(il.data[0:2])&0x0fff) + 2
		if len(il.data) < l {
			return
		}
		il.data = il.data[l:]
		if len(il.data) < 2 {
			return
		}
		l = int(decodeU16(il.data[0:2])&0x0fff) + 2
		if len(il.data) < l {
			return
		}
		il.data = il.data[2:l]
		il.nit = il.nit[1:]
	}
	if len(il.data) < 6 {
		return
	}
	l := int(decodeU16(il.data[4:6])&0x0fff) + 6
	if len(il.data) < l {
		return
	}
	i = MuxInfo(il.data[:l])
	ril.nit = il.nit
	ril.data = il.data[l:]
	return
}
