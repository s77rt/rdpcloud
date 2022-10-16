//go:build windows && amd64

package fileio

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fileioModelsPb "github.com/s77rt/rdpcloud/proto/go/models/fileio"
	fileioServicePb "github.com/s77rt/rdpcloud/proto/go/services/fileio"
	fileioApi "github.com/s77rt/rdpcloud/server/go/api/fileio"
	secauthzApi "github.com/s77rt/rdpcloud/server/go/api/secauthz"
	"github.com/s77rt/rdpcloud/server/go/auth"
)

func (s *Server) GetQuotaState(ctx context.Context, in *fileioServicePb.GetQuotaStateRequest) (*fileioServicePb.GetQuotaStateResponse, error) {
	quotaState, err := fileioApi.GetQuotaState(in.GetVolumePath())
	if err != nil {
		return nil, err
	}

	return &fileioServicePb.GetQuotaStateResponse{
		QuotaState: quotaState,
	}, nil
}

func (s *Server) SetQuotaState(ctx context.Context, in *fileioServicePb.SetQuotaStateRequest) (*fileioServicePb.SetQuotaStateResponse, error) {
	if err := fileioApi.SetQuotaState(in.GetVolumePath(), in.GetQuotaState()); err != nil {
		return nil, err
	}

	return &fileioServicePb.SetQuotaStateResponse{}, nil
}

func (s *Server) GetDefaultQuota(ctx context.Context, in *fileioServicePb.GetDefaultQuotaRequest) (*fileioServicePb.GetDefaultQuotaResponse, error) {
	defaultQuota, err := fileioApi.GetDefaultQuota(in.GetVolumePath())
	if err != nil {
		return nil, err
	}

	return &fileioServicePb.GetDefaultQuotaResponse{
		DefaultQuota: defaultQuota,
	}, nil
}

func (s *Server) SetDefaultQuota(ctx context.Context, in *fileioServicePb.SetDefaultQuotaRequest) (*fileioServicePb.SetDefaultQuotaResponse, error) {
	if err := fileioApi.SetDefaultQuota(in.GetVolumePath(), in.GetDefaultQuota()); err != nil {
		return nil, err
	}

	return &fileioServicePb.SetDefaultQuotaResponse{}, nil
}

func (s *Server) GetUsersQuotaEntries(ctx context.Context, in *fileioServicePb.GetUsersQuotaEntriesRequest) (*fileioServicePb.GetUsersQuotaEntriesResponse, error) {
	quotaEntries, err := fileioApi.GetUsersQuotaEntries(in.GetVolumePath())
	if err != nil {
		return nil, err
	}

	return &fileioServicePb.GetUsersQuotaEntriesResponse{
		QuotaEntries: quotaEntries,
	}, nil
}

func (s *Server) GetUserQuotaEntry(ctx context.Context, in *fileioServicePb.GetUserQuotaEntryRequest) (*fileioServicePb.GetUserQuotaEntryResponse, error) {
	quotaEntry, err := fileioApi.GetUserQuotaEntry(in.GetVolumePath(), in.GetUser())
	if err != nil {
		return nil, err
	}

	return &fileioServicePb.GetUserQuotaEntryResponse{
		QuotaEntry: quotaEntry,
	}, nil
}

func (s *Server) GetMyUserQuotaEntry(ctx context.Context, in *fileioServicePb.GetMyUserQuotaEntryRequest) (*fileioServicePb.GetMyUserQuotaEntryResponse, error) {
	userClaims, err := auth.UserClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid context user claims (%s)", err.Error())
	}

	username, err := secauthzApi.LookupAccountUsernameBySid(userClaims.UserSID)
	if err != nil {
		return nil, err
	}

	quotaEntry, err := fileioApi.GetUserQuotaEntry(in.GetVolumePath(), &fileioModelsPb.User_1{Username: username})
	if err != nil {
		return nil, err
	}

	return &fileioServicePb.GetMyUserQuotaEntryResponse{
		QuotaEntry: quotaEntry,
	}, nil
}

func (s *Server) SetUserQuotaEntry(ctx context.Context, in *fileioServicePb.SetUserQuotaEntryRequest) (*fileioServicePb.SetUserQuotaEntryResponse, error) {
	if err := fileioApi.SetUserQuotaEntry(in.GetVolumePath(), in.GetUser(), in.GetQuotaEntry()); err != nil {
		return nil, err
	}

	return &fileioServicePb.SetUserQuotaEntryResponse{}, nil
}

func (s *Server) DeleteUserQuotaEntry(ctx context.Context, in *fileioServicePb.DeleteUserQuotaEntryRequest) (*fileioServicePb.DeleteUserQuotaEntryResponse, error) {
	if err := fileioApi.DeleteUserQuotaEntry(in.GetVolumePath(), in.GetUser()); err != nil {
		return nil, err
	}

	return &fileioServicePb.DeleteUserQuotaEntryResponse{}, nil
}
