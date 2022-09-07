//go:build windows && amd64

package netmgmt

import (
	"context"

	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	netmgmtApi "github.com/s77rt/rdpcloud/server/go/api/netmgmt"
)

func (s *Server) GetUsers(ctx context.Context, in *netmgmtServicePb.GetUsersRequest) (*netmgmtServicePb.GetUsersResponse, error) {
	users, err := netmgmtApi.GetUsers()
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetUsersResponse{
		Users: users,
	}, nil
}

func (s *Server) AddUser(ctx context.Context, in *netmgmtServicePb.AddUserRequest) (*netmgmtServicePb.AddUserResponse, error) {
	if err := netmgmtApi.AddUser(in.GetUser()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.AddUserResponse{}, nil
}
