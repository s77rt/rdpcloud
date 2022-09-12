//go:build windows && amd64

package sysinfo

var (
	procCloseHandle = modkernel32.NewProc("CloseHandle")
)

func CloseHandle(hObject uintptr) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procCloseHandle.Call(hObject)
	return
}
