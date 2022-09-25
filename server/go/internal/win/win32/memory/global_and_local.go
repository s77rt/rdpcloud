//go:build windows && amd64

package memory

var (
	procLocalFree = modkernel32.NewProc("LocalFree")
)

func LocalFree(hMem uintptr) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procLocalFree.Call(hMem)
	return
}
