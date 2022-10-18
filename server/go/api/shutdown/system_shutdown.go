//go:build windows && amd64

package shutdown

import (
	"runtime"
	"sync"

	"golang.org/x/sys/windows"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	procthreadInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/procthread"
	secauthzInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthz"
	shutdownInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/shutdown"
	sysinfoInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/sysinfo"
)

func initiateSystemShutdown(wg *sync.WaitGroup, message string, timeout uint32, force bool, reboot bool, reason uint32, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	handle, _, _ := procthreadInternalApi.GetCurrentProcess() // The pseudo handle need not be closed when it is no longer needed

	var tokenHandle uintptr

	ret, _, lastErr := secauthzInternalApi.OpenProcessToken(
		handle,
		secauthzInternalApi.TOKEN_QUERY|secauthzInternalApi.TOKEN_ADJUST_PRIVILEGES,
		&tokenHandle,
	)
	if ret == 0 {
		switch lastErr {
		default:
			*err = status.Errorf(codes.Unknown, "Failed to initiate system shutdown (OpenProcessToken) (error: %v)", lastErr)
		}
		return
	}
	defer func() { sysinfoInternalApi.CloseHandle(tokenHandle); tokenHandle = 0 }()

	tokenPrivileges := &secauthzInternalApi.TOKEN_PRIVILEGES{}

	lpName, rerr := encode.UTF16PtrFromString("SeShutdownPrivilege")
	if rerr != nil {
		*err = status.Errorf(codes.FailedPrecondition, "Unable to encode privilege name to UTF16")
		return
	}

	ret, _, lastErr = secauthzInternalApi.LookupPrivilegeValueW(
		nil, // local
		lpName,
		&tokenPrivileges.Privileges[0].Luid,
	)
	if ret == 0 {
		switch lastErr {
		default:
			*err = status.Errorf(codes.Unknown, "Failed to initiate system shutdown (LookupPrivilegeValueW) (error: %v)", lastErr)
		}
		return
	}

	tokenPrivileges.PrivilegeCount = 1 // one privilege to set
	tokenPrivileges.Privileges[0].Attributes = secauthzInternalApi.SE_PRIVILEGE_ENABLED

	ret, _, lastErr = secauthzInternalApi.AdjustTokenPrivileges(
		tokenHandle,
		0,
		tokenPrivileges,
		0,
		nil,
		nil,
	)
	if ret == 0 {
		switch lastErr {
		default:
			*err = status.Errorf(codes.Unknown, "Failed to initiate system shutdown (AdjustTokenPrivileges) (1) (error: %v)", lastErr)
		}
		return
	}
	if lastErr != windows.ERROR_SUCCESS {
		*err = status.Errorf(codes.Unknown, "Failed to initiate system shutdown (AdjustTokenPrivileges) (2) (error: %v)", lastErr)
		return
	}

	var lpMessage *uint16
	if message != "" {
		var rerr error
		lpMessage, rerr = encode.UTF16PtrFromString(message)
		if rerr != nil {
			*err = status.Errorf(codes.InvalidArgument, "Unable to encode message to UTF16")
			return
		}
	}

	var dwTimeout uint32 = timeout

	var bForceAppsClosed int32
	if force {
		bForceAppsClosed = 1
	}

	var bRebootAfterShutdown int32
	if reboot {
		bRebootAfterShutdown = 1
	}

	var dwReason uint32 = reason

	ret, _, lastErr = shutdownInternalApi.InitiateSystemShutdownExW(
		nil, // local
		lpMessage,
		dwTimeout,
		bForceAppsClosed,
		bRebootAfterShutdown,
		dwReason,
	)

	if ret == 0 {
		switch lastErr {
		case windows.ERROR_SHUTDOWN_IS_SCHEDULED:
			*err = status.Errorf(codes.AlreadyExists, "Shutdown is already scheduled")
		case windows.ERROR_SHUTDOWN_IN_PROGRESS:
			*err = status.Errorf(codes.AlreadyExists, "Shutdown is already in progress")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to initiate system shutdown (error: %v)", lastErr)
		}
		return
	}

	return
}

func InitiateSystemShutdown(message string, timeout uint32, force bool, reboot bool, reason uint32) error {
	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go initiateSystemShutdown(&wg, message, timeout, force, reboot, reason, &err)
	wg.Wait()

	return err
}

func abortSystemShutdown(wg *sync.WaitGroup, err *error) {
	defer wg.Done()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	handle, _, _ := procthreadInternalApi.GetCurrentProcess() // The pseudo handle need not be closed when it is no longer needed

	var tokenHandle uintptr

	ret, _, lastErr := secauthzInternalApi.OpenProcessToken(
		handle,
		secauthzInternalApi.TOKEN_QUERY|secauthzInternalApi.TOKEN_ADJUST_PRIVILEGES,
		&tokenHandle,
	)
	if ret == 0 {
		switch lastErr {
		default:
			*err = status.Errorf(codes.Unknown, "Failed to abort system shutdown (OpenProcessToken) (error: %v)", lastErr)
		}
		return
	}
	defer func() { sysinfoInternalApi.CloseHandle(tokenHandle); tokenHandle = 0 }()

	tokenPrivileges := &secauthzInternalApi.TOKEN_PRIVILEGES{}

	lpName, rerr := encode.UTF16PtrFromString("SeShutdownPrivilege")
	if rerr != nil {
		*err = status.Errorf(codes.FailedPrecondition, "Unable to encode privilege name to UTF16")
		return
	}

	ret, _, lastErr = secauthzInternalApi.LookupPrivilegeValueW(
		nil, // local
		lpName,
		&tokenPrivileges.Privileges[0].Luid,
	)
	if ret == 0 {
		switch lastErr {
		default:
			*err = status.Errorf(codes.Unknown, "Failed to abort system shutdown (LookupPrivilegeValueW) (error: %v)", lastErr)
		}
		return
	}

	tokenPrivileges.PrivilegeCount = 1 // one privilege to set
	tokenPrivileges.Privileges[0].Attributes = secauthzInternalApi.SE_PRIVILEGE_ENABLED

	ret, _, lastErr = secauthzInternalApi.AdjustTokenPrivileges(
		tokenHandle,
		0,
		tokenPrivileges,
		0,
		nil,
		nil,
	)
	if ret == 0 {
		switch lastErr {
		default:
			*err = status.Errorf(codes.Unknown, "Failed to abort system shutdown (AdjustTokenPrivileges) (1) (error: %v)", lastErr)
		}
		return
	}
	if lastErr != windows.ERROR_SUCCESS {
		*err = status.Errorf(codes.Unknown, "Failed to abort system shutdown (AdjustTokenPrivileges) (2) (error: %v)", lastErr)
		return
	}

	ret, _, lastErr = shutdownInternalApi.AbortSystemShutdownW(
		nil, // local
	)

	if ret == 0 {
		switch lastErr {
		case windows.ERROR_NO_SHUTDOWN_IN_PROGRESS:
			*err = status.Errorf(codes.FailedPrecondition, "No shutdown is in progress")
		default:
			*err = status.Errorf(codes.Unknown, "Failed to abort system shutdown (error: %v)", lastErr)
		}
		return
	}

	return
}

func AbortSystemShutdown() error {
	var (
		err error
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go abortSystemShutdown(&wg, &err)
	wg.Wait()

	return err
}
