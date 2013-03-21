package psi

import (
	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/ts"
	"log"
)

var (
	ErrSectionLength  = dvb.TemporaryError("incorrect value of section_length field")
	ErrSectionPointer = dvb.TemporaryError("incorrect pointer_field")
	ErrSectionSpace   = dvb.TemporaryError("no free space for section decoding")
	ErrSectionCRC     = dvb.TemporaryError("section has incorrect CRC")
	ErrSectionData    = dvb.TemporaryError("too few data to decode section")
)

// SectionDecoder can decode section from stream of packets
type SectionDecoder struct {
	r        ts.PktReplacer
	pkt      *ts.ArrayPkt
	buffered bool // Not processed data in pkt
	checkCRC bool
}

// NewSectionDecoder creates section decoder. You can use r == nil and
// set source of packets lather using SetPktReplacer or SetPktReader method.
func NewSectionDecoder(r ts.PktReplacer, checkCRC bool) *SectionDecoder {
	return &SectionDecoder{r: r, pkt: new(ts.ArrayPkt), checkCRC: checkCRC}
}

// SetPktReplacer sets ts.PktReplacer that will be used as data source
func (d *SectionDecoder) SetPktReplacer(r ts.PktReplacer) {
	d.r = r
}

// SetPktReader sets ts.PktReader that will be usea as data source
func (d *SectionDecoder) SetPktReader(r ts.PktReader) {
	d.r = ts.PktReaderAsReplacer{r}
}

// Reset resets internal state of decoder (discards possible buffered data for
// next section decoding)
func (d *SectionDecoder) Reset() {
	d.buffered = false
}

// ReadSection decodes one section.
func (d *SectionDecoder) ReadSection(s Section) error {
	if len(s) < 8 {
		panic("section length should be >= 8")
	}
	var (
		p   []byte
		err error
	)

	// Waiting for packet where a payload starts
	for {
		if d.buffered {
			d.buffered = false
		} else {
			if d.pkt, err = d.r.ReplacePkt(d.pkt); err != nil {
				return err
			}
		}
		p = d.pkt.Payload()
		if len(p) == 0 {
			continue
		}
		if d.pkt.Flags().PayloadStart() {
			offset := int(p[0]) + 1
			if offset >= len(p) {
				return ErrSectionPointer
			}
			p = p[offset:]
			break
		}
	}
	// p contains begining of section, n number of bytes copied to s
	n := copy(s[:3], p)
	p = p[n:]
	for n < 3 {
		// p doesn't contain enough data
		if d.pkt, err = d.r.ReplacePkt(d.pkt); err != nil {
			return err
		}
		d.buffered = true // d.pkt contains next packet
		p = d.pkt.Payload()
		k := copy(s[n:3], p)
		p = p[k:]
		n += k
	}
	l := s.Len()
	if l == -1 {
		return ErrSectionLength
	}
	if l > len(s) {
		log.Println("l > len(s)", l, len(s))
		return ErrSectionSpace
	}
	// Now we know section length, so we can copy data and read next packets
	for {
		k := copy(s[n:l], p)
		n += k
		if n == l {
			break
		}
		if d.pkt, err = d.r.ReplacePkt(d.pkt); err != nil {
			return err
		}
		d.buffered = true // d.pkt contains next packet
		p = d.pkt.Payload()
		if d.pkt.PayloadStart() {
			p = p[1:] // skip pointer_field
		}
	}
	// We read all needed data.
	if d.buffered && !d.pkt.Flags().PayloadStart() {
		d.buffered = false // d.pkt doesn't contain begining of next section.
	}
	if !d.checkCRC {
		return nil
	}
	if s.CheckCRC() {
		return nil
	}
	return ErrSectionCRC
}

/*func (d *SectionDecoder) ReadSection(s Section) error {
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
}*/
