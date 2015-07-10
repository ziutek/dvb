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

type TDT Section

func MakeTDT() TDT {
	s := MakeEmptySection(TDTSectionLen, false)
	s.SetTableId(0x70)
	s.SetPrivateSyntax(true)
	s.Alloc(5)
	return TDT(s)
}

// SetTime converts UTC time t to MJD and stores it in tdt. TDT has no
// CRC sum so modified TDT is valid.
func (tdt TDT) SetTime(t time.Time) {
	encodeMJDUTC(tdt[3:8], t)
}

type TOT Section

func MakeTOT() TOT {
	// Alloc section that can storex one local_time_offset_descriptor.
	s := MakeEmptySection(3+5+2+15+4, false)
	s.SetTableId(0x73)
	s.SetPrivateSyntax(true)
	/// s.Alloc() ...
	return TOT(s)
}
