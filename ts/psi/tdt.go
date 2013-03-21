package psi

import (
	"github.com/ziutek/dvb"
	"time"
)

const TDTLen = 3 + 5

type TDT Section

func NewTDT() TDT {
	return TDT(make(Section, TDTLen))
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
	mjd := float64(int(t[3])<<8 | int(t[4]))
	year := int((mjd - 15078.2) / 365.25)
	month := int((mjd - 14956.1 - float64(int(float64(year)*365.25))) / 30.6001)
	day := int(mjd) - 14956 - int(float64(year)*365.25) - int(float64(month)*30.6001)
	if month == 14 || month == 15 {
		year++
		month -= 12
	}
	month--
	year += 1900
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
	if s.Len() != TDTLen || s.TableId() != 0x70 {
		return ErrTDTSectionSyntax
	}
	return nil
}
