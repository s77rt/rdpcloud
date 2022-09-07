//go:build windows && amd64

package netmgmt

import (
	"unsafe"
)

var (
	procNetUserAdd            = modnetapi32.NewProc("NetUserAdd")
	procNetUserChangePassword = modnetapi32.NewProc("NetUserChangePassword")
	procNetUserDel            = modnetapi32.NewProc("NetUserDel")
	procNetUserEnum           = modnetapi32.NewProc("NetUserEnum")
	procNetUserGetInfo        = modnetapi32.NewProc("NetUserGetInfo")
	procNetUserSetInfo        = modnetapi32.NewProc("NetUserSetInfo")
)

func NetUserAdd(servername *uint16, level uint32, buf *byte, parm_err *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserAdd.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(parm_err)),
	)
	return
}

func NetUserEnum(servername *uint16, level uint32, filter uint32, buf **byte, prefmaxlen uint32, entriesread *uint32, totalentries *uint32, resume_handle *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserEnum.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(level),
		uintptr(filter),
		uintptr(unsafe.Pointer(buf)),
		uintptr(prefmaxlen),
		uintptr(unsafe.Pointer(entriesread)),
		uintptr(unsafe.Pointer(totalentries)),
		uintptr(unsafe.Pointer(resume_handle)),
	)
	return
}
