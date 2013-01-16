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

func (p *ArrayPkt) Pid() uint16 {
	return uint16(p[1]&0x1f)<<8 | uint16(p[2])
}

func (p *ArrayPkt) CC() byte {
	return p[3] & 0xf
}

func (p *ArrayPkt) Flags() PktFlags {
	return p.Slice().Flags()
}

func (p *ArrayPkt) AF() AF {
	return p.Slice().AF()
}

func (p *ArrayPkt) Payload() []byte {
	return p.Slice().Payload()
}

// Replacer is interface to replace one packet to another one. After Replace
// old content of p should not be used any more by caller. If Replace returns
// an error it is guaranteed that r == p (but content of p can be modified).
// Generally Replace need to be used in this way:
//
//    p, err = q.Replace(p)
//    if err != nil {
//        ...
//    }
type PktReplacer interface {
	// Replace can return any error but if it returns ErrSync or
	// dvb.ErrOverflow you can try to replace packet one more time.
	Replace(p *ArrayPkt) (r *ArrayPkt, e error)
}
