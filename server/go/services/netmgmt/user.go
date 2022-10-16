//go:build windows && amd64

package netmgmt

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	netmgmtModelsPb "github.com/s77rt/rdpcloud/proto/go/models/netmgmt"
	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	netmgmtApi "github.com/s77rt/rdpcloud/server/go/api/netmgmt"
	secauthzApi "github.com/s77rt/rdpcloud/server/go/api/secauthz"
	"github.com/s77rt/rdpcloud/server/go/auth"
)

func (s *Server) AddUser(ctx context.Context, in *netmgmtServicePb.AddUserRequest) (*netmgmtServicePb.AddUserResponse, error) {
	if err := netmgmtApi.AddUser(in.GetUser()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.AddUserResponse{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *netmgmtServicePb.DeleteUserRequest) (*netmgmtServicePb.DeleteUserResponse, error) {
	if err := netmgmtApi.DeleteUser(in.GetUser()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.DeleteUserResponse{}, nil
}

func (s *Server) GetUsers(ctx context.Context, in *netmgmtServicePb.GetUsersRequest) (*netmgmtServicePb.GetUsersResponse, error) {
	users, err := netmgmtApi.GetUsers()
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetUsersResponse{
		Users: users,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, in *netmgmtServicePb.GetUserRequest) (*netmgmtServicePb.GetUserResponse, error) {
	user, err := netmgmtApi.GetUser(in.GetUser())
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetUserResponse{
		User: user,
	}, nil
}

func (s *Server) GetMyUser(ctx context.Context, in *netmgmtServicePb.GetMyUserRequest) (*netmgmtServicePb.GetMyUserResponse, error) {
	userClaims, err := auth.UserClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid context user claims (%s)", err.Error())
	}

	username, err := secauthzApi.LookupAccountUsernameBySid(userClaims.UserSID)
	if err != nil {
		return nil, err
	}

	user, err := netmgmtApi.GetUser(&netmgmtModelsPb.User_1{Username: username})
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetMyUserResponse{
		User: user,
	}, nil
}

func (s *Server) GetUserLocalGroups(ctx context.Context, in *netmgmtServicePb.GetUserLocalGroupsRequest) (*netmgmtServicePb.GetUserLocalGroupsResponse, error) {
	localGroups, err := netmgmtApi.GetUserLocalGroups(in.GetUser())
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetUserLocalGroupsResponse{
		LocalGroups: localGroups,
	}, nil
}

func (s *Server) GetMyUserLocalGroups(ctx context.Context, in *netmgmtServicePb.GetMyUserLocalGroupsRequest) (*netmgmtServicePb.GetMyUserLocalGroupsResponse, error) {
	userClaims, err := auth.UserClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid context user claims (%s)", err.Error())
	}

	username, err := secauthzApi.LookupAccountUsernameBySid(userClaims.UserSID)
	if err != nil {
		return nil, err
	}

	localGroups, err := netmgmtApi.GetUserLocalGroups(&netmgmtModelsPb.User_1{Username: username})
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetMyUserLocalGroupsResponse{
		LocalGroups: localGroups,
	}, nil
}

func (s *Server) ChangeUserPassword(ctx context.Context, in *netmgmtServicePb.ChangeUserPasswordRequest) (*netmgmtServicePb.ChangeUserPasswordResponse, error) {
	if err := netmgmtApi.ChangeUserPassword(in.GetUser()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.ChangeUserPasswordResponse{}, nil
}

func (s *Server) ChangeMyUserPassword(ctx context.Context, in *netmgmtServicePb.ChangeMyUserPasswordRequest) (*netmgmtServicePb.ChangeMyUserPasswordResponse, error) {
	userClaims, err := auth.UserClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid context user claims (%s)", err.Error())
	}

	username, err := secauthzApi.LookupAccountUsernameBySid(userClaims.UserSID)
	if err != nil {
		return nil, err
	}

	user := &netmgmtModelsPb.User_3{
		Username: username,
		Password: in.GetUser().GetPassword(),
	}

	if err := netmgmtApi.ChangeUserPassword(user); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.ChangeMyUserPasswordResponse{}, nil
}

func (s *Server) EnableUser(ctx context.Context, in *netmgmtServicePb.EnableUserRequest) (*netmgmtServicePb.EnableUserResponse, error) {
	if err := netmgmtApi.EnableUser(in.GetUser()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.EnableUserResponse{}, nil
}

func (s *Server) DisableUser(ctx context.Context, in *netmgmtServicePb.DisableUserRequest) (*netmgmtServicePb.DisableUserResponse, error) {
	if err := netmgmtApi.DisableUser(in.GetUser()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.DisableUserResponse{}, nil
}
