package localfile

import (
	"os"
	"syscall"
	"unsafe"
)

var lockFileEx = syscall.NewLazyDLL("kernel32.dll").NewProc("LockFileEx")

func Lock(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	const exclusiveAndFailImmediately = 3
	var offset syscall.Overlapped
	ok, _, err := lockFileEx.Call(f.Fd(), exclusiveAndFailImmediately, 0, 1, 0, uintptr(unsafe.Pointer(&offset)))
	if ok == 0 {
		f.Close()
		return nil, err
	}
	return f, nil
}
