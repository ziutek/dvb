package ts

import (
	"io"
)

// PktQueue represents queue of TS packets.
type PktQueue struct {
	empty, full chan *ArrayPkt
}

// NewPktQueue creates new queue with internall buffer of size length packets.
func NewPktQueue(length int) *PktQueue {
	q := &PktQueue{
		empty: make(chan *ArrayPkt, length),
		full:  make(chan *ArrayPkt, length),
	}
	for i := 0; i < length; i++ {
		q.empty <- new(ArrayPkt)
	}
	return q
}

func (q *PktQueue) Len() int {
	return len(q.empty)
}

// ReadPart returns read part of q that can be used only to read packets from
// q.
func (q *PktQueue) ReadPart() *PktReadQueue {
	return (*PktReadQueue)(q)
}

// WritePart returns write part of q that can be used to write packets to
// q and to close q.
func (q *PktQueue) WritePart() *PktWriteQueue {
	return (*PktWriteQueue)(q)
}

// PacketReadQueue represenst read part of PktQueue and implements PktReplacer
// interface. If reader uses raw channels insteed of ReplacePkt method it
// should write empty packet to Empty channel and next read full packet from
// Full channel.
type PktReadQueue PktQueue

// Empty returns a channel that can be used to pass empty packets to q.
func (q *PktReadQueue) Empty() chan<- *ArrayPkt {
	return q.empty
}

// Full returns a channel that can be used to obtain full packets from q.
func (q *PktReadQueue) Full() <-chan *ArrayPkt {
	return q.full
}

// ReplacePkt pass empty pkt to q and obtain full packet from q.
func (q *PktReadQueue) ReplacePkt(pkt *ArrayPkt) (*ArrayPkt, error) {
	q.empty <- pkt
	p, ok := <-q.full
	if !ok {
		return pkt, io.EOF
	}
	return p, nil
}

// PacketWriteQueue represenst write part of PktQueue and implements PktReplacer
// interface. If writer uses raw channels insteed of ReplacePkt method it
// should read empty packet from Empty channel and next write full packet to
// Full channel.
type PktWriteQueue PktQueue

// Empty returns a channel that can be used to obtain empty packets from q.
func (q *PktWriteQueue) Empty() <-chan *ArrayPkt {
	return q.empty
}

// Full returns a channel that can be used to pass full packets to q.
func (q *PktWriteQueue) Full() chan<- *ArrayPkt {
	return q.full
}

// Close closes write part of queue. After close on write part, ReplacePkt
// method on read part returns io.EOF error if there is no more packets to read.
func (q *PktWriteQueue) Close() {
	close(q.full)
}

// ReplacePkt obtain empty packet from q and pass pkt to q.
func (q *PktWriteQueue) ReplacePkt(pkt *ArrayPkt) (*ArrayPkt, error) {
	p := <-q.empty
	q.full <- pkt
	return p, nil
}
