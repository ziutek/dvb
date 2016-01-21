package psi_test

import (
	"github.com/ziutek/dvb/ts"
	"github.com/ziutek/dvb/ts/psi"
	"io"
	"testing"
)

func init() {
}

func TestCodec(t *testing.T) {
	q := ts.NewPktQueue(1)
	go source(t, q.WritePart())
	d := psi.NewSectionDecoder(q.ReadPart(), true)
	s := make(psi.Section, psi.ISOSectionMaxLen)

	for {
		if err := d.ReadSection(s); err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
	}
}

func source(t *testing.T, q *ts.PktWriteQueue) {
	s := psi.MakeEmptySection(psi.ISOSectionMaxLen, true)
	s.SetTableId(79)
	s.SetTableIdExt(333)
	s.SetVersion(3)
	s.SetCurrent(true)
	s.SetNumber(1)
	s.SetLastNumber(1)
	e := psi.NewSectionEncoder(q, 0x321)
	for n := 4000; n > 0; n-- {
		buf := s.Alloc(n, n)
		for i := range buf {
			buf[i] = byte(n + i)
		}
		s.MakeCRC()
		if err := e.WriteSection(s); err != nil {
			t.Fatal(err)
		}
		s.SetEmpty()
	}
	if err := e.Flush(); err != nil {
		t.Fatal(err)
	}
	q.Close()
}
