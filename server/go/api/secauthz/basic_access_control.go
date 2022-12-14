//go:build windows && amd64

package secauthz

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	memoryInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/memory"
	secauthzInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthz"
)

func LookupAccountSidByUsername(username string) (string, error) {
	if username == "" {
		return "", status.Errorf(codes.InvalidArgument, "Username cannot be empty")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	lpAccountName, err := encode.UTF16PtrFromString(fmt.Sprintf("%s\\%s", hostname, username))
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode account name to UTF16")
	}

	var cbSid uint32
	var cchReferencedDomainName uint32
	var peUse uint32

	ret, _, lastErr := secauthzInternalApi.LookupAccountNameW(
		nil, // start lookup on local
		lpAccountName,
		nil,
		&cbSid,
		nil,
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 && lastErr != windows.ERROR_INSUFFICIENT_BUFFER {
		switch lastErr {
		case windows.ERROR_INVALID_ACCOUNT_NAME:
			return "", status.Errorf(codes.InvalidArgument, "Invalid account name")
		case windows.ERROR_NONE_MAPPED:
			return "", status.Errorf(codes.NotFound, "User not found")
		default:
			return "", status.Errorf(codes.Unknown, "Failed to lookup account SID by name (1) (error: %v)", lastErr)
		}
	}

	var Sid = make([]byte, cbSid)
	var ReferencedDomainName = make([]uint16, cchReferencedDomainName)

	ret, _, lastErr = secauthzInternalApi.LookupAccountNameW(
		nil, // start lookup on local
		lpAccountName,
		&Sid[0],
		&cbSid,
		&ReferencedDomainName[0],
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 {
		switch lastErr {
		case windows.ERROR_INVALID_ACCOUNT_NAME:
			return "", status.Errorf(codes.InvalidArgument, "Invalid account name")
		case windows.ERROR_NONE_MAPPED:
			return "", status.Errorf(codes.NotFound, "User not found")
		default:
			return "", status.Errorf(codes.Unknown, "Failed to lookup account SID by name (2) (error: %v)", lastErr)
		}
	}

	domain := encode.UTF16ToString(ReferencedDomainName)
	if !strings.EqualFold(domain, hostname) {
		return "", status.Errorf(codes.Unknown, "Failed to lookup account SID by name (domain mismatch)")
	}

	var StringSid = new(uint16)

	ret, _, lastErr = secauthzInternalApi.ConvertSidToStringSidW(
		&Sid[0],
		&StringSid,
	)

	if ret == 0 {
		return "", status.Errorf(codes.Unknown, "Failed to lookup account SID by name (ConvertSidToStringSidW) (error: %v)", lastErr)
	}

	sidString := encode.UTF16PtrToString(StringSid)

	memoryInternalApi.LocalFree(uintptr(unsafe.Pointer(StringSid)))
	StringSid = nil

	return sidString, nil
}

func LookupAccountUsernameBySid(sidString string) (string, error) {
	if sidString == "" {
		return "", status.Errorf(codes.InvalidArgument, "SID cannot be empty")
	}

	StringSid, err := encode.UTF16PtrFromString(sidString)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode SID to UTF16")
	}

	var Sid = new(byte)

	ret, _, lastErr := secauthzInternalApi.ConvertStringSidToSidW(
		StringSid,
		&Sid,
	)

	if ret == 0 {
		return "", status.Errorf(codes.Unknown, "Failed to lookup account name by SID (ConvertStringSidToSidW) (error: %v)", lastErr)
	}

	defer func() { memoryInternalApi.LocalFree(uintptr(unsafe.Pointer(Sid))); Sid = nil }()

	var cchName uint32
	var cchReferencedDomainName uint32
	var peUse uint32

	ret, _, lastErr = secauthzInternalApi.LookupAccountSidW(
		nil, // start lookup on local
		Sid,
		nil,
		&cchName,
		nil,
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 && lastErr != windows.ERROR_INSUFFICIENT_BUFFER {
		switch lastErr {
		case windows.ERROR_INVALID_SID:
			return "", status.Errorf(codes.InvalidArgument, "Invalid SID")
		case windows.ERROR_NONE_MAPPED:
			return "", status.Errorf(codes.NotFound, "User not found")
		default:
			return "", status.Errorf(codes.Unknown, "Failed to lookup account name by SID (1) (error: %v)", lastErr)
		}
	}

	var Name = make([]uint16, cchName)
	var ReferencedDomainName = make([]uint16, cchReferencedDomainName)

	ret, _, lastErr = secauthzInternalApi.LookupAccountSidW(
		nil, // start lookup on local
		Sid,
		&Name[0],
		&cchName,
		&ReferencedDomainName[0],
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 {
		switch lastErr {
		case windows.ERROR_INVALID_SID:
			return "", status.Errorf(codes.InvalidArgument, "Invalid SID")
		case windows.ERROR_NONE_MAPPED:
			return "", status.Errorf(codes.NotFound, "User not found")
		default:
			return "", status.Errorf(codes.Unknown, "Failed to lookup account name by SID (2) (error: %v)", lastErr)
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	domain := encode.UTF16ToString(ReferencedDomainName)
	if !strings.EqualFold(domain, hostname) {
		return "", status.Errorf(codes.Unknown, "Failed to lookup account name by SID (domain mismatch)")
	}

	username := encode.UTF16ToString(Name)

	return username, nil
}
