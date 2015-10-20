package main

import (
	"bufio"
	"net"
	"os"

	"github.com/ziutek/dvb/ts"
)

type outputFile struct {
	w *bufio.Writer
}

func newOutputFile(path string) ts.PktWriter {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	checkErr(err)
	return outputFile{bufio.NewWriter(f)}
}

func (o outputFile) WritePkt(pkt ts.Pkt) error {
	_, err := o.w.Write(pkt.Bytes())
	return err
}

type outputUDP struct {
	conn *net.UDPConn
	addr *net.UDPAddr
	buf  []byte
}

func newOutputUDP(src, dst string) ts.PktWriter {
	var err error
	saddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	if src != "" {
		saddr, err = net.ResolveUDPAddr("udp", src)
		checkErr(err)
	}
	conn, err := net.ListenUDP("udp", saddr)
	checkErr(err)
	daddr, err := net.ResolveUDPAddr("udp", dst)
	checkErr(err)

	return &outputUDP{
		conn: conn,
		addr: daddr,
		buf:  make([]byte, 0, 7*ts.PktLen),
	}
}

func (o *outputUDP) WritePkt(pkt ts.Pkt) error {
	o.buf = append(o.buf, pkt.Bytes()...)
	if len(o.buf) < 7*ts.PktLen {
		return nil
	}
	_, err := o.conn.WriteToUDP(o.buf, o.addr)
	o.buf = o.buf[:0]
	return err
}

type outputDiscard struct{}

func (_ outputDiscard) WritePkt(_ ts.Pkt) error {
	return nil
}
