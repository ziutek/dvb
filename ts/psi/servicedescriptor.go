package psi

type ServiceType byte

const (
	ReservedServiceType ServiceType = iota
	DigitalTelevisionService
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
	"reserved for future use",
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
	if t > DataBroadcastService {
		return "unknown service type"
	}
	return stn[t]
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
