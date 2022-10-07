//go:build windows && amd64

package sysinfo

import (
	sysinfoInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/sysinfo"
)

func GetUptime() (uint64, error) {
	ret, _, _ := sysinfoInternalApi.GetTickCount64()

	uptime := uint64(ret)

	return uptime, nil
}
