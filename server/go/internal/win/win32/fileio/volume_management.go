//go:build windows && amd64

package fileio

import "unsafe"

var (
	procFindFirstVolumeW                 = modkernel32.NewProc("FindFirstVolumeW")
	procFindNextVolumeW                  = modkernel32.NewProc("FindNextVolumeW")
	procFindVolumeClose                  = modkernel32.NewProc("FindVolumeClose")
	procGetVolumePathNamesForVolumeNameW = modkernel32.NewProc("GetVolumePathNamesForVolumeNameW")
)

func FindFirstVolumeW(lpszVolumeName *uint16, cchBufferLength uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procFindFirstVolumeW.Call(
		uintptr(unsafe.Pointer(lpszVolumeName)),
		uintptr(cchBufferLength),
	)
	return
}

func FindNextVolumeW(hFindVolume uintptr, lpszVolumeName *uint16, cchBufferLength uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procFindNextVolumeW.Call(
		hFindVolume,
		uintptr(unsafe.Pointer(lpszVolumeName)),
		uintptr(cchBufferLength),
	)
	return
}

func FindVolumeClose(hFindVolume uintptr) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procFindVolumeClose.Call(
		hFindVolume,
	)
	return
}

func GetVolumePathNamesForVolumeNameW(lpszVolumeName *uint16, lpszVolumePathNames *uint16, cchBufferLength uint32, lpcchReturnLength *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procGetVolumePathNamesForVolumeNameW.Call(
		uintptr(unsafe.Pointer(lpszVolumeName)),
		uintptr(unsafe.Pointer(lpszVolumePathNames)),
		uintptr(cchBufferLength),
		uintptr(unsafe.Pointer(lpcchReturnLength)),
	)
	return
}
