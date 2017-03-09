package ts

import (
	"errors"
	"time"
)

// AF represents content of adaptation field
type AF []byte

// Flags returns adaptation field flags.
// If len(af) == 0 returns zero flags (all AFFlags methods returns false).
func (a AF) Flags() AFFlags {
	if len(a) == 0 {
		return 0
	}
	return AFFlags(a[0])
}

type AFFlags byte

const (
	Discontinuity       AFFlags = 0x80 // discontinuity_indicator
	RandomAccess        AFFlags = 0x40 // random_access_indicator
	ESPrio              AFFlags = 0x20 // elementary_stream_priority_indicator
	ContainsPCR         AFFlags = 0x10 // PCR_flag
	ContainsOPCR        AFFlags = 0x08 // OPCR_flag
	SplicingPoint       AFFlags = 0x04 // splicing_point_flag
	ContainsPrivateData AFFlags = 0x02 // transport_private_data_flag
	HasExtension        AFFlags = 0x01 // adaptation_field_extension_flag
)

var (
	ErrAFTooShort = errors.New("adaptation field is too short")
	ErrBadPCR     = errors.New("PCR decoding error")
	ErrNotInAF    = errors.New("no such entry in adaptation field")
)

// PCR contains the number of ticks of 27MHz generator. PCR generator counts
// modulo PCRModulo.
type PCR int64

const (
	PCRModulo = (1 << 33) * 300
	PCRFreq   = 27e6 // Hz
)

func decodePCR(a []byte) (PCR, error) {
	b := uint(a[0])<<24 | uint(a[1])<<16 | uint(a[2])<<8 | uint(a[3])
	base := uint64(b)<<1 | uint64(a[4])>>7
	ext := uint(a[4]&1)<<8 | uint(a[5])
	if ext >= 300 {
		return -1, ErrBadPCR
	}
	return PCR(base*300 + uint64(ext)), nil
}

func encodePCR(a []byte, pcr PCR) {
	if uint64(pcr) >= PCRModulo {
		panic("bad PCR value")
	}
	base := pcr / 300
	ext := pcr - base*300
	a[0] = byte(base >> 25)
	a[1] = byte(base >> 17)
	a[2] = byte(base >> 9)
	a[3] = byte(base >> 1)
	a[4] = byte(base<<7) | a[4]&0x7e | byte(ext>>8)
	a[5] = byte(ext)
}

// Nanosec returns (c * 1000 + 13) / 27
func (c PCR) Nanosec() time.Duration {
	return time.Duration(c*1000+13) / 27
}

func (c PCR) Add(ns time.Duration) PCR {
	c += PCR(ns*27+500) / 1000
	for c < 0 {
		c += PCRModulo
	}
	for c > PCRModulo {
		c -= PCRModulo
	}
	return c
}

// PCR returns value of PCR in a. It returns PCR == -1 and not nil
// error if there is no PCR in AF or it can't decode PCR.
func (a AF) PCR() (PCR, error) {
	if a.Flags()&ContainsPCR == 0 {
		return -1, ErrNotInAF
	}
	end := 1 + 6
	if len(a) < end {
		return -1, ErrAFTooShort
	}
	return decodePCR(a[end-6 : end])
}

func (a AF) SetPCR(pcr PCR) error {
	if a.Flags()&ContainsPCR == 0 {
		return ErrNotInAF
	}
	end := 1 + 6
	if len(a) < end {
		return ErrAFTooShort
	}
	encodePCR(a[end-6:end], pcr)
	return nil
}

// OPCR returns value of OPCR in a. It returns OPCR == -1 and not nil
// error if there is no OPCR in AF or it can't decode OPCR.
func (a AF) OPCR() (PCR, error) {
	f := a.Flags()
	if f&ContainsOPCR == 0 {
		return -1, ErrNotInAF
	}
	end := 1 + 7
	if a.Flags()&ContainsPCR != 0 {
		end += 6
	}
	if len(a) < end {
		return 0, ErrAFTooShort
	}
	return decodePCR(a[end-6 : end])
}

func (a AF) SpliceCountdown() (int8, error) {
	f := a.Flags()
	if f&SplicingPoint == 0 {
		return -1, ErrNotInAF
	}
	offset := 1
	if f&ContainsPCR != 0 {
		offset += 6
	}
	if f&ContainsOPCR != 0 {
		offset += 6
	}
	if len(a) < offset+1 {
		return -1, ErrAFTooShort
	}
	return int8(a[offset]), nil
}
