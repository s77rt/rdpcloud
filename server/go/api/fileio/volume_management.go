//go:build windows && amd64

package fileio

import (
	"fmt"

	"golang.org/x/sys/windows"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fileioModelsPb "github.com/s77rt/rdpcloud/proto/go/models/fileio"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/win"
	"github.com/s77rt/rdpcloud/server/go/internal/win/win32/fileio"
)

func getVolumesGuids() ([]string, error) {
	volumesGuids := []string{}

	lpszVolumeName := make([]uint16, 49+1) // expected GUID format is "\\?\Volume{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}\" of length 49, +1 for the terminating null character
	cchBufferLength := uint32(len(lpszVolumeName))

	hFindVolume, _, lastErr := fileio.FindFirstVolumeW(
		&lpszVolumeName[0],
		cchBufferLength,
	)
	if hFindVolume == win.INVALID_HANDLE_VALUE {
		return nil, status.Errorf(codes.Unknown, "Failed to get volumes guids (FindFirstVolumeW) (error: %v)", lastErr)
	}
	defer func() { fileio.FindVolumeClose(hFindVolume); hFindVolume = 0 }()

	volumesGuids = append(volumesGuids, encode.UTF16ToString(lpszVolumeName)[11:47])

	if lastErr == windows.ERROR_MORE_DATA {
		for {
			lpszVolumeName = make([]uint16, 49+1)
			cchBufferLength = uint32(len(lpszVolumeName))

			ret, _, lastErr := fileio.FindNextVolumeW(
				hFindVolume,
				&lpszVolumeName[0],
				cchBufferLength,
			)
			if ret == 0 {
				if lastErr != windows.ERROR_NO_MORE_FILES {
					return nil, status.Errorf(codes.Unknown, "Failed to get volumes guids (FindNextVolumeW) (error: %v)", lastErr)
				}
				break
			}

			volumesGuids = append(volumesGuids, encode.UTF16ToString(lpszVolumeName)[11:47])
		}
	}

	return volumesGuids, nil
}

func getVolumePaths(volumeGuid string) ([]string, error) {
	volumePaths := []string{}

	lpszVolumeName, err := encode.UTF16PtrFromString(fmt.Sprintf("\\\\?\\Volume{%s}\\", volumeGuid))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to encode volume guid to UTF16")
	}

	var lpszVolumePathNames = make([]uint16, 1)
	var cchBufferLength uint32 = uint32(len(lpszVolumePathNames))
	var lpcchReturnLength uint32

	ret, _, lastErr := fileio.GetVolumePathNamesForVolumeNameW(
		lpszVolumeName,
		&lpszVolumePathNames[0],
		cchBufferLength,
		&lpcchReturnLength,
	)
	if ret == 0 && lastErr != windows.ERROR_MORE_DATA {
		return nil, status.Errorf(codes.Unknown, "Failed to get volume paths (1) (error: %v)", lastErr)
	} else if ret != 0 {
		// No assigned paths
		return volumePaths, nil
	}

	lpszVolumePathNames = make([]uint16, lpcchReturnLength)
	cchBufferLength = uint32(len(lpszVolumePathNames))

	ret, _, lastErr = fileio.GetVolumePathNamesForVolumeNameW(
		lpszVolumeName,
		&lpszVolumePathNames[0],
		cchBufferLength,
		&lpcchReturnLength,
	)
	if ret == 0 {
		return nil, status.Errorf(codes.Unknown, "Failed to get volume paths (2) (error: %v)", lastErr)
	}

	offset := 0
	for i, v := range lpszVolumePathNames {
		if v == 0 {
			volumePath := encode.UTF16PtrToString(&lpszVolumePathNames[offset])
			if volumePath != "" {
				volumePaths = append(volumePaths, volumePath)
			}
			offset = i + 1
		}
	}

	return volumePaths, nil
}

func GetVolumes() ([]*fileioModelsPb.Volume, error) {
	volumes := []*fileioModelsPb.Volume{}

	volumesGuids, err := getVolumesGuids()
	if err != nil {
		return nil, err
	}

	for _, volumeGuid := range volumesGuids {
		volumePaths, err := getVolumePaths(volumeGuid)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, &fileioModelsPb.Volume{
			Guid:  volumeGuid,
			Paths: volumePaths,
		})
	}

	return volumes, nil
}
