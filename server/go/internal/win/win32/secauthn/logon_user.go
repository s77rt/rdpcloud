//go:build windows && amd64

package secauthn

import "unsafe"

const LOGON32_LOGON_INTERACTIVE = 2
const LOGON32_LOGON_NETWORK = 3
const LOGON32_LOGON_BATCH = 4
const LOGON32_LOGON_SERVICE = 5
const LOGON32_LOGON_UNLOCK = 7
const LOGON32_LOGON_NETWORK_CLEARTEXT = 8
const LOGON32_LOGON_NEW_CREDENTIALS = 9

const LOGON32_PROVIDER_DEFAULT = 0
const LOGON32_PROVIDER_WINNT35 = 1
const LOGON32_PROVIDER_WINNT40 = 2
const LOGON32_PROVIDER_WINNT50 = 3
const LOGON32_PROVIDER_VIRTUAL = 4

var (
	procLogonUserW = modadvapi32.NewProc("LogonUserW")
)

func LogonUserW(lpszUsername *uint16, lpszDomain *uint16, lpszPassword *uint16, dwLogonType uint32, dwLogonProvider uint32, phToken *uintptr) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procLogonUserW.Call(
		uintptr(unsafe.Pointer(lpszUsername)),
		uintptr(unsafe.Pointer(lpszDomain)),
		uintptr(unsafe.Pointer(lpszPassword)),
		uintptr(dwLogonType),
		uintptr(dwLogonProvider),
		uintptr(unsafe.Pointer(phToken)),
	)
	return
}
