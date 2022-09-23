//go:build windows && amd64

package com

import (
	"syscall"
	"unsafe"
)

type IUnknown struct {
	lpVtbl *IUnknownVtbl
}

type IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

func (x *IUnknown) AddRef() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.AddRef,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

func (x *IUnknown) Release() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.Release,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

var IID_IUnknown = &GUID{0x00000000, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
