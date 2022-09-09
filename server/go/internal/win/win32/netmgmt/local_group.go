//go:build windows && amd64

package netmgmt

import (
	"unsafe"
)

type LOCALGROUP_INFO_0 struct {
	Lgrpi0_name *uint16
}

type LOCALGROUP_MEMBERS_INFO_3 struct {
	Lgrmi3_domainandname *uint16
}

var (
	procNetLocalGroupAddMembers = modnetapi32.NewProc("NetLocalGroupAddMembers")
	procNetLocalGroupDelMembers = modnetapi32.NewProc("NetLocalGroupDelMembers")
	procNetLocalGroupEnum       = modnetapi32.NewProc("NetLocalGroupEnum")
	procNetLocalGroupGetMembers = modnetapi32.NewProc("NetLocalGroupGetMembers")
)

func NetLocalGroupAddMembers(servername *uint16, groupname *uint16, level uint32, buf *byte, totalentries uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetLocalGroupAddMembers.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(groupname)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(totalentries),
	)
	return
}

func NetLocalGroupDelMembers(servername *uint16, groupname *uint16, level uint32, buf *byte, totalentries uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetLocalGroupDelMembers.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(groupname)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(totalentries),
	)
	return
}

func NetLocalGroupEnum(servername *uint16, level uint32, buf **byte, prefmaxlen uint32, entriesread *uint32, totalentries *uint32, resumehandle *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetLocalGroupEnum.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(prefmaxlen),
		uintptr(unsafe.Pointer(entriesread)),
		uintptr(unsafe.Pointer(totalentries)),
		uintptr(unsafe.Pointer(resumehandle)),
	)
	return
}

func NetLocalGroupGetMembers(servername *uint16, groupname *uint16, level uint32, buf **byte, prefmaxlen uint32, entriesread *uint32, totalentries *uint32, resumehandle *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetLocalGroupGetMembers.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(groupname)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(prefmaxlen),
		uintptr(unsafe.Pointer(entriesread)),
		uintptr(unsafe.Pointer(totalentries)),
		uintptr(unsafe.Pointer(resumehandle)),
	)
	return
}
