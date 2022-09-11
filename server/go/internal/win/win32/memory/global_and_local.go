//go:build windows && amd64

package memory

import "unsafe"

var (
	procLocalFree = modkernel32.NewProc("LocalFree")
)

func LocalFree(hMem *byte) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procLocalFree.Call(uintptr(unsafe.Pointer(hMem)))
	return
}
