//go:build windows && amd64

package fileio

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fileioModelsPb "github.com/s77rt/rdpcloud/proto/go/models/fileio"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/win"
	"github.com/s77rt/rdpcloud/server/go/internal/win/win32/com"
	"github.com/s77rt/rdpcloud/server/go/internal/win/win32/fileio"
)

func getQuotaState(wg *sync.WaitGroup, volume string, pdwState *uint32, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get quota state (CoInitializeEx) (error: 0x%x)", ret)
		return
	}
	defer com.CoUninitialize()

	var ppv uintptr
	ret, _, _ = com.CoCreateInstance(
		fileio.CLSID_DiskQuotaControl,
		nil,
		com.CLSCTX_INPROC_SERVER,
		fileio.IID_IDiskQuotaControl,
		&ppv,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get quota state (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volume)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode volume path to UTF16")
		return
	}

	ret, _, _ = diskQuotaControl.Initialize(
		pszPath,
		0, // false => read-only
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.DXERROR_FILE_NOT_FOUND, win.DXERROR_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to get quota state (Initialize) (error: 0x%x)", ret)
		}
		return
	}

	ret, _, _ = diskQuotaControl.GetQuotaState(pdwState)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get quota state (error: 0x%x)", ret)
		return
	}

	return
}

func GetQuotaState(volume string) (uint32, error) {
	if volume == "" {
		return 0, status.Errorf(codes.InvalidArgument, "Volume cannot be empty")
	}

	var (
		quotaState uint32
		err        error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getQuotaState(&wg, volume, &quotaState, &err)
	wg.Wait()

	return quotaState, err
}

func setQuotaState(wg *sync.WaitGroup, volume string, dwState uint32, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set quota state (CoInitializeEx) (error: 0x%x)", ret)
		return
	}
	defer com.CoUninitialize()

	var ppv uintptr
	ret, _, _ = com.CoCreateInstance(
		fileio.CLSID_DiskQuotaControl,
		nil,
		com.CLSCTX_INPROC_SERVER,
		fileio.IID_IDiskQuotaControl,
		&ppv,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set quota state (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volume)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode volume path to UTF16")
		return
	}

	ret, _, _ = diskQuotaControl.Initialize(
		pszPath,
		1, // true => read-write
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.DXERROR_FILE_NOT_FOUND, win.DXERROR_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to set quota state (Initialize) (error: 0x%x)", ret)
		}
		return
	}

	ret, _, _ = diskQuotaControl.SetQuotaState(dwState)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set quota state (error: 0x%x)", ret)
		return
	}

	return
}

func SetQuotaState(volume string, quotaState uint32) error {
	if volume == "" {
		return status.Errorf(codes.InvalidArgument, "Volume cannot be empty")
	}

	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go setQuotaState(&wg, volume, quotaState, &err)
	wg.Wait()

	return err
}

func getDefaultQuota(wg *sync.WaitGroup, volume string, defaultQuota **fileioModelsPb.QuotaEntry_6, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get default quota (CoInitializeEx) (error: 0x%x)", ret)
		return
	}
	defer com.CoUninitialize()

	var ppv uintptr
	ret, _, _ = com.CoCreateInstance(
		fileio.CLSID_DiskQuotaControl,
		nil,
		com.CLSCTX_INPROC_SERVER,
		fileio.IID_IDiskQuotaControl,
		&ppv,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get default quota (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volume)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode volume path to UTF16")
		return
	}

	ret, _, _ = diskQuotaControl.Initialize(
		pszPath,
		0, // false => read-only
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.DXERROR_FILE_NOT_FOUND, win.DXERROR_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to get default quota (Initialize) (error: 0x%x)", ret)
		}
		return
	}

	var pllThreshold int64
	ret, _, _ = diskQuotaControl.GetDefaultQuotaThreshold(&pllThreshold)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get default quota (GetDefaultQuotaThreshold) (error: 0x%x)", ret)
		return
	}

	var pllLimit int64
	ret, _, _ = diskQuotaControl.GetDefaultQuotaLimit(&pllLimit)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get default quota (GetDefaultQuotaLimit) (error: 0x%x)", ret)
		return
	}

	*defaultQuota = &fileioModelsPb.QuotaEntry_6{
		QuotaThreshold: pllThreshold,
		QuotaLimit:     pllLimit,
	}

	return
}

func GetDefaultQuota(volume string) (*fileioModelsPb.QuotaEntry_6, error) {
	if volume == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Volume cannot be empty")
	}

	var (
		defaultQuota *fileioModelsPb.QuotaEntry_6
		err          error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getDefaultQuota(&wg, volume, &defaultQuota, &err)
	wg.Wait()

	return defaultQuota, err
}

func setDefaultQuota(wg *sync.WaitGroup, volume string, defaultQuota *fileioModelsPb.QuotaEntry_6, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set default quota (CoInitializeEx) (error: 0x%x)", ret)
		return
	}
	defer com.CoUninitialize()

	var ppv uintptr
	ret, _, _ = com.CoCreateInstance(
		fileio.CLSID_DiskQuotaControl,
		nil,
		com.CLSCTX_INPROC_SERVER,
		fileio.IID_IDiskQuotaControl,
		&ppv,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set default quota (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volume)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode volume path to UTF16")
		return
	}

	ret, _, _ = diskQuotaControl.Initialize(
		pszPath,
		1, // true => read-wrtie
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.DXERROR_FILE_NOT_FOUND, win.DXERROR_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to set default quota (Initialize) (error: 0x%x)", ret)
		}
		return
	}

	// Unfortunately, the below two methods does not support transactions
	// so this call is not atomic
	// => "Either all of the operations in the transaction succeed or none of the operations persist." DOES NOT apply here
	// SetDefaultQuotaThreshold may succeed but SetDefaultQuotaLimit fail within the same call
	// The good thing is that those two methods are very alike and there is a high probability
	// that they either both fail or both succeed

	ret, _, _ = diskQuotaControl.SetDefaultQuotaThreshold(defaultQuota.QuotaThreshold)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set default quota (SetDefaultQuotaThreshold) (error: 0x%x)", ret)
		return
	}

	ret, _, _ = diskQuotaControl.SetDefaultQuotaLimit(defaultQuota.QuotaLimit)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set default quota (SetDefaultQuotaLimit) (error: 0x%x)", ret)
		return
	}

	return
}

func SetDefaultQuota(volume string, defaultQuota *fileioModelsPb.QuotaEntry_6) error {
	if volume == "" {
		return status.Errorf(codes.InvalidArgument, "Volume cannot be empty")
	}

	if defaultQuota == nil {
		return status.Errorf(codes.InvalidArgument, "DefaultQuota cannot be nil")
	}

	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go setDefaultQuota(&wg, volume, defaultQuota, &err)
	wg.Wait()

	return err
}

func getUsersQuotaEntries(wg *sync.WaitGroup, volume string, quotaEntries *[]*fileioModelsPb.QuotaEntry, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get users quota entries (CoInitializeEx) (error: 0x%x)", ret)
		return
	}
	defer com.CoUninitialize()

	var ppv uintptr
	ret, _, _ = com.CoCreateInstance(
		fileio.CLSID_DiskQuotaControl,
		nil,
		com.CLSCTX_INPROC_SERVER,
		fileio.IID_IDiskQuotaControl,
		&ppv,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get users quota entries (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volume)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode volume path to UTF16")
		return
	}

	ret, _, _ = diskQuotaControl.Initialize(
		pszPath,
		0, // false => read-only
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.DXERROR_FILE_NOT_FOUND, win.DXERROR_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to get users quota entries (Initialize) (error: 0x%x)", ret)
		}
		return
	}

	ppEnum := &fileio.IEnumDiskQuotaUsers{}
	ret, _, _ = diskQuotaControl.CreateEnumUsers(
		nil, // no specific SIDs => all
		0,   // size of the SIDs array
		fileio.DISKQUOTA_USERNAME_RESOLVE_SYNC,
		&ppEnum,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get users quota entries (CreateEnumUsers) (error: 0x%x)", ret)
		return
	}
	defer func() { ppEnum.Release(); ppEnum = nil }()

	hostname, rerr := os.Hostname()
	if rerr != nil {
		*err = status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", rerr.Error())
		return
	}

	for {
		ppUser := &fileio.IDiskQuotaUser{}
		ret, _, _ = ppEnum.Next(
			1, // requesting one user at a time
			&ppUser,
			nil, // not needed
		)
		if ret != win.S_OK {
			break
		}

		var pszLogonName = make([]uint16, win.UNLEN+1)                              // A buffer size of (UNLEN + 1) characters will hold the maximum length user name including the terminating null character
		ppUser.GetName(nil, 0, &pszLogonName[0], uint32(len(pszLogonName)), nil, 0) // AccountContainer, LogonName, DisplayName
		var pllThreshold int64
		ppUser.GetQuotaThreshold(&pllThreshold)
		var pllLimit int64
		ppUser.GetQuotaLimit(&pllLimit)
		var pllUsed int64
		ppUser.GetQuotaUsed(&pllUsed)
		var pdwStatus uint32
		ppUser.GetAccountStatus(&pdwStatus)

		logonName := encode.UTF16PtrToString(&pszLogonName[0])
		logonName_splitted := strings.Split(logonName, "\\")
		if len(logonName_splitted) == 2 {
			domain := logonName_splitted[0]
			if domain == hostname {
				username := logonName_splitted[1]

				var quotaEntry = &fileioModelsPb.QuotaEntry{
					Username:       username,
					QuotaThreshold: pllThreshold,
					QuotaLimit:     pllLimit,
					QuotaUsed:      pllUsed,
					AccountStatus:  pdwStatus,
				}
				*quotaEntries = append(*quotaEntries, quotaEntry)
			}
		}

		ppUser.Release()
		ppUser = nil
	}

	return
}

func GetUsersQuotaEntries(volume string) ([]*fileioModelsPb.QuotaEntry, error) {
	if volume == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Volume cannot be empty")
	}

	var (
		quotaEntries []*fileioModelsPb.QuotaEntry
		err          error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getUsersQuotaEntries(&wg, volume, &quotaEntries, &err)
	wg.Wait()

	return quotaEntries, err
}

func getUserQuotaEntry(wg *sync.WaitGroup, volume string, user *fileioModelsPb.User_1, quotaEntry **fileioModelsPb.QuotaEntry_30, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (CoInitializeEx) (error: 0x%x)", ret)
		return
	}
	defer com.CoUninitialize()

	var ppv uintptr
	ret, _, _ = com.CoCreateInstance(
		fileio.CLSID_DiskQuotaControl,
		nil,
		com.CLSCTX_INPROC_SERVER,
		fileio.IID_IDiskQuotaControl,
		&ppv,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volume)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode volume path to UTF16")
		return
	}

	ret, _, _ = diskQuotaControl.Initialize(
		pszPath,
		0, // false => read-only
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.DXERROR_FILE_NOT_FOUND, win.DXERROR_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (Initialize) (error: 0x%x)", ret)
		}
		return
	}

	hostname, rerr := os.Hostname()
	if rerr != nil {
		*err = status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", rerr.Error())
		return
	}

	pszLogonName, rerr := encode.UTF16PtrFromString(fmt.Sprintf("%s\\%s", hostname, user.Username))
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode logon name to UTF16")
		return
	}

	ppUser := &fileio.IDiskQuotaUser{}
	ret, _, _ = diskQuotaControl.FindUserName(
		pszLogonName,
		&ppUser,
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_NONE_MAPPED, win.DIERR_NOTFOUND:
			*err = status.Errorf(codes.NotFound, "User not found")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (FindUserName) (error: 0x%x)", ret)
		}
		return
	}
	defer func() { ppUser.Release(); ppUser = nil }()

	var pllThreshold int64
	ret, _, _ = ppUser.GetQuotaThreshold(&pllThreshold)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (GetQuotaThreshold) (error: 0x%x)", ret)
		return
	}

	var pllLimit int64
	ret, _, _ = ppUser.GetQuotaLimit(&pllLimit)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (GetQuotaLimit) (error: 0x%x)", ret)
		return
	}

	var pllUsed int64
	ret, _, _ = ppUser.GetQuotaUsed(&pllUsed)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (GetQuotaUsed) (error: 0x%x)", ret)
		return
	}

	var pdwStatus uint32
	ret, _, _ = ppUser.GetAccountStatus(&pdwStatus)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to get user quota entry (GetAccountStatus) (error: 0x%x)", ret)
		return
	}

	*quotaEntry = &fileioModelsPb.QuotaEntry_30{
		QuotaThreshold: pllThreshold,
		QuotaLimit:     pllLimit,
		QuotaUsed:      pllUsed,
		AccountStatus:  pdwStatus,
	}

	return
}

func GetUserQuotaEntry(volume string, user *fileioModelsPb.User_1) (*fileioModelsPb.QuotaEntry_30, error) {
	if volume == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Volume cannot be empty")
	}

	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "User cannot be empty")
	}

	var (
		quotaEntry *fileioModelsPb.QuotaEntry_30
		err        error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getUserQuotaEntry(&wg, volume, user, &quotaEntry, &err)
	wg.Wait()

	return quotaEntry, err
}

func setUserQuotaEntry(wg *sync.WaitGroup, volume string, user *fileioModelsPb.User_1, quotaEntry *fileioModelsPb.QuotaEntry_6, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set user quota entry (CoInitializeEx) (error: 0x%x)", ret)
		return
	}
	defer com.CoUninitialize()

	var ppv uintptr
	ret, _, _ = com.CoCreateInstance(
		fileio.CLSID_DiskQuotaControl,
		nil,
		com.CLSCTX_INPROC_SERVER,
		fileio.IID_IDiskQuotaControl,
		&ppv,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set user quota entry (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volume)
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode volume path to UTF16")
		return
	}

	ret, _, _ = diskQuotaControl.Initialize(
		pszPath,
		1, // true => read-write
	)
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.DXERROR_FILE_NOT_FOUND, win.DXERROR_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to set user quota entry (Initialize) (error: 0x%x)", ret)
		}
		return
	}

	hostname, rerr := os.Hostname()
	if rerr != nil {
		*err = status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", rerr.Error())
		return
	}

	pszLogonName, rerr := encode.UTF16PtrFromString(fmt.Sprintf("%s\\%s", hostname, user.Username))
	if rerr != nil {
		*err = status.Errorf(codes.InvalidArgument, "Unable to encode logon name to UTF16")
		return
	}

	ppUser := &fileio.IDiskQuotaUser{}
	ret, _, _ = diskQuotaControl.AddUserName(
		pszLogonName,
		fileio.DISKQUOTA_USERNAME_RESOLVE_NONE,
		&ppUser,
	)
	if ret != win.S_OK && ret != win.S_FALSE {
		switch ret {
		case win.ERROR_NONE_MAPPED, win.DIERR_NOTFOUND:
			*err = status.Errorf(codes.NotFound, "User not found")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to set user quota entry (AddUserName) (error: 0x%x)", ret)
		}
		return
	}
	defer func() { ppUser.Release(); ppUser = nil }()

	// NOT atomic; read comments on setDefaultQuota()

	ret, _, _ = ppUser.SetQuotaThreshold(quotaEntry.QuotaThreshold, 1)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set user quota entry (SetQuotaThreshold) (error: 0x%x)", ret)
		return
	}

	ret, _, _ = ppUser.SetQuotaLimit(quotaEntry.QuotaLimit, 1)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to set user quota entry (SetQuotaLimit) (error: 0x%x)", ret)
		return
	}

	return
}

func SetUserQuotaEntry(volume string, user *fileioModelsPb.User_1, quotaEntry *fileioModelsPb.QuotaEntry_6) error {
	if volume == "" {
		return status.Errorf(codes.InvalidArgument, "Volume cannot be empty")
	}

	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	if quotaEntry == nil {
		return status.Errorf(codes.InvalidArgument, "QuotaEntry cannot be nil")
	}

	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go setUserQuotaEntry(&wg, volume, user, quotaEntry, &err)
	wg.Wait()

	return err
}
