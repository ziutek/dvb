package psi

import (
	"github.com/ziutek/dvb"
	"time"
)

const TDTSectionLen = 3 + 5

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

var ErrTDTSectionSyntax = dvb.TemporaryError("incorrect TDT section syntax")

func ParseTDT(s Section) (utc time.Time, err error) {
	if s.Len() != TDTSectionLen || s.TableId() != 0x70 {
		err = ErrTDTSectionSyntax
		return
	}
	hour := bcd2bin(s[5])
	if hour == -1 {
		return
	}
	min := bcd2bin(s[6])
	if min == -1 {
		return
	}
	sec := bcd2bin(s[7])
	if sec == -1 {
		return
	}
	mjd := float64(int(s[3])<<8 | int(s[4]))
	year := int((mjd - 15078.2) / 365.25)
	month := int((mjd - 14956.1 - float64(int(float64(year)*365.25))) / 30.6001)
	day := int(mjd) - 14956 - int(float64(year)*365.25) - int(float64(month)*30.6001)
	if month == 14 || month == 15 {
		year++
		month -= 12
	}
	month--
	year += 1900
	utc = time.Date(
		int(year), time.Month(month), int(day),
		hour, min, sec, 0,
		time.UTC,
	)
	return
}
