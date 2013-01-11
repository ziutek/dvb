package ts

type PktWriter interface {
	WritePkt(Pkt) error
}
