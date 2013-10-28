package frontend

// Satellite equipment control (SEC) voltage
type Voltage uint32

const (
	Voltage13 Voltage = iota
	Voltage18
	VoltageOff
)

// Satellite equipment control (SEC) tone
type Tone uint32

const (
	ToneOn Tone = iota
	ToneOff
)

// SecParam calculates intermediate frequency, tone and voltage for given
// absolute frequency (HZ) and polarization ('h' or 'v').
func SecParam(freq uint64, polarization byte) (f uint32, t Tone, v Voltage) {
	switch polarization {
	case 'h':
		v = Voltage18
	case 'v':
		v = Voltage13
	default:
		panic("unknown polarization")
	}
	if freq < 11700e6 {
		f = uint32(freq-9750e6) / 1000
		t = ToneOff
	} else {
		f = uint32(freq-10600e6) / 1000
		t = ToneOn
	}
	return
}

type Error struct {
	Op   string
	What string
	Err  error
}

func (e Error) Error() string {
	return e.Op + " " + e.What + ": " + e.Err.Error()
}
