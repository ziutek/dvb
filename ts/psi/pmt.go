package psi

import (
	"github.com/ziutek/dvb"
)

type PMT Section

func NewPMT() PMT {
	return PMT(NewSection(ISOSectionMaxLen))
}

func (p PMT) Version() byte {
	return Section(p).Version()
}

func (p PMT) Current() bool {
	return Section(p).Current()
}

func (p PMT) ProgId() uint16 {
	return Section(p).TableIdExt()
}

func (p PMT) PidPCR() uint16 {
	return decodeU16(Section(p).Data()[0:2]) & 0x1fff
}

func (p PMT) progInfoLen() int {
	return int(decodeU16(Section(p).Data()[2:4]) & 0x0fff)
}

func (p PMT) ProgramInfo() DescriptorList {
	return DescriptorList(Section(p).Data()[4 : 4+p.progInfoLen()])
}

func (p PMT) ESInfo() ESInfoList {
	return ESInfoList(Section(p).Data()[4+p.progInfoLen():])
}

var (
	ErrPMTSectionSyntax = dvb.TemporaryError("incorrect PMT section syntax")
	ErrPMTProgInfoLen   = dvb.TemporaryError("incorrect PMT program info length")
)

type PMTDecoder struct {
	r SectionReader
}

func NewPMTDecoder(r SectionReader) *PMTDecoder {
	return &PMTDecoder{r}
}

func (d *PMTDecoder) SetSectionReader(r SectionReader) {
	d.r = r
}

func (d *PMTDecoder) ReadPMT(p PMT) error {
	s := Section(p)
	err := d.r.ReadSection(s)
	if err != nil {
		return err
	}
	if s.TableId() != 2 || !s.GenericSyntax() || s.Number() != 0 ||
		s.LastNumber() != 0 {
		return ErrPMTSectionSyntax
	}
	if p.progInfoLen()+4 > len(s.Data()) {
		return ErrPMTProgInfoLen
	}
	return nil
}

type ESInfo []byte

func (i ESInfo) Type() StreamType {
	return StreamType(i[0])
}

func (i ESInfo) Pid() uint16 {
	return decodeU16(i[1:3]) & 0x1fff
}

func (i ESInfo) esInfoLen() uint16 {
	return decodeU16(i[3:5]) & 0x0fff
}

func (i ESInfo) DescriptorList() DescriptorList {
	l := decodeU16(i[3:5])&0x0fff + 5
	return DescriptorList(i[5:l])
}

type ESInfoList []byte

// Pop returns first elementary stream information element in i and remaining
// elements in ril. If an error occurs it returns i == nil. If there is no more
// elements len(ril) == 0.
func (il ESInfoList) Pop() (i ESInfo, ril ESInfoList) {
	if len(il) < 5 {
		return
	}
	l := int(decodeU16(il[3:5])&0x0fff) + 5
	if len(il) < l {
		return
	}
	i = ESInfo(il[:l])
	ril = il[l:]
	return
}

// Append adds i to the end of the il. It works like Go append function so need
// to be used in this way:
//     il = il.Append(i)
func (il ESInfoList) Append(i ESInfo) ESInfoList {
	return append(il, i...)
}
