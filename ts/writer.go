package ts

import (
	"io"
)

// PktWriter is an interface that wraps the WritePkt method.
type PktWriter interface {
	// WritePkt writes one MPEG-TS packet.
	WritePkt(Pkt) error
}

// PktStreamReader wraps any io.Writer interface and returns PktWriter and
// PktReplacer implementation for write MPEG-TS packets as stream of bytes.
type PktStreamWriter struct {
	W io.Writer
}

// WritePkt wraps s.W.Write method
func (s PktStreamWriter) WritePkt(pkt Pkt) error {
	_, err := s.W.Write(pkt.Bytes())
	return err
}

// ReplacePkt works like WritePkt but implements PktReplacer interface.
func (s PktStreamWriter) ReplacePkt(pkt *ArrayPkt) (*ArrayPkt, error) {
	_, err := s.W.Write(pkt.Bytes())
	return pkt, err
}
