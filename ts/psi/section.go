package psi

import (
	"fmt"
)

const (
	SectionMaxLen    = 4096
	ISOSectionMaxLen = 1024
)

type Section []byte

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
	if !s.GenericSyntax() {
		panic("GenericSyntax need for TableIdExt")
	}
	return decodeU16(s[3:5])
}

// Set TableIdExt sets the value of table_id_extension
func (s Section) SetTableIdExt(id uint16) {
	encodeU16(s[3:5], id)
}

// GenericSyntax returns the value of section_syntax_indicator field
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

func (s Section) Reserved() int {
	return (int(s[1]) & 0x30) >> 4
}

/*
// SetReserved: r should be 3.
func (s Section) SetReserved(r int) {
	s[1] = s[1]&^0x30 | byte(r<<4)&0x30
}
*/

// Len returns length of the whole section (section_length + 3) or -1 if
// section_length filed contains incorrect value
func (s Section) Len() int {
	l := int(decodeU16(s[1:3]))&0xfff + 3
	if l > SectionMaxLen {
		return -1
	}
	if s.GenericSyntax() && l < 3+5+4 {
		return -1
	}
	return l
}

func (s Section) Cap() int {
	return len(s)
}

func (s Section) setLen(n int) {
	if n > len(s) {
		panic("psi: too big value for section length")
	}
	n -= 3
	s[1] = s[1]&0xf0 | byte(n>>8)&0x0f
	s[2] = byte(n)
}

// Version returns the value of version_number field.
func (s Section) Version() int8 {
	if !s.GenericSyntax() {
		panic("psi: GenericSyntax need for Version")
	}
	return int8(s[5]>>1) & 0x1f
}

// SetVersion sets the value of version_number field.
// It panic if v > 31
func (s Section) SetVersion(v int8) {
	if uint(v) > 31 {
		panic("psi: value for version_number field is too large")
	}
	s[5] = s[5]&0xc1 | byte(v<<1)
}

// Current returns the value of current_next_indicator field
func (s Section) Current() bool {
	if !s.GenericSyntax() {
		panic("psi: GenericSyntax need for Current")
	}
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
	if !s.GenericSyntax() {
		panic("psi: GenericSyntax need for Number")
	}
	return s[6]
}

// SetNumber sets the value of section_number field
func (s Section) SetNumber(n byte) {
	s[6] = n
}

// LastNumber returns the value of last_section_number field
func (s Section) LastNumber() byte {
	if !s.GenericSyntax() {
		panic("psi: GenericSyntax need for LastNumber")
	}
	return s[7]
}

// SetLastNumber sets the value of last_section_number field
func (s Section) SetLastNumber(n byte) {
	s[7] = n
}

// Data rturns data part of section.
func (s Section) Data() []byte {
	end := s.Len()
	beg := 3
	if s.GenericSyntax() {
		beg += 5 // Generic header.
		end -= 4 // CRC.
	}
	if end < beg || end > len(s) {
		panic("psi: not enough data or section_length has incorrect value")
	}
	return s[beg:end]
}

// CheckCRC returns true if s.Len() is valid and CRC32 of whole
// section is correct
func (s Section) CheckCRC() bool {
	if !s.GenericSyntax() {
		panic("psi: GenericSyntax need for CheckCRC")
	}
	l := s.Len()
	if l < 4 || len(s) < l {
		return false
	}
	crc := decodeU32(s[l-4 : l])
	return mpegCRC32(s[0:l-4]) == crc
}

// MakeCRC calculates CRC32 for whole section and uses it to set CRC_32 field
func (s Section) MakeCRC() {
	if !s.GenericSyntax() {
		panic("psi: GenericSyntax need for MakeCRC")
	}
	l := s.Len()
	if l < 4 || len(s) < l {
		panic("psi: bad section length to calculate CRC sum")
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

// SectionWriter is an interface wraps the WriteSection method.
type SectionWriter interface {
	WriteSection(s Section) error
}

func MakeEmptySection(maxLen int, genericSyntax bool) Section {
	s := make(Section, maxLen)
	s[1] = 0x30 // Reseved bits.
	s.SetGenericSyntax(genericSyntax)
	s.SetEmpty()
	return s
}

// Alloc allocates n bytes in section. It returns nil if there is no room for
// n+m bytes. Alloc invalidates CRC sum if used. Use MakeCRC to
// recalculate it.
func (s Section) Alloc(n, m int) []byte {
	b := s.Len()
	e := b + n
	if e+m > s.Cap() {
		return nil
	}
	s.setLen(e)
	if s.GenericSyntax() {
		// s.Len() includes CRC
		b -= 4
		e -= 4
	}
	return s[b:e]
}

// SetEmpty initializes section lenght, so s becomes empty. SetEmpty
// invalidates CRC sum. Use MakeCRC to recalculate it.
func (s Section) SetEmpty() {
	n := 3
	if s.GenericSyntax() {
		// table_id_extension + reserved + version_number +
		// current_next_indicator + section_number + last_section_number = 5 B
		// CRC = 4 B
		s[5] = 0xc0 // Reserved bits.
		n += 5 + 4
	}
	s.setLen(n)
}
