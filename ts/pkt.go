package ts

const (
	PktLen  = 188
	NullPid = 8191
)

// Pkt is common interface to any MPEG-TS packet implementation
type Pkt interface {
	// Bytes returns content of the packet as byte slice. It guarantees that
	// len of slice is equal to PktLen.
	Bytes() []byte
	// Copy copies conntent of Pkt.
	Copy(Pkt)
	// SyncOK checks does sync byte is OK.
	SyncOK() bool
	// SetSync() sets right sync byte
	SetSync()
	// Pid return value of packet identifier
	Pid() int16
	// SetPid sets packet identifier
	SetPid(int16)
	// CC returns value of continuity counter
	CC() byte
	// SetCC sets the value of continuity counter to byte&0x0f
	SetCC(byte)
	// IncCC increments  continuity counter
	IncCC()
	// Flags returns packet flags
	Flags() PktFlags
	// SetFlags sets packet flags
	SetFlags(PktFlags)
	// AF returns adaptation field bytes. If p doesn't contain AF it returns
	// AF{}. If adaptation_field_length byte has wrong value it returns nil.
	AF() AF
	// Payload returns payload bytes. It returns nil if packet dosn't contain
	// payload or adaptation_field_length byte has incorrect value.
	Payload() []byte

	// Direct access to options/flags
	ContainsError() bool
	SetContainsError(b bool)
	PayloadStart() bool
	SetPayloadStart(b bool)
	Prio() bool
	SetPrio(bool)
	ScramblingCtrl() PktScramblingCtrl
	SetScramblingCtrl(PktScramblingCtrl)
	ContainsAF() bool
	SetContainsAF(bool)
	ContainsPayload() bool
	SetContainsPayload(bool)
}

type PktFlags byte

// ContainsError returns true if transport_error_indicator == 1
func (f PktFlags) ContainsError() bool {
	return f&0x80 != 0
}

// SetContainsError sets transport_error_indicator
func (f *PktFlags) SetContainsError(b bool) {
	if b {
		*f |= 0x80
	} else {
		*f &^= 0x80
	}
}

// PayloadStart returns true if payload_unit_start_indicator == 1
func (f PktFlags) PayloadStart() bool {
	return f&0x40 != 0
}

// SetPayloadStart sets payload_unit_start_indicator
func (f *PktFlags) SetPayloadStart(b bool) {
	if b {
		*f |= 0x40
	} else {
		*f &^= 0x40
	}
}

// Prio returns true if transport_priority == 1
func (f PktFlags) Prio() bool {
	return f&0x20 != 0
}

func (f *PktFlags) SetPrio(b bool) {
	if b {
		*f |= 0x20
	} else {
		*f &^= 0x20
	}
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

func (f *PktFlags) SetScramblingCtrl(sc PktScramblingCtrl) {
	*f = *f&0xf3 | PktFlags(sc&3)<<2
}

// ContainsAF returns true if adaptation_field_control & 2 == 1
func (f PktFlags) ContainsAF() bool {
	return f&2 != 0
}

// SetContainsAF sets first bit of adaptation_field_control
func (f *PktFlags) SetContainsAF(b bool) {
	if b {
		*f |= 2
	} else {
		*f &^= 2
	}
}

// ContainsPayload returns true if adaptation_field_control & 1 == 1
func (f PktFlags) ContainsPayload() bool {
	return f&1 != 0
}

// SetContainsPayload sets second bit of daptation_field_control
func (f *PktFlags) SetContainsPayload(b bool) {
	if b {
		*f |= 1
	} else {
		*f &^= 1
	}
}
