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
	// If it returns ErrSync or dvb.ErrOverflow you can try to read next
	// packet.
	ReadPkt(Pkt) error
}

// PktStreamReader wraps any io.Reader interface and returns PktReader and
// PktReplacer implementation for read MPEG-TS packets from stream of bytes.
// Internally it doesn't allocate any memory so is friendly for real-time
// applications (it doesn't cause GC to run).
//
// Using PktStreamReader you can start read at any point in stream. If the start point
// doesn't match a beginning of a packet, PktReader returns ErrSync and
// tries to synchronize during next read.
type PktStreamReader struct {
	r       io.Reader
	syncBuf [3 * PktLen]byte
	sbStart int
}

// SetReader sets new io.Reader as stream source. Forces resynchronization.
func (s *PktStreamReader) SetReader(r io.Reader) {
	s.r = r
	s.sbStart = -1
}

// NewPktStreamReader is equivalent to:
//	s := new(PktStreamReader); s.SetReader(r)
func NewPktStreamReader(r io.Reader) *PktStreamReader {
	s := new(PktStreamReader)
	s.SetReader(r)
	return s
}

func readFull(r io.Reader, b []byte) error {
	/*
		for len(b) > 0 {
			m, err := r.Read(b)
			b = b[m:]
			if err != nil {
				switch e := err.(type) {
				case syscall.Errno:
					if e != syscall.EINTR {
						return err
					}
				case *os.PathError:
					if e.Err != syscall.EINTR {
						return err
					}
				case *net.OpError:
					if e.Err != syscall.EINTR {
						return err
					}
				default:
					return err
				}
				time.Sleep(time.Second)
			}
		}
		return nil
	*/
	_, err := io.ReadFull(r, b)
	return err
}

func (s *PktStreamReader) synchronize() (err error) {
	b := s.syncBuf[:]
	if s.sbStart == -1 {
		// First try of synchronization - read full buffer (three packets)
		err = readFull(s.r, b)
		s.sbStart = -2
	} else {
		// Subsequent try of synchronization - read next packet
		copy(b, b[PktLen:])
		err = readFull(s.r, b[2*PktLen:])
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
			err = readFull(s.r, b[len(b)-i:])
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
func (s *PktStreamReader) ReadPkt(pkt Pkt) error {
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
	if err := readFull(s.r, pkt.Bytes()); err != nil {
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
func (s *PktStreamReader) ReplacePkt(pkt *ArrayPkt) (*ArrayPkt, error) {
	err := s.ReadPkt(pkt)
	return pkt, err
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

// PktPktReader wraps any io.Reader interface and returns PktReader and
// PktReplacer implementation for read MPEG-TS packets from stream of packets.
// Internally it doesn't allocate any memory so is friendly for real-time
// applications (it doesn't cause GC to run).
//
// PktPktReader assumes that every call of io.Reader's Read method returns
// one transport packet that fits in provided buffer and contains integer
// number of MPEG-TS packets that are properly aligned to the begining of
// transport packet.
type PktPktReader struct {
	r    io.Reader
	buf  []byte
	next int
}

// NewPktPktReader is equivalent to:
//	p := new(PktPktReader); p.SetReader(r); p.SetBuffer(buf)
func NewPktPktReader(r io.Reader, buf []byte) *PktPktReader {
	p := new(PktPktReader)
	p.SetReader(r)
	p.SetBuffer(buf)
	return p
}

// SetReader sets new io.Reader as stream source. Forces resynchronization.
func (p *PktPktReader) SetReader(r io.Reader) {
	p.r = r
}

// SetReader sets new io.Reader as stream source. Forces resynchronization.
func (p *PktPktReader) SetBuffer(buf []byte) {
	p.buf = buf[:0]
}

// ReadPkt reads one MPEG-TS packet
func (p *PktPktReader) ReadPkt(pkt Pkt) error {
	for {
		for {
			end := p.next + PktLen
			if end > len(p.buf) {
				break
			}
			newpkt := SlicePkt(p.buf[p.next:end])
			p.next = end
			if newpkt.SyncOK() {
				pkt.Copy(newpkt)
				return nil
			} else {
				return ErrSync
			}
		}
		p.buf = p.buf[:cap(p.buf)]
		n, err := p.r.Read(p.buf)
		if err != nil {
			// This isn't correct if n!=0 but simplifies implementation.
			return err
		}
		if n == 0 {
			return io.EOF
		}
		p.buf = p.buf[:n]
		p.next = 0
	}
}
