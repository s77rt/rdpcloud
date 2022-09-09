//go:build windows && amd64

package netmgmt

import (
	"context"

	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	netmgmtApi "github.com/s77rt/rdpcloud/server/go/api/netmgmt"
)

func (s *Server) AddUserToLocalGroup(ctx context.Context, in *netmgmtServicePb.AddUserToLocalGroupRequest) (*netmgmtServicePb.AddUserToLocalGroupResponse, error) {
	if err := netmgmtApi.AddUserToLocalGroup(in.GetUser(), in.GetLocalGroup()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.AddUserToLocalGroupResponse{}, nil
}

func (s *Server) RemoveUserFromLocalGroup(ctx context.Context, in *netmgmtServicePb.RemoveUserFromLocalGroupRequest) (*netmgmtServicePb.RemoveUserFromLocalGroupResponse, error) {
	if err := netmgmtApi.RemoveUserFromLocalGroup(in.GetUser(), in.GetLocalGroup()); err != nil {
		return nil, err
	}

	return &netmgmtServicePb.RemoveUserFromLocalGroupResponse{}, nil
}

func (s *Server) GetLocalGroups(ctx context.Context, in *netmgmtServicePb.GetLocalGroupsRequest) (*netmgmtServicePb.GetLocalGroupsResponse, error) {
	localGroups, err := netmgmtApi.GetLocalGroups()
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetLocalGroupsResponse{
		LocalGroups: localGroups,
	}, nil
}

func (s *Server) GetUsersInLocalGroup(ctx context.Context, in *netmgmtServicePb.GetUsersInLocalGroupRequest) (*netmgmtServicePb.GetUsersInLocalGroupResponse, error) {
	users, err := netmgmtApi.GetUsersInLocalGroup(in.GetLocalGroup())
	if err != nil {
		return nil, err
	}

	return &netmgmtServicePb.GetUsersInLocalGroupResponse{
		Users: users,
	}, nil
}
