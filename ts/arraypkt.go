package ts

// ArrayPkt implements Pkt interface and represents array of bytes that
// contains one MPEG-TS packet.
type ArrayPkt [PktLen]byte

// Slice returns refference to the content of p as SlicePkt
func (p *ArrayPkt) Slice() SlicePkt {
	return SlicePkt(p[:])
}

func (p *ArrayPkt) Bytes() []byte {
	return p[:]
}

func (p *ArrayPkt) Copy(pkt Pkt) {
	copy(p[:], pkt.Bytes())
}

func (p *ArrayPkt) SyncOK() bool {
	return p[0] == 0x47
}

func (p *ArrayPkt) SetSync() {
	p[0] = 0x47
}

func (p *ArrayPkt) Pid() int16 {
	return int16(p[1]&0x1f)<<8 | int16(p[2])
}

func (p *ArrayPkt) SetPid(pid int16) {
	if uint(pid) > 8191 {
		panic("Bad PID")
	}
	p[1] = p[1]&0xe0 | byte(pid>>8)
	p[2] = byte(pid)
}

func (p *ArrayPkt) CC() byte {
	return p[3] & 0xf
}

func (p *ArrayPkt) SetCC(b byte) {
	p[3] = p[3]&0xf0 | b&0x0f
}

func (p *ArrayPkt) IncCC() {
	b := p[3]
	p[3] = b&0xf0 | (b+1)&0x0f
}

func (p *ArrayPkt) Flags() PktFlags {
	return PktFlags(p[1]&0xe0 | (p[3] >> 4))
}

func (p *ArrayPkt) SetFlags(f PktFlags) {
	p[1] = p[1]&0x1f | byte(f&0xf0)
	p[3] = p[3]&0x0f | byte(f<<4)
}

func (p *ArrayPkt) AF() AF {
	return p.Slice().AF()
}

func (p *ArrayPkt) Payload() []byte {
	return p.Slice().Payload()
}

func (p *ArrayPkt) ContainsError() bool {
	return p[1]&0x80 != 0
}

func (p *ArrayPkt) SetContainsError(b bool) {
	if b {
		p[1] |= 0x80
	} else {
		p[1] &^= 0x80
	}
}

func (p *ArrayPkt) PayloadStart() bool {
	return p[1]&0x40 != 0
}

func (p *ArrayPkt) SetPayloadStart(b bool) {
	if b {
		p[1] |= 0x40
	} else {
		p[1] &^= 0x40
	}
}

func (p *ArrayPkt) Prio() bool {
	return p[1]&0x20 != 0
}

func (p *ArrayPkt) SetPrio(b bool) {
	if b {
		p[1] |= 0x20
	} else {
		p[1] &^= 0x20
	}
}

func (p *ArrayPkt) ScramblingCtrl() PktScramblingCtrl {
	return PktScramblingCtrl((p[3] >> 6) & 3)
}

func (p *ArrayPkt) SetScramblingCtrl(sc PktScramblingCtrl) {
	p[3] = p[3]&0x3f | byte(sc&3)<<6
}

func (p *ArrayPkt) ContainsAF() bool {
	return p[3]&0x20 != 0
}

func (p *ArrayPkt) SetContainsAF(b bool) {
	if b {
		p[3] |= 0x20
	} else {
		p[3] &^= 0x20
	}
}

func (p *ArrayPkt) ContainsPayload() bool {
	return p[3]&0x10 != 0
}

func (p *ArrayPkt) SetContainsPayload(b bool) {
	if b {
		p[3] |= 0x10
	} else {
		p[3] &^= 0x10
	}
}
