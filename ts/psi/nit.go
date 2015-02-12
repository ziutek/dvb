package psi

import (
	"github.com/ziutek/dvb"
)

type NIT Table

func NewNIT() *NIT {
	return (*NIT)(NewTable(ISOSectionMaxLen))
}

func (nit *NIT) t() *Table {
	return (*Table)(nit)
}

func (nit *NIT) Version() byte {
	return nit.t().Version()
}

func (nit *NIT) Current() bool {
	return nit.t().Current()
}

func (nit *NIT) NetId() uint16 {
	return nit.t().TableIdExt()
}

func (nit *NIT) s() Section {
	ss := nit.t().Sections()
	if len(ss) == 0 {
		panic("NIT doesn't contain valid data")
	}
	return ss[0]
}

var ErrNITSectionLen = dvb.TemporaryError("incorrect NIT section length")

// Update reads next NIT from r
func (nit *NIT) Update(r SectionReader, actualMux bool, current bool) error {
	tableId := byte(0x41)
	if actualMux {
		tableId = 0x40
	}
	return nit.t().Update(r, tableId, true, current)
}

// Descriptors returns network descriptors list
func (nit *NIT) Descriptors() TableDescriptors {
	return nit.t().Descriptors(0)
}

func (nit *NIT) MuxInfo() MuxInfoList {
	return MuxInfoList{ss: nit.t().Sections()}
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
	ss   []Section
	data []byte
}

func (il MuxInfoList) IsEmpty() bool {
	return len(il.ss) == 0 && len(il.data) == 0
}

// Pop returns first multiplex information element in i and remaining
// elements in ril.  If there is no more elements then len(ril) == 0. If an
// error occurs i == nil.
func (il MuxInfoList) Pop() (i MuxInfo, ril MuxInfoList) {
	if len(il.data) == 0 {
		if len(il.ss) == 0 {
			return
		}
		il.data = il.ss[0].Data()
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
		il.ss = il.ss[1:]
	}
	if len(il.data) < 6 {
		return
	}
	l := int(decodeU16(il.data[4:6])&0x0fff) + 6
	if len(il.data) < l {
		return
	}
	i = MuxInfo(il.data[:l])
	ril.ss = il.ss
	ril.data = il.data[l:]
	return
}
