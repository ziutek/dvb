package internal

import (
	"errors"
	"net"
	"strings"

	"github.com/ziutek/dvb/ts"
)

func ListenUDP(laddr string, pids ...int16) (ts.PktReader, error) {
	la, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUDP("udp", la)
	if err != nil {
		return nil, err
	}
	err = c.SetReadBuffer(2 * 1024 * 1024)
	if err != nil {
		return nil, err
	}
	return &PidFilter{
		r:    ts.NewPktPktReader(c, make([]byte, 7*ts.PktLen)),
		pids: pids,
	}, nil
}

var ErrNotMulticast = errors.New("not a multicast address")

func ListenMulticastUDP(group string, pids ...int16) (ts.PktReader, error) {
	var interf string
	if n := strings.IndexByte(group, '@'); n >= 0 {
		interf = group[n+1:]
		group = group[:n]
	}
	gaddr, err := net.ResolveUDPAddr("udp", group)
	if err != nil {
		return nil, err
	}
	if !gaddr.IP.IsMulticast() {
		return nil, ErrNotMulticast
	}
	ifi, err := net.InterfaceByName(interf)
	if err != nil {
		return nil, err
	}
	c, err := net.ListenMulticastUDP("udp", ifi, gaddr)
	if err != nil {
		return nil, err
	}
	err = c.SetReadBuffer(2 * 1024 * 1024)
	if err != nil {
		return nil, err
	}
	return &PidFilter{
		r:    ts.NewPktPktReader(c, make([]byte, 7*ts.PktLen)),
		pids: pids,
	}, nil
}
