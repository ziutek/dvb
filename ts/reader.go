package ts

import (
	"github.com/ziutek/dvb"
	"io"
	"os"
	"syscall"
)

var (
	// ErrSync means a lost of MPEG-TS synchronization.
	ErrSync = dvb.TemporaryError("MPEG-TS synchronization error")
)

// PktReader is an interface that wraps the ReadPkt method.
type PktReader interface {
	// ReadPkt reads one MPEG-TS packet.
	// If it returns ErrSync or dvb.ErrOverflow you can try to Read next
	// packet.
	ReadPkt(Pkt) error
}

// PktReplacer is an interface that wraps the ReplacePkt method. After
// ReplacePkt call caller should not reffer to p content any more.
// If ReplacePkt returns an error it is guaranteed that r == p (but content
// reffered by p can be modified). Generally ReplacePkt should be used in this
// way:
//
//    p, err = q.ReplacePkt(p)
//    if err != nil {
//        ...
//    }
type PktReplacer interface {
	// ReplacePkt consumes packet reffered by p an returns other packet reffered
	// by r.
	// If it returns ErrSync or dvb.ErrOverflow you can try to call ReplacePkt
	// one more time.
	ReplacePkt(p *ArrayPkt) (r *ArrayPkt, e error)
}

// PktReaderAsReplacer converts any PktReader to PktReplacer
type PktReaderAsReplacer struct {
	PktReader
}

func (r PktReaderAsReplacer) ReplacePkt(p *ArrayPkt) (*ArrayPkt, error) {
	err := r.PktReader.ReadPkt(p)
	return p, err
}

// PktStream wraps any io.Reader interface and returns PktReader and
// PktReplacer implementation for read MPEG-TS packets from stream of bytes.
// Internally it doesn't allocate any memory so is friendly for real-time
// applications (it doesn't cause GC to run).
//
// Using PktStream you can start read at any point in stream. If the start point
// doesn't match a beginning of a packet, PktReader returns ErrSync and
// tries to synchronize during next read.
type PktStream struct {
	r       io.Reader
	syncBuf [3 * PktLen]byte
	sbStart int
}

// SetReader sets new io.Reader as stream source. Forces resynchronization.
func (p *PktStream) SetReader(r io.Reader) {
	p.r = r
	p.sbStart = -1
}

// NewPktStreame is equivalent to p := new(PktStream); p.SetReader(r)
func NewPktStream(r io.Reader) *PktStream {
	s := new(PktStream)
	s.SetReader(r)
	return s
}

func (s *PktStream) synchronize() (err error) {
	b := s.syncBuf[:]
	if s.sbStart == -1 {
		// First try of synchronization - read full buffer (three packets)
		_, err = io.ReadFull(s.r, b)
		s.sbStart = -2
	} else {
		// Subsequent try of synchronization - read next packet
		copy(b, b[PktLen:])
		_, err = io.ReadFull(s.r, b[2*PktLen:])
	}
	if err != nil {
		return
	}
	// Try to find a sync point in syncBuffer
	for i := 0; i < PktLen; i++ {
		if b[i] == 0x47 && b[i+PktLen] == 0x47 && b[i+2*PktLen] == 0x47 {
			// Sync point found
			s.sbStart = 0
			copy(b, b[i:])
			_, err = io.ReadFull(s.r, b[len(b)-i:])
			return
		}
	}
	return ErrSync
}

func convertEoverflow(err error) error {
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.EOVERFLOW {
		return dvb.ErrOverflow
	}
	return err
}

// ReadPkt reads one MPEG-TS packet directly to provided buffer with exception
// for out of sync state when it reads more than one packet to internal buffer
// and tries to synchronize. ReadPkt check len(pkt) and panics if it isn't
// PktLen (usefull if bound checking is disabled at compile time).
// ReadPkt converts os.PathError{Err: syscall.EOVERFLOW} to dvb.ErrOverflow.
func (s *PktStream) ReadPkt(pkt Pkt) error {
	if s.sbStart < 0 {
		if err := s.synchronize(); err != nil {
			return convertEoverflow(err)
		}
	}
	if s.sbStart != len(s.syncBuf) {
		// Copy packet from sync buffer
		copy(pkt.Bytes(), s.syncBuf[s.sbStart:])
		s.sbStart += PktLen
		return nil
	}
	// Read packet from io.Reader
	if _, err := io.ReadFull(s.r, pkt.Bytes()); err != nil {
		return convertEoverflow(err)
	}
	if !pkt.SyncOK() {
		s.sbStart = -1
		copy(s.syncBuf[:], pkt.Bytes())
		return ErrSync
	}
	return nil
}

// ReplacePkt works like ReadPkt but implements PktReplacer interface.
func (s *PktStream) ReplacePkt(p *ArrayPkt) (*ArrayPkt, error) {
	err := s.ReadPkt(p)
	return p, err
}

// Reader wraps PktReplacer or PktReader to implement a standard io.Reader
// interface. Internally it doesn't allocate any memory so it is friendly
// for real-time applications (it doesn't cause GC to run).
type Reader struct {
	p   PktReader
	pkt [PktLen]byte
	i   int
}

// SetPktReader sets new PktReader as packets source. If internal buffer
// contains some data from previous source they will be returned before any new
// read from new source.
func (r *Reader) SetPktReader(p PktReader) {
	r.p = p
	r.i = PktLen
}

// Read allow to read from MPEG-TS packet stream as from ordinary byte stream.
// It reads no more than 2*PktLen-1 bytes. If len(buf) >= 2*PktLen-1, buf always
// contains one packet at the end of data, plus (possibly) some data from
// previously not fully read packet. If you always use len(buf) >= PktLen it
// always read one MPEG-TS packet without internal buffering. If you need to
// fill big buffer for multiple packets, use io.ReadFull helper function.
func (r *Reader) Read(buf []byte) (n int, err error) {
	// Try to copy remaining data from internal buffer
	if r.i != PktLen {
		n = copy(buf, r.pkt[r.i:])
		buf = buf[n:]
		r.i += n
	}
	// If buf is long enough try read one packet directly to it
	if len(buf) >= PktLen {
		err = r.p.ReadPkt(AsPkt(buf))
		if err == nil {
			n += PktLen
		}
		return
	}
	// If there is place in buf read packet to internal buffer and copy from it
	// to fill buf.
	if len(buf) > 0 {
		err = r.p.ReadPkt(AsPkt(r.pkt[:]))
		if err == nil {
			r.i = copy(buf, r.pkt[:])
			n += r.i
		}
	}
	return
}
