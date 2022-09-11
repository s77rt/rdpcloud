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
	procLogonUserExA = modadvapi32.NewProc("LogonUserExA")
)

func LogonUserExA(lpszUsername *uint8, lpszDomain *uint8, lpszPassword *uint8, dwLogonType uint32, dwLogonProvider uint32, phToken *byte, ppLogonSid **byte, ppProfileBuffer **byte, pdwProfileLength *uint32, pQuotaLimits *byte) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procLogonUserExA.Call(
		uintptr(unsafe.Pointer(lpszUsername)),
		uintptr(unsafe.Pointer(lpszDomain)),
		uintptr(unsafe.Pointer(lpszPassword)),
		uintptr(dwLogonType),
		uintptr(dwLogonProvider),
		uintptr(unsafe.Pointer(phToken)),
		uintptr(unsafe.Pointer(ppLogonSid)),
		uintptr(unsafe.Pointer(ppProfileBuffer)),
		uintptr(unsafe.Pointer(pdwProfileLength)),
		uintptr(unsafe.Pointer(pQuotaLimits)),
	)
	return
}
