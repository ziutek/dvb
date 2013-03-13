package ts

// SlicePkt implements Pkt interface and represents MPEG-TS packet that can be
// a slice of some more data (eg. slice of buffer that contains more packets).
// Use it only when you can't use ArrayPkt. Assigning variable of type SlicePkt
// to variable of type Pkt causes memory allocation (you can use *SlicePkt to
// avoid this).
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

func (p SlicePkt) Copy(pkt Pkt) {
	copy(p, pkt.Bytes())
}

func (p SlicePkt) SyncOK() bool {
	return p[0] == 0x47
}

func (p SlicePkt) SetSync() {
	p[0] = 0x47
}

func (p SlicePkt) Pid() uint16 {
	return uint16(p[1]&0x1f)<<8 | uint16(p[2])
}

func (p SlicePkt) SetPid(pid uint16) {
	if pid > 8191 {
		panic("Bad PID")
	}
	p[1] = p[1]&0xe0 | byte(pid>>8)
	p[2] = byte(pid)
}

func (p SlicePkt) CC() byte {
	return p[3] & 0x0f
}

func (p SlicePkt) SetCC(b byte) {
	p[3] = p[3]&0xf0 | b&0x0f
}

func (p SlicePkt) IncCC() {
	b := p[3]
	p[3] = b&0xf0 | (b+1)&0x0f
}

func (p SlicePkt) Flags() PktFlags {
	return PktFlags(p[1]&0xe0 | (p[3] >> 4))
}

func (p SlicePkt) SetFlags(f PktFlags) {
	p[1] = p[1]&0x1f | byte(f&0xf0)
	p[3] = p[3]&0x0f | byte(f<<4)
}

func (p SlicePkt) AF() AF {
	if !p.ContainsAF() {
		return AF{}
	}
	alen := p[4]
	if p.ContainsPayload() {
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
	if !p.ContainsPayload() {
		return nil
	}
	offset := 4
	if p.ContainsAF() {
		af := p.AF()
		if af == nil {
			return nil
		}
		offset += len(af) + 1
	}
	return p[offset:]
}

func (p SlicePkt) ContainsError() bool {
	return p[1]&0x80 != 0
}

func (p SlicePkt) SetContainsError(b bool) {
	if b {
		p[1] |= 0x80
	} else {
		p[1] &^= 0x80
	}
}

func (p SlicePkt) PayloadStart() bool {
	return p[1]&0x40 != 0
}

func (p SlicePkt) SetPayloadStart(b bool) {
	if b {
		p[1] |= 0x40
	} else {
		p[1] &^= 0x40
	}
}

func (p SlicePkt) Prio() bool {
	return p[1]&0x20 != 0
}

func (p SlicePkt) SetPrio(b bool) {
	if b {
		p[1] |= 0x20
	} else {
		p[1] &^= 0x20
	}
}

func (p SlicePkt) ScramblingCtrl() PktScramblingCtrl {
	return PktScramblingCtrl((p[3] >> 6) & 3)
}

func (p SlicePkt) SetScramblingCtrl(sc PktScramblingCtrl) {
	p[3] = p[3]&0x3f | byte(sc&3)<<6
}

func (p SlicePkt) ContainsAF() bool {
	return p[3]&0x20 != 0
}

func (p SlicePkt) SetContainsAF(b bool) {
	if b {
		p[3] |= 0x20
	} else {
		p[3] &^= 0x20
	}
}

func (p SlicePkt) ContainsPayload() bool {
	return p[3]&0x10 != 0
}

func (p SlicePkt) SetContainsPayload(b bool) {
	if b {
		p[3] |= 0x10
	} else {
		p[3] &^= 0x10
	}
}
