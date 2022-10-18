//go:build windows && amd64

package shell

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	shellInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/shell"
)

func DeleteProfile(sidString string) error {
	if sidString == "" {
		return status.Errorf(codes.InvalidArgument, "SID cannot be empty")
	}

	lpSidString, err := encode.UTF16PtrFromString(sidString)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode SID to UTF16")
	}

	ret, _, lastErr := shellInternalApi.DeleteProfileW(
		lpSidString,
		nil, // obtain the profile path from the registry
		nil, // local
	)

	if ret != 1 {
		return status.Errorf(codes.Unknown, "Failed to delete profile (error: %v)", lastErr)
	}

	return nil
}
