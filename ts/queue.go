package ts

import (
	"io"
)

// PktQueue represents queue of TS packets.
type PktQueue struct {
	empty, full chan *ArrayPkt
}

// NewPktQueue creates new queue with internall buffer of size length packets.
func NewPktQueue(length int) PktQueue {
	q := PktQueue{
		empty: make(chan *ArrayPkt, length),
		full:  make(chan *ArrayPkt, length),
	}
	for i := 0; i < length; i++ {
		q.empty <- new(ArrayPkt)
	}
	return q
}

func (q PktQueue) Len() int {
	return len(q.empty)
}

// ReadPart returns read part of q that can be used only to read packets from
// q.
func (q PktQueue) ReadPart() PktReadQueue {
	return PktReadQueue{Empty: q.empty, Full: q.full}
}

// WritePart returns write part of q that can be used to write packets to
// q and to close q.
func (q PktQueue) WritePart() PktWriteQueue {
	return PktWriteQueue{Empty: q.empty, Full: q.full}
}

// PacketReadQueue represenst read part of PktQueue. If reader uses raw
// channels insteed of Replace method it should write empty packet to Empty
// channel and next read full packet from Full channel.
type PktReadQueue struct {
	Empty chan<- *ArrayPkt
	Full  <-chan *ArrayPkt
}

func (q PktReadQueue) Replace(pkt *ArrayPkt) (*ArrayPkt, error) {
	q.Empty <- pkt
	p, ok := <-q.Full
	if !ok {
		return pkt, io.EOF
	}
	return p, nil
}

// PacketWriteQueue represenst write part of PktQueue. If writer uses raw
// channels insteed of Replace method it should read empty packet from Empty
// channel and next write full packet to Full channel.
type PktWriteQueue struct {
	Empty <-chan *ArrayPkt
	Full  chan<- *ArrayPkt
}

// Close closes write part of queue. After close on write part, Replace method
// on read part returns io.EOF error if there is no more packets to read.
func (q PktWriteQueue) Close() {
	close(q.Full)
}

func (q PktWriteQueue) Replace(pkt *ArrayPkt) (*ArrayPkt, error) {
	p := <-q.Empty
	q.Full <- pkt
	return p, nil
}
