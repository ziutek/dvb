package psi

type Descriptor []byte

func (d Descriptor) Tag() DescriptorTag {
	return DescriptorTag(d[0])
}

func (d Descriptor) Data() []byte {
	return d[2 : 2+d[1]]
}

type DescriptorList []byte

// Pop returns first descriptor in d and remaining descriptors in rdl.
// If there is no more descriptors len(rdl) == 0. If an error occurs
// d == nil.
func (dl DescriptorList) Pop() (d Descriptor, rdl DescriptorList) {
	if len(dl) < 2 {
		return
	}
	l := int(dl[1]) + 2
	if len(dl) < l {
		return
	}
	d = Descriptor(dl[2:l])
	rdl = dl[l:]
	return
}

// Append adds d to the end of the dl. It works like Go append function so need
// to be used in this way:
//     dl = dl.Append(d)
func (dl DescriptorList) Append(d Descriptor) DescriptorList {
	return append(dl, d...)
}

type DescriptorTag byte

const (
	NetworkNameTag               DescriptorTag = 0x40
	ServiceListTag               DescriptorTag = 0x41
	StuffingTag                  DescriptorTag = 0x42
	SatelliteDeliverySystemTag   DescriptorTag = 0x43
	CableDeliverySystemTag       DescriptorTag = 0x44
	BouquetNameTag               DescriptorTag = 0x47
	ServiceTag                   DescriptorTag = 0x48
	CountryAvailabilityTag       DescriptorTag = 0x49
	LinkageTag                   DescriptorTag = 0x4A
	NVODReferenceTag             DescriptorTag = 0x4B
	TimeShiftedServiceTag        DescriptorTag = 0x4C
	ShortEventTag                DescriptorTag = 0x4D
	ExtendedEventTag             DescriptorTag = 0x4E
	TimeShiftedEventTag          DescriptorTag = 0x4F
	ComponentTag                 DescriptorTag = 0x50
	MosaicTag                    DescriptorTag = 0x51
	StreamIdentifierTag          DescriptorTag = 0x52
	CAIdentifierTag              DescriptorTag = 0x53
	ContentTag                   DescriptorTag = 0x54
	ParentalRatingTag            DescriptorTag = 0x55
	TeletextTag                  DescriptorTag = 0x56
	TelephoneTag                 DescriptorTag = 0x57
	LocalTimeOffsetTag           DescriptorTag = 0x58
	SubtitlingTag                DescriptorTag = 0x59
	TerrestrialDeliverySystemTag DescriptorTag = 0x5A
	MultilingualNetworkName      DescriptorTag = 0x5B
	MultilingualBouquetName      DescriptorTag = 0x5C
	MultilingualServiceName      DescriptorTag = 0x5D
	MultilingualComponent        DescriptorTag = 0x5E
	PrivateDataSpecifier         DescriptorTag = 0x5F
	ServiceMove                  DescriptorTag = 0x60
	ShortSmoothingBuffer         DescriptorTag = 0x61
	FrequencyList                DescriptorTag = 0x62
	PartialTransportStreamTag    DescriptorTag = 0x63
	DataBroadcastTag             DescriptorTag = 0x64
	CASystemTag                  DescriptorTag = 0x65
	DataBroadcastIdTag           DescriptorTag = 0x66
)
