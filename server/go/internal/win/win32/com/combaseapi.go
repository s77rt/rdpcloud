//go:build windows && amd64

package com

import "unsafe"

const COINIT_APARTMENTTHREADED = 0x2

const CLSCTX_INPROC_SERVER = 0x1
const CLSCTX_INPROC_HANDLER = 0x2
const CLSCTX_LOCAL_SERVER = 0x4

var (
	procCoInitializeEx   = modole32.NewProc("CoInitializeEx")
	procCoUninitialize   = modole32.NewProc("CoUninitialize")
	procCoCreateInstance = modole32.NewProc("CoCreateInstance")
)

func CoInitializeEx(pvReserved uintptr, dwCoInit uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procCoInitializeEx.Call(
		uintptr(pvReserved),
		uintptr(dwCoInit),
	)
	return
}

func CoUninitialize() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procCoUninitialize.Call()
	return
}

func CoCreateInstance(rclsid *GUID, pUnkOuter *byte, dwClsContext uint32, riid *GUID, ppv *uintptr) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(rclsid)),
		uintptr(unsafe.Pointer(pUnkOuter)),
		uintptr(dwClsContext),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppv)),
	)
	return
}
