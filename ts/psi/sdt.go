package psi

import (
	"github.com/ziutek/dvb"
)

type SDT Table

func (sdt SDT) Version() int8 {
	return Table(sdt).Version()
}

func (sdt SDT) Current() bool {
	return Table(sdt).Current()
}

func (sdt SDT) MuxId() uint16 {
	return Table(sdt).TableIdExt()
}

var ErrSDTSectionLen = dvb.TemporaryError("incorrect SDT section length")

// Update reads next SDT from r
func (sdt *SDT) Update(r SectionReader, actualMux bool, current bool) error {
	tableId := byte(0x46)
	if actualMux {
		tableId = 0x42
	}
	t := (*Table)(sdt)
	err := t.Update(r, tableId, true, current, ISOSectionMaxLen)
	if err != nil {
		return err
	}
	for _, s := range *t {
		if len(s.Data()) < 2 {
			t.Reset()
			return ErrSDTSectionLen
		}
	}
	return nil
}

// OrgNetId returns original_network_id
func (sdt SDT) OrgNetId() uint16 {
	return decodeU16(sdt[0].Data()[0:2])
}

// Info returns list of ifnormation about services (programs)
func (sdt SDT) ServiceInfo() ServiceInfoList {
	return ServiceInfoList{Table(sdt).Cursor()}
}

func (sdt *SDT) SetEmpty() {
	(*Table)(sdt).SetEmpty()
}

var (
	sdtCfgActual = &TableConfig{
		TableId:        0x42,
		SectionMaxLen:  ISOSectionMaxLen,
		SectionHeadLen: 3,
		GenericSyntax:  true,
		PrivateSyntax:  true,
	}
	sdtCfgOther = &TableConfig{
		TableId:        0x46,
		SectionMaxLen:  ISOSectionMaxLen,
		SectionHeadLen: 3,
		GenericSyntax:  true,
		PrivateSyntax:  true,
	}
)

func (sdt *SDT) Append(onid uint16, si ServiceInfo) {
	var head [3]byte
	encodeU16(head[0:2], onid)
	head[2] = 0xff
	data := (*Table)(sdt).Alloc(len(si), sdtCfgActual, 0, head[:])
	copy(data, si)
}

func (sdt SDT) Close(tsid uint16, actualMux, current bool, version int8) {
	sdtcfg := sdtCfgOther
	if actualMux {
		sdtcfg = sdtCfgActual
	}
	Table(sdt).Close(sdtcfg, tsid, current, version)
}

type ServiceInfoList struct {
	TableCursor
}

// Pop returns first ServiceInfo element from sl. If there is no more data to
// read Pop returns empty ServiceInfoList. If an error occurs it returns nil
// ServiceInfo.
func (sl ServiceInfoList) Pop() (ServiceInfo, ServiceInfoList) {
	if len(sl.Data) == 0 {
		if len(sl.Tab) == 0 {
			return nil, sl
		}
		sl.TableCursor = sl.NextSection()
		// Skip original_network_id.
		if len(sl.Data) < 3 {
			return nil, sl
		}
		sl.Data = sl.Data[3:]
	}
	if len(sl.Data) < 5 {
		return nil, sl
	}
	n := int(decodeU16(sl.Data[3:5])&0x0fff) + 5
	if len(sl.Data) < n {
		return nil, sl
	}
	data := sl.Data[:n]
	sl.Data = sl.Data[n:]
	return data, sl
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

func (si ServiceInfo) SetServiceId(sid uint16) {
	encodeU16(si[0:2], sid)
}

// EITSchedule returns the value of EIT_schedule_flag field.
func (si ServiceInfo) EITSchedule() bool {
	return si[2]&0x02 != 0
}

func (si ServiceInfo) SetEITSchedule(b bool) {
	if b {
		si[2] |= 0x02
	} else {
		si[2] &^= 0x02
	}
}

// EITPresentFollowing returns the value of EIT_present_following_flag field.
func (si ServiceInfo) EITPresentFollowing() bool {
	return si[2]&0x01 != 0
}

// Status returns the value of running_status field.
func (si ServiceInfo) Status() ServiceStatus {
	return ServiceStatus(si[3] >> 5)
}

// SetStatus sets running_status field.
func (si ServiceInfo) SetStatus(s ServiceStatus) {
	si[3] = si[3]&0x1f | byte(s<<5)
}

// Scrambled returns the value of free_CA_mode field.
func (si ServiceInfo) Scrambled() bool {
	return si[3]&0x10 != 0
}

func (si ServiceInfo) descrLoopLen() int {
	return loopLen(si[3:5])
}

func (si ServiceInfo) setDescrLoopLen(n int) {
	setLoopLen(si[3:5], n)
}

func (si ServiceInfo) Descriptors() DescriptorList {
	return DescriptorList(si[5 : 5+si.descrLoopLen()])
}

func MakeServiceInfo() ServiceInfo {
	si := make(ServiceInfo, 5)
	si[2] = 0xfc
	return si
}

// ClearDescriptors clears descriptors_loop_length field.
func (si ServiceInfo) ClearDescriptors() {
	si.setDescrLoopLen(0)
}

func (si *ServiceInfo) AppendDescriptors(ds ...Descriptor) {
	n := si.descrLoopLen()
	for _, d := range ds {
		*si = append((*si)[:5+n], d...)
		n += len(d)
	}
	si.setDescrLoopLen(n)
}
