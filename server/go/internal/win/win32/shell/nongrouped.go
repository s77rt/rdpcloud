//go:build windows && amd64

package shell

import (
	"unsafe"
)

var (
	procDeleteProfileW = moduserenv.NewProc("DeleteProfileW")
)

func DeleteProfileW(lpSidString *uint16, lpProfilePath *uint16, lpComputerName *uint16) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procDeleteProfileW.Call(
		uintptr(unsafe.Pointer(lpSidString)),
		uintptr(unsafe.Pointer(lpProfilePath)),
		uintptr(unsafe.Pointer(lpComputerName)),
	)
	return
}
