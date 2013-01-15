package ts

// SlicePkt implements Pkt interface and represents MPEG-TS packet that can be
// a slice of some more data (eg. slice of buffer that contains more packets)
type SlicePkt []byte

// AsPkt returns beginning of buf as SlicePkt. It panics if len(buf) < PktLen.
func AsPkt(buf []byte) SlicePkt {
	if len(buf) < PktLen {
		panic("buffer too small to be treat as MPEG-TS packet")
	}
	return SlicePkt(buf[:PktLen])
}

func (p SlicePkt) Bytes() []byte {
	if len(p) != PktLen {
		panic("wrong MPEG-TS packet length")
	}
	return p
}

func (p SlicePkt) SyncOK() bool {
	return p[0] == 0x47
}

func (p SlicePkt) Pid() uint16 {
	return uint16(p[1]&0x1f)<<8 | uint16(p[2])
}

func (p SlicePkt) CC() byte {
	return p[3] & 0xf
}

func (p SlicePkt) Flags() PktFlags {
	return PktFlags(p[1]&0xe0 | (p[3] >> 4))
}

func (p SlicePkt) AF() AF {
	f := p.Flags()
	if !f.ContainsAF() {
		return AF{}
	}
	alen := p[4]
	if f.ContainsPayload() {
		if alen > 182 {
			return nil
		}
	} else {
		if alen != 183 {
			return nil
		}
	}
	return AF(p[5 : 5+alen])
}

func (p SlicePkt) Payload() []byte {
	f := p.Flags()
	if !f.ContainsPayload() {
		return nil
	}
	offset := 4
	if f.ContainsAF() {
		af := p.AF()
		if af == nil {
			return nil
		}
		offset += len(af) + 1
	}
	return p[offset:]
}
