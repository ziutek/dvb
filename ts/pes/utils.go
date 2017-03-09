package pes

func decodeU16(b []byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

func decodeU40(b []byte) uint64 {
	return uint64(b[0])<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 |
		uint64(b[3])<<8 | uint64(b[4])
}
