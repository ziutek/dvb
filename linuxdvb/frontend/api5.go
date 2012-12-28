package frontend

// APIv5 commands
const (
	dtvUndefined = iota
	dtvTune
	dtvClear
	dtvFreqency
	dtvModulation
	dtvBandwidthHz
	dtvInversion
	dtvDiseqcMaster
	dtvSymbolRate
	dtvInnerFec
	dtvVoltage
	dtvTone
	dtvPilot
	dtvRollOff
	dtvDiseqcSlaveReply
	dtvFeCapabilityCount
	dtvFeCapability
	dtvDeliverySystem
)

// Frontend delivery sytem
const (
	sysUndefined = iota
	sysDVBCAnnexAC
	sysDVBCAnnexB
	sysDVBT
	sysDSS
	sysDVBS
	sysDVBS2
	sysDVBH
	sysISDBT
	sysISDBS
	sysISDBC
	sysATSC
	sysATSCMH
	sysDMBTH
	sysCMMB
	sysDAB
	sysDVBT2
	sysTURBO
)
