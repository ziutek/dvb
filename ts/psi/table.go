package psi

import (
	"github.com/ziutek/dvb"
)

var (
	ErrTableSectionNumber = dvb.TemporaryError("Table: incorrect section number")
	ErrTableSyntax        = dvb.TemporaryError("Table: incorrect section syntax")
)

type Table struct {
	ss         []Section
	m          int
	sectionLen int
}

func NewTable(sectionLen int) *Table {
	return &Table{sectionLen: sectionLen}
}

func (t *Table) Sections() []Section {
	return t.ss[:t.m]
}

func (t *Table) check() {
	if t.m == 0 {
		panic("table doesn't contain valid data")
	}
}

func (t *Table) Reset() {
	t.m = 0
}

func (t *Table) TableId() byte {
	t.check()
	return t.ss[0].TableId()
}

func (t *Table) Version() byte {
	t.check()
	return t.ss[0].Version()
}

func (t *Table) Current() bool {
	t.check()
	return t.ss[0].Current()
}

func (t *Table) TableIdExt() uint16 {
	t.check()
	return t.ss[0].TableIdExt()
}

// Update reads next table from r.
func (t *Table) Update(r SectionReader, tableId byte, private, current bool) error {
	var rd uint64
	t.m = 0
	m := 0
	for {
		var s Section
		if m < len(t.ss) {
			s = t.ss[m]
		} else {
			s = make(Section, t.sectionLen)
			t.ss = append(t.ss, s)
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
		}
	}
	// TODO: sort sections according to the s.Number()
	t.m = m
	return nil
}

type TableDescriptors struct {
	ss     []Section
	dl     DescriptorList
	offset int
}

func (td TableDescriptors) IsEmpty() bool {
	return len(td.ss) == 0 && len(td.dl) == 0
}

// Descriptors handles table global descriptors (if exists). offset is an
// offest from begining of section data part to descriptor length word.
func (t *Table) Descriptors(offset int) TableDescriptors {
	return TableDescriptors{ss: t.Sections(), offset: offset}
}

func (td TableDescriptors) Pop() (d Descriptor, rtd TableDescriptors) {
	if len(td.dl) == 0 {
		if len(td.ss) == 0 {
			return
		}
		data := td.ss[0].Data()
		if len(data) < td.offset+2 {
			return
		}
		l := int(decodeU16(data[td.offset:td.offset+2]) & 0x0fff)
		data = data[td.offset+2:]
		if len(data) < l {
			return
		}
		td.dl = DescriptorList(data[:l])
		td.ss = td.ss[1:]
	}
	d, td.dl = td.dl.Pop()
	rtd = td
	return
}
