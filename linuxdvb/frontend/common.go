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
