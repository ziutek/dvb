package pes

import (
	"time"

	"github.com/ziutek/dvb/ts"
)

type Header []byte

const (
	programStreamMap       = 0xbc
	privateStream1         = 0xbd
	paddingStream          = 0xbe
	privateStream2         = 0xbf
	ecmStream              = 0xf0
	emmStream              = 0xf1
	dsmccStream            = 0xf2
	h2221TypeEStream       = 0xf8
	programStreamDirectory = 0xff
)

func (h Header) StreamId() byte {
	return h[3]
}

func (h Header) PktLen() int {
	return int(decodeU16(h[4:6]))
}

type HeaderFlags uint16

const (
	ScramblingControl HeaderFlags = 3 << 12 // PES_scrambling_control
	Priority          HeaderFlags = 1 << 11 // PES_priority
	DataAlignment     HeaderFlags = 1 << 10 // data_alignment_indicator
	Copyright         HeaderFlags = 1 << 9  // copyright
	Original          HeaderFlags = 1 << 8  // original_or_copy
	HasPTS            HeaderFlags = 1 << 7  // PTS_DTS_flags
	HasDTS            HeaderFlags = 1 << 6  // PTS_DTS_flags
	HasESCR           HeaderFlags = 1 << 5  // ESCR_flag
	HasESRate         HeaderFlags = 1 << 4  // ES_rate_flag
	HasDSMTrickMode   HeaderFlags = 1 << 3  // DSM_trick_mode_flag
	HasAdditCopyInfo  HeaderFlags = 1 << 2  // additional_copy_info_flag
	HasCRC            HeaderFlags = 1 << 1  // PES_CRC_flag
	HasExtension      HeaderFlags = 1 << 0  // PES_extension_flag
)

func (h Header) HasFlags() bool {
	sid := h.StreamId()
	return !(sid == programStreamMap ||
		sid == privateStream1 ||
		sid == paddingStream ||
		sid == privateStream2 ||
		sid == ecmStream ||
		sid == emmStream ||
		sid == dsmccStream ||
		sid == h2221TypeEStream ||
		sid == programStreamDirectory)
}

func (h Header) Flags() HeaderFlags {
	if !h.HasFlags() {
		return 0
	}
	return HeaderFlags(decodeU16(h[6:8]))
}

func (h Header) optLen() int {
	return int(h[8])
}

func (h Header) IsValid() bool {
	if len(h) < 6 || h[0] != 0 || h[1] != 0 || h[2] != 1 {
		return false
	}
	if !h.HasFlags() {
		return true
	}
	if len(h) < 9 {
		return false
	}
	flags := h.Flags()
	if flags>>14 != 2 {
		return false
	}
	var optLen int
	switch flags&(HasPTS | HasDTS) {
	case HasPTS:
		optLen = 5
	case HasPTS | HasDTS:
		optLen = 10
	default:
		return false
	}
	return len(h) >= h.optLen()+9 && h.optLen() >= optLen
	// BUG: Most optional fields are not implemented.
}

func (h Header) PTS() TimeStamp {
	if fl := h.Flags(); fl&HasPTS == 0 {
		return -1
	}
	t := decodeU40(h[9:14])
	if t&0x2100010001 != 0x2100010001 {
		return -1
	}
	return TimeStamp(t&0xfffe>>1 | t&0xfffe0000>>2 | t&0xe00000000>>3)
}

type TimeStamp int64

func (t TimeStamp) PCR() ts.PCR {
	return ts.PCR(t * 300)
}

func (t TimeStamp) Nanosec() time.Duration {
	return time.Duration(t*1e5+4) / 9
}
