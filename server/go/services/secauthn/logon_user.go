//go:build windows && amd64

package secauthn

import (
	"context"

	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
	secauthnApi "github.com/s77rt/rdpcloud/server/go/api/secauthn"
)

func (s *Server) LogonUser(ctx context.Context, in *secauthnServicePb.LogonUserRequest) (*secauthnServicePb.LogonUserResponse, error) {
	sidString, err := secauthnApi.LogonUser(in.GetUser())
	if err != nil {
		return nil, err
	}

	token := sidString

	return &secauthnServicePb.LogonUserResponse{
		Token: token,
	}, nil
}
