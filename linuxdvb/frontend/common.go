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
	TonefOff
)

// SecParam calculates intermediate frequency, tone and voltage for given
// absolute frequency (HZ) and polarization ('h' or 'v').
func SecParam(freq uint64, polarization byte) (f uint64, t Tone, v Voltage) {
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
		t = TonefOff
	} else {
		f = freq - 10600000
		t = TonefOn
	}
	return
}
