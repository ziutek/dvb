package psi

import (
	"github.com/ziutek/dvb"
)

type SDT struct {
	MuxId    uint16
	OrgNetId uint16
	Version  byte
	Valid    bool
}

var ErrSDTSectionSyntax = dvb.TemporaryError("incorrect SDT section syntax")

type SDTDecoder struct {
	s      Section
	r      *TableReader
	actual bool
}

func NewSDTDecoder(r SectionReader, actual bool) *SDTDecoder {
	tableId := byte(0x46)
	if actual {
		tableId = 0x42
	}
	return &SDTDecoder{
		s:      NewSection(ISOSectionMaxLen),
		r:      NewTableReader(r, tableId, true),
		actual: actual,
	}
}

// ReadSDT updates sdt using data from stream of sections provided by internal
// SectionReader. Only sections with Current flag set are processed.
// If ReadSDT returns error sdt.Valid == false, otherwise sdt.Valid == true.
func (d *SDTDecoder) ReadSDT(sdt *SDT) error {
	s := d.s
	for {
		done, err := d.r.ReadTableSection(s)
		if err != nil {
			return err
		}

		if done {
			break
		}
	}
	return nil
}

func (d *SDTDecoder) update(s Section) error {
	return nil
}
