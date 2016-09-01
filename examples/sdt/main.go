package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ziutek/dvb"
	"github.com/ziutek/dvb/examples/internal"
	"github.com/ziutek/dvb/ts"
	"github.com/ziutek/dvb/ts/psi"
)

func checkErr(err error) {
	if err != nil {
		if _, ok := err.(dvb.TemporaryError); ok {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s ADDR:PORT|MADDR:PORT[@IFNAME]|FILE\n",
			filepath.Base(os.Args[0]),
		)
		os.Exit(1)
	}
	src := os.Args[1]

	var rp ts.PktReplacer

	if strings.IndexByte(src, ':') == -1 {
		f, err := os.Open(src)
		checkErr(err)
		rp = ts.NewPktStreamReader(f)
	} else {
		r, err := internal.ListenMulticastUDP(src, 17)
		if err == internal.ErrNotMulticast {
			r, err = internal.ListenUDP(src, 17)
		}
		checkErr(err)
		rp = ts.PktReaderAsReplacer{r}
	}
	d := psi.NewSectionDecoder(rp, true)
	fmt.Println("SID Provider Name Type Status Scrambled EIT(PresentFollowing/Schedule)")
	var sdt psi.SDT
	for {
		checkErr(sdt.Update(d, true, true))
		sl := sdt.ServiceInfo()
		for !sl.IsEmpty() {
			var si psi.ServiceInfo
			si, sl = sl.Pop()
			if si == nil {
				os.Stderr.WriteString("Error: demaged service list\n")
				break
			}
			sid := si.ServiceId()
			status := si.Status()
			scrambled := si.Scrambled()
			eitPF := si.EITPresentFollowing()
			eitSched := si.EITSchedule()
			var (
				name     string
				provider string
				typ      psi.ServiceType
			)
			dl := si.Descriptors()
			for len(dl) != 0 {
				var d psi.Descriptor
				d, dl = dl.Pop()
				if d == nil {
					os.Stderr.WriteString("Error: demaged descriptor list\n")
					break
				}
				if d.Tag() == psi.ServiceTag {
					sd, ok := psi.ParseServiceDescriptor(d)
					if !ok {
						os.Stderr.WriteString("Error: bad service descriptor\n")
						break
					}
					typ = sd.Type
					name = psi.DecodeText(sd.ServiceName)
					provider = psi.DecodeText(sd.ProviderName)
					break
				}
			}
			fmt.Printf(
				"%d \"%s\" \"%s\" \"%v\" %v %t %t/%t\n",
				sid, provider, name, typ, status, scrambled, eitPF, eitSched,
			)
		}
	}
}
