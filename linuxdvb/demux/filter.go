package demux

import (
	"os"
	"sync"
	"syscall"
	"unsafe"
)

type fdset syscall.FdSet

func (s *fdset) Sys() *syscall.FdSet {
	return (*syscall.FdSet)(s)
}

func (s *fdset) Set(fd uintptr) {
	bits := 8 * unsafe.Sizeof(s.Bits[0])
	if fd >= bits*uintptr(len(s.Bits)) {
		panic("fdset: fd out of range")
	}
	n := fd / bits
	m := fd % bits
	s.Bits[n] |= 1 << m
}

func (s *fdset) IsSet(fd uintptr) bool {
	bits := 8 * unsafe.Sizeof(s.Bits[0])
	if fd >= bits*uintptr(len(s.Bits)) {
		panic("fdset: fd out of range")
	}
	n := fd / bits
	m := fd % bits
	return s.Bits[n]&(1<<m) != 0
}

type syn struct {
	pr, pw *os.File
	m      sync.Mutex
}

func (s *syn) pread() error {
	var b [1]byte
	_, err := s.pr.Read(b[:])
	return err
}

var nl = []byte{'\n'}

func (s *syn) pwrite() error {
	_, err := s.pw.Write(nl)
	return err
}

func newSyn() (*syn, error) {
	s := new(syn)
	var err error
	s.pr, s.pw, err = os.Pipe()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *syn) Close() {
	s.pw.Close()
	s.pr.Close()
}

// WaitRead returns true if f can be readed without blocking or false if not or
// error.
func (s *syn) WaitRead(f *os.File) (bool, error) {
	pfd := s.pr.Fd()
	ffd := f.Fd()
	nfd := 1
	if pfd < ffd {
		nfd += int(ffd)
	} else {
		nfd += int(pfd)
	}
	s.m.Lock()
	for {
		var r fdset
		r.Set(ffd)
		r.Set(pfd)
		n, err := syscall.Select(nfd, r.Sys(), nil, nil, nil)
		if err != nil {
			return false, err
		}
		if n > 0 {
			if r.IsSet(pfd) {
				// Command waits for access f.
				s.m.Unlock()
				return false, nil
			}
			return true, nil
		}
	}
}

func (s *syn) Done() {
	s.m.Unlock()
}

func (s *syn) WaitCmd() error {
	if err := s.pwrite(); err != nil {
		return err
	}
	s.m.Lock()
	return s.pread()
}

// Filter implements common functionality for all demux filters.
type Filter struct {
	data *os.File
	s    *syn
}

func newFilter(d Device, typ uintptr, p unsafe.Pointer, dvr bool) (*Filter, error) {
	f, err := os.Open(string(d))
	if err != nil {
		return nil, err
	}
	if _, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		typ,
		uintptr(p),
	); e != 0 {
		return nil, e
	}
	if dvr {
		return &Filter{data: f}, nil
	}
	s, err := newSyn()
	if err != nil {
		return nil, err
	}
	return &Filter{data: f, s: s}, nil
}

func (f *Filter) Close() error {
	if f.s != nil {
		f.s.Close()
	}
	return f.data.Close()
}

func (f *Filter) Read(buf []byte) (int, error) {
	if f.s != nil {
		if ok, err := f.s.WaitRead(f.data); !ok {
			return 0, err
		}
		defer f.s.Done()
	}
	return f.data.Read(buf)
}

func (f *Filter) start() error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(f.data.Fd()), _DMX_START, 0,
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f *Filter) stop() error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(f.data.Fd()), _DMX_STOP, 0,
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f *Filter) setBufferSize(n int) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(f.data.Fd()),
		_DMX_SET_BUFFER_SIZE, uintptr(n),
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f *Filter) addPid(pid int16) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.data.Fd()),
		_DMX_ADD_PID,
		uintptr(unsafe.Pointer(&pid)),
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f *Filter) delPid(pid int16) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(f.data.Fd()),
		_DMX_REMOVE_PID,
		uintptr(unsafe.Pointer(&pid)),
	)
	if e != 0 {
		return e
	}
	return nil
}

func (f *Filter) Start() error {
	if f.s != nil {
		if err := f.s.WaitCmd(); err != nil {
			return err
		}
		defer f.s.Done()
	}
	return f.start()
}

func (f *Filter) Stop() error {
	if f.s != nil {
		if err := f.s.WaitCmd(); err != nil {
			return err
		}
		defer f.s.Done()
	}
	return f.stop()
}

func (f *Filter) SetBufferSize(n int) error {
	if f.s != nil {
		if err := f.s.WaitCmd(); err != nil {
			return err
		}
		defer f.s.Done()
	}
	return f.setBufferSize(n)
}

func (f *Filter) AddPid(pid int16) error {
	if f.s != nil {
		if err := f.s.WaitCmd(); err != nil {
			return err
		}
		defer f.s.Done()
	}
	return f.addPid(pid)
}

func (f *Filter) DelPid(pid int16) error {
	if f.s != nil {
		if err := f.s.WaitCmd(); err != nil {
			return err
		}
		defer f.s.Done()
	}
	return f.delPid(pid)
}
