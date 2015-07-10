package psi

import (
	"github.com/ziutek/dvb/ts"
)

// SectionEncoder can encode sections into stream of MPEG-TS packets. It has
// one TS packet internal buffer and writes it down only when it need more
// space or Flush method is called.
type SectionEncoder struct {
	r      ts.PktReplacer
	pid    int16
	cc     int8
	flags  ts.PktFlags
	pkt    *ts.ArrayPkt
	offset int // offset in pkt.Payload()
}

func (e *SectionEncoder) setupPktHeader() {
	e.pkt.SetSync()
	e.pkt.SetPid(e.pid)
	e.pkt.SetFlags(e.flags)
	e.pkt.SetCC(e.cc)
	e.cc++
}

// NewSectionEncoder creates section encoder. You can use r == nil and
// set it lather using SetPktReplacer or SetPktWriter method.
func NewSectionEncoder(r ts.PktReplacer, pid int16) *SectionEncoder {
	e := &SectionEncoder{
		r:   r,
		pid: pid,
		pkt: new(ts.ArrayPkt),
	}
	e.flags.SetContainsPayload(true)
	e.setupPktHeader() // e.pkt should allways contain a valid header
	return e
}

func (e *SectionEncoder) write() error {
	e.offset = 0
	var err error
	e.pkt, err = e.r.ReplacePkt(e.pkt)
	e.setupPktHeader() // e.pkt should allways contain a valid header
	return err
}

// Flush writes internally buffered packet adding stuffing bytes at end if need.
func (e *SectionEncoder) Flush() error {
	if e.offset == 0 {
		return nil
	}
	p := e.pkt.Payload()[e.offset:]
	// Add padding
	for i := range p {
		p[i] = 0xff
	}
	return e.write()
}

// WriteSection encodes one section into one or more MPEG-TS packets.
func (e *SectionEncoder) WriteSection(s Section) error {
	s = s[:s.Len()]
	if len(s) == 0 {
		return nil
	}
	if e.offset > 0 {
		// e.pkt contains some data from previous section
		p := e.pkt.Payload()
		if e.pkt.PayloadStart() || e.offset+2 >= len(p) {
			// Previous section starts and ends in this packet or there is no
			// place for even one byte in it.
			if err := e.Flush(); err != nil {
				return err
			}
			e.pkt.SetPayloadStart(true) // section will start in new packet
			e.pkt.Payload()[0] = 0      // set pointer_field in new packet
			e.offset++
		} else {
			// There is place for at least one byte of new section in this
			// packet.
			copy(p[1:], p[:e.offset])   // move data to add pointer field
			p[0] = byte(e.offset)       // set pointer_field
			e.pkt.SetPayloadStart(true) // section starts in current packet
			e.offset++
			n := copy(p[e.offset:], s)
			s = s[n:]
			e.offset += n
			// Write this packet (with pading if there is no enough data in
			// section to fill it fully)
			if err := e.Flush(); err != nil {
				return err
			}
		}
	} else {
		// e.pkt is empty.
		e.pkt.SetPayloadStart(true) // section will start in e.pkt
		e.pkt.Payload()[0] = 0      // pointer_field
		e.offset++
	}
	// At this point we allways have an empty e.pkt with valid header,
	// properly set payload_start_indicator and pointer_field

	for len(s) > 0 {
		p := e.pkt.Payload()[e.offset:]
		n := copy(p, s)
		s = s[n:]
		e.offset += n
		if n == len(p) {
			// Packet is full
			if err := e.write(); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetPktReplacer sets ts.PktReplacer that will be used to write packets.
func (e *SectionEncoder) SetPktReplacer(r ts.PktReplacer) {
	e.r = r
}

// SetPktWriter sets ts.PktWriter that will be used to write packets.
func (e *SectionEncoder) SetPktWriter(r ts.PktWriter) {
	e.r = ts.PktWriterAsReplacer{r}
}
