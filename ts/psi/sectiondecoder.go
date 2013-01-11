package psi

import (
	"github.com/ziutek/dvb/ts"
)

// SectionReader is interface for read one MPEG-TS section. len(s) should be
// equal to MaxSectionLen or MaxISOSectionLen (if you read standard PSI tables).
// You can use shorter s (but not shorter that 8 bytes) if you are sure that
// read section should fit in it. If ReadSection returned error of TemporaryError
// type you can try read next section.
type SectionReader interface {
	ReadSection(s Section) error
}

var (
	ErrSectionLength  = TemporaryError("incorrect value of section_length field")
	ErrSectionPointer = TemporaryError("incorrect pointer_field")
	ErrSectionSpace   = TemporaryError("no free space for section decoding")
	ErrSectionCRC     = TemporaryError("section has incorrect CRC")
	ErrSectionData    = TemporaryError("too few data to decode section")
)

// SectionDecoder can decode section from stream of packets
type SectionDecoder struct {
	r   ts.PktReader
	pkt ts.Pkt
	b   bool
}

// NewSectionDecoder create decoder to decode sections from stream of TS packets
// readed from r. You can use r == nil and set it later using SetPktReader method.
func NewSectionDecoder(r ts.PktReader) *SectionDecoder {
	d := new(SectionDecoder)
	d.r = r
	d.pkt = ts.NewPkt()
	return d
}

// SetPktReader sets ts.PktReader that will be usea as data source
func (d *SectionDecoder) SetPktReader(r ts.PktReader) {
	d.r = r
}

// Reset resets internal state of decoder (discards possible buffered data for next
// section decoding)
func (d *SectionDecoder) Reset() {
	d.b = false
}

// ReadSection decodes one section.
func (d *SectionDecoder) ReadSection(s Section) error {
	if len(s) < 8 {
		panic("section length should be >= 8")
	}
	n, limit := 0, -1
	for {
		if d.b {
			d.b = false
		} else {
			if err := d.r.ReadPkt(d.pkt); err != nil {
				return err
			}
		}

		p := d.pkt.Payload()
		if len(p) == 0 {
			continue
		}

		offset := 0
		if d.pkt.Flags().PayloadStart() {
			offset = int(p[0]) + 1
			if offset >= len(p) {
				return ErrSectionPointer
			}
		}

		if n == 0 {
			// Decoding isn't started yet
			if offset == 0 {
				// Section doesn't start in this packet
				continue
			}
			p = p[offset:]
		} else {
			if offset != 0 {
				p = p[1:offset]
			}
		}

		if limit == -1 {
			// Copy only up to the section_length byte
			k := copy(s[n:3], p)
			n += k
			if n < 3 {
				// Can't decode section_length - need more bytes
				continue
			}
			limit = s.Len()
			if limit == -1 {
				return ErrSectionLength
			}
			if limit > len(s) {
				return ErrSectionSpace
			}
			p = p[k:]
		}

		n += copy(s[n:limit], p)

		if n < limit {
			if offset == 0 {
				continue
			}
			// New section begins in this packet
			return ErrSectionData
		}

		// Whole section read
		if s.CheckCRC() {
			break
		}
		return ErrSectionCRC
	}
	return nil
}
