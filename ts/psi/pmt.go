package psi

type PMT struct {
	s Section
}

func NewPMT() PMT {
	return PMT{NewSection(ISOSectionMaxLen)}
}

func (p PMT) Version() byte {
	return p.s.Version()
}

func (p PMT) Current() bool {
	return p.s.Current()
}

func (p PMT) ProgId() uint16 {
	return p.s.TableIdExt()
}

func (p PMT) PCRPid() uint16 {
	return decodeU16(p.s.Data()[0:2]) & 0x1fff
}

func (p PMT) progInfoLen() uint16 {
	return decodeU16(p.s.Data()[2:4]) & 0x0fff
}

type StreamInfo []byte

func (i StreamInfo) Type() StreamType {
	return StreamType(i[0])
}

type StreamType byte

const (
	Reserved StreamType = iota
	MPEG1Video
	MPEG2Video
	MPEG1Audio
	MPEG2Audio
	PrivateSect
	PrivatePES
	MHEG
	DSMCC
	H222_1
	DSMCC_A
	DSMCC_B
	DSMCC_C
	DSMCC_D
	MPEG2Aux
	AAC
	MPEG4Video
	MPEG4Audio
	DSMCC_SDP
	SPSPES
	SPSSect
	MetaDataPES
	MetaDataSect
	MetaDataDC
	MetaDataOC
	MetaDataDL
	MPEG2IPMP
	H264Video
)

var streamTypes = []string{
	"Reserved",
	"MPEG1Video",
	"MPEG2Video",
	"MPEG1Audio",
	"MPEG2Audio",
	"PrivateSect",
	"PrivatePES",
	"MHEG",
	"DSMCC",
	"H222_1",
	"DSMCC_A",
	"DSMCC_B",
	"DSMCC_C",
	"DSMCC_D",
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
}

func (t StreamType) String() string {
	if int(t) >= len(streamTypes) {
		return "unknown"
	}
	return streamTypes[t]
}
