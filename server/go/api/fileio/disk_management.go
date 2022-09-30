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

func getQuotaState(wg *sync.WaitGroup, volumePath string, pdwState *uint32, err *error) {
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

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
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

func GetQuotaState(volumePath string) (uint32, error) {
	if volumePath == "" {
		return 0, status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
	}

	var (
		quotaState uint32
		err        error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getQuotaState(&wg, volumePath, &quotaState, &err)
	wg.Wait()

	return quotaState, err
}

func setQuotaState(wg *sync.WaitGroup, volumePath string, dwState uint32, err *error) {
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

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
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

func SetQuotaState(volumePath string, quotaState uint32) error {
	if volumePath == "" {
		return status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
	}

	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go setQuotaState(&wg, volumePath, quotaState, &err)
	wg.Wait()

	return err
}

func getDefaultQuota(wg *sync.WaitGroup, volumePath string, defaultQuota **fileioModelsPb.QuotaEntry_6, err *error) {
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

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
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

func GetDefaultQuota(volumePath string) (*fileioModelsPb.QuotaEntry_6, error) {
	if volumePath == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
	}

	var (
		defaultQuota *fileioModelsPb.QuotaEntry_6
		err          error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getDefaultQuota(&wg, volumePath, &defaultQuota, &err)
	wg.Wait()

	return defaultQuota, err
}

func setDefaultQuota(wg *sync.WaitGroup, volumePath string, defaultQuota *fileioModelsPb.QuotaEntry_6, err *error) {
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

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
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

func SetDefaultQuota(volumePath string, defaultQuota *fileioModelsPb.QuotaEntry_6) error {
	if volumePath == "" {
		return status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
	}

	if defaultQuota == nil {
		return status.Errorf(codes.InvalidArgument, "DefaultQuota cannot be nil")
	}

	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go setDefaultQuota(&wg, volumePath, defaultQuota, &err)
	wg.Wait()

	return err
}

func getUsersQuotaEntries(wg *sync.WaitGroup, volumePath string, quotaEntries *[]*fileioModelsPb.QuotaEntry, err *error) {
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

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
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

func GetUsersQuotaEntries(volumePath string) ([]*fileioModelsPb.QuotaEntry, error) {
	if volumePath == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
	}

	var (
		quotaEntries []*fileioModelsPb.QuotaEntry
		err          error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getUsersQuotaEntries(&wg, volumePath, &quotaEntries, &err)
	wg.Wait()

	return quotaEntries, err
}

func getUserQuotaEntry(wg *sync.WaitGroup, volumePath string, user *fileioModelsPb.User_1, quotaEntry **fileioModelsPb.QuotaEntry_30, err *error) {
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

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
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
		case win.ERROR_NONE_MAPPED, win.E_FILE_NOT_FOUND:
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

func GetUserQuotaEntry(volumePath string, user *fileioModelsPb.User_1) (*fileioModelsPb.QuotaEntry_30, error) {
	if volumePath == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
	}

	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	var (
		quotaEntry *fileioModelsPb.QuotaEntry_30
		err        error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go getUserQuotaEntry(&wg, volumePath, user, &quotaEntry, &err)
	wg.Wait()

	return quotaEntry, err
}

func setUserQuotaEntry(wg *sync.WaitGroup, volumePath string, user *fileioModelsPb.User_1, quotaEntry *fileioModelsPb.QuotaEntry_6, err *error) {
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

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
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
		case win.ERROR_NONE_MAPPED, win.E_FILE_NOT_FOUND:
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

func SetUserQuotaEntry(volumePath string, user *fileioModelsPb.User_1, quotaEntry *fileioModelsPb.QuotaEntry_6) error {
	if volumePath == "" {
		return status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
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
	go setUserQuotaEntry(&wg, volumePath, user, quotaEntry, &err)
	wg.Wait()

	return err
}

func deleteUserQuotaEntry(wg *sync.WaitGroup, volumePath string, user *fileioModelsPb.User_1, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := com.CoInitializeEx(
		0, // must be NULL
		com.COINIT_APARTMENTTHREADED,
	)
	if ret != win.S_OK {
		*err = status.Errorf(codes.Unknown, "Failed to delete user quota entry (CoInitializeEx) (error: 0x%x)", ret)
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
		*err = status.Errorf(codes.Unknown, "Failed to delete user quota entry (CoCreateInstance) (error: 0x%x)", ret)
		return
	}

	diskQuotaControl := (*fileio.IDiskQuotaControl)(unsafe.Pointer(ppv))
	defer func() { diskQuotaControl.Release(); diskQuotaControl = nil }()

	pszPath, rerr := encode.UTF16PtrFromString(volumePath)
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
		case win.ERROR_BAD_PATHNAME, win.ERROR_INVALID_NAME, win.E_PATH_TOO_LONG, win.E_INVALID_NAME:
			*err = status.Errorf(codes.InvalidArgument, "Bad volume path")
		case win.ERROR_FILE_NOT_FOUND, win.ERROR_PATH_NOT_FOUND, win.E_FILE_NOT_FOUND, win.E_PATH_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "Volume not found")
		case win.ERROR_NOT_SUPPORTED, win.E_NOT_SUPPORTED:
			*err = status.Errorf(codes.FailedPrecondition, "Volume does not support quotas")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to delete user quota entry (Initialize) (error: 0x%x)", ret)
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

	// To delete a quota entry we need a IDiskQuotaUser interface (ppUser)
	// and we get that by first adding the user to the quota entries using AddUserName
	//
	// Note: FindUserName cannot be used, although it does provide us the IDiskQuotaUser interface (ppUser)
	// calling DeleteUser with that ppUser will fail if the user never had a record in that volume
	// and that is not the desired behaviour and it's kinda confusing:
	// 1: FindUserName->DeleteUser (User have a quota entry in the volume) => OK
	// 2: FindUserName->DeleteUser (User had a quota entry in the volume but not anymore) => OK
	// 3: FindUserName->DeleteUser (User never had a quota entry in the volume) => NOK
	// in both cases 2 and 3, the record does not actually exist, but in the former it's ok and in the latter it's not
	// Using AddUserName clear up this confusing:
	// 1: AddUserName->DeleteUser (User have a quota entry in the volume) => OK
	// 2: AddUserName->DeleteUser (User had a quota entry in the volume but not anymore) => OK
	// 3: AddUserName->DeleteUser (User never had a quota entry in the volume) => OK

	ppUser := &fileio.IDiskQuotaUser{}
	ret, _, _ = diskQuotaControl.AddUserName(
		pszLogonName,
		fileio.DISKQUOTA_USERNAME_RESOLVE_NONE,
		&ppUser,
	)
	if ret != win.S_OK && ret != win.S_FALSE {
		switch ret {
		case win.ERROR_NONE_MAPPED, win.E_FILE_NOT_FOUND:
			*err = status.Errorf(codes.NotFound, "User not found")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to delete user quota entry (AddUserName) (error: 0x%x)", ret)
		}
		return
	}
	defer func() { ppUser.Release(); ppUser = nil }()

	ppUser.AddRef() // Increment the internal counter before handing the object to another function
	ret, _, _ = diskQuotaControl.DeleteUser(ppUser)
	ppUser.Release() // As per the Microsoft docs, DeleteUser does not release the object thus we have to release it manually
	if ret != win.S_OK {
		switch ret {
		case win.ERROR_FILE_EXISTS, win.E_FILE_EXISTS:
			*err = status.Errorf(codes.FailedPrecondition, "User owns files on the volume")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to delete user quota entry (error: 0x%x)", ret)
		}
		return
	}

	return
}

func DeleteUserQuotaEntry(volumePath string, user *fileioModelsPb.User_1) error {
	if volumePath == "" {
		return status.Errorf(codes.InvalidArgument, "Volume path cannot be empty")
	}

	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go deleteUserQuotaEntry(&wg, volumePath, user, &err)
	wg.Wait()

	return err
}
