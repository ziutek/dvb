package frontend

// API3
const (
	_FE_GET_INFO                = 0x80a86f3d
	_FE_GET_FRONTEND            = 0x80246f4d
	_FE_SET_FRONTEND            = 0x40246f4c
	_FE_GET_EVENT               = 0x80286f4e
	_FE_READ_STATUS             = 0x80046f45
	_FE_READ_BER                = 0x80046f46
	_FE_READ_SIGNAL_STRENGTH    = 0x80026f47
	_FE_READ_SNR                = 0x80026f48
	_FE_READ_UNCORRECTED_BLOCKS = 0x80046f49

	_FE_SET_TONE                = 0x6f42
	_FE_SET_VOLTAGE             = 0x6f42
	_FE_ENABLE_HIGH_LNB_VOLTAGE = 0x6f42
)

// API5
const (
	// Values for kernel 3.16 (check using checksyscall/main.c)
	_FE_SET_PROPERTY = 0x40106f52
	_FE_GET_PROPERTY = 0x80106f53
)
