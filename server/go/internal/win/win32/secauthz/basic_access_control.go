//go:build windows && amd64

package secauthz

import "unsafe"

type LUID struct {
	LowPart  uint32
	HighPart int32
}

type LUID_AND_ATTRIBUTES struct {
	Luid       LUID
	Attributes uint32
}

type TOKEN_PRIVILEGES struct {
	PrivilegeCount uint32
	Privileges     [ANYSIZE_ARRAY]LUID_AND_ATTRIBUTES
}

const ANYSIZE_ARRAY = 1

const TOKEN_QUERY = 0x0008
const TOKEN_ADJUST_PRIVILEGES = 0x0020

const SE_PRIVILEGE_ENABLED_BY_DEFAULT = 0x00000001
const SE_PRIVILEGE_ENABLED = 0x00000002
const SE_PRIVILEGE_REMOVED = 0x00000004
const SE_PRIVILEGE_USED_FOR_ACCESS = 0x80000000
const SE_PRIVILEGE_VALID_ATTRIBUTES = 0x80000007

var (
	procAdjustTokenPrivileges  = modadvapi32.NewProc("AdjustTokenPrivileges")
	procConvertSidToStringSidW = modadvapi32.NewProc("ConvertSidToStringSidW")
	procConvertStringSidToSidW = modadvapi32.NewProc("ConvertStringSidToSidW")
	procLookupAccountNameW     = modadvapi32.NewProc("LookupAccountNameW")
	procLookupAccountSidW      = modadvapi32.NewProc("LookupAccountSidW")
	procLookupPrivilegeValueW  = modadvapi32.NewProc("LookupPrivilegeValueW")
	procOpenProcessToken       = modadvapi32.NewProc("OpenProcessToken")
)

func AdjustTokenPrivileges(TokenHandle uintptr, DisableAllPrivileges int32, NewState *TOKEN_PRIVILEGES, BufferLength uint32, PreviousState *TOKEN_PRIVILEGES, ReturnLength *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procAdjustTokenPrivileges.Call(
		uintptr(TokenHandle),
		uintptr(DisableAllPrivileges),
		uintptr(unsafe.Pointer(NewState)),
		uintptr(BufferLength),
		uintptr(unsafe.Pointer(PreviousState)),
		uintptr(unsafe.Pointer(ReturnLength)),
	)
	return
}

func ConvertSidToStringSidW(Sid *byte, StringSid **uint16) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procConvertSidToStringSidW.Call(
		uintptr(unsafe.Pointer(Sid)),
		uintptr(unsafe.Pointer(StringSid)),
	)
	return
}

func ConvertStringSidToSidW(StringSid *uint16, Sid **byte) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procConvertStringSidToSidW.Call(
		uintptr(unsafe.Pointer(StringSid)),
		uintptr(unsafe.Pointer(Sid)),
	)
	return
}

func LookupAccountNameW(lpSystemName *uint16, lpAccountName *uint16, Sid *byte, cbSid *uint32, ReferencedDomainName *uint16, cchReferencedDomainName *uint32, peUse *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procLookupAccountNameW.Call(
		uintptr(unsafe.Pointer(lpSystemName)),
		uintptr(unsafe.Pointer(lpAccountName)),
		uintptr(unsafe.Pointer(Sid)),
		uintptr(unsafe.Pointer(cbSid)),
		uintptr(unsafe.Pointer(ReferencedDomainName)),
		uintptr(unsafe.Pointer(cchReferencedDomainName)),
		uintptr(unsafe.Pointer(peUse)),
	)
	return
}

func LookupAccountSidW(lpSystemName *uint16, Sid *byte, Name *uint16, cchName *uint32, ReferencedDomainName *uint16, cchReferencedDomainName *uint32, peUse *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procLookupAccountSidW.Call(
		uintptr(unsafe.Pointer(lpSystemName)),
		uintptr(unsafe.Pointer(Sid)),
		uintptr(unsafe.Pointer(Name)),
		uintptr(unsafe.Pointer(cchName)),
		uintptr(unsafe.Pointer(ReferencedDomainName)),
		uintptr(unsafe.Pointer(cchReferencedDomainName)),
		uintptr(unsafe.Pointer(peUse)),
	)
	return
}

func LookupPrivilegeValueW(lpSystemName *uint16, lpName *uint16, lpLuid *LUID) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procLookupPrivilegeValueW.Call(
		uintptr(unsafe.Pointer(lpSystemName)),
		uintptr(unsafe.Pointer(lpName)),
		uintptr(unsafe.Pointer(lpLuid)),
	)
	return
}

func OpenProcessToken(ProcessHandle uintptr, DesiredAccess uint32, TokenHandle *uintptr) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procOpenProcessToken.Call(
		uintptr(ProcessHandle),
		uintptr(DesiredAccess),
		uintptr(unsafe.Pointer(TokenHandle)),
	)
	return
}
