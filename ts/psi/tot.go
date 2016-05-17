package psi

import (
	"time"

	"github.com/ziutek/dvb"
)

var ErrTOTSectionSyntax = dvb.TemporaryError("incorrect TOT section syntax")

func ParseTOT(s Section) (utc time.Time, tod LocalTimeOffsetDescriptor, err error) {
	slen := s.Len()
	if s.TableId() != 0x73 || s.GenericSyntax() || slen < 3+5+2+4 {
		err = ErrTOTSectionSyntax
		return
	}
	crc := decodeU32(s[slen-4 : slen])
	if mpegCRC32(s[0:slen-4]) != crc {
		err = ErrSectionCRC
		return
	}
	utc, err = decodeMJDUTC(s[3:8])
	if err != nil {
		return
	}
	l := int(decodeU16(s[8:10]) & 0x0fff)
	if slen-(3+5+2+4) != l {
		err = ErrTOTSectionSyntax
		return
	}
	dl := DescriptorList(s[10 : 10+l])
	for len(dl) != 0 {
		var d Descriptor
		d, dl = dl.Pop()
		if d == nil {
			err = ErrSectionSyntax
			break
		}
		var ok bool
		if tod, ok = ParseLocalTimeOffsetDescriptor(d); ok {
			break
		}
	}
	return
}

type TOT Section

// MakeTOT creates TOT that can hold n time offests.
func MakeTOT(n int) TOT {
	s := MakeEmptySection(3+5+2+2+n*13+4, false)
	s.SetTableId(0x73)
	s.SetPrivateSyntax(true)
	s.Alloc(5+2+2+n*13+4, 0)
	encodeU16(s[8:10], 0xf000|uint16(2+n*13))
	s[10] = byte(LocalTimeOffsetTag)
	s[11] = byte(n * 13)
	return TOT(s)
}

func (tot TOT) MakeCRC() {
	l := len(tot)
	crc := mpegCRC32(tot[0 : l-4])
	encodeU32(tot[l-4:l], crc)
}

// SetTime converts UTC time t to MJD and stores it in tot. TOT has CRC sum so
// modified TOT is invalid. Use MakeCRC method to recalculate CRC.
func (tot TOT) SetTime(t time.Time) {
	encodeMJDUTC(tot[3:8], t)
}

func (tot TOT) SetLTO(n int, lto LocalTimeOffset) {
	d := LocalTimeOffsetDescriptor(tot[12 : len(tot)-4])
	d.Set(n, lto)
}
