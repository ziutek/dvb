package psi

import (
	"errors"
	"time"
)

func checkPid(pid int16) {
	if uint(pid) > 8191 {
		panic("Bad PID")
	}
}

func decodeU32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func decodeU24(b []byte) uint32 {
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

func decodeU16(b []byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

func encodeU32(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func encodeU16(b []byte, v uint16) {
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

func decodeBCD(bcd byte) int {
	h := int(bcd) >> 4
	if h > 9 {
		return -1
	}
	l := int(bcd) & 0x0f
	if l > 9 {
		return -1
	}
	return h*10 + l
}

func encodeBCD(d int) byte {
	if uint(d) > 99 {
		panic("psi: encodeBCD(>99)")
	}
	return byte(d/10<<4 | d%10)

}

var ErrBadMJDUTC = errors.New("bad MJD UTC time")

func decodeMJDUTC(b []byte) (utc time.Time, err error) {
	if len(b) != 5 {
		panic("dcodeMJDUTC with len(b) != 5")
	}
	hour := decodeBCD(b[2])
	if hour < 0 {
		err = ErrBadMJDUTC
		return
	}
	min := decodeBCD(b[3])
	if min < 0 {
		err = ErrBadMJDUTC
		return
	}
	sec := decodeBCD(b[4])
	if sec < 0 {
		err = ErrBadMJDUTC
		return
	}
	mjd := float64(int(b[0])<<8 | int(b[1]))
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

func encodeMJDUTC(b []byte, t time.Time) {
	if len(b) != 5 {
		panic("encodeMJDUTC with len(b) != 5")
	}
	t = t.UTC()
	var l int
	y, m, d := t.Date()
	if m <= 2 {
		l = 1
	}
	mjd := 14956 + int(d) +
		int(float64(y-1900-l)*365.25) +
		int(float64(int(m)+1+l*12)*30.6001)
	b[0] = byte(mjd >> 8)
	b[1] = byte(mjd)
	b[2] = encodeBCD(t.Hour())
	b[3] = encodeBCD(t.Minute())
	b[4] = encodeBCD(t.Second())
}

func setLoopLen(b []byte, n int) {
	if uint(n) > 0xfff {
		panic("psi: Bad descriptors loop length to set")
	}
	b[0] = 0xf0 | byte(n>>8)
	b[1] = byte(n)
}

func clearLoopLen(b []byte) {
	b[0] = 0xf0
	b[1] = 0
}

func loopLen(b []byte) int {
	return int(b[0]&0x0f)<<8 + int(b[1])
}

var crcTable [256]uint32

func mpegCRC32(buf []byte) uint32 {
	crc := uint32(0xffffffff)
	for _, b := range buf {
		crc = crcTable[byte(crc>>24)^b] ^ (crc << 8)
	}
	return crc
}

func init() {
	poly := uint32(0x04c11db7)
	for i := 0; i < 256; i++ {
		crc := uint32(i) << 24
		for j := 0; j < 8; j++ {
			if crc&0x80000000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		crcTable[i] = crc
	}
}
