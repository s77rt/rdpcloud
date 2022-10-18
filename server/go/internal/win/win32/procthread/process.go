//go:build windows && amd64

package procthread

var (
	procGetCurrentProcess = modkernel32.NewProc("GetCurrentProcess")
)

func GetCurrentProcess() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procGetCurrentProcess.Call()
	return
}
