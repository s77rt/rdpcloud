//go:build windows && amd64

package shutdown

import "unsafe"

var (
	procInitiateSystemShutdownExW = modadvapi32.NewProc("InitiateSystemShutdownExW")
	procAbortSystemShutdownW      = modadvapi32.NewProc("AbortSystemShutdownW")
)

func InitiateSystemShutdownExW(lpMachineName *uint16, lpMessage *uint16, dwTimeout uint32, bForceAppsClosed int32, bRebootAfterShutdown int32, dwReason uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procInitiateSystemShutdownExW.Call(
		uintptr(unsafe.Pointer(lpMachineName)),
		uintptr(unsafe.Pointer(lpMessage)),
		uintptr(dwTimeout),
		uintptr(bForceAppsClosed),
		uintptr(bRebootAfterShutdown),
		uintptr(dwReason),
	)
	return
}

func AbortSystemShutdownW(lpMachineName *uint16) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procAbortSystemShutdownW.Call(
		uintptr(unsafe.Pointer(lpMachineName)),
	)
	return
}
