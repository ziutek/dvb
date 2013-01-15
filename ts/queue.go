package ts

import (
	"io"
)

// PktQueue represents queue of TS packets.
type PacketQueue struct {
	empty, full chan *Packet
}

// NewPktQueue creates new queue with internall buffer of size length packets.
func NewPacketQueue(length int) PacketQueue {
	q := PacketQueue{
		empty: make(chan *Packet, length),
		full:  make(chan *Packet, length),
	}
	for i := 0; i < length; i++ {
		q.empty <- new(Packet)
	}
	return q
}

func (q PacketQueue) Len() int {
	return len(q.empty)
}

// ReadPart returns read part of q that can be used only to read packets from
// q.
func (q PacketQueue) ReadPart() PacketReadQueue {
	return PacketReadQueue{Empty: q.empty, Full: q.full}
}

// WritePart returns write part of q that can be used to write packets to
// q and to close q.
func (q PacketQueue) WritePart() PacketWriteQueue {
	return PacketWriteQueue{Empty: q.empty, Full: q.full}
}

// PacketReadQueue represenst read part of PacketQueue. If reader uses raw
// channels insteed of Swap method it should write empty buffer to Empty
// channel and next read full buffer from Full channel.
type PacketReadQueue struct {
	Empty chan<- *Packet
	Full  <-chan *Packet
}

func (q PacketReadQueue) Swap(pkt *Packet) (*Packet, error) {
	q.Empty <- pkt
	p, ok := <-q.Full
	if !ok {
		return pkt, io.EOF
	}
	return p, nil
}

// PacketWriteQueue represenst write part of PacketQueue. If writer uses raw
// channels insteed of Swap method it should read empty buffer from Empty
// channel and next write full buffer to Full channel.
type PacketWriteQueue struct {
	Empty <-chan *Packet
	Full  chan<- *Packet
}

// Close closes write part of queue. After close on write part, Swap method
// on read part returns io.EOF error if there is no more packets to read.
func (q PacketWriteQueue) Close() {
	close(q.Full)
}

func (q PacketWriteQueue) Swap(pkt *Packet) (*Packet, error) {
	p := <-q.Empty
	q.Full <- pkt
	return p, nil
}
