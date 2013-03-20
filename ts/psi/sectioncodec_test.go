package psi_test

import (
	"github.com/ziutek/dvb/ts"
	"github.com/ziutek/dvb/ts/psi"
	"io"
	"log"
	"testing"
)

func init() {
}

func TestCodec(t *testing.T) {
	q := ts.NewPktQueue(1)
	go source(t, q.WritePart())
	d := psi.NewSectionDecoder(q.ReadPart())
	s := make(psi.Section, psi.ISOSectionMaxLen)

	for {
		if err := d.ReadSection(s); err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
		}
		log.Println(s.Len())
	}
}

func source(t *testing.T, q *ts.PktWriteQueue) {
	s := make(psi.Section, psi.ISOSectionMaxLen)
	for i := range s {
		s[i] = byte(i)
	}

	e := psi.NewSectionEncoder(q, 0x321)

	for l := 7; l < psi.ISOSectionMaxLen; l++ {
		s.SetLen(l)
		s.MakeCRC()
		if err := e.WriteSection(s); err != nil {
			t.Error(err)
		}
	}
	if err := e.Flush(); err != nil {
		t.Error(err)
	}
	q.Close()
}
