// +build linux,cgo

package demux

/*
#include <sys/ioctl.h>
#include <linux/dvb/dmx.h>
*/
import "C"

const (
	_DMX_START           = C.DMX_START
	_DMX_STOP            = C.DMX_STOP
	_DMX_SET_BUFFER_SIZE = C.DMX_SET_BUFFER_SIZE
	_DMX_SET_FILTER      = C.DMX_SET_FILTER
	_DMX_SET_PES_FILTER  = C.DMX_SET_PES_FILTER
	_DMX_GET_PES_PIDS    = C.DMX_GET_PES_PIDS
	_DMX_GET_CAPS        = C.DMX_GET_CAPS
	_DMX_SET_SOURCE      = C.DMX_SET_SOURCE
	_DMX_GET_STC         = C.DMX_GET_STC
	_DMX_ADD_PID         = C.DMX_ADD_PID
	_DMX_REMOVE_PID      = C.DMX_REMOVE_PID
)
