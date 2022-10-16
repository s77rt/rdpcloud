//go:build windows && amd64

package termserv

import (
	"os"
	"unsafe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	termservModelsPb "github.com/s77rt/rdpcloud/proto/go/models/termserv"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	termservInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/termserv"
)

func LogoffUser(user *termservModelsPb.User_1) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	var pLevel uint32 = 1 // This parameter is reserved. Always set this parameter to one
	var ppSessionInfo = &termservInternalApi.WTS_SESSION_INFO_1W{}
	var pCount uint32

	ret, _, _ := termservInternalApi.WTSEnumerateSessionsExW(
		termservInternalApi.WTS_CURRENT_SERVER_HANDLE,
		&pLevel,
		0, // This parameter is reserved. Always set this parameter to zero
		&ppSessionInfo,
		&pCount,
	)
	if ret == 0 {
		return status.Errorf(codes.Unknown, "Failed to logoff user (WTSEnumerateSessionsExW) (error: 0x%x)", ret)
	}

	defer func() {
		termservInternalApi.WTSFreeMemoryExW(
			termservInternalApi.WTSTypeSessionInfoLevel1,
			uintptr(unsafe.Pointer(ppSessionInfo)),
			uint64(pCount),
		)
		ppSessionInfo = nil
	}()

	hostname, err := os.Hostname()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}
	username := user.Username

	var bufDataSample termservInternalApi.WTS_SESSION_INFO_1W
	var SessionIds []uint32

	var bufPtr = unsafe.Pointer(ppSessionInfo)
	for i := uint32(0); i < pCount; i++ {
		var bufData = (*termservInternalApi.WTS_SESSION_INFO_1W)(unsafe.Pointer(uintptr(bufPtr) + uintptr(i)*unsafe.Sizeof(bufDataSample)))
		if bufData.PUserName != nil && bufData.PDomainName != nil {
			if username == encode.UTF16PtrToString(bufData.PUserName) && hostname == encode.UTF16PtrToString(bufData.PDomainName) {
				SessionIds = append(SessionIds, bufData.SessionId)
			}
		}
	}

	if len(SessionIds) == 0 {
		return status.Errorf(codes.NotFound, "User not found / User not logged in")
	}

	for _, SessionId := range SessionIds {
		ret, _, _ := termservInternalApi.WTSLogoffSession(
			termservInternalApi.WTS_CURRENT_SERVER_HANDLE,
			SessionId,
			1, // true => sync
		)
		if ret == 0 {
			return status.Errorf(codes.Unknown, "Failed to logoff user (error: 0x%x)", ret)
		}
	}

	return nil
}
