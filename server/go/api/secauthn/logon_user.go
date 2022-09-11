//go:build windows && amd64

package secauthn

import (
	"fmt"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	secauthnModelsPb "github.com/s77rt/rdpcloud/proto/go/models/secauthn"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/secure"
	memoryInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/memory"
	secauthnInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthn"
	secauthzInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthz"
)

func LogonUser(user *secauthnModelsPb.User) (string, error) {
	if user == nil {
		return "", status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", status.Errorf(codes.FailedPrecondition, fmt.Sprintf("Unable to get hostname (%s)", err.Error()))
	}

	lpszUsername, err := encode.Windows1252PtrFromString(user.Username)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode username to Windows1252")
	}
	lpszDomain, err := encode.Windows1252PtrFromString(hostname)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode domain to Windows1252")
	}
	lpszPassword, err := encode.Windows1252PtrFromString(user.Password)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode password to Windows1252")
	}

	var ppLogonSid = new(byte)

	ret, _, _ := secauthnInternalApi.LogonUserExA(
		lpszUsername,
		lpszDomain,
		lpszPassword,
		secauthnInternalApi.LOGON32_LOGON_NETWORK,
		secauthnInternalApi.LOGON32_PROVIDER_DEFAULT,
		nil, // phToken not required
		&ppLogonSid,
		nil, // ppProfileBuffer not required
		nil, // pdwProfileLength not required
		nil, // pQuotaLimits not required
	)

	var sidString = new(uint8)

	var ret2 uintptr

	if ret != uintptr(0) {
		ret2, _, _ = secauthzInternalApi.ConvertSidToStringSidA(
			ppLogonSid,
			&sidString,
		)
		memoryInternalApi.LocalFree(ppLogonSid)
	}

	user.Password = ""
	secure.ZeroMemoryUint8FromPtr(lpszPassword)

	if ret == uintptr(0) {
		return "", status.Errorf(codes.Unknown, "Failed to logon user")
	}
	if ret2 == uintptr(0) {
		return "", status.Errorf(codes.Unknown, "Failed to logon user (ConvertSidToStringSidA failure)")
	}

	return encode.Windows1252PtrToString(sidString), nil
}
