package psi

import (
	"github.com/ziutek/dvb/ts"
)

type PAT Table

func NewPAT() *PAT {
	return (*PAT)(NewTable(ISOSectionMaxLen))
}

func (pat *PAT) t() *Table {
	return (*Table)(pat)
}

func (pat *PAT) Version() byte {
	return pat.t().Version()
}

func (pat *PAT) Current() bool {
	return pat.t().Current()
}

func (pat *PAT) MuxId() uint16 {
	return pat.t().TableIdExt()
}

// Update reads next PAT from r.
func (pat *PAT) Update(r SectionReader, current bool) error {
	return pat.t().Update(r, 0, false, current)
}

// ProgramList returns list of programs
func (pat *PAT) ProgramList() ProgramList {
	return ProgramList{ss: pat.t().Sections()}
}

// FindPMT returns PMT PID for given progid. If there is no such progId it
// returns pid == ts.NullPid. If an error occurs pid > ts.NullPid.
func (pat *PAT) FindPMT(progid uint16) (pmtpid uint16) {
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
func (pat *PAT) FindProgId(pmtpid uint16) (progid uint16, ok bool) {
	pl := pat.ProgramList()
	for !pl.IsEmpty() {
		var pid uint16
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
// If an error occurs pid > ts.NullPid
func (pl ProgramList) Pop() (progId, pmtpid uint16, rpl ProgramList) {
	if len(pl.data) == 0 {
		if len(pl.ss) == 0 {
			return
		}
		pl.data = pl.ss[0].Data()
		pl.ss = pl.ss[1:]
	}
	if len(pl.data) < 4 {
		pmtpid = ts.NullPid + 1
		return
	}
	progId = decodeU16(pl.data[0:2])
	pmtpid = decodeU16(pl.data[2:4]) & 0x1fff
	rpl.ss = pl.ss
	rpl.data = pl.data[4:]
	return
}
