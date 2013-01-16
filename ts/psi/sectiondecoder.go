package psi

import (
	"github.com/ziutek/dvb/ts"
	"log"
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
	r        ts.PktReplacer
	pkt      *ts.ArrayPkt
	buffered bool // Not processed data in pkt
}

// NewSectionDecoder creates section decoder. You can use r == nil and
// set source of packets lather using SetPktReplacer or SetPktReader method.
func NewSectionDecoder(r ts.PktReplacer) *SectionDecoder {
	return &SectionDecoder{r: r, pkt: new(ts.ArrayPkt)}
}

// SetPktReplacer sets ts.PktReplacer that will be used as data source
func (d *SectionDecoder) SetPktReplacer(r ts.PktReplacer) {
	d.r = r
}

// SetPktReader sets ts.PktReader that will be usea as data source
func (d *SectionDecoder) SetPktReader(r ts.PktReader) {
	d.r = ts.PktReaderAsReplacer{r}
}

// Reset resets internal state of decoder (discards possible buffered data for next
// section decoding)
func (d *SectionDecoder) Reset() {
	d.buffered = false
}

// ReadSection decodes one section.
func (d *SectionDecoder) ReadSection(s Section) error {
	if len(s) < 8 {
		panic("section length should be >= 8")
	}
	n, limit := 0, -1
	for {
		if d.buffered {
			d.buffered = false
		} else {
			var err error
			if d.pkt, err = d.r.ReplacePkt(d.pkt); err != nil {
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
				d.buffered = true
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

		log.Println(s)
		log.Println(d.pkt.Pid(), "n:", n, "limit:", limit, "offset:", offset, "buffered:", d.buffered)
		if n < limit {
			if d.buffered {
				// New section begins in this packet
				return ErrSectionData
			}
			continue
		}

		// Whole section was read
		if s.CheckCRC() {
			break
		}
		return ErrSectionCRC
	}
	return nil
}
