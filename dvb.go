package dvb

type TemporaryError string

func (e TemporaryError) Error() string {
	return string(e)
}

var (
	// ErrOverflow means that some buffer has been overflowed and some data
	// has been lost.
	ErrOverflow = TemporaryError("dvb: buffer overflow")
)
