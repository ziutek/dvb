package psi

import (
	"github.com/ziutek/dvb"
)

type PAT struct {
	MuxId   uint16
	Version byte
	Progs   map[uint16]uint16
	Valid   bool
}

func clearProgs(progs map[uint16]uint16) {
	for k := range progs {
		delete(progs, k)
	}
}

var (
	ErrPATSectionSyntax = dvb.TemporaryError("incorrect PAT section syntax")
	ErrPATSectionNumber = dvb.TemporaryError("incorrect PAT section number")
	ErrPATMuxId         = dvb.TemporaryError("incorrect PAT mux id")
	ErrPATDataLength    = dvb.TemporaryError("incorrect PAT data length")
)

// PATDecoder decodes PAT data from stream of sections.
type PATDecoder struct {
	s Section
	r SectionReader
}

func NewPATDecoder(r SectionReader) *PATDecoder {
	return &PATDecoder{s: NewSection(ISOSectionMaxLen), r: r}
}

func (d *PATDecoder) SetSectionReader(r SectionReader) {
	d.r = r
}

// ReadPAT updates p using data from stream of sections provided by internal
// SectionReader. Only sections with Current flag set are processed.
// TODO: This implementation assumes PAT occupies no more than 64 sections
// (standard permits 256 sections). Rewrite it to permit 256 sections.
func (d *PATDecoder) ReadPAT(p *PAT) error {
	var rd uint64
	s := d.s

	for {
		if err := d.r.ReadSection(s); err != nil {
			return err
		}
		if !s.Current() {
			continue
		}

		if s.TableId() != 0 || !s.GenericSyntax() {
			return ErrPATSectionSyntax
		}

		// Always update maxN because sometimes provider update PAT content
		// without update version. In this case we can block waiting for section
		// number that will never appear.
		maxN := s.LastNumber()
		if maxN > 63 {
			maxN = 63 // BUG: this doesn't permit more than 64 sections
		}
		tord := uint64(1) << (maxN - 1)
		n := s.Number()
		if n > maxN {
			return ErrPATSectionNumber
		}
		rd |= uint64(1) << (n - 1)

		muxId := decodeU16(s[3:5])
		if p.Progs == nil {
			// Initial state
			p.MuxId = muxId
			p.Version = s.Version()
			p.Progs = make(map[uint16]uint16)
		} else {
			if p.Version != s.Version() {
				p.MuxId = muxId
				clearProgs(p.Progs)
				p.Valid = false
			} else if p.MuxId != muxId {
				return ErrPATMuxId
			}
		}

		d := s.Data()
		if len(d)%4 != 0 {
			return ErrPATDataLength
		}
		for i := 0; i < len(d); i += 4 {
			p.Progs[decodeU16(d[i:i+2])] = decodeU16(d[i+2:i+4]) & 0x1fff
		}

		if rd&tord == tord {
			break
		}
	}

	p.Valid = true
	return nil
}
