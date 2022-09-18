//go:build windows && amd64

package secauthz

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	secauthzModelsPb "github.com/s77rt/rdpcloud/proto/go/models/secauthz"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	memoryInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/memory"
	secauthzInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/secauthz"
)

// LookupAccountByName lookup for an account by it's name and returns the corresponding SID
func LookupAccountByName(user *secauthzModelsPb.User_1) (string, error) {
	if user == nil {
		return "", status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	lpAccountName, err := encode.UTF16PtrFromString(fmt.Sprintf("%s\\%s", hostname, user.Username))
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Unable to encode account name to UTF16")
	}

	var cbSid uint32
	var cchReferencedDomainName uint32
	var peUse uint32

	ret, _, lasterr := secauthzInternalApi.LookupAccountNameW(
		nil, // start lookup on local
		lpAccountName,
		nil,
		&cbSid,
		nil,
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 && lasterr != windows.ERROR_INSUFFICIENT_BUFFER {
		return "", status.Errorf(codes.Unknown, "Failed to lookup account by name (1) (error: %d)", lasterr)
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
		return "", status.Errorf(codes.Unknown, "Failed to lookup account by name (2) (error: %d)", lasterr)
	}

	domain := encode.UTF16ToString(ReferencedDomainName)
	if domain != hostname {
		return "", status.Errorf(codes.Unknown, "Failed to lookup account by name (domain mismatch)")
	}

	var StringSid = new(uint16)

	ret, _, lasterr = secauthzInternalApi.ConvertSidToStringSidW(
		&Sid[0],
		&StringSid,
	)

	if ret == 0 {
		return "", status.Errorf(codes.Unknown, "Failed to lookup account by name (ConvertSidToStringSidW) (error: %d)", lasterr)
	}

	sidString := encode.UTF16PtrToString(StringSid)

	memoryInternalApi.LocalFree(StringSid)

	return sidString, nil
}

// LookupAccountBySid lookup for an account by it's SID and returns the corresponding name
func LookupAccountBySid(sidString string) (*secauthzModelsPb.User_1, error) {
	if sidString == "" {
		return nil, status.Errorf(codes.InvalidArgument, "SID cannot be empty")
	}

	StringSid, err := encode.UTF16PtrFromString(sidString)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to encode SID to UTF16")
	}

	var Sid = new(byte)

	ret, _, lasterr := secauthzInternalApi.ConvertStringSidToSidW(
		StringSid,
		&Sid,
	)

	if ret == 0 {
		return nil, status.Errorf(codes.Unknown, "Failed to lookup account by SID (ConvertStringSidToSidW) (error: %d)", lasterr)
	}

	defer memoryInternalApi.LocalFree(Sid)

	var cchName uint32
	var cchReferencedDomainName uint32
	var peUse uint32

	ret, _, lasterr = secauthzInternalApi.LookupAccountSidW(
		nil, // start lookup on local
		Sid,
		nil,
		&cchName,
		nil,
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 && lasterr != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, status.Errorf(codes.Unknown, "Failed to lookup account by SID (1) (error: %d)", lasterr)
	}

	var Name = make([]uint16, cchName)
	var ReferencedDomainName = make([]uint16, cchReferencedDomainName)

	ret, _, lasterr = secauthzInternalApi.LookupAccountSidW(
		nil, // start lookup on local
		Sid,
		&Name[0],
		&cchName,
		&ReferencedDomainName[0],
		&cchReferencedDomainName,
		&peUse,
	)

	if ret == 0 {
		return nil, status.Errorf(codes.Unknown, "Failed to lookup account by SID (2) (error: %d)", lasterr)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	domain := encode.UTF16ToString(ReferencedDomainName)
	if domain != hostname {
		return nil, status.Errorf(codes.Unknown, "Failed to lookup account by SID (domain mismatch)")
	}

	user := &secauthzModelsPb.User_1{
		Username: encode.UTF16PtrToString(&Name[0]),
	}

	return user, nil
}
