package psi

import (
	"github.com/ziutek/dvb"
)

type SDT Table

func NewSDT() *SDT {
	return (*SDT)(NewTable(ISOSectionMaxLen))
}

func (sdt *SDT) t() *Table {
	return (*Table)(sdt)
}

func (sdt *SDT) Version() int8 {
	return sdt.t().Version()
}

func (sdt *SDT) Current() bool {
	return sdt.t().Current()
}

func (sdt *SDT) MuxId() uint16 {
	return sdt.t().TableIdExt()
}

func (sdt *SDT) s() Section {
	ss := sdt.t().Sections()
	if len(ss) == 0 {
		panic("SDT doesn't contain valid data")
	}
	return ss[0]
}

var ErrSDTSectionLen = dvb.TemporaryError("incorrect SDT section length")

// Update reads next SDT from r
func (sdt *SDT) Update(r SectionReader, actualMux bool, current bool) error {
	tableId := byte(0x46)
	if actualMux {
		tableId = 0x42
	}
	err := sdt.t().Update(r, tableId, true, current)
	if err != nil {
		return err
	}
	if len(sdt.s().Data()) < 2 {
		sdt.t().Reset()
		return ErrSDTSectionLen
	}
	return nil
}

// OrgNetId returns original_network_id
func (sdt *SDT) OrgNetId() uint16 {
	return decodeU16(sdt.s().Data()[0:2])
}

// Info returns list of ifnormation about services (programs)
func (sdt *SDT) ServiceInfo() ServiceInfoList {
	return ServiceInfoList{ss: sdt.t().Sections()}
}

type ServiceInfoList struct {
	ss   []Section
	data []byte
}

func (sl ServiceInfoList) IsEmpty() bool {
	return len(sl.ss) == 0 && len(sl.data) == 0
}

// Pop returns first service information element from sl. Remaining elements
// are returned in rsl. If there is no more informations to read rsl is empty.
// If an error occurs si == nil.
func (sl ServiceInfoList) Pop() (si ServiceInfo, rsl ServiceInfoList) {
	if len(sl.data) == 0 {
		if len(sl.ss) == 0 {
			return
		}
		sl.data = sl.ss[0].Data()
		if len(sl.data) < 3 {
			return
		}
		sl.data = sl.data[3:]
		sl.ss = sl.ss[1:]
	}
	if len(sl.data) < 5 {
		return
	}
	l := int(decodeU16(sl.data[3:5])&0x0fff) + 5
	if len(sl.data) < l {
		return
	}
	si = sl.data[:l]
	rsl.ss = sl.ss
	rsl.data = sl.data[l:]
	return
}

type ServiceStatus byte

const (
	StatusUndefined ServiceStatus = iota
	NotRunnind
	StartsInFewSeconds
	Pausing
	Running
)

var ssn = []string{
	"undefined",
	"not runnind",
	"starts in few seconds",
	"pausing",
	"running",
}

func (ss ServiceStatus) String() string {
	if ss > Running {
		return "unknown"
	}
	return ssn[ss]
}

type ServiceInfo []byte

// ServiceId return id of service (program) that this information applies to.
func (si ServiceInfo) ServiceId() uint16 {
	return decodeU16(si[0:2])
}

// EITSchedule returns the value of EIT_schedule_flag field.
func (si ServiceInfo) EITSchedule() bool {
	return si[2]&0x02 != 0
}

// EITPresentFollowing returns the value of EIT_present_following_flag field.
func (si ServiceInfo) EITPresentFollowing() bool {
	return si[2]&0x01 != 0
}

// Status returns the value of running_status field.
func (si ServiceInfo) Status() ServiceStatus {
	return ServiceStatus(si[3] >> 5)
}

// Scrambled returns the value of free_CA_mode field.
func (si ServiceInfo) Scrambled() bool {
	return si[3]&0x10 != 0
}

func (si ServiceInfo) Descriptors() DescriptorList {
	l := int(decodeU16(si[3:5]) & 0x0fff)
	return DescriptorList(si[5 : 5+l])
}
