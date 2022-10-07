//go:build windows && amd64

package shell

import (
	"context"

	shellServicePb "github.com/s77rt/rdpcloud/proto/go/services/shell"
	shellApi "github.com/s77rt/rdpcloud/server/go/api/shell"
)

type Server struct {
	shellServicePb.UnimplementedShellServer
}

func (s *Server) DeleteProfile(ctx context.Context, in *shellServicePb.DeleteProfileRequest) (*shellServicePb.DeleteProfileResponse, error) {
	if err := shellApi.DeleteProfile(in.GetSid()); err != nil {
		return nil, err
	}

	return &shellServicePb.DeleteProfileResponse{}, nil
}
