//go:build windows && amd64

package msi

import "unsafe"

var (
	procMsiEnumProductsW   = modmsi.NewProc("MsiEnumProductsW")
	procMsiGetProductInfoW = modmsi.NewProc("MsiGetProductInfoW")
)

func MsiEnumProductsW(iProductIndex uint32, lpProductBuf *uint16) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procMsiEnumProductsW.Call(
		uintptr(iProductIndex),
		uintptr(unsafe.Pointer(lpProductBuf)),
	)
	return
}

func MsiGetProductInfoW(szProduct *uint16, szAttribute *uint16, lpValueBuf *uint16, pcchValueBuf *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procMsiGetProductInfoW.Call(
		uintptr(unsafe.Pointer(szProduct)),
		uintptr(unsafe.Pointer(szAttribute)),
		uintptr(unsafe.Pointer(lpValueBuf)),
		uintptr(unsafe.Pointer(pcchValueBuf)),
	)
	return
}
