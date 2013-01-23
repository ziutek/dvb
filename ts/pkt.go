package ts

const (
	PktLen = 188
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
	// Pid return value of packet identifier
	Pid() uint16
	// CC returns value of continuity counter
	CC() byte
	// Flags returns packet flags
	Flags() PktFlags
	// AF returns adaptation field bytes. If p doesn't contain AF it returns
	// AF{}. If adaptation_field_length byte has wrong value it returns nil.
	AF() AF
	// Payload returns payload bytes. It returns nil if packet dosn't contain
	// payload or adaptation_field_length byte has incorrect value.
	Payload() []byte
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
