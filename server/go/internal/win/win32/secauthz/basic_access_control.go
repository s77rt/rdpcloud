//go:build windows && amd64

package secauthz

import "unsafe"

var (
	procConvertSidToStringSidW = modadvapi32.NewProc("ConvertSidToStringSidW")
	procLookupAccountNameW     = modadvapi32.NewProc("LookupAccountNameW")
)

func ConvertSidToStringSidW(Sid *byte, StringSid **uint16) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procConvertSidToStringSidW.Call(
		uintptr(unsafe.Pointer(Sid)),
		uintptr(unsafe.Pointer(StringSid)),
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
