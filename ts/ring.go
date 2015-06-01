package ts

// PktRing represents ring buffer of TS packets. PktRing is allways ready
// to insert new packet. If PktRing is full oldest packet is overwritten. If
// reader reads first overwritten packet ErrSync is returned.
type PktRing struct {
	empty, filled chan *ArrayPkt
}

// ReadPart returns read part of r that can be used only to read packets from
// r.
func (r *PktRing) ReadPart() *PktReadRing {
	return (*PktReadRing)(r)
}

// WritePart returns write part of r that can be used to write packets to
// r and to close r.
func (r *PktRing) WritePart() *PktWriteRing {
	return (*PktWriteRing)(r)
}

// TODO:
type PktWriteRing PktRing

// TODO:
type PktReadRing PktRing