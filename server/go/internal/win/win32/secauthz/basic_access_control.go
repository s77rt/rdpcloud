//go:build windows && amd64

package secauthz

import "unsafe"

var (
	procConvertSidToStringSidW = modadvapi32.NewProc("ConvertSidToStringSidW")
	procConvertStringSidToSidW = modadvapi32.NewProc("ConvertStringSidToSidW")
	procLookupAccountNameW     = modadvapi32.NewProc("LookupAccountNameW")
	procLookupAccountSidW      = modadvapi32.NewProc("LookupAccountSidW")
)

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
