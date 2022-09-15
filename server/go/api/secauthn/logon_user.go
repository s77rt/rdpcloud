//go:build windows && amd64

package secauthn

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	secauthnModelsPb "github.com/s77rt/rdpcloud/proto/go/models/secauthn"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/secure"
	memoryInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/memory"
	secauthnInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthn"
	secauthzInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthz"
	sysinfoInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/sysinfo"
)

func LogonUser(user *secauthnModelsPb.User_3) (string, error) {
	if user == nil {
		return "", status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Unable to get hostname (%s)", err.Error()))
	}

	lpszUsername, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}
	lpszDomain, err := encode.UTF16PtrFromString(hostname)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode domain to UTF16")
	}
	lpszPassword, err := encode.UTF16PtrFromString(user.Password)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode password to UTF16")
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
			return "", status.Errorf(codes.Unauthenticated, "Login failure")
		}
		return "", status.Errorf(codes.Unknown, fmt.Sprintf("Failed to logon user (error: %d)", lasterr))
	}

	sysinfoInternalApi.CloseHandle(phToken)

	lpAccountName, err := encode.UTF16PtrFromString(fmt.Sprintf("%s\\%s", hostname, user.Username))
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode account name to UTF16")
	}

	var cbSid uint32
	var cchReferencedDomainName uint32
	var peUse uint32

	ret, _, lasterr = secauthzInternalApi.LookupAccountNameW(
		nil, // start lookup on local
		lpAccountName,
		nil,
		&cbSid,
		nil,
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 && lasterr != windows.ERROR_INSUFFICIENT_BUFFER {
		return "", status.Errorf(codes.Unknown, fmt.Sprintf("Failed to logon user [LookupAccountNameW STEP1] (error: %d)", lasterr))
	}

	var Sid = make([]byte, cbSid)
	var ReferencedDomainName = make([]uint16, cchReferencedDomainName)

	ret, _, lasterr = secauthzInternalApi.LookupAccountNameW(
		nil, // start lookup on local
		lpAccountName,
		&Sid[0],
		&cbSid,
		&ReferencedDomainName[0],
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 {
		return "", status.Errorf(codes.Unknown, fmt.Sprintf("Failed to logon user [LookupAccountNameW STEP2] (error: %d)", lasterr))
	}

	domain := encode.UTF16ToString(ReferencedDomainName)
	if domain != hostname {
		return "", status.Errorf(codes.Unknown, "Failed to logon user (domain mismatch)")
	}

	var StringSid = new(uint16)

	ret, _, lasterr = secauthzInternalApi.ConvertSidToStringSidW(
		&Sid[0],
		&StringSid,
	)

	if ret == 0 {
		return "", status.Errorf(codes.Unknown, fmt.Sprintf("Failed to logon user [ConvertSidToStringSidW] (error: %d)", lasterr))
	}

	sidString := encode.UTF16PtrToString(StringSid)

	memoryInternalApi.LocalFree(StringSid)

	return sidString, nil
}
