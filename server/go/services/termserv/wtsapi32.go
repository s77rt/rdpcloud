//go:build windows && amd64

package termserv

import (
	"context"

	termservServicePb "github.com/s77rt/rdpcloud/proto/go/services/termserv"
	termservApi "github.com/s77rt/rdpcloud/server/go/api/termserv"
)

func (s *Server) LogoffUser(ctx context.Context, in *termservServicePb.LogoffUserRequest) (*termservServicePb.LogoffUserResponse, error) {
	if err := termservApi.LogoffUser(in.GetUser()); err != nil {
		return nil, err
	}

	return &termservServicePb.LogoffUserResponse{}, nil
}
