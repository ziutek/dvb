package psi

import (
	"github.com/ziutek/dvb/ts"
	"log"
)

// SectionEncoder can encode sections into stream of MPEG-TS packets. It has
// one TS packet internal buffer and writes it down only when it need more
// space or Flush method is called.
type SectionEncoder struct {
	r      ts.PktReplacer
	pid    uint16
	cc     byte
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

// NewSectionDecoder creates section decoder. You can use r == nil and
// set it lather using SetPktReplacer or SetPktWriter method.
func NewSectionEncoder(r ts.PktReplacer, pid uint16) *SectionEncoder {
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
		var err error
		// e.pkt contains some data from previous section
		if e.pkt.PayloadStart() {
			// Previous section starts and ends in this packet. We can't add any
			// data to it.
			err = e.Flush()
			e.pkt.SetPayloadStart(true) // section will start in new packet
			e.pkt.Payload()[0] = 0      // pointer_field
			e.offset++
		} else {
			// Previous section ends in this packet, but doesn't start in it.
			// We can add begineng of new section to it.
			e.pkt.SetPayloadStart(true)
			n := copy(e.pkt.Payload()[e.offset:], s)
			s = s[n:]
			err = e.Flush() // we allways need to flush such packet
		}
		if err != nil {
			return err
		}
	} else {
		// e.pkt is empty.
		e.pkt.SetPayloadStart(true) // section will start in e.pkt
		e.pkt.Payload()[0] = 0      // pointer_field
		e.offset++
	}
	// At this point we allways have an empty e.pkt with valid header and
	// properly set payload_start_indicator and pointer_field

	for len(s) > 0 {
		log.Print("len(s):", len(s))
		p := e.pkt.Payload()
		if len(s)+1 < len(p) && !e.pkt.PayloadStart() {
			// This is the last packet of current section and there will be a
			// place for at least one byte of next section in current packet.
			p[0] = len(s) // pointer_field
			p = p[1:]
		}
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
