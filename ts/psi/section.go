package psi

import (
	"fmt"
)

const (
	SectionMaxLen    = 4096
	ISOSectionMaxLen = 1024
)

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

// TableIdExt returns the value of table_id_extension
func (s Section) TableIdExt() uint16 {
	return decodeU16(s[3:5])
}

// Set TableIdExt sets the value of table_id_extension
func (s Section) SetTableIdExt(id uint16) {
	encodeU16(s[3:5], id)
}

// SyntaxIndicator returns the value of section_syntax_indicator field
func (s Section) GenericSyntax() bool {
	return s[1]&0x80 != 0
}

// SetGenericSyntax sets the value of section_syntax_indicator field
func (s Section) SetGenericSyntax(si bool) {
	if si {
		s[1] |= 0x80
	} else {
		s[1] &^= 0x80
	}
}

// PrivateIndicator returns the value of private_syntax_indicator field
func (s Section) PrivateSyntax() bool {
	return s[1]&0x40 != 0
}

// SetPrivateSyntax sets the value of private_syntax_indicator field
func (s Section) SetPrivateSyntax(pi bool) {
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
func (s Section) SetLen(l int) {
	if l < 4+3 || l > SectionMaxLen {
		panic("incorrect value for section_length field")
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

// Number returns the value of section_number field
func (s Section) Number() byte {
	return s[6]
}

// SetNumber sets the value of section_number field
func (s Section) SetNumber(n byte) {
	s[6] = n
}

// LastNumber returns the value of last_section_number field
func (s Section) LastNumber() byte {
	return s[7]
}

// SetLastNumber sets the value of last_section_number field
func (s Section) SetLastNumber(n byte) {
	s[7] = n
}

// Data rturns data part of section
func (s Section) Data() []byte {
	end := s.Len() - 4
	if end == -1-4 || end > len(s) {
		panic("there is no enough data or section_length has incorrect value")
	}
	return s[8:end]
}

// CheckCRC returns true if s.Length() is valid and CRC32 of whole
// section is correct
func (s Section) CheckCRC() bool {
	l := s.Len()
	if l == -1 || len(s) < l {
		return false
	}
	crc := decodeU32(s[l-4 : l])
	return mpegCRC32(s[0:l-4]) == crc
}

// MakeCRC calculates CRC32 for whole section and uses it to set CRC_32 field
func (s Section) MakeCRC() {
	l := s.Len()
	if l == -1 || len(s) < l {
		panic("bad section length to calculate CRC sum")
	}
	crc := mpegCRC32(s[0 : l-4])
	encodeU32(s[l-4:l], crc)
}

func (s Section) String() string {
	return fmt.Sprintf(
		"TableId: %d Syntax: generic=%t private=%t Len: %d Version: %d "+
			"Current: %t Number: %d/%d",
		s.TableId(), s.GenericSyntax(), s.PrivateSyntax(), s.Len(), s.Version(),
		s.Current(), s.Number(), s.LastNumber(),
	)
}

// SectionReader is an interface that wraps the ReadSection method.
type SectionReader interface {
	// ReadSection reads one section into s. len(s) should be equal to
	// MaxSectionLen or MaxISOSectionLen (if you read standard PSI tables). You
	// can use shorter s (but not shorter that 8 bytes) if you are sure that
	// read section should fit in it. If ReadSection returned error of
	// dvb.TemporaryError type you can try read next section.
	ReadSection(s Section) error
}
