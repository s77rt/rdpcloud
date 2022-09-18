//go:build windows && amd64

package secauthz

import (
	"context"

	secauthzServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthz"
	secauthzApi "github.com/s77rt/rdpcloud/server/go/api/secauthz"
)

func (s *Server) LookupAccountByName(ctx context.Context, in *secauthzServicePb.LookupAccountByNameRequest) (*secauthzServicePb.LookupAccountByNameResponse, error) {
	sidString, err := secauthzApi.LookupAccountByName(in.GetUser())
	if err != nil {
		return nil, err
	}

	return &secauthzServicePb.LookupAccountByNameResponse{
		Sid: sidString,
	}, nil
}

func (s *Server) LookupAccountBySid(ctx context.Context, in *secauthzServicePb.LookupAccountBySidRequest) (*secauthzServicePb.LookupAccountBySidResponse, error) {
	user, err := secauthzApi.LookupAccountBySid(in.GetSid())
	if err != nil {
		return nil, err
	}

	return &secauthzServicePb.LookupAccountBySidResponse{
		User: user,
	}, nil
}
