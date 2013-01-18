package psi

type PMT struct {
	s Section
}

func NewPMT() PMT {
	return PMT{NewSection(ISOSectionMaxLen)}
}

func (p PMT) Version() byte {
	return p.s.Version()
}

func (p PMT) Current() bool {
	return p.s.Current()
}

func (p PMT) ProgId() uint16 {
	return p.s.TableIdExt()
}

func (p PMT) PCRPid() uint16 {
	return decodeU16(p.s.Data()[0:2]) & 0x1fff
}

func (p PMT) progInfoLen() int {
	return int(decodeU16(p.s.Data()[2:4]) & 0x0fff)
}

func (p PMT) ProgramInfo() DescriptorList {
	return DescriptorList(p.s.Data()[4 : 4+p.progInfoLen()])
}

func (p PMT) StreamsInfo() ESInfoList {
	return ESInfoList(p.s.Data()[4+p.progInfoLen():])
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
	i = ESInfo(il[5:l])
	ril = il[l:]
	return
}

// Append adds i to the end of the il. It works like Go append function so need
// to be used in this way:
//     il = il.Append(i)
func (il ESInfoList) Append(i ESInfo) ESInfoList {
	return append(il, i...)
}
