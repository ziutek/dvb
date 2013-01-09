package psi

type PAT struct {
	MuxId   uint16
	Version byte
	Current bool
	Progs   map[uint16]uint16

	n    byte
	maxN byte
}

func NewPAT() *PAT {
	p := new(PAT)
	p.Progs = make(map[uint16]uint16)
	return p
}

func (p *PAT) resetProgs() {
	for k := range p.Progs {
		delete(p.Progs, k)
	}
}

func (p *PAT) resetState() {
	p.n = 0
	p.maxN = 0
}

func (p *PAT) Reset() {
	p.resetProgs()
	p.resetState()
}

type PATError string

func (e PATError) Error() string {
	return string(e)
}

var (
	ErrPATSectionSyntax = PATError("incorrect PAT section syntax")
	ErrPATSectionNumber = PATError("incorrect PAT section number")
	ErrPATSectionMatch  = PATError("subsequent PAT section doesn't match first one")
	ErrPATDataLength    = PATError("incorrect PAT data length")
)

func (p *PAT) Decode(s Section) (ok bool, err error) {
	defer func() {
		if ok || err != nil {
			p.resetState()
		}
	}()

	if s.TableId() != 0 || !s.GenericSyntax() {
		err = ErrPATSectionSyntax
		return
	}

	muxId := decodeU16(s[3:5])
	if p.n == 0 {
		// Initial state: wait for section_number == 0
		if s.Number() != 0 {
			return
		}
		p.MuxId = muxId
		p.Version = s.Version()
		p.Current = s.Current()
		p.resetProgs()
		p.n = 1
		p.maxN = s.LastNumber()
	} else {
		if s.Number() != p.n || s.LastNumber() != p.maxN {
			err = ErrPATSectionNumber
			return
		}
		if muxId != p.MuxId || p.Version != s.Version() ||
			p.Current != s.Current() {
			err = ErrPATSectionMatch
			return
		}
		p.n++
	}

	d := s.Data()
	if len(d)%4 != 0 {
		err = ErrPATDataLength
		return
	}
	for i := 0; i < len(d); i += 4 {
		p.Progs[decodeU16(d[i:i+2])] = decodeU16(d[i+2:i+4]) & 0x1fff
	}

	if s.Number() == p.maxN {
		ok = true
	}
	return
}
