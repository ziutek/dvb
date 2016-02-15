package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/ts"
)

func die(i interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", i)
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		die(err)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: stream_file FILE BITRATE DESTINATION\n")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	checkErr(err)
	r := ts.NewPktStreamReader(bufio.NewReader(file))

	bitrate, err := strconv.ParseUint(os.Args[2], 0, 64)
	checkErr(err)

	addr, err := net.ResolveUDPAddr("udp", os.Args[3])
	checkErr(err)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	checkErr(err)

	var (
		buf [7 * ts.PktLen]byte
		n   int
	)
	bps := time.Duration(bitrate)
	pktperiod := (time.Second*ts.PktLen*8 + bps/2) / bps
	sendt := time.Now()
	for {
		pkt := ts.AsPkt(buf[n:])
		err := r.ReadPkt(pkt)
		if err == nil {
			n += ts.PktLen
			if n == len(buf) {
				n = 0
				sendt = sendt.Add(pktperiod)
				time.Sleep(sendt.Sub(time.Now()))
				conn.WriteToUDP(buf[:], addr)
				//fmt.Printf("PID: %d\n", pkt.Pid())
			}
		} else if err == io.EOF {
			_, err := file.Seek(0, os.SEEK_SET)
			checkErr(err)
			os.Stdout.WriteString("seek\n")
		} else if _, ok := err.(dvb.TemporaryError); ok {
			fmt.Fprintf(os.Stderr, "temporary error:%v\n", err)
		} else {
			die(err)
		}
	}
}
