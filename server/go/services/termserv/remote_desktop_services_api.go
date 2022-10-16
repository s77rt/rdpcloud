//go:build windows && amd64

package termserv

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	termservModelsPb "github.com/s77rt/rdpcloud/proto/go/models/termserv"
	termservServicePb "github.com/s77rt/rdpcloud/proto/go/services/termserv"
	secauthzApi "github.com/s77rt/rdpcloud/server/go/api/secauthz"
	termservApi "github.com/s77rt/rdpcloud/server/go/api/termserv"
	"github.com/s77rt/rdpcloud/server/go/auth"
)

func (s *Server) LogoffUser(ctx context.Context, in *termservServicePb.LogoffUserRequest) (*termservServicePb.LogoffUserResponse, error) {
	if err := termservApi.LogoffUser(in.GetUser()); err != nil {
		return nil, err
	}

	return &termservServicePb.LogoffUserResponse{}, nil
}

func (s *Server) LogoffMyUser(ctx context.Context, in *termservServicePb.LogoffMyUserRequest) (*termservServicePb.LogoffMyUserResponse, error) {
	userClaims, err := auth.UserClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid context user claims (%s)", err.Error())
	}

	username, err := secauthzApi.LookupAccountUsernameBySid(userClaims.UserSID)
	if err != nil {
		return nil, err
	}

	if err := termservApi.LogoffUser(&termservModelsPb.User_1{Username: username}); err != nil {
		return nil, err
	}

	return &termservServicePb.LogoffMyUserResponse{}, nil
}
