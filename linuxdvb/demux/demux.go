// Package demux provides interface to Linux DVB demux device
package demux

import (
	"os"
	"syscall"
	"unsafe"
)

// Parameters for StreamFilter

type Input uint32

const (
	InFrontend Input = iota
	InDvr
)

type Output uint32

const (
	OutDecoder Output = iota
	OutTap
	OutTSTap
	OutTSDemuxTap
)

type StreamType uint32

const (
	Audi StreamType = iota
	Video
	Teletext
	Subtitle
	PCR
)

const (
	Audio0 StreamType = iota
	Video0
	Teletext0
	Subtitle0
	PCR0

	Audio1
	Video1
	Teletext1
	Subtitle1
	PCR1

	Audio2
	Video2
	Teletext2
	Subtitle2
	PCR2

	Audio3
	Video3
	Teletext3
	Subtitle3
	PCR3

	Other
)

type Flags uint32

const (
	CheckCRC Flags = 1 << iota
	Oneshot
	ImmediateStart

	KernelClient Flags = 0x8000
)

type StreamFilterParam struct {
	Pid   int16
	In    Input
	Out   Output
	Type  StreamType
	Flags Flags
}

// StreamFilter represents PES filter configured in Linux kernel
type StreamFilter struct {
	*Filter
}

func (f StreamFilter) AddPid(pid int16) error {
	return f.addPid(pid)
}

func (f StreamFilter) DelPid(pid int16) error {
	return f.delPid(pid)
}

// Parameters for SectionFilter

type Pattern struct {
	Bits [16]byte
	Mask [16]byte
	Mode [16]byte
}

type SectionFilterParam struct {
	Pid     int16
	Pattern Pattern
	Timeout uint32
	Flags   Flags
}

// SectionFilter represents filter configured in Linux kernel
type SectionFilter struct {
	*Filter
}

// Dev represents Linux DVB demux device
type Device string

// Returns a handler to elementary stream filter.
func (d Device) NewStreamFilter(p *StreamFilterParam) (StreamFilter, error) {
	f, err := newFilter(d, _DMX_SET_PES_FILTER, unsafe.Pointer(p), p.Out == OutTSTap)
	return StreamFilter{f}, err
}

// Returns a handler to section filter.
func (d Device) NewSectionFilter(p *SectionFilterParam) (SectionFilter, error) {
	f, err := newFilter(d, _DMX_SET_FILTER, unsafe.Pointer(p), false)
	return SectionFilter{f}, err
}

type DVR struct {
	*os.File
}

func (dvr DVR) SetBufferSize(n int) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(dvr.File.Fd()),
		_DMX_SET_BUFFER_SIZE, uintptr(n),
	)
	if e != 0 {
		return e
	}
	return nil
}
