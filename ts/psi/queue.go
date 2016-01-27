package psi

import (
	"io"
)

// SectionQueue represents queue of Sections.
type SectionQueue struct {
	empty, filled chan Section
}

// NewSectionQueue creates new queue with internall buffer of size n sections. Every section in
// buffer has secap capacity. Sections in buffer are not inintialized.
func NewSectionQueue(n, seclen int) *SectionQueue {
	q := &SectionQueue{
		empty:  make(chan Section, n),
		filled: make(chan Section, n),
	}
	for i := 0; i < n; i++ {
		q.empty <- make(Section, seclen)
	}
	return q
}

// Cap returns capacity of q.
func (q *SectionQueue) Cap() int {
	return cap(q.filled)
}

// Len returns number of sections queued in q.
func (q *SectionQueue) Len() int {
	return len(q.filled)
}

// ReadPart returns read part of q that can be used only to read sections from
// q.
func (q *SectionQueue) ReadPart() *SectionReadQueue {
	return (*SectionReadQueue)(q)
}

// WritePart returns write part of q that can be used to write sections to q and
// to close q.
func (q *SectionQueue) WritePart() *SectionWriteQueue {
	return (*SectionWriteQueue)(q)
}

// SectionReadQueue represenst read part of SectionQueue and implements
// SectionReplacer interface. If reader uses raw channels insteed of
// ReplaceSection method it should first read filled sections from the Filled
// channel and next write empty sections to the Empty channel.
type SectionReadQueue SectionQueue

// Empty returns a channel that can be used to pass empty sections to q.
func (q *SectionReadQueue) Empty() chan<- Section {
	return q.empty
}

// Filled returns a channel that can be used to obtain filled sections from q.
func (q *SectionReadQueue) Filled() <-chan Section {
	return q.filled
}

// ReplaceSection obtains filled section from q and next pass empty section to
// q. It returns io.EOF error when queue was closed and there is no more
// sections to read.
func (q *SectionReadQueue) ReplaceSection(s Section) (Section, error) {
	fs, ok := <-q.filled
	if !ok {
		return s, io.EOF
	}
	q.empty <- s
	return fs, nil
}

// Cap returns capacity of q.
func (q *SectionReadQueue) Cap() int {
	return cap(q.filled)
}

// Len returns number of sections queued in q.
func (q *SectionReadQueue) Len() int {
	return len(q.filled)
}

// SectionWriteQueue represenst write part of SectionQueue and implements
// SectionReplacer interface. If writer uses raw channels insteed of
// ReplaceSection method it should read empty section from Empty channel and
// next write filled section to Filled channel.
type SectionWriteQueue SectionQueue

// Empty returns a channel that can be used to obtain empty sections from q.
func (q *SectionWriteQueue) Empty() <-chan Section {
	return q.empty
}

// Filled returns a channel that can be used to pass filled sections to q.
func (q *SectionWriteQueue) Filled() chan<- Section {
	return q.filled
}

// Close closes write part of queue. After close on write part, ReplaceSection
// method on read part returns io.EOF error if there is no more sections to read
// from q.
func (q *SectionWriteQueue) Close() {
	close(q.filled)
}

// ReplaceSection obtains empty section from q and next pass pkt to q. It always
// returns nil error.
func (q *SectionWriteQueue) ReplaceSection(s Section) (Section, error) {
	es := <-q.empty
	q.filled <- s
	return es, nil
}

// Cap returns capacity of q.
func (q *SectionWriteQueue) Cap() int {
	return cap(q.filled)
}

// Len returns number of sections queued in q.
func (q *SectionWriteQueue) Len() int {
	return len(q.filled)
}
