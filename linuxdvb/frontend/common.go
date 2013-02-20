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
func SecParam(freq uint32, polarization byte) (f uint32, t Tone, v Voltage) {
	switch polarization {
	case 'h':
		v = Voltage18
	case 'v':
		v = Voltage13
	default:
		panic("unknown polarization")
	}
	if freq < 11700000 {
		f = freq - 9750000
		t = ToneOff
	} else {
		f = freq - 10600000
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
