package psi

func decodeU32(b []byte) uint32 {
	if len(b) != 4 {
		panic("decodeU32 with len(b) != 4")
	}
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 |
		uint32(b[3])
}

func decodeU16(b []byte) uint16 {
	if len(b) != 2 {
		panic("decodeU16 with len(b) != 2")
	}
	return uint16(b[0])<<8 | uint16(b[1])
}

func encodeU32(b []byte, v uint32) {
	if len(b) != 4 {
		panic("encodeU32 with len(b) != 4")
	}
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func encodeU16(b []byte, v uint16) {
	if len(b) != 2 {
		panic("encodeU16 with len(b) != 2")
	}
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

var crcTable [256]uint32

func mpegCRC32(buf []byte) uint32 {
	crc := uint32(0xffffffff)
	for _, b := range buf {
		crc = crcTable[byte(crc>>24)^b] ^ (crc << 8)
	}
	return crc
}

func init() {
	poly := uint32(0x04c11db7)
	for i := 0; i < 256; i++ {
		crc := uint32(i) << 24
		for j := 0; j < 8; j++ {
			if crc&0x80000000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		crcTable[i] = crc
	}
}
