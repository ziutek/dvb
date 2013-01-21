package psi

import (
	"github.com/ziutek/dvb"
)

var (
	ErrTableSectionNumber  = dvb.TemporaryError("incorrect PSI table section number")
	ErrTableTSectionSyntax = dvb.TemporaryError("incorrect PSI section syntax")
)

type TableReader struct {
	r       SectionReader
	tableId byte
	current bool

	rd         uint64
	tableIdExt uint16
	version    byte
}

func NewTableReader(r SectionReader, tableId byte, current bool) *TableReader {
	return &TableReader{r: r, tableId: tableId, current: current}
}

// ReadTableSecion returns true if it reads all sections for specified table. It
// can read some sections from previous table so you should check Version and
// TableIdExt.
func (tr *TableReader) ReadTableSection(s Section) (done bool, err error) {
	if err = tr.r.ReadSection(s); err != nil {
		return
	}
	if s.TableId() != tr.tableId || s.Current() != tr.current {
		return
	}
	if !s.GenericSyntax() {
		err = ErrTableTSectionSyntax
		return
	}

	// Always update maxN because sometimes provider updates table content
	// without update version. In this case we can block waiting for section
	// number that will never appear.
	maxN := s.LastNumber()
	if maxN > 63 {
		maxN = 63 // BUG: this doesn't permit more than 64 sections
	}
	tord := uint64(1) << (maxN - 1)
	n := s.Number()
	if n > maxN {
		err = ErrTableSectionNumber
		return
	}

	if tr.rd == 0 || s.TableIdExt() != tr.tableIdExt ||
		s.Version() != tr.version {
		// Initialize tableIdExt, version and reset rd if first time oro
		// some changes occurs during read.
		tr.tableIdExt = s.TableIdExt()
		tr.version = s.Version()
		tr.rd = 0
	}

	tr.rd |= uint64(1) << (n - 1)

	done = tr.rd&tord == tord
	return
}
