//go:build windows && amd64

package secauthn

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	netmgmtModelsPb "github.com/s77rt/rdpcloud/proto/go/models/netmgmt"
	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
	netmgmtApi "github.com/s77rt/rdpcloud/server/go/api/netmgmt"
	secauthnApi "github.com/s77rt/rdpcloud/server/go/api/secauthn"
	secauthzApi "github.com/s77rt/rdpcloud/server/go/api/secauthz"

	"github.com/s77rt/rdpcloud/server/go/auth"
)

func (s *Server) LogonUser(ctx context.Context, in *secauthnServicePb.LogonUserRequest) (*secauthnServicePb.LogonUserResponse, error) {
	user := in.GetUser()

	if err := secauthnApi.LogonUser(user); err != nil {
		return nil, err
	}

	sidString, err := secauthzApi.LookupAccountSidByUsername(user.GetUsername())
	if err != nil {
		return nil, err
	}

	fetchedUser, err := netmgmtApi.GetUser(&netmgmtModelsPb.User_1{Username: user.GetUsername()})
	if err != nil {
		return nil, err
	}

	token := auth.NewTokenWithUserClaims(auth.UserClaims{
		UserSID:           sidString, // Subject
		PreferredUsername: fetchedUser.GetUsername(),
		Privilege:         fetchedUser.GetPrivilege(),
	})

	tokenString, err := auth.TokenSignedString(token)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Unable to sign the token (%s)", err.Error()))
	}

	return &secauthnServicePb.LogonUserResponse{
		Token: tokenString,
	}, nil
}
