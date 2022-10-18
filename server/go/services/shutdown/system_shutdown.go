//go:build windows && amd64

package shutdown

import (
	"context"

	shutdownServicePb "github.com/s77rt/rdpcloud/proto/go/services/shutdown"
	shutdownApi "github.com/s77rt/rdpcloud/server/go/api/shutdown"
)

func (s *Server) InitiateSystemShutdown(ctx context.Context, in *shutdownServicePb.InitiateSystemShutdownRequest) (*shutdownServicePb.InitiateSystemShutdownResponse, error) {
	if err := shutdownApi.InitiateSystemShutdown(in.GetMessage(), in.GetTimeout(), in.GetForce(), in.GetReboot(), in.GetReason()); err != nil {
		return nil, err
	}

	return &shutdownServicePb.InitiateSystemShutdownResponse{}, nil
}

func (s *Server) AbortSystemShutdown(ctx context.Context, in *shutdownServicePb.AbortSystemShutdownRequest) (*shutdownServicePb.AbortSystemShutdownResponse, error) {
	if err := shutdownApi.AbortSystemShutdown(); err != nil {
		return nil, err
	}

	return &shutdownServicePb.AbortSystemShutdownResponse{}, nil
}
