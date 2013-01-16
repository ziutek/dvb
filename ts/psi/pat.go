package psi

import (
	"github.com/ziutek/dvb"
)

type PAT struct {
	MuxId   uint16
	Version byte
	Current bool
	Progs   map[uint16]uint16
}

func NewPAT() *PAT {
	p := new(PAT)
	p.Progs = make(map[uint16]uint16)
	return p
}

func clearProgs(progs map[uint16]uint16) {
	for k := range progs {
		delete(progs, k)
	}
}

var (
	ErrPATSectionSyntax = dvb.TemporaryError("incorrect PAT section syntax")
	ErrPATSectionNumber = dvb.TemporaryError("incorrect PAT section number")
	ErrPATSectionMatch  = dvb.TemporaryError("subsequent PAT section doesn't match first one")
	ErrPATDataLength    = dvb.TemporaryError("incorrect PAT data length")
)

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

// Update updates p using r to read data. Returns true if there are any changes
// in p.
func (d *PATDecoder) ReadPAT(p *PAT) error {
	var n, maxN byte
	s := d.s

	for n <= maxN {
		if err := d.r.ReadSection(s); err != nil {
			return err
		}

		if s.TableId() != 0 || !s.GenericSyntax() {
			return ErrPATSectionSyntax
		}

		muxId := decodeU16(s[3:5])
		if n == 0 {
			// Initial state: wait for section_number == 0
			if s.Number() != 0 {
				continue
			}
			p.MuxId = muxId
			p.Version = s.Version()
			p.Current = s.Current()
			clearProgs(p.Progs)
			n = 1
			maxN = s.LastNumber()
		} else {
			if s.Number() != n || s.LastNumber() != maxN {
				return ErrPATSectionNumber
			}
			if muxId != p.MuxId || p.Version != s.Version() || p.Current != s.Current() {
				return ErrPATSectionMatch
			}
			n++
		}

		d := s.Data()
		if len(d)%4 != 0 {
			return ErrPATDataLength
		}
		for i := 0; i < len(d); i += 4 {
			p.Progs[decodeU16(d[i:i+2])] = decodeU16(d[i+2:i+4]) & 0x1fff
		}
	}

	return nil
}
