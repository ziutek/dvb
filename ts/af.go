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

// Discontinuity returns true if discontinuity_indicator == 1
func (f AFFlags) Discontinuity() bool {
	return f&0x80 != 0
}

// RandomAccess returns true if random_access_indicator == 1
func (f AFFlags) RandomAccess() bool {
	return f&0x40 != 0
}

// ESPrio returns true if elementary_stream_priority_indicator == 1
func (f AFFlags) ESPrio() bool {
	return f&0x20 != 0
}

// ContainsPCR returns true if PCR_flag == 1
func (f AFFlags) ContainsPCR() bool {
	return f&0x10 != 0
}

// ContainsOPCR returns true if OPCR_flag == 1
func (f AFFlags) ContainsOPCR() bool {
	return f&8 != 0
}

// SplicingPoint returns true if splicing_point_flag == 1
func (f AFFlags) SplicingPoint() bool {
	return f&4 != 0
}

// PrivateData returns true if transport_private_data_flag == 1
func (f AFFlags) ContainsPrivateData() bool {
	return f&2 != 0
}

// HasExtension returns true if adaptation_field_extension_flag == 1
func (f AFFlags) HasExtension() bool {
	return f&1 != 0
}

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
	if !a.Flags().ContainsPCR() {
		return -1, ErrNotInAF
	}
	end := 1 + 6
	if len(a) < end {
		return -1, ErrAFTooShort
	}
	return decodePCR(a[end-6 : end])
}

func (a AF) SetPCR(pcr PCR) error {
	if !a.Flags().ContainsPCR() {
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
	if !f.ContainsOPCR() {
		return -1, ErrNotInAF
	}
	end := 1 + 7
	if f.ContainsPCR() {
		end += 6
	}
	if len(a) < end {
		return 0, ErrAFTooShort
	}
	return decodePCR(a[end-6 : end])
}

func (a AF) SpliceCountdown() (int8, error) {
	f := a.Flags()
	if !f.SplicingPoint() {
		return -1, ErrNotInAF
	}
	offset := 1
	if f.ContainsPCR() {
		offset += 6
	}
	if f.ContainsOPCR() {
		offset += 6
	}
	if len(a) < offset+1 {
		return -1, ErrAFTooShort
	}
	return int8(a[offset]), nil
}
