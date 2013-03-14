package psi

import (
	"io"
)

// SectionReader is an interface that wraps the ReadSection method.
type SectionReader interface {
	// ReadSection reads one section into s. len(s) should be equal to
	// MaxSectionLen or MaxISOSectionLen (if you read standard PSI tables). You
	// can use shorter s (but not shorter that 8 bytes) if you are sure that
	// read section should fit in it. If ReadSection returned error of
	// dvb.TemporaryError type you can try read next section.
	ReadSection(Section) error
}

// SectionStreamReader wraps any io.Reader interface and returns SectionReader
// implementation for read MPEG-TS sections from stream of bytes.
type SectionStreamReader struct {
	r        io.Reader
	checkCRC bool
}

func NewSectionStreamReader(r io.Reader, checkCRC bool) *SectionStreamReader {
	return &SectionStreamReader{r, checkCRC}
}

func (sr *SectionStreamReader) ReadSection(s Section) error {
	if len(s) < 8 {
		panic("section length should be >= 8")
	}
	_, err := io.ReadFull(sr.r, s[:3])
	if err != nil {
		return err
	}
	l := s.Len()
	if l == -1 {
		return ErrSectionLength
	}
	if l > len(s) {
		return ErrSectionSpace
	}
	_, err = io.ReadFull(sr.r, s[3:l])
	if err != nil {
		return err
	}
	if sr.checkCRC && !s.CheckCRC() {
		return ErrSectionCRC
	}
	return nil
}
