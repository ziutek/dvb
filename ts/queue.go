package ts

import (
	"io"
)

// PktQueue represents queue of TS packets.
type PktQueue struct {
	empty, filled chan *ArrayPkt
}

// NewPktQueue creates new queue with internall buffer of size length packets.
func NewPktQueue(length int) *PktQueue {
	q := &PktQueue{
		empty:  make(chan *ArrayPkt, length),
		filled: make(chan *ArrayPkt, length),
	}
	for i := 0; i < length; i++ {
		q.empty <- new(ArrayPkt)
	}
	return q
}

// Cap returns capacity of q.
func (q *PktQueue) Cap() int {
	return cap(q.filled)
}

// Len returns number of packet queued in q.
func (q *PktQueue) Len() int {
	return len(q.filled)
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
// should first read filled packet from the Filled channel and next write
// empty packet to the Empty channel.
type PktReadQueue PktQueue

// Empty returns a channel that can be used to pass empty packets to q.
func (q *PktReadQueue) Empty() chan<- *ArrayPkt {
	return q.empty
}

// Filled returns a channel that can be used to obtain filled packets from q.
func (q *PktReadQueue) Filled() <-chan *ArrayPkt {
	return q.filled
}

// ReplacePkt pass empty pkt to q and obtain filled packet from q. It returns
// io.EOF error when queue was closed.
func (q *PktReadQueue) ReplacePkt(pkt *ArrayPkt) (*ArrayPkt, error) {
	p, ok := <-q.filled
	if !ok {
		return pkt, io.EOF
	}
	q.empty <- pkt
	return p, nil
}

// Cap returns capacity of q.
func (q *PktReadQueue) Cap() int {
	return cap(q.filled)
}

// Len returns number of packet queued in q.
func (q *PktReadQueue) Len() int {
	return len(q.filled)
}

// PacketWriteQueue represenst write part of PktQueue and implements PktReplacer
// interface. If writer uses raw channels insteed of ReplacePkt method it
// should read empty packet from Empty channel and next write filled packet
// to Filled channel.
type PktWriteQueue PktQueue

// Empty returns a channel that can be used to obtain empty packets from q.
func (q *PktWriteQueue) Empty() <-chan *ArrayPkt {
	return q.empty
}

// Filled returns a channel that can be used to pass filled packets to q.
func (q *PktWriteQueue) Filled() chan<- *ArrayPkt {
	return q.filled
}

// Close closes write part of queue. After close on write part, ReplacePkt
// method on read part returns io.EOF error if there is no more packets to read.
func (q *PktWriteQueue) Close() {
	close(q.filled)
}

// ReplacePkt obtain empty packet from q and pass pkt to q. It always return
// nil error.
func (q *PktWriteQueue) ReplacePkt(pkt *ArrayPkt) (*ArrayPkt, error) {
	p := <-q.empty
	q.filled <- pkt
	return p, nil
}

// Cap returns capacity of q.
func (q *PktWriteQueue) Cap() int {
	return cap(q.filled)
}

// Len returns number of packet queued in q.
func (q *PktWriteQueue) Len() int {
	return len(q.filled)
}
