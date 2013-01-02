package dvb

import (
	"errors"
)

var (
	// ErrOverflow means that some buffer has been overflowed and some data
	// has been lost.
	ErrOverflow = errors.New("buffering overflow")
)
