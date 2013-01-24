package psi

type ServiceType byte

const (
	DigitalTelevisionService ServiceType = iota + 1
	DigitalRadioSoundService
	TeletextService
	NVODReferenceService
	NVODTimeShiftedService
	MosaicService
	PALCodedSignal
	SCAMCodedSignal
	D2MAC
	FMRadio
	NTSCCodedSignal
	DataBroadcastService
)

var stn = []string{
	"digital television service",
	"digital radio sound service",
	"Teletext service",
	"NVOD reference service",
	"NVOD time-shifted service",
	"mosaic service",
	"PAL coded signal",
	"SECAM coded signal",
	"D/D2-MAC",
	"FM Radio",
	"NTSC coded signal",
	"data broadcast service",
}

func (t ServiceType) String() string {
	if t == 0 || t > DataBroadcastService && t <= 0x7F {
		return "reserved"
	}
	if t > 0x7F && t <= 0xfe {
		return "user defined"
	}
	return stn[t-1]
}

type ServiceDescriptor struct {
	Type         ServiceType
	ProviderName []byte
	ServiceName  []byte
}

func ParseServiceDescriptor(d Descriptor) (sd ServiceDescriptor, ok bool) {
	data := d.Data()
	if len(data) < 2 {
		return
	}
	sd.Type = ServiceType(data[0])
	providerNameLen := data[1]
	data = data[2:]

	if len(data) < int(providerNameLen+1) {
		return
	}
	sd.ProviderName = data[:providerNameLen]
	serviceNameLen := data[providerNameLen]
	data = data[providerNameLen+1:]

	if len(data) < int(serviceNameLen) {
		return
	}
	sd.ServiceName = data[:serviceNameLen]

	ok = true
	return
}
