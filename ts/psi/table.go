package psi

import (
	"github.com/ziutek/dvb"
)

var (
	ErrTableSectionNumber = dvb.TemporaryError("Table: incorrect section number")
	ErrTableSyntax        = dvb.TemporaryError("Table: incorrect section syntax")
)

type Table []Section

func (t Table) check() {
	if len(t) == 0 {
		panic("table doesn't contain valid data")
	}
}

func (t *Table) Reset() {
	*t = (*t)[:0]
}

func (t Table) TableId() byte {
	t.check()
	return t[0].TableId()
}

func (t Table) SetTableId(id byte) {
	t.check()
	for _, s := range t {
		s.SetTableId(id)
	}
}

func (t Table) Version() int8 {
	t.check()
	return t[0].Version()
}

func (t Table) Current() bool {
	t.check()
	return t[0].Current()
}

func (t Table) TableIdExt() uint16 {
	t.check()
	return t[0].TableIdExt()
}

// Update reads next table from r.
func (t *Table) Update(r SectionReader, tableId byte, private, current bool, sectionMaxLen int) error {
	var rd uint64
	t.Reset()
	m := 0
	for {
		var s Section
		if m < cap(*t) {
			*t = (*t)[:m+1]
			s = (*t)[m]
		} else {
			s = make(Section, sectionMaxLen)
			*t = append(*t, s)
		}
		if err := r.ReadSection(s); err != nil {
			return err
		}
		if s.TableId() != tableId || s.Current() != current {
			continue
		}
		if !s.GenericSyntax() || s.PrivateSyntax() != private {
			return ErrTableSyntax
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
			return ErrTableSectionNumber
		}

		bitn := uint64(1) << (n - 1)

		if rd&bitn != 0 {
			// Section readed before
			continue
		}

		// New section readed
		m++

		rd |= bitn
		if rd&tord == tord {
			// All sections readed
			break
		}

		if m == 1 {
			continue
		}
		if s.Version() != t.Version() || s.TableIdExt() != t.TableIdExt() {
			// Old table can never appear
			rd, m = 0, 0
			t.Reset()
		}
	}
	// TODO: sort sections according to the s.Number()
	return nil
}

// Cursor returns TableCursor that can be used to obtain data from table.
func (t Table) Cursor() TableCursor {
	return TableCursor{Tab: t}
}

// Descriptors handles table global descriptors (if exists). offset is an
// offest from begining of section data part to descriptor length word.
func (t Table) Descriptors(offset int) TableDescriptors {
	return TableDescriptors{tab: t, offset: offset}
}

type TableDescriptors struct {
	tab    Table
	dl     DescriptorList
	offset int
}

func (td TableDescriptors) IsEmpty() bool {
	return len(td.tab) == 0 && len(td.dl) == 0
}

func (td TableDescriptors) Pop() (d Descriptor, rtd TableDescriptors) {
	if len(td.dl) == 0 {
		if len(td.tab) == 0 {
			return
		}
		data := td.tab[0].Data()
		if len(data) < td.offset+2 {
			return
		}
		l := int(decodeU16(data[td.offset:td.offset+2]) & 0x0fff)
		data = data[td.offset+2:]
		if len(data) < l {
			return
		}
		td.dl = DescriptorList(data[:l])
		td.tab = td.tab[1:]
	}
	d, td.dl = td.dl.Pop()
	rtd = td
	return
}

type TableCursor struct {
	Tab  Table
	Data []byte
}

func (tc TableCursor) IsEmpty() bool {
	return len(tc.Data) == 0 && len(tc.Tab) == 0
}

func (tc TableCursor) NextSection() TableCursor {
	tc.Data = tc.Tab[0].Data()
	tc.Tab = tc.Tab[1:]
	return tc
}

func (t *Table) SetEmpty() {
	*t = (*t)[:0]
}

type TableConfig struct {
	TableId        byte
	GenericSyntax  bool
	PrivateSyntax  bool
	SectionMaxLen  int
	SectionHeadLen int
	NumLenFields   int // Number of length fields in section.
}

func getlf(sec Section, hlen, uself int) []byte {
	data := sec.Data()[hlen:]
	for uself > 0 && len(data) > 0 {
		lf := loopLen(data[:2])
		data = data[2+lf:]
		uself--
	}
	if len(data) == 0 {
		lfs := sec.Alloc((uself+1)*2, 0)
		for i := 0; i < len(lfs); i += 2 {
			clearLoopLen(lfs[i : i+2])
		}
		return lfs[len(lfs)-2:]
	}
	return data[:2]
}

func lfadd(sec Section, hlen, uself, n int) {
	lf := getlf(sec, hlen, uself)
	setLoopLen(lf, loopLen(lf)+n)
}

func (t *Table) Alloc(n int, cfg *TableConfig, uself int, sectionHeader []byte) []byte {
	var (
		sec  Section
		data []byte
	)
	postlf := (cfg.NumLenFields - uself - 1) * 2
	if postlf < 0 {
		postlf = 0
	}
	m := len(*t)
	if m > 0 {
		sec = (*t)[m-1]
		if cfg.NumLenFields > 0 {
			lfadd(sec, cfg.SectionHeadLen, uself, n)
		}
		data = sec.Alloc(n, postlf)
	}
	if sec == nil || data == nil {
		if sec != nil && postlf > 0 {
			// Insert zero length fields after data.
			tail := sec.Alloc(postlf, 0)
			for i := 0; i < postlf; i += 2 {
				clearLoopLen(tail[i : i+2])
			}
		}
		if m < cap(*t) {
			*t = (*t)[:m+1]
			sec = (*t)[m]
			sec.SetEmpty()
		} else {
			sec = MakeEmptySection(cfg.SectionMaxLen, cfg.GenericSyntax)
			sec.SetPrivateSyntax(cfg.PrivateSyntax)
			*t = append(*t, sec)
		}
		if cfg.SectionHeadLen > 0 {
			head := sec.Alloc(cfg.SectionHeadLen, postlf)
			copy(head, sectionHeader)
		}
		if cfg.NumLenFields > 0 {
			lfadd(sec, cfg.SectionHeadLen, uself, n)
		}
		data = sec.Alloc(n, postlf)
	}
	return data
}

// Close recalculates section numbers and makes CRC sums for all sections.
func (t Table) Close(cfg *TableConfig, tableIdExt uint16, current bool, version int8) {
	lastnum := byte(len(t) - 1)
	if cfg.NumLenFields > 0 {
		// Properly initialize unused length fields in last section.
		getlf(t[lastnum], cfg.SectionHeadLen, cfg.NumLenFields-1)
	}
	for num, s := range t {
		s.SetTableId(cfg.TableId)
		s.SetTableIdExt(tableIdExt)
		s.SetCurrent(current)
		s.SetVersion(version)
		s.SetNumber(byte(num))
		s.SetLastNumber(lastnum)
		s.MakeCRC()
	}
}
