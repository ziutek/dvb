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
	s := make(psi.Section, psi.ISOSectionMaxLen)
	for i := range s[:7-1-4] {
		s[i] = byte(i)
	}
	s.SetReserved(3)
	e := psi.NewSectionEncoder(q, 0x321)

	for l := 7; l < psi.ISOSectionMaxLen; l++ {
		s[l-1-4] = byte(l - 1 - 4)
		s.SetLen(l)
		s.MakeCRC()
		if err := e.WriteSection(s); err != nil {
			t.Fatal(err)
		}
	}
	if err := e.Flush(); err != nil {
		t.Fatal(err)
	}
	q.Close()
}
