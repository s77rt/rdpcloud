//go:build windows && amd64

package sysinfo

var (
	procGetTickCount64 = modkernel32.NewProc("GetTickCount64")
)

func GetTickCount64() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procGetTickCount64.Call()
	return
}
