package psi

import (
	"github.com/ziutek/dvb/ts"
)

type PAT Table

func (pat PAT) Version() int8 {
	return Table(pat).Version()
}

func (pat PAT) Current() bool {
	return Table(pat).Current()
}

func (pat PAT) MuxId() uint16 {
	return Table(pat).TableIdExt()
}

// Update reads next PAT from r.
func (pat *PAT) Update(r SectionReader, current bool) error {
	return (*Table)(pat).Update(r, 0, false, current, ISOSectionMaxLen)
}

// ProgramList returns list of programs
func (pat PAT) ProgramList() ProgramList {
	return ProgramList{ss: []Section(pat)}
}

// FindPMT returns PMT PID for given progid. If there is no such progId it
// returns pid == ts.NullPid. If an error occurs FindPMT retuns -1
func (pat PAT) FindPMT(progid uint16) (pmtpid int16) {
	pl := pat.ProgramList()
	for !pl.IsEmpty() {
		var id uint16
		id, pmtpid, pl = pl.Pop()
		if pmtpid > ts.NullPid {
			return // Error
		}
		if id == progid {
			return // Found
		}
	}
	return ts.NullPid
}

// FindProgId returns first found program number that corresponds to pmtpid.
// Returns ok == false if not found or error.
func (pat PAT) FindProgId(pmtpid int16) (progid uint16, ok bool) {
	pl := pat.ProgramList()
	for !pl.IsEmpty() {
		var pid int16
		progid, pid, pl = pl.Pop()
		if pid > ts.NullPid {
			return // Error
		}
		ok = (pid == pmtpid)
		if ok {
			return // Found
		}
	}
	return
}

type ProgramList struct {
	ss   []Section
	data []byte
}

func (pl ProgramList) IsEmpty() bool {
	return len(pl.ss) == 0 && len(pl.data) == 0
}

// Pop returns first (progId, pid) pair from pl. Remaining pairs are returned
// in rpl. If there is no more programs to read rpl is empty.
// If an error occurs pmtpid == -1
func (pl ProgramList) Pop() (progId uint16, pmtpid int16, rpl ProgramList) {
	if len(pl.data) == 0 {
		if len(pl.ss) == 0 {
			return
		}
		pl.data = pl.ss[0].Data()
		pl.ss = pl.ss[1:]
	}
	if len(pl.data) < 4 {
		pmtpid = -1
		return
	}
	progId = decodeU16(pl.data[0:2])
	pmtpid = int16(decodeU16(pl.data[2:4]) & 0x1fff)
	rpl.ss = pl.ss
	rpl.data = pl.data[4:]
	return
}
