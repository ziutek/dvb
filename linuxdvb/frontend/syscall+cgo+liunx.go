// +build cgo,linux

package frontend

/*
#include <sys/ioctl.h>
#include <linux/dvb/frontend.h>
*/
import "C"

// API3
const (
	_FE_GET_INFO                = C.FE_GET_INFO
	_FE_GET_FRONTEND            = C.FE_GET_FRONTEND
	_FE_SET_FRONTEND            = C.FE_SET_FRONTEND
	_FE_GET_EVENT               = C.FE_GET_EVENT
	_FE_READ_STATUS             = C.FE_READ_STATUS
	_FE_READ_BER                = C.FE_READ_BER
	_FE_READ_SIGNAL_STRENGTH    = C.FE_READ_SIGNAL_STRENGTH
	_FE_READ_SNR                = C.FE_READ_SNR
	_FE_READ_UNCORRECTED_BLOCKS = C.FE_READ_UNCORRECTED_BLOCKS

	_FE_SET_TONE                = C.FE_SET_TONE
	_FE_SET_VOLTAGE             = C.FE_SET_VOLTAGE
	_FE_ENABLE_HIGH_LNB_VOLTAGE = C.FE_ENABLE_HIGH_LNB_VOLTAGE
)

// API5
const (
	_FE_SET_PROPERTY = C.FE_SET_PROPERTY
	_FE_GET_PROPERTY = C.FE_GET_PROPERTY
)
