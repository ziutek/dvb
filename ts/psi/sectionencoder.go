package psi

import (
	"github.com/ziutek/dvb/ts"
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
	offset int
}

// NewSectionDecoder creates section decoder. You can use r == nil and
// set it lather using SetPktReplacer or SetPktWriter method.
func NewSectionEncoder(r ts.PktReplacer, pid uint16) *SectionEncoder {
	s := &SectionEncoder{
		r:   r,
		pid: pid,
		pkt: new(ts.ArrayPkt),
	}
	s.flags.SetContainsPayload(true)
	return s
}

// SetPktReplacer sets ts.PktReplacer that will be used to write packets.
func (e *SectionEncoder) SetPktReplacer(r ts.PktReplacer) {
	e.r = r
}

// SetPktWriter sets ts.PktWriter that will be used to write packets.
func (e *SectionEncoder) SetPktWriter(r ts.PktWriter) {
	e.r = ts.PktWriterAsReplacer{r}
}

// Flush writes internally buffered packet adding stuffing bytes at end if need.
func (e *SectionEncoder) Flush() error {
	if e.offset == 0 {
		return nil
	}
	// Setup header (packet obtained from PktReplacer can contain random data)
	e.pkt.SetSync()
	e.pkt.SetPid(e.pid)
	e.pkt.SetFlags(e.flags)
	e.pkt.SetCC(e.cc)
	e.cc++
	// Add padding
	for i := e.offset; i < ts.PktLen; i++ {
		e.pkt[i] = 0xff
	}
	e.offset = 0
	var err error
	e.pkt, err = e.r.ReplacePkt(e.pkt)
	return err
}

// WriteSection encodes one section into one or more MPEG-TS packets.
func (e *SectionEncoder) WriteSection(s Section) error {

	return nil
}
