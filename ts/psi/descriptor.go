package psi

type Descriptor []byte

const descrDataTooLong = "psi: data to long to alloc descriptor"

func MakeDescriptor(tag DescriptorTag, datalen int) Descriptor {
	if uint(datalen) > 255 {
		panic(descrDataTooLong)
	}
	d := make(Descriptor, 1+1+datalen)
	d.SetTag(tag)
	d[1] = byte(datalen)
	return d
}

func (d Descriptor) Tag() DescriptorTag {
	return DescriptorTag(d[0])
}

func (d Descriptor) SetTag(tag DescriptorTag) {
	d[0] = byte(tag)
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
	d = Descriptor(dl[:l])
	rdl = dl[l:]
	return
}

// Alloc allocates descriptor for datalen bytes of data + 2 bytes for tag and
// data length.
func (dl *DescriptorList) Alloc(datalen int) Descriptor {
	if uint(datalen) > 255 {
		panic("descrDataTooLong")
	}
	m := len(*dl)
	n := m + datalen + 2
	if n <= cap(*dl) {
		*dl = (*dl)[:n]
	} else {
		ndl := make(DescriptorList, n, n+m)
		copy(ndl, *dl)
		*dl = ndl
	}
	d := Descriptor((*dl)[m:n])
	d[1] = byte(datalen)
	return d
}

type DescriptorTag byte

const (
	VideoStreamTag                DescriptorTag = 0x02
	AudioStreamTag                DescriptorTag = 0x03 //PMT
	HierarchyTag                  DescriptorTag = 0x04
	RegistrationTag               DescriptorTag = 0x05 //PMT
	DataStreamAlignmentTag        DescriptorTag = 0x06
	TargetBackgroundGridTag       DescriptorTag = 0x07
	VideoWindowTag                DescriptorTag = 0x08
	CATag                         DescriptorTag = 0x09
	ISO639LangTag                 DescriptorTag = 0x0A //PMT
	SystemClockTag                DescriptorTag = 0x0B
	MultiplexBufferUtilizationTag DescriptorTag = 0x0C
	CopyrightTag                  DescriptorTag = 0x0D
	MaximumBitrateTag             DescriptorTag = 0x0E //PMT
	PrivateDataIndicatorTag       DescriptorTag = 0x0F
	SmoothingBufferTag            DescriptorTag = 0x10
	STDTag                        DescriptorTag = 0x11
	IBPTag                        DescriptorTag = 0x12

	MPEG4VideoTag      DescriptorTag = 0x1b
	MPEG4AudioTag      DescriptorTag = 0x1c
	IODTag             DescriptorTag = 0x1d
	SLTag              DescriptorTag = 0x1f
	FMCTag             DescriptorTag = 0x20
	ExternalESIDTag    DescriptorTag = 0x21
	MuxCodeTag         DescriptorTag = 0x22
	FmxBufferSizeTag   DescriptorTag = 0x23
	MultiplexBufferTag DescriptorTag = 0x24

	NetworkNameTag               DescriptorTag = 0x40 //NIT
	ServiceListTag               DescriptorTag = 0x41 //NIT
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
	StreamIdentifierTag          DescriptorTag = 0x52 //PMT
	CAIdentifierTag              DescriptorTag = 0x53
	ContentTag                   DescriptorTag = 0x54
	ParentalRatingTag            DescriptorTag = 0x55
	TeletextTag                  DescriptorTag = 0x56
	TelephoneTag                 DescriptorTag = 0x57
	LocalTimeOffsetTag           DescriptorTag = 0x58
	SubtitlingTag                DescriptorTag = 0x59 //PMT
	TerrestrialDeliverySystemTag DescriptorTag = 0x5A //NIT
	MultilingualNetworkNameTag   DescriptorTag = 0x5B
	MultilingualBouquetNameTag   DescriptorTag = 0x5C
	MultilingualServiceNameTag   DescriptorTag = 0x5D
	MultilingualComponentTag     DescriptorTag = 0x5E
	PrivateDataSpecifierTag      DescriptorTag = 0x5F //NIT
	ServiceMoveTag               DescriptorTag = 0x60
	ShortSmoothingBufferTag      DescriptorTag = 0x61
	FrequencyListTag             DescriptorTag = 0x62 //NIT
	PartialTransportStreamTag    DescriptorTag = 0x63
	DataBroadcastTag             DescriptorTag = 0x64
	CASystemTag                  DescriptorTag = 0x65
	DataBroadcastIdTag           DescriptorTag = 0x66
	TransportStreamTag           DescriptorTag = 0x67
	DSNGTag                      DescriptorTag = 0x68
	PDCTag                       DescriptorTag = 0x69
	AC3Tag                       DescriptorTag = 0x6a //PMT
	AncillaryDataTag             DescriptorTag = 0x6b
	CellListTag                  DescriptorTag = 0x6c
	CellFrequencyLinkTag         DescriptorTag = 0x6d
	AnnouncementSupportTag       DescriptorTag = 0x6e
	ApplicationSignallingTag     DescriptorTag = 0x6f //PMT

	EnhancedAC3Tag DescriptorTag = 0x7a //PMT
	DTSTag         DescriptorTag = 0x7b
	AACTag         DescriptorTag = 0x7c

	LogicalChannelTag DescriptorTag = 0x83 //NIT
)

var dtagstr = [...]string{
	VideoStreamTag:                "VideoStream",
	AudioStreamTag:                "AudioStream",
	HierarchyTag:                  "Hierarchy",
	RegistrationTag:               "Registration",
	DataStreamAlignmentTag:        "DataStreamAlignment",
	TargetBackgroundGridTag:       "TargetBackgroundGrid",
	VideoWindowTag:                "VideoWindow",
	CATag:                         "CA",
	ISO639LangTag:                 "ISO639Lang",
	SystemClockTag:                "SystemClock",
	MultiplexBufferUtilizationTag: "MultiplexBufferUtilization",
	CopyrightTag:                  "Copyright",
	MaximumBitrateTag:             "MaximumBitrate",
	PrivateDataIndicatorTag:       "PrivateDataIndicator",
	SmoothingBufferTag:            "SmoothingBuffer",
	STDTag:                        "STD",
	IBPTag:                        "IBP",

	MPEG4VideoTag:      "MPEG4Video",
	MPEG4AudioTag:      "MPEG4Audio",
	IODTag:             "IOD",
	SLTag:              "SL",
	FMCTag:             "FMC",
	ExternalESIDTag:    "ExternalESID",
	MuxCodeTag:         "MuxCode",
	FmxBufferSizeTag:   "FmxBufferSize",
	MultiplexBufferTag: "MultiplexBuffer",

	NetworkNameTag:               "NetworkName",
	ServiceListTag:               "ServiceList",
	StuffingTag:                  "Stuffing",
	SatelliteDeliverySystemTag:   "SatelliteDeliverySystem",
	CableDeliverySystemTag:       "CableDeliverySystem",
	BouquetNameTag:               "BouquetName",
	ServiceTag:                   "Service",
	CountryAvailabilityTag:       "CountryAvailability",
	LinkageTag:                   "Linkage",
	NVODReferenceTag:             "NVODReference",
	TimeShiftedServiceTag:        "TimeShiftedService",
	ShortEventTag:                "ShortEvent",
	ExtendedEventTag:             "ExtendedEvent",
	TimeShiftedEventTag:          "TimeShiftedEvent",
	ComponentTag:                 "Component",
	MosaicTag:                    "Mosaic",
	StreamIdentifierTag:          "StreamIdentifier",
	CAIdentifierTag:              "CAIdentifier",
	ContentTag:                   "Content",
	ParentalRatingTag:            "ParentalRating",
	TeletextTag:                  "Teletext",
	TelephoneTag:                 "Telephone",
	LocalTimeOffsetTag:           "LocalTimeOffset",
	SubtitlingTag:                "Subtitling",
	TerrestrialDeliverySystemTag: "TerrestrialDeliverySystem",
	MultilingualNetworkNameTag:   "MultilingualNetworkName",
	MultilingualBouquetNameTag:   "MultilingualBouquetName",
	MultilingualServiceNameTag:   "MultilingualServiceName",
	MultilingualComponentTag:     "MultilingualComponent",
	PrivateDataSpecifierTag:      "PrivateDataSpecifier",
	ServiceMoveTag:               "ServiceMove",
	ShortSmoothingBufferTag:      "ShortSmoothingBuffer",
	FrequencyListTag:             "FrequencyList",
	PartialTransportStreamTag:    "PartialTransportStream",
	DataBroadcastTag:             "DataBroadcast",
	CASystemTag:                  "CASystem",
	DataBroadcastIdTag:           "DataBroadcastId",
	TransportStreamTag:           "TransportStream",
	DSNGTag:                      "DSNG",
	PDCTag:                       "PDC",
	AC3Tag:                       "AC3",
	AncillaryDataTag:             "AncillaryData",
	CellListTag:                  "CellList",
	CellFrequencyLinkTag:         "CellFrequencyLink",
	AnnouncementSupportTag:       "AnnouncementSupport",
	ApplicationSignallingTag:     "ApplicationSignalling",

	EnhancedAC3Tag: "EnhancedAC3",
	DTSTag:         "DTS",
	AACTag:         "AAC",

	LogicalChannelTag: "LogicalChannel",
}

func (tag DescriptorTag) String() string {
	if int(tag) >= len(dtagstr) {
		return ""
	}
	return dtagstr[tag]
}
