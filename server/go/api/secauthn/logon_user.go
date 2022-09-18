//go:build windows && amd64

package secauthn

import (
	"os"

	"golang.org/x/sys/windows"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	secauthnModelsPb "github.com/s77rt/rdpcloud/proto/go/models/secauthn"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/secure"
	secauthnInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthn"
	sysinfoInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/sysinfo"
)

func LogonUser(user *secauthnModelsPb.User_3) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	lpszUsername, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}
	lpszDomain, err := encode.UTF16PtrFromString(hostname)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode domain to UTF16")
	}
	lpszPassword, err := encode.UTF16PtrFromString(user.Password)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode password to UTF16")
	}

	var phToken uintptr

	ret, _, lasterr := secauthnInternalApi.LogonUserW(
		lpszUsername,
		lpszDomain,
		lpszPassword,
		secauthnInternalApi.LOGON32_LOGON_NETWORK,
		secauthnInternalApi.LOGON32_PROVIDER_DEFAULT,
		&phToken,
	)

	user.Password = ""
	secure.ZeroMemoryUint16FromPtr(lpszPassword)

	if ret == 0 {
		if lasterr == windows.ERROR_LOGON_FAILURE {
			return status.Errorf(codes.Unauthenticated, "Login failure")
		}
		return status.Errorf(codes.Unknown, "Failed to logon user (error: %d)", lasterr)
	}

	sysinfoInternalApi.CloseHandle(phToken)

	return nil
}
