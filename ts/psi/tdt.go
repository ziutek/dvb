package psi

import (
	"github.com/ziutek/dvb"
	"time"
)

const TDTSectionLen = 3 + 5

var ErrTDTSectionSyntax = dvb.TemporaryError("incorrect TDT section syntax")

func ParseTDT(s Section) (time.Time, error) {
	if s.TableId() != 0x70 || s.GenericSyntax() || !s.PrivateSyntax() || s.Len() != TDTSectionLen {
		return time.Time{}, ErrTDTSectionSyntax
	}
	return decodeMJDUTC(s[3:8])
}
