//go:build windows && amd64

package sysinfo

import (
	sysinfoServicePb "github.com/s77rt/rdpcloud/proto/go/services/sysinfo"
)

type Server struct {
	sysinfoServicePb.UnimplementedSysinfoServer
}
