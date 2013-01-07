package ts

const PktLen = 188

type Pkt []byte

func NewPkt() Pkt {
	return make(Pkt, PktLen)
}

func AsPkt(buf []byte) Pkt {
	if len(buf) < PktLen {
		panic("buffer too small to be treat as MPEG-TS packet")
	}
	return Pkt(buf[:PktLen])
}

func (p Pkt) SyncOK() bool {
	return p[0] == 0x47
}

func (p Pkt) Pid() uint16 {
	return uint16(p[1]&0x1f)<<8 | uint16(p[2])
}

func (p Pkt) CC() byte {
	return p[3] & 0xf
}

func (p Pkt) Flags() PktFlags {
	return PktFlags(p[1]&0xe0 | (p[3] >> 4))
}

type PktFlags byte

// Error returns true if transport_error_indicator == 1
func (f PktFlags) ContainsError() bool {
	return f&0x80 != 0
}

// PayloadStart returns true if payload_unit_start_indicator == 1
func (f PktFlags) PayloadStart() bool {
	return f&0x40 != 0
}

// Prio returns true if transport_priority == 1
func (f PktFlags) Prio() bool {
	return f&0x20 != 0
}

type PktScramblingCtrl byte

const (
	PktNotScrambled PktScramblingCtrl = iota
	PktScrambled1
	PktScrambled2
	PktScrambled3
)

func (f PktFlags) ScramblingCtrl() PktScramblingCtrl {
	return PktScramblingCtrl((f >> 2) & 3)
}

// ContainsAF returns true if adaptation_field_control & 2 == 1
func (f PktFlags) ContainsAF() bool {
	return f&2 != 0
}

// ContainsPayload returns true if adaptation_field_control & 1 == 1
func (f PktFlags) ContainsPayload() bool {
	return f&1 != 0
}

// AF returns adaptation field bytes It returns nil if ContainsAF() == false
// or adaptation_field_length byte is too big
func (p Pkt) AF() AF {
	if !p.Flags().ContainsAF() {
		return nil
	}
	alen := p[4]
	if p.Flags().ContainsPayload() {
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
