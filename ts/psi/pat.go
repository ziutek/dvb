package psi

import (
	"github.com/ziutek/dvb"
)

type PAT struct {
	Progs   map[uint16]uint16
	MuxId   uint16
	Version byte
	Valid   bool
}

func clearProgs(progs map[uint16]uint16) {
	for k := range progs {
		delete(progs, k)
	}
}

var ErrPATDataLength = dvb.TemporaryError("incorrect PAT data length")

// PATDecoder decodes PAT data from stream of sections.
type PATDecoder struct {
	s Section
	r *TableReader
}

func NewPATDecoder(r SectionReader) *PATDecoder {
	return &PATDecoder{
		s: NewSection(ISOSectionMaxLen),
		r: NewTableReader(r, 0, true),
	}
}

func (d *PATDecoder) SetSectionReader(r SectionReader) {
	d.r.SetSectionReader(r)
}

// ReadPAT updates pat using data from stream of sections provided by internal
// SectionReader. Only sections with Current flag set are processed.
// If ReadPAT returns error pat.Valid == false, otherwise pat.Valid == true.
func (d *PATDecoder) ReadPAT(pat *PAT) error {
	s := d.s

	for {
		done, err := d.r.ReadTableSection(s)
		if err != nil {
			return err
		}

		muxId := s.TableIdExt()
		version := s.Version()
		if pat.Progs == nil {
			// Initial state
			pat.MuxId = muxId
			pat.Version = version
			pat.Progs = make(map[uint16]uint16)
		} else if pat.Version != s.Version() || pat.MuxId != muxId {
			pat.MuxId = muxId
			pat.Version = version
			clearProgs(pat.Progs)
			pat.Valid = false
		}

		d := s.Data()
		if len(d)%4 != 0 {
			return ErrPATDataLength
		}
		for i := 0; i < len(d); i += 4 {
			pat.Progs[decodeU16(d[i:i+2])] = decodeU16(d[i+2:i+4]) & 0x1fff
		}

		if done {
			break
		}
	}

	pat.Valid = true
	return nil
}
