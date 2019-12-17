package psi

type StreamType byte

const (
	ZeroStreamType StreamType = 0x00
	MPEG1Video     StreamType = 0x01
	MPEG2Video     StreamType = 0x02
	MPEG1Audio     StreamType = 0x03
	MPEG2Audio     StreamType = 0x04
	PrivSect       StreamType = 0x05
	PrivPES        StreamType = 0x06
	MHEG           StreamType = 0x07
	DSMCC          StreamType = 0x08
	H222_1         StreamType = 0x09
	DSMCC_A        StreamType = 0x0A
	DSMCC_B        StreamType = 0x0B
	DSMCC_C        StreamType = 0x0C
	DSMCC_D        StreamType = 0x0D
	MPEG2Aux       StreamType = 0x0E
	AAC            StreamType = 0x0F
	MPEG4Video     StreamType = 0x10
	MPEG4Audio     StreamType = 0x11
	DSMCC_SDP      StreamType = 0x12
	SPSPES         StreamType = 0x13
	SPSSect        StreamType = 0x14
	MetaDataPES    StreamType = 0x15
	MetaDataSect   StreamType = 0x16
	MetaDataDC     StreamType = 0x17
	MetaDataOC     StreamType = 0x18
	MetaDataDL     StreamType = 0x19
	MPEG2IPMP      StreamType = 0x1A
	H264Video      StreamType = 0x1B
	MPEG4RawAudio  StreamType = 0x1C
	MPEG4Text      StreamType = 0x1D
	MPEG4Aux       StreamType = 0x1E
	MPEG4SVC       StreamType = 0x1E
	MPEG4MVC       StreamType = 0x20
	JPEG2000Video  StreamType = 0x21
	_              StreamType = 0x22
	_              StreamType = 0x23
	H265Video      StreamType = 0x24
)

var streamTypes = []string{
	"ZeroStreamType",
	"MPEG1Video",
	"MPEG2Video",
	"MPEG1Audio",
	"MPEG2Audio",
	"PrivSect",
	"PrivPES",
	"MHEG",
	"DSMCC",
	"H222_1",
	"DSMCCA",
	"DSMCCB",
	"DSMCCC",
	"DSMCCD",
	"MPEG2Aux",
	"AAC",
	"MPEG4Video",
	"MPEG4Audio",
	"DSMCC_SDP",
	"SPSPES",
	"SPSSect",
	"MetaDataPES",
	"MetaDataSect",
	"MetaDataDC",
	"MetaDataOC",
	"MetaDataDL",
	"MPEG2IPMP",
	"H264Video",
	"MPEG4RawAudio",
	"MPEG4Text",
	"MPEG4Aux",
	"MPEG4SVC",
	"MPEG4MVC",
	"JPEG2000Video",
	"reserved",
	"reserved",
	"H265Video",
}

func (t StreamType) String() string {
	if int(t) >= len(streamTypes) {
		return "reserved"
	}
	return streamTypes[t]
}

func ParseStreamType(s string) StreamType {
	for i, t := range streamTypes {
		if t == s {
			return StreamType(i)
		}
	}
	return ZeroStreamType
}
