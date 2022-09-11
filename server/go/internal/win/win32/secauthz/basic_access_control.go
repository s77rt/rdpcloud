//go:build windows && amd64

package secauthz

import "unsafe"

var (
	procConvertSidToStringSidA = modadvapi32.NewProc("ConvertSidToStringSidA")
)

func ConvertSidToStringSidA(Sid *byte, StringSid **uint8) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procConvertSidToStringSidA.Call(
		uintptr(unsafe.Pointer(Sid)),
		uintptr(unsafe.Pointer(StringSid)),
	)
	return
}
