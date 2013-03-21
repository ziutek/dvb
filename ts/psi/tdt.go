package psi

import (
	"github.com/ziutek/dvb"
	"time"
)

type TDT Section

func NewTDT() TDT {
	return TDT(make(Section, 12))
}

// Returns -1 if error
func bcd2bin(bcd byte) int {
	h := int(bcd >> 4)
	if h > 9 {
		return -1
	}
	l := int(bcd & 0x0f)
	if l > 9 {
		return -1
	}
	return h*10 + l
}

// UTC returns zero time if error
func (t TDT) UTC() (utc time.Time) {
	h := bcd2bin(t[5])
	if h == -1 {
		return
	}
	m := bcd2bin(t[6])
	if m == -1 {
		return
	}
	s := bcd2bin(t[7])
	if s == -1 {
		return
	}
	mjd := (uint(t[3])<<8 | uint(t[4])) * 10000
	year := (mjd - 150782000) / 3652500
	month := (mjd - 149561000 - year*3652500) / 306001
	day := (mjd - 149560000 - year*3652500 - month*306001) / 10000
	if month == 14 || month == 15 {
		year++
		month -= 12
	}
	month--
	return time.Date(int(year), time.Month(month), int(day), h, m, s, 0, time.UTC)
}

var (
	ErrTDTSectionSyntax = dvb.TemporaryError("incorrect TDT section syntax")
)

func (t TDT) Update(r SectionReader) error {
	s := Section(t)
	err := r.ReadSection(s)
	if err != nil {
		return err
	}
	if s.Len() != 12 || s.TableId() != 0x70 || !s.GenericSyntax() ||
		s.Number() != 0 || s.LastNumber() != 0 {
		return ErrTDTSectionSyntax
	}
	return nil
}
