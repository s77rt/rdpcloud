//go:build windows && amd64

package secauthn

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang-jwt/jwt/v4"

	netmgmtModelsPb "github.com/s77rt/rdpcloud/proto/go/models/netmgmt"
	secauthzModelsPb "github.com/s77rt/rdpcloud/proto/go/models/secauthz"
	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
	netmgmtApi "github.com/s77rt/rdpcloud/server/go/api/netmgmt"
	secauthnApi "github.com/s77rt/rdpcloud/server/go/api/secauthn"
	secauthzApi "github.com/s77rt/rdpcloud/server/go/api/secauthz"
	"github.com/s77rt/rdpcloud/server/go/config"
	customJWT "github.com/s77rt/rdpcloud/server/go/jwt"
)

func (s *Server) LogonUser(ctx context.Context, in *secauthnServicePb.LogonUserRequest) (*secauthnServicePb.LogonUserResponse, error) {
	user := in.GetUser()

	if err := secauthnApi.LogonUser(user); err != nil {
		return nil, err
	}

	sidString, err := secauthzApi.LookupAccountByName(&secauthzModelsPb.User_1{Username: user.GetUsername()})
	if err != nil {
		return nil, err
	}

	fetchedUser, err := netmgmtApi.GetUser(&netmgmtModelsPb.User_1{Username: user.GetUsername()})
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customJWT.UserClaims{
		PreferredUsername: fetchedUser.GetUsername(),
		Privilege:         fetchedUser.GetPrivilege(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sidString,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenLifetime * time.Second)),
		},
	})

	tokenString, err := token.SignedString(config.Secret)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Unable to sign the token (%s)", err.Error()))
	}

	return &secauthnServicePb.LogonUserResponse{
		Token: tokenString,
	}, nil
}
