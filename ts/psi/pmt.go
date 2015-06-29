package psi

import (
	"github.com/ziutek/dvb"
)

type PMT Section

func (p PMT) Version() int8 {
	return Section(p).Version()
}

func (p PMT) SetVersion(v int8) {
	Section(p).SetVersion(v)
}

func (p PMT) Current() bool {
	return Section(p).Current()
}

func (p PMT) MakeCRC() {
	Section(p).MakeCRC()
}

func (p PMT) ProgId() uint16 {
	return Section(p).TableIdExt()
}

func (p PMT) PidPCR() int16 {
	return int16(decodeU16(Section(p).Data()[0:2]) & 0x1fff)
}

func (p PMT) SetPidPCR(pid int16) {
	checkPid(pid)
	d := Section(p).Data()
	d[0] = d[0]&0xe0 | byte(pid>>8)
	d[1] = byte(pid)
}

func (p PMT) progInfoLen() int {
	return int(decodeU16(Section(p).Data()[2:4]) & 0x0fff)
}

func (p PMT) ProgramDescriptors() DescriptorList {
	return DescriptorList(Section(p).Data()[4 : 4+p.progInfoLen()])
}

func (p PMT) ESInfo() ESInfoList {
	return ESInfoList(Section(p).Data()[4+p.progInfoLen():])
}

var (
	ErrPMTSectionSyntax = dvb.TemporaryError("incorrect PMT section syntax")
	ErrPMTProgInfoLen   = dvb.TemporaryError("incorrect PMT program info length")
)

// AsPMT returns s as PMT or error if s isn't PMT section. This works because
// PMT should fit in one section (other tables occupy multiple sections.
func AsPMT(s Section) (PMT, error) {
	if s.TableId() != 2 || !s.GenericSyntax() || s.Number() != 0 ||
		s.LastNumber() != 0 {
		return nil, ErrPMTSectionSyntax
	}
	p := PMT(s)
	if p.progInfoLen()+4 > len(s.Data()) {
		return nil, ErrPMTProgInfoLen
	}
	return p, nil
}

// Update reads one section into p and runs AsPMT to check its syntax.
func (p PMT) Update(r SectionReader) error {
	err := r.ReadSection(p.Section())
	if err != nil {
		return err
	}
	_, err = AsPMT(p.Section())
	return err
}

func (p PMT) Section() Section {
	return Section(p)
}

type ESInfo []byte

func (i ESInfo) Type() StreamType {
	return StreamType(i[0])
}

func (i ESInfo) Pid() int16 {
	return int16(decodeU16(i[1:3]) & 0x1fff)
}

func (i ESInfo) SetPid(pid int16) {
	checkPid(pid)
	i[1] = i[1]&0xe0 | byte(pid>>8)
	i[2] = byte(pid)
}

func (i ESInfo) esInfoLen() uint16 {
	return decodeU16(i[3:5]) & 0x0fff
}

func (i ESInfo) Descriptors() DescriptorList {
	l := decodeU16(i[3:5])&0x0fff + 5
	return DescriptorList(i[5:l])
}

type ESInfoList []byte

// Pop returns first elementary stream information element in i and remaining
// elements in ril.  If there is no more elements then len(ril) == 0. If an
// error occurs i == nil.
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
