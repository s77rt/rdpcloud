//go:build windows && amd64

package secauthz

import (
	"context"

	secauthzServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthz"
	secauthzApi "github.com/s77rt/rdpcloud/server/go/api/secauthz"
)

func (s *Server) LookupAccountSidByUsername(ctx context.Context, in *secauthzServicePb.LookupAccountSidByUsernameRequest) (*secauthzServicePb.LookupAccountSidByUsernameResponse, error) {
	sidString, err := secauthzApi.LookupAccountSidByUsername(in.GetUsername())
	if err != nil {
		return nil, err
	}

	return &secauthzServicePb.LookupAccountSidByUsernameResponse{
		Sid: sidString,
	}, nil
}

func (s *Server) LookupAccountUsernameBySid(ctx context.Context, in *secauthzServicePb.LookupAccountUsernameBySidRequest) (*secauthzServicePb.LookupAccountUsernameBySidResponse, error) {
	username, err := secauthzApi.LookupAccountUsernameBySid(in.GetSid())
	if err != nil {
		return nil, err
	}

	return &secauthzServicePb.LookupAccountUsernameBySidResponse{
		Username: username,
	}, nil
}
