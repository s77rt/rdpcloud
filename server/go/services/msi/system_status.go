//go:build windows && amd64

package msi

import (
	"context"

	msiServicePb "github.com/s77rt/rdpcloud/proto/go/services/msi"
	msiApi "github.com/s77rt/rdpcloud/server/go/api/msi"
)

func (s *Server) GetProducts(ctx context.Context, in *msiServicePb.GetProductsRequest) (*msiServicePb.GetProductsResponse, error) {
	products, err := msiApi.GetProducts()
	if err != nil {
		return nil, err
	}

	return &msiServicePb.GetProductsResponse{
		Products: products,
	}, nil
}
