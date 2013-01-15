package ts

// Packet implements Pkt interface and represents independent array of
// bytes that contains one MPEG-TS packet.
type Packet [PktLen]byte

// Slice returns refference to the content of p as SlicePkt
func (p *Packet) Slice() SlicePkt {
	return SlicePkt(p[:])
}

func (p *Packet) Bytes() []byte {
	return p[:]
}

func (p *Packet) Copy(pkt Pkt) {
	copy(p[:], pkt.Bytes())
}

func (p *Packet) SyncOK() bool {
	return p[0] == 0x47
}

func (p *Packet) Pid() uint16 {
	return uint16(p[1]&0x1f)<<8 | uint16(p[2])
}

func (p *Packet) CC() byte {
	return p[3] & 0xf
}

func (p *Packet) Flags() PktFlags {
	return p.Slice().Flags()
}

func (p *Packet) AF() AF {
	return p.Slice().AF()
}

func (p *Packet) Payload() []byte {
	return p.Slice().Payload()
}

// PktSwapper is interface to swap/exchange one independent packet to
// another one. If SwapPkt returns an error it is guaranteed that r == p so
// you can safely use it in this way:
//
//    p, err = q.SwapPkt(p)
//    if err != nil {
//        ...
//    }
type PacketSwapper interface {
	Swap(p *Packet) (r *Packet, e error)
}
