package ts

import (
	"time"
)

// AF represents content of adaptation field
type AF []byte

const BadPCR = PCR(1<<64 - 1)

// PCR returns BadPCR if can't decode PCR in adaptation field
func (a AF) PCR() PCR {
	if len(a) < 7 {
		return BadPCR
	}
	return decodePCR(a[1:7])
}

// OPCR returns BadPCR if can't decode OPCR in adaptation field
func (a AF) OPCR() PCR {
	if len(a) < 7+6 {
		return BadPCR
	}
	return decodePCR(a[7:13])
}

type PCR uint64

func decodePCR(a []byte) PCR {
	b := uint(a[0])<<24 | uint(a[1])<<16 | uint(a[2])<<8 | uint(a[3])
	base := uint64(b)<<1 | uint64(a[4])>>7
	ext := uint(a[4]&1)<<8 | uint(a[5])
	if ext >= 300 {
		return BadPCR
	}
	return PCR(base*300 + uint64(ext))
}

// Nanosec returns c * 1000 / 27
func (c PCR) Nanosec() time.Duration {
	return time.Duration(c * 1000 / 27)
}

// Flags returns adaptation field flags.
// If len(af) == 0 returns zero flags (all methods returns false).
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
