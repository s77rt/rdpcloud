//go:build windows && amd64

package fileio

import (
	fileioServicePb "github.com/s77rt/rdpcloud/proto/go/services/fileio"
)

type Server struct {
	fileioServicePb.UnimplementedFileioServer
}
