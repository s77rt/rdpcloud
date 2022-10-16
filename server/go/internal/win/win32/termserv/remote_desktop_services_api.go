//go:build windows && amd64

package termserv

import "unsafe"

type WTS_CONNECTSTATE_CLASS int32

const (
	WTSActive WTS_CONNECTSTATE_CLASS = iota
	WTSConnected
	WTSConnectQuery
	WTSShadow
	WTSDisconnected
	WTSIdle
	WTSListen
	WTSReset
	WTSDown
	WTSInit
)

type WTS_TYPE_CLASS int32

const (
	WTSTypeProcessInfoLevel0 WTS_TYPE_CLASS = iota
	WTSTypeProcessInfoLevel1
	WTSTypeSessionInfoLevel1
)

type WTS_SESSION_INFO_1W struct {
	ExecEnvId    uint32
	State        WTS_CONNECTSTATE_CLASS
	SessionId    uint32
	PSessionName *uint16
	PHostName    *uint16
	PUserName    *uint16
	PDomainName  *uint16
	PFarmName    *uint16
}

const WTS_CURRENT_SERVER_HANDLE = 0

var (
	procWTSEnumerateSessionsExW = modwtsapi32.NewProc("WTSEnumerateSessionsExW")
	procWTSLogoffSession        = modwtsapi32.NewProc("WTSLogoffSession")
	procWTSFreeMemoryExW        = modwtsapi32.NewProc("WTSFreeMemoryExW")
)

func WTSEnumerateSessionsExW(hServer uintptr, pLevel *uint32, Filter uint32, ppSessionInfo **WTS_SESSION_INFO_1W, pCount *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procWTSEnumerateSessionsExW.Call(
		hServer,
		uintptr(unsafe.Pointer(pLevel)),
		uintptr(Filter),
		uintptr(unsafe.Pointer(ppSessionInfo)),
		uintptr(unsafe.Pointer(pCount)),
	)
	return
}

func WTSLogoffSession(hServer uintptr, SessionId uint32, bWait int32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procWTSLogoffSession.Call(
		hServer,
		uintptr(SessionId),
		uintptr(bWait),
	)
	return
}

func WTSFreeMemoryExW(WTSTypeClass WTS_TYPE_CLASS, pMemory uintptr, NumberOfEntries uint64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procWTSFreeMemoryExW.Call(
		uintptr(WTSTypeClass),
		pMemory,
		uintptr(NumberOfEntries),
	)
	return
}
