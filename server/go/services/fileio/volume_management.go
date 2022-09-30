//go:build windows && amd64

package fileio

import (
	"context"

	fileioServicePb "github.com/s77rt/rdpcloud/proto/go/services/fileio"
	fileioApi "github.com/s77rt/rdpcloud/server/go/api/fileio"
)

func (s *Server) GetVolumes(ctx context.Context, in *fileioServicePb.GetVolumesRequest) (*fileioServicePb.GetVolumesResponse, error) {
	volumes, err := fileioApi.GetVolumes()
	if err != nil {
		return nil, err
	}

	return &fileioServicePb.GetVolumesResponse{
		Volumes: volumes,
	}, nil
}
