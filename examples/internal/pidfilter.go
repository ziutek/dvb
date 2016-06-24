package internal

import (
	"github.com/ziutek/dvb/ts"
)

type PidFilter struct {
	r    ts.PktReader
	pids []int16
}

func (f *PidFilter) ReadPkt(pkt ts.Pkt) error {
	for {
		if err := f.r.ReadPkt(pkt); err != nil {
			return err
		}
		pid := pkt.Pid()
		// TODO: sort f.pids to use more effecitve search method.
		for _, p := range f.pids {
			if p == 8192 || p == pid {
				return nil
			}
		}
	}
}
