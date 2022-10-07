//go:build windows && amd64

package sysinfo

import (
	"context"

	sysinfoServicePb "github.com/s77rt/rdpcloud/proto/go/services/sysinfo"
	sysinfoApi "github.com/s77rt/rdpcloud/server/go/api/sysinfo"
)

func (s *Server) GetUptime(ctx context.Context, in *sysinfoServicePb.GetUptimeRequest) (*sysinfoServicePb.GetUptimeResponse, error) {
	uptime, err := sysinfoApi.GetUptime()
	if err != nil {
		return nil, err
	}

	return &sysinfoServicePb.GetUptimeResponse{
		Uptime: uptime,
	}, nil
}
