// Package demux provides interface to Linux DVB demux device
package demux

import (
	"os"
	"syscall"
	"unsafe"
)

// Filter implements common functionality for all specific filters
type Filter struct {
	file *os.File
}

func (f Filter) Close() error {
	return f.file.Close()
}

func (f Filter) Read(buf []byte) (int, error) {
	return f.file.Read(buf)
}

func (f Filter) Start() error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(f.file.Fd()), _DMX_START, 0,
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f Filter) Stop() error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(f.file.Fd()), _DMX_STOP, 0,
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f Filter) SetBufferLen(n uint32) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(f.file.Fd()),
		_DMX_SET_BUFFER_SIZE, uintptr(n),
	)
	if e != 0 {
		return e
	}
	return nil
}

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
	Filter
}

func (f StreamFilter) AddPid(pid int16) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.file.Fd()),
		_DMX_ADD_PID,
		uintptr(unsafe.Pointer(&pid)),
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f StreamFilter) RemovePid(pid int16) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.file.Fd()),
		_DMX_REMOVE_PID,
		uintptr(unsafe.Pointer(&pid)),
	)
	if e != 0 {
		return e
	}
	return nil
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
	Filter
}

// Dev represents Linux DVB demux device
type Device string

// Returns a handler to elementary stream filter.
func (d Device) StreamFilter(p *StreamFilterParam) (f StreamFilter, err error) {
	f.file, err = os.Open(string(d))
	if err != nil {
		return
	}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.file.Fd()),
		_DMX_SET_PES_FILTER,
		uintptr(unsafe.Pointer(p)),
	)
	if e != 0 {
		err = e
	}
	return
}

// Returns a handler to section filter.
func (d Device) SectionFilter(p *SectionFilterParam) (f SectionFilter, err error) {
	f.file, err = os.Open(string(d))
	if err != nil {
		return
	}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.file.Fd()),
		_DMX_SET_FILTER,
		uintptr(unsafe.Pointer(p)),
	)
	if e != 0 {
		err = e
	}
	return
}
