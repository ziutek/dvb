package frontend

type Inversion uint32

const (
	InversionOff Inversion = iota
	InversionOn
	InversionAuto
)

type CodeRate uint32

const (
	FECNone CodeRate = iota
	FEC12
	FEC23
	FEC34
	FEC45
	FEC56
	FEC67
	FEC78
	FEC89
	FECAuto
	FEC35
	FEC910
)

type Modulation uint32

const (
	QPSK Modulation = iota
	QAM16
	QAM32
	QAM64
	QAM128
	QAM256
	QAMAuto
	VSB8
	VSB16
	PSK8
	APSK16
	APSK32
	DQPSK
)

type TxMode uint32

const (
	TxMode2k TxMode = iota
	TxMode8k
	TxModeAuto
	TxMode4k
	TxMode1k
	TxMode16k
	TxMode32k
)

type GuardInt uint32

const (
	GuardInt32 GuardInt = iota // 1/32
	GuardInt16                 // 1/16
	GuardInt8                  // 1/8
	GuardInt4                  // 1/4
	GuardIntAuto
	GuardInt128  // 1/128
	GuardIntN128 // 19/128
	GuardIntN256 // 19/128
)

type Hierarchy uint32

const (
	HierarchyNone Hierarchy = iota
	Hierarchy1
	Hierarchy2
	Hierarchy4
	HierarchyAuto
)

/*type Rolloff uint32

const (
	Rolloff35 Rolloff = iota // Implied value in DVB-S, default for DVB-S2
	Rolloff20
	Rolloff25
	RolloffAuto
)*/
