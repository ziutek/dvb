package psi

import (
	"errors"
	"github.com/ziutek/dvb/ts"
)

const SectionMaxLen = 4096

type Section []byte

func NewSection(l int) Section {
	return make(Section, l)
}

// TableId returns the value of table_id field
func (s Section) TableId() byte {
	return s[0]
}

// SetTableId sets the value of table_id field
func (s Section) SetTableId(id byte) {
	s[0] = id
}

// SyntaxIndicator returns the value of section_syntax_indicator field
func (s Section) SyntaxIndicator() bool {
	return s[1]&0x80 != 0
}

// SetSyntaxIndicator sets the value of section_syntax_indicator field
func (s Section) SetSyntaxIndicator(si bool) {
	if si {
		s[1] |= 0x80
	} else {
		s[1] &^= 0x80
	}
}

// PrivateIndicator returns the value of private_syntax_indicator field
func (s Section) PrivateIndicator() bool {
	return s[1]&0x40 != 0
}

// SetPrivateIndicator sets the value of private_syntax_indicator field
func (s Section) SetPrivateIndicator(pi bool) {
	if pi {
		s[1] |= 0x40
	} else {
		s[1] &^= 0x40
	}
}

// Len returns length of the whole section (section_length + 3) or -1 if
// section_length filed contains incorrect value
func (s Section) Len() int {
	l := ((int(s[1]&0x0f) << 8) | int(s[2])) + 3
	if l < 4+3 || l > SectionMaxLen {
		return -1
	}
	return l
}

// SetLen sets the value of section_length field to l-3.
// It panics if l < 7 or l > SectionMaxLen
func (s Section) SetLenField(l int) {
	if l < 4+3 || l > SectionMaxLen {
		panic("incorrect value for section length field")
	}
	l -= 3
	h := byte(l>>8) & 0x0f
	s[1] = s[1]&0xf0 | h
	s[2] = byte(l)
}

// Version returns the value of version_number field
func (s Section) Version() byte {
	return (s[5] >> 1) & 0x1f
}

// SetVersion sets the value of version_number field.
// It panic if v > 31
func (s Section) SetVerison(v byte) {
	if v > 31 {
		panic("value for version_number field is too large")
	}
	s[5] = s[5]&0x3e | (v << 1)
}

// Current returns the value of current_next_indicator field
func (s Section) Current() bool {
	return s[5]&0x01 != 0
}

// SetCurrent sets the value of current_next_indicator field
func (s Section) SetCurrent(c bool) {
	if c {
		s[5] |= 0x01
	} else {
		s[5] &^= 0x01
	}
}

// Number returns the vale of section_number field
func (s Section) Number() byte {
	return s[6]
}

// SetNumber sets the vale of section_number field
func (s Section) SetNumber(n byte) {
	s[6] = n
}

// LastNumber returns the vale of last_section_number field
func (s Section) LastNumber() byte {
	return s[7]
}

// SetLastNumber sets the vale of last_section_number field
func (s Section) SetLastNumber(n byte) {
	s[7] = n
}

// CheckCRC returns true if s.Length() is valid and IEEE CRC32 of whole
// section is correct
func (s Section) CheckCRC() bool {
	l := s.Len()
	if l == -1 || len(s) < l {
		return false
	}
	crc := decodeU32(s[l-4 : l])
	return mpegCRC32(s[0:l-4]) == crc
}

var ErrSectionBadLength = errors.New("incorrect value of section_length field")

// MakeCRC calculates CRC32 for whole section and uses it to set CRC_32 field
func (s Section) MakeCRC() {
	l := s.Len()
	if l == -1 || len(s) < l {
		panic(ErrSectionBadLength)
	}
	crc := mpegCRC32(s[0 : l-4])
	encodeU32(s[l-4:l], crc)
}

// SectionDecoder can decode section from stream of packets
type SectionDecoder struct {
	s     Section
	n     int
	limit int
}

var (
	ErrDecodePointerField = errors.New("incorrect pointer_field")
	ErrDecodeNoSpace      = errors.New("no free space for section decoding")
	ErrDecodeCRC          = errors.New("section has incorrect CRC")
	ErrDecodeTooFewData   = errors.New("too few data to decode section")
)

// Init initializes decoder to decode into s. It panics if len(s) < 8
func (d *SectionDecoder) Init(s Section) {
	if len(s) < 8 {
		panic("section length need to be >= 8")
	}
	d.s = s
	d.Reset()
}

func (d *SectionDecoder) Reset() {
	d.n = 0
	d.limit = -1
}

func (d *SectionDecoder) Section() Section {
	return d.s
}

// Decode decodes first section from packets passed to it. All packets
// should contain the same PID. If Decode returns true or error, the section is
// decoded or can't be decoded but pkt can contain beginning of the next
// section so it should be used as start packet for next section decoding.
func (d *SectionDecoder) Decode(pkt ts.Pkt) (ok bool, err error) {
	defer func() {
		if ok == true || err != nil {
			d.Reset()
		}
	}()

	p := pkt.Payload()
	if len(p) == 0 {
		return
	}

	offset := 0
	if pkt.Flags().PayloadStart() {
		offset = int(p[0]) + 1
		if offset >= len(p) {
			err = ErrDecodePointerField
			return
		}
	}

	if d.n == 0 {
		// Decoding isn't started yet
		if offset == 0 {
			// Section doesn't start in this packet
			return
		}
		p = p[offset:]
	} else {
		if offset != 0 {
			p = p[1:offset]
		}
	}

	if d.limit == -1 {
		// Copy only up to section_length byte
		n := copy(d.s[d.n:3], p)
		d.n += n
		if d.n < 3 {
			// Can't decode section_length - need more bytes
			return
		}
		d.limit = d.s.Len()
		if d.limit == -1 {
			err = ErrSectionBadLength
			return
		}
		if d.limit > len(d.s) {
			err = ErrDecodeNoSpace
			return
		}
		p = p[n:]
	}

	d.n += copy(d.s[d.n:d.limit], p)

	if d.n < d.limit {
		if offset == 0 {
			return
		}
		// New section begins in this packet
		err = ErrDecodeTooFewData
		return
	}

	// Whole section read
	if d.s.CheckCRC() {
		ok = true
		return
	}
	err = ErrDecodeCRC
	return
}
