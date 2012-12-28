package dvb

import (
	"errors"
)

var (
	// OverflowError means that some buffer has been overflowed and some data
	// has been lost.
	OverflowError = errors.New("buffering overflow")
)
