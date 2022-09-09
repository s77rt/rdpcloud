//go:build windows && amd64

package netmgmt

import (
	"context"

	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	netmgmtApi "github.com/s77rt/rdpcloud/server/go/api/netmgmt"
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

func (s *Server) GetUserLocalGroups(ctx context.Context, in *netmgmtServicePb.GetUserLocalGroupsRequest) (*netmgmtServicePb.GetUserLocalGroupsResponse, error) {
	localGroups, err := netmgmtApi.GetUserLocalGroups(in.GetUser())
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetUserLocalGroupsResponse{
		LocalGroups: localGroups,
	}, nil
}

func (s *Server) ChangeUserPassword(ctx context.Context, in *netmgmtServicePb.ChangeUserPasswordRequest) (*netmgmtServicePb.ChangeUserPasswordResponse, error) {
	if err := netmgmtApi.ChangeUserPassword(in.GetUser()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.ChangeUserPasswordResponse{}, nil
}
