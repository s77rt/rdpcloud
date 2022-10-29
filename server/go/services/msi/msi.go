//go:build windows && amd64

package msi

import (
	msiServicePb "github.com/s77rt/rdpcloud/proto/go/services/msi"
)

type Server struct {
	msiServicePb.UnimplementedMsiServer
}
