package psi

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

// SyntaxIndicator returns the value of section_syntax_indicator field
func (s Section) GenericSyntax() bool {
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

// Data rturns data part of section
func (s Section) Data() []byte {
	end := s.Len() - 4
	if end == -1-4 || end > len(s) {
		panic("there is no enough data or section_length has incorrect value")
	}
	return s[8:end]
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

// MakeCRC calculates CRC32 for whole section and uses it to set CRC_32 field
func (s Section) MakeCRC() {
	l := s.Len()
	if l == -1 || len(s) < l {
		panic("bad section length to calculate CRC sum")
	}
	crc := mpegCRC32(s[0 : l-4])
	encodeU32(s[l-4:l], crc)
}
