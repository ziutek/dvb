package psi

import (
	"github.com/ziutek/dvb"
	"time"
)

var ErrTOTSectionSyntax = dvb.TemporaryError("incorrect TOT section syntax")

func ParseTOT(s Section, checkCRC bool) (utc time.Time, err error) {
	if checkCRC && !s.CheckCRC() {
		err = ErrSectionCRC
		return
	}
	if s.TableId() != 0x73 || s.GenericSyntax() || s.Len() < 3+7+4 {
		err = ErrTOTSectionSyntax
		return
	}
	utc, err = decodeMJDUTC(s[3:8])
	if err != nil {
		return
	}
	l := int(decodeU16(s[8:10]) & 0x0fff)
	if s.Len()-(3+7+4) != l {
		err = ErrTOTSectionSyntax
		return
	}
	dl := DescriptorList(s[10 : 10+l])
	for len(dl) != 0 {
		var d Descriptor
		d, dl = dl.Pop()
		if d == nil {
			err = ErrSectionCRC
			return
		}
		// if d.Tag() == LocalTimeOffsetTag
	}
	return
}
