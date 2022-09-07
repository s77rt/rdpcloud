//go:build windows && amd64

package netmgmt

import (
	"unsafe"
)

var (
	procNetApiBufferFree = modnetapi32.NewProc("NetApiBufferFree")
)

func NetApiBufferFree(buf *byte) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetApiBufferFree.Call(uintptr(unsafe.Pointer(buf)))
	return
}
