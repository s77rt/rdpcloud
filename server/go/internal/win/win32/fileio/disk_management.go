//go:build windows && amd64

package fileio

import (
	"syscall"
	"unsafe"

	comInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/com"
)

type IDiskQuotaControl struct {
	lpVtbl *IDiskQuotaControlVtbl
}

type IDiskQuotaControlVtbl struct {
	// IUnknown
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	// IConnectionPointContainer
	EnumConnectionPoints uintptr
	FindConnectionPoint  uintptr

	Initialize                     uintptr
	SetQuotaState                  uintptr
	GetQuotaState                  uintptr
	SetQuotaLogFlags               uintptr
	GetQuotaLogFlags               uintptr
	SetDefaultQuotaThreshold       uintptr
	GetDefaultQuotaThreshold       uintptr
	GetDefaultQuotaThresholdText   uintptr
	SetDefaultQuotaLimit           uintptr
	GetDefaultQuotaLimit           uintptr
	GetDefaultQuotaLimitText       uintptr
	AddUserSid                     uintptr
	AddUserName                    uintptr
	DeleteUser                     uintptr
	FindUserSid                    uintptr
	FindUserName                   uintptr
	CreateEnumUsers                uintptr
	CreateUserBatch                uintptr
	InvalidateSidNameCache         uintptr
	GiveUserNameResolutionPriority uintptr
	ShutdownNameResolution         uintptr
}

func (x *IDiskQuotaControl) AddRef() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.AddRef,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

func (x *IDiskQuotaControl) Release() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.Release,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

func (x *IDiskQuotaControl) Initialize(pszPath *uint16, bReadWrite int32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.Initialize,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pszPath)),
		uintptr(bReadWrite),
	)
	return
}

func (x *IDiskQuotaControl) SetQuotaState(dwState uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.SetQuotaState,
		uintptr(unsafe.Pointer(x)),
		uintptr(dwState),
	)
	return
}

func (x *IDiskQuotaControl) GetQuotaState(pdwState *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetQuotaState,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pdwState)),
	)
	return
}

func (x *IDiskQuotaControl) SetDefaultQuotaThreshold(llThreshold int64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.SetDefaultQuotaThreshold,
		uintptr(unsafe.Pointer(x)),
		uintptr(llThreshold),
	)
	return
}

func (x *IDiskQuotaControl) GetDefaultQuotaThreshold(pllThreshold *int64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetDefaultQuotaThreshold,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pllThreshold)),
	)
	return
}

func (x *IDiskQuotaControl) SetDefaultQuotaLimit(llLimit int64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.SetDefaultQuotaLimit,
		uintptr(unsafe.Pointer(x)),
		uintptr(llLimit),
	)
	return
}

func (x *IDiskQuotaControl) GetDefaultQuotaLimit(pllLimit *int64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetDefaultQuotaLimit,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pllLimit)),
	)
	return
}

func (x *IDiskQuotaControl) AddUserName(pszLogonName *uint16, fNameResolution uint32, ppUser **IDiskQuotaUser) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.AddUserName,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pszLogonName)),
		uintptr(fNameResolution),
		uintptr(unsafe.Pointer(ppUser)),
	)
	return
}

func (x *IDiskQuotaControl) DeleteUser(ppUser *IDiskQuotaUser) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.DeleteUser,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(ppUser)),
	)
	return
}

func (x *IDiskQuotaControl) FindUserName(pszLogonName *uint16, ppUser **IDiskQuotaUser) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.FindUserName,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pszLogonName)),
		uintptr(unsafe.Pointer(ppUser)),
	)
	return
}

func (x *IDiskQuotaControl) CreateEnumUsers(rgpUserSids **byte, cpSids uint32, fNameResolution uint32, ppEnum **IEnumDiskQuotaUsers) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.CreateEnumUsers,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(rgpUserSids)),
		uintptr(cpSids),
		uintptr(fNameResolution),
		uintptr(unsafe.Pointer(ppEnum)),
	)
	return
}

type IDiskQuotaUser struct {
	lpVtbl *IDiskQuotaUserVtbl
}

type IDiskQuotaUserVtbl struct {
	// IUnknown
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetID                 uintptr
	GetName               uintptr
	GetSidLength          uintptr
	GetSid                uintptr
	GetQuotaThreshold     uintptr
	GetQuotaThresholdText uintptr
	GetQuotaLimit         uintptr
	GetQuotaLimitText     uintptr
	GetQuotaUsed          uintptr
	GetQuotaUsedText      uintptr
	GetQuotaInformation   uintptr
	SetQuotaThreshold     uintptr
	SetQuotaLimit         uintptr
	Invalidate            uintptr
	GetAccountStatus      uintptr
}

func (x *IDiskQuotaUser) AddRef() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.AddRef,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

func (x *IDiskQuotaUser) Release() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.Release,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

func (x *IDiskQuotaUser) GetName(pszAccountContainer *uint16, cchAccountContainer uint32, pszLogonName *uint16, cchLogonName uint32, pszDisplayName *uint16, cchDisplayName uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetName,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pszAccountContainer)),
		uintptr(cchAccountContainer),
		uintptr(unsafe.Pointer(pszLogonName)),
		uintptr(cchLogonName),
		uintptr(unsafe.Pointer(pszDisplayName)),
		uintptr(cchDisplayName),
	)
	return
}

func (x *IDiskQuotaUser) GetQuotaThreshold(pllThreshold *int64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetQuotaThreshold,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pllThreshold)),
	)
	return
}

func (x *IDiskQuotaUser) GetQuotaLimit(pllLimit *int64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetQuotaLimit,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pllLimit)),
	)
	return
}

func (x *IDiskQuotaUser) GetQuotaUsed(pllUsed *int64) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetQuotaUsed,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pllUsed)),
	)
	return
}

func (x *IDiskQuotaUser) SetQuotaThreshold(llThreshold int64, fWriteThrough int32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.SetQuotaThreshold,
		uintptr(unsafe.Pointer(x)),
		uintptr(llThreshold),
		uintptr(fWriteThrough),
	)
	return
}

func (x *IDiskQuotaUser) SetQuotaLimit(llLimit int64, fWriteThrough int32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.SetQuotaLimit,
		uintptr(unsafe.Pointer(x)),
		uintptr(llLimit),
		uintptr(fWriteThrough),
	)
	return
}

func (x *IDiskQuotaUser) GetAccountStatus(pdwStatus *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.GetAccountStatus,
		uintptr(unsafe.Pointer(x)),
		uintptr(unsafe.Pointer(pdwStatus)),
	)
	return
}

type IEnumDiskQuotaUsers struct {
	lpVtbl *IEnumDiskQuotaUsersVtbl
}

type IEnumDiskQuotaUsersVtbl struct {
	// IUnknown
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	Next  uintptr
	Skip  uintptr
	Reset uintptr
	Clone uintptr
}

func (x *IEnumDiskQuotaUsers) AddRef() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.AddRef,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

func (x *IEnumDiskQuotaUsers) Release() (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.Release,
		uintptr(unsafe.Pointer(x)),
	)
	return
}

func (x *IEnumDiskQuotaUsers) Next(cUsers uint32, rgUsers **IDiskQuotaUser, pcUsersFetched *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = syscall.SyscallN(
		x.lpVtbl.Next,
		uintptr(unsafe.Pointer(x)),
		uintptr(cUsers),
		uintptr(unsafe.Pointer(rgUsers)),
		uintptr(unsafe.Pointer(pcUsersFetched)),
	)
	return
}

const DISKQUOTA_STATE_DISABLED = 0x00000000
const DISKQUOTA_STATE_TRACK = 0x00000001
const DISKQUOTA_STATE_ENFORCE = 0x00000002
const DISKQUOTA_STATE_MASK = 0x00000003
const DISKQUOTA_FILESTATE_INCOMPLETE = 0x00000100
const DISKQUOTA_FILESTATE_REBUILDING = 0x00000200
const DISKQUOTA_FILESTATE_MASK = 0x00000300

const DISKQUOTA_USERNAME_RESOLVE_NONE = 0
const DISKQUOTA_USERNAME_RESOLVE_SYNC = 1
const DISKQUOTA_USERNAME_RESOLVE_ASYNC = 2

const DISKQUOTA_USER_ACCOUNT_RESOLVED = 0
const DISKQUOTA_USER_ACCOUNT_UNAVAILABLE = 1
const DISKQUOTA_USER_ACCOUNT_DELETED = 2
const DISKQUOTA_USER_ACCOUNT_INVALID = 3
const DISKQUOTA_USER_ACCOUNT_UNKNOWN = 4
const DISKQUOTA_USER_ACCOUNT_UNRESOLVED = 5

var CLSID_DiskQuotaControl = &comInternalApi.GUID{0x7988b571, 0xec89, 0x11cf, [8]byte{0x9c, 0x0, 0x0, 0xaa, 0x0, 0xa1, 0x4f, 0x56}}
var IID_IDiskQuotaControl = &comInternalApi.GUID{0x7988b572, 0xec89, 0x11cf, [8]byte{0x9c, 0x0, 0x0, 0xaa, 0x0, 0xa1, 0x4f, 0x56}}
var IID_IDiskQuotaUser = &comInternalApi.GUID{0x7988b574, 0xec89, 0x11cf, [8]byte{0x9c, 0x0, 0x0, 0xaa, 0x0, 0xa1, 0x4f, 0x56}}
var IID_IEnumDiskQuotaUsers = &comInternalApi.GUID{0x7988b577, 0xec89, 0x11cf, [8]byte{0x9c, 0x0, 0x0, 0xaa, 0x0, 0xa1, 0x4f, 0x56}}
