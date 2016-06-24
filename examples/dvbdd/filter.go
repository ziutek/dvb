package main

import (
	"os"

	"github.com/ziutek/dvb/linuxdvb/demux"
	"github.com/ziutek/dvb/ts"
)

func setFilter(dmxpath, dvrpath string, pids []int16) (ts.PktReader, demux.StreamFilter) {
	filterParam := demux.StreamFilterParam{
		Pid:  pids[0],
		In:   demux.InFrontend,
		Out:  demux.OutTSDemuxTap,
		Type: demux.Other,
	}
	if dvrpath != "" {
		filterParam.Out = demux.OutTSTap
	}
	filter, err := demux.Device(dmxpath).NewStreamFilter(&filterParam)
	checkErr(err)
	for _, pid := range pids[1:] {
		checkErr(filter.AddPid(pid))
	}
	if dvrpath == "" {
		checkErr(filter.SetBufferSize(1024 * 188))
		checkErr(filter.Start())
		return ts.NewPktStreamReader(filter), filter
	}
	dvr, err := os.Open(dvrpath)
	checkErr(err)
	checkErr(filter.Start())
	return ts.NewPktStreamReader(dvr), filter
}
