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
	return ProgramList{Table(pat).Cursor()}
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
	TableCursor
}

// Pop returns first (progId, pid) pair from pl. Remaining pairs are returned
// in rpl. If there is no more programs to read rpl is empty. If an error
// occurs pmtpid == -1
func (pl ProgramList) Pop() (progId uint16, pmtpid int16, rpl ProgramList) {
	if len(pl.Data) == 0 {
		if len(pl.Tab) == 0 {
			return 0, -1, pl
		}
		pl.TableCursor = pl.NextSection()
	}
	if len(pl.Data) < 4 {
		return 0, -1, pl
	}
	progId = decodeU16(pl.Data[0:2])
	pmtpid = int16(decodeU16(pl.Data[2:4]) & 0x1fff)
	rpl.Tab = pl.Tab
	rpl.Data = pl.Data[4:]
	return
}

func (pat *PAT) SetEmpty() {
	(*Table)(pat).SetEmpty()
}

var patCfg = &TableConfig{
	TableId:       0,
	SectionMaxLen: ISOSectionMaxLen,
	GenericSyntax: true,
}

// Append appends next program to PAT. After Append pat is in invalid state.
// Use Close to recalculate all section numbers and CRCs.
func (pat *PAT) Append(progId uint16, pmtpid int16) {
	data := (*Table)(pat).Alloc(4, patCfg, 0, nil)
	encodeU16(data[0:2], progId)
	encodeU16(data[2:4], uint16(pmtpid)|0xe000)
}

func (pat PAT) Close(tsid uint16, current bool, version int8) {
	Table(pat).Close(patCfg, tsid, current, version)
}
