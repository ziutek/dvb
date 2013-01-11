package ts

import (
	"io"
)

type QPkt struct {
	Pkt
	Error error
}

// PktQueue represents queue of MPEG-TS packets. PktQueue is primarily intended
// for use by two gorutines (one writer and one reader), but can be used
// concurently by multiple gorutines if necessary (in this case errors
// propagation model should be well thought).
type PktQueue chan QPkt

// NewPktQueue creates new queue with internall buffer of size length.
func NewPktQueue(length int) PktQueue {
	return make(PktQueue, length)
}

func (q PktQueue) ReadPkt(pkt Pkt) error {
	qp, ok := <-q
	if !ok {
		return io.EOF
	}
	if qp.Error != nil {
		return qp.Error
	}
	copy(pkt, qp.Pkt)
	return nil
}

func (q PktQueue) WritePkt(pkt Pkt) error {
	q <- QPkt{Pkt: pkt}
	return nil
}

func (q PktQueue) WriteError(e error) {
	q <- QPkt{Error: e}
}
