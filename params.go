package dvb

// Device delivery sytem
type DeliverySystem uint32

const (
	SysUndefined DeliverySystem = iota
	SysDVBCAnnexA
	SysDVBCAnnexB
	SysDVBT
	SysDSS
	SysDVBS
	SysDVBS2
	SysDVBH
	SysISDBT
	SysISDBS
	SysISDBC
	SysATSC
	SysATSCMH
	SysDMBTH
	SysCMMB
	SysDAB
	SysDVBT2
	SysTURBO
	SysDVBCAnnexC
)

var dsn = []string{
	"Undefined",
	"DVB-C Annex AC",
	"DVB-C Annex B",
	"DVB-T",
	"DSS",
	"DVB-S",
	"DVB-S2",
	"DVB-H",
	"ISDB-T",
	"ISDB-S",
	"ISDB-C",
	"ATSC",
	"ATSC-MH",
	"DMBT-H",
	"CMMB",
	"DAB",
	"DVB-T2",
	"TURBO",
}

func (ds DeliverySystem) String() string {
	if ds > DeliverySystem(len(dsn)) {
		return "unknown"
	}
	return dsn[ds]
}

type Inversion uint32

const (
	InversionOff Inversion = iota
	InversionOn
	InversionAuto
)

var inversionNames = []string{
	"off",
	"on",
	"auto",
}

func (i Inversion) String() string {
	if i > InversionAuto {
		return "unknown"
	}
	return inversionNames[i]
}

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

var codeRateNames = []string{
	"none",
	"1/2",
	"2/3",
	"3/4",
	"4/5",
	"5/6",
	"6/7",
	"7/8",
	"8/9",
	"auto",
	"3/5",
	"9/10",
}

func (cr CodeRate) String() string {
	if cr > FEC910 {
		return "unknown"
	}
	return codeRateNames[cr]
}

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

var modulationNames = []string{
	"QPSK",
	"QAM16",
	"QAM32",
	"QAM64",
	"QAM128",
	"QAM256",
	"QAMAuto",
	"VSB8",
	"VSB16",
	"PSK8",
	"APSK16",
	"APSK32",
	"DQPSK",
}

func (m Modulation) String() string {
	if m > DQPSK {
		return "unknown"
	}
	return modulationNames[m]
}

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

var txModeNames = []string{
	"2k",
	"8k",
	"auto",
	"4k",
	"1k",
	"16k",
	"32k",
}

func (tm TxMode) String() string {
	if tm > TxMode32k {
		return "unknown"
	}
	return txModeNames[tm]
}

type Guard uint32

const (
	Guard32 Guard = iota // 1/32
	Guard16              // 1/16
	Guard8               // 1/8
	Guard4               // 1/4
	GuardAuto
	Guard128  // 1/128
	GuardN128 // 19/128
	GuardN256 // 19/256
)

var guardNames = []string{
	"1/32",
	"1/16",
	"1/8",
	"1/4",
	"auto",
	"1/128",
	"19/128",
	"19/256",
}

func (gi Guard) String() string {
	if gi > GuardN256 {
		return "unknown"
	}
	return guardNames[gi]
}

type Hierarchy uint32

const (
	HierarchyNone Hierarchy = iota
	Hierarchy1
	Hierarchy2
	Hierarchy4
	HierarchyAuto
)

var hierarchyNames = []string{
	"none",
	"uniform",
	"HP/LP=2",
	"HP/LP=4",
	"auto",
}

func (h Hierarchy) String() string {
	if h > HierarchyAuto {
		return "unknown"
	}
	return hierarchyNames[h]
}

// DVB-S2 pilot
type Pilot uint32

const (
	PilotOn Pilot = iota
	PilotOff
	PilotAuto
)

// DVB-S2 rolloff
type Rolloff uint32

const (
	Rolloff35 Rolloff = iota // Implied value in DVB-S, default for DVB-S2
	Rolloff20
	Rolloff25
	RolloffAuto
)
